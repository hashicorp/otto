// localaddr is a helper library for allocating local IP addresses.
//
// localaddr does its best to ensure that the IP addresses allocated are
// both predictable (if possible) and on a subnet that isn't in use.
package localaddr

import (
	"fmt"
	"log"
	"net"
)

// RFC 1918
var privateCIDR = []string{"172.16.0.0/12", "10.0.0.0/8", "192.168.0.0/16"}
var privateIPNets = make([]*net.IPNet, len(privateCIDR))

func init() {
	for i, cidr := range privateCIDR {
		_, ipnet, err := net.ParseCIDR(cidr)
		if err != nil {
			panic(err)
		}

		privateIPNets[i] = ipnet
	}
}

// UsableSubnet returns a /24 CIDR block of usable network addresses in the
// RFC private address space that also isn't in use by any network interface
// on this machine currently.
func UsableSubnet() (*net.IPNet, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	// First find all the taken private IP spaces
	taken := make([]*net.IPNet, 0, 5)
	for _, i := range interfaces {
		addrs, err := i.Addrs()
		if err != nil {
			return nil, err
		}

		for _, addr := range addrs {
			if addr.Network() != "ip+net" {
				log.Printf("[DEBUG] ignoring non ip+net addr: %s", addr)
				continue
			}

			// Parse the CIDR block of this interface
			ip, ipnet, err := net.ParseCIDR(addr.String())
			if err != nil {
				return nil, err
			}

			// We only do IPv4 for now
			if ip.To4() == nil {
				log.Printf("[DEBUG] ignoring non IPv4 addr: %s", addr)
				continue
			}

			// If the addr isn't even in the private IP space, then ignore it
			private := false
			for _, ipnet := range privateIPNets {
				if ipnet.Contains(ip) {
					private = true
					break
				}
			}
			if !private {
				log.Printf("[DEBUG] ignoring non-private IP space: %s", addr)
				continue
			}

			log.Printf("[DEBUG] occupied private IP space: %s", addr)
			taken = append(taken, ipnet)
		}
	}

	// Now go through and find a space that we can use
	for _, ipnet := range privateIPNets {
		// Get the first IP, and add one since we always want to start at
		// x.x.1.x since x.x.0.x sometimes means something special.
		ip := ipnet.IP.Mask(ipnet.Mask)
		for ip[2] = 1; ip[2] <= 255; ip[2]++ {
			// Determine if one of our taken CIDRs has this IP. We can
			// probably do this way more efficiently but this is fine for now.
			bad := false
			for _, ipnet := range taken {
				if ipnet.Contains(ip) {
					bad = true
					break
				}
			}
			if bad {
				continue
			}

			// This is a good address space
			_, ipnet, err := net.ParseCIDR(fmt.Sprintf("%s/24", ip))
			if err != nil {
				return nil, err
			}

			log.Printf("[DEBUG] found usable subnet: %s", ipnet)
			return ipnet, nil
		}
	}

	return nil, fmt.Errorf("no usable subnet found in private space")
}
