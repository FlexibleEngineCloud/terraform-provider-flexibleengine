package instances

import "github.com/chnsz/golangsdk"

type Instance struct {
	ID             string     `json:"id"`
	Name           string     `json:"name"`
	Version        string     `json:"version"`
	FlavorRef      string     `json:"flavorRef"`
	Status         string     `json:"status"`
	CurrentTask    string     `json:"currentTask"`
	Network        Network    `json:"network"`
	MrsCluster     MrsCluster `json:"mrsCluster"`
	Created        string     `json:"created"`
	Updated        string     `json:"updated"`
	InnerEndPoint  string     `json:"innerEndPoint"`
	PublicEndPoint string     `json:"publicEndPoint"`
}

type Network struct {
	VpcId           string   `json:"vpcId"`
	SubnetId        string   `json:"subnetId"`
	SecurityGroupId string   `json:"securityGroupId"`
	AvailableZone   string   `json:"availableZone"`
	PublicIP        PublicIP `json:"publicIP"`
}

type PublicIP struct {
	BindType string `json:"bindType"`
	EipId    string `json:"eipId"`
}

type MrsCluster struct {
	Id string `json:"id"`
}

type instanceResult struct {
	golangsdk.Result
}

// Extract will get the Instance object out of the commonResult object.
func (r instanceResult) Extract() (*Instance, error) {
	var s Instance
	err := r.ExtractInto(&s)
	return &s, err
}

func (r instanceResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoStructPtr(v, "instance")
}

type CreateResult struct {
	instanceResult
}

type GetResult struct {
	instanceResult
}

type DeleteResult struct {
	golangsdk.ErrResult
}
