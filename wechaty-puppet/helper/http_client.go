package helper

import (
	"net/http"
	"time"
)

var HttpClient = http.Client{
	Transport:     nil,
	CheckRedirect: nil,
	Jar:           nil,
	Timeout:       5 * time.Second,
}
