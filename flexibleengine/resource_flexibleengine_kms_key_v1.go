package flexibleengine

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/kms/v1/keys"
)

const (
	WaitingForEnableState = "1"
	EnabledState          = "2"
	DisabledState         = "3"
	PendingDeletionState  = "4"
)

func resourceKmsKeyV1() *schema.Resource {
	return &schema.Resource{
		Create: resourceKmsKeyV1Create,
		Read:   resourceKmsKeyV1Read,
		Update: resourceKmsKeyV1Update,
		Delete: resourceKmsKeyV1Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"key_alias": {
				Type:     schema.TypeString,
				Required: true,
			},
			"key_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"realm": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"domain_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creation_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"scheduled_deletion_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"default_key_flag": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"expiration_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"origin": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"pending_days": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "7",
			},
		},
	}
}

func resourceKmsKeyV1Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	kmsKeyV1Client, err := config.kmsKeyV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine kms key client: %s", err)
	}

	createOpts := &keys.CreateOpts{
		KeyAlias:       d.Get("key_alias").(string),
		KeyDescription: d.Get("key_description").(string),
		Realm:          d.Get("realm").(string),
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	v, err := keys.Create(kmsKeyV1Client, createOpts).ExtractKeyInfo()
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine key: %s", err)
	}
	log.Printf("[INFO] Key ID: %s", v.KeyID)

	// Wait for the key to become enabled.
	log.Printf("[DEBUG] Waiting for key (%s) to become enabled", v.KeyID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{WaitingForEnableState, DisabledState},
		Target:     []string{EnabledState},
		Refresh:    keyV1StateRefreshFunc(kmsKeyV1Client, v.KeyID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for key (%s) to become ready: %s",
			v.KeyID, err)
	}

	if !d.Get("is_enabled").(bool) {
		key, err := keys.DisableKey(kmsKeyV1Client, v.KeyID).ExtractKeyInfo()
		if err != nil {
			return fmt.Errorf("Error disabling key: %s.", err)
		}

		if key.KeyState != DisabledState {
			return fmt.Errorf("Error disabling key, the key state is: %s", key.KeyState)
		}
	}

	// Store the key ID now
	d.SetId(v.KeyID)

	return resourceKmsKeyV1Read(d, meta)
}

func resourceKmsKeyV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	kmsKeyV1Client, err := config.kmsKeyV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine kms key client: %s", err)
	}

	v, err := keys.Get(kmsKeyV1Client, d.Id()).ExtractKeyInfo()
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Kms key %s: %+v", d.Id(), v)
	if v.KeyState == PendingDeletionState {
		log.Printf("[WARN] Removing KMS key %s because it's already gone", d.Id())
		d.SetId("")
		return nil
	}

	d.SetId(v.KeyID)
	d.Set("domain_id", v.DomainID)
	d.Set("key_alias", v.KeyAlias)
	d.Set("realm", v.Realm)
	d.Set("key_description", v.KeyDescription)
	d.Set("creation_date", v.CreationDate)
	d.Set("scheduled_deletion_date", v.ScheduledDeletionDate)
	d.Set("is_enabled", v.KeyState == EnabledState)
	d.Set("default_key_flag", v.DefaultKeyFlag)
	d.Set("expiration_time", v.ExpirationTime)
	d.Set("origin", v.Origin)

	return nil
}

func resourceKmsKeyV1Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	kmsKeyV1Client, err := config.kmsKeyV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine kms key client: %s", err)
	}

	if d.HasChange("key_alias") {
		updateAliasOpts := keys.UpdateAliasOpts{
			KeyID:    d.Id(),
			KeyAlias: d.Get("key_alias").(string),
		}
		_, err = keys.UpdateAlias(kmsKeyV1Client, updateAliasOpts).ExtractKeyInfo()
		if err != nil {
			return fmt.Errorf("Error updating FlexibleEngine key: %s", err)
		}
	}

	if d.HasChange("key_description") {
		updateDesOpts := keys.UpdateDesOpts{
			KeyID:          d.Id(),
			KeyDescription: d.Get("key_description").(string),
		}
		_, err = keys.UpdateDes(kmsKeyV1Client, updateDesOpts).ExtractKeyInfo()
		if err != nil {
			return fmt.Errorf("Error updating FlexibleEngine key: %s", err)
		}
	}

	if d.HasChange("is_enabled") {
		v, err := keys.Get(kmsKeyV1Client, d.Id()).ExtractKeyInfo()
		if err != nil {
			return fmt.Errorf("DescribeKey got an error: %s.", err)
		}

		if d.Get("is_enabled").(bool) && v.KeyState == DisabledState {
			key, err := keys.EnableKey(kmsKeyV1Client, d.Id()).ExtractKeyInfo()
			if err != nil {
				return fmt.Errorf("Error enabling key: %s.", err)
			}
			if key.KeyState != EnabledState {
				return fmt.Errorf("Error enabling key, the key state is: %s", key.KeyState)
			}
		}

		if !d.Get("is_enabled").(bool) && v.KeyState == EnabledState {
			key, err := keys.DisableKey(kmsKeyV1Client, d.Id()).ExtractKeyInfo()
			if err != nil {
				return fmt.Errorf("Error disabling key: %s.", err)
			}
			if key.KeyState != DisabledState {
				return fmt.Errorf("Error disabling key, the key state is: %s", key.KeyState)
			}
		}
	}

	return resourceKmsKeyV1Read(d, meta)
}

func resourceKmsKeyV1Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	kmsKeyV1Client, err := config.kmsKeyV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine kms key client: %s", err)
	}

	v, err := keys.Get(kmsKeyV1Client, d.Id()).ExtractKeyInfo()
	if err != nil {
		return CheckDeleted(d, err, "key")
	}

	deleteOpts := &keys.DeleteOpts{
		KeyID: d.Id(),
	}
	if v, ok := d.GetOk("pending_days"); ok {
		deleteOpts.PendingDays = v.(string)
	}

	// It's possible that this key was used as a boot device and is currently
	// in a pending deletion state from when the instance was terminated.
	// If this is true, just move on. It'll eventually delete.
	if v.KeyState != PendingDeletionState {
		v, err = keys.Delete(kmsKeyV1Client, deleteOpts).Extract()
		if err != nil {
			return err
		}

		if v.KeyState != PendingDeletionState {
			return fmt.Errorf("failed to delete key")
		}
	}

	log.Printf("[DEBUG] KMS Key %s deactivated.", d.Id())
	d.SetId("")
	return nil
}

func keyV1StateRefreshFunc(client *golangsdk.ServiceClient, keyID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		v, err := keys.Get(client, keyID).ExtractKeyInfo()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return v, PendingDeletionState, nil
			}
			return nil, "", err
		}

		return v, v.KeyState, nil
	}
}
