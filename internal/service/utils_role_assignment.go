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

	RoleRemovalRetryable = func(ctx context.Context, r *http.Response, p1error *model.P1Error) bool {
		return false
	}
)
