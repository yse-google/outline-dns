// Package udp is for UDP DNS lookup. This most probably returns invalid
// results, which can be used to check the result of other engines.
package udp

import (
	"context"
	"log"
	"net"
	"time"
)

// Query DNS
func Query(domain string) []net.IP {
	r := &net.Resolver{
		// Force not using cgo implementation which would bypass using our supplied Dial.
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: time.Millisecond * time.Duration(10000),
			}
			return d.DialContext(ctx, "tcp4", "8.8.8.8:53")
		},
	}
	ips, error := r.LookupIP(context.Background(), "ip4", domain)
	if error != nil {
		log.Fatal(error.Error())
	}

	return ips
}
