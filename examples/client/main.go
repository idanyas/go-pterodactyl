package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/idanyas/go-pterodactyl"
	"github.com/idanyas/go-pterodactyl/pagination"
)

func main() {
	// Get configuration from environment variables
	panelURL := os.Getenv("PTERODACTYL_URL")
	apiKey := os.Getenv("PTERODACTYL_CLIENT_API_KEY")

	if panelURL == "" || apiKey == "" {
		log.Fatal("Please set PTERODACTYL_URL and PTERODACTYL_CLIENT_API_KEY environment variables")
	}

	// Create a new Pterodactyl client
	client, err := pterodactyl.New(panelURL, pterodactyl.WithAPIKey(apiKey))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get the Client API
	clientAPI := client.Client()

	// Example 1: Get account information
	fmt.Println("=== Account Information ===")
	account, err := clientAPI.GetAccount(ctx)
	if err != nil {
		log.Fatalf("Failed to get account: %v", err)
	}
	fmt.Printf("Username: %s\n", account.Username)
	fmt.Printf("Email: %s\n", account.Email)
	fmt.Printf("Admin: %v\n", account.RootAdmin)
	fmt.Println()

	// Example 2: List all servers
	fmt.Println("=== Your Servers ===")
	servers, paginator, err := clientAPI.ListServers(ctx, pagination.ListOptions{
		PerPage: 10,
		Include: []string{"allocations"},
	})
	if err != nil {
		log.Fatalf("Failed to list servers: %v", err)
	}

	for _, server := range servers {
		fmt.Printf("Server: %s (ID: %s)\n", server.Name, server.Identifier)
		fmt.Printf("  Status: %s\n", safeString(server.Status))
		fmt.Printf("  Memory: %d MB\n", server.Limits.Memory)
		fmt.Printf("  Disk: %d MB\n", server.Limits.Disk)
		if server.Relationships != nil && len(server.Relationships.Allocations.Data) > 0 {
			alloc := server.Relationships.Allocations.Data[0]
			fmt.Printf("  Primary IP: %s:%d\n", alloc.IP, alloc.Port)
		}
		fmt.Println()
	}

	// Check if there are more pages
	if paginator.HasMorePages() {
		fmt.Println("There are more servers available. Use paginator.NextPage() to fetch them.")
	}

	// Example 3: Get system permissions
	fmt.Println("=== System Permissions ===")
	perms, err := clientAPI.GetSystemPermissions(ctx)
	if err != nil {
		log.Fatalf("Failed to get permissions: %v", err)
	}
	for groupName, group := range perms.Permissions {
		fmt.Printf("%s: %s\n", groupName, group.Description)
		for key, desc := range group.Keys {
			fmt.Printf("  - %s: %s\n", key, desc)
		}
	}
	fmt.Println()

	// Example 4: Work with a specific server (if available)
	if len(servers) > 0 {
		serverID := servers[0].Identifier
		fmt.Printf("=== Working with Server: %s ===\n", serverID)

		// Get server resources
		resources, err := clientAPI.GetServerResources(ctx, serverID)
		if err != nil {
			log.Printf("Failed to get server resources: %v", err)
		} else {
			fmt.Printf("Current State: %s\n", resources.CurrentState)
			fmt.Printf("Memory Usage: %d / %d MB\n",
				resources.Resources.MemoryBytes/(1024*1024),
				resources.Resources.MemoryLimitBytes/(1024*1024))
			fmt.Printf("CPU Usage: %.2f%%\n", resources.Resources.CPUAbsolute)
		}
		fmt.Println()

		// List server files
		files, err := clientAPI.ListFiles(ctx, serverID, "/")
		if err != nil {
			log.Printf("Failed to list files: %v", err)
		} else {
			fmt.Println("Root directory files:")
			for _, file := range files {
				fileType := "file"
				if !file.IsFile {
					fileType = "dir"
				}
				fmt.Printf("  [%s] %s (%d bytes)\n", fileType, file.Name, file.Size)
			}
		}
		fmt.Println()

		// List backups
		backups, err := clientAPI.ListBackups(ctx, serverID)
		if err != nil {
			log.Printf("Failed to list backups: %v", err)
		} else {
			fmt.Printf("Server has %d backup(s)\n", len(backups))
			for _, backup := range backups {
				fmt.Printf("  - %s (%d MB) - Created: %s\n",
					backup.Name,
					backup.Bytes/(1024*1024),
					backup.CreatedAt.Format(time.RFC3339))
			}
		}
	}
}

func safeString(s *string) string {
	if s == nil {
		return "unknown"
	}
	return *s
}
