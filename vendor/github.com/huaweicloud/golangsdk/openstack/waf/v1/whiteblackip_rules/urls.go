package whiteblackip_rules

import "github.com/huaweicloud/golangsdk"

func rootURL(c *golangsdk.ServiceClient, policyID string) string {
	return c.ServiceURL("policy", policyID, "whiteblackip")
}

func resourceURL(c *golangsdk.ServiceClient, policyID, id string) string {
	return c.ServiceURL("policy", policyID, "whiteblackip", id)
}
