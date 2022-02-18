package internal

import (
	"fmt"
	"github.com/nicseltzer/ddns/pkg/cloudflare"
)

type Service interface {
	UpdateDNS()
}

type service struct {
	config Config
}

func NewService() Service {
	return newService()
}

func newService() *service {
	return &service{
		config: NewConfig(),
	}
}

func (s *service) UpdateDNS() {
	client := cloudflare.NewClient("foo", "bar", "baz")
	fmt.Println(client)
}
