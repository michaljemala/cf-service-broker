package broker

type Options struct {
	Host        string
	Port        int
	Username    string
	Password    string
	CatalogFile string
	Debug       bool
	LogFile     string
	Trace       bool
	PidFile     string
}
