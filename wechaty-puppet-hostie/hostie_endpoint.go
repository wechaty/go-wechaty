package puppethostie

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

var (
	// ErrNotToken token not found error
	ErrNotToken = errors.New("wechaty-puppet-hostie: token not found. See: <https://github.com/wechaty/wechaty-puppet-hostie#1-wechaty_puppet_hostie_token>")
)

// HostieEndPoint api.chatie.io endpoint api response
type HostieEndPoint struct {
	IP   string `json:"ip"`
	Port int    `json:"port,omitempty"`
}

// IsValid EndPoint is valid
func (p *HostieEndPoint) IsValid() bool {
	return p.IP == "" || p.IP == "0.0.0.0"
}

// Target Export IP+Port
func (p *HostieEndPoint) Target() string {
	port := p.Port
	if p.Port == 0 {
		port = 8788
	}
	return fmt.Sprintf("%s:%d", p.IP, port)
}

// discoverHostieEndPoint discover hostie endpoint ip and port
func (p *PuppetHostie) discoverHostieEndPoint() (endPoint HostieEndPoint, err error) {
	const hostieEndpoint = "https://api.chatie.io/v0/hosties/%s"

	if p.Token == "" {
		return endPoint, ErrNotToken
	}

	client := &http.Client{}
	if p.Timeout > 0 {
		client = &http.Client{
			Timeout: p.Timeout,
		}
	}

	resp, err := client.Get(fmt.Sprintf(hostieEndpoint, p.Token))
	if err != nil {
		return endPoint, fmt.Errorf("discoverHostieIP() err: %w", err)
	}

	if resp.StatusCode == 200 {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		var endPoint HostieEndPoint
		err = json.Unmarshal(body, &endPoint)
		if err != nil {
			return endPoint, fmt.Errorf("discoverHostieIP() err: %w", err)
		}
		return endPoint, nil
	}
	return endPoint, fmt.Errorf("discoverHostieIP() err: %w", err)
}
