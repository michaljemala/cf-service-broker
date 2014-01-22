package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"rabbitmq-service-broker/broker"
)

const (
	version = "1.0.0"
)

func main() {
	o := parseAndProcessOptions()
	b, err := broker.NewBroker(o)
	if err != nil {
		log.Fatal(err)
	}
	h := broker.NewHandler(b)
	r := broker.NewRouter(h)
	go func() {
		addr := fmt.Sprintf("%v:%v", o.Host, o.Port)
		log.Printf("Broker started: Listening at [%v]", addr)
		err := http.ListenAndServe(addr, r)
		if err != nil {
			log.Fatal(err)
		}
	}()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	log.Print("Broker shutdown gracefully")
}

func parseAndProcessOptions() broker.Options {
	opts := broker.Options{}

	flag.StringVar(&opts.Host, "host", "", "")
	flag.IntVar(&opts.Port, "port", 9999, "")
	flag.StringVar(&opts.Username, "user", "admin", "")
	flag.StringVar(&opts.Password, "pass", "secret", "")
	flag.StringVar(&opts.CatalogFile, "catalog", "", "")
	flag.BoolVar(&opts.Debug, "D", false, "")
	flag.StringVar(&opts.LogFile, "L", "", "")
	flag.BoolVar(&opts.Trace, "V", false, "")
	flag.StringVar(&opts.PidFile, "P", "", "")

	flag.StringVar(&opts.RabbitHost, "rabbit-host", "127.0.0.1", "")
	flag.IntVar(&opts.RabbitPort, "rabbit-port", 5672, "")
	flag.IntVar(&opts.RabbitMgmtPort, "rabbit-mgmt-port", 15672, "")
	flag.StringVar(&opts.RabbitUsername, "rabbit-user", "guest", "")
	flag.StringVar(&opts.RabbitPassword, "rabbit-pass", "guest", "")
	flag.BoolVar(&opts.RabbitTrace, "R", false, "")

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

	return opts
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
