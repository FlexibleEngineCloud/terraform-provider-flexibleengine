package topics

import (
	"encoding/json"
	"strconv"

	"github.com/chnsz/golangsdk"
)

// CreateResponse is a struct that contains the create response
type CreateResponse struct {
	Name string `json:"id"`
}

// Topic includes the parameters of an topic
type Topic struct {
	Name             string `json:"id"`
	Partition        int    `json:"partition"`
	Replication      int    `json:"replication"`
	RetentionTime    int    `json:"retention_time"`
	TopicType        int    `json:"topic_type"`
	PoliciesOnly     bool   `json:"policiesOnly"`
	SyncReplication  bool   `json:"-"`
	SyncMessageFlush bool   `json:"-"`
}

// UnmarshalJSON to override default
func (r *Topic) UnmarshalJSON(b []byte) error {
	type tmp Topic
	var s struct {
		tmp
		SyncReplication  interface{} `json:"sync_replication"`
		SyncMessageFlush interface{} `json:"sync_message_flush"`
	}
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}

	*r = Topic(s.tmp)

	switch t := s.SyncReplication.(type) {
	case string:
		enabled, _ := strconv.ParseBool(s.SyncReplication.(string))
		r.SyncReplication = enabled
	case bool:
		r.SyncReplication = t
	}

	switch t := s.SyncMessageFlush.(type) {
	case string:
		enabled, _ := strconv.ParseBool(s.SyncMessageFlush.(string))
		r.SyncMessageFlush = enabled
	case bool:
		r.SyncMessageFlush = t
	}

	return err
}

// ListResponse is a struct that contains the list response
type ListResponse struct {
	Total            int     `json:"total"`
	Size             int     `json:"size"`
	RemainPartitions int     `json:"remain_partitions"`
	MaxPartitions    int     `json:"max_partitions"`
	Topics           []Topic `json:"topics"`
}

// CreateResult is a struct that contains all the return parameters of creation
type CreateResult struct {
	golangsdk.Result
}

// Extract from CreateResult
func (r CreateResult) Extract() (*CreateResponse, error) {
	var s CreateResponse
	err := r.Result.ExtractInto(&s)
	return &s, err
}

// ListResult contains the body of getting detailed
type ListResult struct {
	golangsdk.Result
}

// Extract from ListResult
func (r ListResult) Extract() ([]Topic, error) {
	var s ListResponse
	err := r.Result.ExtractInto(&s)
	return s.Topics, err
}

// DeleteResult is a struct which contains the result of deletion
type DeleteResult struct {
	golangsdk.Result
}

// DeleteResponse is a struct that contains the deletion response
type DeleteResponse struct {
	Name    string `json:"id"`
	Success bool   `json:"success"`
}

// Extract from DeleteResult
func (r DeleteResult) Extract() ([]DeleteResponse, error) {
	var s struct {
		Topics []DeleteResponse `json:"topics"`
	}
	err := r.Result.ExtractInto(&s)
	return s.Topics, err
}
