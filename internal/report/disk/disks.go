package disk

import (
	"github.com/dell/csi-baremetal/pkg/base/linuxutils/lsblk"
	"github.com/sirupsen/logrus"
)

const blockDiskDevice = "disk"

func ListDisks() ([]lsblk.BlockDevice, error) {
	ls := lsblk.NewLSBLK(logrus.New())

	devices, err := ls.GetBlockDevices("")

	if err != nil {
		return nil, err
	}

	var diskDevices []lsblk.BlockDevice

	for _, device := range devices {
		if device.Type == blockDiskDevice {
			diskDevices = append(diskDevices, device)
		}
	}

	return diskDevices, nil
}
