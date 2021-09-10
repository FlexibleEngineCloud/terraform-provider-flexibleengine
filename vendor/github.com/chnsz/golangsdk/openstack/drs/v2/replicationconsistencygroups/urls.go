package replicationconsistencygroups

import "github.com/chnsz/golangsdk"

// endpoint/os-vendor-replication-consistency-groups
const resourcePath = "os-vendor-replication-consistency-groups"

// createURL will build the rest query url of creation
func createURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL(resourcePath)
}

// deleteURL will build the url of deletion
func deleteURL(client *golangsdk.ServiceClient, id string) string {
	return client.ServiceURL(resourcePath, id)
}

// getURL will build the get url of get function
func getURL(client *golangsdk.ServiceClient, id string) string {
	return client.ServiceURL(resourcePath, id)
}

// listURL will build the list url of list function
func listURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL(resourcePath)
}

// updateURL will build the url of update
func updateURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(resourcePath, id)
}

// actionURL will build the url of action
func actionURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(resourcePath, id, "action")
}
