package reactive

import "fmt"

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
	allNodes := n.GetRuntimeNodes()
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
