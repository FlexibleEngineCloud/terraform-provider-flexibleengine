package flexibleengine

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/identity/federatedauth/mappings"
	"github.com/chnsz/golangsdk/openstack/identity/federatedauth/metadatas"
	"github.com/chnsz/golangsdk/openstack/identity/federatedauth/oidcconfig"
	"github.com/chnsz/golangsdk/openstack/identity/federatedauth/protocols"
	"github.com/chnsz/golangsdk/openstack/identity/federatedauth/providers"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const (
	protocolSAML = "saml"
	protocolOIDC = "oidc"

	scopeSpilt = " "
)

func resourceIdentityProvider() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentityProviderCreate,
		ReadContext:   resourceIdentityProviderRead,
		UpdateContext: resourceIdentityProviderUpdate,
		DeleteContext: resourceIdentityProviderDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[\w-]{1,64}$`),
					"The maximum length is 64 characters. "+
						"Only letters, digits, underscores (_), and hyphens (-) are allowed"),
			},
			"protocol": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{protocolSAML, protocolOIDC}, false),
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"metadata": {
				Type:      schema.TypeString,
				Optional:  true,
				StateFunc: utils.HashAndHexEncode,
			},
			"openid_connect_config": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"access_type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"program", "program_console"}, false),
						},
						"provider_url": {
							Type:     schema.TypeString,
							Required: true,
						},
						"client_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"signing_key": {
							Type:     schema.TypeString,
							Required: true,
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								equal, _ := utils.CompareJsonTemplateAreEquivalent(old, new)
								return equal
							},
						},
						"authorization_endpoint": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"scopes": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 10,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"response_type": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "id_token",
						},
						"response_mode": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "form_post",
							ValidateFunc: validation.StringInSlice([]string{"fragment", "form_post"}, false),
						},
					},
				},
			},
			"conversion_rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"local": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"username": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"group": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"remote": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"attribute": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"condition": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"value": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
					},
				},
			},
			"sso_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"login_link": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceIdentityProviderCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*Config)
	client, err := conf.IAMNoVersionClient(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("Error creating FlexibleEngine IAM without version client: %s", err)
	}

	if conf.DomainID == "" {
		return diag.Errorf("the domain_id must be specified in the provider configuration")
	}

	// Create a SAML protocol provider.
	opts := providers.CreateProviderOpts{
		Enabled:     d.Get("enabled").(bool),
		Description: d.Get("description").(string),
	}
	name := d.Get("name").(string)
	log.Printf("[DEBUG] Create identity options %s : %#v", name, opts)
	provider, err := providers.Create(client, name, opts)
	if err != nil {
		return diag.Errorf("Failed to create identity provider: %s", err)
	}

	d.SetId(provider.ID)

	// Create protocol and default mapping
	protocol := d.Get("protocol").(string)
	err = createProtocol(d, client)
	if err != nil {
		return diag.Errorf("error in creating provider protocol: %s", err)
	}

	// Import metadata, metadata only worked on saml protocol providers
	if protocol == protocolSAML {
		err = importMetadata(d, client, conf.DomainID)
		if err != nil {
			return diag.Errorf("error importing matedata into identity provider: %s", err)
		}
	} else if ac, ok := d.GetOk("openid_connect_config"); ok {
		// Create access config for oidc provider.
		accessConfigArr := ac.([]interface{})
		accessConfig := accessConfigArr[0].(map[string]interface{})

		accessType := accessConfig["access_type"].(string)
		createAccessTypeOpts := oidcconfig.CreateOpts{
			AccessMode: accessType,
			IdpURL:     accessConfig["provider_url"].(string),
			ClientID:   accessConfig["client_id"].(string),
			SigningKey: accessConfig["signing_key"].(string),
		}

		if accessType == "program_console" {
			scopes := utils.ExpandToStringList(accessConfig["scopes"].([]interface{}))
			createAccessTypeOpts.Scope = strings.Join(scopes, scopeSpilt)
			createAccessTypeOpts.AuthorizationEndpoint = accessConfig["authorization_endpoint"].(string)
			createAccessTypeOpts.ResponseType = accessConfig["response_type"].(string)
			createAccessTypeOpts.ResponseMode = accessConfig["response_mode"].(string)
		}
		log.Printf("[DEBUG] Create access type of provider: %#v", opts)

		_, err = oidcconfig.Create(client, provider.ID, createAccessTypeOpts)
		if err != nil {
			return diag.Errorf("Error creating the provider access config: %s", err)
		}
	}

	return resourceIdentityProviderRead(ctx, d, meta)
}

// importMetadata import metadata to provider, overwrite if it exists.
func importMetadata(d *schema.ResourceData, client *golangsdk.ServiceClient, domainID string) error {
	metadata := d.Get("metadata").(string)
	if metadata == "" {
		return nil
	}

	providerID := d.Get("name").(string)
	opts := metadatas.ImportOpts{
		DomainID: domainID,
		Metadata: metadata,
	}
	if _, err := metadatas.Import(client, providerID, protocolSAML, opts); err != nil {
		return fmt.Errorf("failed to import metadata: %s", err)
	}
	return nil
}

// createProtocol create protocol and default mapping
func createProtocol(d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	providerID := d.Get("name").(string)

	// Create default mapping
	defaultConversionRules := getDefaultConversionOpts()
	conversionRuleID := "mapping_" + providerID
	_, err := mappings.Create(client, conversionRuleID, *defaultConversionRules)
	if err != nil {
		return fmt.Errorf("error in creating default conversion rule: %s", err)
	}

	// Create protocol
	protocolName := d.Get("protocol").(string)
	_, err = protocols.Create(client, providerID, protocolName, conversionRuleID)
	if err != nil {
		// If fails to create protocols, then delete the mapping.
		mErr := multierror.Append(
			nil,
			err,
			mappings.Delete(client, conversionRuleID),
		)
		log.Printf("[ERROR] Error creating protocol, and the mapping that has been created. Error: %s", mErr)
		return fmt.Errorf("error creating identity provider protocol: %s", mErr.Error())
	}
	return nil
}

func getDefaultConversionOpts() *mappings.MappingOption {
	localRules := []mappings.LocalRule{
		{
			User: mappings.LocalRuleVal{
				Name: "FederationUser",
			},
		},
	}
	remoteRules := []mappings.RemoteRule{
		{
			Type: "__NAMEID__",
		},
	}

	opts := mappings.MappingOption{
		Rules: []mappings.MappingRule{
			{
				Local:  localRules,
				Remote: remoteRules,
			},
		},
	}
	return &opts
}

func resourceIdentityProviderRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*Config)
	client, err := conf.IAMNoVersionClient(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("Error creating FlexibleEngine IAM client without version number: %s", err)
	}

	provider, err := providers.Get(client, d.Id())
	if err != nil {
		return CheckDeletedDiag(d, err, "error obtaining identity provider")
	}

	// Query the protocol name from cloud
	protocol := queryProtocolName(d, client)
	url := generateLoginLink(conf, provider.ID, protocol)

	mErr := multierror.Append(err,
		d.Set("name", provider.ID),
		d.Set("protocol", protocol),
		d.Set("sso_type", provider.SsoType),
		d.Set("enabled", provider.Enabled),
		d.Set("login_link", url),
		d.Set("description", provider.Description),
	)

	// Query and set conversion rules
	conversionRuleID := "mapping_" + d.Id()
	conversions, err := mappings.Get(client, conversionRuleID)
	if err == nil {
		conversionRules := buildConversionRulesAttr(conversions)
		err = d.Set("conversion_rules", conversionRules)
		mErr = multierror.Append(mErr, err)
	}

	// Query and set metadata of the protocol SAML provider
	if protocol == protocolSAML {
		r, err := metadatas.Get(client, d.Id(), protocolSAML)
		if err == nil {
			err = d.Set("metadata", utils.HashAndHexEncode(r.Data))
			mErr = multierror.Append(mErr, err)
		}
	}

	// Query and set access type of the protocol OIDC provider
	if protocol == protocolOIDC {
		accessType, err := oidcconfig.Get(client, d.Id())
		if err == nil {
			scopes := strings.Split(accessType.Scope, scopeSpilt)
			accessTypeConfig := []interface{}{
				map[string]interface{}{
					"access_type":            accessType.AccessMode,
					"provider_url":           accessType.IdpURL,
					"client_id":              accessType.ClientID,
					"signing_key":            accessType.SigningKey,
					"scopes":                 scopes,
					"response_mode":          accessType.ResponseMode,
					"authorization_endpoint": accessType.AuthorizationEndpoint,
					"response_type":          accessType.ResponseType,
				},
			}

			err = d.Set("openid_connect_config", accessTypeConfig)
			mErr = multierror.Append(mErr, err)
		}
	}

	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting identity provider attributes: %s", err)
	}
	return nil
}

func buildConversionRulesAttr(conversions *mappings.IdentityMapping) []interface{} {
	conversionRules := make([]interface{}, 0, len(conversions.Rules))
	for _, v := range conversions.Rules {
		localRules := make([]map[string]interface{}, 0, len(v.Local))
		for _, localRule := range v.Local {
			username := localRule.User.Name
			r := map[string]interface{}{
				"username": username,
			}
			if localRule.Group != nil {
				r["group"] = localRule.Group.Name
			}
			localRules = append(localRules, r)
		}

		remoteRules := make([]map[string]interface{}, 0, len(v.Remote))
		for _, remoteRule := range v.Remote {
			r := map[string]interface{}{
				"attribute": remoteRule.Type,
			}
			if len(remoteRule.NotAnyOf) > 0 {
				r["condition"] = "not_any_of"
				r["value"] = remoteRule.NotAnyOf
			} else if len(remoteRule.AnyOneOf) > 0 {
				r["condition"] = "any_one_of"
				r["value"] = remoteRule.AnyOneOf
			}
			remoteRules = append(remoteRules, r)
		}

		rule := map[string]interface{}{
			"local":  localRules,
			"remote": remoteRules,
		}
		conversionRules = append(conversionRules, rule)
	}
	return conversionRules
}

// generateLoginLink generate login link base on config.domainID.
func generateLoginLink(conf *Config, id, protocol string) string {
	url := fmt.Sprintf("https://auth.%s/authui/federation/websso?domain_id=%s&idp=%s&protocol=%s",
		conf.Cloud, conf.DomainID, id, protocol)
	return url
}

func queryProtocolName(d *schema.ResourceData, client *golangsdk.ServiceClient) string {
	arr, err := protocols.List(client, d.Id())
	if err != nil {
		return ""
	}

	// The SAML protocol provider may not have protocol data,
	// so the default value is set to SAML.
	protocolName := protocolSAML
	if len(arr) > 0 {
		protocolName = arr[0].ID
	}
	return protocolName
}

func resourceIdentityProviderUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*Config)
	client, err := conf.IAMNoVersionClient(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("Error creating FlexibleEngine IAM client without version number: %s", err)
	}

	mErr := &multierror.Error{}
	if d.HasChanges("enabled", "description") {
		status := d.Get("enabled").(bool)
		description := d.Get("description").(string)
		opts := providers.UpdateOpts{
			Enabled:     &status,
			Description: &description,
		}
		log.Printf("[DEBUG] Update identity options %s : %#v", d.Id(), opts)

		_, err = providers.Update(client, d.Id(), opts)
		if err != nil {
			e := fmt.Errorf("Failed to update identity provider: %s", err)
			mErr = multierror.Append(mErr, e)
		}
	}

	if d.HasChange("metadata") {
		err = importMetadata(d, client, conf.DomainID)
		if err != nil {
			mErr = multierror.Append(mErr, err)
		}
	}

	if d.HasChange("openid_connect_config") && d.Get("protocol") == protocolOIDC {
		err = updateAccessConfig(client, d)
		if err != nil {
			mErr = multierror.Append(mErr, err)
		}
	}

	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error in updating provider: %s", err)
	}

	return resourceIdentityProviderRead(ctx, d, meta)
}

func updateAccessConfig(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	accessConfigArr := d.Get("openid_connect_config").([]interface{})
	if len(accessConfigArr) == 0 {
		return fmt.Errorf("the openid_connect_config is required for the OIDC provider")
	}
	accessConfig := accessConfigArr[0].(map[string]interface{})

	accessType := accessConfig["access_type"].(string)
	opts := oidcconfig.UpdateOpenIDConnectConfigOpts{
		AccessMode: accessType,
		IdpURL:     accessConfig["provider_url"].(string),
		ClientID:   accessConfig["client_id"].(string),
		SigningKey: accessConfig["signing_key"].(string),
	}

	if accessType == "program_console" {
		scopes := utils.ExpandToStringList(accessConfig["scopes"].([]interface{}))
		opts.Scope = strings.Join(scopes, scopeSpilt)
		opts.AuthorizationEndpoint = accessConfig["authorization_endpoint"].(string)
		opts.ResponseType = accessConfig["response_type"].(string)
		opts.ResponseMode = accessConfig["response_mode"].(string)
	}
	log.Printf("[DEBUG] Update access type of provider: %#v", opts)
	providerID := d.Id()
	_, err := oidcconfig.Update(client, providerID, opts)
	return err
}

func resourceIdentityProviderDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*Config)
	client, err := conf.IAMNoVersionClient(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("Error creating FlexibleEngine IAM client without version number: %s", err)
	}

	err = providers.Delete(client, d.Id())
	if err != nil {
		return CheckDeletedDiag(d, err, "Error deleting FlexibleEngine identity provider")
	}
	d.SetId("")
	return nil
}
