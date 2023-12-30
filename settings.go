package reactive

import (
	"errors"
	"fmt"
	"reflect"
)

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
	return n.settings
}

func (n *BaseNode) AddSettings(settings *Settings) {
	n.settings = settings
}

func (n *BaseNode) UpdateSettings(settings *Settings) {
	n.settings = settings
}

func (n *BaseNode) AddData(key string, data any) {
	n.data[key] = data
}

func (n *BaseNode) GetDataByKey(key string, out interface{}) error {
	data, exists := n.data[key]
	if !exists {
		return errors.New(fmt.Sprintf("failed to find by key: %s", key))
	}
	outValue := reflect.ValueOf(out).Elem()
	dataValue := reflect.ValueOf(data)
	fmt.Println(11111, outValue.Type(), dataValue.Type())
	if outValue.Type() != dataValue.Type() {
		return errors.New("type mismatch")
	}
	outValue.Set(dataValue)
	return nil
}

func (n *BaseNode) GetData() map[string]any {
	return n.data
}
