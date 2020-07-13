package wechaty_puppet_padplus

import "os"

var (
	// WechatyPuppetToken ...
	WechatyPuppetToken string
	// WechatyPuppetEndpoint ...
	WechatyPuppetEndpoint string
)

const Endpoint = "padplus.juzibot.com:50051"

func init() {
	WechatyPuppetToken, _ = os.LookupEnv("WECHATY_PUPPET_HOSTIE_TOKEN")
	WechatyPuppetEndpoint, _ = os.LookupEnv("WECHATY_PUPPET_HOSTIE_ENDPOINT")
}
