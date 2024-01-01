package reactive

import "github.com/NubeIO/reactive/tracer"

func (n *BaseNode) GetTracer() *tracer.Tracer {
	return n.tracer
}

func (n *BaseNode) SetTracer(key string) *tracer.Tracer {
	n.tracer.TracerKey(key)
	return n.tracer
}

func (n *BaseNode) InitTracer(t *tracer.Tracer) {
	n.tracer = t
	err := n.tracer.AddTracer(n.GetUUID(), "")
	if err != nil {
		n.GetLogger().Errorf("error on setup node tracer err: %s", err.Error())
		return
	}
}
