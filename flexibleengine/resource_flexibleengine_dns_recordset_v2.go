package flexibleengine

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/common/tags"
	"github.com/chnsz/golangsdk/openstack/dns/v2/recordsets"
	"github.com/chnsz/golangsdk/openstack/dns/v2/zones"
)

func resourceDNSRecordSetV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceDNSRecordSetV2Create,
		Read:   resourceDNSRecordSetV2Read,
		Update: resourceDNSRecordSetV2Update,
		Delete: resourceDNSRecordSetV2Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 255),
			},
			"records": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"ttl": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      300,
				ValidateFunc: validation.IntBetween(1, 2147483647),
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"A", "AAAA", "MX", "CNAME", "TXT", "NS", "SRV", "PTR", "CAA",
				}, false),
			},
			"value_specs": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"tags": tagsSchema(),
		},
	}
}

func resourceDNSRecordSetV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	dnsClient, err := config.DnsV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine DNS client: %s", err)
	}

	zoneID := d.Get("zone_id").(string)
	zoneType, err := getZoneTypebyID(dnsClient, zoneID)
	if err != nil {
		return fmt.Errorf("Error retrieving DNS zone %s: %s", zoneID, err)
	}

	recordsraw := d.Get("records").(*schema.Set).List()
	records := make([]string, len(recordsraw))
	for i, recordraw := range recordsraw {
		records[i] = recordraw.(string)
	}

	createOpts := RecordSetCreateOpts{
		recordsets.CreateOpts{
			Name:        d.Get("name").(string),
			Description: d.Get("description").(string),
			Records:     records,
			TTL:         d.Get("ttl").(int),
			Type:        d.Get("type").(string),
		},
		MapValueSpecs(d),
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	n, err := recordsets.Create(dnsClient, zoneID, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine DNS record set: %s", err)
	}

	log.Printf("[DEBUG] Created FlexibleEngine DNS record set %s: %#v", n.ID, n)
	id := fmt.Sprintf("%s/%s", zoneID, n.ID)
	d.SetId(id)

	log.Printf("[DEBUG] Waiting for DNS record set (%s) to become available", n.ID)
	stateConf := &resource.StateChangeConf{
		Target:     []string{"ACTIVE"},
		Pending:    []string{"PENDING"},
		Refresh:    waitForDNSRecordSet(dnsClient, zoneID, n.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()

	if err != nil {
		return fmt.Errorf(
			"Error waiting for record set (%s) to become ACTIVE for creation: %s",
			n.ID, err)
	}

	// set tags
	tagRaw := d.Get("tags").(map[string]interface{})
	if len(tagRaw) > 0 {
		resourceType, err := getDNSRecordSetTagType(zoneType)
		if err != nil {
			return fmt.Errorf("Error getting resource type of DNS record set %s: %s", n.ID, err)
		}

		taglist := expandResourceTags(tagRaw)
		if tagErr := tags.Create(dnsClient, resourceType, n.ID, taglist).ExtractErr(); tagErr != nil {
			return fmt.Errorf("Error setting tags of DNS record set %s: %s", n.ID, tagErr)
		}
	}

	return resourceDNSRecordSetV2Read(d, meta)
}

func resourceDNSRecordSetV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	dnsClient, err := config.DnsV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine DNS client: %s", err)
	}

	// Obtain relevant info from parsing the ID
	zoneID, recordsetID, err := parseDNSV2RecordSetId(d.Id())
	if err != nil {
		return err
	}
	// get zone type: public or priavte
	zoneType, err := getZoneTypebyID(dnsClient, zoneID)
	if err != nil {
		return fmt.Errorf("Error retrieving DNS zone %s: %s", zoneID, err)
	}

	time.Sleep(2 * time.Second)
	n, err := recordsets.Get(dnsClient, zoneID, recordsetID).Extract()
	if err != nil {
		return CheckDeleted(d, err, "record_set")
	}

	log.Printf("[DEBUG] Retrieved  record set %s: %#v", recordsetID, n)

	d.Set("name", n.Name)
	d.Set("description", n.Description)
	d.Set("ttl", n.TTL)
	d.Set("type", n.Type)
	if err := d.Set("records", n.Records); err != nil {
		return fmt.Errorf("Error saving records to state for FlexibleEngine DNS record set (%s): %s", d.Id(), err)
	}
	d.Set("region", GetRegion(d, config))
	d.Set("zone_id", zoneID)

	// save tags
	if resourceType, err := getDNSRecordSetTagType(zoneType); err == nil {
		resourceTags, err := tags.Get(dnsClient, resourceType, recordsetID).Extract()
		if err == nil {
			tagmap := tagsToMap(resourceTags.Tags)
			d.Set("tags", tagmap)
		} else {
			log.Printf("[WARN] Error fetching FlexibleEngine DNS record set tags: %s", err)
		}
	}

	return nil
}

func resourceDNSRecordSetV2Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	dnsClient, err := config.DnsV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine DNS client: %s", err)
	}

	// Obtain relevant info from parsing the ID
	zoneID, recordsetID, err := parseDNSV2RecordSetId(d.Id())
	if err != nil {
		return err
	}
	// get zone type: public or priavte
	zoneType, err := getZoneTypebyID(dnsClient, zoneID)
	if err != nil {
		return fmt.Errorf("Error retrieving DNS zone %s: %s", zoneID, err)
	}

	if d.HasChanges("description", "ttl", "records") {
		var updateOpts recordsets.UpdateOpts

		// fix #703
		// API issue: `UpdateOpts.Records` field should not be empty
		// "code":"DNS.0308", "message":"Attribute 'records' is invalid, records is null or empty."
		// if you want to change it, please verify again.
		recordsraw := d.Get("records").(*schema.Set).List()
		records := make([]string, len(recordsraw))
		for i, recordraw := range recordsraw {
			records[i] = recordraw.(string)
		}
		updateOpts.Records = records

		if d.HasChange("ttl") {
			updateOpts.TTL = d.Get("ttl").(int)
		}

		if d.HasChange("description") {
			updateOpts.Description = d.Get("description").(string)
		}

		log.Printf("[DEBUG] Updating record set %s with options: %#v", recordsetID, updateOpts)
		_, err = recordsets.Update(dnsClient, zoneID, recordsetID, updateOpts).Extract()
		if err != nil {
			return fmt.Errorf("Error updating FlexibleEngine DNS record set: %s", err)
		}

		log.Printf("[DEBUG] Waiting for DNS record set (%s) to update", recordsetID)
		stateConf := &resource.StateChangeConf{
			Target:     []string{"ACTIVE"},
			Pending:    []string{"PENDING"},
			Refresh:    waitForDNSRecordSet(dnsClient, zoneID, recordsetID),
			Timeout:    d.Timeout(schema.TimeoutUpdate),
			Delay:      5 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, err = stateConf.WaitForState()
		if err != nil {
			return fmt.Errorf(
				"Error waiting for record set (%s) to become ACTIVE for updation: %s",
				recordsetID, err)
		}
	}

	// update tags
	resourceType, err := getDNSRecordSetTagType(zoneType)
	if err != nil {
		return fmt.Errorf("Error getting resource type of DNS record set %s: %s", d.Id(), err)
	}

	tagErr := UpdateResourceTags(dnsClient, d, resourceType, recordsetID)
	if tagErr != nil {
		return fmt.Errorf("Error updating tags of DNS record set %s: %s", d.Id(), tagErr)
	}

	return resourceDNSRecordSetV2Read(d, meta)
}

func resourceDNSRecordSetV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	dnsClient, err := config.DnsV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine DNS client: %s", err)
	}

	// Obtain relevant info from parsing the ID
	zoneID, recordsetID, err := parseDNSV2RecordSetId(d.Id())
	if err != nil {
		return err
	}

	err = recordsets.Delete(dnsClient, zoneID, recordsetID).ExtractErr()
	if err != nil {
		return fmt.Errorf("Error deleting FlexibleEngine DNS record set: %s", err)
	}

	log.Printf("[DEBUG] Waiting for DNS record set (%s) to be deleted", recordsetID)
	stateConf := &resource.StateChangeConf{
		Target:     []string{"DELETED"},
		Pending:    []string{"ACTIVE", "PENDING", "ERROR"},
		Refresh:    waitForDNSRecordSet(dnsClient, zoneID, recordsetID),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for record set (%s) to become DELETED for deletion: %s",
			recordsetID, err)
	}

	d.SetId("")
	return nil
}

func parseStatus(rawStatus string) string {
	splits := strings.Split(rawStatus, "_")
	// rawStatus maybe one of PENDING_CREATE, PENDING_UPDATE, PENDING_DELETE, ACTIVE, or ERROR
	return splits[0]
}

func waitForDNSRecordSet(dnsClient *golangsdk.ServiceClient, zoneID, recordsetId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		recordset, err := recordsets.Get(dnsClient, zoneID, recordsetId).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return recordset, "DELETED", nil
			}

			return nil, "", err
		}

		log.Printf("[DEBUG] FlexibleEngine DNS record set (%s) current status: %s", recordset.ID, recordset.Status)
		return recordset, parseStatus(recordset.Status), nil
	}
}

func parseDNSV2RecordSetId(id string) (string, string, error) {
	idParts := strings.Split(id, "/")
	if len(idParts) != 2 {
		return "", "", fmt.Errorf("Unable to determine DNS record set ID from raw ID: %s", id)
	}

	zoneID := idParts[0]
	recordsetID := idParts[1]

	return zoneID, recordsetID, nil
}

func getZoneTypebyID(dnsClient *golangsdk.ServiceClient, zoneID string) (string, error) {
	n, err := zones.Get(dnsClient, zoneID).Extract()
	if err != nil {
		return "", err
	}

	return n.ZoneType, nil
}
