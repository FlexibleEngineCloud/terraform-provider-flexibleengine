package elbaas

import (
	"github.com/gophercloud/gophercloud"
)

type Job struct {
	Uri   string `json: "uri"`
	JobId string `json:"job_id"`
}

type JobResult struct {
	gophercloud.Result
}

func (r JobResult) Extract() (*Job, error) {
	j := &Job{}
	err := r.ExtractInto(j)
	return j, err
}

type JobInfo struct {
	Status     string                 `json:"status"`
	Entities   map[string]interface{} `json:"entities"`
	JobId      string                 `json:"job_id"`
	JobType    string                 `json:"job_type"`
	ErrorCode  string                 `json:"error_code"`
	FailReason string                 `json:"fail_reason"`
}

type JobInfoResult struct {
	gophercloud.Result
}

func (r JobInfoResult) Extract() (*JobInfo, error) {
	j := &JobInfo{}
	err := r.ExtractInto(j)
	return j, err
}

func QueryJobInfo(c *gophercloud.ServiceClient, jobId string) (r JobInfoResult) {
	_, r.Err = c.Get(c.ServiceURL("jobs", jobId), &r.Body, nil)
	return
}
