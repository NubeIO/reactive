package reactive

func (n *BaseNode) SetDetails(details *Details) {
	n.nodeDetails = details
}

func (n *BaseNode) GetDetails() *Details {
	return n.nodeDetails
}
