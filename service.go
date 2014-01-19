package broker

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
	admin *rabbitAdmin
}

func (b *rabbitServiceBroker) Catalog() (Catalog, error) {
	return nil, nil
}

func (b *rabbitServiceBroker) Provision(ProvisioningRequest) (string, error) {
	return "", nil
}

func (b *rabbitServiceBroker) Deprovision(ProvisioningRequest) error {
	return nil
}

func (b *rabbitServiceBroker) Bind(BindingRequest) (Credentials, string, error) {
	return nil, "", nil

}

func (b *rabbitServiceBroker) Unbind(BindingRequest) error {
	return nil
}
