package instances

import "github.com/huaweicloud/golangsdk"

func createURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("instances")
}

func resourceURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL("instances", id)
}
