package broker

// The BrokerService defines the internal API used by the broker's HTTP endpoints.
type BrokerService interface {

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

const (
	// Raised by Broker Service if service instance or service instance binding already exists
	ErrCodeConflict = 10
	// Raised by Broker Service if service instance or service instance binding cannot be found
	ErrCodeGone = 20
	// Raised by Broker Service for any other issues
	ErrCodeOther = 99
)

type BrokerServiceError interface {
	Code() int
	Error() string
}

// See http://docs.cloudfoundry.com/docs/running/architecture/services/api.html#provisioning
type ProvisioningRequest struct {
	InstanceId string `json:"-"`
	ServiceId  string `json:"service_id"`
	PlanId     string `json:"plan_id"`
	OrgId      string `json:"organization_guid"`
	SpaceId    string `json:"space_guid"`
}

// See http://docs.cloudfoundry.com/docs/running/architecture/services/api.html#binding
type BindingRequest struct {
	InstanceId string `json:"-"`
	BindingId  string `json:"-"`
	ServiceId  string `json:"service_id"`
	PlanId     string `json:"plan_id"`
	AppId      string `json:"app_guid"`
}

type Credentials map[string]interface{}

// See http://docs.cloudfoundry.com/docs/running/architecture/services/api.html#catalog-mgmt
type Catalog struct {
	Services []Service `json:"services"`
}

// See http://docs.cloudfoundry.com/docs/running/architecture/services/api.html#catalog-mgmt
type Service struct {
	Id          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Bindable    bool                   `json:"bindable"`
	Tags        []string               `json:"tags,omitempty"`
	Requires    []string               `json:"requires,omitempty"`
	Plans       []Plan                 `json:"plans"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// See http://docs.cloudfoundry.com/docs/running/architecture/services/api.html#catalog-mgmt
type Plan struct {
	Id          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// Other types
type BrokerError struct {
	Description string
}
