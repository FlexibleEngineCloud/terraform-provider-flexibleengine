package deprecated

import (
	"fmt"
	"net"
	"reflect"
	"regexp"
	"strings"
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
