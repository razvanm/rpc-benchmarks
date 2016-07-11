package main

import (
	"flag"
	"fmt"
	"net"

	"github.com/razvanm/rpc-benchmarks/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	address  = flag.String("address", ":0", "Local address to listen on")
	certFile = flag.String("cert", "certs/server.pem", "TLS cert file")
	keyFile  = flag.String("key", "certs/server.key", "TLS key file")
)

func main() {
	flag.Parse()

	listener, err := net.Listen("tcp", *address)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Listening on %v\n", listener.Addr())

	creds, err := credentials.NewServerTLSFromFile(*certFile, *keyFile)
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer(grpc.Creds(creds))
	sink.RegisterSinkServer(s, &impl{})
	s.Serve(listener)
}
