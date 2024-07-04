package base

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework"
	"github.com/pingidentity/terraform-provider-pingone/internal/framework/customtypes/pingonetypes"
)

type gatewayResourceModelV0 struct {
	Id            pingonetypes.ResourceIDValue `tfsdk:"id"`
	EnvironmentId pingonetypes.ResourceIDValue `tfsdk:"environment_id"`
	Name          types.String                 `tfsdk:"name"`
	Description   types.String                 `tfsdk:"description"`
	Type          types.String                 `tfsdk:"type"`
	Enabled       types.Bool                   `tfsdk:"enabled"`

	// LDAP
	BindDN                                types.String `tfsdk:"bind_dn"`
	BindPassword                          types.String `tfsdk:"bind_password"`
	ConnectionSecurity                    types.String `tfsdk:"connection_security"`
	KerberosServiceAccountPassword        types.String `tfsdk:"kerberos_service_account_password"`
	KerberosServiceAccountUpn             types.String `tfsdk:"kerberos_service_account_upn"`
	KerberosRetailPreviousCredentialsMins types.Int64  `tfsdk:"kerberos_retain_previous_credentials_mins"`
	Servers                               types.Set    `tfsdk:"servers"`
	ValidateTLSCertificates               types.Bool   `tfsdk:"validate_tls_certificates"`
	Vendor                                types.String `tfsdk:"vendor"`
	UserType                              types.Set    `tfsdk:"user_type"`

	// Radius
	RadiusClient              types.Set                    `tfsdk:"radius_client"`
	RadiusDavinciPolicyId     pingonetypes.ResourceIDValue `tfsdk:"radius_davinci_policy_id"`
	RadiusDefaultSharedSecret types.String                 `tfsdk:"radius_default_shared_secret"`
}

type gatewayUserTypeResourceModelV0 struct {
	Id                        pingonetypes.ResourceIDValue `tfsdk:"id"`
	Name                      types.String                 `tfsdk:"name"`
	PasswordAuthority         types.String                 `tfsdk:"password_authority"`
	SearchBaseDN              types.String                 `tfsdk:"search_base_dn"`
	UserLinkAttributes        types.List                   `tfsdk:"user_link_attributes"`
	UserMigration             types.List                   `tfsdk:"user_migration"`
	PushPasswordChangesToLdap types.Bool                   `tfsdk:"push_password_changes_to_ldap"`
}

type gatewayUserTypeUserMigrationResourceModelV0 struct {
	AttributeMappings   types.Set                    `tfsdk:"attribute_mapping"`
	LookupFilterPattern types.String                 `tfsdk:"lookup_filter_pattern"`
	PopulationId        pingonetypes.ResourceIDValue `tfsdk:"population_id"`
}

func (r *GatewayResource) UpgradeState(ctx context.Context) map[int64]resource.StateUpgrader {
	return map[int64]resource.StateUpgrader{
		// State upgrade implementation from 0 (prior state version) to 1 (Schema.Version)
		0: {
			PriorSchema: &schema.Schema{
				Attributes: map[string]schema.Attribute{
					"id": framework.Attr_ID(),

					"environment_id": framework.Attr_LinkID(
						framework.SchemaAttributeDescriptionFromMarkdown(""),
					),

					"name": schema.StringAttribute{
						Required: true,
					},

					"description": schema.StringAttribute{
						Optional: true,
					},

					"type": schema.StringAttribute{
						Required: true,
					},

					"enabled": schema.BoolAttribute{
						Required: true,
					},

					"bind_dn": schema.StringAttribute{
						Optional: true,
					},

					"bind_password": schema.StringAttribute{
						Optional:  true,
						Sensitive: true,
					},

					"connection_security": schema.StringAttribute{
						Optional: true,
						Computed: true,

						Default: stringdefault.StaticString(string(management.ENUMGATEWAYTYPELDAPSECURITY_NONE)),
					},

					"kerberos_service_account_password": schema.StringAttribute{
						Optional:  true,
						Sensitive: true,
					},

					"kerberos_service_account_upn": schema.StringAttribute{
						Optional: true,
					},

					"kerberos_retain_previous_credentials_mins": schema.Int64Attribute{
						Optional: true,
					},

					"servers": schema.SetAttribute{
						Optional: true,

						ElementType: types.StringType,
					},

					"validate_tls_certificates": schema.BoolAttribute{
						Optional: true,
						Computed: true,

						Default: booldefault.StaticBool(true),
					},

					"vendor": schema.StringAttribute{
						Optional: true,
					},

					"radius_davinci_policy_id": schema.StringAttribute{
						Optional: true,

						CustomType: pingonetypes.ResourceIDType{},
					},

					"radius_default_shared_secret": schema.StringAttribute{
						Optional:  true,
						Sensitive: true,
					},
				},

				Blocks: map[string]schema.Block{
					"user_type": schema.SetNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									Computed: true,

									CustomType: pingonetypes.ResourceIDType{},
								},

								"name": schema.StringAttribute{
									Required: true,
								},

								"password_authority": schema.StringAttribute{
									Required: true,
								},

								"search_base_dn": schema.StringAttribute{
									Required: true,
								},

								"user_link_attributes": schema.ListAttribute{
									Required: true,

									ElementType: types.StringType,
								},

								"push_password_changes_to_ldap": schema.BoolAttribute{
									Optional: true,
									Computed: true,

									Default: booldefault.StaticBool(false),
								},
							},

							Blocks: map[string]schema.Block{
								"user_migration": schema.ListNestedBlock{
									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"lookup_filter_pattern": schema.StringAttribute{
												Required: true,
											},

											"population_id": schema.StringAttribute{
												Required: true,

												CustomType: pingonetypes.ResourceIDType{},
											},
										},

										Blocks: map[string]schema.Block{
											"attribute_mapping": schema.SetNestedBlock{
												NestedObject: schema.NestedBlockObject{
													Attributes: map[string]schema.Attribute{
														"name": schema.StringAttribute{
															Required: true,
														},

														"value": schema.StringAttribute{
															Required: true,
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},

					"radius_client": schema.SetNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"ip": schema.StringAttribute{
									Required: true,
								},

								"shared_secret": schema.StringAttribute{
									Optional:  true,
									Sensitive: true,
								},
							},
						},
					},
				},
			},
			StateUpgrader: func(ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse) {
				var d diag.Diagnostics
				var priorStateData gatewayResourceModelV0

				resp.Diagnostics.Append(req.State.Get(ctx, &priorStateData)...)

				if resp.Diagnostics.HasError() {
					return
				}

				kerberos, d := priorStateData.schemaUpgradeKerberosV0toV1(ctx)
				resp.Diagnostics.Append(d...)

				userTypes, d := priorStateData.schemaUpgradeUserTypesV0toV1(ctx)
				resp.Diagnostics.Append(d...)

				upgradedStateData := gatewayResourceModelV1{
					Id:                        priorStateData.Id,
					EnvironmentId:             priorStateData.EnvironmentId,
					Name:                      priorStateData.Name,
					Description:               priorStateData.Description,
					Type:                      priorStateData.Type,
					Enabled:                   priorStateData.Enabled,
					BindDN:                    priorStateData.BindDN,
					BindPassword:              priorStateData.BindPassword,
					ConnectionSecurity:        priorStateData.ConnectionSecurity,
					FollowReferrals:           types.BoolNull(),
					Kerberos:                  kerberos,
					Servers:                   priorStateData.Servers,
					ValidateTLSCertificates:   priorStateData.ValidateTLSCertificates,
					Vendor:                    priorStateData.Vendor,
					UserTypes:                 userTypes,
					RadiusClients:             priorStateData.RadiusClient,
					RadiusDavinciPolicyId:     priorStateData.RadiusDavinciPolicyId,
					RadiusDefaultSharedSecret: priorStateData.RadiusDefaultSharedSecret,
					RadiusNetworkPolicyServer: types.ObjectNull(gatewayRadiusNetworkPolicyServerTFObjectTypes),
				}

				resp.Diagnostics.Append(resp.State.Set(ctx, upgradedStateData)...)
			},
		},
	}
}

func (p *gatewayResourceModelV0) schemaUpgradeKerberosV0toV1(ctx context.Context) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributeTypes := gatewayKerberosTFObjectTypes

	if p.KerberosRetailPreviousCredentialsMins.IsNull() && p.KerberosServiceAccountPassword.IsNull() && p.KerberosServiceAccountUpn.IsNull() {
		return types.ObjectNull(attributeTypes), diags
	} else if p.KerberosRetailPreviousCredentialsMins.IsUnknown() || p.KerberosServiceAccountPassword.IsUnknown() || p.KerberosServiceAccountUpn.IsUnknown() {
		return types.ObjectUnknown(attributeTypes), diags
	} else {

		upgradedStateData := gatewayKerberosResourceModelV1{
			ServiceAccountPassword:        p.KerberosServiceAccountPassword,
			ServiceAccountUPN:             p.KerberosServiceAccountUpn,
			RetainPreviousCredentialsMins: p.KerberosRetailPreviousCredentialsMins,
		}

		returnVar, d := types.ObjectValueFrom(ctx, attributeTypes, upgradedStateData)
		diags.Append(d...)

		return returnVar, diags
	}
}

func (p *gatewayResourceModelV0) schemaUpgradeUserTypesV0toV1(ctx context.Context) (types.Map, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributeTypes := types.ObjectType{AttrTypes: gatewayUserTypesTFObjectTypes}
	planAttribute := p.UserType

	if planAttribute.IsNull() {
		return types.MapNull(attributeTypes), diags
	} else if planAttribute.IsUnknown() {
		return types.MapUnknown(attributeTypes), diags
	} else {
		var priorStateData []gatewayUserTypeResourceModelV0
		d := planAttribute.ElementsAs(ctx, &priorStateData, false)
		diags.Append(d...)
		if diags.HasError() {
			return types.MapNull(attributeTypes), diags
		}

		if len(priorStateData) == 0 {
			return types.MapNull(attributeTypes), diags
		}

		upgradedStateData := make(map[string]gatewayUserTypeResourceModelV1)

		for _, userType := range priorStateData {

			newUserLookup, d := userType.schemaUpgradeNewUserLookupV0toV1(ctx)
			diags.Append(d...)

			upgradedStateData[userType.Name.ValueString()] = gatewayUserTypeResourceModelV1{
				AllowPasswordChanges:                 userType.PushPasswordChangesToLdap,
				Id:                                   userType.Id,
				NewUserLookup:                        newUserLookup,
				PasswordAuthority:                    userType.PasswordAuthority,
				SearchBaseDN:                         userType.SearchBaseDN,
				UpdateUserOnSuccessfulAuthentication: types.BoolValue(false),
				UserLinkAttributes:                   userType.UserLinkAttributes,
			}
		}

		returnVar, d := types.MapValueFrom(ctx, attributeTypes, upgradedStateData)
		diags.Append(d...)

		return returnVar, diags
	}
}

func (p *gatewayUserTypeResourceModelV0) schemaUpgradeNewUserLookupV0toV1(ctx context.Context) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributeTypes := gatewayUserTypesNewUserLookupTFObjectTypes
	planAttribute := p.UserMigration

	if planAttribute.IsNull() {
		return types.ObjectNull(attributeTypes), diags
	} else if planAttribute.IsUnknown() {
		return types.ObjectUnknown(attributeTypes), diags
	} else {
		var priorStateData []gatewayUserTypeUserMigrationResourceModelV0
		d := planAttribute.ElementsAs(ctx, &priorStateData, false)
		diags.Append(d...)
		if diags.HasError() {
			return types.ObjectNull(attributeTypes), diags
		}

		if len(priorStateData) == 0 {
			return types.ObjectNull(attributeTypes), diags
		}

		upgradedStateData := gatewayUserTypeNewUserLookupResourceModelV1{
			AttributeMappings: priorStateData[0].AttributeMappings,
			LDAPFilterPattern: priorStateData[0].LookupFilterPattern,
			PopulationId:      priorStateData[0].PopulationId,
		}

		returnVar, d := types.ObjectValueFrom(ctx, attributeTypes, upgradedStateData)
		diags.Append(d...)

		return returnVar, diags
	}
}
