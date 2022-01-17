package topics

import (
	"github.com/chnsz/golangsdk"
)

// CreateOpsBuilder is an interface which is used for creating a kafka topic
type CreateOpsBuilder interface {
	ToTopicCreateMap() (map[string]interface{}, error)
}

// CreateOps is a struct that contains all the parameters of create function
type CreateOps struct {
	// the name/ID of a topic
	Name string `json:"id" required:"true"`
	// topic partitions, value range: 1-50, Default value:3
	Partition int `json:"partition,omitempty"`
	// topic replications, value range: 1-3, Default value:3
	Replication int `json:"replication,omitempty"`
	// aging time in hours, value range: 1-168, , Default value:72
	RetentionTime int `json:"retention_time,omitempty"`

	SyncMessageFlush bool `json:"sync_message_flush,omitempty"`
	SyncReplication  bool `json:"sync_replication,omitempty"`
}

// ToTopicCreateMap is used for type convert
func (ops CreateOps) ToTopicCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(ops, "")
}

// Create a kafka topic with given parameters
func Create(client *golangsdk.ServiceClient, instanceID string, ops CreateOpsBuilder) (r CreateResult) {
	b, err := ops.ToTopicCreateMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(rootURL(client, instanceID), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})

	return
}

// List all topics belong to the instance id
func List(client *golangsdk.ServiceClient, instanceID string) (r ListResult) {
	_, r.Err = client.Get(rootURL(client, instanceID), &r.Body, nil)
	return
}

// Delete given topics belong to the instance id
func Delete(client *golangsdk.ServiceClient, instanceID string, topics []string) (r DeleteResult) {
	var delOpts = struct {
		Topics []string `json:"topics" required:"true"`
	}{Topics: topics}

	b, err := golangsdk.BuildRequestBody(delOpts, "")
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(deleteURL(client, instanceID), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})

	return
}
