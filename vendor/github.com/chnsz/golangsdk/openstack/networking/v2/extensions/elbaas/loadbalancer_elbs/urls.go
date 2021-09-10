package loadbalancer_elbs

import "github.com/chnsz/golangsdk"

const (
	rootPath     = "elbaas"
	resourcePath = "loadbalancers"
)

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(rootPath, resourcePath)
}

func resourceURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(rootPath, resourcePath, id)
}

func deleteURL(c *golangsdk.ServiceClient, id string, keepEIP bool) string {
	if keepEIP {
		return c.ServiceURL(rootPath, resourcePath, id, "keep-eip")
	}
	return c.ServiceURL(rootPath, resourcePath, id)
}
