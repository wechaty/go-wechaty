package memory_card

import (
  "crypto/md5"
  "fmt"
  storage2 "github.com/wechaty/go-wechaty/wechaty-puppet/memory-card/storage"
  "io"
  "math/rand"
  "testing"
  "time"
)

func TestMemoryCard_GetInt(t *testing.T) {
  mc := newMemoryCard(t)

  t.Run("not value return 0", func(t *testing.T) {
    got := mc.GetInt64("not")
    if got != 0 {
      t.Fatalf("got %d expect 0", got)
    }
  })

  t.Run("set string return 0", func(t *testing.T) {
    key := "set_string_return_0"
    mc.SetString(key, "string")
    got := mc.GetInt64(key)
    if got != 0 {
      t.Fatalf("got %d expect 0", got)
    }
  })

  t.Run("success", func(t *testing.T) {
    key := "success"
    expect := rand.Int63()
    mc.SetInt64(key, expect)
    got := mc.GetInt64(key)
    if got != expect {
      t.Fatalf("got %d expect %d", got, expect)
    }
  })
}

func TestMemoryCard_GetString(t *testing.T) {
  mc := newMemoryCard(t)
  t.Run("not value return empty string", func(t *testing.T) {
    got := mc.GetString("not")
    if got != "" {
      t.Fatalf("got %s expect empty stirng", got)
    }
  })

  t.Run("set int return empty string", func(t *testing.T) {
    key := "set_int_return_empty_string"
    mc.SetInt64(key, 1)
    got := mc.GetString(key)
    if got != "" {
      t.Fatalf("got %s expect empty stirng", got)
    }
  })

  t.Run("success", func(t *testing.T) {
    key := "success"
    expect := randString()
    mc.SetString(key, expect)
    got := mc.GetString(key)
    if got != expect {
      t.Fatalf("got %s expect %s", got, expect)
    }
  })
}

func TestMemoryCard_SetInt(t *testing.T) {
  mc := newMemoryCard(t)
  key := "success"
  expect := rand.Int63()
  mc.SetInt64(key, expect)
  got := mc.GetInt64(key)
  if got != expect {
    t.Fatalf("got %d expect %d", got, expect)
  }
}

func TestMemoryCard_SetString(t *testing.T) {
  mc := newMemoryCard(t)
  key := "success"
  expect := randString()
  mc.SetString(key, expect)
  got := mc.GetString(key)
  if got != expect {
    t.Fatalf("got %s expect %s", got, expect)
  }
}

func TestMemoryCard_Clear(t *testing.T) {
  mc := newMemoryCard(t)
  table := []struct {
    key   string
    value int64
  }{
    {"key1", 1},
    {"key2", 2},
  }
  for _, v := range table {
    mc.SetInt64(v.key, v.value)
  }
  mc.Clear()
  for _, v := range table {
    got := mc.GetInt64(v.key)
    if got != 0 {
      t.Fatalf("got %d expect 0", got)
    }
  }
}

func TestMemoryCard_Delete(t *testing.T) {
  mc := newMemoryCard(t)
  table := []struct {
    key   string
    value int64
  }{
    {"key1", 1},
    {"key2", 2},
  }
  for _, v := range table {
    mc.SetInt64(v.key, v.value)
  }
  mc.Delete(table[0].key)
  got := mc.GetInt64(table[0].key)
  if got != 0 {
    t.Fatalf("got %d expect 0", got)
  }
  got = mc.GetInt64(table[1].key)
  if got != table[1].value {
    t.Fatalf("got %d expect d", table[1].value)
  }
}

func TestMemoryCard_Has(t *testing.T) {
  mc := newMemoryCard(t)
  if false != mc.Has("null") {
    t.Fatalf("got true expect false")
  }
  mc.SetString("key", "value")
  if true != mc.Has("key") {
    t.Fatalf("got false expect false")
  }
}

func TestMemoryCard_SaveAndLoad(t *testing.T) {
  table := []struct {
    key   string
    value int64
  }{
    {"key0", 0},
    {"key1", 1},
    {"key2", 2},
  }
  mc1 := newMemoryCard(t)
  for _, v := range table {
    mc1.set(v.key, v.value)
  }
  err := mc1.Save()
  if err != nil {
    t.Fatal(err.Error())
  }
  mc2 := newMemoryCard(t)
  err = mc2.Load()
  if err != nil {
    t.Fatalf(err.Error())
  }
  for _, v := range table {
    got := mc2.GetInt64(v.key)
    if got != v.value {
      t.Fatalf("got %d expected %d", got, v.value)
    }
  }
}

func randString() string {
  t := time.Now()
  h := md5.New()
  _, _ = io.WriteString(h, t.String())
  return fmt.Sprintf("%x", h.Sum(nil))
}

func newMemoryCard(t *testing.T) *MemoryCard {
  storage, err := storage2.NewFileStorage("testdata/file")
  if err != nil {
    t.Fatalf(err.Error())
  }
  return NewMemoryCard(storage)
}
