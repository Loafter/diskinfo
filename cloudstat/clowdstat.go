package cloudstat
type DiskInfo struct {
	Name  string
	Free  uint64
	Total uint64
}

type HealthData struct {
	ServerName string
	DisksInfo  []DiskInfo
}