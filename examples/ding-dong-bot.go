package main

import (
	"fmt"

	"github.com/wechaty/go-wechaty/wechaty"
)

func main() {
	_ = wechaty.NewWechaty().
		OnScan(func(qrCode, status string) {
			fmt.Printf("Scan QR Code to login: %s\nhttps://api.qrserver.com/v1/create-qr-code/?data=%s\n", status, qrCode)
		}).
		OnLogin(func(user string) { fmt.Printf("User %s logined\n", user) }).
		OnMessage(func(message string) { fmt.Printf("Message: %s\n", message) }).
		Start()
}
