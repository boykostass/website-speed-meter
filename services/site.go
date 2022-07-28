package services

// SiteInfo - структура для описания данных в базе данных
type SiteInfo struct {
	Site        string `json:"site"`
	Date        string `json:"date"`
	Time        string `json:"time"`
	Delay       string `json:"delay"`
	Performance string `json:"performance"`
}
