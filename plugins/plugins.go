package plugins

type Node struct {
	ID       string  `json:"id"`
	Export   string  `json:"-"`
	Children []*Node `json:"children,omitempty"`
}

type Export struct {
	Name        string             `json:"name"`
	Version     string             `json:"version"`
	Description string             `json:"description"`
	Nodes       map[string][]*Node `json:"nodes"`
}

func NewPlugin(name, version, description string) *Export {
	return &Export{
		Name:        name,
		Version:     version,
		Description: description,
		Nodes:       make(map[string][]*Node),
	}
}

func (p *Export) AddCategory(category string) {
	p.Nodes[category] = []*Node{}
}

func (p *Export) AddNode(category, id, export string) {
	node := &Node{
		ID:     id,
		Export: export,
	}
	p.Nodes[category] = append(p.Nodes[category], node)
}

func (p *Export) AddChildNode(category, parentID, childID, childExport string) {
	childNode := &Node{ID: childID, Export: childExport}

	for i, node := range p.Nodes[category] {
		if node.ID == parentID {
			p.Nodes[category][i].Children = append(p.Nodes[category][i].Children, childNode)
			return
		}
	}
}

func (p *Export) GetNode(nodeName string, pluginExport *Export) *Node {
	for _, category := range pluginExport.Nodes {
		for _, node := range category {
			if node.ID == nodeName {
				return node
			}
		}
	}
	return nil
}

func (p *Export) GetNodes(pluginExport *Export) []*Node {
	var nodes []*Node
	for _, category := range pluginExport.Nodes {
		for _, node := range category {
			nodes = append(nodes, node)
			nodes = append(nodes, node.Children...)
		}
	}
	return nodes
}

func (p *Export) GetChildNodes(pluginExport *Export) []*Node {
	var children []*Node
	for _, category := range pluginExport.Nodes {
		for _, node := range category {
			children = append(children, node.Children...)
		}
	}
	return children
}

func (p *Export) GetAllCategories(pluginExport *Export) map[string][]*Node {
	return pluginExport.Nodes
}

func (p *Export) GetCategory(name string, pluginExport *Export) []*Node {
	return pluginExport.Nodes[name]
}
