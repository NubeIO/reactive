package reactive

import (
	"github.com/NubeIO/schema"
)

func (n *BaseNode) AddSchema() {}

func (n *BaseNode) GetSchema() *schema.Generated {
	return n.Schema
}
