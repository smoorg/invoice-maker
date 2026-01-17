package help

import (
	"reflect"

	"github.com/charmbracelet/bubbles/key"
)

func MapToBindingsList(k any) []key.Binding {
	v := reflect.ValueOf(k)

	values := make([]key.Binding, v.NumField())

	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).Type().AssignableTo(reflect.TypeOf(key.Binding{})) {
			if v, ok := v.Field(i).Interface().(key.Binding); ok {
				values = append(values, v)
			}
		}
	}
	return values
}
