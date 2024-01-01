package reactive

func (n *BaseNode) SetDetails(details *Details) {
	n.nodeDetails = details
}

func (n *BaseNode) GetDetails() *Details {
	return n.nodeDetails
}

func (n *BaseNode) SupportsDB() bool {
	return n.nodeDetails.HasDB
}

func (n *BaseNode) SupportsLogging() bool {
	return n.nodeDetails.HasLogger
}
