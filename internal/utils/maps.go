// Copyright © 2025 Ping Identity Corporation

package utils

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func MergeSchemaAttributeMaps(dest, src map[string]schema.Attribute, overwrite bool) {
	for key, value := range src {
		if _, exists := dest[key]; !exists || overwrite {
			dest[key] = value
		}
	}
}

func MergeSchemaAttributeMapsRtn(src ...map[string]schema.Attribute) map[string]schema.Attribute {
	result := make(map[string]schema.Attribute)
	for _, m := range src {
		for key, value := range m {
			if _, exists := result[key]; !exists {
				result[key] = value
			}
		}
	}
	return result
}

func MergeAttributeTypeMapsRtn(src ...map[string]attr.Type) map[string]attr.Type {
	result := make(map[string]attr.Type)
	for _, m := range src {
		for key, value := range m {
			if _, exists := result[key]; !exists {
				result[key] = value
			}
		}
	}
	return result
}

func MergeAttributeValueMapsRtn(src ...map[string]attr.Value) map[string]attr.Value {
	result := make(map[string]attr.Value)
	for _, m := range src {
		for key, value := range m {
			if _, exists := result[key]; !exists {
				result[key] = value
			}
		}
	}
	return result
}
