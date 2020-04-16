package wechaty

import (
  wp "github.com/wechaty/go-wechaty/wechaty-puppet"
  mc "github.com/wechaty/go-wechaty/wechaty-puppet/memory-card"
)

// Wechaty option
type Option struct {
  // wechaty name
  name string
  // puppet instance
  puppet wp.PuppetInterface
  // puppet option
  puppetOption *wp.Option
  // io token
  ioToken string
  // memory card
  memoryCard mc.IMemoryCard
}

// Option func
type OptionFn func(opts *Option)

// with name
func WithName(name string) OptionFn {
  return func(opt *Option) {
    opt.name = name
  }
}

// with puppet impl
func WithPuppet(puppet wp.PuppetInterface) OptionFn {
  return func(opt *Option) {
    opt.puppet = puppet
  }
}

// with puppet option
func WithPuppetOption(puppetOption *wp.Option) OptionFn {
  return func(opt *Option) {
    opt.puppetOption = puppetOption
  }
}

// with io token
func WithIOToken(ioToken string) OptionFn {
  return func(opt *Option) {
    opt.ioToken = ioToken
  }
}

// with memory card
func WithMemoryCard(memoryCard mc.IMemoryCard) OptionFn {
  return func(opt *Option) {
    opt.memoryCard = memoryCard
  }
}
