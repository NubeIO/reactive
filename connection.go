package reactive

import "fmt"

// Connection defines a structure for input subscriptions.
type Connection struct {
	SourceUUID    string        `json:"source"`
	SourcePort    string        `json:"sourceHandle"`
	TargetUUID    string        `json:"target"`
	TargetPort    string        `json:"targetHandle"`
	FlowDirection flowDirection `json:"flowDirection"` // subscriber is if it's in an input and publisher if It's for an output

}

func (n *BaseNode) GetConnections() []*Connection {
	return n.Connections
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
