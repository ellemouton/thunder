package macaroon

import (
	"context"
	"github.com/lightningnetwork/lnd/lntypes"
	"gopkg.in/macaroon.v2"
)

type Client interface {
	Close() error
	Create(ctx context.Context, paymentHash lntypes.Hash, resourceType string, resourceID int64) (*macaroon.Macaroon, error)
	Verify(ctx context.Context, macBytes []byte, preimage []byte, resourceType string, resourceID int64) (bool, error)
}
