package workers

import (
	"log"
	"markat/internal/imapclient"
	"time"
)

func NoopTicker(stop chan bool, ic *imapclient.ImapClient) {
	ticker := time.NewTicker(time.Minute * 1)

	go func() {
		for {
			select {
			case <-stop:
				log.Println("Noop ticker stopped")
				return
			case <-ticker.C:
				err := ic.Noop(stop)
				if err == nil {
					log.Println("NOOP")
				}
			}
		}
	}()
}
