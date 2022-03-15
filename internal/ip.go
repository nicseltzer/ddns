package internal

import (
	"context"
	"github.com/nicseltzer/ddns/pkg/ifconfigme"
)

func (s *service) fetchIP(ctx context.Context) string {
	ip := ifconfigme.NewClient(s.config.Timeout)
	return ip.Ip(ctx)
}
