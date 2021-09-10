package preciseprotection_rules

import (
	"github.com/chnsz/golangsdk"
)

var RequestOpts golangsdk.RequestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToPreciseCreateMap() (map[string]interface{}, error)
}

// CreateOpts contains all the values needed to create a new precise protection rule.
type CreateOpts struct {
	Name       string      `json:"name" required:"true"`
	Time       bool        `json:"time"`
	Start      int64       `json:"start,omitempty"`
	End        int64       `json:"end,omitempty"`
	Conditions []Condition `json:"conditions" required:"true"`
	Action     Action      `json:"action" required:"true"`
	Priority   *int        `json:"priority,omitempty"`
}

type Condition struct {
	Category string   `json:"category" required:"true"`
	Index    string   `json:"index,omitempty"`
	Logic    string   `json:"logic" required:"true"`
	Contents []string `json:"contents" required:"true"`
}

type Action struct {
	Category string `json:"category" required:"true"`
}

// ToPreciseCreateMap builds a create request body from CreateOpts.
func (opts CreateOpts) ToPreciseCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Create will create a new precise protection rule based on the values in CreateOpts.
func Create(c *golangsdk.ServiceClient, policyID string, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToPreciseCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{200}}
	_, r.Err = c.Post(rootURL(c, policyID), b, &r.Body, reqOpt)
	return
}

// Update will update a precise protection rule based on the values in CreateOpts.
// The response code from api is 200
func Update(c *golangsdk.ServiceClient, policyID, ruleID string, opts CreateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToPreciseCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{200}}
	_, r.Err = c.Put(resourceURL(c, policyID, ruleID), b, &r.Body, reqOpt)
	return
}

// Get retrieves a particular precise rule based on its unique ID.
func Get(c *golangsdk.ServiceClient, policyID, ruleID string) (r GetResult) {
	reqOpt := &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	}

	_, r.Err = c.Get(resourceURL(c, policyID, ruleID), &r.Body, reqOpt)
	return
}

// Delete will permanently delete a particular precise rule based on its unique ID.
func Delete(c *golangsdk.ServiceClient, policyID, ruleID string) (r DeleteResult) {
	reqOpt := &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	}

	_, r.Err = c.Delete(resourceURL(c, policyID, ruleID), reqOpt)
	return
}
