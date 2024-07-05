package framework

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/pingidentity/terraform-provider-pingone/internal/utils"
)

type SchemaDescriptionModel struct {
	Description         string
	MarkdownDescription string
}

type SchemaAttributeDescription SchemaDescriptionModel

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

func (r SchemaAttributeDescription) DefaultValue(defaultValue any) SchemaAttributeDescription {
	return r.genericValue("Defaults to", defaultValue)
}

func (r SchemaAttributeDescription) FixedValue(defaultValue any) SchemaAttributeDescription {
	return r.genericValue("Fixed value of", defaultValue)
}

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

func (r SchemaAttributeDescription) AllowedValuesComplex(allowedValuesMap map[string]string) SchemaAttributeDescription {
	allowedValues := make([]string, 0)
	for k, v := range allowedValuesMap {
		allowedValues = append(allowedValues, fmt.Sprintf("`%s` (%s)", k, v))
	}

	sort.Strings(allowedValues)

	return r.AppendMarkdownString(fmt.Sprintf("Options are %s.", strings.Join(allowedValues, ", ")))
}

func (r SchemaAttributeDescription) AllowedValuesEnum(allowedValuesEnumSlice interface{}) SchemaAttributeDescription {
	return r.AllowedValues(utils.EnumSliceToAnySlice(allowedValuesEnumSlice)...)
}

func (r SchemaAttributeDescription) ConflictsWith(fieldPaths []string) SchemaAttributeDescription {
	return r.AppendSliceValues("Conflicts with", fieldPaths)
}

func (r SchemaAttributeDescription) ExactlyOneOf(fieldPaths []string) SchemaAttributeDescription {
	return r.AppendSliceValues("Exactly one of the following must be defined:", fieldPaths)
}

func (r SchemaAttributeDescription) RequiresReplace() SchemaAttributeDescription {
	return r.AppendMarkdownString("This field is immutable and will trigger a replace plan if changed.")
}

func (r SchemaAttributeDescription) RequiresReplaceNestedAttributes() SchemaAttributeDescription {
	return r.AppendMarkdownString("If this object is added or removed, a replacement plan is triggered.  Parameters within the object are subject to their own immutability rules.")
}

func (r SchemaAttributeDescription) UnmodifiableDataLossProtection() SchemaAttributeDescription {
	return r.AppendMarkdownString("This field is immutable and cannot be changed once defined.  To protect against accidental data loss, this resource must be replaced manually (for example, by using Terraform's [plan `-replace` command option](https://developer.hashicorp.com/terraform/cli/commands/plan#replace-address)).  Any data that is stored against this resource must be manually exported before the resource is removed and re-imported once the resource has been replaced.")
}

func (r SchemaAttributeDescription) AppendSliceValues(pretext string, values []string) SchemaAttributeDescription {
	pretext = strings.TrimSpace(pretext)

	return r.AppendMarkdownString(fmt.Sprintf("%s `%s`.", pretext, strings.Join(values, "`, `")))
}

func (r SchemaAttributeDescription) AppendStringValue(pretext string, value string) SchemaAttributeDescription {
	pretext = strings.TrimSpace(pretext)
	value = strings.TrimSpace(value)

	return r.AppendMarkdownString(fmt.Sprintf("%s `%s`.", pretext, value))
}

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

func SchemaAttributeDescriptionFromMarkdown(description string) SchemaAttributeDescription {
	return func() SchemaAttributeDescription {
		return SchemaAttributeDescription{
			Description:         schemaDescriptionMarkdownToUnformatted(description),
			MarkdownDescription: description,
		}
	}().Clean(false)
}
