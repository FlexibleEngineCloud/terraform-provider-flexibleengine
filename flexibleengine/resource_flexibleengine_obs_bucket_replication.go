package flexibleengine

import (
	"fmt"
	"log"

	"github.com/chnsz/golangsdk/openstack/obs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceObsBucketReplication() *schema.Resource {
	return &schema.Resource{
		Create: resourceObsBucketReplicationCreate,
		Update: resourceObsBucketReplicationCreate,
		Read:   resourceObsBucketReplicationRead,
		Delete: resourceObsBucketReplicationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"destination_bucket": {
				Type:     schema.TypeString,
				Required: true,
			},
			"agency": {
				Type:     schema.TypeString,
				Required: true,
			},
			"rule": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"prefix": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"storage_class": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"STANDARD", "WARM", "COLD",
							}, false),
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceObsBucketReplicationCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	obsClient, err := config.ObjectStorageClientWithSignature(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine OBS client: %s", err)
	}

	var replicationRules []obs.ReplicationRule
	bucket := d.Get("bucket").(string)
	destBucket := d.Get("destination_bucket").(string)

	rules := d.Get("rule").([]interface{})
	totalRules := len(rules)
	if totalRules == 0 {
		replicationRules = []obs.ReplicationRule{
			{
				Status:            obs.RuleStatusEnabled,
				DestinationBucket: destBucket,
			},
		}
	} else {
		replicationRules = make([]obs.ReplicationRule, totalRules)
		for i, raw := range rules {
			ruleItem := raw.(map[string]interface{})

			replicationRules[i].DestinationBucket = destBucket

			// rule ID
			replicationRules[i].ID = ruleItem["id"].(string)

			// rule Status
			if val, ok := ruleItem["enabled"].(bool); ok && val {
				replicationRules[i].Status = obs.RuleStatusEnabled
			} else {
				replicationRules[i].Status = obs.RuleStatusDisabled
			}

			// Prefix
			prefix := ruleItem["prefix"].(string)
			if prefix == "" && totalRules > 1 {
				return fmt.Errorf("To apply a rule to all objects, delete all rules that take effect by prefixes first")
			}
			replicationRules[i].Prefix = prefix

			if val, ok := ruleItem["storage_class"].(string); ok {
				replicationRules[i].StorageClass = obs.ParseStringToStorageClassType(val)
			}
		}
	}

	opts := &obs.SetBucketReplicationInput{}
	opts.Bucket = bucket
	opts.BucketReplicationConfiguration = obs.BucketReplicationConfiguration{
		Agency:           d.Get("agency").(string),
		ReplicationRules: replicationRules,
	}
	log.Printf("[DEBUG] set cross-region replication of OBS bucket %s: %#v", bucket, opts)

	_, err = obsClient.SetBucketReplication(opts)
	if err != nil {
		return getObsError("Error setting cross-region replication of OBS bucket", bucket, err)
	}

	// Assign the source bucket name as the resource ID
	d.SetId(bucket)

	return resourceObsBucketReplicationRead(d, meta)
}

func resourceObsBucketReplicationRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	obsClient, err := config.ObjectStorageClientWithSignature(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine OBS client: %s", err)
	}

	return setObsBucketReplicationConfiguration(obsClient, d)
}

func setObsBucketReplicationConfiguration(obsClient *obs.ObsClient, d *schema.ResourceData) error {
	bucket := d.Id()
	output, err := obsClient.GetBucketReplication(bucket)
	if err != nil {
		if obsError, ok := err.(obs.ObsError); ok {
			if obsError.Code == "ReplicationConfigurationNotFoundError" {
				d.SetId("")
				return nil
			}
			return fmt.Errorf("Error getting cross-region replication configuration of OBS bucket %s: %s,\n Reason: %s",
				bucket, obsError.Code, obsError.Message)

		}
		return err
	}

	var destBucket string
	rawRules := output.ReplicationRules
	log.Printf("[DEBUG] getting cross-region replication configuration of OBS bucket %s: %#v", bucket, rawRules)

	rules := make([]map[string]interface{}, 0, len(rawRules))
	for _, replicaRule := range rawRules {
		rule := make(map[string]interface{})

		if destBucket == "" {
			destBucket = replicaRule.DestinationBucket
		}

		// Enabled
		if replicaRule.Status == obs.RuleStatusEnabled {
			rule["enabled"] = true
		} else {
			rule["enabled"] = false
		}

		rule["id"] = replicaRule.ID
		rule["prefix"] = replicaRule.Prefix
		if replicaRule.StorageClass != "" {
			rule["storage_class"] = replicaRule.StorageClass
		}

		rules = append(rules, rule)
	}

	log.Printf("[DEBUG] saving cross-region replication configuration of OBS bucket %s: %#v", bucket, rules)
	if err := d.Set("rule", rules); err != nil {
		return fmt.Errorf("Error saving cross-region replication configuration of OBS bucket %s: %s", bucket, err)
	}

	d.Set("agency", output.Agency)
	d.Set("destination_bucket", destBucket)
	d.Set("bucket", bucket)

	return nil
}

func resourceObsBucketReplicationDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	obsClient, err := config.ObjectStorageClientWithSignature(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine OBS client: %s", err)
	}

	bucket := d.Id()
	log.Printf("[DEBUG] delete cross-region replication configuration of OBS bucket %s", bucket)
	_, err = obsClient.DeleteBucketReplication(bucket)
	if err != nil {
		return getObsError("Error deleting cross-region replication configuration of OBS bucket", bucket, err)
	}

	return nil
}
