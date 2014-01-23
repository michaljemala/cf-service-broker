package rabbitmq

import (
	"cf-service-broker/broker"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"log"
)

// BrokerService implementation for RabbitMQ Server
type rabbitService struct {
	opts  Options
	admin *rabbitAdmin
}

func NewBrokerService(opts Options) (*rabbitService, error) {
	url := fmt.Sprintf("http://%v:%v", opts.MgmtHost, opts.MgmtPort)
	adm, err := newRabbitAdmin(url, opts.MgmtUser, opts.MgmtPass)
	if err != nil {
		return nil, err
	}
	return &rabbitService{opts, adm}, nil
}

func (b *rabbitService) Catalog() (broker.Catalog, error) {
	// TODO: Read catalog from file
	return broker.Catalog{
		Services: []broker.Service{
			broker.Service{
				Id:          "rabbitmq",
				Name:        "RabbitMQ",
				Description: "RabbitMQ Message Broker",
				Bindable:    true,
				Tags:        []string{"rabbitmq", "messaging"},
				Plans: []broker.Plan{
					broker.Plan{
						Id:          "simple",
						Name:        "Simple RabbitMQ Plan",
						Description: "Simple RabbitMQ plan represented as a unique broker's vhost.",
					},
				},
			},
		},
	}, nil
}

func (b *rabbitService) Provision(pr broker.ProvisioningRequest) (string, error) {
	vhost := pr.Id
	if err := b.admin.createVhost(vhost, false); err != nil {
		return "", err
	}
	log.Printf("Service: Virtual host created: [%v]", vhost)

	username := fmt.Sprintf("m-%v", vhost)
	password := generatePassword(pr.Id)
	if err := b.admin.createUser(username, password); err != nil {
		b.admin.deleteVhost(vhost)
		return "", err
	}
	log.Printf("Service: Management user created: [%v]", username)

	if err := b.admin.grantAllPermissionsIn(username, vhost); err != nil {
		b.admin.deleteUser(username)
		b.admin.deleteVhost(vhost)
		return "", err
	}
	log.Printf("Service: All permissions granted to management user: [%v]", username)

	dashboardUrl := fmt.Sprintf("http://%v:%v/#/login/%v/%v", b.opts.Host, b.opts.MgmtPort, username, password)
	log.Printf("Service: Dasboard URL generated: [%v]", dashboardUrl)

	return dashboardUrl, nil
}

func (b *rabbitService) Deprovision(pr broker.ProvisioningRequest) error {
	vhost := pr.Id
	username := fmt.Sprintf("m-%v", vhost)
	if err := b.admin.deleteUser(username); err != nil {
		return err
	}
	log.Printf("Service: Management user deleted: [%v]", username)

	//TODO:Should close existing connections from user 'username'???

	if err := b.admin.deleteVhost(vhost); err != nil {
		return err
	}
	log.Printf("Service: Virtual host deleted: [%v]", vhost)

	return nil
}

func (b *rabbitService) Bind(br broker.BindingRequest) (broker.Credentials, string, error) {
	vhost := br.Id
	username := fmt.Sprintf("u-%v", vhost)
	password := generatePassword(br.Id + br.ServiceId)
	if err := b.admin.createUser(username, password); err != nil {
		return nil, "", err
	}
	log.Printf("Service: User created: [%v]", username)

	if err := b.admin.grantAllPermissionsIn(username, vhost); err != nil {
		b.admin.deleteUser(username)
		return nil, "", err
	}
	log.Printf("Service: All permissions granted for [%v] to user: [%v]", vhost, username)

	amqpUrl := fmt.Sprintf("amqp://%v:%v@%v:%v/%v", username, password, b.opts.Host, b.opts.Port, vhost)
	log.Printf("Service: AMQP URL generated: [%v]", amqpUrl)

	return broker.Credentials{"uri": amqpUrl}, "", nil
}

func (b *rabbitService) Unbind(br broker.BindingRequest) error {
	username := fmt.Sprintf("u-%v", br.Id)
	err := b.admin.deleteUser(username)
	if err != nil {
		return err
	}
	log.Printf("Service: User deleted: [%v]", username)

	//TODO:Should close existing connections from user 'username'???

	return nil
}

func generatePassword(str string) string {
	hash := sha1.New().Sum([]byte(str))
	password := base64.StdEncoding.EncodeToString(hash)
	return password
}
