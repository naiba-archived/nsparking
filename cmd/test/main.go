package main

import (
	"context"
	"log"
	"net"

	"github.com/bobesa/go-domain-util/domainutil"
)

func main() {
	var domains = []string{
		"baidu.com.cn",
		"g.com.ms",
		"x.com.ax",
		"a.baidu.com.cn",
	}
	for i := 0; i < len(domains); i++ {
		log.Println(domains[i], domainutil.Domain(domains[i]))
	}

	resolver := &net.Resolver{
		PreferGo:     true,
		StrictErrors: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{}
			return d.DialContext(ctx, "udp", "223.5.5.5:53")
		},
	}

	ns, err := resolver.LookupNS(context.Background(), "nsparking.tk")
	log.Println(ns[0], err)
}
