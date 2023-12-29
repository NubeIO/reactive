package reactive

import (
	"fmt"
	"github.com/NubeIO/schema"
	"github.com/google/uuid"
	"log"
	"reflect"
	"strings"
	"sync"
)

type portDataType string

const (
	portTypeAny    portDataType = "any"
	portTypeFloat  portDataType = "float"
	portTypeString portDataType = "string"
	portTypeBool   portDataType = "bool"
)

type flowDirection string

const (
	DirectionSubscriber flowDirection = "subscriber"
	DirectionPublisher  flowDirection = "publisher"
)

type portDirection string

const (
	input  portDirection = "input"
	output portDirection = "output"
)

type Details struct {
	Category string  `json:"category"`
	ParentID *string `json:"parentID"`
}

type Node interface {
	New(nodeUUID, name string, bus *EventBus, settings *Settings, opts *Options) Node
	SetDetails(details *Details)
	GetDetails() *Details
	Start()
	Delete()
	GetUUID() string
	GetParentUUID() string
	GetID() string
	GetNodeName() string
	NewPort(port *Port)
	GetInput(id string) *Port
	GetInputs() []*Port
	GetOutputs() []*Port
	SetInputValue(id string, value interface{})
	GetAllNodeValues() []*NodeValue
	GetAllPortValues() []*Port
	GetAllInputValues() []*Port
	GetAllOutputValues() []*Port
	SetLastValue(port *Port)
	GetPortValue(portID string) (*Port, error)
	GetSchema() *schema.Generated
	BuildSchema()
	AddSettings(settings *Settings)
	GetSettings() *Settings
	SetMeta(opts *Options)
	GetMeta() *Meta
	AddConnection(connection *Connection)
	GetConnections() []*Connection
	UpdateConnections(connections []*Connection)
	UpdateSettings(settings *Settings)
	SetHotFix()
	HotFix() bool
	SetLoaded(set bool)
	Loaded() bool
	NotLoaded() bool
	AddRuntime(runtimeNodes map[string]Node)
	GetRuntimeNodes() map[string]Node
	AddToNodeToRuntime(node Node) Node
	RemoveNodeFromRuntime()

	RegisterChildNode(child Node)
	GetChildNodes() []Node
	GetChildNode(uuid string) Node
	GetChildsByType(nodeID string) []Node
	GetPortValuesChildNode(uuid string) []*Port
	SetLastValueChildNode(uuid string, port *Port)
}

func (n *BaseNode) New(nodeUUID, name string, bus *EventBus, settings *Settings, opts *Options) Node {
	return n
}

var runtimeNodesMutex sync.Mutex

func (n *BaseNode) GetUUID() string {
	return n.UUID
}

func (n *BaseNode) GetID() string {
	return n.ID
}

func (n *BaseNode) GetNodeName() string {
	return n.Name
}

func (n *BaseNode) GetInputs() []*Port {
	return n.Inputs
}

func (n *BaseNode) GetOutputs() []*Port {
	return n.Outputs
}

func (n *BaseNode) Start() {}

func (n *BaseNode) Delete() {
	n.RemoveNodeFromRuntime()
}

func (n *BaseNode) HotFix() bool {
	return n.allowHotFix
}

func (n *BaseNode) SetHotFix() {
	n.allowHotFix = true
}

func (n *BaseNode) AddRuntime(runtimeNodes map[string]Node) {
	n.runtimeNodes = runtimeNodes
}

func (n *BaseNode) GetRuntimeNodes() map[string]Node {
	return n.runtimeNodes
}

func (n *BaseNode) AddToNodeToRuntime(node Node) Node {
	runtimeNodesMutex.Lock()
	defer runtimeNodesMutex.Unlock()
	n.runtimeNodes[node.GetUUID()] = node
	return n.runtimeNodes[node.GetUUID()]
}

func (n *BaseNode) RemoveNodeFromRuntime() {
	runtimeNodesMutex.Lock()
	defer runtimeNodesMutex.Unlock()
	delete(n.runtimeNodes, n.UUID)
}

func (n *BaseNode) SetLoaded(set bool) {
	n.loaded = set
}

func (n *BaseNode) Loaded() bool {
	return n.loaded
}

func (n *BaseNode) NotLoaded() bool {
	return !n.loaded
}

// ---------------------------- EVENT BUS -------------------------- //

func (n *BaseNode) setPortTopic(portId string) string {
	return fmt.Sprintf("%s-%s", n.UUID, portId)
}

func (n *BaseNode) PublishMessage(port *Port, setLastValue ...bool) {
	if port.Name == "" {
		log.Fatalf("port name can not be empty")
	}
	topic := n.setPortTopic(port.ID)
	m := &message{
		Port:     port,
		NodeUUID: n.GetUUID(),
		NodeID:   n.GetID(),
	}
	if len(setLastValue) >= 0 {
		go n.SetLastValue(port)
	}
	go n.EventBus.Publish(topic, m)
	fmt.Printf("Published message from node: (name: %s uuid: %s) to topic: %s value %v\n", n.GetID(), n.GetUUID(), topic, printValue(port.Value))
}

type Options struct {
	addToNodesMap bool
	Meta          *Meta
}

func generateShortUUID() string {
	uuidWithoutHyphens := strings.ReplaceAll(uuid.New().String(), "-", "")
	shortUUID := uuidWithoutHyphens[:10]
	return shortUUID
}

// printValue converts a value to a string representation.
func printValue(value interface{}) string {
	val := reflect.ValueOf(value)
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return "value is empty"
		}
		return fmt.Sprintf("%v", val.Elem())
	}
	return fmt.Sprintf("%v", value)
}
