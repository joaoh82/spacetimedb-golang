package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joaoh82/spacetimedb-golang/client"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	dbURL := os.Getenv("SPACETIMEDB_URL")
	if dbURL == "" {
		log.Fatal("SPACETIMEDB_URL environment variable is not set")
	}

	// Create a new client
	spacetimeClient, err := client.NewClient(dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer spacetimeClient.Close()
	fmt.Println("spacetimeClient", spacetimeClient)

	// Create a new identity
	identityResp, err := spacetimeClient.CreateIdentity()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Created identity: %s\n", identityResp.Identity)

	// Create a new client with the identity and token
	spacetimeClient, err = client.NewClient(
		dbURL,
		client.WithToken(identityResp.Token),
		client.WithIdentity(identityResp.Identity),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Verify the identity
	err = spacetimeClient.VerifyIdentity(identityResp.Identity)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Identity verified successfully")

	// Get databases owned by the identity
	databases, err := spacetimeClient.GetDatabases(identityResp.Identity)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Databases: %v\n", databases)

	if len(databases) > 0 {
		// Connect to the first database via WebSocket
		err = spacetimeClient.ConnectWebSocket(databases[0])
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Connected to database via WebSocket")

		// Send a test message
		err = spacetimeClient.SendMessage(map[string]interface{}{
			"type": "test",
			"data": "Hello, SpacetimeDB!",
		})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Sent test message")

		// Listen for messages for 5 seconds
		done := make(chan bool)
		go func() {
			for {
				message, err := spacetimeClient.ReceiveMessage()
				if err != nil {
					log.Printf("Error receiving message: %v\n", err)
					continue
				}
				fmt.Printf("Received message: %+v\n", message)
			}
		}()

		time.Sleep(5 * time.Second)
		done <- true
	}
}
