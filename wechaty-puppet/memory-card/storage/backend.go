package storage

type IStorage interface {
  Save(payload map[string]interface{}) error
  Load() (map[string]interface{}, error)
  Destroy() error
}
