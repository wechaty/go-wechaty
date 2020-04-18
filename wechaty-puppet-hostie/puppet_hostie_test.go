package wechaty_puppet_hostie

import (
  "context"
  "fmt"
  "github.com/wechaty/go-wechaty/wechaty-puppet/events"
  "github.com/wechaty/go-wechaty/wechaty-puppet/option"
  "github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
  "testing"
)

func TestPuppetHostie_Start(t *testing.T) {
  event := events.New()
  event.On(schemas.PuppetEventNameScan, func(i ...interface{}) {
     fmt.Println(i)
  })
  NewPuppetHostie(option.WithToken("donut-test-user-3006"), option.WithEventEmitter(event)).Start(context.Background())
}
