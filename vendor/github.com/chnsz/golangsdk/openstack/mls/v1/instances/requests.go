package instances

import (
	"github.com/chnsz/golangsdk"
)

var RequestOpts golangsdk.RequestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

//CreateOptsBuilder is used for creating instance parameters.
//any struct providing the parameters should implement this interface
type CreateOptsBuilder interface {
	ToInstanceCreateMap() (map[string]interface{}, error)
}

//CreateOpts is a struct that contains all the parameters.
type CreateOpts struct {
	Name string `json:"name" required:"true"`

	Version string `json:"version" required:"true"`

	Network NetworkOpts `json:"network" required:"true"`

	Agency string `json:"agency,omitempty"`

	FlavorRef string `json:"flavorRef" required:"true"`

	MrsCluster MrsClusterOpts `json:"mrsCluster" required:"true"`
}

type NetworkOpts struct {
	VpcId           string       `json:"vpcId" required:"true"`
	SubnetId        string       `json:"subnetId" required:"true"`
	SecurityGroupId string       `json:"securityGroupId,omitempty"`
	AvailableZone   string       `json:"availableZone" required:"true"`
	PublicIP        PublicIPOpts `json:"publicIP" required:"true"`
}

type PublicIPOpts struct {
	BindType string `json:"bindType" required:"true"`
}

type MrsClusterOpts struct {
	Id           string `json:"id" required:"true"`
	UserName     string `json:"userName,omitempty"`
	UserPassword string `json:"userPassword,omitempty"`
}

func (opts CreateOpts) ToInstanceCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "instance")
}

//Create an instance with given parameters.
func Create(client *golangsdk.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToInstanceCreateMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(createURL(client), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes:     []int{202},
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	return
}

//delete an instance via id
func Delete(client *golangsdk.ServiceClient, id string) (r DeleteResult) {
	b := make(map[string]interface{})
	_, r.Err = client.DeleteWithBody(resourceURL(client, id), b, &golangsdk.RequestOpts{
		OkCodes: []int{202},
	})
	return
}

//get an instance with detailed information by id
func Get(client *golangsdk.ServiceClient, id string) (r GetResult) {
	_, r.Err = client.Get(resourceURL(client, id), &r.Body, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	return
}
