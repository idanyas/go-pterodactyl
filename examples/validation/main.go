package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/idanyas/go-pterodactyl"
	"github.com/idanyas/go-pterodactyl/application"
	"github.com/idanyas/go-pterodactyl/validation"
)

func main() {
	panelURL := os.Getenv("PTERODACTYL_URL")
	apiKey := os.Getenv("PTERODACTYL_APPLICATION_API_KEY")

	if panelURL == "" || apiKey == "" {
		log.Fatal("Please set PTERODACTYL_URL and PTERODACTYL_APPLICATION_API_KEY")
	}

	client, err := pterodactyl.New(panelURL, pterodactyl.WithAPIKey(apiKey))
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	app := client.Application()

	// Example 1: Valid request
	fmt.Println("=== Creating user with valid data ===")
	validReq := application.CreateUserRequest{
		Email:     "newuser@example.com",
		Username:  "newuser",
		FirstName: "New",
		LastName:  "User",
	}

	if err := validation.Validate(validReq); err != nil {
		log.Printf("Validation failed: %v\n", err)
	} else {
		fmt.Println("✓ Validation passed")
		// Would create user here in real scenario
		// user, err := app.CreateUser(ctx, validReq)
	}

	// Example 2: Invalid request (missing required fields)
	fmt.Println("\n=== Creating user with missing email ===")
	invalidReq1 := application.CreateUserRequest{
		Username:  "newuser2",
		FirstName: "New",
		LastName:  "User",
		// Email is missing
	}

	if err := validation.Validate(invalidReq1); err != nil {
		fmt.Printf("✓ Validation correctly failed: %v\n", err)
	} else {
		fmt.Println("✗ Validation should have failed but passed")
	}

	// Example 3: Invalid request (invalid email)
	fmt.Println("\n=== Creating user with invalid email ===")
	invalidReq2 := application.CreateUserRequest{
		Email:     "not-an-email",
		Username:  "newuser3",
		FirstName: "New",
		LastName:  "User",
	}

	if err := validation.Validate(invalidReq2); err != nil {
		fmt.Printf("✓ Validation correctly failed: %v\n", err)
	} else {
		fmt.Println("✗ Validation should have failed but passed")
	}

	// Example 4: Validate location
	fmt.Println("\n=== Creating location with valid data ===")
	locationReq := application.CreateLocationRequest{
		Short: "us-west",
		Long:  "United States - West Coast",
	}

	if err := validation.Validate(locationReq); err != nil {
		log.Printf("Validation failed: %v\n", err)
	} else {
		fmt.Println("✓ Validation passed")
		location, err := app.CreateLocation(ctx, locationReq)
		if err != nil {
			log.Printf("Failed to create location: %v\n", err)
		} else {
			fmt.Printf("Created location: %s (ID: %d)\n", location.Short, location.ID)
			// Clean up created location
			app.DeleteLocation(ctx, location.ID)
			fmt.Printf("Cleaned up location ID: %d\n", location.ID)
		}
	}

	// Example 5: Test server creation validation
	fmt.Println("\n=== Testing server creation request validation ===")
	serverReq := application.CreateServerRequest{
		// Missing required fields
		Name: "Test Server",
		User: 1,
		// Egg is missing
	}

	if err := validation.Validate(serverReq); err != nil {
		fmt.Printf("✓ Validation correctly caught missing fields: %v\n", err)
	} else {
		fmt.Println("✗ Validation should have failed")
	}

	// Example 6: Multiple validation errors
	fmt.Println("\n=== Testing multiple validation errors ===")
	multiErrorReq := application.CreateUserRequest{
		Email:     "invalid", // Invalid email
		Username:  "",        // Empty username
		FirstName: "A",       // Too short (if we had min validation)
		LastName:  "",        // Missing
	}

	if err := validation.Validate(multiErrorReq); err != nil {
		fmt.Printf("✓ Multiple validation errors caught:\n%v\n", err)
	}
}
