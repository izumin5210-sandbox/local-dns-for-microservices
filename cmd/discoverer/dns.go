package main

import (
	"context"
	"log"
	"net"
	"strings"

	"github.com/miekg/dns"
)

type DNSServer struct {
	mapping   Mapping
	server    *dns.Server
	localhost net.IP
}

func NewDNSServer(mapping Mapping) *DNSServer {
	return &DNSServer{
		mapping: mapping,
	}
}

func (s *DNSServer) Serve(addr string) error {
	s.updateLocalhostIP()
	s.server = &dns.Server{
		Handler: dns.HandlerFunc(s.handle),
		Addr:    addr,
		Net:     "udp",
	}
	log.Printf("start DNS server on %s\n", addr)
	err := s.server.ListenAndServe()
	if err != nil {
		log.Println(err)
	}
	return err
}

func (s *DNSServer) Shutdown(context.Context) error {
	return s.server.Shutdown()
}

func (s *DNSServer) handle(w dns.ResponseWriter, req *dns.Msg) {
	// log
	log.Printf("received %#v\n", req.Question)

	// handle
	q := req.Question[0]
	resp := new(dns.Msg)
	resp.SetReply(req)

	if q.Qtype == dns.TypeA && q.Qclass == dns.ClassINET && strings.HasSuffix(q.Name, "izumin.local.") {
		resp.Answer = append(resp.Answer, &dns.A{
			Hdr: dns.RR_Header{
				Name:   q.Name,
				Rrtype: dns.TypeA,
				Class:  dns.ClassINET,
				Ttl:    600,
			},
			A: s.localhost,
		})
	} else {
		resp.MsgHdr.Rcode = dns.RcodeNameError
	}

	w.WriteMsg(resp)
}

func (s *DNSServer) updateLocalhostIP() {
	s.localhost = net.IPv4(127, 0, 0, 1)

	ips, _ := net.LookupIP("localhost")

	for _, ip := range ips {
		switch ip.String() {
		case "127.0.0.1", "::1":
			continue
		default:
			s.localhost = ip
			break
		}
	}
}
