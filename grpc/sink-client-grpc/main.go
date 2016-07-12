package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	"crypto/tls"
	"crypto/x509"
	"github.com/razvanm/rpc-benchmarks/grpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"v.io/x/ref/test/benchmark"
)

var (
	duration = flag.Duration("duration", 10*time.Second, "Duration of the benchmark")
	server   = flag.String("server", "", "Name of the server to connect to")
	size     = flag.Uint("size", 0, "Size of the payload")
	stream   = flag.Bool("stream", false, "Use streaming RPCs")
	warmup   = flag.Duration("warmup", time.Second, "Duration of the warmup")
	caFile   = flag.String("ca", "certs/ca.pem", "TLS CA file")
	certFile = flag.String("cert", "certs/server.pem", "TLS cert file")
	keyFile  = flag.String("key", "certs/server.key", "TLS key file")

	client sink.SinkClient
)

func loop(duration time.Duration, payload *sink.Payload) *benchmark.Stats {
	stats := benchmark.NewStats(16)
	end := time.After(duration)
	var err error
	for {
		select {
		case <-end:
			return stats
		default:
			start := time.Now()
			_, err = client.Sink(context.Background(), payload)
			elapsed := time.Since(start)
			if err != nil {
				panic(err)
			}
			stats.Add(elapsed)
		}
	}
}

func loopStream(duration time.Duration, payload *sink.Payload) *benchmark.Stats {
	stats := benchmark.NewStats(16)
	stream, err := client.SinkStream(context.Background())
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

// transportCredentials is a combination of credentials.NewClientTLSFromFile and
// credentials.NewServerTLSFromFile.
func transportCredentials(caCertFile, caName, certFile, keyFile string) credentials.TransportCredentials {
	b, err := ioutil.ReadFile(caCertFile)
	if err != nil {
		panic(err)
	}
	cp := x509.NewCertPool()
	if !cp.AppendCertsFromPEM(b) {
		panic(err)
	}
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		panic(err)
	}
	return credentials.NewTLS(&tls.Config{
		ServerName:   caName,
		RootCAs:      cp,
		Certificates: []tls.Certificate{cert},
	})
}

func main() {
	flag.Parse()
	creds := transportCredentials(*caFile, "server", *certFile, *keyFile)
	opts := grpc.WithTransportCredentials(creds)
	conn, err := grpc.Dial(*server, opts)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client = sink.NewSinkClient(conn)

	b := make([]byte, *size)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	payload := &sink.Payload{Payload: b}

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
