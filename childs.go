package reactive

// RegisterChildNode registers a node as a child
func (n *BaseNode) RegisterChildNode(child Node) {
	n.childNodes[child.GetUUID()] = child
}

// GetChildNodes returns a slice of child nodes
func (n *BaseNode) GetChildNodes() []Node {
	children := make([]Node, 0, len(n.childNodes))
	for _, child := range n.childNodes {
		children = append(children, child)
	}
	return children
}

func (n *BaseNode) GetChildsByType(nodeID string) []Node {
	var childrenByType []Node
	for _, child := range n.childNodes {
		if child.GetID() == nodeID {
			childrenByType = append(childrenByType, child)
		}
	}
	return childrenByType
}

// GetChildNode returns a child node by its UUID
func (n *BaseNode) GetChildNode(uuid string) Node {
	return n.childNodes[uuid]
}

// GetPortValuesChildNode returns the port values of a specific child node
func (n *BaseNode) GetPortValuesChildNode(uuid string) []*Port {
	child, exists := n.childNodes[uuid]
	if !exists {
		return nil
	}
	return child.GetAllPortValues() // Assuming GetAllPortValues is a method that returns []*Port
}

func (n *BaseNode) SetLastValueChildNode(uuid string, port *Port) {
	child, exists := n.childNodes[uuid]
	if !exists {
		return
	}
	child.SetLastValue(port) // Assuming SetLastValue is a method that sets the last value of a port
}
