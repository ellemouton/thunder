package ops

func StartLoops(b Backends) {
	go pingSyncchainForever(b)
}
