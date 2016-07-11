package main

import (
	"golang.org/x/net/context"
	"github.com/razvanm/rpc-benchmarks/grpc"
	"io"
)

type impl struct {
}

func (f *impl) Sink(ctx context.Context, payload *sink.Payload) (*sink.Void, error) {
	return &sink.Void{}, nil
}

func (f *impl) SinkStream(stream sink.Sink_SinkStreamServer) error {
	for {
		_, err := stream.Recv()

		if err == io.EOF {
			return nil
		}

		if err != nil {
			panic(err)
		}

		// Nothing to do beside iterating over the stream.
	}
}