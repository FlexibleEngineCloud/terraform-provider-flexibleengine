package topics

import (
	"github.com/chnsz/golangsdk"
)

const (
	resourcePath = "instances"
	topicPath    = "topics"
)

// rootURL will build the url of create, update and list
func rootURL(client *golangsdk.ServiceClient, instanceID string) string {
	return client.ServiceURL(resourcePath, instanceID, topicPath)
}

// deleteURL will build the url of delete
func deleteURL(client *golangsdk.ServiceClient, instanceID string) string {
	return client.ServiceURL(resourcePath, instanceID, topicPath, "delete")
}
