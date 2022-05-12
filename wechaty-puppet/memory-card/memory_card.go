package memory_card

import (
  "encoding/json"
  storage2 "github.com/wechaty/go-wechaty/wechaty-puppet/memory-card/storage"
  "sync"
)

// memory card interface
type IMemoryCard interface {
  GetInt64(key string) int64
  GetString(key string) string
  SetInt64(key string, value int64)
  Clear()
  Delete(key string)
  Has(key string) bool
  Load() error
  Save() error
  Destroy() error
  SetString(key string, value string)
  Set(key string, value interface{})
}

// TODO: 我将这个地方调整为 把storage的初始化放内部，原实现者可根据情况调整一下
func NewMemoryCard(name string) (IMemoryCard, error) {
  var storage, err = storage2.NewFileStorage(name)
  if err != nil {
    return nil, err
  }
  var memoryCard = &MemoryCard{payload: &sync.Map{}, storage: storage}
  return memoryCard, nil
}

// memory card
type MemoryCard struct {
  payload *sync.Map
  storage storage2.IStorage
}

func (mc *MemoryCard) GetInt64(key string) int64 {
  switch raw := mc.get(key).(type) {
  case json.Number:
    value, _ := raw.Int64()
    return value
  case int64:
    return raw
  default:
    return 0
  }
}

func (mc *MemoryCard) GetString(key string) string {
  value, ok := mc.get(key).(string)
  if ok {
    return value
  }
  return ""
}

func (mc *MemoryCard) get(key string) interface{} {
  v, _ := mc.payload.Load(key)
  return v
}

func (mc *MemoryCard) SetInt64(key string, value int64) {
  mc.Set(key, value)
}

func (mc *MemoryCard) SetString(key string, value string) {
  mc.Set(key, value)
}

func (mc *MemoryCard) Set(key string, value interface{}) {
  mc.payload.Store(key, value)
}

func (mc *MemoryCard) Clear() {
  mc.payload = &sync.Map{}
}

func (mc *MemoryCard) Delete(key string) {
  mc.payload.Delete(key)
}

func (mc *MemoryCard) Has(key string) bool {
  _, ok := mc.payload.Load(key)
  return ok
}

func (mc *MemoryCard) Load() error {
  raw, err := mc.storage.Load()
  if err != nil {
    return err
  }
  for k, v := range raw {
    mc.Set(k, v)
  }
  return nil
}

func (mc *MemoryCard) Save() error {
  raw := map[string]interface{}{}
  mc.payload.Range(func(key, value interface{}) bool {
    raw[key.(string)] = value
    return true
  })
  return mc.storage.Save(raw)
}

func (mc *MemoryCard) Destroy() error {
  mc.Clear()
  return mc.storage.Destroy()
}
