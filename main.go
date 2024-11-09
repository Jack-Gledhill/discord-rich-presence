package main

import (
	"log"
	"os"
	"time"

	"github.com/Jack-Gledhill/discord-rich-presence/api"

	"github.com/hugolgst/rich-go/client"
)

const (
	// RetryInterval is the time to wait between connection attempts
	RetryInterval = 60 * time.Second
	// UpdateInterval is the time between presence updates
	UpdateInterval = 30 * time.Second
)

var (
	// APIBase is the base URL for the presence API
	APIBase string
	// AppID is the Discord application ID that is used to connect to the RPC
	AppID string
)

func init() {
	// Get environment variables and crash if they're not set (because they're required)
	var ok bool
	APIBase, ok = os.LookupEnv("API_BASE")
	if !ok {
		log.Fatalln("FATAL | API_BASE environment variable not set")
	}

	AppID, ok = os.LookupEnv("APP_ID")
	if !ok {
		log.Fatalln("FATAL | APP_ID environment variable not set")
	}

	// Repeatedly try to connect to the Discord RPC
	// This stops the program crashing repeatedly if Discord isn't running
	var connected bool
	var err error

	for !connected {
		err = client.Login(AppID)
		if err != nil {
			log.Println("WARN | Connection to Discord RPC failed, retrying shortly...")
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
