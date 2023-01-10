package dns

import (
	"strconv"
	"testing"
)

func TestQuery(t *testing.T) {
	t.Logf("HERE!")
	var ips = Query("google.com")
	for ip := range ips {
		t.Errorf(strconv.Itoa(ip))
	}
}
