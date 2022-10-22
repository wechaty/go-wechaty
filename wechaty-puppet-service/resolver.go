package puppetservice

import (
	"encoding/json"
	"fmt"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/resolver/manual"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func wechatyResolver() resolver.Builder {
	r := manual.NewBuilderWithScheme("wechaty")
	r.BuildCallback = resolverBuildCallBack
	return r
}

func resolverBuildCallBack(target resolver.Target, conn resolver.ClientConn, options resolver.BuildOptions) {
	// target.URL.Host `api.chatie.io` in `wechaty://api.chatie.io/__token__`
	// target.URL.Path `__token__` in `wechaty://api.chatie.io/__token__`
	log.Println("resolverBuildCallBack()")
	uri := fmt.Sprintf("https://%s/v0/hosties%s", target.URL.Host, target.URL.Path)
	address, err := discoverAPI(uri)
	if err != nil {
		conn.ReportError(err)
		return
	}
	if address == nil || address.Host == "" {
		conn.ReportError(fmt.Errorf(`token %s does not exist`, strings.TrimLeft(target.URL.Path, "/")))
		return
	}
	err = conn.UpdateState(resolver.State{
		Addresses: []resolver.Address{{
			Addr: fmt.Sprintf("%s:%d", address.Host, address.Port),
		}},
	})
	if err != nil {
		log.Println("resolverBuildCallBack UpdateState err: ", err.Error())
		return
	}
}

type serviceAddress struct {
	Host string
	Port int
}

func discoverAPI(uri string) (*serviceAddress, error) {
	response, err := http.Get(uri)
	if err != nil {
		return nil, fmt.Errorf("discoverAPI http.Get() %w", err)
	}
	defer response.Body.Close()

	// 4xx
	if response.StatusCode >= 400 && response.StatusCode < 500 {
		return nil, nil
	}

	// 2xx
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, fmt.Errorf("discoverAPI http.Get() status:%s %w", response.Status, err)
	}
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("discoverAPI ioutil.ReadAll %w", err)
	}

	r := &serviceAddress{}
	if err := json.Unmarshal(data, r); err != nil {
		return nil, fmt.Errorf("discoverAPI json.Unmarshal %w", err)
	}
	return r, nil
}
