package main

import (
	"bitbucket.org/michaljemala/cf-service-broker/broker"
	"bitbucket.org/michaljemala/cf-service-broker/rabbitmq"
	"flag"
	"fmt"
	"log"
	"os"
)

const version = "1.0.0"

func init() {
	flag.BoolVar(&showHelp, "help", false, "")
	flag.BoolVar(&showVersion, "version", false, "")
}

func main() {
	flag.Usage = Usage
	flag.Parse()

	if showHelp {
		Usage()
	}
	if showVersion {
		Version()
	}

	brokerService, err := rabbitmq.New(rabbitmq.Opts)
	if err != nil {
		log.Fatal(err)
	}

	broker := broker.New(broker.Opts, brokerService)
	broker.Start()
}

func Usage() {
	fmt.Print(broker.UsageStr)
	fmt.Print(rabbitmq.UsageStr)
	fmt.Print(usageStr)
	os.Exit(0)
}
func Version() {
	fmt.Print(versionStr)
	os.Exit(0)
}

var (
	showHelp, showVersion bool
	versionStr            = fmt.Sprintf("RabbitMQ Service Broker v%v\n", version)
	usageStr              = fmt.Sprintf(`RabbitMQ Service Broker v%v
Common options:
        --help                         Show this message
        --version                      Show service broker version
`, version)
)
