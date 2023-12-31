package tracer

import (
	"errors"
	"fmt"
	"github.com/NubeIO/reactive/helpers"
	"time"
)

// TS: 2023 UUID: abc134 PATH: my-plugin APP: modbus KEY: nodeUUID

type Message struct {
	UUID       string    `json:"uuid" sql:"uuid" gorm:"type:varchar(255);unique;primaryKey"`
	TracerUUID string    `json:"tracerUUID,omitempty" gorm:"references tracers;not null;default:null"`
	Path       string    `json:"path"` // points
	Text       string    `json:"text"`
	LoggerType string    `json:"type"` // "info", "debug", "error", "warning", "disabled"
	Timestamp  time.Time `json:"timestamp"`
}

// AddMessage adds a new message to the SQLite database.
func (ms *Tracer) AddMessage(path, text, loggerType string) (*Message, error) {
	newMessage := &Message{
		UUID:       helpers.UUID(),
		Path:       path,
		Text:       text,
		LoggerType: loggerType, // Assign the current logger type
		Timestamp:  time.Now(), // Timestamp when the message is added
		TracerUUID: ms.UUID,
	}
	// Log the new message with additional details
	logMessage := fmt.Sprintf("TS:%s UUID: %s: Path: %s ->: %s", newMessage.Timestamp.Format(time.RFC850), newMessage.UUID, newMessage.Path, newMessage.Text)

	switch loggerType {
	case info:
		ms.logger.Info(logMessage)
	case errorType:
		ms.logger.Error(logMessage)
	case debug:
		ms.logger.Debug(logMessage)
	case warning:
		ms.logger.Warning(logMessage)
	}
	ms.unsavedMessages = append(ms.unsavedMessages, newMessage)

	return newMessage, nil
}

func (ms *Tracer) GetInMemoryMessages() []*Message {
	return ms.unsavedMessages
}

func (ms *Tracer) GetAllMessages(fromMemory bool) ([]*Message, error) {

	var messagesFromDB []*Message

	// Retrieve messages from the database
	if err := ms.db.Find(&messagesFromDB).Error; err != nil {
		return nil, fmt.Errorf("error retrieving messages from the database: %v", err)
	}

	var messagesFromMemory []*Message

	// Retrieve messages from memory if alsoFromMemory is true
	if fromMemory {
		for _, msg := range ms.Messages {
			messagesFromMemory = append(messagesFromMemory, msg)
		}
	}

	// Combine messages from memory and database
	allMessages := append(messagesFromDB, messagesFromMemory...)

	return allMessages, nil
}

func (ms *Tracer) SaveMessagesToDB(maxTableSize int) error {
	mux.Lock()
	defer mux.Unlock()

	if ms.UUID == "" {
		return errors.New("tracer-uuid can not be empty")
	}

	// Combine messages from memory and the database
	allMessages, err := ms.GetAllMessages(true)
	if err != nil {
		return err
	}

	// Check if the number of messages exceeds the maximum table size
	if len(allMessages) > maxTableSize {
		numToRemove := len(allMessages) - maxTableSize

		// Limit numToRemove to the available messages
		if numToRemove > len(ms.unsavedMessages) {
			numToRemove = len(ms.unsavedMessages)
		}

		// Delete the oldest messages from the database
		for i := 0; i < numToRemove; i++ {
			if err := ms.db.Delete(allMessages[i]).Error; err != nil {
				return fmt.Errorf("error deleting oldest message from the database: %v", err)
			}
		}

		// Remove the deleted messages from the slice
		ms.unsavedMessages = ms.unsavedMessages[numToRemove:]
	}

	// Save unsaved messages to the database and associate them with the specified tracer
	for _, msg := range ms.unsavedMessages {
		msg.TracerUUID = ms.UUID
		if err := ms.db.Create(msg).Error; err != nil {
			return fmt.Errorf("error saving unsaved message to the database: %v", err)
		}
	}

	// Clear unsaved messages in memory
	ms.unsavedMessages = []*Message{}

	return nil
}

// GetTracerMessages retrieves all messages associated with a tracer from the database.
func (ms *Tracer) GetTracerMessages(uuid string) ([]*Message, error) {
	var messages []*Message
	if err := ms.db.Where("tracer_uuid = ?", uuid).Find(&messages).Error; err != nil {
		return nil, fmt.Errorf("error retrieving messages for tracer: %v", err)
	}
	return messages, nil
}

func (ms *Tracer) Setup(path, application string) {
	ms.path = path
	ms.application = application
}

func (ms *Tracer) Debug(data ...any) *Message {
	message, err := ms.AddMessage(ms.path, joinString(data), debug)
	if err != nil {
		return nil
	}
	return message
}

func (ms *Tracer) Info(data ...any) *Message {
	message, err := ms.AddMessage(ms.path, joinString(data), info)
	if err != nil {
		return nil
	}
	return message
}

func (ms *Tracer) Error(data ...any) *Message {
	message, err := ms.AddMessage(ms.path, joinString(data), errorType)
	if err != nil {
		return nil
	}
	return message
}

func joinString(s ...any) string {
	var result string
	for i, item := range s {
		if i > 0 {
			result += " "
		}
		result += fmt.Sprintf("%v", item)
	}
	return result
}
