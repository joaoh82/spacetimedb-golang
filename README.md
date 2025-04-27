# SpacetimeDB Golang SDK

[![Go Build and Test](https://github.com/joaoh82/spacetimedb-golang/actions/workflows/go.yml/badge.svg)](https://github.com/joaoh82/spacetimedb-golang/actions/workflows/go.yml)
[![Last Commit](https://img.shields.io/github/last-commit/joaoh82/spacetimedb-golang)](https://github.com/joaoh82/spacetimedb-golang/commits/main)

A Golang SDK for interacting with SpacetimeDB instances through both HTTP and WebSocket APIs.

## Installation

```bash
go get github.com/joaoh82/spacetimedb-golang
```

## Usage

### Creating a Client

```go
import "github.com/joaoh82/spacetimedb-golang/client"

// Create a new client
client, err := client.NewClient("https://your-spacetimedb-instance.com")
if err != nil {
    log.Fatal(err)
}
defer client.Close()
```

### Identity Management

```go
// Create a new identity
identityResp, err := client.CreateIdentity()
if err != nil {
    log.Fatal(err)
}

// Use the identity and token for subsequent operations
client = client.WithToken(identityResp.Token).WithIdentity(identityResp.Identity)

// Verify the identity
err = client.VerifyIdentity(identityResp.Identity)
if err != nil {
    log.Fatal(err)
}

// Get databases owned by the identity
databases, err := client.GetDatabases(identityResp.Identity)
if err != nil {
    log.Fatal(err)
}
```

### Database Operations

```go
// Get database information
dbInfo, err := client.GetDatabaseInfo("your-database-address")
if err != nil {
    log.Fatal(err)
}

// Execute a SQL query
result, err := client.ExecuteQuery("your-database-address", "SELECT * FROM your_table")
if err != nil {
    log.Fatal(err)
}
```

### WebSocket Connection

```go
// Connect to a database via WebSocket
err := client.ConnectWebSocket("your-database-address")
if err != nil {
    log.Fatal(err)
}

// Send a message
err = client.SendMessage(map[string]interface{}{
    "type": "your_message_type",
    "data": "your_message_data",
})
if err != nil {
    log.Fatal(err)
}

// Receive messages
for {
    message, err := client.ReceiveMessage()
    if err != nil {
        log.Fatal(err)
    }
    // Handle the message
    fmt.Printf("Received: %+v\n", message)
}
```

## Features

- Identity management (create, verify)
- Database operations (info, queries)
- WebSocket connection for real-time updates
- Token-based authentication
- Error handling and context support

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
