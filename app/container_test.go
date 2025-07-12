package app

import (
	"reflect"
	"testing"
)

func TestNewAppContainer(t *testing.T) {
	container := NewAppContainer()

	// Use reflection to check for nil fields
	val := reflect.ValueOf(container).Elem()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldName := val.Type().Field(i).Name

		if field.IsNil() {
			t.Errorf("Field %s in AppContainer is nil after initialization", fieldName)
		}
	}
}
