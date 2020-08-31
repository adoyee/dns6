package dns6

import (
	"net"

	"github.com/miekg/dns"
)

type DnsMode int

const (
	buffSize = 1024
	lenIP6   = 16
	lenIP4   = 4

	modeProxy DnsMode = iota + 1
	modeDNS64
	modeIVI

	cacheTimeout = 3600
	dialTimeout  = 3
)

var (
	server *dns.Server

	modeString = []string{
		"Proxy", "DNS64", "IVI",
	}

	wellKnownPrefix = "64:ff9b::0"
	defaultLength   = 96
	defaultServer   *Server
)

func (mode DnsMode) String() string {
	return modeString[mode]
}

type Server struct {
	opts    *dnsOptions
	srv     *dns.Server
	udpConn net.PacketConn
}

func (s *Server) applyOpts(args ...optionApply) (err error) {
	if s.opts == nil {
		s.opts = newOptions()
	}

	for _, opt := range args {
		_, err = opt.apply(s.opts)
		if err != nil {
			return err
		}
	}

	return
}

func (s *Server) ListenAndServe(addr string, args ...optionApply) (err error) {
	if err = s.applyOpts(args...); err != nil {
		return
	}

	switch s.opts.mode {
	case modeDNS64, modeIVI:
		s.udpConn, err = net.ListenPacket("udp6", addr)
	default:
		s.udpConn, err = net.ListenPacket("udp", addr)
	}

	if err != nil {
		return
	}

	if _, err = net.ResolveUDPAddr("udp", s.opts.dnsServer); err != nil {
		return
	}

	s.srv = &dns.Server{
		PacketConn: s.udpConn,
		Handler:    newHandler(s.opts),
	}

	return s.srv.ActivateAndServe()
}

func ListenAndServe(addr string, args ...optionApply) (err error) {
	if defaultServer == nil {
		defaultServer = &Server{}
	}
	return defaultServer.ListenAndServe(addr, args...)
}
