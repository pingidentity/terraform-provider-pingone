// Copyright © 2026 Ping Identity Corporation

package utils

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	dataSourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func MergeSchemaAttributeMaps(dest, src map[string]resourceSchema.Attribute, overwrite bool) {
	for key, value := range src {
		if _, exists := dest[key]; !exists || overwrite {
			dest[key] = value
		}
	}
}

func MergeResourceSchemaAttributeMapsRtn(src ...map[string]resourceSchema.Attribute) map[string]resourceSchema.Attribute {
	result := make(map[string]resourceSchema.Attribute)
	for _, m := range src {
		for key, value := range m {
			if _, exists := result[key]; !exists {
				result[key] = value
			}
		}
	}
	return result
}

func MergeDataSourceSchemaAttributeMapsRtn(src ...map[string]dataSourceSchema.Attribute) map[string]dataSourceSchema.Attribute {
	result := make(map[string]dataSourceSchema.Attribute)
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
