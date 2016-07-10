package main

import (
	"flag"
	"time"

	"github.com/razvanm/rpc-benchmarks/vanadium-core"
	"math/rand"
	"v.io/v23"
	_ "v.io/x/ref/runtime/factories/generic"
	"v.io/v23/context"
	"fmt"
	"v.io/x/ref/test/benchmark"
	"os"
)

var (
	server   = flag.String("server", "", "Name of the server to connect to")
	warmup   = flag.Duration("warmup", time.Second, "Duration of the warmup")
	duration = flag.Duration("duration", 10*time.Second, "Duration of the benchmark")
	size     = flag.Uint("size", 0, "Size of the payload")

	stub sync.SyncClientStub
	rootCtx *context.T
)

func loop(duration time.Duration, payload []byte) *benchmark.Stats {
	stats := benchmark.NewStats(16)
	end := time.After(duration)
	var err error
	for {
		select {
		case <-end:
			return stats
		default:
			start := time.Now()
			err = stub.Sync(rootCtx, payload)
			elapsed := time.Since(start)
			if err != nil {
				panic(err)
			}
			stats.Add(elapsed)
		}
	}
}

func main() {
	ctx, shutdown := v23.Init()
	defer shutdown()
	rootCtx = ctx

	stub = sync.SyncClient(*server)

	payload := make([]byte, *size)
	if _, err := rand.Read(payload); err != nil {
		panic(err)
	}

	fmt.Printf("Warming up for %s...\n", *warmup)
	loop(*warmup, payload)
	fmt.Printf("Running the benchmark for %s...\n", *duration)
	loop(*duration, payload).Print(os.Stdout)
}
