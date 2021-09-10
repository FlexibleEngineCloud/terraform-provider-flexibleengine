package falsealarmmasking_rules

import (
	"github.com/chnsz/golangsdk"
)

var RequestOpts golangsdk.RequestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToAlarmMaskingCreateMap() (map[string]interface{}, error)
}

// CreateOpts contains all the values needed to create a new falsealarmmasking rule.
type CreateOpts struct {
	Path    string `json:"path" required:"true"`
	EventID string `json:"event_id" required:"true"`
}

// ToAlarmMaskingCreateMap builds a create request body from CreateOpts.
func (opts CreateOpts) ToAlarmMaskingCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Create will create a new falsealarmmasking rule based on the values in CreateOpts.
func Create(c *golangsdk.ServiceClient, policyID string, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToAlarmMaskingCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{200}}
	_, r.Err = c.Post(rootURL(c, policyID), b, &r.Body, reqOpt)
	return
}

// UpdateOptsBuilder allows extensions to add additional parameters to the
// Update request.
type UpdateOptsBuilder interface {
	ToAlarmMaskingUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts contains all the values needed to update a falsealarmmasking rule.
type UpdateOpts struct {
	Path    string `json:"path,omitempty"`
	EventID string `json:"event_id,omitempty"`
}

// ToAlarmMaskingUpdateMap builds a update request body from UpdateOpts.
func (opts UpdateOpts) ToAlarmMaskingUpdateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Update accepts a UpdateOpts struct and uses the values to update a rule.The response code from api is 200
func Update(c *golangsdk.ServiceClient, policyID, ruleID string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToAlarmMaskingUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{200}}
	_, r.Err = c.Put(resourceURL(c, policyID, ruleID), b, nil, reqOpt)
	return
}

// Get retrieves a particular falsealarmmasking rule based on its unique ID.
func Get(c *golangsdk.ServiceClient, policyID, ruleID string) (r GetResult) {
	reqOpt := &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	}

	_, r.Err = c.Get(resourceURL(c, policyID, ruleID), &r.Body, reqOpt)
	return
}

// List retrieves falsealarmmasking rules.
func List(c *golangsdk.ServiceClient, policyID string) (r ListResult) {
	_, r.Err = c.Get(rootURL(c, policyID), &r.Body, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	return
}

// Delete will permanently delete a particular falsealarmmasking rule based on its unique ID.
func Delete(c *golangsdk.ServiceClient, policyID, ruleID string) (r DeleteResult) {
	reqOpt := &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	}

	_, r.Err = c.Delete(resourceURL(c, policyID, ruleID), reqOpt)
	return
}
