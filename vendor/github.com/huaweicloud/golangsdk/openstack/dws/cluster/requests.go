package cluster

import (
	"log"

	"github.com/huaweicloud/golangsdk"
)

type PublicIpOpts struct {
	EipID          string `json:"eip_id,omitempty"`
	PublicBindType string `json:"public_bind_type,omitempty"`
}

type CreateOpts struct {
	Name             string        `json:"name" required:"true"`
	NumberOfNode     int           `json:"number_of_node" required:"true"`
	AvailabilityZone string        `json:"availability_zone,omitempty"`
	SubnetID         string        `json:"subnet_id" required:"true"`
	UserPwd          string        `json:"user_pwd" required:"true"`
	SecurityGroupID  string        `json:"security_group_id" required:"true"`
	PublicIp         *PublicIpOpts `json:"public_ip,omitempty"`
	NodeType         string        `json:"node_type" required:"true"`
	VpcID            string        `json:"vpc_id" required:"true"`
	UserName         string        `json:"user_name" required:"true"`
	Port             int           `json:"port,omitempty"`
}

type CreateOptsBuilder interface {
	ToClusterCreateMap() (map[string]interface{}, error)
}

func (opts CreateOpts) ToClusterCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "cluster")
}

func Create(c *golangsdk.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToClusterCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	log.Printf("[DEBUG] create url:%q, body=%#v", createURL(c), b)
	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{200}}
	_, r.Err = c.Post(createURL(c), b, &r.Body, reqOpt)
	return
}

func Get(c *golangsdk.ServiceClient, id string) (r GetResult) {
	_, r.Err = c.Get(resourceURL(c, id), &r.Body, nil)
	return
}

func Delete(c *golangsdk.ServiceClient, id string) (r DeleteResult) {
	reqOpt := &golangsdk.RequestOpts{
		OkCodes: []int{202},
		JSONBody: map[string]interface{}{
			"keep_last_manual_snapshot": 0,
		},
	}
	_, r.Err = c.Delete(resourceURL(c, id), reqOpt)
	return
}
