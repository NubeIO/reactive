package reactive

import (
	message "github.com/NubeIO/reactive/tracer"
	"github.com/NubeIO/schema"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"sync"
)

type BaseNode struct {
	EventBus       *EventBus
	ID             string
	UUID           string
	parentUUID     string
	Name           string
	Inputs         []*Port
	Outputs        []*Port
	LastValue      map[string]*Port
	Bus            map[string]chan *Message
	Connections    []*Connection
	settings       *Settings
	data           map[string]any
	nodeDetails    *Details
	Schema         *schema.Generated
	meta           *Meta
	mux            sync.Mutex
	PublishOnTopic bool // if its set to true we will publish its parent info as a topic eg; myFolder/bacnetPoint
	allowHotFix    bool
	loaded         bool
	runtimeNodes   map[string]Node
	childNodes     map[string]Node
	tracer         *message.Tracer
	db             *gorm.DB
	logger         *logrus.Logger
}

// NewBaseNode creates a new BaseNode with the given ID, name, EventBus, and Flow.
func NewBaseNode(id, nodeUUID, name string, bus *EventBus, opts *Options) *BaseNode {
	if nodeUUID == "" {
		nodeUUID = generateShortUUID()
	}
	n := &BaseNode{
		EventBus:    bus,
		ID:          id,
		Name:        name,
		UUID:        nodeUUID,
		Inputs:      []*Port{},
		Outputs:     []*Port{},
		Bus:         make(map[string]chan *Message),
		LastValue:   make(map[string]*Port),
		Connections: nil,
		allowHotFix: false,
		childNodes:  make(map[string]Node),
		data:        make(map[string]any),
	}
	n.setMeta(opts)
	return n
}
