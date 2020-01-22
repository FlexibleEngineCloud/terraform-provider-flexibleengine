package drill

import "github.com/huaweicloud/golangsdk"

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("disaster-recovery-drills")
}

func resourceURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL("disaster-recovery-drills", id)
}
