package main

import (
	"time"
)

func main() {
	sync := make(chan time.Timer)

	select {
	case <-sync:
		// Sync on an interval
	}
}
