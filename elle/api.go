package elle

import "context"

type Client interface {
        RequiresPayment(ctx context.Context, path string) (*AssetDetails, error)
        Close() error
}