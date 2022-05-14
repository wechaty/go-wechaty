// Deprecated:
package puppetservice

import (
	"errors"
	"fmt"
)

var (
	// ErrNotToken token not found error
	// Deprecated
	ErrNotToken = errors.New("wechaty-puppet-service: token not found. See: <https://github.com/wechaty/wechaty-puppet-service#1-wechaty_puppet_service_token>")
)

// ServiceEndPoint api.chatie.io endpoint api response
// Deprecated
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
