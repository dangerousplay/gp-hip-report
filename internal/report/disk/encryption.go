package disk

import (
	"github.com/hashicorp/go-hclog"
	"gp-hip-report/internal/utils"
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
		logger.Warn("Failed to check if cryptsetup is installed: %s", err)
	}

	cryptSetupProd := &EncryptionProductInfoProd{
		Name:    cryptsetup,
		Version: parseCryptSetupVersion(cryptVersion),
		Vendor:  cryptsetupVendor,
	}

	encryptedDrives := EncryptionEntry{
		ProductInfo: EncryptionProductInfo{
			Product: cryptSetupProd,
			Drives:  EncryptionDrives{},
		},
	}

	unencryptedDrives := EncryptionEntry{
		ProductInfo: EncryptionProductInfo{
			Drives: EncryptionDrives{},
		},
	}

	for _, disk := range disks {
		for _, child := range disk.Children {
			diskName := path.Base(child.Name)
			encrypted := IsEncrypted(diskName)
			_, err = os.Stat(child.MountPoint)

			if len(child.MountPoint) < 1 {
				continue
			}

			if err != nil {
				continue
			}

			encDrive := EncryptionDrive{
				MountPoint: child.MountPoint,
			}

			if encrypted {
				encDrive.State = StateEncrypted
				encryptedDrives.ProductInfo.Drives.Entries = append(encryptedDrives.ProductInfo.Drives.Entries, encDrive)
			} else {
				encDrive.State = StateUnencrypted
				unencryptedDrives.ProductInfo.Drives.Entries = append(unencryptedDrives.ProductInfo.Drives.Entries, encDrive)
			}
		}
	}

	if len(encryptedDrives.ProductInfo.Drives.Entries) > 0 {
		enc.List.Entries = append(enc.List.Entries, encryptedDrives)
	}

	if len(unencryptedDrives.ProductInfo.Drives.Entries) > 0 {
		enc.List.Entries = append(enc.List.Entries, unencryptedDrives)
	}

	return &enc, nil
}

func parseCryptSetupVersion(output string) string {
	parts := strings.Split(output, " ")

	if len(parts) < 2 {
		return ""
	}

	return parts[1]
}
