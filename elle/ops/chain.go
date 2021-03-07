package ops

import (
	"context"
	"log"
	"time"
)

func pingSyncchainForever(b Backends) {
	ctx := context.TODO()
	for {
		pong, err := b.SyncchainClient().Ping(ctx, "testing testing 123")
		if err != nil {
			log.Print(err)
		} else {
			log.Print(pong)
		}

		time.Sleep(10 * time.Minute)
	}
}
