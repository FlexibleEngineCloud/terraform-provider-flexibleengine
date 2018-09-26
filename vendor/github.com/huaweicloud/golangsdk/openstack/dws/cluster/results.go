package cluster

import "github.com/huaweicloud/golangsdk"

type Cluster struct {
	Status           string              `json:"status"`
	SubStatus        string              `json:"sub_status"`
	Updated          string              `json:"updated"`
	Endpoints        Endpoints           `json:"endPoints"`
	Name             string              `json:"name"`
	NumberOfNode     int                 `json:"number_of_node"`
	AvailabilityZone string              `json:"availability_zone"`
	SubnetID         string              `json:"subnet_id"`
	PublicEndpoints  PublicEndpoints     `json:"public_endpoints"`
	Created          string              `json:"created"`
	SecurityGroupID  string              `json:"security_group_id"`
	Port             int                 `json:"port"`
	NodeType         string              `json:"node_type"`
	Version          string              `json:"version"`
	PublicIp         PublicIp            `json:"public_ip"`
	FailedReasons    map[string]FailInfo `json:"failed_reasons"`
	VpcID            string              `json:"vpc_id"`
	TaskStatus       string              `json:"task_status"`
	UserName         string              `json:"user_name"`
	ID               string              `json:"id"`
}

type FailInfo struct {
	ErrorCode string `json:"error_code"`
	ErrorMsg  string `json:"error_msg"`
}

type PublicIp struct {
	EipID          string `json:"eip_id"`
	PublicBindType string `json:"public_bind_type"`
}

type FailedReasons struct {
	FailInfo FailInfo `json:"fail_info"`
}

type Endpoints struct {
	ConnectInfo string `json:"connect_info"`
	JdbcUrl     string `json:"jdbc_url"`
}

type PublicEndpoints struct {
	PublicConnectInfo string `json:"public_connect_info"`
	JdbcUrl           string `json:"jdbc_url"`
}

type GetResult struct {
	golangsdk.Result
}

func (r GetResult) Extract() (*Cluster, error) {
	o := &Cluster{}
	return o, r.ExtractIntoStructPtr(o, "cluster")
}

type CreateRsp struct {
	ID string `json:"id"`
}

type CreateResult struct {
	golangsdk.Result
}

func (r CreateResult) Extract() (*CreateRsp, error) {
	o := &CreateRsp{}
	return o, r.ExtractIntoStructPtr(o, "cluster")
}

type DeleteResult struct {
	golangsdk.ErrResult
}
