package puppetservice

import (
	wechatypuppet "github.com/wechaty/go-wechaty/wechaty-puppet"
	"time"
)

// TlsConfig tls config
type TlsConfig struct {
	CaCert     string
	ServerName string

	Disable bool //  only for compatible with old clients/servers
}

// Options puppet-service options
type Options struct {
	wechatypuppet.Option

	GrpcReconnectInterval time.Duration
	Authority             string
	Tls                   TlsConfig
}
