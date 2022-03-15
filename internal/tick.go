package internal

import (
	"context"
	"fmt"
	"time"
)

const Interval = 15 * time.Minute

func (s *service) StartTick() {
	go s.dnsUpdater()
}

func (s *service) dnsUpdater() {
	fmt.Println("starting dns updater tick...")
	ticker := time.NewTicker(Interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			go s.checkAddr()
		}
	}
}

func (s *service) checkAddr() {
	ctx, cancel := context.WithTimeout(context.Background(), s.config.Timeout)
	defer cancel()
	prev := GetCachedIP()
	curr := s.fetchIP(ctx)
	if prev != curr {
		fmt.Printf("ip addr change, updating: prev=%s,curr=%s", prev, curr)
		SetCachedIP(curr)
	}
}
