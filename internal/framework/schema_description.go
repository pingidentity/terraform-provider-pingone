// Copyright © 2025 Ping Identity Corporation

// Package framework provides utilities for Terraform Plugin Framework implementation in the PingOne provider.
// This package contains schema description builders for resource-level documentation and validation.
package framework

import (
	"fmt"
	"regexp"
	"strings"
)

// SchemaDescription provides methods for building resource-level schema descriptions.
// It extends SchemaDescriptionModel with chainable methods for adding resource-specific
// documentation such as environment constraints and configuration warnings.
type SchemaDescription SchemaDescriptionModel

// Clean normalizes and formats the resource schema description by removing trailing punctuation and whitespace.
// It returns a cleaned SchemaDescription with standardized formatting.
// The removeTrailingStop parameter determines whether trailing periods should be removed.
// This method also ensures that both Description and MarkdownDescription fields are populated
// by copying content from one to the other when one is empty.
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

// OnlyOneDefinitionPerEnvironment appends a warning about resource uniqueness per environment.
// It returns a SchemaDescription with information about the resource constraint documented.
// The resourceName parameter specifies the name of the resource type being documented.
// This method adds a warning that only one instance of the resource should be configured per environment
// to prevent conflicts during apply operations.
func (r SchemaDescription) OnlyOneDefinitionPerEnvironment(resourceName string) SchemaDescription {
	return r.AppendMarkdownString(fmt.Sprintf("\n\n~> Only one `%s` resource should be configured for an environment.  If multiple `%s` resource definition for an environment have been defined, these are likely to conflict with each other on apply.", resourceName, resourceName))
}

// AppendSliceValues appends formatted information about a slice of values to the resource schema description.
// It returns a SchemaDescription with the values list documented.
// The pretext parameter provides context for the values (e.g., "Options are", "Conflicts with").
// The values parameter contains the list of values to be formatted and appended.
func (r SchemaDescription) AppendSliceValues(pretext string, values []string) SchemaDescription {
	pretext = strings.TrimSpace(pretext)

	return r.AppendMarkdownString(fmt.Sprintf("%s `%s`.", pretext, strings.Join(values, "`, `")))
}

// AppendStringValue appends formatted information about a single value to the resource schema description.
// It returns a SchemaDescription with the value documented.
// The pretext parameter provides context for the value (e.g., "Defaults to", "Fixed value of").
// The value parameter contains the value to be formatted and appended.
func (r SchemaDescription) AppendStringValue(pretext string, value string) SchemaDescription {
	pretext = strings.TrimSpace(pretext)
	value = strings.TrimSpace(value)

	return r.AppendMarkdownString(fmt.Sprintf("%s `%s`.", pretext, value))
}

// AppendMarkdownString appends markdown-formatted text to the resource schema description.
// It returns a SchemaDescription with the additional text incorporated.
// The text parameter contains the markdown text to append to both description fields.
// This method handles proper punctuation and formatting when combining multiple description elements.
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

// SchemaDescriptionFromMarkdown creates a SchemaDescription from markdown text.
// It returns a SchemaDescription with both markdown and plain text versions populated.
// The description parameter contains the markdown-formatted description text.
// This function automatically generates a plain text version by removing markdown formatting.
func SchemaDescriptionFromMarkdown(description string) SchemaDescription {
	return func() SchemaDescription {
		return SchemaDescription{
			Description:         schemaDescriptionMarkdownToUnformatted(description),
			MarkdownDescription: description,
		}
	}().Clean(false)
}

// schemaDescriptionMarkdownToUnformatted converts markdown-formatted text to plain text.
// It returns a string with markdown formatting removed for use in plain text descriptions.
// The description parameter contains the markdown text to be converted.
// This function currently removes backtick formatting and converts it to quote formatting.
func schemaDescriptionMarkdownToUnformatted(description string) string {
	return strings.ReplaceAll(description, "`", "\"")
}
