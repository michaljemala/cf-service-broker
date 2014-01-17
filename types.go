package broker

// Types based on: http://docs.cloudfoundry.com/docs/running/architecture/services/api.html#catalog-mgmt

type common struct {
	Id          string
	Name        string
	Description string
	Metadata    map[string]interface{}
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
