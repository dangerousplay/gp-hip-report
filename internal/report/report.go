package report

import (
	"context"
	"encoding/xml"
	"github.com/hashicorp/go-hclog"
	"gp-hip-report/internal/report/antithreat"
	"gp-hip-report/internal/report/disk"
	"gp-hip-report/internal/report/dlp"
	"gp-hip-report/internal/report/network"
	"gp-hip-report/internal/report/patch"
	"net/url"
	"time"
)

const (
	hipReportVersion = 4
	hipReportName    = "hip-report"
	timeFormat       = "01/02/2006 15:04:05"

	computerParam = "computer"
	domainParam   = "domain"
	userParam     = "user"
)

type Categories struct {
	Entries []interface{} `xml:"entry"`
}

type HIPReport struct {
	Name             string   `xml:"name,attr"`
	XMLName          xml.Name `xml:"hip-report"`
	GenerateTime     string   `xml:"generate-time"`
	HIPReportVersion int      `xml:"hip-report-version"`

	MD5Sum   string `xml:"md5-sum"`
	Username string `xml:"user-name"`
	Domain   string `xml:"domain"`
	HostName string `xml:"host-name"`
	HostID   string `xml:"host-id"`

	IpAddress   string `xml:"ip-address"`
	IpV6Address string `xml:"ipv6-address"`

	Categories Categories `xml:"categories"`
}

func GenerateReport(ctx context.Context, cookie, md5, clientIpv4, clientIpv6 string) (HIPReport, error) {
	datetime := time.Now().Format(timeFormat)

	params, err := url.ParseQuery(cookie)

	if err != nil {
		return HIPReport{}, err
	}

	user := params.Get(userParam)
	domain := params.Get(domainParam)
	computer := params.Get(computerParam)

	encryption, err := disk.GetDiskEncryptionInfo()

	if err != nil {
		hclog.Default().Error("Failed to get disk encryption", err)
	}

	diskBackup := disk.GetBackupInfo()

	hostInfo, err := GetHostInformation(computer, domain)

	if err != nil {
		hclog.Default().Error("Failed to get host information", err)
	}

	firewall, err := network.GetFirewallInfo()

	if err != nil {
		hclog.Default().Error("Failed to get firewall information", err)
	}

	antivirus, err := antithreat.GetAntiMalware(ctx)

	if err != nil {
		hclog.Default().Error("Failed to get antivirus information", err)
	}

	dlpInfo := dlp.GetDlpInfo()

	patchManagement, err := patch.GetPatchManagement()

	if err != nil {
		hclog.Default().Error("Failed to get antivirus information", err)
	}

	return HIPReport{
		Name:             hipReportName,
		GenerateTime:     datetime,
		HIPReportVersion: hipReportVersion,
		MD5Sum:           md5,
		Username:         user,
		Domain:           domain,
		HostName:         computer,
		HostID:           defaultHostId,
		IpAddress:        clientIpv4,
		IpV6Address:      clientIpv6,
		Categories: Categories{
			Entries: []interface{}{
				hostInfo,
				encryption,
				firewall,
				patchManagement,
				antivirus,
				dlpInfo,
				diskBackup,
			},
		},
	}, err
}
