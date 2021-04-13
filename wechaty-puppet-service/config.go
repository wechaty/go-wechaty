package puppetservice

import "os"

var (
	// WechatyPuppetHostieToken ...
	// Deprecated: please use WechatyPuppetServiceToken TODO:will be deleted in the future
	WechatyPuppetHostieToken string

	// WechatyPuppetHostieEndpoint ...
	// Deprecated: please use WechatyPuppetHostieEndpoint TODO:will be deleted in the future
	WechatyPuppetHostieEndpoint string

	// WechatyPuppetServiceToken ...
	WechatyPuppetServiceToken string

	// WechatyPuppetServiceEndpoint ...
	WechatyPuppetServiceEndpoint string
)

func init() {
	WechatyPuppetHostieToken, _ = os.LookupEnv("WECHATY_PUPPET_HOSTIE_TOKEN")
	WechatyPuppetHostieEndpoint, _ = os.LookupEnv("WECHATY_PUPPET_HOSTIE_ENDPOINT")

	WechatyPuppetServiceToken, _ = os.LookupEnv("WECHATY_PUPPET_SERVICE_TOKEN")
	WechatyPuppetServiceEndpoint, _ = os.LookupEnv("WECHATY_PUPPET_SERVICE_ENDPOINT")
}
