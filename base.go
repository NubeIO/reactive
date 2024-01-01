package reactive

import (
	"github.com/NubeIO/reactive/helpers"
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
	pluginName     string
	application    string // eg modbus-driver
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
	options        *Options
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

type Info struct {
	NodeID      string
	NodeUUID    string
	Name        string
	PluginName  string
	Application string
}

// NodeInfo
//   - NodeID
//   - NodeUUID
//   - Name
//   - PluginName
//   - Application
func NodeInfo(s ...string) *Info {
	setup := &Info{
		NodeID:     "",
		NodeUUID:   "",
		Name:       "",
		PluginName: "",
	}

	// Check the length of s and assign values accordingly
	if len(s) > 0 {
		setup.NodeID = s[0]
	}
	if len(s) > 1 {
		setup.NodeUUID = s[1]
	}
	if len(s) > 2 {
		setup.Name = s[2]
	}
	if len(s) > 3 {
		setup.PluginName = s[3]
	}
	if len(s) > 4 {
		setup.PluginName = s[4]
	}

	return setup
}

// NewBaseNode creates a new BaseNode with the given ID, name, EventBus, and Flow.
func NewBaseNode(n *Info, bus *EventBus, opts *Options) *BaseNode {
	if n == nil {
		n = &Info{}
	}
	if n.NodeUUID == "" {
		n.NodeUUID = helpers.UUID()
	}
	newNode := &BaseNode{
		EventBus:    bus,
		ID:          n.NodeID,
		Name:        n.Name,
		UUID:        n.NodeUUID,
		pluginName:  n.PluginName,
		application: n.Application,
		Inputs:      []*Port{},
		Outputs:     []*Port{},
		Bus:         make(map[string]chan *Message),
		LastValue:   make(map[string]*Port),
		Connections: nil,
		allowHotFix: false,
		childNodes:  make(map[string]Node),
		data:        make(map[string]any),
	}
	newNode.setOptions(opts)
	return newNode
}
