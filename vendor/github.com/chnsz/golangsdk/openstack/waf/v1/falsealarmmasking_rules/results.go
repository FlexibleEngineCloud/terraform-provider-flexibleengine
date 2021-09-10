package falsealarmmasking_rules

import (
	"github.com/chnsz/golangsdk"
)

type ListResponse struct {
	// Total number of the Rules
	Total int `json:"total"`

	// List of AlarmMasking
	Items []AlarmMasking `json:"items"`
}

type AlarmMasking struct {
	// the ID of a false alarm masking rule
	Id string `json:"id"`
	// the policy ID
	PolicyID string `json:"policy_id"`
	// a misreported URL excluding a domain name
	Path string `json:"path"`
	// the event ID
	EventID string `json:"event_id"`
	// the event ID
	EventType string `json:"event_type"`
	// the rule ID, which consists of six digits and cannot be empty
	Rule string `json:"rule"`
	// the time when a false alarm masking rule is added
	TimeStamp int `json:"timestamp"`
}

type commonResult struct {
	golangsdk.Result
}

// Extract is a function that accepts a result and extracts a falsealarmmasking rule.
func (r commonResult) Extract() (*AlarmMasking, error) {
	var response AlarmMasking
	err := r.ExtractInto(&response)
	return &response, err
}

// CreateResult represents the result of a create operation. Call its Extract
// method to interpret it as a False Alarm Masking rule.
type CreateResult struct {
	commonResult
}

type UpdateResult struct {
	commonResult
}

type GetResult struct {
	commonResult
}

type ListResult struct {
	commonResult
}

func (r ListResult) Extract() ([]AlarmMasking, error) {
	var s ListResponse
	err := r.ExtractInto(&s)
	if err != nil {
		return nil, err
	}
	return s.Items, nil
}

// DeleteResult represents the result of a delete operation. Call its ExtractErr
// method to determine if the request succeeded or failed.
type DeleteResult struct {
	golangsdk.ErrResult
}
