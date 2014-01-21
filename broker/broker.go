package broker

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"log"
)

// The ServiceBroker defines the internal API used by the broker's HTTP endpoints.
type ServiceBroker interface {

	// Exposes the catalog of services managed by this broker.
	// Returns the exposed catalog.
	Catalog() (Catalog, error)

	// Creates a service instance of a specified service and plan.
	// Returns the optional management URL.
	Provision(ProvisioningRequest) (string, error)

	// Removes created service instance.
	Deprovision(ProvisioningRequest) error

	// Binds to specified service instance.
	// Returns  credentials necessary to establish connection to this
	// service instance as well as optional syslog drain URL.
	Bind(BindingRequest) (Credentials, string, error)

	// Removes created binding.
	Unbind(BindingRequest) error
}

// RabbitMQ Service Broker impl
type rabbitServiceBroker struct {
	opts  Options
	admin *rabbitAdmin
}

func NewBroker(opts Options) (*rabbitServiceBroker, error) {
	url := fmt.Sprintf("http://%v:%v", opts.RabbitHost, opts.RabbitPort)
	admin, err := NewAdmin(url, opts.RabbitUsername, opts.RabbitPassword)
	if err != nil {
		return nil, err
	}
	return &rabbitServiceBroker{opts, admin}, nil
}

func (b *rabbitServiceBroker) Catalog() (Catalog, error) {
	return Catalog{
		Services: []Service{
			Service{
				Id:          "rabbitmq",
				Name:        "RabbitMQ",
				Description: "RabbitMQ Message Broker",
				Bindable:    true,
				Tags:        []string{"rabbitmq", "messaging"},
				Plans: []Plan{
					Plan{
						Id:          "simple",
						Name:        "Simple RabbitMQ Plan",
						Description: "Simple RabbitMQ plan represented as a unique broker's vhost.",
					},
				},
			},
		},
	}, nil
}

func (b *rabbitServiceBroker) Provision(pr ProvisioningRequest) (string, error) {
	vhost := pr.Id
	if err := b.admin.createVhost(vhost, false); err != nil {
		return "", err
	}
	log.Printf("Broker: Virtual host created: [%v]", vhost)

	username, password := vhost, generatePassword(vhost)
	if err := b.admin.createUser(username, password); err != nil {
		b.admin.deleteVhost(vhost)
		return "", err
	}
	log.Printf("Broker: Management user created: [%v]", username)

	if err := b.admin.grantAllPermissionsIn(username, vhost); err != nil {
		b.admin.deleteUser(username)
		b.admin.deleteVhost(vhost)
		return "", err
	}
	log.Printf("Broker: All permissions granted to management user: [%v]", username)

	dashboardUrl := fmt.Sprintf("http://%v:%v/#/login/%v/%v", b.opts.RabbitHost, b.opts.RabbitPort, username, password)
	log.Printf("Broker: Dasboard URL generated: [%v]", dashboardUrl)

	return dashboardUrl, nil
}

func (b *rabbitServiceBroker) Deprovision(pr ProvisioningRequest) error {
	return nil
}

func (b *rabbitServiceBroker) Bind(br BindingRequest) (Credentials, string, error) {
	return nil, "", nil

}

func (b *rabbitServiceBroker) Unbind(br BindingRequest) error {
	return nil
}

func generatePassword(str string) string {
	hash := sha1.New().Sum([]byte(str))
	return base64.StdEncoding.EncodeToString(hash)
}
