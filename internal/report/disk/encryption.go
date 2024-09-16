package disk

import (
	"github.com/cockroachdb/errors"
	"github.com/dell/csi-baremetal/pkg/base/linuxutils/lsblk"
	"github.com/hashicorp/go-hclog"
	"gp-hip-report/internal/utils"
	"io/fs"
	"os"
	"path"
	"strings"
)

type EncryptionState string

const (
	StateEncrypted   EncryptionState = "encrypted"
	StateUnencrypted EncryptionState = "unencrypted"
)

const (
	cryptsetup       = "cryptsetup"
	cryptsetupVendor = "GitLab Inc."

	diskEncryption = "disk-encryption"
)

type Encryption struct {
	Name string         `xml:"name,attr"`
	List EncryptionList `xml:"list"`
}

type EncryptionList struct {
	Entries []EncryptionEntry `xml:"entry"`
}

type EncryptionEntry struct {
	ProductInfo EncryptionProductInfo
}

type EncryptionProductInfo struct {
	Product *EncryptionProductInfoProd `xml:"Prod"`
	Drives  EncryptionDrives           `xml:"drives"`
}

type EncryptionProductInfoProd struct {
	Name    string `xml:"name,attr"`
	Version string `xml:"version,attr"`
	Vendor  string `xml:"vendor,attr"`
}

type EncryptionDrives struct {
	Entries []EncryptionDrive `xml:"entry"`
}

type EncryptionDrive struct {
	MountPoint string          `xml:"drive-name"`
	State      EncryptionState `xml:"enc-state"`
}

func GetDiskEncryptionInfo() (*Encryption, error) {
	logger := hclog.Default().Named("disk-encryption")

	disks, err := ListDisks()

	if err != nil {
		return nil, err
	}

	enc := Encryption{
		Name: diskEncryption,
	}

	_, cryptVersion, err := utils.CheckExists("cryptsetup", []string{"cryptsetup", "--version"})

	if err != nil {
		logger.Warn("Failed to check if cryptsetup is installed", err)
	}

	cryptSetupProd := &EncryptionProductInfoProd{
		Name:    cryptsetup,
		Version: parseCryptSetupVersion(cryptVersion),
		Vendor:  cryptsetupVendor,
	}

	encryptedDrives := &EncryptionEntry{
		ProductInfo: EncryptionProductInfo{
			Product: cryptSetupProd,
			Drives:  EncryptionDrives{},
		},
	}

	unencryptedDrives := &EncryptionEntry{
		ProductInfo: EncryptionProductInfo{
			Drives: EncryptionDrives{},
		},
	}

	var errs error
	for _, disk := range disks {
		for _, child := range disk.Children {
			if len(child.MountPoint) < 1 && len(child.Children) < 1 {
				continue
			}

			err = appendDisk(child, encryptedDrives, unencryptedDrives)
			errs = errors.CombineErrors(errs, err)
		}
	}

	if errs != nil {
		hclog.Default().Warn("failed to check disk encryption", errs)
	}

	if len(encryptedDrives.ProductInfo.Drives.Entries) > 0 {
		enc.List.Entries = append(enc.List.Entries, *encryptedDrives)
	}

	if len(unencryptedDrives.ProductInfo.Drives.Entries) > 0 {
		enc.List.Entries = append(enc.List.Entries, *unencryptedDrives)
	}

	return &enc, nil
}

func appendDisk(disk lsblk.BlockDevice, encryptedDrives *EncryptionEntry, unencryptedDrives *EncryptionEntry) (errs error) {
	diskName := path.Base(disk.Name)
	encrypted := IsEncrypted(diskName)

	if len(disk.Children) > 0 {
		for _, child := range disk.Children {
			err := appendDisk(child, encryptedDrives, unencryptedDrives)
			errs = errors.CombineErrors(errs, err)
		}
	}

	if len(disk.MountPoint) < 1 {
		hclog.Default().Debug("no mount point defined for disk", "disk", disk.Name)
		return
	}

	_, err := os.Stat(disk.MountPoint)
	if errors.Is(err, fs.ErrNotExist) {
		return
	} else if err != nil {
		errs = errors.CombineErrors(errs, err)
		return
	}

	encDrive := EncryptionDrive{
		MountPoint: disk.MountPoint,
	}

	if encrypted {
		encDrive.State = StateEncrypted
		encryptedDrives.ProductInfo.Drives.Entries = append(encryptedDrives.ProductInfo.Drives.Entries, encDrive)
	} else {
		encDrive.State = StateUnencrypted
		unencryptedDrives.ProductInfo.Drives.Entries = append(unencryptedDrives.ProductInfo.Drives.Entries, encDrive)
	}

	return
}

func parseCryptSetupVersion(output string) string {
	parts := strings.Split(output, " ")

	if len(parts) < 2 {
		return ""
	}

	return parts[1]
}
