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
	recordTtl  uint32
}

func NewDnsServer(address string, network string, baseDomain string, recordTtl uint32) *DnsServer {
	dnsServer := new(DnsServer)
	dnsServer.baseDomain = baseDomain
	dnsServer.server = &dns.Server{Addr: address, Net: network}
	dnsServer.recordTtl = recordTtl
	dns.HandleFunc(baseDomain, dnsServer.handleRequest)
	return dnsServer
}

func (d *DnsServer) StartAsync() {
	err := d.server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func (d *DnsServer) handleRequest(writer dns.ResponseWriter, reply *dns.Msg) {
	msg := new(dns.Msg)
	msg.SetReply(reply)

	for _, q := range reply.Question {
		if !strings.HasSuffix(q.Name, d.baseDomain) {
			println("requested domain", q.Name, "doesn't match the provided baseDomain")
			continue
		}

		hostname := strings.TrimSuffix(q.Name, d.baseDomain)
		hostname = strings.TrimSuffix(hostname, ".")

		typeString := dns.TypeToString[q.Qtype]

		switch typeString {
		case "SOA":
			soa := &dns.SOA{
				Hdr:    dns.RR_Header{Name: d.baseDomain, Rrtype: dns.TypeSOA, Class: dns.ClassINET, Ttl: 3600},
				Ns:     d.baseDomain,
				Mbox:   "root." + d.baseDomain,
				Serial: 20220509, Refresh: 7200, Retry: 3600, Expire: 86400, Minttl: 3600,
			}
			msg.Answer = append(msg.Answer, soa)
			break

		case "NS":
			ns := &dns.NS{
				Hdr: dns.RR_Header{Name: d.baseDomain, Rrtype: dns.TypeNS, Class: dns.ClassINET, Ttl: 3600},
				Ns:  d.baseDomain,
			}
			msg.Answer = append(msg.Answer, ns)
			break

		case "A":
		case "AAAA":
		case "CNAME":
			record, err := db.QueryDnsRecordWithType(hostname, typeString)
			if err != nil {
				println("cannot find dns record for", "\""+q.Name+"\"", "of type", "\""+typeString+"\":", err.Error())
				break
			}

			var ans dns.RR
			switch typeString {
			case "A":
				ans = &dns.A{
					Hdr: dns.RR_Header{Name: q.Name, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: d.recordTtl},
					A:   net.ParseIP(record),
				}
				break
			case "AAAA":
				ans = &dns.AAAA{
					Hdr:  dns.RR_Header{Name: q.Name, Rrtype: dns.TypeAAAA, Class: dns.ClassINET, Ttl: d.recordTtl},
					AAAA: net.ParseIP(record),
				}
			case "CNAME":
				ans = &dns.CNAME{
					Hdr:    dns.RR_Header{Name: q.Name, Rrtype: dns.TypeCNAME, Class: dns.ClassINET, Ttl: d.recordTtl},
					Target: record,
				}
			}
			msg.Answer = append(msg.Answer, ans)
		}
	}
	writer.WriteMsg(msg)
}
