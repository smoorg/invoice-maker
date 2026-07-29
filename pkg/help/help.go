package help

import (
	"reflect"

	"github.com/charmbracelet/bubbles/key"
)

func MapToBindingsList(k any) []key.Binding {
	v := reflect.ValueOf(k)

	values := make([]key.Binding, v.NumField())

	for _, field := range v.Fields() {
		if field.Type().AssignableTo(reflect.TypeFor[key.Binding]()) {
			if v, ok := field.Interface().(key.Binding); ok {
				values = append(values, v)
			}
		}
	}
	return values
}
