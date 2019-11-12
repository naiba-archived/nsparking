package main

import (
	"log"
	"strconv"

	"github.com/miekg/dns"

	"github.com/naiba/nsparking/cmd/parking/controller"
)

func main() {
	go ns()
}

func ns() {
	srv := &dns.Server{Addr: ":" + strconv.Itoa(53), Net: "udp"}
	srv.Handler = &controller.DNSHandler{}
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Failed to set udp listener %s\n", err.Error())
	}
}
