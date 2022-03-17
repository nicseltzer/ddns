package internal

import (
	"fmt"
	"io"
	"net/http"
)

func (s *service) Register() {
	internalMux := s.internalMux

	internalMux.HandleFunc("/v1/ddns/ip", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		ip := s.ipCache.IP()
		if ip == "" {
			w.WriteHeader(http.StatusServiceUnavailable)
			io.WriteString(w, fmt.Sprintf("%d service unavailable", http.StatusServiceUnavailable))
			return
		}
		_, err := io.WriteString(w, ip)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}
