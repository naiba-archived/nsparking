package controller

import (
	"net"
	"strconv"

	"github.com/naiba/nsparking/pkg/log"

	"github.com/miekg/dns"
)

var domainsToAddresses map[string]string = map[string]string{}

// DNSHandler 服务器
type DNSHandler struct{}

// ServeDNS ...
func (dh *DNSHandler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	msg := dns.Msg{}
	msg.SetReply(r)
	switch r.Question[0].Qtype {
	case dns.TypeA:
		msg.Authoritative = true
		domain := msg.Question[0].Name
		address, ok := domainsToAddresses[domain]
		if ok {
			msg.Answer = append(msg.Answer, &dns.A{
				Hdr: dns.RR_Header{Name: domain, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 60},
				A:   net.ParseIP(address),
			})
		}
	}
	w.WriteMsg(&msg)
}

// ServeDNS ..
func ServeDNS() {
	log.Println("Starting DNS server !!")
	srv := &dns.Server{Addr: ":" + strconv.Itoa(53), Net: "udp"}
	srv.Handler = &DNSHandler{}
	if err := srv.ListenAndServe(); err != nil {
		log.Println("Failed to set udp listener ", err.Error())
	}
}
