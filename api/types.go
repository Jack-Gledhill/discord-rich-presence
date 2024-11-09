package api

import (
	"github.com/hugolgst/rich-go/client"
)

// Response is a wrapper for the presence API JSON response
type Response struct {
	Presence client.Activity `json:"presence"`
}
