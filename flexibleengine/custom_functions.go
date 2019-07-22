package flexibleengine

import (
	"strings"
)

func expandRdsInstanceV3CreateRegion(d interface{}, arrayIndex map[string]int) (interface{}, error) {
	return navigateValue(d, []string{"region"}, arrayIndex)
}

func flattenRdsInstanceV3HAReplicationMode(d interface{}, arrayIndex map[string]int, currentValue interface{}) (interface{}, error) {
	v, err := navigateValue(d, []string{"list", "flavor_ref"}, nil)
	if err != nil {
		return nil, err
	}
	if strings.HasSuffix(v.(string), ".ha") {
		return navigateValue(d, []string{"list", "ha", "replication_mode"}, nil)
	}
	return "", nil
}
