package main

import (
	"log"

	"github.com/razvanm/rpc-benchmarks/vanadium-core"
	"v.io/v23"
	"v.io/x/ref/lib/security/securityflag"
	"v.io/x/ref/lib/signals"
)

func main() {
	ctx, shutdown := v23.Init()
	defer shutdown()

	auth := securityflag.NewAuthorizerOrDie()
	_, s, err := v23.WithNewServer(ctx, "", sync.SyncServer(&impl{}), auth)
	if err != nil {
		log.Panic("Error listening: ", err)
	}
	for _, endpoint := range s.Status().Endpoints {
		log.Printf("ENDPOINT=%s\n", endpoint.Name())
	}

	<-signals.ShutdownOnSignals(ctx) // Wait forever.
}
