package main

import (
	"log"
	"time"

	"github.com/Jack-Gledhill/discord-rich-presence/api"

	"github.com/hugolgst/rich-go/client"
)

const (
	// APIBase is the base URL for the presence API
	APIBase = "http://localhost:8080"
	// AppID is the Discord application ID that is used to connect to the RPC
	AppID = "1302660616741326959"
	// RetryInterval is the time to wait between connection attempts
	RetryInterval = 60 * time.Second
	// UpdateInterval is the time between presence updates
	UpdateInterval = 30 * time.Second
)

func init() {
	// Repeatedly try to connect to the Discord RPC
	// This stops the program crashing repeatedly if Discord isn't running
	var connected bool
	var err error

	for !connected {
		err = client.Login(AppID)
		if err != nil {
			log.Println("WARN | Connection to Discord RPC failed, retrying shortly...")
			log.Println("WARN | Error: ", err)
			time.Sleep(RetryInterval)
		} else {
			connected = true
		}
	}
}

// This function isn't called until init completes (i.e. we have successfully connected to the RPC)
// Once a connection is established, the program can continue indefinitely.
// If Discord is closed or crashes, the program will continue to run, and will just print the error
func main() {
	defer client.Logout() // Logout gracefully when the program exits

	// Set initial presence, this tests everything is working
	err := UpdatePresence()
	if err != nil {
		log.Fatalln("FATAL | Failed to set initial presence: ", err)
	}

	// Start presence update loop
	ticker := time.NewTicker(UpdateInterval)
	for range ticker.C {
		err := UpdatePresence()
		if err != nil {
			log.Println("WARN | Failed to complete presence update: ", err)
		}
	}
}

// UpdatePresence contains all the logic for fetching and setting the presence
func UpdatePresence() error {
	// Send GET request to presence API
	r, err := api.FetchPresence(APIBase)
	if err != nil {
		return err
	}

	// Parse API response into presence struct
	p, err := api.ParsePresence(r)
	if err != nil {
		return err
	}

	// Set presence in Discord RPC, but only if the API returned a presence
	if p != nil {
		err = client.SetActivity(*p)
		if err != nil {
			return err
		}
	}

	return nil
}
