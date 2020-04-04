package memory_card

import (
  "encoding/json"
  storage2 "github.com/wechaty/go-wechaty/wechaty-puppet/memory-card/storage"
  "sync"
)

type MemoryCard struct {
  payload *sync.Map
  storage storage2.IStorage
}

func NewMemoryCard(storage storage2.IStorage) *MemoryCard {
  return &MemoryCard{
    payload: &sync.Map{},
    storage: storage,
  }
}

func (mc *MemoryCard) GetInt64(key string) int64 {
  raw := mc.get(key)
  switch raw.(type) {
  case json.Number:
    value, _ := raw.(json.Number).Int64()
    return value
  case int64:
    return raw.(int64)
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
  mc.set(key, value)
}

func (mc *MemoryCard) SetString(key string, value string) {
  mc.set(key, value)
}

func (mc *MemoryCard) set(key string, value interface{}) {
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
    mc.set(k, v)
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
