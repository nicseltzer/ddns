package internal

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type Service interface {
	Register()
	StartTick()
	Start()
}

type service struct {
	config      *Config
	internalMux *mux.Router
}

func NewService() Service {
	return newService()
}

func newService() *service {
	return &service{
		config:      NewConfig(),
		internalMux: mux.NewRouter(),
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

	// Lifecycle
	s.setup()
	defer s.tearDown()
}

func (s *service) setup() {
}

func (s *service) tearDown() {
	fmt.Println("shutting down...")
}
