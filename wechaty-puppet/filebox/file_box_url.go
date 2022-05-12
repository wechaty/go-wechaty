package filebox

import (
	"bytes"
	"fmt"
	"github.com/wechaty/go-wechaty/wechaty-puppet/helper"
	"io"
	"io/ioutil"
	"net/http"
)

type fileBoxUrl struct {
	remoteUrl string
	headers   http.Header
}

func newFileBoxUrl(remoteUrl string, headers http.Header) *fileBoxUrl {
	return &fileBoxUrl{remoteUrl: remoteUrl, headers: headers}
}

func (fb *fileBoxUrl) toJSONMap() (map[string]interface{}, error) {
	if fb.remoteUrl == "" {
		return nil, fmt.Errorf("fileBoxUrl.toJSONMap %w", ErrNoUrl)
	}
	return map[string]interface{}{
		"headers": fb.headers,
		"url":     fb.remoteUrl,
	}, nil
}

func (fb *fileBoxUrl) toBytes() ([]byte, error) { //nolint:unused
	request, err := http.NewRequest(http.MethodGet, fb.remoteUrl, nil)
	if err != nil {
		return nil, err
	}
	request.Header = fb.headers
	response, err := helper.HttpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	return ioutil.ReadAll(response.Body)
}

func (fb *fileBoxUrl) toReader() (io.Reader, error) {
	request, err := http.NewRequest(http.MethodGet, fb.remoteUrl, nil)
	if err != nil {
		return nil, err
	}
	request.Header = fb.headers
	response, err := helper.HttpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	all, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(all), nil
}
