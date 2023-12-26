package reactive

// ---------------------------- NODE SETTINGS -------------------------- //

type Settings struct {
	Value interface{} `json:"value"`
}

func (s *Settings) GetFloat64Value() float64 {
	if s == nil {
		return 0
	}
	if floatValue, ok := s.Value.(float64); ok {
		return floatValue
	}
	return 0
}

func (s *Settings) GetFloat64ValuePointer() *float64 {
	if s == nil {
		return nil
	}
	if floatValue, ok := s.Value.(float64); ok {
		return &floatValue
	}
	return nil
}

func (n *BaseNode) GetSettings() *Settings {
	return n.Settings
}

func (n *BaseNode) AddSettings(settings *Settings) {
	n.Settings = settings
}

func (n *BaseNode) UpdateSettings(settings *Settings) {
	n.Settings = settings
}
