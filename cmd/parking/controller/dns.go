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
	log.Println("Accept DNS query", domain, err)
	if err != nil {
		w.WriteMsg(&msg)
		return
	}
	var a dns.RR
	msg.Authoritative = true
	switch r.Question[0].Qtype {
	case dns.TypeCNAME:
		a = &dns.CNAME{
			Hdr:    dns.RR_Header{Name: domain, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 60},
			Target: p.Value,
		}
	case dns.TypeA:
		switch p.Mode {
		case "cname":
			a = &dns.CNAME{
				Hdr:    dns.RR_Header{Name: domain, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 60},
				Target: p.Value,
			}
		case "a":
			a = &dns.A{
				Hdr: dns.RR_Header{Name: domain, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 60},
				A:   net.ParseIP(p.Value),
			}
		case "url":
			a = &dns.A{
				Hdr: dns.RR_Header{Name: domain, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 60},
				A:   net.ParseIP(model.IP),
			}
		}
	}
	msg.Answer = append(msg.Answer, a)
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
