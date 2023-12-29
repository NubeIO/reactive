package reactive

import (
	"github.com/NubeIO/schema"
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
	Bus            map[string]chan *message
	Connections    []*Connection
	Settings       *Settings
	nodeDetails    *Details
	Schema         *schema.Generated
	meta           *Meta
	mux            sync.Mutex
	PublishOnTopic bool // if its set to true we will publish its parent info as a topic eg; myFolder/bacnetPoint
	allowHotFix    bool
	loaded         bool
	runtimeNodes   map[string]Node
	childNodes     map[string]Node
}

// NewBaseNode creates a new BaseNode with the given ID, name, EventBus, and Flow.
func NewBaseNode(id, nodeUUID, name string, bus *EventBus) *BaseNode {
	if nodeUUID == "" {
		nodeUUID = generateShortUUID()
	}

	return &BaseNode{
		EventBus:    bus,
		ID:          id,
		Name:        name,
		UUID:        nodeUUID,
		Inputs:      []*Port{},
		Outputs:     []*Port{},
		Bus:         make(map[string]chan *message),
		LastValue:   make(map[string]*Port),
		Connections: nil,
		allowHotFix: false,
		childNodes:  make(map[string]Node),
	}
}
