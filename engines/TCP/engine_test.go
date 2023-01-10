package main

import (
	"net"
	"testing"

	"google3/third_party/golang/dns/dns"
)

func emptyReplyDNSQuery(name string, server net.IP) *dns.Msg {
	if name == "google.com" {
		return &dns.Msg{}
	}
	return &dns.Msg{}
}

func simpleAnswerDNSQuery(name string, server net.IP) *dns.Msg {
	answer := &dns.A{
		Hdr: dns.RR_Header{Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 0},
		A:   net.ParseIP("127.0.0.1").To4(),
	}
	if name == "google.com" {
		m := new(dns.Msg)
		m.Answer = append(m.Answer, answer)
		return m
	}
	return &dns.Msg{}
}

func withGlueDNSQuery(name string, server net.IP) *dns.Msg {
	glue := &dns.A{
		Hdr: dns.RR_Header{Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 0},
		A:   net.ParseIP("1.1.1.1").To4(),
	}
	answer := &dns.A{
		Hdr: dns.RR_Header{Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 0},
		A:   net.ParseIP("127.0.0.1").To4(),
	}
	if name == "google.com" {
		if server.Equal(net.ParseIP("198.41.0.4")) {
			m := new(dns.Msg)
			m.Extra = append(m.Extra, glue)
			return m
		} else if server.Equal(net.ParseIP("1.1.1.1")) {
			m := new(dns.Msg)
			m.Answer = append(m.Answer, answer)
			return m
		}
	}
	return nil
}

func withNSDNSQuery(name string, server net.IP) *dns.Msg {
	ns := &dns.NS{
		Hdr: dns.RR_Header{Rrtype: dns.TypeNS, Class: dns.ClassINET, Ttl: 0},
		Ns:  "nameserver.com",
	}
	nsAnswer := &dns.A{
		Hdr: dns.RR_Header{Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 0},
		A:   net.ParseIP("1.1.1.1").To4(),
	}
	answer := &dns.A{
		Hdr: dns.RR_Header{Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 0},
		A:   net.ParseIP("127.0.0.1").To4(),
	}
	if name == "google.com" {
		if server.Equal(net.ParseIP("198.41.0.4")) {
			m := new(dns.Msg)
			m.Ns = append(m.Ns, ns)
			return m
		} else if server.Equal(net.ParseIP("1.1.1.1")) {
			m := new(dns.Msg)
			m.Answer = append(m.Answer, answer)
			return m
		}
	} else if name == "nameserver.com" {
		m := new(dns.Msg)
		m.Answer = append(m.Answer, nsAnswer)
		return m
	}
	return nil
}

func TestResolve(t *testing.T) {
	tests := []struct {
		testName   string
		dnsQuery   func(name string, server net.IP) *dns.Msg
		nameserver net.IP
		wantErr    bool
		wantIP     net.IP
	}{
		{"emptyReplyDNSQuery", emptyReplyDNSQuery, net.ParseIP("198.41.0.4"), true, nil},
		{"simpleAnswerDNSQuery", simpleAnswerDNSQuery, net.ParseIP("198.41.0.4"), false, net.ParseIP("127.0.0.1")},
		{"withGlueDNSQuery", withGlueDNSQuery, net.ParseIP("198.41.0.4"), false, net.ParseIP("127.0.0.1")},
		{"withNSDNSQuery", withNSDNSQuery, net.ParseIP("198.41.0.4"), false, net.ParseIP("127.0.0.1")},
	}

	for _, v := range tests {
		var gotIP net.IP
		var err error
		gotIP, err = resolve("google.com", v.nameserver, v.dnsQuery)
		if v.wantErr && err == nil {
			t.Errorf("Expected error for inputs %v", v)
		}
		if gotIP.String() != v.wantIP.String() {
			t.Errorf("Wrong IP returned. Got %v, wanted %v for dnsquery %v", gotIP, v.wantIP, v.testName)
		}
	}
}
