package reactive

type Position struct {
	PositionY int `json:"positionY"`
	PositionX int `json:"positionX"`
}

type Meta struct {
	Position   Position `json:"position"`
	ParentUUID string   `json:"parentUUID"`
}

func (n *BaseNode) setMeta(opts *Options) {
	if opts != nil {
		n.meta = opts.Meta
		if n.meta != nil {
			n.parentUUID = n.meta.ParentUUID
		}
	}
}

func (n *BaseNode) GetMeta() *Meta {
	return n.meta
}
