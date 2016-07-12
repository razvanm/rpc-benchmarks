package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/razvanm/rpc-benchmarks/vanadium-core"
	"v.io/v23"
	"v.io/v23/context"
	_ "v.io/x/ref/runtime/factories/roaming"
	"v.io/x/ref/test/benchmark"
)

var (
	duration = flag.Duration("duration", 10*time.Second, "Duration of the benchmark")
	server   = flag.String("server", "", "Name of the server to connect to")
	size     = flag.Uint("size", 0, "Size of the payload")
	stream   = flag.Bool("stream", false, "Use streaming RPCs")
	warmup   = flag.Duration("warmup", time.Second, "Duration of the warmup")

	rootCtx *context.T
	stub    sink.SinkClientStub
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
			err = stub.Sink(rootCtx, payload)
			elapsed := time.Since(start)
			if err != nil {
				panic(err)
			}
			stats.Add(elapsed)
		}
	}
}

func loopStream(duration time.Duration, payload []byte) *benchmark.Stats {
	stats := benchmark.NewStats(16)
	call, err := stub.SinkStream(rootCtx)
	if err != nil {
		panic(err)
	}
	stream := call.SendStream()
	end := time.After(duration)
	for {
		select {
		case <-end:
			return stats
		default:
			start := time.Now()
			err = stream.Send(payload)
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

	stub = sink.SinkClient(*server)

	payload := make([]byte, *size)
	if _, err := rand.Read(payload); err != nil {
		panic(err)
	}

	if *stream {
		fmt.Printf("Warming up for %s...\n", *warmup)
		loopStream(*warmup, payload)
		fmt.Printf("Benchmark params: %d bytes payload, %s duration, streaming\n", *size, *duration)
		loopStream(*duration, payload).Print(os.Stdout)
	} else {
		fmt.Printf("Warming up for %s...\n", *warmup)
		loop(*warmup, payload)
		fmt.Printf("Benchmark params: %d bytes payload, %s duration, no streaming\n", *size, *duration)
		loop(*duration, payload).Print(os.Stdout)
	}
}
