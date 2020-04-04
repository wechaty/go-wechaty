package wechaty

import (
	"reflect"
	"testing"
)

func TestNewWechaty(t *testing.T) {
	tests := []struct {
		name string
		want *Wechaty
	}{
		{name: "new", want: NewWechaty()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewWechaty(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewWechaty() = %v, want %v", got, tt.want)
			}
		})
	}
}
