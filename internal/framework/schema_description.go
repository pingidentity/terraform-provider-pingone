// Copyright © 2026 Ping Identity Corporation

package framework

import (
	"fmt"
	"regexp"
	"strings"
)

type SchemaDescription SchemaDescriptionModel

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

func (r SchemaDescription) OnlyOneDefinitionPerEnvironment(resourceName string) SchemaDescription {
	return r.AppendMarkdownString(fmt.Sprintf("\n\n~> Only one `%s` resource should be configured for an environment.  If multiple `%s` resource definition for an environment have been defined, these are likely to conflict with each other on apply.", resourceName, resourceName))
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

func (r SchemaDescription) Beta(text string) SchemaDescription {
	return r.AppendMarkdownString(fmt.Sprintf("**Beta or experimental**. Use of this resource is subject to change at any time and to be used with caution. The API may change without notice which may lead to errors. %s", text))
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
