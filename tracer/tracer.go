package tracer

import (
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
	UUID            string `json:"uuid" sql:"uuid" gorm:"type:varchar(255);unique;primaryKey"`
	path            string // plugin, service name
	application     string // modbus
	key             string
	instanceUUID    string     // node uuid
	Messages        []*Message `json:"messages,omitempty" gorm:"constraint:OnDelete:CASCADE"`
	unsavedMessages []*Message // Store unsaved messages in memory
	db              *gorm.DB
	logger          *logrus.Logger // Logger for logging
}

func NewTracer(path, application string, logger *logrus.Logger, db *gorm.DB) *Tracer {
	return &Tracer{
		path:            path,
		application:     application,
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

// GetTracer retrieves a tracer by UUID from the database.
func (ms *Tracer) GetTracer(uuid string) (*Tracer, error) {
	var tracer Tracer
	if err := ms.db.Preload("Messages").Where("uuid = ?", uuid).First(&tracer).Error; err != nil {
		return nil, fmt.Errorf("error retrieving tracer: %v", err)
	}
	return &tracer, nil
}

// AddTracer adds a new tracer to the database.
func (ms *Tracer) AddTracer(instanceUUID, key string) error {
	tracer := &Tracer{
		UUID:         helpers.UUID(),
		path:         ms.path,
		application:  ms.application,
		key:          key,
		instanceUUID: instanceUUID,
	}
	ms.UUID = tracer.UUID
	if err := ms.db.Create(tracer).Error; err != nil {
		return fmt.Errorf("error creating tracer: %v", err)
	}
	return nil
}

// UpdateTracer updates an existing tracer in the database.
func (ms *Tracer) UpdateKey(key string) {
	ms.key = key
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
