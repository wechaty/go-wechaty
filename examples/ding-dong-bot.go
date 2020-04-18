package main

import (
  "context"
  "fmt"
  "github.com/wechaty/go-wechaty/wechaty"
  wechaty_puppet_hostie "github.com/wechaty/go-wechaty/wechaty-puppet-hostie"
  "github.com/wechaty/go-wechaty/wechaty-puppet/events"
  "github.com/wechaty/go-wechaty/wechaty-puppet/option"
  "github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
  "github.com/wechaty/go-wechaty/wechaty/user"
  "log"
  "os"
  "os/signal"
)

func main() {
  event := events.New()
  event.On(schemas.PuppetEventNameScan, func(i ...interface{}) {
    payload := i[0].(*schemas.EventScanPayload)
    fmt.Println("https://api.qrserver.com/v1/create-qr-code/?data=" + payload.QrCode)
  })
  event.On(schemas.PuppetEventNameMessage, func(i ...interface{}) {
    fmt.Println(i[0].(*schemas.EventMessagePayload))
  })
  e := wechaty_puppet_hostie.NewPuppetHostie(option.WithToken("donut-test-user-3006"), option.WithEventEmitter(event)).Start(context.Background())
  if e != nil {
    panic(e)
  }
  select {

  }
  os.Exit(0)
  var bot = wechaty.NewWechaty()

  bot.
    OnScan(func(qrCode string, status schemas.ScanStatus, data string) {
      fmt.Printf("Scan QR Code to login: %v\nhttps://api.qrserver.com/v1/create-qr-code/?data=%s\n", status, qrCode)
    }).
    OnLogin(func(user string) {
      fmt.Printf("User %s logined\n", user)
    }).
    OnMessage(func(message *user.Message) {
      fmt.Println(fmt.Printf("Message: %v\n", message))
    })

  var err = bot.Start()
  if err != nil {
    panic(err)
  }

  var quitSig = make(chan os.Signal)
  signal.Notify(quitSig, os.Interrupt, os.Kill)

  select {
  case <-quitSig:
    log.Fatal("exit.by.signal")
  }
}
