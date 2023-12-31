package reactive

import "github.com/NubeIO/reactive/tracer"

func (n *BaseNode) AddTracer(stack *tracer.Tracer) {
	n.tracer = stack
}
func (n *BaseNode) Tracer() *tracer.Tracer {
	return n.tracer
}
