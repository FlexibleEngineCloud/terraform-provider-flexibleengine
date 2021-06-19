package domains

import (
	"github.com/huaweicloud/golangsdk"
)

type Domain struct {
	// Domain ID
	Id string `json:"id"`
	// Domain name
	HostName string `json:"hostname"`
	// CNAME value
	Cname string `json:"cname"`
	// TXT record
	TxtCode string `json:"txt_code"`
	// Sub Domain name
	SubDomain string `json:"sub_domain"`
	// Policy ID
	PolicyID string `json:"policy_id"`
	// WAF mode: 0 - disabled, 1 - enabled, -1 - bypassed.
	ProtectStatus int `json:"protect_status"`
	// Whether a domain name is connected to WAF
	AccessStatus int `json:"access_status"`
	// Protocol type
	Protocol string `json:"protocol"`
	// Certificate ID
	CertificateId string `json:"certificate_id"`
	//The original server information
	Servers []Server `json:"server"`
	//Whether proxy is configured
	Proxy bool `json:"proxy"`
	//The type of the source IP header
	SipHeaderName string `json:"sip_header_name"`
	//The HTTP request header for identifying the real source IP.
	SipHeaderList []string `json:"sip_header_list"`
	// the time when the domain is created in unix timestamp
	TimeStamp int `json:"timestamp"`
}

type Server struct {
	// Protocol type of the client
	ClientProtocol string `json:"client_protocol"`
	// Protocol used by WAF to forward client requests to the server
	ServerProtocol string `json:"server_protocol"`
	// IP address or domain name of the web server that the client accesses.
	Address string `json:"address"`
	// Port number used by the web server
	Port int `json:"port"`
}

type commonResult struct {
	golangsdk.Result
}

// Extract is a function that accepts a result and extracts a domain.
func (r commonResult) Extract() (*Domain, error) {
	var response Domain
	err := r.ExtractInto(&response)
	return &response, err
}

// CreateResult represents the result of a create operation. Call its Extract
// method to interpret it as a Domain.
type CreateResult struct {
	commonResult
}

// UpdateResult represents the result of a update operation. Call its Extract
// method to interpret it as a Domain.
type UpdateResult struct {
	commonResult
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a Domain.
type GetResult struct {
	commonResult
}

// DeleteResult represents the result of a delete operation. Call its ExtractErr
// method to determine if the request succeeded or failed.
type DeleteResult struct {
	golangsdk.ErrResult
}
