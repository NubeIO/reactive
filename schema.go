package reactive

import (
	"github.com/NubeIO/schema"
)

func (n *BaseNode) BuildSchema() {}

func (n *BaseNode) GetSchema() *schema.Generated {
	return n.Schema
}
