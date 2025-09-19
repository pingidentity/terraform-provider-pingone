// Copyright © 2025 Ping Identity Corporation

// Package utils provides utility functions for common operations in the PingOne Terraform provider.
package utils

import "github.com/hashicorp/terraform-plugin-framework/resource/schema"

// MergeSchemaAttributeMaps merges schema attributes from a source map into a destination map.
// The dest parameter is the target map that will receive the merged attributes.
// The src parameter is the source map containing attributes to merge.
// The overwrite parameter determines whether existing attributes in dest should be replaced by src attributes.
// When overwrite is false, existing attributes in dest are preserved and only new attributes from src are added.
// When overwrite is true, attributes from src will replace any existing attributes with the same key in dest.
// This function is useful for combining schema definitions from multiple sources in Terraform framework resources.
func MergeSchemaAttributeMaps(dest, src map[string]schema.Attribute, overwrite bool) {
	for key, value := range src {
		if _, exists := dest[key]; !exists || overwrite {
			dest[key] = value
		}
	}
}
