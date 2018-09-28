package cluster

import "github.com/huaweicloud/golangsdk"

const (
	resourcePath = "clusters"
)

func createURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(resourcePath)
}

func resourceURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(resourcePath, id)
}
