package gpu

type OrgClusterGroupsResponse struct {
	BillingAccountID string        `json:"billingAccountId"`
	Clusters         []Cluster     `json:"clusters"`
	RequestStatus    RequestStatus `json:"requestStatus"`
}

type Cluster struct {
	Cluster          string `json:"cluster"`
	GpuType          string `json:"gpuType"`
	InstanceType     string `json:"instanceType"`
	MaxInstances     int    `json:"maxInstances"`
	CurrentInstances int    `json:"currentInstances"`
}

type RequestStatus struct {
	StatusCode string `json:"statusCode"`
}
