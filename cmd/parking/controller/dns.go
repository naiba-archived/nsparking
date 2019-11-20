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
	domain := msg.Question[0].Name
	p, err := getRedirectByDomain(domain)
	log.Println("Accept DNS query", domain, err, r.Question[0].Qtype)
	if err != nil {
		w.WriteMsg(&msg)
		return
	}
	msg.Authoritative = true
	var a []dns.RR
	switch r.Question[0].Qtype {
	case dns.TypeCNAME:
		switch p.Mode {
		case "cname":
			a = append(a, &dns.CNAME{
				Hdr:    dns.RR_Header{Name: domain, Rrtype: dns.TypeCNAME, Class: dns.ClassINET, Ttl: 60},
				Target: p.Value,
			})
		}
	case dns.TypeA:
		switch p.Mode {
		case "cname":
			// msg.Extra = append(msg.Extra, &dns.CNAME{
			// 	Hdr:    dns.RR_Header{Name: domain, Rrtype: dns.TypeCNAME, Class: dns.ClassINET, Ttl: 60},
			// 	Target: p.Value,
			// })
			a = append(a, &dns.CNAME{
				Hdr:    dns.RR_Header{Name: domain, Rrtype: dns.TypeCNAME, Class: dns.ClassINET, Ttl: 60},
				Target: p.Value,
			})
			as, err := getA(p.Value)
			if err == nil {
				a = append(a, as...)
			}
		case "a":
			a = append(a, &dns.A{
				Hdr: dns.RR_Header{Name: domain, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 60},
				A:   net.ParseIP(p.Value),
			})
		case "url":
			a = append(a, &dns.A{
				Hdr: dns.RR_Header{Name: domain, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 60},
				A:   net.ParseIP(model.IP),
			})
		}
	}
	log.Println(p.Mode, a)
	msg.Answer = append(msg.Answer, a...)
	w.WriteMsg(&msg)
	log.Println("Response DNS query", domain)
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
