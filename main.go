package main

import (
	"fmt"
	"log"
	"net/http"
	"rabbitmq-service-broker/broker"
)

const (
	brokerHost = ""
	brokerPort = "9999"
)

func main() {
	admin, err := broker.NewAdmin("http://localhost:15672", "guest", "guest")
	if err != nil {
		log.Fatal(err)
	}

	router := broker.NewRouter(broker.NewHandler(broker.NewBroker(admin)))

	addr := fmt.Sprintf("%v:%v", brokerHost, brokerPort)
	err = http.ListenAndServe(addr, router)
	if err != nil {
		log.Fatal(err)
	}
}
