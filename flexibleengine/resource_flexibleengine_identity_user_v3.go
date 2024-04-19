package flexibleengine

import (
	"context"
	"log"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	iam_users "github.com/chnsz/golangsdk/openstack/identity/v3.0/users"
	"github.com/chnsz/golangsdk/openstack/identity/v3/users"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func resourceIdentityUserV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentityUserV3Create,
		ReadContext:   resourceIdentityUserV3Read,
		UpdateContext: resourceIdentityUserV3Update,
		DeleteContext: resourceIdentityUserV3Delete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"access_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"default", "programmatic", "console"}, false),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"email": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"phone": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"country_code"},
				ValidateFunc: validation.StringMatch(regexp.MustCompile("^[0-9]{0,32}$"),
					"the phone number must have a maximum of 32 digits"),
			},
			"country_code": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"phone"},
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"password_strength": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_login": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceIdentityUserV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	iamClient, err := config.IAMV3Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating FlexibleEngine iam client: %s", err)
	}

	if config.DomainID == "" {
		return diag.Errorf("the domain_id must be specified in the provider configuration")
	}

	enabled := d.Get("enabled").(bool)
	createOpts := iam_users.CreateOpts{
		Name:        d.Get("name").(string),
		AccessMode:  d.Get("access_mode").(string),
		Description: d.Get("description").(string),
		Email:       d.Get("email").(string),
		Phone:       d.Get("phone").(string),
		AreaCode:    d.Get("country_code").(string),
		Enabled:     &enabled,
		DomainID:    config.DomainID,
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	// Add password here so it wouldn't go in the above log entry
	createOpts.Password = d.Get("password").(string)

	user, err := iam_users.Create(iamClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("Error creating FlexibleEngine iam user: %s", err)
	}

	d.SetId(user.ID)

	return resourceIdentityUserV3Read(ctx, d, meta)
}

func resourceIdentityUserV3Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	iamClient, err := config.IAMV3Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating FlexibleEngine iam client: %s", err)
	}

	user, err := iam_users.Get(iamClient, d.Id()).Extract()
	if err != nil {
		return CheckDeletedDiag(d, err, "IAM user")
	}

	log.Printf("[DEBUG] Retrieved FlexibleEngine user: %#v", user)

	d.Set("enabled", user.Enabled)
	d.Set("name", user.Name)
	d.Set("access_mode", user.AccessMode)
	d.Set("description", user.Description)
	d.Set("email", user.Email)
	d.Set("country_code", user.AreaCode)
	d.Set("password_strength", user.PasswordStrength)
	d.Set("create_time", user.CreateAt)
	d.Set("last_login", user.LastLogin)

	phone := strings.Split(user.Phone, "-")
	if len(phone) > 1 {
		d.Set("phone", phone[1])
	} else {
		d.Set("phone", user.Phone)
	}

	return nil
}

func resourceIdentityUserV3Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	iamClient, err := config.IAMV3Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating FlexibleEngine iam client: %s", err)
	}

	var updateOpts iam_users.UpdateOpts

	if d.HasChange("name") {
		updateOpts.Name = d.Get("name").(string)
	}

	if d.HasChange("access_mode") {
		updateOpts.AccessMode = d.Get("access_mode").(string)
	}

	if d.HasChange("description") {
		updateOpts.Description = utils.String(d.Get("description").(string))
	}

	if d.HasChange("email") {
		updateOpts.Email = d.Get("email").(string)
	}

	if d.HasChanges("country_code", "phone") {
		updateOpts.AreaCode = d.Get("country_code").(string)
		updateOpts.Phone = d.Get("phone").(string)
	}

	if d.HasChange("enabled") {
		enabled := d.Get("enabled").(bool)
		updateOpts.Enabled = &enabled
	}

	log.Printf("[DEBUG] Update Options: %#v", updateOpts)

	// Add password here so it wouldn't go in the above log entry
	if d.HasChange("password") {
		updateOpts.Password = d.Get("password").(string)
	}

	_, err = iam_users.Update(iamClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return diag.Errorf("Error updating FlexibleEngine user: %s", err)
	}

	return resourceIdentityUserV3Read(ctx, d, meta)
}

func resourceIdentityUserV3Delete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	identityClient, err := config.IdentityV3Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating FlexibleEngine identity client: %s", err)
	}

	err = users.Delete(identityClient, d.Id()).ExtractErr()
	if err != nil {
		return diag.Errorf("Error deleting FlexibleEngine user: %s", err)
	}

	return nil
}
