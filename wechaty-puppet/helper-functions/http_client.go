package helper_functions

import (
  "net/http"
  "time"
)

var HttpClient = http.Client{
  Transport:     nil,
  CheckRedirect: nil,
  Jar:           nil,
  Timeout:       30 * time.Second,
}
