package keypairs

import (
	"github.com/chnsz/golangsdk"
)

const resourcePath = "os-keypairs"

func listURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(resourcePath)
}
