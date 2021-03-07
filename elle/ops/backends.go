package ops

import "github.com/ellemouton/thunder/syncchain"

type Backends interface {
	SyncchainClient() syncchain.Client
}
