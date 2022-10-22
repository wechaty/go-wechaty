package puppetservice

import "context"

// TLSInsecureServerCertCommonName Common Name:
//   - Server Name Indication (SNI)
//   - Case insensitive
// Wechaty Token format: `${SNI}/${UUIDv4}`
const TLSInsecureServerCertCommonName = "insecure"

// TLSCaCert Wechaty Certificate Authority Repo:
/* https://github.com/wechaty/dotenv/tree/main/ca

The SSL_ROOT_CERT is a root certificate generated by and for wechaty community.

Because it's the default root cert for the puppet service,
 so all the Polyglot Wechaty SDK should set this cert to be trusted by default.

Update:
 - Huan(202108): init, expired in 3650 days (after 2031/07)
*/
const TLSCaCert = `-----BEGIN CERTIFICATE-----
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
-----END CERTIFICATE-----`

type callCredToken struct {
	token string
}

func (r callCredToken) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": "Wechaty " + r.token,
	}, nil
}

func (r callCredToken) RequireTransportSecurity() bool {
	return true
}
