package deprecated

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"reflect"
	"regexp"
	"strings"

	"github.com/chnsz/golangsdk"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func validateName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if len(value) > 64 {
		errors = append(errors, fmt.Errorf(
			"%q cannot be longer than 64 characters: %q", k, value))
	}

	pattern := `^[\.\-_A-Za-z0-9]+$`
	if !regexp.MustCompile(pattern).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"%q doesn't comply with restrictions (%q): %q",
			k, pattern, value))
	}

	return
}

func validateIP(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	ipnet := net.ParseIP(value)

	if ipnet == nil || value != ipnet.String() {
		errors = append(errors, fmt.Errorf(
			"%q must contain a valid network IP address, got %q", k, value))
	}

	return
}

func validateCIDR(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	_, ipnet, err := net.ParseCIDR(value)
	if err != nil {
		errors = append(errors, fmt.Errorf(
			"%q must contain a valid CIDR, got error parsing: %s", k, err))
		return
	}

	if ipnet == nil || strings.ToLower(value) != ipnet.String() {
		errors = append(errors, fmt.Errorf(
			"%q must contain a valid network CIDR, got %q", k, value))
	}

	return
}

func navigateValue(d interface{}, index []string, arrayIndex map[string]int) (interface{}, error) {
	for n, i := range index {
		if d == nil {
			return nil, nil
		}
		if d1, ok := d.(map[string]interface{}); ok {
			d, ok = d1[i]
			if !ok {
				msg := fmt.Sprintf("navigate value with index(%s)", strings.Join(index, "."))
				return nil, fmt.Errorf("%s: '%s' may not exist", msg, i)
			}
		} else {
			msg := fmt.Sprintf("navigate value with index(%s)", strings.Join(index, "."))
			return nil, fmt.Errorf("%s: Can not convert (%s) to map", msg, reflect.TypeOf(d))
		}

		if arrayIndex != nil {
			if j, ok := arrayIndex[strings.Join(index[:n+1], ".")]; ok {
				if d == nil {
					return nil, nil
				}
				if d2, ok := d.([]interface{}); ok {
					if len(d2) == 0 {
						return nil, nil
					}
					if j >= len(d2) {
						msg := fmt.Sprintf("navigate value with index(%s)", strings.Join(index, "."))
						return nil, fmt.Errorf("%s: The index is out of array", msg)
					}

					d = d2[j]
				} else {
					msg := fmt.Sprintf("navigate value with index(%s)", strings.Join(index, "."))
					return nil, fmt.Errorf("%s: Can not convert (%s) to array, index=%s.%v", msg, reflect.TypeOf(d), i, j)
				}
			}
		}
	}

	return d, nil
}

// ErrorResp is the response when API failed
type ErrorResp struct {
	ErrorCode string `json:"error_code"`
	ErrorMsg  string `json:"error_msg"`
}

// ParseErrorMsg is used to unmarshal the error body to ErrorResp
// usage: resp, pErr := ParseErrorMsg(err.Body)
func ParseErrorMsg(body []byte) (ErrorResp, error) {
	resp := ErrorResp{}
	err := json.Unmarshal(body, &resp)
	return resp, err
}

// CheckDeleted checks the error to see if it's a 404 (Not Found) and, if so,
// sets the resource ID to the empty string instead of throwing an error.
func CheckDeleted(d *schema.ResourceData, err error, msg string) error {
	if _, ok := err.(golangsdk.ErrDefault404); ok {
		d.SetId("")
		return nil
	}

	return fmt.Errorf("%s: %s", msg, err)
}

// CheckDeletedDiag checks the error to see if it's a 404 (Not Found) and, if so,
// sets the resource ID to the empty string instead of throwing an error.
func CheckDeletedDiag(d *schema.ResourceData, err error, msg string) diag.Diagnostics {
	if _, ok := err.(golangsdk.ErrDefault404); ok {
		d.SetId("")
		return nil
	}

	return diag.Errorf("%s: %s", msg, err)
}

// GetRegion returns the region that was specified in the resource. If a
// region was not set, the provider-level region is checked. The provider-level
// region can either be set by the region argument or by OS_REGION_NAME.
func GetRegion(d *schema.ResourceData, config *config.Config) string {
	if v, ok := d.GetOk("region"); ok {
		return v.(string)
	}

	return config.Region
}

func checkForRetryableError(err error) *resource.RetryError {
	switch errCode := err.(type) {
	case golangsdk.ErrDefault500:
		return resource.RetryableError(err)
	case golangsdk.ErrUnexpectedResponseCode:
		switch errCode.Actual {
		case 409, 503:
			return resource.RetryableError(err)
		default:
			return resource.NonRetryableError(err)
		}
	default:
		return resource.NonRetryableError(err)
	}
}

func isResourceNotFound(err error) bool {
	if err == nil {
		return false
	}
	_, ok := err.(golangsdk.ErrDefault404)
	return ok
}

func hasFilledOpt(d *schema.ResourceData, param string) bool {
	_, b := d.GetOkExists(param)
	return b
}

// strSliceContains checks if a given string is contained in a slice
// When anybody asks why Go needs generics, here you go.
func strSliceContains(haystack []string, needle string) bool {
	for _, s := range haystack {
		if s == needle {
			return true
		}
	}
	return false
}

func jsonBytesEqual(b1, b2 []byte) bool {
	var o1 interface{}
	if err := json.Unmarshal(b1, &o1); err != nil {
		return false
	}

	var o2 interface{}
	if err := json.Unmarshal(b2, &o2); err != nil {
		return false
	}

	return reflect.DeepEqual(o1, o2)
}

// convertStructToMap converts an instance of struct to a map object, and
// changes each key of fileds to the value of 'nameMap' if the key in it
// or to its corresponding lowercase.
func convertStructToMap(obj interface{}, nameMap map[string]string) (map[string]interface{}, error) {
	b, err := json.Marshal(obj)
	if err != nil {
		return nil, fmt.Errorf("Error converting struct to map, marshal failed:%v", err)
	}

	m, err := regexp.Compile(`"[a-z0-9A-Z_]+":`)
	if err != nil {
		return nil, fmt.Errorf("Error converting struct to map, compile regular express failed")
	}
	nb := m.ReplaceAllFunc(
		b,
		func(src []byte) []byte {
			k := fmt.Sprintf("%s", src[1:len(src)-2])
			v, ok := nameMap[k]
			if !ok {
				v = strings.ToLower(k)
			}
			return []byte(fmt.Sprintf("\"%s\":", v))
		},
	)
	log.Printf("[DEBUG]convertStructToMap:: before change b =%s", b)
	log.Printf("[DEBUG]convertStructToMap:: after change nb=%s", nb)

	p := make(map[string]interface{})
	err = json.Unmarshal(nb, &p)
	if err != nil {
		return nil, fmt.Errorf("Error converting struct to map, unmarshal failed:%v", err)
	}
	log.Printf("[DEBUG]convertStructToMap:: map= %#v\n", p)
	return p, nil
}

func looksLikeJsonString(s interface{}) bool {
	return regexp.MustCompile(`^\s*{`).MatchString(s.(string))
}

func compareJsonTemplateAreEquivalent(tem1, tem2 string) (bool, error) {
	var obj1 interface{}
	err := json.Unmarshal([]byte(tem1), &obj1)
	if err != nil {
		return false, err
	}

	canonicalJson1, _ := json.Marshal(obj1)

	var obj2 interface{}
	err = json.Unmarshal([]byte(tem2), &obj2)
	if err != nil {
		return false, err
	}

	canonicalJson2, _ := json.Marshal(obj2)

	equal := bytes.Compare(canonicalJson1, canonicalJson2) == 0
	if !equal {
		log.Printf("[DEBUG] Canonical template are not equal.\nFirst: %s\nSecond: %s\n",
			canonicalJson1, canonicalJson2)
	}
	return equal, nil
}

// HashStrings hashes a list of strings to a unique hashcode.
func HashStrings(strings []string) string {
	var buf bytes.Buffer

	for _, s := range strings {
		buf.WriteString(fmt.Sprintf("%s-", s))
	}

	return fmt.Sprintf("%d", schema.HashString(buf.String()))
}

func IsUUIDFormat(str string) bool {
	if _, err := uuid.ParseUUID(str); err != nil {
		log.Printf("[WARN] '%s' is not a valid UUID: %s", str, err)
		return false
	}
	return true
}

// isStrContainsSliceElement returns true if the string exists in given slice or contains in one of slice elements when
// open exact flag. Also you can ignore case for this check.
func isStrContainsSliceElement(str string, sl []string, ignoreCase, isExcat bool) bool {
	if ignoreCase {
		str = strings.ToLower(str)
	}
	for _, s := range sl {
		if ignoreCase {
			s = strings.ToLower(s)
		}
		if isExcat && s == str {
			return true
		}
		if !isExcat && strings.Contains(str, s) {
			return true
		}
	}
	return false
}
