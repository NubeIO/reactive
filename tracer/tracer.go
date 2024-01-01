package tracer

import (
	"errors"
	"fmt"
	"github.com/NubeIO/reactive/helpers"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const (
	info      = "info"
	errorType = "error"
	debug     = "debug"
	warning   = "warning"
)

type Tracer struct {
	UUID            string     `json:"uuid" sql:"uuid" gorm:"type:varchar(255);unique;primaryKey"`
	Path            string     // plugin, service name
	Application     string     // modbus
	Key             string     // could like modbus read-coil, something common
	instanceUUID    string     // node uuid
	Messages        []*Message `json:"messages,omitempty" gorm:"constraint:OnDelete:CASCADE"`
	unsavedMessages []*Message // Store unsaved messages in memory
	db              *gorm.DB
	logger          *logrus.Logger // Logger for logging
}

func NewTracer(path, application string, logger *logrus.Logger, db *gorm.DB) *Tracer {
	return &Tracer{
		Path:            path,
		Application:     application,
		Messages:        []*Message{},
		unsavedMessages: []*Message{},
		logger:          logger,
		db:              db,
	}
}

// GetAllTracers retrieves all tracers from the database.
func (ms *Tracer) GetAllTracers() ([]*Tracer, error) {
	var tracers []*Tracer
	if err := ms.db.Preload("Messages").Find(&tracers).Error; err != nil {
		return nil, fmt.Errorf("error retrieving tracers: %v", err)
	}
	return tracers, nil
}

func (ms *Tracer) GetMessagesByTracerUUID(tracerUUID string) ([]*Message, error) {
	if ms.db == nil {
		return nil, errors.New("GetMessagesByTracerUUID() database has not been initialised yet")
	}

	var messages []*Message

	// Retrieve messages with the specified tracerUUID from the database
	if err := ms.db.Where("tracer_uuid = ?", tracerUUID).Find(&messages).Error; err != nil {
		return nil, fmt.Errorf("error retrieving messages for tracerUUID %s: %v", tracerUUID, err)
	}

	return messages, nil
}

// AddTracer adds a new tracer to the database.
func (ms *Tracer) AddTracer(instanceUUID, key string) error {
	tracer := &Tracer{
		UUID:         helpers.UUID(),
		Path:         ms.Path,
		Application:  ms.Application,
		Key:          key,
		instanceUUID: instanceUUID,
	}
	ms.UUID = tracer.UUID
	if err := ms.db.Create(tracer).Error; err != nil {
		return fmt.Errorf("error creating tracer: %v", err)
	}
	return nil
}

// TracerKey key could like modbus read-coil or something common
func (ms *Tracer) TracerKey(key string) {
	ms.Key = key
}

// UpdateTracer updates an existing tracer in the database.
func (ms *Tracer) UpdateTracer(tracer *Tracer) error {
	if err := ms.db.Save(tracer).Error; err != nil {
		return fmt.Errorf("error updating tracer: %v", err)
	}
	return nil
}

// DeleteTracer deletes a tracer by UUID from the database.
func (ms *Tracer) DeleteTracer(uuid string) error {
	if err := ms.db.Where("uuid = ?", uuid).Delete(&Tracer{}).Error; err != nil {
		return fmt.Errorf("error deleting tracer: %v", err)
	}
	return nil
}
