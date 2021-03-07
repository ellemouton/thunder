package syncchain

import "context"

type Client interface {
	Ping(ctx context.Context, ping string) (string, error)
	Close() error
}
