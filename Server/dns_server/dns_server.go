package dns_server

import (
	"log"
	"net"
	"strings"

	"api.mooody.me/db"
	"github.com/miekg/dns"
)

type DnsServer struct {
	Server     *dns.Server
	baseDomain string
	recordTtl  uint32
}

func NewDnsServer(address string, network string, baseDomain string, recordTtl uint32) *DnsServer {
	dnsServer := new(DnsServer)
	dnsServer.baseDomain = baseDomain
	dnsServer.Server = &dns.Server{Addr: address, Net: network}
	dnsServer.recordTtl = recordTtl
	dns.HandleFunc(baseDomain, dnsServer.handleRequest)
	return dnsServer
}

func (d *DnsServer) Close() {
	d.Server.Shutdown()
	log.Println("DNS server closed")
}

func (d *DnsServer) handleRequest(writer dns.ResponseWriter, reply *dns.Msg) {
	msg := new(dns.Msg)
	msg.SetReply(reply)

	for _, q := range reply.Question {
		q.Name = strings.ToLower(q.Name)
		if !strings.HasSuffix(q.Name, d.baseDomain) {
			println("requested domain", q.Name, "doesn't match the provided baseDomain")
			continue
		}

		hostname := strings.TrimSuffix(q.Name, d.baseDomain)
		hostname = strings.TrimSuffix(hostname, ".")

		typeString := dns.TypeToString[q.Qtype]
		soa := &dns.SOA{
			Hdr:    dns.RR_Header{Name: d.baseDomain, Rrtype: dns.TypeSOA, Class: dns.ClassINET, Ttl: 3600},
			Ns:     d.baseDomain,
			Mbox:   "root." + d.baseDomain,
			Serial: 20220509, Refresh: 7200, Retry: 3600, Expire: 86400, Minttl: 3600,
		}

		switch typeString {
		case "SOA":
			msg.Answer = append(msg.Answer, soa)
		case "NS":
			ns := &dns.NS{
				Hdr: dns.RR_Header{Name: d.baseDomain, Rrtype: dns.TypeNS, Class: dns.ClassINET, Ttl: 3600},
				Ns:  d.baseDomain,
			}
			msg.Answer = append(msg.Answer, ns)

		case "A", "AAAA", "CNAME":
			record, err := db.QueryDnsRecordWithType(hostname, typeString)
			if err != nil {
				println("cannot find dns record for", "\""+q.Name+"\"", "of type", "\""+typeString+"\":", err.Error())
				msg.Ns = append(msg.Ns, soa)
				// set the error code to "name error"
				msg.Rcode = dns.RcodeNameError
				break
			}

			var ans dns.RR
			switch typeString {
			case "A":
				ans = &dns.A{
					Hdr: dns.RR_Header{Name: q.Name, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: d.recordTtl},
					A:   net.ParseIP(record),
				}

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
