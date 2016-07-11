package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/razvanm/rpc-benchmarks/grpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"v.io/x/ref/test/benchmark"
	"google.golang.org/grpc/credentials"
)

var (
	duration = flag.Duration("duration", 10*time.Second, "Duration of the benchmark")
	server   = flag.String("server", "", "Name of the server to connect to")
	size     = flag.Uint("size", 0, "Size of the payload")
	stream   = flag.Bool("stream", false, "Use streaming RPCs")
	warmup   = flag.Duration("warmup", time.Second, "Duration of the warmup")
	caFile = flag.String("ca", "certs/ca.pem", "TLS CA file")

	client sync.SyncClient
)

func loop(duration time.Duration, payload *sync.Payload) *benchmark.Stats {
	stats := benchmark.NewStats(16)
	end := time.After(duration)
	var err error
	for {
		select {
		case <-end:
			return stats
		default:
			start := time.Now()
			_, err = client.Sync(context.Background(), payload)
			elapsed := time.Since(start)
			if err != nil {
				panic(err)
			}
			stats.Add(elapsed)
		}
	}
}

func loopStream(duration time.Duration, payload *sync.Payload) *benchmark.Stats {
	stats := benchmark.NewStats(16)
	stream, err := client.SyncStream(context.Background())
	if err != nil {
		panic(err)
	}
	end := time.After(duration)
	for {
		select {
		case <-end:
			stream.CloseSend()
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
	flag.Parse()
	creds, err := credentials.NewClientTLSFromFile(*caFile, "server")
	if err != nil {
		panic(err)
	}
	opts := grpc.WithTransportCredentials(creds)
	conn, err := grpc.Dial(*server, opts)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client = sync.NewSyncClient(conn)

	b := make([]byte, *size)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	payload := &sync.Payload{Payload: b}

	if *stream {
		fmt.Printf("Warming up for %s...\n", *warmup)
		loopStream(*warmup, payload)
		fmt.Printf("Running the benchmark for %s...\n", *duration)
		loopStream(*duration, payload).Print(os.Stdout)
	} else {
		fmt.Printf("Warming up for %s...\n", *warmup)
		loop(*warmup, payload)
		fmt.Printf("Running the benchmark for %s...\n", *duration)
		loop(*duration, payload).Print(os.Stdout)
	}
}
