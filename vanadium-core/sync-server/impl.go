package main

import (
	"fmt"

	"github.com/razvanm/rpc-benchmarks/vanadium-core"
	"v.io/v23/context"
	"v.io/v23/rpc"
	_ "v.io/x/ref/runtime/factories/generic"
)

type impl struct {
}

func (f *impl) Sync(_ *context.T, _ rpc.ServerCall, payload []byte) error {
	fmt.Printf("Sync: %d bytes\n", len(payload))
	return nil
}

func (f *impl) SyncStream(_ *context.T, call sync.SyncSyncStreamServerCall) error {
	stream := call.RecvStream()
	for stream.Advance() {
		// Nothing to do beside iterating over the stream.
	}
	return nil
}
