package websocket

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/coder/websocket"
)

func TestWebSocket_ConnectAndReceive(t *testing.T) {
	// Mock WebSocket Server
	var wg sync.WaitGroup
	wg.Add(1)
	wsServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := websocket.Accept(w, r, nil)
		if err != nil {
			t.Logf("websocket accept error: %v", err)
			return
		}
		defer c.Close(websocket.StatusInternalError, "internal error")

		// Expect auth message
		_, authData, err := c.Read(r.Context())
		if err != nil {
			t.Logf("ws read auth error: %v", err)
			return
		}
		var authMsg message
		json.Unmarshal(authData, &authMsg)
		if authMsg.Event != "auth" || authMsg.Args[0] != "test-token" {
			t.Errorf("expected auth event with test-token, got %+v", authMsg)
			return
		}

		// Send some events
		events := []message{
			{Event: "status", Args: []string{"running"}},
			{Event: "console output", Args: []string{"Server started!"}},
		}
		for _, event := range events {
			data, _ := json.Marshal(event)
			c.Write(r.Context(), websocket.MessageText, data)
		}

		// Wait for client to close connection
		wg.Wait()
		c.Close(websocket.StatusNormalClosure, "")
	}))
	defer wsServer.Close()

	wsURL := "ws" + strings.TrimPrefix(wsServer.URL, "http")

	// Test logic
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ws, err := NewConn(ctx, wsURL, "test-token", nil)
	if err != nil {
		t.Fatalf("NewConn failed: %v", err)
	}

	var receivedStatus, receivedConsole bool
	for i := 0; i < 2; i++ {
		select {
		case event := <-ws.Events():
			switch e := event.(type) {
			case *StatusEvent:
				if e.Status == "running" {
					receivedStatus = true
				}
			case *ConsoleOutputEvent:
				if e.Line == "Server started!" {
					receivedConsole = true
				}
			default:
				t.Errorf("unexpected event type: %T", e)
			}
		case <-ctx.Done():
			t.Fatal("test timed out")
		}
	}

	ws.Close()
	wg.Done() // Signal WS server it can close

	if !receivedStatus || !receivedConsole {
		t.Errorf("did not receive all expected events (status: %v, console: %v)", receivedStatus, receivedConsole)
	}
}
