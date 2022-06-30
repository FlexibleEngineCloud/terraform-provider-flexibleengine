package flexibleengine

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/kms/v1/keys"
	"github.com/chnsz/golangsdk/openstack/kms/v1/rotation"
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
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
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
			"pending_days": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "7",
			},
			"is_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"rotation_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"rotation_interval": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				RequiredWith: []string{"rotation_enabled"},
				ValidateFunc: validation.IntBetween(30, 365),
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
			"default_key_flag": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"origin": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"rotation_number": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceKmsKeyValidation(d *schema.ResourceData) error {
	_, rotationEnabled := d.GetOk("rotation_enabled")
	_, hasInterval := d.GetOk("rotation_interval")

	if !rotationEnabled && hasInterval {
		return fmt.Errorf("invalid arguments: rotation_interval is only valid when rotation is enabled")
	}
	return nil
}

func resourceKmsKeyV1Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	kmsClient, err := config.KmsKeyV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine kms key client: %s", err)
	}

	createOpts := &keys.CreateOpts{
		KeyAlias:       d.Get("key_alias").(string),
		KeyDescription: d.Get("key_description").(string),
		Realm:          d.Get("realm").(string),
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	v, err := keys.Create(kmsClient, createOpts).ExtractKeyInfo()
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine key: %s", err)
	}

	// Store the key ID now
	d.SetId(v.KeyID)

	// Wait for the key to become enabled.
	log.Printf("[DEBUG] Waiting for key (%s) to become enabled", v.KeyID)
	stateConf := &resource.StateChangeConf{
		Pending:    []string{WaitingForEnableState, DisabledState},
		Target:     []string{EnabledState},
		Refresh:    keyV1StateRefreshFunc(kmsClient, v.KeyID),
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
		key, err := keys.DisableKey(kmsClient, v.KeyID).ExtractKeyInfo()
		if err != nil {
			return fmt.Errorf("Error disabling key: %s", err)
		}

		if key.KeyState != DisabledState {
			return fmt.Errorf("Error disabling key, the key state is: %s", key.KeyState)
		}
	}

	// enable rotation and change interval if necessary
	if _, ok := d.GetOk("rotation_enabled"); ok {
		rotationOpts := &rotation.RotationOpts{
			KeyID: v.KeyID,
		}
		err := rotation.Enable(kmsClient, rotationOpts).ExtractErr()
		if err != nil {
			return fmt.Errorf("failed to enable key rotation: %s", err)
		}

		if i, ok := d.GetOk("rotation_interval"); ok {
			intervalOpts := &rotation.IntervalOpts{
				KeyID:    v.KeyID,
				Interval: i.(int),
			}
			err := rotation.Update(kmsClient, intervalOpts).ExtractErr()
			if err != nil {
				return fmt.Errorf("failed to change key rotation interval: %s", err)
			}
		}
	}

	return resourceKmsKeyV1Read(d, meta)
}

func resourceKmsKeyV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	kmsClient, err := config.KmsKeyV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine kms key client: %s", err)
	}

	v, err := keys.Get(kmsClient, d.Id()).ExtractKeyInfo()
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
	d.Set("is_enabled", v.KeyState == EnabledState)
	d.Set("default_key_flag", v.DefaultKeyFlag)
	d.Set("origin", v.Origin)

	// for import, set pending_days default to 7
	if _, ok := d.GetOk("pending_days"); !ok {
		d.Set("pending_days", "7")
	}

	// Set KMS rotation
	rotationOpts := &rotation.RotationOpts{
		KeyID: v.KeyID,
	}
	r, err := rotation.Get(kmsClient, rotationOpts).Extract()
	if err == nil {
		d.Set("rotation_enabled", r.Enabled)
		d.Set("rotation_interval", r.Interval)
		d.Set("rotation_number", r.NumberOfRotations)
	} else {
		log.Printf("[WARN] Error fetching details about key rotation: %s", err)
	}

	return nil
}

func resourceKmsKeyV1Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	kmsClient, err := config.KmsKeyV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine kms key client: %s", err)
	}

	keyID := d.Id()
	if d.HasChange("key_alias") {
		updateAliasOpts := keys.UpdateAliasOpts{
			KeyID:    keyID,
			KeyAlias: d.Get("key_alias").(string),
		}
		_, err = keys.UpdateAlias(kmsClient, updateAliasOpts).ExtractKeyInfo()
		if err != nil {
			return fmt.Errorf("Error updating FlexibleEngine key alias: %s", err)
		}
	}

	if d.HasChange("key_description") {
		updateDesOpts := keys.UpdateDesOpts{
			KeyID:          keyID,
			KeyDescription: d.Get("key_description").(string),
		}
		_, err = keys.UpdateDes(kmsClient, updateDesOpts).ExtractKeyInfo()
		if err != nil {
			return fmt.Errorf("Error updating FlexibleEngine key description: %s", err)
		}
	}

	if d.HasChange("is_enabled") {
		v, err := keys.Get(kmsClient, keyID).ExtractKeyInfo()
		if err != nil {
			return fmt.Errorf("Error fetching FlexibleEngine key: %s", err)
		}

		if d.Get("is_enabled").(bool) && v.KeyState == DisabledState {
			key, err := keys.EnableKey(kmsClient, keyID).ExtractKeyInfo()
			if err != nil {
				return fmt.Errorf("Error enabling key: %s", err)
			}
			if key.KeyState != EnabledState {
				return fmt.Errorf("Error enabling key, the key state is: %s", key.KeyState)
			}
		}

		if !d.Get("is_enabled").(bool) && v.KeyState == EnabledState {
			key, err := keys.DisableKey(kmsClient, keyID).ExtractKeyInfo()
			if err != nil {
				return fmt.Errorf("Error disabling key: %s", err)
			}
			if key.KeyState != DisabledState {
				return fmt.Errorf("Error disabling key, the key state is: %s", key.KeyState)
			}
		}
	}

	_, rotationEnabled := d.GetOk("rotation_enabled")
	if d.HasChange("rotation_enabled") {
		var rotationErr error
		rotationOpts := &rotation.RotationOpts{
			KeyID: keyID,
		}
		if rotationEnabled {
			rotationErr = rotation.Enable(kmsClient, rotationOpts).ExtractErr()
		} else {
			rotationErr = rotation.Disable(kmsClient, rotationOpts).ExtractErr()
		}

		if rotationErr != nil {
			return fmt.Errorf("failed to update key rotation status: %s", err)
		}
	}

	if rotationEnabled && d.HasChange("rotation_interval") {
		intervalOpts := &rotation.IntervalOpts{
			KeyID:    keyID,
			Interval: d.Get("rotation_interval").(int),
		}
		err := rotation.Update(kmsClient, intervalOpts).ExtractErr()
		if err != nil {
			return fmt.Errorf("failed to change key rotation interval: %s", err)
		}
	}

	return resourceKmsKeyV1Read(d, meta)
}

func resourceKmsKeyV1Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	kmsClient, err := config.KmsKeyV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine kms key client: %s", err)
	}

	v, err := keys.Get(kmsClient, d.Id()).ExtractKeyInfo()
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
		v, err = keys.Delete(kmsClient, deleteOpts).Extract()
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
