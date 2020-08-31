package dns6

import (
	"log"

	"github.com/miekg/dns"
)

var (
	dnsMux *dns.ServeMux
)

type dnsHandler struct {
	opts *dnsOptions
}

func newHandler(opts *dnsOptions) (mux *dns.ServeMux) {
	handler := &dnsHandler{
		opts: opts,
	}
	mux = dns.NewServeMux()
	mux.HandleFunc(".", handler.serveDNS)
	return
}

func (h *dnsHandler) serveDNS(w dns.ResponseWriter, msg *dns.Msg) {
	defer w.Close()

	if len(msg.Question) == 0 {
		return
	}

	var (
		resp *dns.Msg
		err  error
	)

	question := &msg.Question[0]
	switch question.Qtype {
	case dns.TypeA:
		resp, err = h.requestTypeA(msg)
	case dns.TypeAAAA:
		resp, err = h.requestTypeAAAA(msg)
	default:
		resp, err = h.defaultType(msg)
	}

	if err != nil {
		log.Println(err)
		return
	}

	if err = w.WriteMsg(resp); err != nil {
		log.Println(err)
	}
}

func (h *dnsHandler) requestTypeA(req *dns.Msg) (resp *dns.Msg, err error) {
	client := new(dns.Client)
	resp, _, err = client.Exchange(req, h.opts.dnsServer)
	return
}

func (h *dnsHandler) requestTypeAAAA(req *dns.Msg) (resp *dns.Msg, err error) {
	question6to4(req)
	client := new(dns.Client)
	if resp, _, err = client.Exchange(req, h.opts.dnsServer); err != nil {
		return
	}
	rr4to6(resp, h.opts.prefix, h.opts.prefixLength)
	return
}

func (h *dnsHandler) defaultType(req *dns.Msg) (resp *dns.Msg, err error) {
	client := new(dns.Client)
	resp, _, err = client.Exchange(req, h.opts.dnsServer)
	return
}
