package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/idanyas/go-pterodactyl/websocket"
)

// ConnectWebSocket establishes a WebSocket connection to a server.
// For automatic reconnection, pass reconnectOpts. Pass nil to disable reconnection.
func (c *client) ConnectWebSocket(ctx context.Context, serverID string) (*websocket.Conn, error) {
	return c.ConnectWebSocketWithReconnect(ctx, serverID, nil)
}

// ConnectWebSocketWithReconnect establishes a WebSocket connection with custom reconnection options.
func (c *client) ConnectWebSocketWithReconnect(ctx context.Context, serverID string, reconnectOpts *websocket.ReconnectOptions) (*websocket.Conn, error) {
	if serverID == "" {
		return nil, fmt.Errorf("server ID cannot be empty")
	}

	path := fmt.Sprintf("client/servers/%s/websocket", serverID)
	var response struct {
		Data struct {
			Token  string `json:"token"`
			Socket string `json:"socket"`
		} `json:"data"`
	}
	_, err := c.client.Do(ctx, http.MethodGet, path, nil, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to get websocket credentials: %w", err)
	}

	return websocket.NewConn(ctx, response.Data.Socket, response.Data.Token, reconnectOpts)
}
