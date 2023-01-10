// Package tcp is for TCP DNS lookup.
package main

import (
	"fmt"
	"net"
	"os"
	"strings"

	"google3/third_party/golang/dns/dns"
)

// Query returns the resolved IP address for the given domain.
func Query(domain string) (net.IP, error) {
	// We always start with a root nameserver.
	// TODO: Nameservers can become unavailable. Choose and cache the right one.
	nameserver := net.ParseIP("198.41.0.4")
	return resolve(domain, nameserver, dnsQuery)
}

func resolve(name string, nameserver net.IP, dnsQueryFunc func(string, net.IP) *dns.Msg) (net.IP, error) {
	for {
		reply := dnsQueryFunc(name, nameserver)
		if reply == nil {
			return nil, fmt.Errorf("DNS returned nil for domain %v from nameserver %v", name, nameserver)
		}
		if ip := getAnswer(reply); ip != nil {
			// Best case: we get an answer to our query and we're done
			return ip, nil
		} else if nsIP := getGlue(reply); nsIP != nil {
			// Second best: we get a "glue record" with the *IP address* of another nameserver to query
			nameserver = nsIP
		} else if domain := getNS(reply); domain != "" {
			// Third best: we get the *domain name* of another nameserver to query, which we can look up the IP for
			var err error
			if nameserver, err = resolve(domain, nameserver, dnsQueryFunc); err != nil {
				return nil, err
			}
		} else {
			return nil, fmt.Errorf("TCP resovler failed for domain %v from nameserver %v", name, nameserver)
		}
	}
}

func getAnswer(reply *dns.Msg) net.IP {
	for _, record := range reply.Answer {
		if record.Header().Rrtype == dns.TypeA {
			fmt.Println("  ", record, " <- answer")
			return record.(*dns.A).A
		}
	}
	return nil
}

func getGlue(reply *dns.Msg) net.IP {
	for _, record := range reply.Extra {
		if record.Header().Rrtype == dns.TypeA {
			fmt.Println("  ", record, " <- glue")
			return record.(*dns.A).A
		}
	}
	return nil
}

func getNS(reply *dns.Msg) string {
	for _, record := range reply.Ns {
		if record.Header().Rrtype == dns.TypeNS {
			fmt.Println("  ", record, " <- NS")
			return record.(*dns.NS).Ns
		}
	}
	return ""
}

func dnsQuery(name string, server net.IP) *dns.Msg {
	fmt.Printf("dig -r @%s %s\n", server.String(), name)
	msg := new(dns.Msg)
	msg.SetQuestion(name, dns.TypeA)
	c := new(dns.Client)
	c.Net = "tcp"
	reply, _, _ := c.Exchange(msg, server.String()+":53")
	return reply
}

func main() {
	name := os.Args[1]
	nameserver := os.Args[2]
	if !strings.HasSuffix(name, ".") {
		name = name + "."
	}
	if ip, err := resolve(name, net.ParseIP(nameserver), dnsQuery); err == nil {
		fmt.Println("Result:", ip)
	} else {
		fmt.Println(err.Error())
	}
}
