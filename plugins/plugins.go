package plugins

import (
	"errors"
)

type Node struct {
	ID       string  `json:"id"`
	Export   string  `json:"-"`
	Children []*Node `json:"children,omitempty"`
}

type Category struct {
	Name  string  `json:"name"`
	Nodes []*Node `json:"nodes,omitempty"`
}

type Export struct {
	Name        string      `json:"name"`
	Version     string      `json:"version"`
	Path        string      `json:"path"` // its file name
	Description string      `json:"description"`
	Categories  []*Category `json:"categories,omitempty"`
}

func NewPlugin(name, version, description string) *Export {
	return &Export{
		Name:        name,
		Version:     version,
		Description: description,
		Categories:  make([]*Category, 0),
	}
}

func (p *Export) AddCategory(categoryName string) {
	category := &Category{Name: categoryName, Nodes: make([]*Node, 0)}
	p.Categories = append(p.Categories, category)
}

func (p *Export) GetCategory(categoryName string) (*Category, error) {
	for _, category := range p.Categories {
		if category.Name == categoryName {
			return category, nil
		}
	}
	return nil, errors.New("category not found")
}

func (p *Export) AddNode(categoryName, nodeID, export string) error {
	category, err := p.GetCategory(categoryName)
	if err != nil {
		return err
	}

	node := &Node{ID: nodeID, Export: export}
	category.Nodes = append(category.Nodes, node)
	return nil
}

func (p *Export) AddChildNode(categoryName, parentID, childID, childExport string) error {
	category, err := p.GetCategory(categoryName)
	if err != nil {
		return err
	}

	parentNode := p.findNodeByID(category.Nodes, parentID)
	if parentNode == nil {
		return errors.New("parent node not found")
	}

	childNode := &Node{ID: childID, Export: childExport}
	parentNode.Children = append(parentNode.Children, childNode)
	return nil
}

func (p *Export) findNodeByID(nodes []*Node, id string) *Node {
	for _, node := range nodes {
		if node.ID == id {
			return node
		}
		if foundNode := p.findNodeByID(node.Children, id); foundNode != nil {
			return foundNode
		}
	}
	return nil
}

func (p *Export) GetAllNodes() []*Node {
	var allNodes []*Node

	for _, category := range p.Categories {
		allNodes = append(allNodes, p.getAllNodesInCategory(category.Nodes)...)
	}

	return allNodes
}

func (p *Export) getAllNodesInCategory(nodes []*Node) []*Node {
	var allNodes []*Node

	for _, node := range nodes {
		allNodes = append(allNodes, node)
		allNodes = append(allNodes, p.getAllNodesInCategory(node.Children)...)
	}

	return allNodes
}

func (p *Export) GelNodes(export *Export) []*Node {
	var allNodes []*Node

	for _, category := range export.Categories {
		allNodes = append(allNodes, p.getAllNodesInCategory(category.Nodes)...)
	}

	return allNodes
}
