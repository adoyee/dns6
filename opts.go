package dns6

import (
	"errors"
	"net"
)

type dnsOptions struct {
	local        string
	dnsServer    string
	mode         DnsMode
	cacheTimeout int
	dialTimeout  int
	prefix       net.IP
	prefixLength int
}

type optionApply interface {
	apply(*dnsOptions) (*dnsOptions, error)
}

type withCacheTimeout struct {
	timeout int
}

func (o *withCacheTimeout) apply(opts *dnsOptions) (*dnsOptions, error) {
	opts.dialTimeout = o.timeout
	return opts, nil
}

type withDialTimeout struct {
	timeout int
}

func (o *withDialTimeout) apply(opts *dnsOptions) (*dnsOptions, error) {
	opts.dialTimeout = o.timeout
	return opts, nil
}

type withMode struct {
	mode DnsMode
}

func (o *withMode) apply(opts *dnsOptions) (*dnsOptions, error) {
	opts.mode = o.mode
	return opts, nil
}

type withPrefix struct {
	prefix string
}

func (o *withPrefix) apply(opts *dnsOptions) (*dnsOptions, error) {
	ip := net.ParseIP(o.prefix)
	if ip == nil || len(ip) != lenIP6 {
		return nil, errors.New(o.prefix + "is not an ip6 address")
	}
	opts.prefix = ip
	return opts, nil
}

type withPrefixLen struct {
	length int
}

func (o *withPrefixLen) apply(opts *dnsOptions) (*dnsOptions, error) {
	if o.length > 96 || o.length < 32 {
		return nil, errors.New("prefix length must be between from 32 to 96")
	}
	if o.length%8 != 0 {
		return nil, errors.New("prefix length must be multiple of 8 ")
	}
	opts.prefixLength = o.length
	return opts, nil
}

func WithCacheTimeout(t int) *withCacheTimeout {
	if t <= 0 || t >= cacheTimeout {
		t = cacheTimeout
	}
	return &withCacheTimeout{timeout: t}
}

func WithDialTimeout(t int) *withDialTimeout {
	if t <= 0 || t >= dialTimeout {
		t = dialTimeout
	}
	return &withDialTimeout{timeout: t}
}

func WithPrefix(prefix string) *withPrefix {
	return &withPrefix{
		prefix: prefix,
	}
}

func WithPrefixLength(length int) *withPrefixLen {
	return &withPrefixLen{length: length}
}

var (
	WithModeProxy = &withMode{mode: modeProxy}
	WithModeDNS64 = &withMode{mode: modeDNS64}
	WithModeIVI   = &withMode{mode: modeIVI}
)

func newOptions() *dnsOptions {
	opts := &dnsOptions{
		dnsServer:    "8.8.8.8:53",
		mode:         modeIVI,
		cacheTimeout: cacheTimeout,
		dialTimeout:  dialTimeout,
		prefix:       net.ParseIP(wellKnownPrefix),
		prefixLength: defaultLength,
	}
	return opts
}
