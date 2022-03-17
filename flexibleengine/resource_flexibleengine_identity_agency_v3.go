package flexibleengine

import (
	"bytes"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/identity/v3/agency"
	sdkprojects "github.com/chnsz/golangsdk/openstack/identity/v3/projects"
	sdkroles "github.com/chnsz/golangsdk/openstack/identity/v3/roles"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceIdentityAgencyV3() *schema.Resource {
	return &schema.Resource{
		Create: resourceIdentityAgencyV3Create,
		Read:   resourceIdentityAgencyV3Read,
		Update: resourceIdentityAgencyV3Update,
		Delete: resourceIdentityAgencyV3Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"delegated_domain_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"delegated_service_name"},
			},
			"delegated_service_name": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile("^op_svc_[A-Za-z]+"),
					"the value must start with op_svc_, for example, op_svc_obs"),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"duration": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "FOREVER",
			},
			"expire_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project_role": {
				Type:         schema.TypeSet,
				Optional:     true,
				AtLeastOneOf: []string{"project_role", "domain_roles"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"roles": {
							Type:     schema.TypeSet,
							Required: true,
							MinItems: 1,
							MaxItems: 25,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
						},
						"project": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
				Set: resourceIdentityAgencyProRoleHash,
			},
			"domain_roles": {
				Type:     schema.TypeSet,
				Optional: true,
				MinItems: 1,
				MaxItems: 25,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
		},
	}
}

func resourceIdentityAgencyProRoleHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	buf.WriteString(fmt.Sprintf("%s-", m["project"].(string)))

	r := m["roles"].(*schema.Set).List()
	s := make([]string, len(r))
	for i, item := range r {
		s[i] = item.(string)
	}
	buf.WriteString(strings.Join(s, "-"))

	return schema.HashString(buf.String())
}

func listProjectsOfDomain(client *golangsdk.ServiceClient, domainID string) (map[string]string, error) {
	opts := sdkprojects.ListOpts{
		DomainID: domainID,
	}
	allPages, err := sdkprojects.List(client, &opts).AllPages()
	if err != nil {
		return nil, fmt.Errorf("List projects failed, err=%s", err)
	}

	all, err := sdkprojects.ExtractProjects(allPages)
	if err != nil {
		return nil, fmt.Errorf("Extract projects failed, err=%s", err)
	}

	r := make(map[string]string, len(all))
	for _, item := range all {
		r[item.Name] = item.ID
	}
	log.Printf("[TRACE] projects = %#v\n", r)
	return r, nil
}

func listRolesOfDomain(client *golangsdk.ServiceClient, domainID string) (map[string]string, error) {
	opts := sdkroles.ListOpts{
		DomainID: domainID,
	}
	allPages, err := sdkroles.List(client, &opts).AllPages()
	if err != nil {
		return nil, fmt.Errorf("List roles failed, err=%s", err)
	}

	all, err := sdkroles.ExtractRoles(allPages)
	if err != nil {
		return nil, fmt.Errorf("Extract roles failed, err=%s", err)
	}
	if len(all) == 0 {
		return nil, nil
	}

	r := make(map[string]string, len(all))
	for _, item := range all {
		if name := item.DisplayName; name != "" {
			r[name] = item.ID
		} else {
			log.Printf("[WARN] role %s without displayname", item.Name)
		}
	}
	log.Printf("[TRACE] list roles = %#v, len=%d\n", r, len(r))
	return r, nil
}

func allRolesOfDomain(client *golangsdk.ServiceClient, domainID string) (map[string]string, error) {
	roles, err := listRolesOfDomain(client, "")
	if err != nil {
		return nil, fmt.Errorf("Error listing global roles, err=%s", err)
	}

	customRoles, err := listRolesOfDomain(client, domainID)
	if err != nil {
		return nil, fmt.Errorf("Error listing domain's custom roles, err=%s", err)
	}

	if roles == nil {
		return customRoles, nil
	}

	if customRoles == nil {
		return roles, nil
	}

	for k, v := range customRoles {
		roles[k] = v
	}
	return roles, nil
}

func changeToPRPair(prs *schema.Set) (r map[string]bool) {
	r = make(map[string]bool)
	for _, v := range prs.List() {
		pr := v.(map[string]interface{})

		pn := pr["project"].(string)
		rs := pr["roles"].(*schema.Set)
		for _, role := range rs.List() {
			r[pn+"|"+role.(string)] = true
		}
	}
	return
}

func diffChangeOfProjectRole(old, newv *schema.Set) (delete, add []string) {
	delete = make([]string, 0)
	add = make([]string, 0)

	oldprs := changeToPRPair(old)
	newprs := changeToPRPair(newv)

	for k := range oldprs {
		if _, ok := newprs[k]; !ok {
			delete = append(delete, k)
		}
	}

	for k := range newprs {
		if _, ok := oldprs[k]; !ok {
			add = append(add, k)
		}
	}
	return
}

func resourceIdentityAgencyV3Create(d *schema.ResourceData, meta interface{}) error {
	prs := d.Get("project_role").(*schema.Set)
	drs := d.Get("domain_roles").(*schema.Set)
	if prs.Len() == 0 && drs.Len() == 0 {
		return fmt.Errorf("One or both of project_role and domain_roles must be input")
	}

	config := meta.(*Config)
	region := GetRegion(d, config)
	client, err := config.IAMV3Client(region)
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine IAM client: %s", err)
	}
	identityClient, err := config.IdentityV3Client(region)
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine identity client: %s", err)
	}

	domainID := config.DomainID
	if domainID == "" {
		return fmt.Errorf("the domain_id must be specified in the provider configuration")
	}

	opts := agency.CreateOpts{
		Name:        d.Get("name").(string),
		DomainID:    domainID,
		Description: d.Get("description").(string),
		Duration:    d.Get("duration").(string),
	}

	if v, ok := d.GetOk("delegated_domain_name"); ok {
		opts.DelegatedDomain = v.(string)
	} else {
		opts.DelegatedDomain = d.Get("delegated_service_name").(string)
	}

	log.Printf("[DEBUG] Create Identity-Agency Options: %#v", opts)
	a, err := agency.Create(client, opts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating Identity-Agency: %s", err)
	}

	d.SetId(a.ID)

	projects, err := listProjectsOfDomain(identityClient, domainID)
	if err != nil {
		return fmt.Errorf("Error querying the projects, err=%s", err)
	}

	roles, err := allRolesOfDomain(identityClient, domainID)
	if err != nil {
		return fmt.Errorf("Error querying the roles, err=%s", err)
	}

	agencyID := a.ID
	for _, v := range prs.List() {
		pr := v.(map[string]interface{})
		pn := pr["project"].(string)
		pid, ok := projects[pn]
		if !ok {
			return fmt.Errorf("The project(%s) is not exist", pn)
		}

		rs := pr["roles"].(*schema.Set)
		for _, role := range rs.List() {
			r := role.(string)
			rid, ok := roles[r]
			if !ok {
				return fmt.Errorf("The role(%s) is not exist", r)
			}

			err = agency.AttachRoleByProject(client, agencyID, pid, rid).ExtractErr()
			if err != nil {
				return fmt.Errorf("Error attaching role(%s) by project{%s} to agency(%s), err=%s",
					rid, pid, agencyID, err)
			}
		}
	}

	for _, role := range drs.List() {
		r := role.(string)
		rid, ok := roles[r]
		if !ok {
			return fmt.Errorf("The role(%s) is not exist", r)
		}

		err = agency.AttachRoleByDomain(client, agencyID, domainID, rid).ExtractErr()
		if err != nil {
			return fmt.Errorf("Error attaching role(%s) by domain{%s} to agency(%s), err=%s",
				rid, domainID, agencyID, err)
		}
	}

	return resourceIdentityAgencyV3Read(d, meta)
}

func resourceIdentityAgencyV3Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	region := GetRegion(d, config)
	client, err := config.IAMV3Client(region)
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine IAM client: %s", err)
	}
	identityClient, err := config.IdentityV3Client(region)
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine identity client: %s", err)
	}

	a, err := agency.Get(client, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "Identity-Agency")
	}
	log.Printf("[DEBUG] Retrieved Identity-Agency %s: %#v", d.Id(), a)

	d.Set("name", a.Name)
	d.Set("description", a.Description)
	d.Set("duration", a.Duration)
	d.Set("expire_time", a.ExpireTime)
	d.Set("create_time", a.CreateTime)

	if ok, err := regexp.MatchString("^op_svc_[A-Za-z]+$", a.DelegatedDomainName); err != nil {
		log.Printf("[ERROR] Regexp error, err= %s", err)
	} else if ok {
		d.Set("delegated_service_name", a.DelegatedDomainName)
	} else {
		d.Set("delegated_domain_name", a.DelegatedDomainName)
	}

	projects, err := listProjectsOfDomain(identityClient, a.DomainID)
	if err != nil {
		return fmt.Errorf("Error querying the projects, err=%s", err)
	}
	agencyID := d.Id()
	prs := schema.Set{F: resourceIdentityAgencyProRoleHash}
	for pn, pid := range projects {
		roles, err := agency.ListRolesAttachedOnProject(client, agencyID, pid).ExtractRoles()
		if err != nil && !isResourceNotFound(err) {
			return fmt.Errorf("Error querying the roles attached on project(%s), err=%s", pn, err)
		}
		if len(roles) == 0 {
			continue
		}
		v := schema.Set{F: schema.HashString}
		for _, role := range roles {
			v.Add(role.DisplayName)
		}
		prs.Add(map[string]interface{}{
			"project": pn,
			"roles":   &v,
		})
	}
	err = d.Set("project_role", &prs)
	if err != nil {
		log.Printf("[ERROR]Set project_role failed, err=%s", err)
	}

	roles, err := agency.ListRolesAttachedOnDomain(client, agencyID, a.DomainID).ExtractRoles()
	if err != nil && !isResourceNotFound(err) {
		return fmt.Errorf("Error querying the roles attached on domain, err=%s", err)
	}
	if len(roles) != 0 {
		v := schema.Set{F: schema.HashString}
		for _, role := range roles {
			v.Add(role.DisplayName)
		}
		err = d.Set("domain_roles", &v)
		if err != nil {
			log.Printf("[ERROR]Set domain_roles failed, err=%s", err)
		}
	}

	return nil
}

func resourceIdentityAgencyV3Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	region := GetRegion(d, config)
	client, err := config.IAMV3Client(region)
	if err != nil {
		return fmt.Errorf("Error creating IAM client: %s", err)
	}
	identityClient, err := config.IdentityV3Client(region)
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine identity client: %s", err)
	}

	aID := d.Id()
	domainID := config.DomainID
	if domainID == "" {
		return fmt.Errorf("the domain_id must be specified in the provider configuration")
	}

	if d.HasChange("delegated_domain_name") || d.HasChange("delegated_service_name") ||
		d.HasChange("description") || d.HasChange("duration") {
		updateOpts := agency.UpdateOpts{
			Description: d.Get("description").(string),
			Duration:    d.Get("duration").(string),
		}
		if v, ok := d.GetOk("delegated_domain_name"); ok {
			updateOpts.DelegatedDomain = v.(string)
		} else {
			updateOpts.DelegatedDomain = d.Get("delegated_service_name").(string)
		}

		log.Printf("[DEBUG] Updating Identity-Agency %s with options: %#v", aID, updateOpts)
		timeout := d.Timeout(schema.TimeoutUpdate)
		err = resource.Retry(timeout, func() *resource.RetryError {
			_, err := agency.Update(client, aID, updateOpts).Extract()
			if err != nil {
				return checkForRetryableError(err)
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("Error updating Identity-Agency %s: %s", aID, err)
		}
	}

	var roles map[string]string
	if d.HasChange("project_role") || d.HasChange("domain_roles") {
		roles, err = allRolesOfDomain(identityClient, domainID)
		if err != nil {
			return fmt.Errorf("Error querying the roles, err=%s", err)
		}
	}

	if d.HasChange("project_role") {
		projects, err := listProjectsOfDomain(identityClient, domainID)
		if err != nil {
			return fmt.Errorf("Error querying the projects, err=%s", err)
		}

		o, n := d.GetChange("project_role")
		deleteprs, addprs := diffChangeOfProjectRole(o.(*schema.Set), n.(*schema.Set))
		for _, v := range deleteprs {
			pr := strings.Split(v, "|")
			pid, ok := projects[pr[0]]
			if !ok {
				return fmt.Errorf("The project(%s) is not exist", pr[0])
			}
			rid, ok := roles[pr[1]]
			if !ok {
				return fmt.Errorf("The role(%s) is not exist", pr[1])
			}

			err = agency.DetachRoleByProject(client, aID, pid, rid).ExtractErr()
			if err != nil && !isResourceNotFound(err) {
				return fmt.Errorf("Error detaching role(%s) by project{%s} from agency(%s), err=%s",
					rid, pid, aID, err)
			}
		}

		for _, v := range addprs {
			pr := strings.Split(v, "|")
			pid, ok := projects[pr[0]]
			if !ok {
				return fmt.Errorf("The project(%s) is not exist", pr[0])
			}
			rid, ok := roles[pr[1]]
			if !ok {
				return fmt.Errorf("The role(%s) is not exist", pr[1])
			}

			err = agency.AttachRoleByProject(client, aID, pid, rid).ExtractErr()
			if err != nil {
				return fmt.Errorf("Error attaching role(%s) by project{%s} to agency(%s), err=%s",
					rid, pid, aID, err)
			}
		}
	}

	if d.HasChange("domain_roles") {
		o, n := d.GetChange("domain_roles")
		oldr := o.(*schema.Set)
		newr := n.(*schema.Set)

		for _, r := range oldr.Difference(newr).List() {
			rid, ok := roles[r.(string)]
			if !ok {
				return fmt.Errorf("The role(%s) is not exist", r.(string))
			}

			err = agency.DetachRoleByDomain(client, aID, domainID, rid).ExtractErr()
			if err != nil && !isResourceNotFound(err) {
				return fmt.Errorf("Error detaching role(%s) by domain{%s} from agency(%s), err=%s",
					rid, domainID, aID, err)
			}
		}

		for _, r := range newr.Difference(oldr).List() {
			rid, ok := roles[r.(string)]
			if !ok {
				return fmt.Errorf("The role(%s) is not exist", r.(string))
			}

			err = agency.AttachRoleByDomain(client, aID, domainID, rid).ExtractErr()
			if err != nil {
				return fmt.Errorf("Error attaching role(%s) by domain{%s} to agency(%s), err=%s",
					rid, domainID, aID, err)
			}
		}
	}
	return resourceIdentityAgencyV3Read(d, meta)
}

func resourceIdentityAgencyV3Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.IAMV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating IAM client: %s", err)
	}

	rID := d.Id()
	log.Printf("[DEBUG] Deleting Identity-Agency %s", rID)

	timeout := d.Timeout(schema.TimeoutDelete)
	err = resource.Retry(timeout, func() *resource.RetryError {
		err := agency.Delete(client, rID).ExtractErr()
		if err != nil {
			return checkForRetryableError(err)
		}
		return nil
	})
	if err != nil {
		if isResourceNotFound(err) {
			log.Printf("[INFO] deleting an unavailable Identity-Agency: %s", rID)
			return nil
		}
		return fmt.Errorf("Error deleting Identity-Agency %s: %s", rID, err)
	}

	return nil
}
