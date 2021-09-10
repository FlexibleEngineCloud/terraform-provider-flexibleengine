package replicationconsistencygroups

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

// CreateOpsBuilder is used for creating replication consistency group parameters.
// any struct providing the parameters should implement this interface
type CreateOpsBuilder interface {
	ToReplicationConsistencyGroupCreateMap() (map[string]interface{}, error)
}

// CreateOps is a struct that contains all the parameters.
type CreateOps struct {
	// The name of the replication consistency group.
	// The name can contain a maximum of 255 bytes.
	Name string `json:"name,omitempty"`

	// The description of the replication consistency group.
	// The description can contain a maximum of 255 bytes.
	Description string `json:"description,omitempty"`

	// The IDs of the EVS replication pairs used to
	// create the replication consistency group.
	ReplicationIDs []string `json:"replication_ids" required:"true"`

	// The primary AZ of the replication consistency group.
	// That is the AZ where the production disk belongs.
	PriorityStation string `json:"priority_station" required:"true"`

	// The type of the created replication consistency group.
	// Currently only type hypermetro is supported.
	ReplicationModel string `json:"replication_model" required:"true"`
}

// ToReplicationConsistencyGroupCreateMap is used for type convert
func (ops CreateOps) ToReplicationConsistencyGroupCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(ops, "replication_consistency_group")
}

// Create a replication consistency group with given parameters.
func Create(client *golangsdk.ServiceClient, ops CreateOpsBuilder) (r CreateResult) {
	b, err := ops.ToReplicationConsistencyGroupCreateMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(createURL(client), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{202},
	})

	return
}

// Delete a replication consistency group by id
func Delete(client *golangsdk.ServiceClient, id string) (r DeleteResult) {
	_, r.Err = client.Delete(deleteURL(client, id), nil)
	return
}

// Get a replication consistency group with detailed information by id
func Get(client *golangsdk.ServiceClient, id string) (r GetResult) {
	_, r.Err = client.Get(getURL(client, id), &r.Body, nil)
	return
}

// ListOptsBuilder is an interface by which can be able to
// build the query string of the list function
type ListOptsBuilder interface {
	ToReplicationConsistencyGroupListQuery() (string, error)
}

// ListOpts is a struct that contains all the parameters.
type ListOpts struct {
	Marker          string `q:"marker"`
	Limit           int    `q:"limit"`
	SortKey         string `q:"sort_key"`
	SortDir         string `q:"sort_dir"`
	Offset          int    `q:"offset"`
	ChangesSince    string `q:"changes-since"`
	Name            string `q:"name"`
	Status          string `q:"status"`
	PriorityStation string `q:"priority_station"`
	VolumeID        string `q:"volume_id"`
}

// ToReplicationConsistencyGroupListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToReplicationConsistencyGroupListQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	return q.String(), err
}

// List all the replication consistency groups
func List(client *golangsdk.ServiceClient, ops ListOptsBuilder) pagination.Pager {
	url := listURL(client)
	if ops != nil {
		q, err := ops.ToReplicationConsistencyGroupListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += q
	}

	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return ReplicationConsistencyGroupPage{pagination.SinglePageBase(r)}
	})
}

//UpdateOptsBuilder is an interface which can build the map paramter of update function
type UpdateOptsBuilder interface {
	ToReplicationConsistencyGroupUpdateMap() (map[string]interface{}, error)
}

//UpdateOpts is a struct which represents the parameters of update function
type UpdateOpts struct {
	// The name of the replication consistency group.
	// The name can contain a maximum of 255 bytes.
	Name string `json:"name,omitempty"`

	// The description of the replication consistency group.
	// The description can contain a maximum of 255 bytes.
	Description string `json:"description,omitempty"`

	// The type of the created replication consistency group.
	// Currently only type hypermetro is supported.
	ReplicationModel string `json:"replication_model,omitempty"`

	// The IDs of the EVS replication pairs to be added.
	AddReplicationIDs []string `json:"add_replication_ids,omitempty"`

	// The IDs of the EVS replication pairs to be removeed.
	RemoveReplicationIDs []string `json:"remove_replication_ids,omitempty"`
}

// ToReplicationConsistencyGroupUpdateMap is used for type convert
func (opts UpdateOpts) ToReplicationConsistencyGroupUpdateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "replication_consistency_group")
}

// Update is a method which can be able to update the replication consistency group
// via accessing to the service with Put method and parameters
func Update(client *golangsdk.ServiceClient, id string, opts UpdateOptsBuilder) (r UpdateResult) {
	body, err := opts.ToReplicationConsistencyGroupUpdateMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Put(updateURL(client, id), body, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{202},
	})
	return
}

const (
	// OsFailoverReplicationConsistencyGroup is performing a failover for a replication consistency group.
	OsFailoverReplicationConsistencyGroup = "os-failover-replication-consistency-group"

	// OsSyncReplicationConsistencyGroup is synchronizing a replication consistency group.
	OsSyncReplicationConsistencyGroup = "os-sync-replication-consistency-group"

	// OsReverseReplicationConsistencyGroup is performing a primary/secondary switchover for a replication consistency group.
	OsReverseReplicationConsistencyGroup = "os-reverse-replication-consistency-group"

	// OsStopReplicationConsistencyGroup is pausing a replication consistency group.
	OsStopReplicationConsistencyGroup = "os-stop-replication-consistency-group"

	// OsReprotectReplicationConsistencyGroup is reprotecting a replication consistency group.
	OsReprotectReplicationConsistencyGroup = "os-reprotect-replication-consistency-group"

	// OsExtendReplicationVolume is expanding EVS disks in a replication consistency group.
	OsExtendReplicationVolume = "os-extend-replication-volumes"
)

// FailOver is performing a failover for a replication consistency group.
func FailOver(client *golangsdk.ServiceClient, id string) (r ActionResult) {
	_, r.Err = client.Post(
		actionURL(client, id),
		map[string]interface{}{OsFailoverReplicationConsistencyGroup: nil},
		nil,
		&golangsdk.RequestOpts{
			OkCodes: []int{200},
		})

	return
}

// Sync is synchronizing a replication consistency group.
func Sync(client *golangsdk.ServiceClient, id string) (r ActionResult) {
	_, r.Err = client.Post(
		actionURL(client, id),
		map[string]interface{}{OsSyncReplicationConsistencyGroup: nil},
		nil,
		&golangsdk.RequestOpts{
			OkCodes: []int{200},
		})

	return
}

// Reverse is performing a primary/secondary switchover for a replication consistency group.
func Reverse(client *golangsdk.ServiceClient, id string) (r ActionResult) {
	_, r.Err = client.Post(
		actionURL(client, id),
		map[string]interface{}{OsReverseReplicationConsistencyGroup: nil},
		nil,
		&golangsdk.RequestOpts{
			OkCodes: []int{200},
		})

	return
}

// Stop is pausing a replication consistency group.
func Stop(client *golangsdk.ServiceClient, id string) (r ActionResult) {
	_, r.Err = client.Post(
		actionURL(client, id),
		map[string]interface{}{OsStopReplicationConsistencyGroup: nil},
		nil,
		&golangsdk.RequestOpts{
			OkCodes: []int{200},
		})

	return
}

// Reprotect is reprotecting a replication consistency group.
func Reprotect(client *golangsdk.ServiceClient, id string) (r ActionResult) {
	_, r.Err = client.Post(
		actionURL(client, id),
		map[string]interface{}{OsReprotectReplicationConsistencyGroup: nil},
		nil,
		&golangsdk.RequestOpts{
			OkCodes: []int{200},
		})

	return
}

// ExtendReplicationVolumesOpsBuilder is used for expanding
// EVS disks in a replication consistency group parameters.
// any struct providing the parameters should implement this interface
type ExtendReplicationVolumesOpsBuilder interface {
	ToExtendReplicationVolumesMap() (map[string]interface{}, error)
}

// ReplicationsOps specifies the expansion information
// of one or multiple EVS replication pairs.
type ReplicationsOps struct {
	// The IDs of EVS replication pairs
	ID string `json:"id" required:"true"`

	// The disk capacity after expansion in the EVS replication pair.
	// The unit is GB.
	NewSize int `json:"new_size" required:"true"`
}

// ExtendReplicationVolumesOps is a struct that contains all the parameters.
type ExtendReplicationVolumesOps struct {
	// The expansion information of one or multiple EVS replication pairs.
	Replications []ReplicationsOps `json:"replications" required:"true"`
}

// ToExtendReplicationVolumesMap is used for type convert
func (ops ReplicationsOps) ToExtendReplicationVolumesMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(ops, OsExtendReplicationVolume)
}

// Extend is expanding EVS disks in a replication consistency group.
func Extend(client *golangsdk.ServiceClient, id string, ops ExtendReplicationVolumesOpsBuilder) (r ActionResult) {
	b, err := ops.ToExtendReplicationVolumesMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(actionURL(client, id), b, nil, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})

	return
}
