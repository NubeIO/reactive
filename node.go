package reactive

import (
	"fmt"
	"github.com/NubeIO/reactive/schemas"
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

type Node interface {
	New(nodeUUID, name string, bus *EventBus, opts *NodeOpts) Node
	Start()
	Delete()
	GetUUID() string
	GetID() string
	GetNodeName() string
	NewPort(port *Port)
	GetInput(id string) *Port
	GetInputs() []*Port
	GetOutputs() []*Port
	SetInputValue(id string, value interface{})
	GetAllNodes() []Node
	GetAllNodeValues() []*NodeValue
	GetAllPortValues() []*Port
	GetAllInputValues() []*Port
	GetAllOutputValues() []*Port
	SetLastValue(port *Port)
	GetPortValue(portID string) (*Port, error)
	BuildSchema() *schemas.Schema
	AddSettings(settings *Settings)
	GetSettings() *Settings
	GetMeta() *Meta
	AddConnection(connection *Connection)
	GetConnections() []*Connection
	UpdateConnections(connections []*Connection)
	UpdateSettings(settings *Settings)
	SetHotFix()
	HotFix() bool
}

type message struct {
	Port     *Port
	NodeUUID string
	NodeID   string
}

// ---------------------------- BASE NODE -------------------------- //

type BaseNode struct {
	EventBus       *EventBus
	ID             string
	UUID           string
	Name           string
	Inputs         []*Port
	Outputs        []*Port
	LastValue      map[string]*Port
	Bus            map[string]chan *message
	Connections    []*Connection
	Settings       *Settings
	schema         *schemas.Schema
	Meta           *Meta
	mux            sync.Mutex
	PublishOnTopic bool // if its set to true we will publish its parent info as a topic eg; myFolder/bacnetPoint
	allowHotFix    bool
	nodes          map[string]Node
	Node           Node
}

func (n *BaseNode) New(nodeUUID, name string, bus *EventBus, opts *NodeOpts) Node {
	return n
}

var runtimeNodesMutex sync.Mutex
var runtimeNodes = make(map[string]Node)

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
		nodes:       runtimeNodes, // Initialize the nodes map
		allowHotFix: false,
	}
}

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

func (n *BaseNode) Start() {
	fmt.Println("IN BASE")
}

func (n *BaseNode) Delete() {
	n.RemoveFromNodesMap()
}

func (n *BaseNode) HotFix() bool {
	return n.allowHotFix
}

func (n *BaseNode) SetHotFix() {
	n.allowHotFix = true
}

func (n *BaseNode) GetAllNodes() []Node {
	nodes := make([]Node, 0, len(n.nodes))
	for _, node := range n.nodes {
		nodes = append(nodes, node)
	}
	return nodes
}

func (n *BaseNode) AddToNodesMap(nodeUUID string, node Node) {
	runtimeNodesMutex.Lock()
	defer runtimeNodesMutex.Unlock()
	n.nodes[nodeUUID] = node
}

func (n *BaseNode) RemoveFromNodesMap() {
	runtimeNodesMutex.Lock()
	defer runtimeNodesMutex.Unlock()
	delete(n.nodes, n.UUID)
	delete(runtimeNodes, n.UUID)
}

// ---------------------------- NODE PORTS -------------------------- //

// Port represents a data port with an ID, Name, and Value.
type Port struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Value       interface{}   `json:"value,omitempty"`
	LastUpdated string        `json:"lastUpdated,omitempty"` // last time it got a message
	Direction   portDirection `json:"direction"`
	DataType    portDataType  `json:"dataType"`
}

type remoteConnection struct {
	hostUUID string
	ip       string
}

// Connection defines a structure for input subscriptions.
type Connection struct {
	SourceUUID    string        `json:"source"`
	SourcePort    string        `json:"sourceHandle"`
	TargetUUID    string        `json:"target"`
	TargetPort    string        `json:"targetHandle"`
	FlowDirection flowDirection `json:"flowDirection"` // subscriber is if it's in an input and publisher if It's for an output
	// remoteConnection may add this later, it would be that a message can be sent to a remote server
}

func (n *BaseNode) GetConnections() []*Connection {
	return n.Connections
}

func (n *BaseNode) NewInputPort(id, name string, dataType portDataType) {
	port := &Port{
		ID:        id,
		Name:      name,
		Value:     nil,
		Direction: input,
		DataType:  dataType,
	}
	n.NewPort(port)
}

func (n *BaseNode) NewOutputPort(id, name string, dataType portDataType) {
	port := &Port{
		ID:        id,
		Name:      name,
		Value:     nil,
		Direction: output,
		DataType:  dataType,
	}
	n.NewPort(port)
}

func (n *BaseNode) NewPort(port *Port) {
	if port.Direction == input {
		n.Inputs = append(n.Inputs, port)
		n.Bus[port.ID] = make(chan *message, 1)
	} else if port.Direction == output {
		n.Outputs = append(n.Outputs, port)
	}
}

// ---------------------------- NODE SETTINGS -------------------------- //

type Settings struct {
	Value interface{} `json:"value"`
}

func (s *Settings) GetFloat64Value() float64 {
	if s == nil {
		return 0
	}
	if floatValue, ok := s.Value.(float64); ok {
		return floatValue
	}
	return 0
}

func (s *Settings) GetFloat64ValuePointer() *float64 {
	if s == nil {
		return nil
	}
	if floatValue, ok := s.Value.(float64); ok {
		return &floatValue
	}
	return nil
}

type position struct {
	PositionY int `json:"positionY"`
	PositionX int `json:"positionX"`
}

type Meta struct {
	Position position `json:"position"`
}

func (n *BaseNode) GetMeta() *Meta {
	return n.Meta
}

func (n *BaseNode) GetSettings() *Settings {
	return n.Settings
}

func (n *BaseNode) BuildSchema() *schemas.Schema {
	return n.schema
}

func (n *BaseNode) AddSettings(settings *Settings) {
	n.Settings = settings
}

func (n *BaseNode) UpdateSettings(settings *Settings) {
	n.Settings = settings
}

// ---------------------------- NODE VALUES -------------------------- //

type NodeValue struct {
	NodeId   string  `json:"nodeId"`
	NodeUUID string  `json:"nodeUUID"`
	Ports    []*Port `json:"ports"`
}

type AllNodeValue struct {
	NodeId   string  `json:"nodeId"`
	NodeUUID string  `json:"nodeUUID"`
	Inputs   []*Port `json:"inputs,omitempty"`
	Outputs  []*Port `json:"outputs,omitempty"`
}

func (n *BaseNode) GetAllNodeValues() []*NodeValue {
	allNodes := n.GetAllNodes()
	nodeValues := make([]*NodeValue, len(allNodes))
	for _, node := range allNodes {
		nv := node.GetAllPortValues()
		if nv == nil {
			continue
		}
		portValue := &NodeValue{
			NodeId:   node.GetID(),
			NodeUUID: node.GetUUID(),
			Ports:    nv,
		}
		nodeValues = append(nodeValues, portValue)
	}
	return nodeValues
}

func (n *BaseNode) GetAllPortValues() []*Port {
	var out []*Port
	out = append(out, n.GetAllInputValues()...)
	out = append(out, n.GetAllOutputValues()...)
	return out
}

func (n *BaseNode) GetAllInputValues() []*Port {
	var out []*Port
	for _, port := range n.GetInputs() {
		value, err := n.GetPortValue(port.ID)
		if err != nil {
		} else {
			out = append(out, value)
		}
	}
	return out
}

func (n *BaseNode) GetAllOutputValues() []*Port {
	var out []*Port
	for _, port := range n.GetOutputs() {
		value, err := n.GetPortValue(port.ID)
		if err != nil {
		} else {
			out = append(out, value)
		}
	}
	return out
}

func (n *BaseNode) GetPortValue(portID string) (*Port, error) {
	lastValues := n.LastValue
	port, exists := lastValues[portID]
	if !exists {
		return nil, fmt.Errorf("port with ID %s not found", portID)
	}
	return port, nil
}

func (n *BaseNode) GetInput(id string) *Port {
	ports := n.GetInputs()
	for _, port := range ports {
		if port.ID == id {
			return port
		}
	}
	return nil
}

func (n *BaseNode) SetInputValue(id string, value interface{}) {
	port := n.GetInput(id)
	if port != nil {
		port.Value = value
	}
}

func (n *BaseNode) SetLastValue(port *Port) {
	n.mux.Lock() // Lock the mutex before accessing the shared resource
	defer n.mux.Unlock()
	// Check if the port exists in LastValue
	if existingPort, ok := n.LastValue[port.ID]; ok {
		// Update the value of the existing port
		existingPort.Value = port.Value
		n.LastValue[port.ID] = existingPort
	} else {
		// If the port doesn't exist, create a new port entry
		n.LastValue[port.ID] = port
	}
}

// ---------------------------- CONNECTIONS -------------------------- //

func (n *BaseNode) AddConnection(connection *Connection) {
	if connection == nil {
		panic("node connection can not be empty")
	}

	sourceNodeUUID := connection.SourceUUID
	targetNodeUUID := connection.TargetUUID
	sourceOutput := connection.SourcePort
	targetInput := connection.TargetPort

	sourceTopic := fmt.Sprintf("%s-%s", sourceNodeUUID, sourceOutput)
	fmt.Printf("Add new connection type: %s from: (%s-%s) to: (%s-%s) \n", connection.FlowDirection, sourceNodeUUID, sourceOutput, targetNodeUUID, targetInput)
	n.Bus[targetNodeUUID] = make(chan *message, 1)
	n.EventBus.Subscribe(sourceTopic, n.Bus[targetInput])
	subscriber := &Connection{
		SourceUUID:    sourceNodeUUID,
		SourcePort:    sourceOutput,
		TargetUUID:    targetNodeUUID,
		TargetPort:    targetInput,
		FlowDirection: connection.FlowDirection,
	}

	n.Connections = append(n.Connections, subscriber)
}

func (n *BaseNode) UpdateConnections(connections []*Connection) {
	// Iterate through existing connections
	for i := len(n.Connections) - 1; i >= 0; i-- {
		existingConn := n.Connections[i]
		// Check if the existing connection exists in the download connections
		found := false
		for _, uploadedConn := range connections {
			if existingConn.SourceUUID == uploadedConn.SourceUUID &&
				existingConn.SourcePort == uploadedConn.SourcePort &&
				existingConn.TargetUUID == uploadedConn.TargetUUID {
				found = true
				break
			}
		}
		// If the existing connection is not found in the download connections, remove it
		if !found {
			sourceTopic := fmt.Sprintf("%s-%s", existingConn.SourceUUID, existingConn.SourcePort)
			delete(n.Bus, existingConn.TargetUUID)
			n.EventBus.Unsubscribe(sourceTopic, n.Bus[existingConn.TargetUUID])
			// Remove the connection from the slice
			n.Connections = append(n.Connections[:i], n.Connections[i+1:]...)
			fmt.Printf("Removed input subscriber for topic: %s\n", sourceTopic)
		}
	}
	// Add or update the new connections
	for _, connection := range connections {
		n.AddConnection(connection)
	}
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

type NodeOpts struct {
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
