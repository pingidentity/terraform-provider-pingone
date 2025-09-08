// Copyright © 2025 Ping Identity Corporation

// This file relates to a beta feature described in CDI-492

//go:build beta

package beta

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
)

type ApplicationOIDCOptionsResourceModelV1Beta struct {
	ClientId            types.String `tfsdk:"client_id"`
	InitialClientSecret types.String `tfsdk:"initial_client_secret"`
}

var ApplicationOidcOptionsTFObjectTypes = map[string]attr.Type{
	"client_id":             types.StringType,
	"initial_client_secret": types.StringType,
}

func ClientIdClientSecretSchemaItems() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"client_id": schema.StringAttribute{
			Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the application client ID used to authenticate to the authorization server. If left undefined, the service will generate a value.").RequiresReplace().Beta("To modify the value of this field, the environment must be enabled with the feature flag to allow importing applications with administrator defined client ID and client secret values.").Description,
			Optional:    true,
			Computed:    true,

			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"initial_client_secret": schema.StringAttribute{
			Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the initial import client secret used to authenticate to the authorization server. If left undefined, the service will generate a value. Note this field's value will not change if the secret is rotated. After initial import, the `pingone_application_secret` resource or data source should be used to refer to the current active value of the application's secret.").RequiresReplace().Beta("To modify the value of this field, the environment must be enabled with the feature flag to allow importing applications with administrator defined client ID and client secret values.").Description,
			Optional:    true,
			Sensitive:   true,
		},
	}
}

func AddBeta(data *management.ApplicationOIDC, plan ApplicationOIDCOptionsResourceModelV1Beta) {
	if !plan.ClientId.IsNull() && !plan.ClientId.IsUnknown() {
		data.SetClientId(plan.ClientId.ValueString())
	}
	if !plan.InitialClientSecret.IsNull() && !plan.InitialClientSecret.IsUnknown() {
		data.SetClientSecret(plan.InitialClientSecret.ValueString())
	}
}

func SchemaUpgradeV0toV1(clientId pingonetypes.ResourceIDValue) ApplicationOIDCOptionsResourceModelV1Beta {
	return ApplicationOIDCOptionsResourceModelV1Beta{
		ClientId:            types.StringValue(clientId.ValueString()),
		InitialClientSecret: types.StringNull(),
	}
}

func ApplicationBetaToTF(apiObject *management.ApplicationOIDC, stateValue ApplicationOIDCOptionsResourceModelV1Beta) map[string]attr.Value {
	clientId := framework.StringOkToTF(apiObject.GetClientIdOk())

	if clientId.IsNull() {
		clientId = framework.StringOkToTF(apiObject.GetIdOk())
	}

	return map[string]attr.Value{
		"client_id":             clientId,
		"initial_client_secret": stateValue.InitialClientSecret,
	}
}
