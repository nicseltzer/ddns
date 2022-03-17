package internal

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nicseltzer/ddns/pkg/ipcache"
	"net/http"
)

type Service interface {
	Register()
	StartTick()
	Start()
	Scan()
}

type service struct {
	config      *Config
	internalMux *mux.Router
	ipCache     ipcache.IPCache
}

func NewService() Service {
	return newService()
}

func newService() *service {
	return &service{
		config:      NewConfig(),
		internalMux: mux.NewRouter(),
		ipCache:     ipcache.NewIPCache(),
	}
}

func (s *service) Start() {
	// API
	getListenAddress := func(host string, port int) string {
		if host == "*" {
			return fmt.Sprintf(":%d", port)
		}
		return fmt.Sprintf("%s:%d", host, port)
	}

	listenAddress := getListenAddress("*", s.config.Port)
	fmt.Println("listening on: " + listenAddress)
	go http.ListenAndServe(listenAddress, s.internalMux)

	fmt.Println("starting ddns...")
	ip := s.fetchExternalIP(context.Background())
	s.ipCache.SetIP(ip)
	s.Scan()
}
