package main

import (
	"github.com/razvanm/rpc-benchmarks/vanadium-core"
	"v.io/v23/context"
	"v.io/v23/rpc"
	_ "v.io/x/ref/runtime/factories/generic"
)

type impl struct {
}

func (f *impl) Sink(_ *context.T, _ rpc.ServerCall, payload []byte) error {
	return nil
}

func (f *impl) SinkStream(_ *context.T, call sink.SinkSinkStreamServerCall) error {
	stream := call.RecvStream()
	for stream.Advance() {
		// Nothing to do beside iterating over the stream.
	}
	return nil
}
