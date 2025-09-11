// Copyright © 2025 Ping Identity Corporation

package helpers

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/patrickcping/pingone-go-sdk-v2/risk"
	"github.com/pingidentity/terraform-provider-pingone/internal/utils"
)

// CheckCompositeConditionStructure validates the structure of a composite condition JSON string.
// It returns diagnostics indicating whether the JSON structure is valid and supported.
// The ctx parameter provides context for logging validation operations.
// The jsonValue parameter contains the JSON string to validate as a composite condition.
// This function unmarshals the JSON through SDK objects to verify compatibility and detect unsupported fields.
func CheckCompositeConditionStructure(ctx context.Context, jsonValue string) diag.Diagnostics {
	var diags diag.Diagnostics

	// Unmarshal and marshal through the SDK objects to strip out unsupported fields.
	var condition risk.RiskPredictorCompositeConditionBase
	err := json.Unmarshal([]byte(jsonValue), &condition)
	if err != nil {
		tflog.Error(ctx, "Cannot parse the condition input JSON", map[string]interface{}{
			"err": err,
		})
		diags.AddError(
			"Cannot parse the condition input JSON",
			"The JSON string passed to the condition parameter cannot be parsed as a condition object.  Please check the policy is a valid structure.",
		)
		return diags
	}

	jsonBytes, err := json.Marshal(condition)
	if err != nil {
		tflog.Error(ctx, "Failed to marshal condition object to bytes", map[string]interface{}{
			"err": err,
		})
		diags.AddError(
			"Failed to marshal condition object to bytes",
			"The condition object cannot be converted back to string.  Please report this to the provider maintainers.",
		)

		return diags
	}

	// Check equality of the JSON to see if anything stripped out.  This indicates an unsupported field value.
	if !utils.DeepEqualJSON([]byte(jsonValue), jsonBytes) {
		tflog.Warn(ctx, "Condition object has unsupported fields", map[string]interface{}{
			"err": err,
		})
		diags.AddWarning(
			"Composite condition import has unsupported fields",
			"The composite condition import contains unsupported fields.  Unpredictable results may occur.",
		)
		return diags
	}

	return diags
}

// NormaliseCompositeCondition normalizes a composite condition JSON string by adding missing type fields.
// It returns a pointer to the normalized JSON string and any diagnostics encountered during processing.
// The ctx parameter provides context for logging normalization operations.
// The jsonValue parameter contains the JSON string to normalize.
// This function walks the JSON tree and adds "type" fields to "and", "or", and "not" objects where missing.
func NormaliseCompositeCondition(ctx context.Context, jsonValue string) (*string, diag.Diagnostics) {
	var diags diag.Diagnostics

	var condition map[string]interface{}
	err := json.Unmarshal([]byte(jsonValue), &condition)
	if err != nil {
		tflog.Error(ctx, "Cannot parse the condition input JSON", map[string]interface{}{
			"err": err,
		})
		diags.AddError(
			"Cannot parse the condition input JSON",
			"The JSON string passed to the condition parameter cannot be parsed as a map.  Please check the policy is a valid structure.",
		)
		return nil, diags
	}

	// Walk the JSON tree and add "type" to "and", "or" and "not" objects if not set.
	tflog.Debug(ctx, "Normalising condition object", map[string]interface{}{})

	condition = WalkAggregatedCondition(condition)

	newJsonBytes, err := json.Marshal(condition)
	if err != nil {
		tflog.Warn(ctx, "Failed to marshal new condition map to bytes", map[string]interface{}{
			"err": err,
		})
		diags.AddError(
			"Failed to marshal new condition map to bytes",
			"The condition map cannot be converted back to string.  Please report this to the provider maintainers.",
		)

		return nil, diags
	}

	returnVar := string(newJsonBytes)

	tflog.Debug(ctx, "Normalised condition object", map[string]interface{}{
		"jsonValue::normalised": returnVar,
		"jsonValue::input":      jsonValue,
	})

	return &returnVar, diags
}

// WalkAggregatedCondition recursively walks through a condition map and adds missing type fields to aggregated conditions.
// It returns the modified condition map with proper type fields added to "and", "or", and "not" objects.
// The condition parameter is a map[string]interface{} representing a single condition object in the JSON structure.
// This function identifies aggregated condition types and ensures they have the correct "type" field values.
func WalkAggregatedCondition(condition map[string]interface{}) map[string]interface{} {

	evaluateAggregateConditions := []struct {
		mapIndex  string
		typeValue string
		sliceType bool
	}{
		{
			mapIndex:  "and",
			typeValue: "AND",
			sliceType: true,
		},
		{
			mapIndex:  "or",
			typeValue: "OR",
			sliceType: true,
		},
		{
			mapIndex:  "not",
			typeValue: "NOT",
			sliceType: false,
		},
	}

	for _, evaluateCondition := range evaluateAggregateConditions {
		if condition[evaluateCondition.mapIndex] != nil {
			condition["type"] = evaluateCondition.typeValue

			if evaluateCondition.sliceType {
				condition[evaluateCondition.mapIndex] = WalkAggregatedListCondition(condition[evaluateCondition.mapIndex].([]interface{}))
			} else {
				condition[evaluateCondition.mapIndex] = WalkAggregatedCondition(condition[evaluateCondition.mapIndex].(map[string]interface{}))
			}
		}
	}

	return condition
}

// WalkAggregatedListCondition recursively processes a list of condition objects by applying normalization to each.
// It returns a slice of interface{} containing the normalized condition objects.
// The conditionList parameter is a slice of interface{} representing multiple condition objects in a list.
// This function is used to process the condition arrays found in "and" and "or" aggregated conditions.
func WalkAggregatedListCondition(conditionList []interface{}) []interface{} {

	conditionReturnList := make([]interface{}, 0)

	for _, c := range conditionList {
		conditionReturnList = append(conditionReturnList, WalkAggregatedCondition(c.(map[string]interface{})))
	}

	return conditionReturnList
}
