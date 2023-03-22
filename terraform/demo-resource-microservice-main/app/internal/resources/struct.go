package resources

type resource struct {
	Name   string
	Type   string
	Region string
}

type resourceDTO struct {
	Type   string `json:"resource_type"`
	Region string `json:"region"`
}
