package broker

type Options struct {
	// Broker Settings
	Host        string
	Port        int
	Username    string
	Password    string
	CatalogFile string
	Debug       bool
	LogFile     string
	Trace       bool
	PidFile     string
	// Rabbit Settings
	RabbitHost     string
	RabbitPort     int
	RabbitUsername string
	RabbitPassword string
	RabbitTrace    bool
}
