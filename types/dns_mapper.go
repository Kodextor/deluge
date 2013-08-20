// Steve Phillips / elimisteve
// 2013.04.28

package types

import (
	"net"
)

func NewDNSMapper() *DNSMapper {
	return &DNSMapper{
		ips: map[string]*net.TCPAddr{},
	}
}

type DNSMapper struct {
	ips map[string]*net.TCPAddr
}

