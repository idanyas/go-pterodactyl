# go‑pterodactyl

An unofficial go SDK for the pterodactyl panel.
It covers both the Application (admin) and Client(user) APIs.


## Features
- Complete coverage of Application & Client endpoints
- ⚙Generics remove response‑unmarshalling boiler‑plate
- Context passed through every method for cancellation, time‑outs & tracing
- Helper methods (ListAll, etc.) hide pagination loops

## Quick Start

```go
client, err := pterodactyl.NewClient(
    "https://panel.example.com",
    os.Getenv("PTERO_TOKEN"),
    pterodactyl.ApplicationKey,
)
if err != nil { log.Fatal(err) }

nodes, err := client.ApplicationAPI.Nodes.ListAll(ctx, 0)
if err != nil { log.Fatal(err) }

fmt.Printf("There are %d nodes\n", len(nodes))
```

## Installation 
go get github.com/your‑org/go‑pterodactyl@latest