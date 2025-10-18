package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/idanyas/go-pterodactyl"
	"github.com/idanyas/go-pterodactyl/client"
	"github.com/idanyas/go-pterodactyl/helpers"
)

func main() {
	panelURL := os.Getenv("PTERODACTYL_URL")
	apiKey := os.Getenv("PTERODACTYL_CLIENT_API_KEY")
	serverID := os.Getenv("PTERODACTYL_SERVER_ID")

	if panelURL == "" || apiKey == "" || serverID == "" {
		log.Fatal("Please set PTERODACTYL_URL, PTERODACTYL_CLIENT_API_KEY, and PTERODACTYL_SERVER_ID")
	}

	pClient, err := pterodactyl.New(panelURL, pterodactyl.WithAPIKey(apiKey))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	clientAPI := pClient.Client()

	// Example 1: Wait for server to reach running state
	fmt.Println("=== Waiting for Server to Start ===")

	// Start the server
	err = clientAPI.SendPowerAction(ctx, serverID, "start")
	if err != nil {
		log.Printf("Failed to start server: %v", err)
	}

	// Wait for it to reach running state
	waiter := helpers.NewStateWaiter(clientAPI)
	waitCtx, waitCancel := context.WithTimeout(ctx, 2*time.Minute)
	defer waitCancel()

	err = waiter.WaitForState(waitCtx, serverID, "running", 5*time.Second)
	if err != nil {
		log.Printf("Failed to wait for server: %v", err)
	} else {
		fmt.Println("✓ Server is now running!")
	}

	// Example 2: Create backup and wait for completion
	fmt.Println("\n=== Creating Backup ===")

	manager := helpers.NewBackupManager(clientAPI)
	backup, err := manager.CreateAndWait(ctx, serverID, client.CreateBackupRequest{
		Name: fmt.Sprintf("Auto Backup - %s", time.Now().Format("2006-01-02 15:04")),
	}, 10*time.Second)

	if err != nil {
		log.Printf("Failed to create backup: %v", err)
	} else {
		fmt.Printf("✓ Backup completed: %s (%d MB)\n", backup.Name, backup.Bytes/1024/1024)
	}

	// Example 3: Download a file
	fmt.Println("\n=== Downloading File ===")

	downloader := helpers.NewFileDownloader(clientAPI)
	outputFile, err := os.Create("downloaded-server.properties")
	if err != nil {
		log.Printf("Failed to create output file: %v", err)
	} else {
		defer outputFile.Close()

		err = downloader.DownloadToWriter(ctx, serverID, "/server.properties", outputFile)
		if err != nil {
			log.Printf("Failed to download file: %v", err)
		} else {
			fmt.Println("✓ File downloaded successfully!")
		}
	}
}
