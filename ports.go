package reactive

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
