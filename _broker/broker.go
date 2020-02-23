package _broker

import (
	"sync"
	"time"

	"github.com/nats-io/nats.go"
)

type Broker struct {
	sync.Mutex
	ServiceName string
	connection  *nats.Conn
	registry    Registry // Will hold a list of our services.
}

type Message struct {
	From      string      `json:"sender"`              // Used for logs.
	To        string      `json:"receiver"`            // Used for logs.
	ReplyWith string      `json:"replyWith,omitempty"` // If there's no reply with, it means we're not responding to anyone the message.
	Action    string      `json:"action"`              // Used for internal controller call
	Payload   interface{} `json:"content"`             // Payload to the internal controller we're acting on.
	Meta      interface{} `json:"meta"`                // Any additional information
	Token     string      `json:"token,omitempty"`     // Token optional for auth. for services that requires auth
	TimeSent  time.Time   `json:"timeSent"`            // Used for logging and terminating requests beyond a specific threshold.
}

type ActionHandler func(p map[string]interface{}) utils.Result

type Registry map[string]ActionHandler
