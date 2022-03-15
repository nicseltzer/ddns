package main

import (
	"github.com/nicseltzer/ddns/internal"
)

func main() {
	service := internal.NewService()
	service.StartTick()
	service.Register()
	service.Start()

	select {}
}
