package framework

import (
	"fmt"
	"regexp"
	"sort"
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

func (r SchemaAttributeDescription) DefaultValue(defaultValue string) SchemaAttributeDescription {
	return r.AppendStringValue("Defaults to", defaultValue)
}

func (r SchemaAttributeDescription) AllowedValues(allowedValues []string) SchemaAttributeDescription {
	sort.Strings(allowedValues)

	return r.AppendSliceValues("Options are", allowedValues)
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
	return r.AllowedValues(utils.EnumSliceToStringSlice(allowedValuesEnumSlice))
}

func (r SchemaAttributeDescription) ConflictsWith(fieldPaths []string) SchemaAttributeDescription {
	return r.AppendSliceValues("Conflicts with", fieldPaths)
}

func (r SchemaAttributeDescription) ExactlyOneOf(fieldPaths []string) SchemaAttributeDescription {
	return r.AppendSliceValues("At least one of the following must be defined:", fieldPaths)
}

func (r SchemaAttributeDescription) RequiresReplace() SchemaAttributeDescription {
	return r.AppendMarkdownString("This field is immutable and will trigger a replace plan if changed.")
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
