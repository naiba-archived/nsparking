package controller

import (
	"net"
	"strconv"

	"github.com/naiba/nsparking/model"

	"github.com/naiba/nsparking/pkg/log"

	"github.com/miekg/dns"
)

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
		_, err := getRedirectByDomain(domain)
		if err == nil {
			msg.Answer = append(msg.Answer, &dns.A{
				Hdr: dns.RR_Header{Name: domain, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 60},
				A:   net.ParseIP(model.IP),
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
