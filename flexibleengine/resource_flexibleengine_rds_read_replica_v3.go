package flexibleengine

import (
	"fmt"
	"log"
	"time"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/common/tags"
	"github.com/chnsz/golangsdk/openstack/rds/v3/instances"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceRdsReadReplicaInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceRdsReadReplicaInstanceCreate,
		Read:   resourceRdsReadReplicaInstanceRead,
		Update: resourceRdsReadReplicaInstanceUpdate,
		Delete: resourceRdsInstanceDelete,

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
				ForceNew: true,
			},

			"availability_zone": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"replica_of_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"flavor": {
				Type:     schema.TypeString,
				Required: true,
			},

			"volume": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"size": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"disk_encryption_id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},

			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"subnet_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"security_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"private_ips": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"public_ips": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"db": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"user_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"tags": tagsSchema(),
		},
	}
}

func buildRdsReplicaInstanceVolume(d *schema.ResourceData) *instances.Volume {
	var volume *instances.Volume
	volumeRaw := d.Get("volume").([]interface{})

	if len(volumeRaw) == 1 {
		volume = new(instances.Volume)
		volume.Type = volumeRaw[0].(map[string]interface{})["type"].(string)
		volume.Size = volumeRaw[0].(map[string]interface{})["size"].(int)
		// the size is optional and invalid for replica, but it's required in sdk
		// so just set 100 if not specified
		if volume.Size == 0 {
			volume.Size = 100
		}
	}
	log.Printf("[DEBUG] volume: %+v", volume)
	return volume
}

func resourceRdsReadReplicaInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.RdsV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine RDS client: %s ", err)
	}

	createOpts := instances.CreateReplicaOpts{
		Name:             d.Get("name").(string),
		ReplicaOfId:      d.Get("replica_of_id").(string),
		FlavorRef:        d.Get("flavor").(string),
		Region:           GetRegion(d, config),
		AvailabilityZone: d.Get("availability_zone").(string),
		Volume:           buildRdsReplicaInstanceVolume(d),
		DiskEncryptionId: d.Get("volume.0.disk_encryption_id").(string),
	}
	log.Printf("[DEBUG] Create replica instance Options: %#v", createOpts)

	resp, err := instances.CreateReplica(client, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating replica instance: %s ", err)
	}

	instance := resp.Instance
	d.SetId(instance.Id)
	instanceID := d.Id()
	if err := checkRDSInstanceJobFinish(client, resp.JobId, d.Timeout(schema.TimeoutCreate)); err != nil {
		return fmt.Errorf("Error creating instance (%s): %s", instanceID, err)
	}

	tagRaw := d.Get("tags").(map[string]interface{})
	if len(tagRaw) > 0 {
		tagList := expandResourceTags(tagRaw)
		err := tags.Create(client, "instances", instanceID, tagList).ExtractErr()
		if err != nil {
			return fmt.Errorf("Error setting tags of Rds read replica instance %s: %s", instanceID, err)
		}
	}

	return resourceRdsReadReplicaInstanceRead(d, meta)
}

func resourceRdsReadReplicaInstanceRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.RdsV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine RDS client: %s", err)
	}

	instanceID := d.Id()
	instance, err := getRdsInstanceByID(client, instanceID)
	if err != nil {
		return err
	}
	if instance.Id == "" {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] Retrieved rds read replica instance %s: %#v", instanceID, instance)
	d.Set("name", instance.Name)
	d.Set("flavor", instance.FlavorRef)
	d.Set("region", instance.Region)
	d.Set("private_ips", instance.PrivateIps)
	d.Set("public_ips", instance.PublicIps)
	d.Set("vpc_id", instance.VpcId)
	d.Set("subnet_id", instance.SubnetId)
	d.Set("security_group_id", instance.SecurityGroupId)
	d.Set("type", instance.Type)
	d.Set("status", instance.Status)
	d.Set("tags", tagsToMap(instance.Tags))

	az := expandAvailabilityZone(instance)
	d.Set("availability_zone", az)

	if primaryInstanceID, err := expandPrimaryInstanceID(instance); err == nil {
		d.Set("replica_of_id", primaryInstanceID)
	} else {
		return err
	}

	volumeList := make([]map[string]interface{}, 0, 1)
	volume := map[string]interface{}{
		"type":               instance.Volume.Type,
		"size":               instance.Volume.Size,
		"disk_encryption_id": instance.DiskEncryptionId,
	}
	volumeList = append(volumeList, volume)
	if err := d.Set("volume", volumeList); err != nil {
		return fmt.Errorf("[DEBUG] Error saving volume to RDS read replica instance (%s): %s", instanceID, err)
	}

	dbList := make([]map[string]interface{}, 0, 1)
	database := map[string]interface{}{
		"type":      instance.DataStore.Type,
		"version":   instance.DataStore.Version,
		"port":      instance.Port,
		"user_name": instance.DbUserName,
	}
	dbList = append(dbList, database)
	if err := d.Set("db", dbList); err != nil {
		return fmt.Errorf("[DEBUG] Error saving data base to RDS read replica instance (%s): %s", instanceID, err)
	}

	return nil
}

func resourceRdsReadReplicaInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.RdsV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine RDS client: %s ", err)
	}

	instanceID := d.Id()
	if err := updateRdsInstanceFlavor(d, client, instanceID); err != nil {
		return fmt.Errorf("[ERROR] %s", err)
	}

	if d.HasChange("tags") {
		tagErr := UpdateResourceTags(client, d, "instances", instanceID)
		if tagErr != nil {
			return fmt.Errorf("Error updating tags of RDS read replica instance: %s, err: %s", instanceID, tagErr)
		}
	}

	return resourceRdsReadReplicaInstanceRead(d, meta)
}

func updateRdsInstanceFlavor(d *schema.ResourceData, client *golangsdk.ServiceClient, instanceID string) error {
	if !d.HasChange("flavor") {
		return nil
	}

	resizeFlavor := instances.SpecCode{
		Speccode: d.Get("flavor").(string),
	}
	var resizeFlavorOpts instances.ResizeFlavorOpts
	resizeFlavorOpts.ResizeFlavor = &resizeFlavor

	_, err := instances.Resize(client, resizeFlavorOpts, instanceID).Extract()
	if err != nil {
		return fmt.Errorf("Error updating instance Flavor from result: %s ", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"MODIFYING"},
		Target:       []string{"ACTIVE"},
		Refresh:      rdsInstanceStateRefreshFunc(client, instanceID),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        15 * time.Second,
		PollInterval: 15 * time.Second,
	}
	if _, err = stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for instance (%s) flavor to be Updated: %s ", instanceID, err)
	}
	return nil
}

func resourceRdsInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.RdsV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine RDS client: %s ", err)
	}

	id := d.Id()
	log.Printf("[DEBUG] Deleting Instance %s", id)
	result := instances.Delete(client, id)
	if result.Err != nil {
		return result.Err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    rdsInstanceStateRefreshFunc(client, id),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      15 * time.Second,
		MinTimeout: 5 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for rds instance (%s) to be deleted: %s",
			id, err)
	}

	log.Printf("[DEBUG] Successfully deleted rds instance %s", id)
	return nil
}

func getRdsInstanceByID(client *golangsdk.ServiceClient, instanceID string) (*instances.RdsInstanceResponse, error) {
	listOpts := instances.ListOpts{
		Id: instanceID,
	}
	pages, err := instances.List(client, listOpts).AllPages()
	if err != nil {
		return nil, fmt.Errorf("An error occured while querying rds instance %s: %s", instanceID, err)
	}

	resp, err := instances.ExtractRdsInstances(pages)
	if err != nil {
		return nil, err
	}

	instanceList := resp.Instances
	if len(instanceList) == 0 {
		// return an empty rds instance
		log.Printf("[WARN] can not find the specified rds instance %s", instanceID)
		instance := new(instances.RdsInstanceResponse)
		return instance, nil
	}

	if len(instanceList) > 1 {
		return nil, fmt.Errorf("retrieving more than one rds instance by %s", instanceID)
	}
	if instanceList[0].Id != instanceID {
		return nil, fmt.Errorf("the id of rds instance was expected %s, but got %s",
			instanceID, instanceList[0].Id)
	}

	return &instanceList[0], nil
}

func expandAvailabilityZone(resp *instances.RdsInstanceResponse) string {
	node := resp.Nodes[0]
	return node.AvailabilityZone
}

func expandPrimaryInstanceID(resp *instances.RdsInstanceResponse) (string, error) {
	relatedInst := resp.RelatedInstance
	for _, relate := range relatedInst {
		if relate.Type == "replica_of" {
			return relate.Id, nil
		}
	}
	return "", fmt.Errorf("Error getting primary instance id for replica %s", resp.Id)
}

func checkRDSInstanceJobFinish(client *golangsdk.ServiceClient, jobID string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"Running"},
		Target:       []string{"Completed", "Failed"},
		Refresh:      rdsInstanceJobRefreshFunc(client, jobID),
		Timeout:      timeout,
		Delay:        20 * time.Second,
		PollInterval: 10 * time.Second,
	}
	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for RDS instance (%s) job to be completed: %s ", jobID, err)
	}
	return nil
}

func rdsInstanceJobRefreshFunc(client *golangsdk.ServiceClient, jobID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		jobOpts := instances.RDSJobOpts{
			JobID: jobID,
		}
		jobList, err := instances.GetRDSJob(client, jobOpts).Extract()
		if err != nil {
			return nil, "FOUND ERROR", err
		}

		return jobList.Job, jobList.Job.Status, nil
	}
}

func rdsInstanceStateRefreshFunc(client *golangsdk.ServiceClient, instanceID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		instance, err := getRdsInstanceByID(client, instanceID)
		if err != nil {
			return nil, "FOUND ERROR", err
		}
		if instance.Id == "" {
			return instance, "DELETED", nil
		}

		return instance, instance.Status, nil
	}
}
