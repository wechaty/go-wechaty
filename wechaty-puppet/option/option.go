package option

import (
  "github.com/wechaty/go-wechaty/wechaty-puppet/events"
  "github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
  "time"
)

// Option puppet option
type Option struct {
  Endpoint string
  Timeout  time.Duration
  Token    string
  EmitChan chan<- schemas.EmitStruct
  events.EventEmitter
}

// OptionFn func
type OptionFn func(opts *Option)

// WithEndpoint with Endpoint
func WithEndpoint(endpoint string) OptionFn {
  return func(opts *Option) {
    opts.Endpoint = endpoint
  }
}

// WithTimeout with Timeout
func WithTimeout(duration time.Duration) OptionFn {
  return func(opts *Option) {
    opts.Timeout = duration
  }
}

// WithToken with Token
func WithToken(token string) OptionFn {
  return func(opts *Option) {
    opts.Token = token
  }
}

// WithEmitChan with EmitChan
func WithEmitChan(emitChan chan<- schemas.EmitStruct) OptionFn {
  return func(opts *Option) {
    opts.EmitChan = emitChan
  }
}

func WithEventEmitter(eventEmitter events.EventEmitter) OptionFn {
  return func(opts *Option) {
    opts.EventEmitter = eventEmitter
  }
}
