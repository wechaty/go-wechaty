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
	wechaty.OnHeartbeat(func(context *Context, data string) {
		got = data
	})
	wechaty.emit(schemas.PuppetEventNameHeartbeat, NewContext(), expect)
	if got != expect {
		log.Fatalf("got %s expect %s", got, expect)
	}
}

func TestWechaty_On(t *testing.T) {
	wechaty := NewWechaty()
	got := ""
	expect := "ding"
	wechaty.OnDong(func(context *Context, data string) {
		got = data
	})
	wechaty.emit(schemas.PuppetEventNameDong, NewContext(), expect)
	if got != expect {
		log.Fatalf("got %s expect %s", got, expect)
	}
}

func TestWechatyPluginDisableOnce(t *testing.T) {
	testMessage := "abc"
	received := false

	plugin := NewPlugin()
	plugin.OnHeartbeat(func(context *Context, data string) {
		if data == testMessage {
			received = true
		}
	})

	wechaty := NewWechaty()
	wechaty.OnHeartbeat(func(context *Context, data string) {
		if data == testMessage {
			context.DisableOnce(plugin)
		}
	})
	wechaty.Use(plugin)
	wechaty.emit(schemas.PuppetEventNameHeartbeat, NewContext(), testMessage)

	if received == true {
		t.Fatalf("disable plugin method failed.(Context.DisableOnce())")
	}
}

func TestWechatyPluginSetEnable(t *testing.T) {
	testMessage := "abc"
	received := false

	plugin := NewPlugin()
	plugin.OnHeartbeat(func(context *Context, data string) {
		if data == testMessage {
			received = true
		}
	})

	plugin.SetEnable(false)

	wechaty := NewWechaty()
	wechaty.Use(plugin)
	wechaty.emit(schemas.PuppetEventNameHeartbeat, NewContext(), testMessage)

	if received == true {
		t.Fatalf("disable plugin method failed.(Plugin.Disable())")
	}
}

func TestPluginPassingData(t *testing.T) {
	testData := "hello"

	p1 := NewPlugin()
	p1.OnHeartbeat(func(context *Context, data string) {
		context.SetData("helloStr", testData)
	})

	p2 := NewPlugin()
	p2.OnHeartbeat(func(context *Context, data string) {
		if testData != context.GetData("helloStr").(string) {
			t.Fatal("SetData() / GetData() not working.")
		}
	})

	wechaty := NewWechaty()
	wechaty.Use(p1).Use(p2)
	wechaty.emit(schemas.PuppetEventNameHeartbeat, NewContext(), "Data")
}

func TestPluginAbort(t *testing.T) {
	plugin := NewPlugin()
	plugin.OnHeartbeat(func(context *Context, data string) {
		t.Fatal("Context.Abort() not working.")
	})

	wechaty := NewWechaty()
	wechaty.OnHeartbeat(func(context *Context, data string) {
		context.Abort()
	})
	wechaty.Use(plugin)
	wechaty.emit(schemas.PuppetEventNameHeartbeat, NewContext(), "Data")
}
