package workers

import (
	"log"
	"markat/internal/imapclient"
	"time"
)

func NoopTicker(done chan bool, ic *imapclient.ImapClient) {
	ticker := time.NewTicker(time.Minute * 1)

	go func() {
		for {
			select {
			case <-done:
				log.Println("Noop ticker stopped")
				return
			case <-ticker.C:
				ic.Noop()
				log.Println("NOOP")
			}
		}
	}()
}
