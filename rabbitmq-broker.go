package main

import (
	"cf-service-broker/broker"
	"cf-service-broker/rabbitmq"
	"flag"
	"fmt"
	"log"
	"os"
)

const (
	version = "1.0.0"
)

func main() {
	opts := broker.Options{}
	flag.StringVar(&opts.Host, "bh", "", "")
	flag.StringVar(&opts.Host, "broker-host", "", "")
	flag.IntVar(&opts.Port, "br", 9999, "")
	flag.IntVar(&opts.Port, "broker-port", 9999, "")
	flag.StringVar(&opts.Username, "bu", "admin", "")
	flag.StringVar(&opts.Username, "broker-user", "admin", "")
	flag.StringVar(&opts.Password, "bp", "secret", "")
	flag.StringVar(&opts.Password, "broker-password", "secret", "")
	flag.StringVar(&opts.CatalogFile, "broker-catalog", "", "")
	flag.BoolVar(&opts.Debug, "D", false, "")
	flag.StringVar(&opts.LogFile, "L", "", "")
	flag.BoolVar(&opts.Trace, "V", false, "")
	flag.StringVar(&opts.PidFile, "P", "", "")

	rmqOpts := rabbitmq.Options{}
	flag.StringVar(&rmqOpts.Host, "rh", "127.0.0.1", "")
	flag.StringVar(&rmqOpts.Host, "rabbit-host", "127.0.0.1", "")
	flag.IntVar(&rmqOpts.Port, "rr", 5672, "")
	flag.IntVar(&rmqOpts.Port, "rabbit-port", 5672, "")
	flag.StringVar(&rmqOpts.MgmtHost, "rmh", "127.0.0.1", "")
	flag.StringVar(&rmqOpts.MgmtHost, "rabbit-mgmt-host", "127.0.0.1", "")
	flag.IntVar(&rmqOpts.MgmtPort, "rmr", 15672, "")
	flag.IntVar(&rmqOpts.MgmtPort, "rabbit-mgmt-port", 15672, "")
	flag.StringVar(&rmqOpts.MgmtUser, "rmu", "guest", "")
	flag.StringVar(&rmqOpts.MgmtUser, "rabbit-mgmt-user", "guest", "")
	flag.StringVar(&rmqOpts.MgmtPass, "rmp", "guest", "")
	flag.StringVar(&rmqOpts.MgmtPass, "rabbit-mgmt-pass", "guest", "")
	flag.BoolVar(&rmqOpts.Trace, "R", false, "")

	var showHelp, showVersion bool
	flag.BoolVar(&showHelp, "help", false, "")
	flag.BoolVar(&showVersion, "version", false, "")

	flag.Usage = Usage
	flag.Parse()

	if showHelp {
		Usage()
	}
	if showVersion {
		Version()
	}

	brokerService, err := rabbitmq.NewBrokerService(rmqOpts)
	if err != nil {
		log.Fatal(err)
	}

	broker := broker.New(opts, brokerService)
	broker.Start()
}

var versionStr = fmt.Sprintf(`RabbitMQ Service Broker v%v
`, version)

var usageStr = fmt.Sprintf(`RabbitMQ Service Broker v%v

Broker options:
        --host HOST                 Bind to HOST address (default: 0.0.0.0)
        --port PORT                 Use PORT (default: 9999)
        --user USERNAME             User required to authenticate requests (default: admin)
        --pass PASSWORD             Password for the USERNAME user (default: secret)
        --catalog CATALOG           File CATALOG to load the broker's catalog from
    -D                              Enable debugging output
    -L FILE                         File to redirect log output to
    -V                              Trace the incoming service broker's HTTP requests
    -P FILE                         File to store broker's PID to

RabbitMQ options:
        --rabbit-host HOST          Hostname of RabbitMQ server (default: 127.0.0.1)
        --rabbit-port PORT          Port on which RabbitMQ server listens for messages  (default: 5672)
        --rabbit-mgmt-port PORT     Port on which RabbitMQ server listens for management requests (default: 15672)
        --rabbit-user USERNAME      Username of the RabbitMQ server user with 'administrator' tag assigned (default: guest)
        --rabbit-pass PASSWORD      Password for the USERNAME user (default: guest)
    -R                              Trace the outgoing RabbitMQ server management requests

Common options:
        --help                      Show this message
        --version                   Show service broker version
`, version)

func Usage() {
	fmt.Print(usageStr)
	os.Exit(0)
}

func Version() {
	fmt.Print(versionStr)
	os.Exit(0)
}
