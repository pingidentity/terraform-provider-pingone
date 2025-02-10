// Copyright © 2025 Ping Identity Corporation

package filter

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func BuildScimFilter(filterSet []interface{}, attributeMapping map[string]string) string {

	generalMapping := "%s eq \"%s\""

	returnFilterList := make([]string, len(filterSet))

	for i, v := range filterSet {

		valueObj := v.(map[string]interface{})

		var valueFilter []string

		switch v := valueObj["values"].(type) {
		case []string:
			valueFilter = v
		case []*string:
			for i := range v {
				valueFilter = append(valueFilter, *v[i])
			}
		case *schema.Set:
			valueFilter = make([]string, len(v.List()))
			for i, c := range v.List() {
				valueFilter[i] = c.(string)
			}
		}

		valueFilterList := make([]string, len(valueFilter))
		for j, value := range valueFilter {

			if val, ok := attributeMapping[valueObj["name"].(string)]; ok {
				valueFilterList[j] = fmt.Sprintf(fmt.Sprintf("(%s)", val), value)
			} else {
				valueFilterList[j] = fmt.Sprintf(fmt.Sprintf("(%s)", generalMapping), valueObj["name"].(string), value)
			}
		}

		returnFilterList[i] = fmt.Sprintf("(%s)", strings.Join(valueFilterList, " OR "))

	}

	return strings.Join(returnFilterList, " AND ")

}
