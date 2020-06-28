package wechaty

import (
  wp "github.com/wechaty/go-wechaty/wechaty-puppet"
  mc "github.com/wechaty/go-wechaty/wechaty-puppet/memory-card"
)

// Option wechaty option
type Option struct {
  // wechaty name
  name string
  // puppet instance
  puppet wp.IPuppetAbstract
  // puppet option
  puppetOption wp.Option
  // io token
  ioToken string
  // memory card
  memoryCard mc.IMemoryCard
}

// OptionFn func
type OptionFn func(opts *Option)

// WithName with name
func WithName(name string) OptionFn {
  return func(opt *Option) {
    opt.name = name
  }
}

// WithPuppet with puppet impl
func WithPuppet(puppet wp.IPuppetAbstract) OptionFn {
  return func(opt *Option) {
    opt.puppet = puppet
  }
}

// WithPuppetOption with puppet option
func WithPuppetOption(puppetOption wp.Option) OptionFn {
  return func(opt *Option) {
    opt.puppetOption = puppetOption
  }
}

// WithIOToken with io token
func WithIOToken(ioToken string) OptionFn {
  return func(opt *Option) {
    opt.ioToken = ioToken
  }
}

// WithMemoryCard with memory card
func WithMemoryCard(memoryCard mc.IMemoryCard) OptionFn {
  return func(opt *Option) {
    opt.memoryCard = memoryCard
  }
}
