package reactive

type Position struct {
	PositionY int `json:"positionY"`
	PositionX int `json:"positionX"`
}

type Meta struct {
	Position Position `json:"position"`
}

func (n *BaseNode) GetMeta() *Meta {
	return n.Meta
}
