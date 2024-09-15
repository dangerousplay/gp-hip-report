package dlp

const dlpEntryName = "data-loss-prevention"

type Report struct {
	Name string `xml:"name,attr"`
	List List   `xml:"list"`
}

type List struct {
}

func GetDlpInfo() Report {
	return Report{Name: dlpEntryName}
}
