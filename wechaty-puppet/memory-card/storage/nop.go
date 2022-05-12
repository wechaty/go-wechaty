package storage

type NopStorage struct {
}

func (n NopStorage) Save(payload map[string]interface{}) error {
	return nil
}

func (n NopStorage) Load() (map[string]interface{}, error) {
	return map[string]interface{}{}, nil
}

func (n NopStorage) Destroy() {}
