package puppetservice

import (
	"errors"
	"os"
)

// ErrTokenNotFound err token not found
var ErrTokenNotFound = errors.New("wechaty-puppet-service: WECHATY_PUPPET_SERVICE_TOKEN not found")

func envServiceToken(token string) (string, error) {
	if token != "" {
		return token, nil
	}

	token = os.Getenv("WECHATY_PUPPET_SERVICE_TOKEN")
	if token != "" {
		return token, nil
	}

	token = os.Getenv("WECHATY_PUPPET_HOSTIE_TOKEN")
	if token != "" {
		log.Trace("WECHATY_PUPPET_HOSTIE_TOKEN has been deprecated," +
			"please use WECHATY_PUPPET_SERVICE_TOKEN instead.")
		return token, nil
	}

	return "", ErrTokenNotFound
}

func envEndpoint(endpoint string) string {
	if endpoint != "" {
		return endpoint
	}

	endpoint = os.Getenv("WECHATY_PUPPET_SERVICE_ENDPOINT")
	if endpoint != "" {
		return endpoint
	}

	endpoint = os.Getenv("WECHATY_PUPPET_HOSTIE_ENDPOINT")
	if endpoint != "" {
		log.Println("WECHATY_PUPPET_HOSTIE_ENDPOINT has been deprecated," +
			"please use WECHATY_PUPPET_SERVICE_ENDPOINT instead.")
		return endpoint
	}
	return ""
}

func envAuthority(authority string) string {
	if authority != "" {
		return authority
	}

	authority = os.Getenv("WECHATY_PUPPET_SERVICE_AUTHORITY")
	if authority != "" {
		return authority
	}

	return "api.chatie.io"
}

func envNoTLSInsecureClient(disable bool) bool {
	return disable || os.Getenv("WECHATY_PUPPET_SERVICE_NO_TLS_INSECURE_CLIENT") == "true"
}

func envTLSServerName(serverName string) string {
	if serverName != "" {
		return serverName
	}

	return os.Getenv("WECHATY_PUPPET_SERVICE_TLS_SERVER_NAME")
}

func envTLSCaCert(caCert string) string {
	if caCert != "" {
		return caCert
	}
	caCert = os.Getenv("WECHATY_PUPPET_SERVICE_TLS_CA_CERT")
	if caCert != "" {
		return caCert
	}
	return TLSCaCert
}
