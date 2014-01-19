package broker

// Types based on http://docs.cloudfoundry.com/docs/running/architecture/services/api.html

type ProvisioningRequest struct {
	Id        string
	ServiceId string `json:"service_id"`
	PlanId    string `json:"plan_id"`
	OrgId     string `json:"organization_guid"`
	SpaceId   string `json:"space_guid"`
}

type BindingRequest struct {
	Id        string
	ServiceId string `json:"service_id"`
	PlanId    string `json:"plan_id"`
	AppId     string `json:"app_guid"`
}

type Credentials map[string]interface{}

type Catalog struct {
	services []Service
}

type Service struct {
	common
	Bindable bool
	Tags     []string
	Requires []string
	Plans    []Plan
}

type Plan struct {
	common
}

type common struct {
	Id          string
	Name        string
	Description string
	Metadata    map[string]interface{}
}
