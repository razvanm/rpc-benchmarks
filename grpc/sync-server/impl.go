package main

import (
	"golang.org/x/net/context"
	"github.com/razvanm/rpc-benchmarks/grpc"
	"io"
)

type impl struct {
}

func (f *impl) Sync(ctx context.Context, payload *sync.Payload) (*sync.Void, error) {
	return &sync.Void{}, nil
}

func (f *impl) SyncStream(stream sync.Sync_SyncStreamServer) error {
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