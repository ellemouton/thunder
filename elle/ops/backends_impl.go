package ops

import (
	"testing"

	"github.com/ellemouton/thunder/syncchain"
)

type backendsImpl struct {
	syncchainClient syncchain.Client
}

func MakeBackends(
	syncchainClient syncchain.Client,
) *backendsImpl {
	return &backendsImpl{
		syncchainClient: syncchainClient,
	}
}

type backendsOption func(*backendsImpl)

func NewBackendsForTesting(_ testing.TB, opts ...backendsOption) *backendsImpl {
	var b backendsImpl
	for _, opt := range opts {
		opt(&b)
	}
	return &b
}

func WithSyncchainClient(sccli syncchain.Client) backendsOption {
	return func(b *backendsImpl) {
		b.syncchainClient = sccli
	}
}
