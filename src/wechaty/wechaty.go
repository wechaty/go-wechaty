package wechaty

// Wechaty
type Wechaty struct {
}

// NewWechaty
// instance by golang.
func NewWechaty() *Wechaty {
	return &Wechaty{}
}

func (w *Wechaty) OnScan(f func(qrCode, status string)) *Wechaty {
	return w
}

// OnLogin
// todo:: fake code. user should be struct
func (w *Wechaty) OnLogin(func(user string)) *Wechaty {
	return w
}

// OnMessage
// todo:: fake code. message should be struct
func (w *Wechaty) OnMessage(func(message string)) *Wechaty {
	return w
}

// Start
func (w *Wechaty) Start() *Wechaty {
	return w
}
