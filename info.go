package reactive

func (n *BaseNode) GetID() string {
	return n.ID
}

func (n *BaseNode) GetNodeName() string {
	return n.Name
}

func (n *BaseNode) GetPluginName() string {
	return n.pluginName
}

func (n *BaseNode) GetApplicationUse() string {
	return n.application
}
