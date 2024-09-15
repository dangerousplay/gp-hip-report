package disk

const (
	backupEntryName = "disk-backup"
)

type Backup struct {
	Name string     `xml:"name,attr"`
	List BackupList `xml:"list"`
}

type BackupList struct {
}

func GetBackupInfo() Backup {
	return Backup{Name: backupEntryName}
}
