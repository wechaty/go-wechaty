package main

import (
  "fmt"
  "github.com/wechaty/go-wechaty/wechaty"
  "github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
  "github.com/wechaty/go-wechaty/wechaty/user"
)

func main() {
  err := wechaty.NewWechaty().
    SetToken("").
    OnScan(func(qrCode string, status schemas.ScanStatus, data string) {
      fmt.Printf("Scan QR Code to login: %v\nhttps://api.qrserver.com/v1/create-qr-code/?data=%s\n", status, qrCode)
    }).
    OnLogin(func(user string) {
      fmt.Printf("User %s logined\n", user)
    }).
    OnMessage(func(message *user.Message) {
      fmt.Println(fmt.Printf("Message: %v\n", message))
    }).
    Start()
  panic(err)
}
