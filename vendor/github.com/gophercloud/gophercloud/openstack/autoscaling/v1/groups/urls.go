package groups

import (
	"github.com/gophercloud/gophercloud"
	"log"
)

const resourcePath = "scaling_group"

func createURL(c *gophercloud.ServiceClient) string {
	ur := c.ServiceURL(c.ProjectID, resourcePath)
	log.Printf("[DEBUG] Create URL is: %#v", ur)
	return ur
}

func deleteURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(c.ProjectID, resourcePath, id)
}

func getURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(c.ProjectID, resourcePath, id)
}

func listURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(c.ProjectID, resourcePath)
}

func enableURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(c.ProjectID, resourcePath, id, "action")
}

func updateURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(c.ProjectID, resourcePath, id)
}
