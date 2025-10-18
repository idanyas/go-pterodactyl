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
	apiKey := os.Getenv("PTERODACTYL_APPLICATION_API_KEY")

	if panelURL == "" || apiKey == "" {
		log.Fatal("Please set PTERODACTYL_URL and PTERODACTYL_APPLICATION_API_KEY environment variables")
	}

	// Create a new Pterodactyl client
	client, err := pterodactyl.New(panelURL, pterodactyl.WithAPIKey(apiKey))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get the Application API
	appAPI := client.Application()

	// Example 1: List all users
	fmt.Println("=== Panel Users ===")
	users, userPaginator, err := appAPI.ListUsers(ctx, pagination.ListOptions{PerPage: 10})
	if err != nil {
		log.Fatalf("Failed to list users: %v", err)
	}
	fmt.Printf("Found %d user(s):\n", len(users))
	for _, user := range users {
		fmt.Printf("  - %s (%s) - Admin: %v\n", user.Username, user.Email, user.RootAdmin)
	}
	fmt.Println()

	if userPaginator.HasMorePages() {
		fmt.Println("There are more users. Use paginator.NextPage() to fetch them.")
	}

	// Example 2: List all locations
	fmt.Println("=== Locations ===")
	locations, _, err := appAPI.ListLocations(ctx, pagination.ListOptions{})
	if err != nil {
		log.Fatalf("Failed to list locations: %v", err)
	}
	for _, location := range locations {
		fmt.Printf("  - [%s] %s (ID: %d)\n", location.Short, location.Long, location.ID)
	}
	fmt.Println()

	// Example 3: List all nodes
	fmt.Println("=== Nodes ===")
	nodes, _, err := appAPI.ListNodes(ctx, pagination.ListOptions{})
	if err != nil {
		log.Fatalf("Failed to list nodes: %v", err)
	}
	for _, node := range nodes {
		fmt.Printf("  - %s (%s)\n", node.Name, node.FQDN)
		fmt.Printf("    Memory: %d MB (overallocate: %d%%)\n", node.Memory, node.MemoryOverallocate)
		fmt.Printf("    Disk: %d MB (overallocate: %d%%)\n", node.Disk, node.DiskOverallocate)
		fmt.Printf("    Location ID: %d\n", node.LocationID)
	}
	fmt.Println()

	// Example 4: List all nests and their eggs
	fmt.Println("=== Nests and Eggs ===")
	nests, _, err := appAPI.ListNests(ctx, pagination.ListOptions{})
	if err != nil {
		log.Fatalf("Failed to list nests: %v", err)
	}
	for _, nest := range nests {
		fmt.Printf("  - %s (ID: %d): %s\n", nest.Name, nest.ID, nest.Description)

		eggs, _, err := appAPI.ListNestEggs(ctx, nest.ID, pagination.ListOptions{})
		if err != nil {
			log.Printf("    Failed to list eggs: %v", err)
			continue
		}
		for _, egg := range eggs {
			fmt.Printf("    * %s (ID: %d)\n", egg.Name, egg.ID)
		}
	}
	fmt.Println()

	// Example 5: List all servers
	fmt.Println("=== Servers ===")
	servers, _, err := appAPI.ListServers(ctx, pagination.ListOptions{
		PerPage: 10,
	})
	if err != nil {
		log.Fatalf("Failed to list servers: %v", err)
	}
	fmt.Printf("Found %d server(s):\n", len(servers))
	for _, server := range servers {
		fmt.Printf("  - %s (ID: %d, UUID: %s)\n", server.Name, server.ID, server.UUID)
		fmt.Printf("    Owner User ID: %d\n", server.UserID)
		fmt.Printf("    Node ID: %d\n", server.NodeID)
		fmt.Printf("    Memory: %d MB, Disk: %d MB, CPU: %d%%\n",
			server.Limits.Memory, server.Limits.Disk, server.Limits.CPU)
		fmt.Printf("    Suspended: %v, Installing: %v\n", server.Suspended, server.Installing)
	}
	fmt.Println()

	// Example 6: Create a new location (commented out to avoid side effects)
	/*
		fmt.Println("=== Creating a New Location ===")
		newLocation, err := appAPI.CreateLocation(ctx, application.CreateLocationRequest{
			Short: "us-west",
			Long:  "United States - West Coast",
		})
		if err != nil {
			log.Printf("Failed to create location: %v", err)
		} else {
			fmt.Printf("Created location: %s (ID: %d)\n", newLocation.Short, newLocation.ID)
		}
	*/

	// Example 7: Get deployable nodes (if needed for server creation)
	if len(nodes) > 0 {
		fmt.Println("=== Deployable Nodes ===")
		deployableNodes, err := appAPI.GetDeployableNodes(ctx, 1024, 5120)
		if err != nil {
			log.Printf("Failed to get deployable nodes: %v", err)
		} else {
			fmt.Printf("Nodes that can accept a server with 1GB RAM and 5GB disk:\n")
			for _, node := range deployableNodes {
				fmt.Printf("  - %s (ID: %d)\n", node.Name, node.ID)
			}
		}
		fmt.Println()
	}
}
