package patch

import (
	"encoding/xml"
)

const (
	patchManagementEntryName = "patch-management"
)

type Management struct {
	Name         string         `xml:"name,attr"`
	List         ManagementList `xml:"list"`
	MissingPatch MissingPatch   `xml:"missing-patches"`
}

type ManagementList struct {
	Entries []ManagementEntry `xml:"entry"`
}

type ManagementEntry struct {
	ProductInfo ManagementProductInfo `xml:"ProductInfo"`
}

type ManagementProductInfo struct {
	XMLName xml.Name          `xml:"ProductInfo"`
	Product ManagementProduct `xml:"Prod"`
	Enabled string            `xml:"is-enabled"`
}

type ManagementProduct struct {
	Vendor  string `xml:"vendor,attr"`
	Name    string `xml:"name,attr"`
	Version string `xml:"version,attr"`
}

func GetPatchManagement() (*Management, error) {
	tools, err := getPatchManagementTools()

	if err != nil {
		return nil, err
	}

	return &Management{
		Name: patchManagementEntryName,
		List: ManagementList{Entries: tools},
	}, nil
}
