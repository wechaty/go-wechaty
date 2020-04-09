package wechaty

import (
  "github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
  "log"
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

func TestWechaty_Emit(t *testing.T) {
  wechaty := NewWechaty()
  got := ""
  expect := "test"
  wechaty.OnHeartbeat(func(data string) {
    got = data
  })
  wechaty.Emit(schemas.PuppetEventNameHeartbeat, expect)
  if got != expect {
    log.Fatalf("got %s expect %s", got, expect)
  }
}

func TestWechaty_On(t *testing.T) {
  wechaty := NewWechaty()
  got := ""
  expect := "ding"
  wechaty.OnDong(func(data string) {
    got = data
  })
  wechaty.Emit(schemas.PuppetEventNameDong, expect)
  if got != expect {
    log.Fatalf("got %s expect %s", got, expect)
  }
}
