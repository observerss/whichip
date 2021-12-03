package main

import (
	"fmt"
	"net"
	"strings"
)

func usableIPNets() []*net.IPNet {
	result := make([]*net.IPNet, 0)
	ifaces, err := net.Interfaces()
	if err != nil {
		fmt.Print(fmt.Errorf("localAddresses: %+v\n", err.Error()))
		return nil
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			fmt.Print(fmt.Errorf("localAddresses: %+v\n", err.Error()))
			continue
		}
		for _, a := range addrs {
			switch v := a.(type) {

			case *net.IPNet:
				// 1) not loopback
				// 2) is IPv4 address
				// 3) not a DHCP failure address
				if (!v.IP.IsLoopback() &&
					strings.Count(v.IP.String(), ":") < 2) &&
					!strings.HasPrefix(v.IP.String(), "169.254") {
					result = append(result, v)
				}
			}

		}
	}
	return result
}
