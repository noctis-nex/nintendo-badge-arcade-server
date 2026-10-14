package main

import (
	"sync"

	"github.com/Protarium-Network/nintendo-badge-arcade-server/secure/nex"
)

var wg sync.WaitGroup

func main() {
	wg.Add(1)

	// TODO - Add gRPC server
	go nex.StartNEXServer()

	wg.Wait()
}
