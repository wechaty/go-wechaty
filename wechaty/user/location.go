package user

import (
	"fmt"
	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
)

type Location struct {
	payload *schemas.LocationPayload
}

func NewLocation(payload *schemas.LocationPayload) *Location {
	return &Location{
		payload: payload,
	}
}

func (l *Location) Payload() schemas.LocationPayload {
	return *l.payload
}

func (l *Location) String() string {
	return fmt.Sprintf("Location<%s>", l.payload.Name)
}

func (l *Location) Address() string {
	return l.payload.Address
}

func (l *Location) Latitude() float64 {
	return l.payload.Latitude
}

func (l *Location) longitude() float64 {
	return l.payload.Longitude
}

func (l *Location) Name() string {
	return l.payload.Name
}

func (l *Location) Accuracy() float32 {
	return l.payload.Accuracy
}
