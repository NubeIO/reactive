package reactive

import (
	"github.com/NubeIO/schema"
)

func (n *BaseNode) BuildSchema() *schema.Generated {
	return n.schema
}
