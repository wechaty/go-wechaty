package file_box

import (
  helper_functions "github.com/wechaty/go-wechaty/wechaty-puppet/helper-functions"
  "io/ioutil"
  "net/http"
)

type fileBoxUrl struct {
  remoteUrl string
  headers   http.Header
}

func NewFileBoxUrl(remoteUrl string, headers http.Header) *fileBoxUrl {
  return &fileBoxUrl{remoteUrl: remoteUrl, headers: headers}
}

func (fb *fileBoxUrl) toJSONMap() map[string]interface{} {
  return map[string]interface{}{
    "headers":   fb.headers,
    "remoteUrl": fb.remoteUrl,
  }
}

func (fb *fileBoxUrl) toBytes() ([]byte, error) {
  request, err := http.NewRequest(http.MethodGet, fb.remoteUrl, nil)
  if err != nil {
    return nil, err
  }
  request.Header = fb.headers
  response, err := helper_functions.HttpClient.Do(request)
  if err != nil {
    return nil, err
  }
  defer response.Body.Close()
  return ioutil.ReadAll(response.Body)
}
