// Copyright © 2025 Ping Identity Corporation

// Package filter provides SCIM filter building utilities for the PingOne Terraform provider.
// This package contains functions for constructing SCIM (System for Cross-domain Identity Management)
// filter expressions used in PingOne API queries.
package filter

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// BuildScimFilter constructs a SCIM filter expression from a filter set and attribute mapping.
// It returns a properly formatted SCIM filter string that can be used in PingOne API queries.
// The filterSet parameter contains filter specifications with name and values fields.
// The attributeMapping parameter maps filter names to custom SCIM attribute expressions,
// allowing for complex filtering patterns beyond simple equality checks.
// Multiple values within a filter are combined with OR logic, while multiple filters are combined with AND logic.
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
