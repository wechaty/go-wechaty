package wechaty

import (
  "github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
  "log"
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
      _ = NewWechaty()
    })
  }
}

func TestWechaty_Emit(t *testing.T) {
  wechaty := NewWechaty()
  got := ""
  expect := "test"
  wechaty.OnHeartbeat(func(context *PluginContext, data string) {
    got = data
  })
  wechaty.emit(schemas.PuppetEventNameHeartbeat, NewPluginContext(), expect)
  if got != expect {
    log.Fatalf("got %s expect %s", got, expect)
  }
}

func TestWechaty_On(t *testing.T) {
  wechaty := NewWechaty()
  got := ""
  expect := "ding"
  wechaty.OnDong(func(context *PluginContext, data string) {
    got = data
  })
  wechaty.emit(schemas.PuppetEventNameDong, NewPluginContext(), expect)
  if got != expect {
    log.Fatalf("got %s expect %s", got, expect)
  }
}
