package puppethostie

import "os"

var (
	// WechatyPuppetHostieToken ...
	WechatyPuppetHostieToken string
	// WechatyPuppetHostieEndpoint ...
	WechatyPuppetHostieEndpoint string
)

func init() {
	WechatyPuppetHostieToken, _ = os.LookupEnv("WECHATY_PUPPET_HOSTIE_TOKEN")
	WechatyPuppetHostieEndpoint, _ = os.LookupEnv("WECHATY_PUPPET_HOSTIE_ENDPOINT")
}
