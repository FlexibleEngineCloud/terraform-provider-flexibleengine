package flexibleengine

import (
	"bytes"
	"encoding/json"
	"log"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/jen20/awspolicyequivalence"
)

func suppressEquivalentAwsPolicyDiffs(k, old, new string, d *schema.ResourceData) bool {
	equivalent, err := awspolicy.PoliciesAreEquivalent(old, new)
	if err != nil {
		return false
	}

	return equivalent
}

// Suppress all changes?
func suppressDiffAll(k, old, new string, d *schema.ResourceData) bool {
	return true
}

// Suppress equivalent device name changes, only compare string after first two characters
func suppressDiffDevice(k, old, new string, d *schema.ResourceData) bool {
	// If too short (shouldn't happen, but to be safe), suppress diff
	if len(old) < 2 || len(new) < 2 {
		return true
	}
	return (old[2:] == new[2:])
}

// Suppress changes if we get a computed min_disk_gb if value is unspecified (default 0)
func suppressMinDisk(k, old, new string, d *schema.ResourceData) bool {
	return new == "0" || old == new
}

// Suppress changes if we don't specify an external gateway, but one is specified for us
func suppressExternalGateway(k, old, new string, d *schema.ResourceData) bool {
	return new == "" || old == new
}

// Suppress changes if we get a fixed ip when not expecting one, if we have a floating ip (generates fixed ip).
func suppressComputedFixedWhenFloatingIp(k, old, new string, d *schema.ResourceData) bool {
	if v, ok := d.GetOk("floating_ip"); ok && v != "" {
		return new == "" || old == new
	}
	return false
}

// Suppresses minor version changes to the db_instance engine_version attribute
func suppressAwsDbEngineVersionDiffs(k, old, new string, d *schema.ResourceData) bool {
	// First check if the old/new values are nil.
	// If both are nil, we have no state to compare the values with, so register a diff.
	// This populates the attribute field during a plan/apply with fresh state, allowing
	// the attribute to still be used in future resources.
	// See https://github.com/hashicorp/terraform/issues/11881
	if old == "" && new == "" {
		return false
	}

	if v, ok := d.GetOk("auto_minor_version_upgrade"); ok {
		if v.(bool) {
			// If we're set to auto upgrade minor versions
			// ignore a minor version diff between versions
			if strings.HasPrefix(old, new) {
				log.Printf("[DEBUG] Ignoring minor version diff")
				return true
			}
		}
	}

	// Throw a diff by default
	return false
}

func suppressEquivalentJsonDiffs(k, old, new string, d *schema.ResourceData) bool {
	ob := bytes.NewBufferString("")
	if err := json.Compact(ob, []byte(old)); err != nil {
		return false
	}

	nb := bytes.NewBufferString("")
	if err := json.Compact(nb, []byte(new)); err != nil {
		return false
	}

	return jsonBytesEqual(ob.Bytes(), nb.Bytes())
}

func suppressOpenIdURL(k, old, new string, d *schema.ResourceData) bool {
	oldUrl, err := url.Parse(old)
	if err != nil {
		return false
	}

	newUrl, err := url.Parse(new)
	if err != nil {
		return false
	}

	oldUrl.Scheme = "https"

	return oldUrl.String() == newUrl.String()
}

func suppressAutoscalingGroupAvailabilityZoneDiffs(k, old, new string, d *schema.ResourceData) bool {
	// If VPC zone identifiers are provided then there is no need to explicitly
	// specify availability zones.
	if _, ok := d.GetOk("vpc_zone_identifier"); ok {
		return true
	}

	return false
}

func suppressRdsNameDiffs(k, old, new string, d *schema.ResourceData) bool {
	if strings.HasPrefix(old, new) && strings.HasSuffix(old, "_node0") {
		return true
	}
	return false
}
