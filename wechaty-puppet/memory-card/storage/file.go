package storage

import (
	"encoding/json"
	helper_functions "github.com/wechaty/go-wechaty/wechaty-puppet/helper"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type FileStorage struct {
	absFileName string
}

func NewFileStorage(absFileName string) (*FileStorage, error) {
	absFileName, err := handleAbsFileName(absFileName)
	if err != nil {
		return nil, err
	}
	return &FileStorage{absFileName: absFileName}, nil
}

func (f *FileStorage) Save(payload map[string]interface{}) error {
	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(f.absFileName, jsonBytes, os.ModePerm)
}

func (f *FileStorage) Load() (map[string]interface{}, error) {
	if !helper_functions.FileExists(f.absFileName) {
		return map[string]interface{}{}, nil
	}
	file, err := os.Open(f.absFileName)
	if err != nil {
		return nil, err
	}
	result := map[string]interface{}{}
	decoder := json.NewDecoder(file)
	decoder.UseNumber()
	if err := decoder.Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

func (f *FileStorage) Destroy() error {
	return os.Remove(f.absFileName)
}

func handleAbsFileName(absFileName string) (string, error) {
	const suffix = ".memory-card.json"
	if !strings.HasSuffix(absFileName, suffix) {
		absFileName = absFileName + suffix
	}
	if !filepath.IsAbs(absFileName) {
		dir, err := os.Getwd()
		if err != nil {
			return "", err
		}
		absFileName = filepath.Join(dir, absFileName)
	}
	return absFileName, nil
}
