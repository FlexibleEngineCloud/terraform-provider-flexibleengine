package flexibleengine

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/css/v1/snapshots"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCssSnapshotV1() *schema.Resource {
	return &schema.Resource{
		Create: resourceCssSnapshotV1Create,
		Read:   resourceCssSnapshotV1Read,
		Delete: resourceCssSnapshotV1Delete,
		Importer: &schema.ResourceImporter{
			State: resourceCssSnapshotImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"cluster_id": {
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
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"indices": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cluster_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"backup_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceCssSnapshotV1Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	cssClient, err := config.CssV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine CSS client: %s", err)
	}

	clusterId := d.Get("cluster_id").(string)
	createOpts := &snapshots.CreateOpts{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Indices:     d.Get("indices").(string),
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	snap, err := snapshots.Create(cssClient, createOpts, clusterId).Extract()
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine CSS snapshot: %s", err)
	}

	// Store the snapshot ID
	d.SetId(snap.ID)

	log.Printf("[DEBUG] Waiting for snapshot (%s) to complete", d.Id())
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"BUILDING"},
		Target:     []string{"COMPLETED"},
		Refresh:    cssSnapshotStateRefreshFunc(cssClient, clusterId, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 5 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for snapshot (%s) to complete: %s",
			d.Id(), err)
	}

	return resourceCssSnapshotV1Read(d, meta)
}

func resourceCssSnapshotV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	cssClient, err := config.CssV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine CSS client: %s", err)
	}

	clusterId := d.Get("cluster_id").(string)
	snapList, err := snapshots.List(cssClient, clusterId).Extract()
	if err != nil {
		return CheckDeleted(d, err, "snapshot")
	}

	// find the snapshot by ID
	var snap snapshots.Snapshot
	for _, v := range snapList {
		if v.ID == d.Id() {
			snap = v
			break
		}
	}
	if snap.ID == "" {
		log.Printf("[INFO] the snapshot %s does not exist", d.Id())
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] Retrieved the sanpshot %s: %+v", d.Id(), snap)

	d.Set("name", snap.Name)
	d.Set("description", snap.Description)
	d.Set("status", snap.Status)
	d.Set("indices", snap.Indices)
	d.Set("cluster_id", snap.ClusterID)
	d.Set("cluster_name", snap.ClusterName)
	// Method is more suitable for backup_type
	d.Set("backup_type", snap.Method)

	return nil
}

func resourceCssSnapshotV1Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	cssClient, err := config.CssV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine CSS storage client: %s", err)
	}

	clusterId := d.Get("cluster_id").(string)
	if err := snapshots.Delete(cssClient, clusterId, d.Id()).ExtractErr(); err != nil {
		return CheckDeleted(d, err, "snapshot")
	}

	d.SetId("")
	return nil
}

// cssSnapshotStateRefreshFunc returns a resource.StateRefreshFunc that is used to watch
// an CSS cluster snapshot.
func cssSnapshotStateRefreshFunc(client *golangsdk.ServiceClient, clusterId, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		snapList, err := snapshots.List(client, clusterId).Extract()
		if err != nil {
			return nil, "FAILED", err
		}

		// find the snapshot by ID
		var snap snapshots.Snapshot
		for _, v := range snapList {
			if v.ID == id {
				snap = v
				break
			}
		}

		if snap.ID == "" {
			return nil, "NOTEXIST", fmt.Errorf("The specified snapshot %s not exist", id)
		}

		return snap, snap.Status, nil
	}
}

func resourceCssSnapshotImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		err := fmt.Errorf("Invalid format specified for CSS snapshot. Format must be <cluster id>/<snapshot id>")
		return nil, err
	}
	clusterID := parts[0]
	snapshotID := parts[1]

	config := meta.(*Config)
	client, err := config.CssV1Client(GetRegion(d, config))
	if err != nil {
		return nil, fmt.Errorf("Error creating css client, err=%s", err)
	}

	// check the css cluster whether exists
	d.SetId(clusterID)
	if _, err := sendCssClusterV1ReadRequest(d, client); err != nil {
		return nil, err
	}

	d.Set("cluster_id", clusterID)
	d.SetId(snapshotID)

	return []*schema.ResourceData{d}, nil
}
