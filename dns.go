// Package dns is a DNS lookup library. This can do DNS lookups using multiple engines
// and select the result that's most trustable.
package dns

// List of engines we have available for DNS lookup
const (
	HTTPS
	TLS
	DNSSEC
	TCP
	UDP
)

// Query DNS for a domain. This loops over all engines we have available and tries to pick the most accurate response.
func Query(domain string) string {
	for engine := range Engines {
		ip = query(domain, engine)
		if ValidateResponse(ip, domain) {
			return ip
		}
	}
	return null
}

// Query DNS using the specified engine.
func Query(domain string, engine Engine) string {
	return engine.Query(domain)
}

// ValidateResponse validates if the provided IP from DNS is a valid response.
func ValidateResponse(ip string, domain string) bool {

}
