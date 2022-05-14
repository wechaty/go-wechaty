package puppetservice

import (
	pbwechaty "github.com/wechaty/go-grpc/wechaty"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

func (p *PuppetService) startGrpcClient() error {
	var err error
	var creds credentials.TransportCredentials
	var callOptions []grpc.CallOption
	if p.disableTLS {
		log.Println("PuppetService.startGrpcClient TLS: disabled (INSECURE)")
		creds = insecure.NewCredentials()
	} else {
		callOptions = append(callOptions, grpc.PerRPCCredentials(callCredToken{token: p.token}))
		creds, err = p.createCred()
		if err != nil {
			return err
		}
	}

	dialOptions := []grpc.DialOption{
		grpc.WithTransportCredentials(creds),
		grpc.WithDefaultCallOptions(callOptions...),
		grpc.WithResolvers(wechatyResolver()),
	}

	if p.disableTLS {
		// Deprecated: this block will be removed after Dec 21, 2022.
		dialOptions = append(dialOptions, grpc.WithAuthority(p.token))
	}

	conn, err := grpc.Dial(p.endpoint, dialOptions...)
	if err != nil {
		return err
	}
	p.grpcConn = conn

	go p.autoReconnectGrpcConn()

	p.grpcClient = pbwechaty.NewPuppetClient(conn)
	return nil
}

func (p *PuppetService) autoReconnectGrpcConn() {
	<-p.started
	interval := 2 * time.Second
	if p.opts.GrpcReconnectInterval > 0 {
		interval = p.opts.GrpcReconnectInterval
	}
	ticker := time.NewTicker(interval)
	isClose := false
	for {
		select {
		case <-ticker.C:
			connState := p.grpcConn.GetState()
			// 重新连接成功
			if isClose && connectivity.Ready == connState {
				isClose = false
				log.Printf("PuppetService.autoReconnectGrpcConn grpc reconnection successful")
				if err := p.startGrpcStream(); err != nil {
					log.Printf("PuppetService.autoReconnectGrpcConn startGrpcStream err:%s", err.Error())
				}
			}

			if p.grpcConn.GetState() == connectivity.Idle {
				isClose = true
				p.grpcConn.Connect()
				log.Printf("PuppetService.autoReconnectGrpcConn grpc reconnection...")
			}
		case <-p.stop:
			return
		}
	}
}
