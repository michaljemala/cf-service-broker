package rabbitmq

type Options struct {
	Host     string
	Port     int
	MgmtHost string
	MgmtPort int
	MgmtUser string
	MgmtPass string
	Trace    bool
}
