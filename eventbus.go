package reactive

import (
	"fmt"
	"sync"
)

type message struct {
	Port     *Port
	NodeUUID string
	NodeID   string
}

// EventBus manages event subscriptions and publishes events.
type EventBus struct {
	mu          sync.Mutex
	handlers    map[string][]chan *message
	subscribers map[chan *message]string
}

// NewEventBus creates a new EventBus.
func NewEventBus() *EventBus {
	return &EventBus{
		handlers:    make(map[string][]chan *message),
		subscribers: make(map[chan *message]string),
	}
}

func (eb *EventBus) Subscribe(topic string, ch chan *message) {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	eb.handlers[topic] = append(eb.handlers[topic], ch)
	eb.subscribers[ch] = topic
}

// Unsubscribe unsubscribes a channel from a topic.
func (eb *EventBus) Unsubscribe(topic string, ch chan *message) {
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
func (eb *EventBus) Publish(topic string, data *message) {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	for _, ch := range eb.handlers[topic] {
		go func(ch chan *message) {
			ch <- data
		}(ch)
	}
}
