package flexibleengine

import (
	"fmt"
	"net"
	"regexp"
	"strings"
	"time"
)

func ValidateStringList(v interface{}, k string, l []string) (ws []string, errors []error) {
	value := v.(string)
	for i := range l {
		if value == l[i] {
			return
		}
	}
	errors = append(errors, fmt.Errorf("%q must be one of %v", k, l))
	return
}

// Not currently used.
func ValidateInsensitiveStringList(v interface{}, k string, l []string) (ws []string, errors []error) {
	value := v.(string)
	for i := range l {
		if strings.EqualFold(value, l[i]) {
			return
		}
	}
	errors = append(errors, fmt.Errorf("%q must be one of %v", k, l))
	return
}

func ValidateIntRange(v interface{}, k string, l int, h int) (ws []string, errors []error) {
	i, ok := v.(int)
	if !ok {
		errors = append(errors, fmt.Errorf("%q must be an integer", k))
		return
	}
	if i < l || i > h {
		errors = append(errors, fmt.Errorf("%q must be between %d and %d", k, l, h))
		return
	}
	return
}

func validateTrueOnly(v interface{}, k string) (ws []string, errors []error) {
	if b, ok := v.(bool); ok && b {
		return
	}
	if v, ok := v.(string); ok && v == "true" {
		return
	}
	errors = append(errors, fmt.Errorf("%q must be true", k))
	return
}

func validateS3BucketLifecycleTimestamp(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	_, err := time.Parse(time.RFC3339, fmt.Sprintf("%sT00:00:00Z", value))
	if err != nil {
		errors = append(errors, fmt.Errorf(
			"%q cannot be parsed as RFC3339 Timestamp Format", value))
	}

	return
}

func validateS3BucketLifecycleExpirationDays(v interface{}, k string) (ws []string, errors []error) {
	if v.(int) <= 0 {
		errors = append(errors, fmt.Errorf(
			"%q must be greater than 0", k))
	}

	return
}

func validateS3BucketLifecycleRuleId(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if len(value) > 255 {
		errors = append(errors, fmt.Errorf(
			"%q cannot exceed 255 characters", k))
	}
	return
}

func validateJsonString(v interface{}, k string) (ws []string, errors []error) {
	if _, err := normalizeJsonString(v); err != nil {
		errors = append(errors, fmt.Errorf("%q contains an invalid JSON: %s", k, err))
	}
	return
}

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

	if ipnet == nil || value != ipnet.String() {
		errors = append(errors, fmt.Errorf(
			"%q must contain a valid network CIDR, got %q", k, value))
	}

	return
}

func validateStackTemplate(v interface{}, k string) (ws []string, errors []error) {
	if looksLikeJsonString(v) {
		if _, err := normalizeJsonString(v); err != nil {
			errors = append(errors, fmt.Errorf("%q contains an invalid JSON: %s", k, err))
		}
	} else {
		if _, err := checkYamlString(v); err != nil {
			errors = append(errors, fmt.Errorf("%q contains an invalid YAML: %s", k, err))
		}
	}
	return
}

func validateVBSPolicyName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if strings.HasPrefix(strings.ToLower(value), "default") {
		errors = append(errors, fmt.Errorf(
			"%q cannot start with default: %q", k, value))
	}

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

func validateVBSPolicyFrequency(v interface{}, k string) (ws []string, errors []error) {
	value := v.(int)
	if value < 1 || value > 14 {
		errors = append(errors, fmt.Errorf(
			"%q should be in the range of 1-14: %d", k, value))
	}
	return
}

func validateVBSPolicyStatus(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if value != "ON" && value != "OFF" {
		errors = append(errors, fmt.Errorf(
			"%q should be either ON or OFF: %q", k, value))
	}
	return
}

func validateVBSPolicyRetentionNum(v interface{}, k string) (ws []string, errors []error) {
	value := v.(int)
	if value < 2 {
		errors = append(errors, fmt.Errorf(
			"%q cannot be less than 2: %d", k, value))
	}
	return
}

func validateVBSPolicyRetainBackup(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if value != "Y" && value != "N" {
		errors = append(errors, fmt.Errorf(
			"%q should be either N or Y: %q", k, value))
	}
	return
}

func validateVBSTagKey(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	if len(value) > 36 {
		errors = append(errors, fmt.Errorf(
			"%q cannot be longer than 36 characters: %q", k, value))
	}
	pattern := `^[\.\-_A-Za-z0-9]+$`
	if !regexp.MustCompile(pattern).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"%q doesn't comply with restrictions (%q): %q",
			k, pattern, value))
	}
	return
}

func validateVBSTagValue(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	if len(value) > 43 {
		errors = append(errors, fmt.Errorf(
			"%q cannot be longer than 43 characters: %q", k, value))
	}
	pattern := `^[\.\-_A-Za-z0-9]+$`
	if !regexp.MustCompile(pattern).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"%q doesn't comply with restrictions (%q): %q",
			k, pattern, value))
	}
	return
}

func validateVBSBackupName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if strings.HasPrefix(strings.ToLower(value), "autobk") {
		errors = append(errors, fmt.Errorf(
			"%q cannot start with autobk: %q", k, value))
	}

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

func validateVBSBackupDescription(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if len(value) > 64 {
		errors = append(errors, fmt.Errorf(
			"%q cannot be longer than 64 characters: %q", k, value))
	}
	pattern := `^[^<>]+$`
	if !regexp.MustCompile(pattern).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"%q doesn't comply with restrictions (%q): %q",
			k, pattern, value))
	}
	return
}

func validateAntiDdosTrafficPosID(v interface{}, k string) (ws []string, errors []error) {
	value := v.(int)
	if value < 1 || value > 9 {
		errors = append(errors, fmt.Errorf(
			"%q should be in the range of 1-9: %d", k, value))
	}
	return
}

func validateAntiDdosHttpRequestPosID(v interface{}, k string) (ws []string, errors []error) {
	value := v.(int)
	if value < 1 || value > 15 {
		errors = append(errors, fmt.Errorf(
			"%q should be in the range of 1-15: %d", k, value))
	}
	return
}

func validateAntiDdosCleaningAccessPosID(v interface{}, k string) (ws []string, errors []error) {
	value := v.(int)
	if value < 1 || value > 8 {
		errors = append(errors, fmt.Errorf(
			"%q should be in the range of 1-8: %d", k, value))
	}
	return
}

func validateAntiDdosAppTypeID(v interface{}, k string) (ws []string, errors []error) {
	value := v.(int)
	if value < 0 || value > 1 {
		errors = append(errors, fmt.Errorf(
			"%q should be 0 or 1: %d", k, value))
	}
	return
}

func validateECSTagValue(v interface{}, k string) (ws []string, errors []error) {
	tagmap := v.(map[string]interface{})
	vv := regexp.MustCompile(`^[0-9a-zA-Z-_]+$`)
	for k, v := range tagmap {
		value := v.(string)
		if !vv.MatchString(value) {
			err := fmt.Errorf("Tag value must be string only contains digits, "+
				"letters, underscores(_) and hyphens(-), but got %s=%s", k, value)
			errors = append(errors, err)
			break
		}
	}
	return
}
