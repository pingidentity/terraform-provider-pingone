package utils

import "github.com/hashicorp/terraform-plugin-framework/resource/schema"

func MergeSchemaAttributeMaps(dest, src map[string]schema.Attribute, overwrite bool) {
	for key, value := range src {
		if _, exists := dest[key]; !exists || overwrite {
			dest[key] = value
		}
	}
}
