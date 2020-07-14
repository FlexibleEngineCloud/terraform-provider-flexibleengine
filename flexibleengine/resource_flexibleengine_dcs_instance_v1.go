package flexibleengine

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/dcs/v1/instances"
	"github.com/huaweicloud/golangsdk/openstack/dcs/v1/products"
)

func resourceDcsInstanceV1() *schema.Resource {
	return &schema.Resource{
		Create: resourceDcsInstancesV1Create,
		Read:   resourceDcsInstancesV1Read,
		Update: resourceDcsInstancesV1Update,
		Delete: resourceDcsInstancesV1Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"engine": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"engine_version": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"capacity": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"password": {
				Type:      schema.TypeString,
				Sensitive: true,
				Required:  true,
				ForceNew:  true,
			},
			"access_user": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"subnet_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"network_id"},
				Deprecated:    "use network_id instead",
			},
			"network_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"subnet_id"},
			},
			"available_zones": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"instance_type": {
				Type:          schema.TypeString,
				ConflictsWith: []string{"product_id"},
				ForceNew:      true,
				Optional:      true,
			},
			"product_id": {
				Type:          schema.TypeString,
				ConflictsWith: []string{"instance_type"},
				ForceNew:      true,
				Optional:      true,
				Computed:      true,
			},
			"maintain_begin": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"maintain_end": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"save_days": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"backup_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"begin_at": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"period_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"backup_at": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"order_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"security_group_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subnet_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"resource_spec_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"used_memory": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"internal_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"max_memory": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func getInstanceBackupPolicy(d *schema.ResourceData) *instances.InstanceBackupPolicy {
	backupAts := d.Get("backup_at").([]interface{})
	ats := make([]int, len(backupAts))
	for i, at := range backupAts {
		ats[i] = at.(int)
	}

	periodicalBackupPlan := instances.PeriodicalBackupPlan{
		BeginAt:    d.Get("begin_at").(string),
		PeriodType: d.Get("period_type").(string),
		BackupAt:   ats,
	}

	instanceBackupPolicy := &instances.InstanceBackupPolicy{
		SaveDays:             d.Get("save_days").(int),
		BackupType:           d.Get("backup_type").(string),
		PeriodicalBackupPlan: periodicalBackupPlan,
	}

	return instanceBackupPolicy
}

func getDcsProductId(client *golangsdk.ServiceClient, instanceType string) (string, error) {
	v, err := products.Get(client).Extract()
	if err != nil {
		return "", err
	}
	log.Printf("[DEBUG] Dcs get products : %+v", v)
	var FilteredPd []products.Product
	for _, pd := range v.Products {
		if instanceType != "" && pd.SpecCode != instanceType {
			continue
		}
		FilteredPd = append(FilteredPd, pd)
	}

	if len(FilteredPd) < 1 {
		return "", fmt.Errorf("Your query returned no results. Please change your filters and try again.")
	}

	return FilteredPd[0].ProductID, nil
}

func resourceDcsInstancesV1Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	dcsV1Client, err := config.dcsV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine dcs instance client: %s", err)
	}

	no_password_access := "true"
	if d.Get("access_user").(string) != "" || d.Get("password").(string) != "" {
		no_password_access = "false"
	}

	createOpts := &instances.CreateOps{
		Name:             d.Get("name").(string),
		Description:      d.Get("description").(string),
		Engine:           d.Get("engine").(string),
		EngineVersion:    d.Get("engine_version").(string),
		Capacity:         d.Get("capacity").(int),
		NoPasswordAccess: no_password_access,
		Password:         d.Get("password").(string),
		AccessUser:       d.Get("access_user").(string),
		VPCID:            d.Get("vpc_id").(string),
		SecurityGroupID:  d.Get("security_group_id").(string),
		AvailableZones:   getAllAvailableZones(d),
		MaintainBegin:    d.Get("maintain_begin").(string),
		MaintainEnd:      d.Get("maintain_end").(string),
	}

	subnet_id, subnet_ok := d.GetOk("subnet_id")
	network_id, network_ok := d.GetOk("network_id")
	if !subnet_ok && !network_ok {
		return fmt.Errorf("one of subnet_id or network_id must be configured")
	}
	if subnet_ok {
		createOpts.SubnetID = subnet_id.(string)
	} else {
		createOpts.SubnetID = network_id.(string)
	}

	product_id, product_ok := d.GetOk("product_id")
	instance_type, type_ok := d.GetOk("instance_type")
	if !product_ok && !type_ok {
		return fmt.Errorf("one of product_id or instance_type must be configured")
	}
	if product_ok {
		createOpts.ProductID = product_id.(string)
	} else {
		// Get Product ID
		createOpts.ProductID, err = getDcsProductId(dcsV1Client, instance_type.(string))
		if err != nil {
			return fmt.Errorf("Error get product id for dcs instance client: %s", err)
		}
	}

	if hasFilledOpt(d, "save_days") {
		createOpts.InstanceBackupPolicy = getInstanceBackupPolicy(d)
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	v, err := instances.Create(dcsV1Client, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine instance: %s", err)
	}
	log.Printf("[INFO] instance ID: %s", v.InstanceID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"CREATING"},
		Target:     []string{"RUNNING"},
		Refresh:    DcsInstancesV1StateRefreshFunc(dcsV1Client, v.InstanceID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for instance (%s) to become ready: %s",
			v.InstanceID, err)
	}

	// Store the instance ID now
	d.SetId(v.InstanceID)

	return resourceDcsInstancesV1Read(d, meta)
}

func resourceDcsInstancesV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	dcsV1Client, err := config.dcsV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine dcs instance client: %s", err)
	}
	v, err := instances.Get(dcsV1Client, d.Id()).Extract()
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Dcs instance %s: %+v", d.Id(), v)

	d.SetId(v.InstanceID)
	d.Set("name", v.Name)
	d.Set("engine", v.Engine)
	d.Set("engine_version", v.EngineVersion)
	d.Set("capacity", v.Capacity)
	d.Set("used_memory", v.UsedMemory)
	d.Set("max_memory", v.MaxMemory)
	d.Set("port", v.Port)
	d.Set("status", v.Status)
	d.Set("description", v.Description)
	d.Set("resource_spec_code", v.ResourceSpecCode)
	d.Set("internal_version", v.InternalVersion)
	d.Set("vpc_id", v.VPCID)
	d.Set("vpc_name", v.VPCName)
	d.Set("created_at", v.CreatedAt)
	d.Set("product_id", v.ProductID)
	d.Set("security_group_id", v.SecurityGroupID)
	d.Set("security_group_name", v.SecurityGroupName)
	d.Set("subnet_name", v.SubnetName)
	d.Set("user_id", v.UserID)
	d.Set("user_name", v.UserName)
	d.Set("order_id", v.OrderID)
	d.Set("maintain_begin", v.MaintainBegin)
	d.Set("maintain_end", v.MaintainEnd)
	d.Set("access_user", v.AccessUser)
	d.Set("ip", v.IP)

	if _, ok := d.GetOk("subnet_id"); ok {
		d.Set("subnet_id", v.SubnetID)
	} else {
		d.Set("network_id", v.SubnetID)
	}

	return nil
}

func resourceDcsInstancesV1Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	dcsV1Client, err := config.dcsV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error updating FlexibleEngine dcs instance client: %s", err)
	}
	var updateOpts instances.UpdateOpts
	if d.HasChange("name") {
		updateOpts.Name = d.Get("name").(string)
	}
	if d.HasChange("description") {
		description := d.Get("description").(string)
		updateOpts.Description = &description
	}
	if d.HasChange("maintain_begin") {
		maintain_begin := d.Get("maintain_begin").(string)
		updateOpts.MaintainBegin = maintain_begin
	}
	if d.HasChange("maintain_end") {
		maintain_end := d.Get("maintain_end").(string)
		updateOpts.MaintainEnd = maintain_end
	}
	if d.HasChange("security_group_id") {
		security_group_id := d.Get("security_group_id").(string)
		updateOpts.SecurityGroupID = security_group_id
	}

	err = instances.Update(dcsV1Client, d.Id(), updateOpts).Err
	if err != nil {
		return fmt.Errorf("Error updating FlexibleEngine Dcs Instance: %s", err)
	}

	return resourceDcsInstancesV1Read(d, meta)
}

func resourceDcsInstancesV1Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	dcsV1Client, err := config.dcsV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine dcs instance client: %s", err)
	}

	_, err = instances.Get(dcsV1Client, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "instance")
	}

	err = instances.Delete(dcsV1Client, d.Id()).ExtractErr()
	if err != nil {
		return fmt.Errorf("Error deleting FlexibleEngine instance: %s", err)
	}

	// Wait for the instance to delete before moving on.
	log.Printf("[DEBUG] Waiting for instance (%s) to delete", d.Id())

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"DELETING", "RUNNING"},
		Target:     []string{"DELETED"},
		Refresh:    DcsInstancesV1StateRefreshFunc(dcsV1Client, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for instance (%s) to delete: %s",
			d.Id(), err)
	}

	log.Printf("[DEBUG] Dcs instance %s deactivated.", d.Id())
	d.SetId("")
	return nil
}

func DcsInstancesV1StateRefreshFunc(client *golangsdk.ServiceClient, instanceID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		v, err := instances.Get(client, instanceID).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return v, "DELETED", nil
			}
			return nil, "", err
		}

		return v, v.Status, nil
	}
}
