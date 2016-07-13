package main

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/razvanm/rpc-benchmarks/vanadium-core"
	"v.io/v23"
	"v.io/v23/context"
	"v.io/x/ref/lib/security/securityflag"
	"v.io/x/ref/test/benchmark"
)

var (
	rootCtx *context.T
	stub    sink.SinkClientStub
)

func loop(b *testing.B, payload []byte) *benchmark.Stats {
	stats := benchmark.NewStats(16)
	b.ResetTimer()
	var err error
	for i := 0; i < b.N; i++ {
		b.StartTimer()
		start := time.Now()
		err = stub.Sink(rootCtx, payload)
		elapsed := time.Since(start)
		b.StopTimer()
		if err != nil {
			panic(err)
		}
		stats.Add(elapsed)
	}
	return stats
}

func Benchmark(b *testing.B) {
	payload := make([]byte, 0)
	loop(b, payload).Print(os.Stdout)
}

func TestMain(m *testing.M) {
	ctx, shutdown := v23.Init()
	defer shutdown()
	rootCtx = ctx

	auth := securityflag.NewAuthorizerOrDie()
	_, s, err := v23.WithNewServer(ctx, "", sink.SinkServer(&impl{}), auth)
	if err != nil {
		log.Panic("Error listening: ", err)
	}
	for _, endpoint := range s.Status().Endpoints {
		log.Printf("ENDPOINT=%s\n", endpoint.Name())
	}

	stub = sink.SinkClient(s.Status().Endpoints[0].Name())

	os.Exit(m.Run())
}
