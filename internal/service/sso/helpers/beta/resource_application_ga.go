// Copyright © 2025 Ping Identity Corporation

// This file relates to a beta feature described in CDI-492 and should be modified or removed on completion of CDI-631

//go:build !beta

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

type ApplicationOIDCOptionsResourceModelV1 struct {
	ClientId types.String `tfsdk:"client_id"`
}

var ApplicationOidcOptionsTFObjectTypes = map[string]attr.Type{
	"client_id": types.StringType,
}

func ResourceSchemaItems() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"client_id": schema.StringAttribute{
			Description: framework.SchemaAttributeDescriptionFromMarkdown("A string that specifies the application ID used to authenticate to the authorization server.").Description,
			Computed:    true,

			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseNonNullStateForUnknown(),
			},
		},
	}
}

// no-op
func AddBeta(data *management.ApplicationOIDC, plan ApplicationOIDCOptionsResourceModelV1) {}

func SchemaUpgradeV0toV1(clientId pingonetypes.ResourceIDValue) ApplicationOIDCOptionsResourceModelV1 {
	return ApplicationOIDCOptionsResourceModelV1{
		ClientId: types.StringValue(clientId.ValueString()),
	}
}

func ApplicationBetaToTF(apiObject *management.ApplicationOIDC, stateValue ApplicationOIDCOptionsResourceModelV1) map[string]attr.Value {
	return map[string]attr.Value{
		"client_id": framework.StringOkToTF(apiObject.GetIdOk()),
	}
}
