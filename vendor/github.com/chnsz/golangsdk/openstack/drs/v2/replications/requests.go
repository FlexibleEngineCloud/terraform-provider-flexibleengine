package replications

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

// CreateOpsBuilder is used for creating replication parameters.
// any struct providing the parameters should implement this interface
type CreateOpsBuilder interface {
	ToReplicationCreateMap() (map[string]interface{}, error)
}

// CreateOps is a struct that contains all the parameters.
type CreateOps struct {
	// The name of the EVS replication pair.
	// The name can contain a maximum of 255 bytes.
	Name string `json:"name,omitempty"`

	// The description of the EVS replication pair.
	// The description can contain a maximum of 255 bytes.
	Description string `json:"description,omitempty"`

	// The IDs of the EVS disks used to create the EVS replication pair.
	VolumeIDs []string `json:"volume_ids" required:"true"`

	// The primary AZ of the EVS replication pair.
	// That is the AZ where the production disk belongs.
	PriorityStation string `json:"priority_station" required:"true"`

	// The type of the EVS replication pair.
	// Currently only type hypermetro is supported.
	ReplicationModel string `json:"replication_model" required:"true"`
}

// ToReplicationCreateMap is used for type convert
func (ops CreateOps) ToReplicationCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(ops, "replication")
}

// Create a replication with given parameters.
func Create(client *golangsdk.ServiceClient, ops CreateOpsBuilder) (r CreateResult) {
	b, err := ops.ToReplicationCreateMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(createURL(client), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{202},
	})

	return
}

// Delete a replication by id
func Delete(client *golangsdk.ServiceClient, id string) (r DeleteResult) {
	_, r.Err = client.Delete(deleteURL(client, id), nil)
	return
}

// Get a replication with detailed information by id
func Get(client *golangsdk.ServiceClient, id string) (r GetResult) {
	_, r.Err = client.Get(getURL(client, id), &r.Body, nil)
	return
}

// ListOptsBuilder is an interface by which can be able to
// build the query string of the list function
type ListOptsBuilder interface {
	ToReplicationListQuery() (string, error)
}

// ListOpts is a struct that contains all the parameters.
type ListOpts struct {
	Marker                        string `q:"marker"`
	Limit                         int    `q:"limit"`
	SortKey                       string `q:"sort_key"`
	SortDir                       string `q:"sort_dir"`
	Offset                        int    `q:"offset"`
	ChangesSince                  string `q:"changes-since"`
	Name                          string `q:"name"`
	Status                        string `q:"status"`
	ReplicationConsistencyGroupID string `q:"replication_consistency_group_id"`
	VolumeIDs                     string `q:"volume_ids"`
	VolumeID                      string `q:"volume_id"`
	PriorityStation               string `q:"priority_station"`
}

// ToReplicationListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToReplicationListQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	return q.String(), err
}

// List all the replications
func List(client *golangsdk.ServiceClient, ops ListOptsBuilder) pagination.Pager {
	url := listURL(client)
	if ops != nil {
		q, err := ops.ToReplicationListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += q
	}

	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return ReplicationPage{pagination.SinglePageBase(r)}
	})
}
