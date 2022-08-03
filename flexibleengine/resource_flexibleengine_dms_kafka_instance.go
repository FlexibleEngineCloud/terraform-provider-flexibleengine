package flexibleengine

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/dms/v1/instances"
	dmsv2 "github.com/chnsz/golangsdk/openstack/dms/v2/kafka/instances"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

const (
	ssdSpecCode = "dms.physical.storage.ultra"
	sasSpecCode = "dms.physical.storage.high"
)

func resourceDmsKafkaInstances() *schema.Resource {
	return &schema.Resource{
		Create: resourceDmsKafkaInstancesCreate,
		Read:   resourceDmsKafkaInstancesRead,
		Update: resourceDmsKafkaInstancesUpdate,
		Delete: resourceDmsKafkaInstancesDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"engine_version": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "2.3.0",
			},

			"bandwidth": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice(
					[]string{"100MB", "300MB", "600MB", "1200MB"},
					false,
				),
			},
			"product_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"storage_space": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"storage_spec_code": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  ssdSpecCode,
				ValidateFunc: validation.StringInSlice(
					[]string{ssdSpecCode, sasSpecCode},
					false,
				),
			},

			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"network_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"availability_zones": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"manager_user": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				RequiredWith: []string{"manager_password"},
			},
			"manager_password": {
				Type:         schema.TypeString,
				Sensitive:    true,
				Optional:     true,
				ForceNew:     true,
				RequiredWith: []string{"manager_user"},
			},
			"access_user": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				RequiredWith: []string{"password"},
			},
			"password": {
				Type:         schema.TypeString,
				Sensitive:    true,
				Optional:     true,
				ForceNew:     true,
				RequiredWith: []string{"access_user"},
			},
			"maintain_begin": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				RequiredWith: []string{"maintain_end"},
			},
			"maintain_end": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				RequiredWith: []string{"maintain_begin"},
			},

			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"engine": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"engine_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"product_spec_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"partition_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"vpc_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subnet_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"security_group_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"node_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"manegement_connect_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"connect_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"ssl_enable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"used_storage_space": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceDmsKafkaInstancesCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	dmsV1Client, err := config.DmsV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine DMS v1 client: %s", err)
	}

	sslEnable := false
	if d.Get("access_user").(string) != "" || d.Get("password").(string) != "" {
		sslEnable = true
	}
	createOpts := &instances.CreateOps{
		Engine:          "kafka",
		EngineVersion:   d.Get("engine_version").(string),
		Name:            d.Get("name").(string),
		Description:     d.Get("description").(string),
		Specification:   d.Get("bandwidth").(string),
		ProductID:       d.Get("product_id").(string),
		StorageSpace:    d.Get("storage_space").(int),
		StorageSpecCode: d.Get("storage_spec_code").(string),
		VPCID:           d.Get("vpc_id").(string),
		SubnetID:        d.Get("network_id").(string),
		SecurityGroupID: d.Get("security_group_id").(string),
		AvailableZones:  utils.ExpandToStringList(d.Get("availability_zones").([]interface{})),
		MaintainBegin:   d.Get("maintain_begin").(string),
		MaintainEnd:     d.Get("maintain_end").(string),

		AccessUser:       d.Get("access_user").(string),
		KafkaManagerUser: d.Get("manager_user").(string),
		SslEnable:        sslEnable,
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	// Add password here so it wouldn't go in the above log entry
	createOpts.Password = d.Get("password").(string)
	createOpts.KafkaManagerPassword = d.Get("manager_password").(string)

	v, err := instances.Create(dmsV1Client, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine DMS instance: %s", err)
	}
	log.Printf("[INFO] instance ID: %s", v.InstanceID)

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"CREATING"},
		Target:       []string{"RUNNING"},
		Refresh:      dmsKafkaInstancesStateRefreshFunc(dmsV1Client, v.InstanceID),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        120 * time.Second,
		PollInterval: 15 * time.Second,
	}
	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for instance (%s) to become ready: %s",
			v.InstanceID, err)
	}

	// Store the instance ID now
	d.SetId(v.InstanceID)

	return resourceDmsKafkaInstancesRead(d, meta)
}

func resourceDmsKafkaInstancesRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	region := GetRegion(d, config)
	dmsV2Client, err := config.DmsV2Client(region)
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine DMS v2 client: %s", err)
	}

	v, err := dmsv2.Get(dmsV2Client, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "DMS instance")
	}

	log.Printf("[DEBUG] Dms instance %s: %+v", d.Id(), v)

	partitionNum, _ := strconv.Atoi(v.PartitionNum)

	var mErr *multierror.Error
	mErr = multierror.Append(err,
		d.Set("region", region),
		d.Set("availability_zones", v.AvailableZones),
		d.Set("status", v.Status),
		d.Set("name", v.Name),
		d.Set("description", v.Description),
		d.Set("engine", v.Engine),
		d.Set("engine_version", v.EngineVersion),
		d.Set("engine_type", v.Type),
		d.Set("bandwidth", v.Specification),
		d.Set("partition_num", partitionNum),
		d.Set("product_id", v.ProductID),
		d.Set("product_spec_code", v.ResourceSpecCode),
		d.Set("storage_spec_code", v.StorageSpecCode),
		d.Set("storage_space", v.TotalStorageSpace),
		d.Set("used_storage_space", v.UsedStorageSpace),
		d.Set("vpc_id", v.VPCID),
		d.Set("vpc_name", v.VPCName),
		d.Set("network_id", v.SubnetID),
		d.Set("subnet_name", v.SubnetName),
		d.Set("security_group_id", v.SecurityGroupID),
		d.Set("security_group_name", v.SecurityGroupName),
		d.Set("node_num", v.NodeNum),
		d.Set("manegement_connect_address", v.ManagementConnectAddress),
		d.Set("connect_address", v.ConnectAddress),
		d.Set("port", v.Port),
		d.Set("maintain_begin", v.MaintainBegin),
		d.Set("maintain_end", v.MaintainEnd),
		d.Set("ssl_enable", v.SslEnable),
		d.Set("created_at", setResourceTimestamp(v.CreatedAt)),
	)

	if mErr.ErrorOrNil() != nil {
		return fmt.Errorf("Error setting DMS product attributes: %s", mErr)
	}

	return nil
}

func setResourceTimestamp(stamp string) string {
	mSeconds, err := strconv.ParseInt(stamp, 10, 64)
	if err != nil {
		log.Printf("[WARN] failed to convert string %s to int64: %s", stamp, err)
		return ""
	}
	stampTime := time.Unix(mSeconds/1000, 0)
	return stampTime.Format(time.RFC3339)
}

func resourceDmsKafkaInstancesUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	dmsV1Client, err := config.DmsV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error updating FlexibleEngine DMS instance client: %s", err)
	}

	//lintignore:R019
	if d.HasChanges("name", "description", "maintain_begin", "maintain_end", "security_group_id") {
		var updateOpts instances.UpdateOpts
		if d.HasChange("name") {
			updateOpts.Name = d.Get("name").(string)
		}
		if d.HasChange("description") {
			description := d.Get("description").(string)
			updateOpts.Description = &description
		}
		if d.HasChange("maintain_begin") {
			updateOpts.MaintainBegin = d.Get("maintain_begin").(string)
		}
		if d.HasChange("maintain_end") {
			updateOpts.MaintainEnd = d.Get("maintain_end").(string)
		}
		if d.HasChange("security_group_id") {
			updateOpts.SecurityGroupID = d.Get("security_group_id").(string)
		}

		err = instances.Update(dmsV1Client, d.Id(), updateOpts).Err
		if err != nil {
			return fmt.Errorf("Error updating FlexibleEngine Dms Instance: %s", err)
		}
	}

	return resourceDmsKafkaInstancesRead(d, meta)
}

func resourceDmsKafkaInstancesDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	dmsV1Client, err := config.DmsV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine DMS instance client: %s", err)
	}

	err = instances.Delete(dmsV1Client, d.Id()).ExtractErr()
	if err != nil {
		return fmt.Errorf("Error deleting FlexibleEngine instance: %s", err)
	}

	// Wait for the instance to delete before moving on.
	log.Printf("[DEBUG] Waiting for instance (%s) to delete", d.Id())

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"DELETING", "RUNNING"},
		Target:     []string{"DELETED"},
		Refresh:    dmsKafkaInstancesStateRefreshFunc(dmsV1Client, d.Id()),
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

	log.Printf("[DEBUG] Dms instance %s deactivated.", d.Id())
	d.SetId("")
	return nil
}

func dmsKafkaInstancesStateRefreshFunc(client *golangsdk.ServiceClient, instanceID string) resource.StateRefreshFunc {
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
