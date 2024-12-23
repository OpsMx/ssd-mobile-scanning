package static

type UploadResult struct {
	Analyzer string `json:"analyzer"`
	Status   string `json:"status"`
	Hash     string `json:"hash"`
	ScanType string `json:"scan_type"`
	FileName string `json:"file_name"`
}

type ScanLog struct {
	Timestamp string  `json:"timestamp"`
	Status    string  `json:"status"`
	Exception *string `json:"exception"` // Using pointer to handle potential null values
}

type ScanContent struct {
	Analyzer    string    `json:"ANALYZER"`
	ScanType    string    `json:"SCAN_TYPE"`
	FileName    string    `json:"FILE_NAME"`
	AppName     string    `json:"APP_NAME"`
	PackageName string    `json:"PACKAGE_NAME"`
	VersionName string    `json:"VERSION_NAME"`
	MD5         string    `json:"MD5"`
	Timestamp   string    `json:"TIMESTAMP"`
	ScanLogs    []ScanLog `json:"SCAN_LOGS"`
}

type ScanLogsResult struct {
	Logs []*ScanLog `json:"logs"`
}

type ScansResult struct {
	Content  []ScanContent `json:"content"`
	Count    int           `json:"count"`
	NumPages int           `json:"num_pages"`
}

type DeleteResult struct {
	Deleted string `json:"deleted"`
}
