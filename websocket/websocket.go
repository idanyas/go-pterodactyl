// Package websocket provides a real-time client for Pterodactyl servers.
package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/coder/websocket"
	"github.com/idanyas/go-pterodactyl/models"
)

// message represents the JSON structure of a WebSocket message.
type message struct {
	Event string   `json:"event"`
	Args  []string `json:"args,omitempty"`
}

// Event is an interface for events received from the WebSocket.
type Event interface {
	isEvent()
}

// ConsoleOutputEvent represents a line of console output.
type ConsoleOutputEvent struct {
	Line string
}

func (e *ConsoleOutputEvent) isEvent() {}

// StatsEvent represents a server resource usage update.
type StatsEvent struct {
	Stats models.Resources
}

func (e *StatsEvent) isEvent() {}

// StatusEvent represents a change in the server's power state.
type StatusEvent struct {
	Status string
}

func (e *StatusEvent) isEvent() {}

// TokenExpiredEvent indicates the WebSocket JWT has expired.
type TokenExpiredEvent struct{}

func (e *TokenExpiredEvent) isEvent() {}

// ReconnectOptions configures automatic reconnection behavior.
type ReconnectOptions struct {
	// Enable enables automatic reconnection.
	Enable bool
	// MaxAttempts is the maximum number of reconnection attempts (0 = unlimited).
	MaxAttempts int
	// InitialDelay is the initial delay before the first reconnection attempt.
	InitialDelay time.Duration
	// MaxDelay is the maximum delay between reconnection attempts.
	MaxDelay time.Duration
	// Multiplier is the exponential backoff multiplier.
	Multiplier float64
}

// DefaultReconnectOptions returns sensible defaults for reconnection.
func DefaultReconnectOptions() ReconnectOptions {
	return ReconnectOptions{
		Enable:       true,
		MaxAttempts:  10,
		InitialDelay: 1 * time.Second,
		MaxDelay:     30 * time.Second,
		Multiplier:   2.0,
	}
}

// Conn represents an active WebSocket connection to a server.
type Conn struct {
	socketURL string
	token     string
	conn      *websocket.Conn
	ctx       context.Context
	cancel    context.CancelFunc
	eventChan chan Event
	closeOnce sync.Once

	// Reconnection
	reconnectOpts ReconnectOptions
	reconnecting  bool
	mu            sync.RWMutex
}

// NewConn establishes a new WebSocket connection with optional reconnection.
// Pass nil for reconnectOpts to disable automatic reconnection.
func NewConn(ctx context.Context, socketURL, token string, reconnectOpts *ReconnectOptions) (*Conn, error) {
	wsConnCtx, cancel := context.WithCancel(context.Background())

	ws := &Conn{
		socketURL: socketURL,
		token:     token,
		ctx:       wsConnCtx,
		cancel:    cancel,
		eventChan: make(chan Event, 100),
	}

	if reconnectOpts != nil {
		ws.reconnectOpts = *reconnectOpts
	}

	if err := ws.connect(ctx); err != nil {
		cancel()
		return nil, err
	}

	go ws.readLoop()

	return ws, nil
}

// connect establishes the WebSocket connection and authenticates.
func (ws *Conn) connect(ctx context.Context) error {
	conn, _, err := websocket.Dial(ctx, ws.socketURL, nil)
	if err != nil {
		return fmt.Errorf("failed to dial websocket: %w", err)
	}
	conn.SetReadLimit(10 * 1024 * 1024) // 10MB limit

	// Authenticate
	authMsg := message{
		Event: "auth",
		Args:  []string{ws.token},
	}
	authBytes, _ := json.Marshal(authMsg)
	if err := conn.Write(ctx, websocket.MessageText, authBytes); err != nil {
		conn.Close(websocket.StatusAbnormalClosure, "failed to send auth")
		return fmt.Errorf("failed to send auth message: %w", err)
	}

	ws.mu.Lock()
	ws.conn = conn
	ws.mu.Unlock()

	return nil
}

// readLoop continuously reads messages from the WebSocket and dispatches them as events.
func (ws *Conn) readLoop() {
	defer close(ws.eventChan)

	attempt := 0
	backoff := ws.reconnectOpts.InitialDelay

	for {
		ws.mu.RLock()
		conn := ws.conn
		ws.mu.RUnlock()

		if conn == nil {
			return
		}

		_, data, err := conn.Read(ws.ctx)
		if err != nil {
			// Check if we should reconnect
			if ws.reconnectOpts.Enable && !ws.reconnecting {
				ws.mu.Lock()
				ws.reconnecting = true
				ws.mu.Unlock()

				// Close old connection
				conn.Close(websocket.StatusAbnormalClosure, "read error")

				// Attempt reconnection
				for {
					if ws.reconnectOpts.MaxAttempts > 0 && attempt >= ws.reconnectOpts.MaxAttempts {
						// Max attempts reached
						ws.mu.Lock()
						ws.reconnecting = false
						ws.mu.Unlock()
						return
					}

					select {
					case <-ws.ctx.Done():
						return
					case <-time.After(backoff):
						attempt++

						// Try to reconnect
						ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
						err := ws.connect(ctx)
						cancel()

						if err == nil {
							// Reconnected successfully
							ws.mu.Lock()
							ws.reconnecting = false
							ws.mu.Unlock()
							attempt = 0
							backoff = ws.reconnectOpts.InitialDelay
							break
						}

						// Exponential backoff with jitter
						backoff = time.Duration(float64(backoff) * ws.reconnectOpts.Multiplier)
						if backoff > ws.reconnectOpts.MaxDelay {
							backoff = ws.reconnectOpts.MaxDelay
						}
						jitter := time.Duration(rand.Float64() * float64(backoff) * 0.1)
						backoff += jitter
					}

					if !ws.reconnecting {
						break
					}
				}

				continue
			}

			// No reconnection, just exit
			return
		}

		// Reset backoff on successful read
		backoff = ws.reconnectOpts.InitialDelay
		attempt = 0

		var msg message
		if err := json.Unmarshal(data, &msg); err != nil {
			continue // Ignore malformed messages
		}

		var event Event
		switch msg.Event {
		case "console output":
			if len(msg.Args) > 0 {
				event = &ConsoleOutputEvent{Line: msg.Args[0]}
			}
		case "stats":
			if len(msg.Args) > 0 {
				var stats models.Resources
				if json.Unmarshal([]byte(msg.Args[0]), &stats) == nil {
					event = &StatsEvent{Stats: stats}
				}
			}
		case "status":
			if len(msg.Args) > 0 {
				event = &StatusEvent{Status: msg.Args[0]}
			}
		case "jwt error", "token expiring", "token expired":
			event = &TokenExpiredEvent{}
		}

		if event != nil {
			select {
			case ws.eventChan <- event:
			case <-ws.ctx.Done():
				return
			}
		}
	}
}

// Events returns a read-only channel for receiving WebSocket events.
func (ws *Conn) Events() <-chan Event {
	return ws.eventChan
}

// SendCommand sends a command to the server's console.
func (ws *Conn) SendCommand(command string) error {
	return ws.sendEvent("send command", command)
}

// SetState changes the power state of the server.
// Valid states are "start", "stop", "restart", "kill".
func (ws *Conn) SetState(state string) error {
	return ws.sendEvent("set state", state)
}

// sendEvent is a helper to marshal and send a message to the WebSocket.
func (ws *Conn) sendEvent(event, arg string) error {
	ws.mu.RLock()
	conn := ws.conn
	ws.mu.RUnlock()

	if conn == nil {
		return fmt.Errorf("websocket connection is closed")
	}

	msg := message{Event: event, Args: []string{arg}}
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ws.ctx, 10*time.Second)
	defer cancel()

	return conn.Write(ctx, websocket.MessageText, data)
}

// Close gracefully closes the WebSocket connection.
func (ws *Conn) Close() {
	ws.closeOnce.Do(func() {
		ws.cancel()

		ws.mu.RLock()
		conn := ws.conn
		ws.mu.RUnlock()

		if conn != nil {
			conn.Close(websocket.StatusNormalClosure, "")
		}
	})
}

// IsReconnecting returns true if the connection is currently attempting to reconnect.
func (ws *Conn) IsReconnecting() bool {
	ws.mu.RLock()
	defer ws.mu.RUnlock()
	return ws.reconnecting
}
