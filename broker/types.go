package broker

// Types based on http://docs.cloudfoundry.com/docs/running/architecture/services/api.html

type ProvisioningRequest struct {
	Id        string `json:"-"`
	ServiceId string `json:"service_id"`
	PlanId    string `json:"plan_id"`
	OrgId     string `json:"organization_guid"`
	SpaceId   string `json:"space_guid"`
}

type BindingRequest struct {
	Id         string `json:"-"`
	InstanceId string `json:"-"`
	ServiceId  string `json:"service_id"`
	PlanId     string `json:"plan_id"`
	AppId      string `json:"app_guid"`
}

type Credentials map[string]interface{}

type Catalog struct {
	Services []Service `json:"services"`
}

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

type Plan struct {
	Id          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

//Other types
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
	RabbitMgmtPort int
	RabbitUsername string
	RabbitPassword string
	RabbitTrace    bool
}
