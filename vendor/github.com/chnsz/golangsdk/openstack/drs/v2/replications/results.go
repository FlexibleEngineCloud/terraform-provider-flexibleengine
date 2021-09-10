package replications

import (
	"encoding/json"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

// ReplicationCreate response
type ReplicationCreate struct {
	ID                            string `json:"id"`
	Name                          string `json:"name"`
	Description                   string `json:"description"`
	Status                        string `json:"status"`
	ReplicationConsistencyGroupID string `json:"replication_consistency_group_id"`
	VolumeIDs                     string `json:"volume_ids"`
	PriorityStation               string `json:"priority_station"`
	CreatedAt                     string `json:"created_at"`
	UpdatedAt                     string `json:"updated_at"`
}

// CreateResult is a struct that contains all the return parameters of creation
type CreateResult struct {
	golangsdk.Result
}

// Extract from CreateResult
func (r CreateResult) Extract() (*ReplicationCreate, error) {
	var s struct {
		Replication *ReplicationCreate `json:"replication"`
	}
	err := r.Result.ExtractInto(&s)
	return s.Replication, err
}

// DeleteResult is a struct which contains the result of deletion
type DeleteResult struct {
	golangsdk.ErrResult
}

// Replication response
type Replication struct {
	ID                            string `json:"id"`
	Name                          string `json:"name"`
	Description                   string `json:"description"`
	Status                        string `json:"status"`
	ReplicationConsistencyGroupID string `json:"replication_consistency_group_id"`
	VolumeIDs                     string `json:"volume_ids"`
	PriorityStation               string `json:"priority_station"`
	CreatedAt                     string `json:"created_at"`
	UpdatedAt                     string `json:"updated_at"`
	ReplicationModel              string `json:"replication_model"`
	ReplicationStatus             string `json:"replication_status"`
	Progress                      string `json:"progress"`
	FailureDetail                 string `json:"failure_detail"`
	// RecordMetadata includes volume_type and multiattach currently.
	RecordMetadata json.RawMessage `json:"record_metadata"`
	FaultLevel     string          `json:"fault_level"`
}

// GetResult contains the body of getting detailed
type GetResult struct {
	golangsdk.Result
}

// Extract from GetResult
func (r GetResult) Extract() (*Replication, error) {
	var s struct {
		Replication *Replication `json:"replication"`
	}
	err := r.Result.ExtractInto(&s)
	return s.Replication, err
}

// ReplicationPage may be embedded in a Page
// that contains all of the results from an operation at once.
type ReplicationPage struct {
	pagination.SinglePageBase
}

// IsEmpty returns true if a ListResult contains no replications.
func (r ReplicationPage) IsEmpty() (bool, error) {
	rs, err := ExtractReplications(r)
	return len(rs) == 0, err
}

// ExtractReplications from List
func ExtractReplications(r pagination.Page) ([]Replication, error) {
	var s struct {
		Replications []Replication `json:"replications"`
	}
	err := (r.(ReplicationPage)).ExtractInto(&s)
	return s.Replications, err
}
