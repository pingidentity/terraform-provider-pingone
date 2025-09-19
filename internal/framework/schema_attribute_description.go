// Copyright © 2025 Ping Identity Corporation

// Package framework provides common utilities and helpers for Terraform Plugin Framework implementations.
// This package contains schema description builders, validators, and other utilities that support
// the Framework-based provider implementation.
package framework

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/pingidentity/terraform-provider-pingone/internal/utils"
)

// SchemaDescriptionModel represents the basic structure for schema descriptions.
// It contains both plain text and markdown-formatted versions of descriptions.
type SchemaDescriptionModel struct {
	// Description is the plain text version of the schema description
	Description string
	// MarkdownDescription is the markdown-formatted version of the schema description
	MarkdownDescription string
}

// SchemaAttributeDescription provides methods for building rich schema attribute descriptions.
// It extends SchemaDescriptionModel with chainable methods for adding contextual information
// such as default values, allowed values, conflicts, and validation rules.
type SchemaAttributeDescription SchemaDescriptionModel

// Clean normalizes and formats the schema description by removing trailing punctuation and whitespace.
// It returns a cleaned SchemaAttributeDescription with standardized formatting.
// The removeTrailingStop parameter determines whether trailing periods should be removed.
// This method also ensures that both Description and MarkdownDescription fields are populated
// by copying content from one to the other when one is empty.
func (r SchemaAttributeDescription) Clean(removeTrailingStop bool) SchemaAttributeDescription {

	// Trim trailing fullstop
	if removeTrailingStop {
		trailingDot := regexp.MustCompile(`(\.\s*)$`)
		r.Description = trailingDot.ReplaceAllString(r.Description, "")
		r.MarkdownDescription = trailingDot.ReplaceAllString(r.MarkdownDescription, "")
	}

	r.Description = strings.TrimSpace(r.Description)
	r.MarkdownDescription = strings.TrimSpace(r.MarkdownDescription)

	if r.MarkdownDescription == "" && r.Description != "" {
		// Prefil the blank markdown description with the description value
		r.MarkdownDescription = r.Description
	}

	if r.MarkdownDescription != "" && r.Description == "" {
		// Prefil the blank description with the markdown description value, ignoring MD formatting
		r.Description = r.MarkdownDescription
	}

	return r
}

// DefaultValue appends default value information to the schema description.
// It returns a SchemaAttributeDescription with the default value documented.
// The defaultValue parameter can be of various types (string, int, bool, etc.)
// and will be automatically converted to an appropriate string representation.
func (r SchemaAttributeDescription) DefaultValue(defaultValue any) SchemaAttributeDescription {
	return r.genericValue("Defaults to", defaultValue)
}

// FixedValue appends fixed value information to the schema description.
// It returns a SchemaAttributeDescription with the fixed value documented.
// This is used for attributes that have a constant, unchangeable value.
// The defaultValue parameter will be formatted and displayed as a fixed constraint.
func (r SchemaAttributeDescription) FixedValue(defaultValue any) SchemaAttributeDescription {
	return r.genericValue("Fixed value of", defaultValue)
}

// genericValue is a helper method that formats and appends value information to the schema description.
// It returns a SchemaAttributeDescription with the value information documented.
// The text parameter provides the prefix text (e.g., "Defaults to", "Fixed value of").
// The defaultValue parameter is converted to a string representation and appended to the description.
func (r SchemaAttributeDescription) genericValue(text string, defaultValue any) SchemaAttributeDescription {
	var defaultValueString string
	switch v := defaultValue.(type) {
	case string:
		defaultValueString = v
	case int:
		defaultValueString = strconv.Itoa(v)
	case int32:
		defaultValueString = strconv.Itoa(int(v))
	case int64:
		defaultValueString = strconv.FormatInt(v, 10)
	case bool:
		defaultValueString = strconv.FormatBool(v)
	default:
		defaultValueString = "DOC ERROR: Unknown default data type"
	}
	return r.AppendStringValue(text, defaultValueString)
}

// AllowedValues appends information about allowed values to the schema description.
// It returns a SchemaAttributeDescription with the allowed values documented and sorted alphabetically.
// The allowedValues parameter accepts multiple values of various types that will be converted to strings.
// This method is useful for documenting enumerated values or constrained input options.
func (r SchemaAttributeDescription) AllowedValues(allowedValues ...any) SchemaAttributeDescription {

	allowedValuesParsed := make([]string, 0)
	for _, allowedValue := range allowedValues {
		switch v := allowedValue.(type) {
		case string:
			allowedValuesParsed = append(allowedValuesParsed, v)
		case int:
			allowedValuesParsed = append(allowedValuesParsed, strconv.Itoa(v))
		case int32:
			allowedValuesParsed = append(allowedValuesParsed, strconv.Itoa(int(v)))
		case int64:
			allowedValuesParsed = append(allowedValuesParsed, strconv.FormatInt(v, 10))
		default:
			allowedValuesParsed = append(allowedValuesParsed, fmt.Sprintf("DOC ERROR: Unknown allowed value type: %s", v))
		}
	}

	sort.Strings(allowedValuesParsed)

	return r.AppendSliceValues("Options are", allowedValuesParsed)
}

// AllowedValuesComplex appends information about allowed values with descriptions to the schema description.
// It returns a SchemaAttributeDescription with the allowed values and their explanations documented.
// The allowedValuesMap parameter maps value strings to their descriptive explanations.
// This method is useful for documenting complex enumerated values where each option needs additional context.
func (r SchemaAttributeDescription) AllowedValuesComplex(allowedValuesMap map[string]string) SchemaAttributeDescription {
	allowedValues := make([]string, 0)
	for k, v := range allowedValuesMap {
		allowedValues = append(allowedValues, fmt.Sprintf("`%s` (%s)", k, v))
	}

	sort.Strings(allowedValues)

	return r.AppendMarkdownString(fmt.Sprintf("Options are %s.", strings.Join(allowedValues, ", ")))
}

// AllowedValuesEnum appends information about enum-based allowed values to the schema description.
// It returns a SchemaAttributeDescription with the enum values documented.
// The allowedValuesEnumSlice parameter should be a slice of enum values that will be converted to strings.
// This method is specifically designed for PingOne SDK enum types and uses utility functions for conversion.
func (r SchemaAttributeDescription) AllowedValuesEnum(allowedValuesEnumSlice interface{}) SchemaAttributeDescription {
	return r.AllowedValues(utils.EnumSliceToAnySlice(allowedValuesEnumSlice)...)
}

// ConflictsWith appends information about conflicting attributes to the schema description.
// It returns a SchemaAttributeDescription with the conflicting field paths documented.
// The fieldPaths parameter lists the attribute paths that cannot be used together with this attribute.
// This method helps document mutual exclusion rules between different configuration options.
func (r SchemaAttributeDescription) ConflictsWith(fieldPaths []string) SchemaAttributeDescription {
	return r.AppendSliceValues("Conflicts with", fieldPaths)
}

// ExactlyOneOf appends information about exactly-one-of validation rules to the schema description.
// It returns a SchemaAttributeDescription with the exactly-one-of constraint documented.
// The fieldPaths parameter lists the attribute paths where exactly one must be defined.
// This method helps document validation rules for mutually exclusive but required attribute groups.
func (r SchemaAttributeDescription) ExactlyOneOf(fieldPaths []string) SchemaAttributeDescription {
	return r.AppendSliceValues("Exactly one of the following must be defined:", fieldPaths)
}

// RequiresReplace appends information about immutable attributes that trigger resource replacement.
// It returns a SchemaAttributeDescription with the replacement behavior documented.
// This method is used for attributes that cannot be modified in-place and require
// destroying and recreating the resource when changed.
func (r SchemaAttributeDescription) RequiresReplace() SchemaAttributeDescription {
	return r.AppendMarkdownString("This field is immutable and will trigger a replace plan if changed.")
}

// RequiresReplaceNestedAttributes appends information about nested objects that trigger replacement when added or removed.
// It returns a SchemaAttributeDescription with the nested replacement behavior documented.
// This method is used for nested attributes where the object itself triggers replacement
// but individual parameters within the object may have different immutability rules.
func (r SchemaAttributeDescription) RequiresReplaceNestedAttributes() SchemaAttributeDescription {
	return r.AppendMarkdownString("If this object is added or removed, a replacement plan is triggered.  Parameters within the object are subject to their own immutability rules.")
}

// UnmodifiableDataLossProtection appends information about attributes that require manual replacement due to data loss risk.
// It returns a SchemaAttributeDescription with the data protection warning documented.
// This method is used for critical attributes where modification could result in data loss
// and requires explicit user action using Terraform's replace command option.
func (r SchemaAttributeDescription) UnmodifiableDataLossProtection() SchemaAttributeDescription {
	return r.AppendMarkdownString("This field is immutable and cannot be changed once defined.  To protect against accidental data loss, this resource must be replaced manually (for example, by using Terraform's [plan `-replace` command option](https://developer.hashicorp.com/terraform/cli/commands/plan#replace-address)).  Any data that is stored against this resource must be manually exported before the resource is removed and re-imported once the resource has been replaced.")
}

// AppendSliceValues appends formatted information about a slice of values to the schema description.
// It returns a SchemaAttributeDescription with the values list documented.
// The pretext parameter provides context for the values (e.g., "Options are", "Conflicts with").
// The values parameter contains the list of values to be formatted and appended.
func (r SchemaAttributeDescription) AppendSliceValues(pretext string, values []string) SchemaAttributeDescription {
	pretext = strings.TrimSpace(pretext)

	return r.AppendMarkdownString(fmt.Sprintf("%s `%s`.", pretext, strings.Join(values, "`, `")))
}

// AppendStringValue appends formatted information about a single value to the schema description.
// It returns a SchemaAttributeDescription with the value documented.
// The pretext parameter provides context for the value (e.g., "Defaults to", "Fixed value of").
// The value parameter contains the value to be formatted and appended.
func (r SchemaAttributeDescription) AppendStringValue(pretext string, value string) SchemaAttributeDescription {
	pretext = strings.TrimSpace(pretext)
	value = strings.TrimSpace(value)

	return r.AppendMarkdownString(fmt.Sprintf("%s `%s`.", pretext, value))
}

// AppendMarkdownString appends markdown-formatted text to the schema description.
// It returns a SchemaAttributeDescription with the additional text incorporated.
// The text parameter contains the markdown text to append to both description fields.
// This method handles proper punctuation and formatting when combining multiple description elements.
func (r SchemaAttributeDescription) AppendMarkdownString(text string) SchemaAttributeDescription {
	text = strings.TrimSpace(text)

	r = r.Clean(true)

	if r.Description != "" {
		r.Description = fmt.Sprintf("%s.  ", r.Description)
	}
	r.Description = fmt.Sprintf("%s%s", r.Description, schemaDescriptionMarkdownToUnformatted(text))

	if r.MarkdownDescription != "" {
		r.MarkdownDescription = fmt.Sprintf("%s.  ", r.MarkdownDescription)
	}
	r.MarkdownDescription = fmt.Sprintf("%s%s", r.MarkdownDescription, text)

	return r
}

// SchemaAttributeDescriptionFromMarkdown creates a SchemaAttributeDescription from markdown text.
// It returns a SchemaAttributeDescription with both markdown and plain text versions populated.
// The description parameter contains the markdown-formatted description text.
// This function automatically generates a plain text version by removing markdown formatting.
func SchemaAttributeDescriptionFromMarkdown(description string) SchemaAttributeDescription {
	return func() SchemaAttributeDescription {
		return SchemaAttributeDescription{
			Description:         schemaDescriptionMarkdownToUnformatted(description),
			MarkdownDescription: description,
		}
	}().Clean(false)
}
