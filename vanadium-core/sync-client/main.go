package main

import (
	"flag"
	"time"

	"fmt"
	"github.com/razvanm/rpc-benchmarks/vanadium-core"
	"v.io/v23"
	"v.io/v23/context"
	_ "v.io/x/ref/runtime/factories/generic"
)

var (
	server = flag.String("server", "", "Name of the server to connect to")
)

func main() {
	ctx, shutdown := v23.Init()
	defer shutdown()

	f := sync.SyncClient(*server)

	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	payload := make([]byte, 0)
	err := f.Sync(ctx, payload)
	fmt.Printf("err: %v\n", err)
}
