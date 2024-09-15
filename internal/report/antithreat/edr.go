package antithreat

import (
	"context"
	"encoding/xml"
	"fmt"
	"github.com/cockroachdb/errors"
	"github.com/coreos/go-systemd/v22/dbus"
	"gp-hip-report/internal/constants"
	"gp-hip-report/internal/utils"
	"strings"
	"time"
)

const (
	falconSensorService  = "falcon-sensor.service"
	antiMalwareEntryName = "anti-malware"
	activeStateKey       = "ActiveState"
	loadStateKey         = "LoadState"
	edrFalcon            = "falcon-sensor"
)

type loadState string

const (
	loadStateNotFound loadState = "not-found"
	loadStateLoaded   loadState = "loaded"
)

type unitState string

const (
	activeState   unitState = "active"
	inactiveState unitState = "inactive"
	failedState   unitState = "failed"
	exitedState   unitState = "exited"
)

var knownEdrs = map[string]Product{
	edrFalcon: {
		Name:   "CrowdStrike Falcon",
		Vendor: constants.VendorCrowdStrike,
	},
}

type AntiMalwareReport struct {
	XMLName xml.Name `xml:"entry"`
	Name    string   `xml:"name,attr"`
	List    []Entry  `xml:"list>entry"`
}

type Entry struct {
	ProductInfo ProductInfo `xml:"ProductInfo"`
}

type ProductInfo struct {
	Product            Product `xml:"Prod"`
	RealTimeProtection string  `xml:"real-time-protection"`
	LastFullScanTime   string  `xml:"last-full-scan-time"`
}

type Product struct {
	Name     string `xml:"name,attr"`
	Version  string `xml:"version,attr"`
	DefVer   string `xml:"defver,attr"`
	ProdType string `xml:"prodType,attr"`
	EngVer   string `xml:"engver,attr"`
	OsType   string `xml:"osType,attr"`
	Vendor   string `xml:"vendor,attr"`
	DateDay  string `xml:"dateday,attr"`
	DateYear string `xml:"dateyear,attr"`
	DateMon  string `xml:"datemon,attr"`
}

func GetAntiMalware(ctx context.Context) (*AntiMalwareReport, error) {
	report := &AntiMalwareReport{
		Name: antiMalwareEntryName,
		List: []Entry{},
	}

	falcon, err := checkFalconSensor(ctx)

	if err != nil {
		return nil, errors.Wrap(err, "failed to check falcon sensor")
	}

	if falcon != nil {
		report.List = append(report.List, Entry{
			ProductInfo: *falcon,
		})
	}

	return report, nil
}

func checkFalconSensor(ctx context.Context) (*ProductInfo, error) {
	conn, err := dbus.NewWithContext(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to systemd")
	}
	defer conn.Close()

	unitStatus, err := conn.GetUnitPropertiesContext(ctx, falconSensorService)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get falcon sensor unit")
	}

	if unitStatus[loadStateKey] != string(loadStateLoaded) {
		return nil, nil
	}

	state := unitStatus[activeStateKey]
	if state != string(activeState) {
		return nil, nil
	}

	_, versionOutput, err := utils.CheckExists("falconctl", []string{"falconctl", "-g", "--version"})

	if err != nil {
		return nil, errors.Wrap(err, "failed to check falcon version")
	}

	version := strings.TrimSpace(strings.Replace(versionOutput, "version = ", "", -1))

	now := time.Now()

	falcon := knownEdrs[edrFalcon]

	product := &ProductInfo{
		Product: Product{
			Name:     falcon.Name,
			Version:  version,
			DefVer:   "",
			EngVer:   now.Format("2006.01.02"),
			ProdType: "1",
			OsType:   "1",
			Vendor:   falcon.Vendor,
			DateDay:  fmt.Sprint(now.Day()),
			DateYear: fmt.Sprint(now.Year()),
			DateMon:  fmt.Sprint(int(now.Month())),
		},
		RealTimeProtection: constants.Yes,
		LastFullScanTime:   constants.NotApplicable,
	}

	return product, nil
}
