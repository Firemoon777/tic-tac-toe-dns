package main

import (
	"fmt"
	"log"
	"strings"
	"strconv"

	"github.com/miekg/dns"
		"io/ioutil"
	"net"
)

func parseQuery(m *dns.Msg, ip string) {
	for _, q := range m.Question {
		switch q.Qtype {
		case dns.TypeTXT:
			log.Printf("/ TXT / %s / %s \n", q.Name, ip)
			r, err := ioutil.ReadFile("./res/" + q.Name + "txt")
			if err == nil {
				str := string(r)
				for _, s := range strings.Split(str, "\n") {
					if s != "" {
						t := &dns.TXT{
							Hdr: dns.RR_Header{Name: q.Name, Rrtype: dns.TypeTXT, Class: dns.ClassINET, Ttl: 0},
							Txt: strings.Split(s, "|"),
						}
						m.Answer = append(m.Answer, t)
					}
				}
			}
		case dns.TypeA:
			log.Printf("/ A / %s / %s \n", q.Name, ip)
			rr, err := dns.NewRR(fmt.Sprintf("%s A 127.0.0.1", q.Name))
			if err == nil {
				m.Answer = append(m.Answer, rr)
			}
		}
	}
}

func handleDnsRequest(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)
	m.Compress = false

	var ip_str string

	if ip, ok := w.RemoteAddr().(*net.TCPAddr); ok {
		if ip.IP.To4() != nil {
			ip_str = ip.IP.To4().String()
		} else {
			ip_str = ip.IP.String()
		}
	}


	switch r.Opcode {
	case dns.OpcodeQuery:
		parseQuery(m, ip_str)
	}

	w.WriteMsg(m)
}

func main() {
	// attach request handler func
	dns.HandleFunc("game.f1remoon.com.", handleDnsRequest)

	// start server
	port := 5353
	server := &dns.Server{Addr: ":" + strconv.Itoa(port), Net: "udp"}
	log.Printf("Starting at %d\n", port)
	err := server.ListenAndServe()
	defer server.Shutdown()
	if err != nil {
		log.Fatalf("Failed to start server: %s\n ", err.Error())
	}
}
