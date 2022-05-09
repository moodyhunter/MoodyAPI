package dns_server

import (
	"log"
	"net"
	"strings"

	"api.mooody.me/db"
	"github.com/miekg/dns"
)

type DnsServer struct {
	server     *dns.Server
	baseDomain string
}

var DNSRecordTTL uint32 = 30

func NewDnsServer(address string, network string, baseDomain string) *DnsServer {
	dnsServer := new(DnsServer)
	dnsServer.baseDomain = baseDomain
	dnsServer.server = &dns.Server{Addr: address, Net: network}
	dns.HandleFunc(baseDomain, dnsServer.handleRequest)
	return dnsServer
}

func (d *DnsServer) StartAsync() {
	err := d.server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func (d *DnsServer) handleRequest(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)

	for _, q := range r.Question {
		if !strings.HasSuffix(q.Name, d.baseDomain) {
			println("requested domain", q.Name, "doesn't match the provided baseDomain")
			continue
		}

		hostname := strings.TrimSuffix(q.Name, d.baseDomain)
		hostname = strings.TrimSuffix(hostname, ".")

		typeString := dns.TypeToString[q.Qtype]

		if typeString == "SOA" {
			soa := &dns.SOA{
				Hdr:    dns.RR_Header{Name: d.baseDomain, Rrtype: dns.TypeSOA, Class: dns.ClassINET, Ttl: 3600},
				Ns:     d.baseDomain,
				Mbox:   "root." + d.baseDomain,
				Serial: 20220509, Refresh: 7200, Retry: 3600, Expire: 86400, Minttl: 3600,
			}
			m.Answer = append(m.Answer, soa)
			continue
		}

		if typeString == "NS" {
			ns := &dns.NS{
				Hdr: dns.RR_Header{Name: d.baseDomain, Rrtype: dns.TypeNS, Class: dns.ClassINET, Ttl: 3600},
				Ns:  d.baseDomain,
			}
			m.Answer = append(m.Answer, ns)
			continue
		}

		record, err := db.QueryDnsRecordWithType(hostname, typeString)
		if err != nil {
			println("cannot find dns record for", "\""+q.Name+"\"", "of type", "\""+typeString+"\":", err.Error())
			continue
		}

		var ans dns.RR

		switch q.Qtype {
		case dns.TypeA:
			ans = &dns.A{
				Hdr: dns.RR_Header{Name: q.Name, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: DNSRecordTTL},
				A:   net.ParseIP(record),
			}
			break
		case dns.TypeAAAA:
			ans = &dns.AAAA{
				Hdr:  dns.RR_Header{Name: q.Name, Rrtype: dns.TypeAAAA, Class: dns.ClassINET, Ttl: DNSRecordTTL},
				AAAA: net.ParseIP(record),
			}
			break
		case dns.TypeCNAME:
			ans = &dns.CNAME{
				Hdr:    dns.RR_Header{Name: q.Name, Rrtype: dns.TypeCNAME, Class: dns.ClassINET, Ttl: DNSRecordTTL},
				Target: record,
			}
			break
		}

		if ans != nil {
			m.Answer = append(m.Answer, ans)
		}
	}
	w.WriteMsg(m)
}
