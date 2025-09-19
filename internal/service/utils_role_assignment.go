// Copyright © 2025 Ping Identity Corporation

// Package service provides utility functions and common types for handling role assignment operations
// across different PingOne services in the Terraform provider.
package service

import (
	"context"
	"fmt"
	"net/http"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
)

// roleAssignmentScopeDescription generates a schema attribute description for role assignment scope fields.
// It returns a framework.SchemaAttributeDescription with appropriate text for the specified scope type.
// The scopeType parameter specifies the type of scope (e.g., "application", "environment") for the description.
// The someRolesCannotBeScoped parameter indicates whether to include a warning about roles that cannot be scoped to this type.
func roleAssignmentScopeDescription(scopeType string, someRolesCannotBeScoped bool) framework.SchemaAttributeDescription {
	description := framework.SchemaAttributeDescriptionFromMarkdown(
		fmt.Sprintf("Limit the scope of the admin role assignment to the specified %s ID.  Must be a valid PingOne resource ID.", scopeType),
	)

	if someRolesCannotBeScoped {
		description = description.AppendMarkdownString(
			fmt.Sprintf("Some roles cannot be scoped to the %s.", scopeType),
		)
	}

	return description.RequiresReplace()
}

// RoleAssignmentScopeSchema returns a map of schema attributes for role assignment scope configuration.
// It returns a map[string]schema.Attribute containing all scope-related attributes with proper validation and constraints.
// This function defines mutually exclusive scope attributes for application, environment, organization, and population scopes.
// Each scope attribute is configured with appropriate validators to ensure exactly one scope type is specified.
func RoleAssignmentScopeSchema() map[string]schema.Attribute {

	exactlyOneOfParameters := []string{"scope_application_id", "scope_organization_id", "scope_environment_id", "scope_population_id"}

	applicationIdDescription := roleAssignmentScopeDescription("application", true).ExactlyOneOf(exactlyOneOfParameters)
	environmentIdDescription := roleAssignmentScopeDescription("environment", true).ExactlyOneOf(exactlyOneOfParameters)
	organizationIdDescription := roleAssignmentScopeDescription("organization", true).ExactlyOneOf(exactlyOneOfParameters)
	populationIdDescription := roleAssignmentScopeDescription("population", true).ExactlyOneOf(exactlyOneOfParameters)

	validators := []validator.String{
		stringvalidator.ExactlyOneOf(
			path.MatchRelative().AtParent().AtName("scope_application_id"),
			path.MatchRelative().AtParent().AtName("scope_environment_id"),
			path.MatchRelative().AtParent().AtName("scope_organization_id"),
			path.MatchRelative().AtParent().AtName("scope_population_id"),
		),
	}

	planModifiers := []planmodifier.String{
		stringplanmodifier.RequiresReplace(),
	}

	return map[string]schema.Attribute{
		"scope_application_id": schema.StringAttribute{
			Description:         applicationIdDescription.Description,
			MarkdownDescription: applicationIdDescription.MarkdownDescription,
			Optional:            true,

			CustomType: pingonetypes.ResourceIDType{},

			PlanModifiers: planModifiers,

			Validators: validators,
		},

		"scope_environment_id": schema.StringAttribute{
			Description:         environmentIdDescription.Description,
			MarkdownDescription: environmentIdDescription.MarkdownDescription,
			Optional:            true,

			CustomType: pingonetypes.ResourceIDType{},

			PlanModifiers: planModifiers,

			Validators: validators,
		},

		"scope_organization_id": schema.StringAttribute{
			Description:         organizationIdDescription.Description,
			MarkdownDescription: organizationIdDescription.MarkdownDescription,
			Optional:            true,

			CustomType: pingonetypes.ResourceIDType{},

			PlanModifiers: planModifiers,

			Validators: validators,
		},

		"scope_population_id": schema.StringAttribute{
			Description:         populationIdDescription.Description,
			MarkdownDescription: populationIdDescription.MarkdownDescription,
			Optional:            true,

			CustomType: pingonetypes.ResourceIDType{},

			PlanModifiers: planModifiers,

			Validators: validators,
		},
	}
}

// ExpandRoleAssignmentScope extracts and validates role assignment scope information from Terraform configuration values.
// It returns the scope ID, scope type string, and any diagnostics encountered during processing.
// The function examines the provided scope values to determine which type of scope is configured.
// All scope parameters (scopeEnvironmentID, scopeOrganizationID, scopePopulationID, scopeApplicationID) are ResourceIDValue types from Terraform configuration.
// Exactly one of the scope parameters must be set to a non-null, non-unknown value.
func ExpandRoleAssignmentScope(scopeEnvironmentID, scopeOrganizationID, scopePopulationID, scopeApplicationID pingonetypes.ResourceIDValue) (scopeId, scopeType string, diags diag.Diagnostics) {

	if scopeApplicationID != pingonetypes.ResourceIDNull() && scopeApplicationID != pingonetypes.ResourceIDUnknown() {
		return scopeApplicationID.ValueString(), "APPLICATION", diags
	}

	if scopeEnvironmentID != pingonetypes.ResourceIDNull() && scopeEnvironmentID != pingonetypes.ResourceIDUnknown() {
		return scopeEnvironmentID.ValueString(), "ENVIRONMENT", diags
	}

	if scopeOrganizationID != pingonetypes.ResourceIDNull() && scopeOrganizationID != pingonetypes.ResourceIDUnknown() {
		return scopeOrganizationID.ValueString(), "ORGANIZATION", diags
	}

	if scopePopulationID != pingonetypes.ResourceIDNull() && scopePopulationID != pingonetypes.ResourceIDUnknown() {
		return scopePopulationID.ValueString(), "POPULATION", diags
	}

	diags.AddError(
		"Invalid configuration",
		"One of scope_application_id, scope_organization_id, scope_environment_id or scope_population_id must be set",
	)

	return "", "", diags

}

// RoleAssignmentScopeOkToTF converts a PingOne API role assignment scope object to Terraform Framework resource ID values.
// It returns individual ResourceIDValue instances for each possible scope type (application, environment, organization, population).
// The roleAssignmentScope parameter contains the scope information from the PingOne API response.
// The ok parameter indicates whether the scope data was successfully retrieved from the API.
// Only one of the returned values will be populated based on the scope type, while others will be null.
func RoleAssignmentScopeOkToTF(roleAssignmentScope *management.RoleAssignmentScope, ok bool) (scopeEnvironmentId, scopeOrganizationId, scopePopulationId, scopeApplicationId pingonetypes.ResourceIDValue) {
	scopeApplicationId = pingonetypes.NewResourceIDNull()
	scopeEnvironmentId = pingonetypes.NewResourceIDNull()
	scopeOrganizationId = pingonetypes.NewResourceIDNull()
	scopePopulationId = pingonetypes.NewResourceIDNull()

	if ok {
		if scopeType, ok := roleAssignmentScope.GetTypeOk(); ok {

			switch *scopeType {
			case management.ENUMROLEASSIGNMENTSCOPETYPE_APPLICATION:
				scopeApplicationId = framework.PingOneResourceIDOkToTF(roleAssignmentScope.GetIdOk())
			case management.ENUMROLEASSIGNMENTSCOPETYPE_ENVIRONMENT:
				scopeEnvironmentId = framework.PingOneResourceIDOkToTF(roleAssignmentScope.GetIdOk())
			case management.ENUMROLEASSIGNMENTSCOPETYPE_ORGANIZATION:
				scopeOrganizationId = framework.PingOneResourceIDOkToTF(roleAssignmentScope.GetIdOk())
			case management.ENUMROLEASSIGNMENTSCOPETYPE_POPULATION:
				scopePopulationId = framework.PingOneResourceIDOkToTF(roleAssignmentScope.GetIdOk())
			}

		}
	}

	return
}

// CreateRoleAssignmentErrorFunc provides custom error handling for role assignment creation operations.
// It returns a function that processes HTTP responses and PingOne API errors to provide user-friendly error messages.
// This error handler specifically looks for invalid role/scope combination errors and formats them appropriately.
// The returned function signature matches the CustomError parameter expected by framework.ParseResponse.
var (
	CreateRoleAssignmentErrorFunc = func(_ *http.Response, p1Error *model.P1Error) diag.Diagnostics {
		var diags diag.Diagnostics

		// Invalid role/scope combination
		if details, ok := p1Error.GetDetailsOk(); ok && details != nil && len(details) > 0 {
			if target, ok := details[0].GetTargetOk(); ok && *target == "scope" {
				diags.AddError(
					"Incompatible role and scope combination",
					details[0].GetMessage(),
				)

				return diags
			}
		}

		return diags
	}

	// RoleAssignmentRetryable provides custom retry logic for role assignment operations.
	// It returns a function that determines whether a failed API call should be retried based on the error response.
	// This function handles authorization errors that may occur due to permission propagation delays.
	// The ctx parameter provides context for logging retry decisions.
	// The r parameter contains the HTTP response from the failed API call.
	// The p1error parameter contains the parsed PingOne API error details.
	RoleAssignmentRetryable = func(ctx context.Context, r *http.Response, p1error *model.P1Error) bool {

		if p1error != nil {
			var err error

			// Permissions may not have propagated by this point (1)
			m, err := regexp.MatchString("^The actor attempting to perform the request is not authorized.", p1error.GetMessage())
			if err == nil && m {
				tflog.Warn(ctx, "Insufficient PingOne privileges detected")
				return true
			}
			if err != nil {
				tflog.Warn(ctx, "Cannot match error string for retry")
				return false
			}

			// Permissions may not have propagated by this point (2)
			if details, ok := p1error.GetDetailsOk(); ok && details != nil && len(details) > 0 {
				m, err := regexp.MatchString("^Must have role at the same or broader scope", details[0].GetMessage())
				if err == nil && m {
					tflog.Warn(ctx, "Insufficient PingOne privileges detected")
					return true
				}
				if err != nil {
					tflog.Warn(ctx, "Cannot match error string for retry")
					return false
				}
			}

		}

		return false
	}

	// RoleRemovalRetryable provides retry logic for role removal operations.
	// It returns a function that determines whether a failed role removal API call should be retried.
	// Currently, this function always returns false, indicating that role removal operations should not be retried.
	// The ctx parameter provides context for the operation.
	// The r parameter contains the HTTP response from the failed API call.
	// The p1error parameter contains the parsed PingOne API error details.
	RoleRemovalRetryable = func(ctx context.Context, r *http.Response, p1error *model.P1Error) bool {
		return false
	}
)
