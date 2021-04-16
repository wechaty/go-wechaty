package puppetservice

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	// ErrNotToken token not found error
	ErrNotToken = errors.New("wechaty-puppet-service: token not found. See: <https://github.com/wechaty/wechaty-puppet-service#1-wechaty_puppet_service_token>")
)

// ServiceEndPoint api.chatie.io endpoint api response
type ServiceEndPoint struct {
	IP   string `json:"ip"`
	Port int    `json:"port,omitempty"`
}

// IsValid EndPoint is valid
func (p *ServiceEndPoint) IsValid() bool {
	return len(p.IP) > 0 && p.IP != "0.0.0.0"
}

// Target Export IP+Port
func (p *ServiceEndPoint) Target() string {
	port := p.Port
	if p.Port == 0 {
		port = 8788
	}
	return fmt.Sprintf("%s:%d", p.IP, port)
}

// discoverServiceEndPoint discover service endpoint ip and port
func (p *PuppetService) discoverServiceEndPoint() (endPoint ServiceEndPoint, err error) {
	const serviceEndpoint = "https://api.chatie.io/v0/hosties/%s"

	if p.Token == "" {
		return endPoint, ErrNotToken
	}

	client := &http.Client{}
	if p.Timeout > 0 {
		client = &http.Client{
			Timeout: p.Timeout,
		}
	}
	resp, err := client.Get(fmt.Sprintf(serviceEndpoint, p.Token))
	if err != nil {
		return endPoint, fmt.Errorf("discoverServiceEndPoint() err: %w", err)
	}

	if resp.StatusCode == http.StatusOK {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		err = json.Unmarshal(body, &endPoint)
		if err != nil {
			return endPoint, fmt.Errorf("discoverServiceEndPoint() err: %w", err)
		}
		return endPoint, nil
	}
	if resp.StatusCode == http.StatusNotFound {
		log.Printf("discoverServiceEndPoint() err: http.Status:%s\n", resp.Status)
		return endPoint, nil
	}
	return endPoint, fmt.Errorf("discoverServiceEndPoint() err: http.Status:%s", resp.Status)
}
