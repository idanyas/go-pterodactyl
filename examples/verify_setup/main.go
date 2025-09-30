package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/idanyas/go-pterodactyl"
)

func main() {
	baseURL := os.Getenv("PTERO_BASE_URL")
	apiKey := os.Getenv("PTERO_API_KEY")
	if baseURL == "" || apiKey == "" {
		log.Fatal("Set PTERO_BASE_URL and PTERO_API_KEY environment variables")
	}

	client, err := pterodactyl.NewClient(baseURL, apiKey, pterodactyl.ApplicationKey)
	if err != nil {
		log.Fatalf("new client: %v", err)
	}

	fmt.Println("--- Verifying Panel Setup ---")

	// 1. List all nodes
	nodes, err := client.ApplicationAPI.Nodes.ListAll(context.Background())
	if err != nil {
		log.Fatalf("Fatal: Could not list nodes: %v", err)
	}

	if len(nodes) == 0 {
		log.Fatal("Fatal: No nodes are configured on this panel.")
	}

	fmt.Printf("\nFound %d Node(s):\n", len(nodes))

	// 2. For each node, list its allocations
	for _, node := range nodes {
		fmt.Printf("\n====================================\n")
		fmt.Printf("Node ID: %d | Name: %s\n", node.ID, node.Name)
		fmt.Printf("====================================\n")

		allocations, err := client.ApplicationAPI.Nodes.Allocations(context.Background(), node.ID).ListAll(context.Background())
		if err != nil {
			log.Printf("Warning: Could not list allocations for node %d: %v", node.ID, err)
			continue
		}

		if len(allocations) == 0 {
			fmt.Println(" -> This node has NO allocations configured.")
			continue
		}

		fmt.Println("  ID   | IP:Port              | Assigned?")
		fmt.Println("--------------------------------------------")
		for _, alloc := range allocations {
			// This is the correct way to format the output
			addr := alloc.IP + ":" + strconv.Itoa(alloc.Port)
			fmt.Printf("  %-4d | %-20s | %t\n", alloc.ID, addr, alloc.Assigned)
		}
	}
}
