package main

import (
	"github.com/nicseltzer/ddns/internal"
)

func main() {
	service := internal.NewService()
	service.UpdateDNS()
}
