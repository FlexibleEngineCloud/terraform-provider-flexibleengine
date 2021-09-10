package ccattackprotection_rules

import "github.com/chnsz/golangsdk"

func rootURL(c *golangsdk.ServiceClient, policyID string) string {
	return c.ServiceURL("policy", policyID, "cc")
}

func resourceURL(c *golangsdk.ServiceClient, policyID, id string) string {
	return c.ServiceURL("policy", policyID, "cc", id)
}
