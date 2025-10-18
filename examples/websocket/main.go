package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/idanyas/go-pterodactyl"
	"github.com/idanyas/go-pterodactyl/websocket"
)

func main() {
	// Get configuration from environment variables
	panelURL := os.Getenv("PTERODACTYL_URL")
	apiKey := os.Getenv("PTERODACTYL_CLIENT_API_KEY")
	serverID := os.Getenv("PTERODACTYL_SERVER_ID")

	if panelURL == "" || apiKey == "" || serverID == "" {
		log.Fatal("Please set PTERODACTYL_URL, PTERODACTYL_CLIENT_API_KEY, and PTERODACTYL_SERVER_ID environment variables")
	}

	// Create a new Pterodactyl client
	client, err := pterodactyl.New(panelURL, pterodactyl.WithAPIKey(apiKey))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Set up context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Set up signal handling for graceful shutdown
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		<-signalChan
		fmt.Println("\nReceived interrupt signal, shutting down...")
		cancel()
	}()

	// Get the Client API
	clientAPI := client.Client()

	// Connect to the server's WebSocket
	fmt.Printf("Connecting to WebSocket for server %s...\n", serverID)
	ws, err := clientAPI.ConnectWebSocket(ctx, serverID)
	if err != nil {
		log.Fatalf("Failed to connect to WebSocket: %v", err)
	}
	defer ws.Close()

	fmt.Println("WebSocket connection established!")
	fmt.Println("Listening for events... (Press Ctrl+C to exit)")
	fmt.Println()

	// Track stats for periodic display
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	var lastStats *websocket.StatsEvent

	// Listen for events
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Context cancelled, closing connection...")
			return

		case <-ticker.C:
			// Periodically display stats summary
			if lastStats != nil {
				fmt.Printf("[Summary] CPU: %.2f%%, Memory: %d MB / %d MB, Uptime: %s\n",
					lastStats.Stats.CPUAbsolute,
					lastStats.Stats.MemoryBytes/(1024*1024),
					lastStats.Stats.MemoryLimitBytes/(1024*1024),
					formatUptime(lastStats.Stats.Uptime))
			}

		case event, ok := <-ws.Events():
			if !ok {
				fmt.Println("Event channel closed, connection terminated.")
				return
			}

			// Handle different event types
			switch e := event.(type) {
			case *websocket.ConsoleOutputEvent:
				// Print console output
				fmt.Printf("[Console] %s\n", e.Line)

			case *websocket.StatsEvent:
				// Store latest stats (displayed periodically)
				lastStats = e

			case *websocket.StatusEvent:
				// Print status changes
				fmt.Printf("[Status] Server state changed to: %s\n", e.Status)

			case *websocket.TokenExpiredEvent:
				fmt.Println("[Error] WebSocket token expired. Reconnection needed.")
				return

			default:
				fmt.Printf("[Unknown Event] Type: %T\n", e)
			}
		}
	}
}

// formatUptime converts milliseconds to a human-readable format
func formatUptime(uptimeMs int64) string {
	d := time.Duration(uptimeMs) * time.Millisecond
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	if hours > 0 {
		return fmt.Sprintf("%dh %dm %ds", hours, minutes, seconds)
	} else if minutes > 0 {
		return fmt.Sprintf("%dm %ds", minutes, seconds)
	}
	return fmt.Sprintf("%ds", seconds)
}
