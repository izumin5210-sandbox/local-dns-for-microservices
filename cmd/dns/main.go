package main

import (
	"log"
	"net"
	"strings"

	"github.com/miekg/dns"
)

func main() {
	mux := dns.NewServeMux()

	var localhost net.IP
	ips, err := net.LookupIP("localhost")
	if err != nil {
		log.Fatal(err.Error())
	}
	for _, ip := range ips {
		switch ip.String() {
		case "127.0.0.1", "::1":
			continue
		default:
			localhost = ip
			break
		}
	}
	if len(localhost) == 0 {
		localhost = net.IPv4(127, 0, 0, 1)
	}

	mux.HandleFunc(".", func(w dns.ResponseWriter, req *dns.Msg) {
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
				A: localhost,
			})
		} else {
			resp.MsgHdr.Rcode = dns.RcodeNameError
		}

		log.Printf("<-- %#v\n", resp)

		w.WriteMsg(resp)
	})

	h := dns.HandlerFunc(func(w dns.ResponseWriter, r *dns.Msg) {
		log.Printf("--> %#v\n", r)
		mux.ServeDNS(w, r)
	})

	if err := dns.ListenAndServe(":53", "udp", h); err != nil {
		log.Fatal(err.Error())
	}
}
