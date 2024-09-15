package report

import (
	"github.com/zcalusic/sysinfo"
	"gp-hip-report/internal/report/network"
	"runtime"
	"strings"
)

const (
	defaultClientVersion = "5.1.5-8"
	defaultHostId        = "deadbeef-dead-beef-dead-beefdeadbeef"
	hostInfoEntryName    = "host-info"
	internalDomain       = ".internal"
)

type HostEntry struct {
	Name          string `xml:"name,attr"`
	ClientVersion string `xml:"client-version"`
	OS            string `xml:"os"`
	OSVendor      string `xml:"os-vendor"`
	Domain        string `xml:"domain"`
	HostName      string `xml:"host-name"`
	HostID        string `xml:"host-id"`

	Network *network.Interfaces `xml:"network-interface"`
}

func GetHostInformation(computer, domain string) (HostEntry, error) {
	var osName string
	var vendor string

	switch runtime.GOOS {
	case "linux":
		var si sysinfo.SysInfo
		si.GetSysInfo()

		vendor = strings.Title(runtime.GOOS)
		osName = vendor + " " + si.OS.Name
	}

	networkInfo, err := network.GetNetworkInterfaces()

	return HostEntry{
		Name:          hostInfoEntryName,
		ClientVersion: defaultClientVersion,
		OS:            osName,
		OSVendor:      vendor,
		Domain:        domain + internalDomain,
		HostName:      computer,
		HostID:        defaultHostId,

		Network: networkInfo,
	}, err
}
