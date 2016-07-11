package main

import (
	"flag"
	"fmt"
	"net"

	"github.com/razvanm/rpc-benchmarks/grpc"
	"google.golang.org/grpc"
)

var (
	address = flag.String("address", ":0", "What local address to listen")
)

func main() {
	flag.Parse()
	listener, err := net.Listen("tcp", *address)
	fmt.Printf("Listening on %v", listener.Addr())
	if err != nil {
		panic(err)
	}
	s := grpc.NewServer()
	sync.RegisterSyncServer(s, &impl{})
	s.Serve(listener)
}
