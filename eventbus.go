package reactive

import (
	"fmt"
	"github.com/gorilla/websocket"
	"sync"
)

type Message struct {
	Port     *Port
	NodeUUID string
	NodeID   string
}

// EventBus manages event subscriptions and publishes events.
type EventBus struct {
	mu          sync.Mutex
	handlers    map[string][]chan *Message
	subscribers map[chan *Message]string
	WS          *WSHub
}

// NewEventBus creates a new EventBus.
func NewEventBus() *EventBus {
	ws := NewWSHub()
	go ws.Run()
	return &EventBus{
		handlers:    make(map[string][]chan *Message),
		subscribers: make(map[chan *Message]string),
		WS:          ws,
	}
}

func (eb *EventBus) Subscribe(topic string, ch chan *Message) {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	eb.handlers[topic] = append(eb.handlers[topic], ch)
	eb.subscribers[ch] = topic
}

// Unsubscribe unsubscribes a channel from a topic.
func (eb *EventBus) Unsubscribe(topic string, ch chan *Message) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	// Remove the channel from the list of subscribers for the topic
	subscribers := eb.handlers[topic]
	for i, sub := range subscribers {
		if sub == ch {
			close(sub)                 // Close the channel to stop the goroutine
			subscribers[i] = nil       // Set the channel to nil
			delete(eb.subscribers, ch) // Remove the subscriber entry
			break
		}
	}
	eb.handlers[topic] = subscribers // Update the subscribers list for the topic
	fmt.Printf("Unsubscribed from topic: %s\n", topic)
}

// Publish publishes an event to all subscribers of a topic.
func (eb *EventBus) Publish(topic string, data *Message) {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	for _, ch := range eb.handlers[topic] {
		go func(ch chan *Message) {
			ch <- data
		}(ch)
	}
}

type WSConnection struct {
	Conn *websocket.Conn
	Send chan *Message // Channel to Send messages

}

type WSHub struct {
	Connections map[*WSConnection]bool
	Register    chan *WSConnection
	Unregister  chan *WSConnection
	Broadcast   chan *Message
	mu          sync.Mutex // Mutex to protect concurrent access to connections

}

func NewWSConnection(conn *websocket.Conn) *WSConnection {
	fmt.Printf("ADD NEW WS Conntection :s", conn.LocalAddr().String())
	return &WSConnection{
		Conn: conn,
		Send: make(chan *Message),
	}
}

func (h *WSHub) Unsubscribe(conn *WSConnection) {
	h.mu.Lock() // Use a mutex to handle concurrent access
	defer h.mu.Unlock()

	if _, ok := h.Connections[conn]; ok {
		delete(h.Connections, conn)
		close(conn.Send)
	}
}

func NewWSHub() *WSHub {
	return &WSHub{
		Connections: make(map[*WSConnection]bool),
		Register:    make(chan *WSConnection),
		Unregister:  make(chan *WSConnection),
		Broadcast:   make(chan *Message),
	}
}

func (h *WSHub) Run() {
	for {
		select {
		case conn := <-h.Register:
			h.mu.Lock()
			h.Connections[conn] = true
			h.mu.Unlock()
		case msg := <-h.Broadcast:
			h.mu.Lock()
			for conn := range h.Connections {
				select {
				case conn.Send <- msg:
				default:
					h.Unsubscribe(conn)
				}
			}
			h.mu.Unlock()
		}
	}
}
