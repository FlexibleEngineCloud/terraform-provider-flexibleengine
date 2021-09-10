package flexibleengine

import (
	"fmt"
	"log"
	"time"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/networking/v2/extensions/elbaas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func waitForELBJobSuccess(networkingClient *golangsdk.ServiceClient, j *elbaas.Job, timeout time.Duration) (*elbaas.JobInfo, error) {
	jobId := j.JobId
	target := "SUCCESS"
	pending := []string{"INIT", "RUNNING"}

	log.Printf("[DEBUG] Waiting for elbaas job %s to become %s.", jobId, target)

	ji, err := waitForELBResource(networkingClient, "job", j.JobId, target, pending, timeout, getELBJobInfo)
	if err == nil {
		return ji.(*elbaas.JobInfo), nil
	}
	return nil, err
}

func getELBJobInfo(networkingClient *golangsdk.ServiceClient, uri string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		info, err := elbaas.QueryJobInfo(networkingClient, uri).Extract()
		if err != nil {
			return nil, "", err
		}

		return info, info.Status, nil
	}
}

type getELBResource func(networkingClient *golangsdk.ServiceClient, id string) resource.StateRefreshFunc

func waitForELBResource(networkingClient *golangsdk.ServiceClient, name string, id string, target string, pending []string, timeout time.Duration, f getELBResource) (interface{}, error) {

	stateConf := &resource.StateChangeConf{
		Target:     []string{target},
		Pending:    pending,
		Refresh:    f(networkingClient, id),
		Timeout:    timeout,
		Delay:      5 * time.Second,
		MinTimeout: 1 * time.Second,
	}

	o, err := stateConf.WaitForState()
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			return nil, fmt.Errorf("Error: elbaas %s %s not found: %s", name, id, err)
		}
		return nil, fmt.Errorf("Error waiting for elbaas %s %s to become %s: %s", name, id, target, err)
	}

	return o, nil
}
