package dns6

import (
	"net"

	"github.com/miekg/dns"
)

func composeIP6(prefix net.IP, l int, ip4 net.IP) (ip6 net.IP) {
	ip6 = make(net.IP, net.IPv6len)
	n := l / 8
	for i := 0; i < n; i++ {
		ip6[i] = prefix[i]
	}

	ip6[n+0] = ip4[0]
	ip6[n+1] = ip4[1]
	ip6[n+2] = ip4[2]
	ip6[n+3] = ip4[3]
	return
}

func header6to4(msg *dns.Msg) {
	for i := 0; i < len(msg.Question); i++ {
		q := &msg.Question[i]
		if q.Qtype != dns.TypeAAAA {
			continue
		}
		q.Qtype = dns.TypeA
	}
}

func header4to6(msg *dns.Msg) {
	for i := 0; i < len(msg.Question); i++ {
		q := &msg.Question[i]
		if q.Qtype != dns.TypeA {
			continue
		}
		q.Qtype = dns.TypeAAAA
	}
}

func question6to4(msg *dns.Msg) {
	header6to4(msg)
}

func rr4to6(resp *dns.Msg, prefix net.IP, l int) {
	header4to6(resp)
	for i := 0; i < len(resp.Answer); i++ {
		answer, ok := resp.Answer[i].(*dns.A)
		if !ok || answer.A == nil || answer.Hdr.Rdlength != net.IPv4len {
			continue
		}
		rr := &dns.AAAA{
			Hdr: answer.Hdr,
		}
		rr.AAAA = composeIP6(prefix, l, answer.A)
		rr.Hdr.Rrtype = dns.TypeAAAA
		rr.Hdr.Rdlength = net.IPv6len
		resp.Answer[i] = rr
	}
}
