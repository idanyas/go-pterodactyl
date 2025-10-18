# go-pterodactyl

[![Go Reference](https://pkg.go.dev/badge/github.com/idanyas/go-pterodactyl.svg)](https://pkg.go.dev/github.com/idanyas/go-pterodactyl)
[![Go Report Card](https://goreportcard.com/badge/github.com/idanyas/go-pterodactyl)](https://goreportcard.com/report/github.com/idanyas/go-pterodactyl)

A comprehensive, type-safe Go client library for the [Pterodactyl Panel](https://pterodactyl.io/) API. Supports both the Application API (admin) and Client API (user), plus real-time WebSocket connections.

## Features

- ✅ **Complete API Coverage**: Application API, Client API, and WebSocket support
- ✅ **Type-Safe**: Strongly-typed request/response structures with validation
- ✅ **Context-Aware**: All operations support context for timeouts and cancellation
- ✅ **Automatic Retries**: Configurable retry logic with exponential backoff
- ✅ **Rate Limiting**: Built-in rate limit detection and handling
- ✅ **Pagination**: Easy-to-use pagination for list endpoints
- ✅ **WebSocket**: Real-time server console with automatic reconnection
- ✅ **No Dependencies**: Minimal external dependencies (only websocket and validator)
- ✅ **Well-Tested**: Comprehensive test coverage
- ✅ **Production-Ready**: Used in production environments

## Installation

```bash
go get github.com/idanyas/go-pterodactyl
```

## Quick Start

### Application API (Admin)

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/idanyas/go-pterodactyl"
    "github.com/idanyas/go-pterodactyl/application"
)

func main() {
    // Create client
    client, err := pterodactyl.New(
        "https://panel.example.com",
        pterodactyl.WithAPIKey("ptla_your_application_api_key"),
    )
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()
    app := client.Application()

    // List all users
    users, _, err := app.ListUsers(ctx, pterodactyl.ListOptions{})
    if err != nil {
        log.Fatal(err)
    }

    for _, user := range users {
        fmt.Printf("User: %s (%s)\n", user.Username, user.Email)
    }

    // Create a new server
    server, err := app.CreateServer(ctx, application.CreateServerRequest{
        Name: "My Game Server",
        User: 1,
        Egg:  1,
        Limits: models.Limits{
            Memory: 1024,
            Disk:   5120,
            CPU:    100,
        },
        FeatureLimits: models.FeatureLimits{
            Databases:   2,
            Allocations: 1,
            Backups:     5,
        },
        Allocation: application.CreateServerAllocation{
            Default: 1,
        },
    })
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Created server: %s (UUID: %s)\n", server.Name, server.UUID)
}
```

### Client API (User)

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/idanyas/go-pterodactyl"
)

func main() {
    // Create client
    client, err := pterodactyl.New(
        "https://panel.example.com",
        pterodactyl.WithAPIKey("ptlc_your_client_api_key"),
    )
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()
    clientAPI := client.Client()

    // List your servers
    servers, _, err := clientAPI.ListServers(ctx, pterodactyl.ListOptions{})
    if err != nil {
        log.Fatal(err)
    }

    for _, server := range servers {
        fmt.Printf("Server: %s (ID: %s)\n", server.Name, server.Identifier)

        // Get resource usage
        resources, err := clientAPI.GetServerResources(ctx, server.Identifier)
        if err != nil {
            log.Printf("Failed to get resources: %v", err)
            continue
        }

        fmt.Printf("  Status: %s\n", resources.CurrentState)
        fmt.Printf("  Memory: %d MB / %d MB\n",
            resources.Resources.MemoryBytes/1024/1024,
            resources.Resources.MemoryLimitBytes/1024/1024)
    }

    // Send a power action
    err = clientAPI.SendPowerAction(ctx, "server-id", "restart")
    if err != nil {
        log.Fatal(err)
    }
}
```

### WebSocket (Real-time Console)

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "os/signal"

    "github.com/idanyas/go-pterodactyl"
    "github.com/idanyas/go-pterodactyl/websocket"
)

func main() {
    client, err := pterodactyl.New(
        "https://panel.example.com",
        pterodactyl.WithAPIKey("ptlc_your_client_api_key"),
    )
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()
    clientAPI := client.Client()

    // Connect to WebSocket
    ws, err := clientAPI.ConnectWebSocket(ctx, "server-id")
    if err != nil {
        log.Fatal(err)
    }
    defer ws.Close()

    // Handle interrupt
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, os.Interrupt)

    // Listen for events
    for {
        select {
        case event := <-ws.Events():
            switch e := event.(type) {
            case *websocket.ConsoleOutputEvent:
                fmt.Printf("[Console] %s\n", e.Line)
            case *websocket.StatusEvent:
                fmt.Printf("[Status] %s\n", e.Status)
            case *websocket.StatsEvent:
                fmt.Printf("[Stats] CPU: %.2f%%, Memory: %d MB\n",
                    e.Stats.CPUAbsolute,
                    e.Stats.MemoryBytes/1024/1024)
            }
        case <-sigChan:
            return
        }
    }
}
```

## API Coverage

### Application API (Admin)

- ✅ **Users**: List, Get, Create, Update, Delete
- ✅ **Servers**: List, Get, Create, Update (details/build/startup), Suspend, Unsuspend, Reinstall, Delete
- ✅ **Databases**: List, Get, Create, Update, Reset Password, Delete
- ✅ **Nodes**: List, Get, Create, Update, Delete, Get Configuration, Allocations, Deployable Nodes
- ✅ **Locations**: List, Get, Create, Update, Delete
- ✅ **Nests & Eggs**: List Nests, Get Nest, List Eggs, Get Egg

### Client API (User)

- ✅ **Account**: Get Details, 2FA, Email, Password, API Keys, SSH Keys, Activity Logs
- ✅ **Servers**: List, Get, Get Resources, Power Actions, Send Commands, Activity Logs
- ✅ **Files**: List, Read, Write, Upload, Download, Create/Delete/Rename, Compress, Decompress, Chmod, Pull
- ✅ **Databases**: List, Create, Rotate Password, Delete
- ✅ **Backups**: List, Get, Create, Download, Restore, Lock/Unlock, Delete
- ✅ **Schedules**: List, Get, Create, Update, Delete, Execute, Tasks
- ✅ **Network**: List Allocations, Assign, Set Primary, Update Notes, Delete
- ✅ **Subusers**: List, Get, Create, Update, Delete
- ✅ **Startup**: Get Configuration, Update Variables
- ✅ **Settings**: Rename, Reinstall, Update Docker Image
- ✅ **Permissions**: Get System Permissions

### WebSocket

- ✅ **Events**: Console Output, Stats, Status, Token Expiration
- ✅ **Commands**: Send Command, Set State (power actions)
- ✅ **Reconnection**: Automatic reconnection with exponential backoff (opt-in)

## Advanced Usage

### Input Validation

All request types support validation:

```go
req := application.CreateUserRequest{
    Email:     "user@example.com",
    Username:  "johndoe",
    FirstName: "John",
    LastName:  "Doe",
}

// Validate before sending
if err := validation.Validate(req); err != nil {
    log.Fatal(err)
}
```

### Pagination

```go
// Initial request
users, paginator, err := app.ListUsers(ctx, pterodactyl.ListOptions{
    PerPage: 50,
})

// Iterate through pages
for paginator.HasMorePages() {
    nextUsers, err := paginator.NextPage(ctx)
    if err != nil {
        log.Fatal(err)
    }
    users = append(users, nextUsers...)
}
```

### Custom HTTP Client

```go
import "net/http"

httpClient := &http.Client{
    Timeout: 30 * time.Second,
    // Add custom transport, proxy, etc.
}

client, err := pterodactyl.New(
    "https://panel.example.com",
    pterodactyl.WithAPIKey("your-api-key"),
    pterodactyl.WithHTTPClient(httpClient),
)
```

### Helper Methods

```go
import "github.com/idanyas/go-pterodactyl/helpers"

// Wait for server to reach a state
waiter := helpers.NewStateWaiter(clientAPI)
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
defer cancel()

err := waiter.WaitForState(ctx, "server-id", "running", 5*time.Second)

// Download file to writer
downloader := helpers.NewFileDownloader(clientAPI)
file, _ := os.Create("backup.tar.gz")
defer file.Close()

err = downloader.DownloadToWriter(ctx, "server-id", "/backup.tar.gz", file)

// Create backup and wait for completion
manager := helpers.NewBackupManager(clientAPI)
backup, err := manager.CreateAndWait(ctx, "server-id", client.CreateBackupRequest{
    Name: "My Backup",
}, 10*time.Second)
```

## Error Handling

The library provides detailed error information:

```go
server, err := app.GetServer(ctx, 9999)
if err != nil {
    if apiErr, ok := err.(*pterodactyl.APIError); ok {
        fmt.Printf("API Error: Status %d\n", apiErr.StatusCode)
        for _, e := range apiErr.Errors {
            fmt.Printf("  - %s: %s\n", e.Code, e.Detail)
        }
    } else {
        fmt.Printf("Error: %v\n", err)
    }
}
```

## Rate Limiting

The library automatically handles rate limits by waiting until the reset time:

```go
// The client will automatically wait if rate limited
for i := 0; i < 1000; i++ {
    _, err := app.GetUser(ctx, i)
    // Will automatically pause when rate limit is hit
}
```

## Testing

Run tests:

```bash
go test ./...
```

Run tests with coverage:

```bash
go test -cover ./...
```

## License

MIT License - see LICENSE file for details
