package main

import (
	"log"

	"github.com/bobesa/go-domain-util/domainutil"
	"github.com/miekg/dns"
)

func main() {
	var domains = []string{
		"baidu.com.cn",
		"g.com.ms",
		"x.com.ax",
		"a.baidu.com.cn",
	}
	for i := 0; i < len(domains); i++ {
		log.Println(domains[i], domainutil.Domain(domains[i]), dns.SplitDomainName(domains[i]))
	}

	var resolver dns.Client
	q := new(dns.Msg)
	q.SetQuestion("nsparking.tk.", dns.TypeNS)
	q.RecursionDesired = true
	msg, _, _ := resolver.Exchange(q, "223.5.5.5:53")
	for i := 0; i < len(msg.Answer); i++ {
		log.Println(msg.Answer[i].(*dns.NS).Ns)
	}
}
