package tracer

import (
	"errors"
	"fmt"
	"github.com/NubeIO/reactive/helpers"
	"time"
)

func (ms *Tracer) LoggerType() string {
	return ms.logger.Level.String()
}

type Message struct {
	UUID       string    `json:"uuid" sql:"uuid" gorm:"type:varchar(255);unique;primaryKey"`
	TracerUUID string    `json:"tracerUUID,omitempty" gorm:"references tracers;not null;default:null"`
	Path       string    `json:"path"` // points
	Text       string    `json:"text"`
	AddToDisk  bool      `json:"-" gorm:"-"`
	LoggerType string    `json:"type"` // "info", "debug", "error", "warning", "disabled"
	Timestamp  time.Time `json:"timestamp"`
}

// AddMessage adds a new message to the SQLite database.
func (ms *Tracer) AddMessage(path, text, loggerType string, addToDisk bool) (*Message, error) {
	newMessage := &Message{
		UUID:       helpers.UUID(),
		TracerUUID: ms.UUID,
		Path:       path,
		Text:       text,
		AddToDisk:  addToDisk,
		LoggerType: loggerType, // Assign the current logger type
		Timestamp:  time.Now(), // Timestamp when the message is added
	}
	// Log the new message with additional details
	logMessage := fmt.Sprintf("TS:%s UUID: %s: Path: %s ->: %s", newMessage.Timestamp.Format(time.DateTime), newMessage.UUID, newMessage.Path, newMessage.Text)

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

func (ms *Tracer) getInMemoryMessagesNoDisk() []*Message {
	var unsavedMessages []*Message
	for _, message := range ms.unsavedMessages {
		if message.AddToDisk {
			unsavedMessages = append(unsavedMessages, message)
		}
	}
	return unsavedMessages

}

func (ms *Tracer) getAllMessagesByTracer(tracerUUID string) ([]*Message, error) {
	t, err := ms.GetMessagesByTracerUUID(tracerUUID)
	if err != nil {
		return nil, err
	}
	if t == nil {
		return nil, errors.New("tracer is nil")
	}
	return t, err
}

func (ms *Tracer) GetAllMessages() ([]*Message, error) {
	if ms.db == nil {
		return nil, errors.New("GetAllMessages() database has not been initialised yet")
	}

	var messagesFromDB []*Message

	// Retrieve messages from the database
	if err := ms.db.Find(&messagesFromDB).Error; err != nil {
		return nil, fmt.Errorf("error retrieving messages from the database: %v", err)
	}
	return messagesFromDB, nil
}

func (ms *Tracer) GetAllMessagesCombine(byTracerUUIID string) ([]*Message, error) {
	if ms.db == nil {
		return nil, errors.New("GetAllMessages() database has not been initialised yet")
	}

	// Retrieve messages from the database
	messagesFromDB, err := ms.getAllMessagesByTracer(byTracerUUIID)
	if err != nil {
		return nil, err
	}
	var messagesFromMemory = ms.getInMemoryMessagesNoDisk()

	// Combine messages from memory and database
	allMessages := append(messagesFromDB, messagesFromMemory...)

	return allMessages, nil
}

// SaveMessagesToDB saves unsaved messages to the database.
func (ms *Tracer) SaveMessagesToDB(maxTableSize int) error {
	mux.Lock()
	defer mux.Unlock()

	if ms.UUID == "" {
		return errors.New("SaveMessagesToDB() tracer-uuid can not be empty")
	}
	if ms.db == nil {
		return errors.New("SaveMessagesToDB() database has not been initialised yet")
	}

	// Combine messages from memory and the database
	allMessages, err := ms.GetAllMessagesCombine(ms.UUID)
	if err != nil {
		return err
	}
	allMessagesDisk, err := ms.GetMessagesByTracerUUID(ms.UUID)
	if err != nil {
		return err
	}

	ms.DebugfNotify("tracer messages on disk: %d messages in memory: %d & max table size: %d", len(allMessagesDisk), len(ms.getInMemoryMessagesNoDisk()), maxTableSize)

	// Check if the number of messages exceeds the maximum table size
	if len(allMessages) > maxTableSize {
		numToRemove := len(allMessages) - maxTableSize

		// Limit numToRemove to the available messages
		if numToRemove > len(ms.unsavedMessages) {
			numToRemove = len(ms.unsavedMessages)
		}

		// Calculate the timestamp threshold for retaining messages
		thresholdTime := time.Now().Add(-time.Duration(maxTableSize) * time.Second)

		// Remove the oldest messages from memory based on their timestamp
		for i := 0; i < numToRemove; i++ {
			if allMessages[i].Timestamp.Before(thresholdTime) {
				numToRemove--
			} else {
				break
			}
		}

		// Construct a slice of UUIDs for the oldest messages to be deleted
		oldestMessageUUIDs := make([]string, numToRemove)
		for i := 0; i < numToRemove; i++ {
			oldestMessageUUIDs[i] = allMessages[i].UUID
		}
		ms.DebugfNotify("tracer messages to delete from DB count: %d", len(oldestMessageUUIDs))

		// Bulk delete the oldest messages from the database based on their UUIDs
		if err := ms.db.Where("uuid IN ?", oldestMessageUUIDs).Delete(&Message{}).Error; err != nil {
			return fmt.Errorf("error bulk deleting oldest messages from the database: %v", err)
		}

		// Remove the deleted messages from the slice
		ms.unsavedMessages = ms.unsavedMessages[numToRemove:]
	}

	ms.DebugfNotify("tracer messages to add to DB count: %d", len(ms.unsavedMessages))
	if len(ms.unsavedMessages) > 0 {
		// Bulk save unsaved messages to the database and associate them with the specified tracer
		for _, msg := range ms.unsavedMessages {
			msg.TracerUUID = ms.UUID
		}
		if err := ms.db.Create(&ms.unsavedMessages).Error; err != nil {
			return fmt.Errorf("error bulk saving unsaved messages to the database: %v", err)
		}
	}

	// Clear unsaved messages in memory
	ms.unsavedMessages = []*Message{}

	return nil
}

func (ms *Tracer) saveMessagesNonBatch(maxTableSize int) error {
	mux.Lock()
	defer mux.Unlock()

	if ms.UUID == "" {
		return errors.New("SaveMessagesToDB() tracer-uuid can not be empty")
	}
	if ms.db == nil {
		return errors.New("SaveMessagesToDB() database has not been initialised yet")
	}

	// Combine messages from memory and the database
	allMessages, err := ms.GetAllMessages()
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
	ms.Path = path
	ms.Application = application
}

func (ms *Tracer) Debugf(format string, args ...any) *Message {
	message, err := ms.AddMessage(ms.Path, fmt.Sprintf(format, args...), debug, true)
	if err != nil {
		return nil
	}
	return message
}

func (ms *Tracer) Debug(data ...any) *Message {
	message, err := ms.AddMessage(ms.Path, joinString(data), debug, true)
	if err != nil {
		return nil
	}
	return message
}

func (ms *Tracer) DebugfNotify(format string, args ...any) *Message {
	message, err := ms.AddMessage(ms.Path, fmt.Sprintf(format, args...), debug, false)
	if err != nil {
		return nil
	}
	return message
}

func (ms *Tracer) DebugNotify(data ...any) *Message {
	message, err := ms.AddMessage(ms.Path, joinString(data), debug, false)
	if err != nil {
		return nil
	}
	return message
}

func (ms *Tracer) ErrorfNotify(format string, args ...any) *Message {
	message, err := ms.AddMessage(ms.Path, fmt.Sprintf(format, args...), info, false)
	if err != nil {
		return nil
	}
	return message
}

func (ms *Tracer) ErrorNotify(data ...any) *Message {
	message, err := ms.AddMessage(ms.Path, joinString(data), errorType, true)
	if err != nil {
		return nil
	}
	return message
}

func (ms *Tracer) InfoNotify(args ...any) *Message {
	message, err := ms.AddMessage(ms.Path, joinString(args), info, false)
	if err != nil {
		return nil
	}
	return message
}

func (ms *Tracer) InfofNotify(format string, args ...any) *Message {
	message, err := ms.AddMessage(ms.Path, fmt.Sprintf(format, args...), info, false)
	if err != nil {
		return nil
	}
	return message
}

func (ms *Tracer) Infof(format string, args ...any) *Message {
	message, err := ms.AddMessage(ms.Path, fmt.Sprintf(format, args...), info, true)
	if err != nil {
		return nil
	}
	return message
}

func (ms *Tracer) Info(args ...any) *Message {
	message, err := ms.AddMessage(ms.Path, joinString(args), info, true)
	if err != nil {
		return nil
	}
	return message
}

func (ms *Tracer) Errorf(format string, args ...any) *Message {
	message, err := ms.AddMessage(ms.Path, fmt.Sprintf(format, args...), errorType, true)
	if err != nil {
		return nil
	}
	return message
}

func (ms *Tracer) Error(args ...any) *Message {
	message, err := ms.AddMessage(ms.Path, joinString(args), errorType, true)
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
