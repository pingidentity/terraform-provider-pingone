package framework

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/pingidentity/terraform-provider-pingone/internal/utils"
)

type SchemaDescription struct {
	Description         string
	MarkdownDescription string
}

func (r SchemaDescription) Clean(removeTrailingStop bool) SchemaDescription {

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

func (r SchemaDescription) DefaultValue(defaultValue string) SchemaDescription {
	return r.AppendStringValue("Defaults to", defaultValue)
}

func (r SchemaDescription) AllowedValues(allowedValues []string) SchemaDescription {
	sort.Strings(allowedValues)

	return r.AppendSliceValues("Options are", allowedValues)
}

func (r SchemaDescription) AllowedValuesComplex(allowedValuesMap map[string]string) SchemaDescription {
	allowedValues := make([]string, 0)
	for k, v := range allowedValuesMap {
		allowedValues = append(allowedValues, fmt.Sprintf("`%s` (%s)", k, v))
	}

	sort.Strings(allowedValues)

	return r.AppendMarkdownString(fmt.Sprintf("Options are %s.", strings.Join(allowedValues, ", ")))
}

func (r SchemaDescription) AllowedValuesEnum(allowedValuesEnumSlice interface{}) SchemaDescription {
	return r.AllowedValues(utils.EnumSliceToStringSlice(allowedValuesEnumSlice))
}

func (r SchemaDescription) ConflictsWith(fieldPaths []string) SchemaDescription {
	return r.AppendSliceValues("Conflicts with", fieldPaths)
}

func (r SchemaDescription) ExactlyOneOf(fieldPaths []string) SchemaDescription {
	return r.AppendSliceValues("At least one of the following must be defined:", fieldPaths)
}

func (r SchemaDescription) RequiresReplace() SchemaDescription {
	return r.AppendMarkdownString("This field is immutable and will trigger a replace plan if changed.")
}

func (r SchemaDescription) AppendSliceValues(pretext string, values []string) SchemaDescription {
	pretext = strings.TrimSpace(pretext)

	return r.AppendMarkdownString(fmt.Sprintf("%s `%s`.", pretext, strings.Join(values, "`, `")))
}

func (r SchemaDescription) AppendStringValue(pretext string, value string) SchemaDescription {
	pretext = strings.TrimSpace(pretext)
	value = strings.TrimSpace(value)

	return r.AppendMarkdownString(fmt.Sprintf("%s `%s`.", pretext, value))
}

func (r SchemaDescription) AppendMarkdownString(text string) SchemaDescription {
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

func SchemaDescriptionFromMarkdown(description string) SchemaDescription {
	return func() SchemaDescription {
		return SchemaDescription{
			Description:         schemaDescriptionMarkdownToUnformatted(description),
			MarkdownDescription: description,
		}
	}().Clean(false)
}

func schemaDescriptionMarkdownToUnformatted(description string) string {
	return strings.ReplaceAll(description, "`", "\"")
}
