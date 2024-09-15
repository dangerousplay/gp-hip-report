package patch

import (
	"github.com/cockroachdb/errors"
	"github.com/zcalusic/sysinfo"
	"gp-hip-report/internal/constants"
	"gp-hip-report/internal/utils"
	"strings"
)

const (
	aptName = "Advanced Packaging Tool"
)

func getPatchManagementTools() (tools []ManagementEntry, err error) {
	var si sysinfo.SysInfo
	si.GetSysInfo()

	switch si.OS.Vendor {
	case "ubuntu":
		fallthrough
	case "debian":
		apt, e := getAptTool()

		if e != nil {
			err = e
			return
		}

		tools = append(tools, ManagementEntry{*apt})
	}

	return
}

func getAptVersion() (string, error) {
	_, versionOut, err := utils.CheckExists("apt", []string{"apt", "--version"})

	if err != nil {
		return "", errors.Wrap(err, "failed to get apt version")
	}

	parts := strings.Split(versionOut, " ")

	if len(parts) < 2 {
		return "", nil
	}

	return parts[1], nil
}

func getAptTool() (*ManagementProductInfo, error) {
	version, err := getAptVersion()

	if err != nil {
		return nil, err
	}

	return &ManagementProductInfo{
		Product: ManagementProduct{
			Vendor:  constants.VendorGNU,
			Name:    aptName,
			Version: version,
		},
		Enabled: constants.Yes,
	}, nil
}
