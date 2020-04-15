package main

import (
  "fmt"
  "github.com/wechaty/go-wechaty/wechaty"
  "github.com/wechaty/go-wechaty/wechaty-puppet/errorx"
  "github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
  "github.com/wechaty/go-wechaty/wechaty/user"
  "log"
  "os"
  "os/signal"
)

func main() {
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

  var errChan = errorx.NewChanErr(1)

  var err = bot.Start()
  if err != nil {
    errChan.Put(err)
  }

  var quitSig = make(chan os.Signal)
  signal.Notify(quitSig, os.Interrupt, os.Kill)

  select {
  case <-quitSig:
    log.Fatal("exit.by.signal")
  case threadErr := <-errChan.WaitErr():
    log.Fatalf("exit.by.err: %+v", threadErr)
  }
}
