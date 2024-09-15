package network

import (
	"encoding/xml"
	"github.com/hashicorp/go-hclog"
	"github.com/life4/genesis/slices"
	"net"
	"strings"
)

var ignoreInterfaces = []string{"lo", "tun", "gpd"}

type Interfaces struct {
	XMLName xml.Name       `xml:"network-interface"`
	Entries []NetworkEntry `xml:"entry"`
}

type NetworkEntry struct {
	Name        string      `xml:"name,attr"`
	Description string      `xml:"description"`
	MacAddress  string      `xml:"mac-address"`
	IPAddress   IPAddresses `xml:"ip-address"`
	IPv6Address IPAddresses `xml:"ipv6-address"`
}

type IPAddresses struct {
	Entries []IPEntry `xml:"entry"`
}

type IPEntry struct {
	Name string `xml:"name,attr"`
}

func GetNetworkInterfaces() (*Interfaces, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	var entries []NetworkEntry

	for _, iface := range interfaces {
		shouldIgnore := slices.Any(ignoreInterfaces, func(ifname string) bool {
			return strings.HasPrefix(iface.Name, ifname)
		})

		if shouldIgnore {
			continue
		}

		ipv6Addrs := IPAddresses{Entries: make([]IPEntry, 0)}
		ipv4Addrs := IPAddresses{Entries: make([]IPEntry, 0)}

		ifIpAddrs, err := iface.Addrs()
		if err != nil {
			hclog.Default().Warn("Error retrieving IPv6 addresses for interface %s: %s", iface.Name, err)
			continue
		}

		for _, addr := range ifIpAddrs {
			ip, _, err := net.ParseCIDR(addr.String())
			if err != nil {
				continue
			}

			if ip.To4() == nil {
				ipv6Addrs.Entries = append(ipv6Addrs.Entries, IPEntry{Name: ip.String()})
			} else {
				ipv4Addrs.Entries = append(ipv4Addrs.Entries, IPEntry{Name: ip.String()})
			}
		}

		entries = append(entries, NetworkEntry{
			Name:        iface.Name,
			Description: iface.Name,
			MacAddress:  iface.HardwareAddr.String(),
			IPAddress:   ipv4Addrs,
			IPv6Address: ipv6Addrs,
		})
	}

	return &Interfaces{
		Entries: entries,
	}, nil
}
