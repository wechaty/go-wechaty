package helper

import "os"

func FileExists(path string) bool {
  _, err := os.Stat(path)
  if err != nil {
    if os.IsExist(err) {
      return true
    }
    return false
  }
  return true
}
