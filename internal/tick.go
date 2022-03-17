package internal

import (
	"time"
)

const Interval = 3 * time.Minute

func (s *service) StartTick() {
	go s.dnsUpdater()
}

func (s *service) dnsUpdater() {
	ticker := time.NewTicker(Interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			go s.Scan()
		}
	}
}
