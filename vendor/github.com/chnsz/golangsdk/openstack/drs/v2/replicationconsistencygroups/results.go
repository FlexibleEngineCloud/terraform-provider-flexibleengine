package replicationconsistencygroups

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

// ReplicationConsistencyGroupCreateorUpdate response
type ReplicationConsistencyGroupCreateorUpdate struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	Status          string `json:"status"`
	PriorityStation string `json:"priority_station"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
}

// CreateResult is a struct that contains all the return parameters of creation
type CreateResult struct {
	golangsdk.Result
}

// Extract from CreateResult
func (r CreateResult) Extract() (*ReplicationConsistencyGroupCreateorUpdate, error) {
	var s struct {
		ReplicationConsistencyGroup *ReplicationConsistencyGroupCreateorUpdate `json:"replication_consistency_group"`
	}
	err := r.Result.ExtractInto(&s)
	return s.ReplicationConsistencyGroup, err
}

// DeleteResult is a struct which contains the result of deletion
type DeleteResult struct {
	golangsdk.ErrResult
}

// ReplicationConsistencyGroup response
type ReplicationConsistencyGroup struct {
	ID                string   `json:"id"`
	Name              string   `json:"name"`
	Description       string   `json:"description"`
	Status            string   `json:"status"`
	PriorityStation   string   `json:"priority_station"`
	ReplicationModel  string   `json:"replication_model"`
	ReplicationStatus string   `json:"replication_status"`
	ReplicationIDs    []string `json:"replication_ids"`
	CreatedAt         string   `json:"created_at"`
	UpdatedAt         string   `json:"updated_at"`
	FailureDetail     string   `json:"failure_detail"`
	FaultLevel        string   `json:"fault_level"`
}

// GetResult contains the body of getting detailed
type GetResult struct {
	golangsdk.Result
}

// Extract from GetResult
func (r GetResult) Extract() (*ReplicationConsistencyGroup, error) {
	var s struct {
		ReplicationConsistencyGroup *ReplicationConsistencyGroup `json:"replication_consistency_group"`
	}
	err := r.Result.ExtractInto(&s)
	return s.ReplicationConsistencyGroup, err
}

// ReplicationConsistencyGroupPage may be embedded in a Page
// that contains all of the results from an operation at once.
type ReplicationConsistencyGroupPage struct {
	pagination.SinglePageBase
}

// IsEmpty returns true if a ListResult contains no replications.
func (r ReplicationConsistencyGroupPage) IsEmpty() (bool, error) {
	rs, err := ExtractReplicationConsistencyGroups(r)
	return len(rs) == 0, err
}

// ExtractReplicationConsistencyGroups from List
func ExtractReplicationConsistencyGroups(r pagination.Page) ([]ReplicationConsistencyGroup, error) {
	var s struct {
		ReplicationConsistencyGroups []ReplicationConsistencyGroup `json:"replication_consistency_groups"`
	}
	err := (r.(ReplicationConsistencyGroupPage)).ExtractInto(&s)
	return s.ReplicationConsistencyGroups, err
}

// UpdateResult is a struct from which can get the result of update method
type UpdateResult struct {
	golangsdk.Result
}

// Extract from UpdateResult
func (r UpdateResult) Extract() (*ReplicationConsistencyGroupCreateorUpdate, error) {
	var s struct {
		ReplicationConsistencyGroup *ReplicationConsistencyGroupCreateorUpdate `json:"replication_consistency_group"`
	}
	err := r.Result.ExtractInto(&s)
	return s.ReplicationConsistencyGroup, err
}

// ActionResult is the result of action operations
type ActionResult struct {
	golangsdk.ErrResult
}
