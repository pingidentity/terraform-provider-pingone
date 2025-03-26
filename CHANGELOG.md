## 1.6.0 (Unreleased)

## 1.5.0 (26 February 2025)

NOTES:

* bump `github.com/hashicorp/terraform-plugin-framework-validators` 0.16.0 => 0.17.0 ([#994](https://github.com/pingidentity/terraform-provider-pingone/issues/994))
* bump `github.com/hashicorp/terraform-plugin-framework` 1.13.0 => 1.14.1 ([#997](https://github.com/pingidentity/terraform-provider-pingone/issues/997))
* bump `github.com/hashicorp/terraform-plugin-sdk/v2` 2.35.0 => 2.36.0 ([#981](https://github.com/pingidentity/terraform-provider-pingone/issues/981))
* bump `github.com/hashicorp/terraform-plugin-sdk/v2` 2.36.0 => 2.36.1 ([#996](https://github.com/pingidentity/terraform-provider-pingone/issues/996))
* bump `github.com/patrickcping/pingone-go-sdk-v2/authorize` 0.7.0 => 0.8.0 ([#987](https://github.com/pingidentity/terraform-provider-pingone/issues/987))
* bump `github.com/patrickcping/pingone-go-sdk-v2/credentials` 0.10.0 => 0.11.0 ([#985](https://github.com/pingidentity/terraform-provider-pingone/issues/985))
* bump `github.com/patrickcping/pingone-go-sdk-v2/management` 0.47.0 => 0.49.0 ([#986](https://github.com/pingidentity/terraform-provider-pingone/issues/986))
* bump `github.com/patrickcping/pingone-go-sdk-v2/management` 0.49 => 0.51.0 ([#999](https://github.com/pingidentity/terraform-provider-pingone/issues/999))
* bump `github.com/patrickcping/pingone-go-sdk-v2/mfa` 0.22.0 => 0.23.0 ([#982](https://github.com/pingidentity/terraform-provider-pingone/issues/982))
* bump `github.com/patrickcping/pingone-go-sdk-v2/risk` 0.18.0 => 0.19.0 ([#984](https://github.com/pingidentity/terraform-provider-pingone/issues/984))
* bump `github.com/patrickcping/pingone-go-sdk-v2/verify` 0.8.0 => 0.9.0 ([#982](https://github.com/pingidentity/terraform-provider-pingone/issues/982))
* bump `github.com/patrickcping/pingone-go-sdk-v2` 0.12.7 => 0.12.9 ([#982](https://github.com/pingidentity/terraform-provider-pingone/issues/982))
* bump `github.com/patrickcping/pingone-go-sdk-v2` 0.12.9 => 0.12.11 ([#999](https://github.com/pingidentity/terraform-provider-pingone/issues/999))

ENHANCEMENTS:

* `resource/identity_provider`: Added support for the `microsoft.tenant_id` field. ([#999](https://github.com/pingidentity/terraform-provider-pingone/issues/999))

BUG FIXES:

* Fix `service_endpoints` and `global_options` provider parameters when defined explicitly in the HCL provider configuration. ([#998](https://github.com/pingidentity/terraform-provider-pingone/issues/998))
* `resource/pingone_risk_policy`: Fix policy naming regex validation to allow numbers. ([#991](https://github.com/pingidentity/terraform-provider-pingone/issues/991))

## 1.4.0 (29 January 2025)

NOTES:

* bump `github.com/hashicorp/terraform-plugin-framework-timeouts` 0.4.1 => 0.5.0 ([#968](https://github.com/pingidentity/terraform-provider-pingone/issues/968))
* bump `github.com/hashicorp/terraform-plugin-go` 0.25.0 => 0.26.0 ([#971](https://github.com/pingidentity/terraform-provider-pingone/issues/971))
* bump `github.com/hashicorp/terraform-plugin-mux` 0.17.0 => 0.18.0 ([#971](https://github.com/pingidentity/terraform-provider-pingone/issues/971))

FEATURES:

* **New Data Source:** `pingone_custom_role` ([#965](https://github.com/pingidentity/terraform-provider-pingone/issues/965))
* **New Data Source:** `pingone_custom_roles` ([#965](https://github.com/pingidentity/terraform-provider-pingone/issues/965))
* **New Resource:** `pingone_custom_role` ([#965](https://github.com/pingidentity/terraform-provider-pingone/issues/965))

ENHANCEMENTS:

* `data-source/role`: Added support for retrieving the `Advanced Identity Cloud Super Admin`,  `Advanced Identity Cloud Tenant Admin`, and `Custom Roles Admin` roles by name. ([#969](https://github.com/pingidentity/terraform-provider-pingone/issues/969))
* `data-source/trust_email_domain_ownership`: Added support for the `environment_dns_record` field. ([#969](https://github.com/pingidentity/terraform-provider-pingone/issues/969))

BUG FIXES:

* `resource/pingone_verify_policy`: Fixed the handling of otp default values in the email and phone property objects when verify set to DISABLED. ([#967](https://github.com/pingidentity/terraform-provider-pingone/issues/967))

## 1.3.1 (7 January 2025)

NOTES:

* bump `golang.org/x/net` 0.28.0 => 0.33.0 ([#954](https://github.com/pingidentity/terraform-provider-pingone/issues/954))

## 1.3.0 (19 December 2024)

BREAKING CHANGES:

* `resource/pingone_risk_predictor`: To ensure correct composite predictor and Terraform behaviours, the `predictor_composite.composition` field has been removed and replaced with `predictor_composite.compositions` field. ([#952](https://github.com/pingidentity/terraform-provider-pingone/issues/952))

NOTES:

* bump `github.com/hashicorp/terraform-plugin-framework-validators` 0.15.0 => 0.16.0 ([#953](https://github.com/pingidentity/terraform-provider-pingone/issues/953))
* bump `github.com/patrickcping/pingone-go-sdk-v2/management` 0.44.0 => 0.45.0 ([#953](https://github.com/pingidentity/terraform-provider-pingone/issues/953))
* bump `github.com/patrickcping/pingone-go-sdk-v2/risk` 0.17.0 => 0.18.0 ([#953](https://github.com/pingidentity/terraform-provider-pingone/issues/953))
* bump `github.com/patrickcping/pingone-go-sdk-v2` 0.12.4 => 0.12.5 ([#953](https://github.com/pingidentity/terraform-provider-pingone/issues/953))

ENHANCEMENTS:

* `resource/pingone_alert_channel`: Added support for the `SUSPICIOUS_TRAFFIC` alert type. ([#953](https://github.com/pingidentity/terraform-provider-pingone/issues/953))
* `resource/pingone_risk_predictor`: Support multiple root level conditions for composite predictors. ([#952](https://github.com/pingidentity/terraform-provider-pingone/issues/952))

BUG FIXES:

* `resource/pingone_risk_predictor`: Fix "Error when calling ReadOneRiskPredictor: data failed to match schemas in oneOf(RiskPredictorCompositeCondition)" when using IP range and IP comparison composite predictors. ([#953](https://github.com/pingidentity/terraform-provider-pingone/issues/953))

## 1.2.1 (13 December 2024)

NOTES:

* bump `github.com/hashicorp/terraform-plugin-sdk/v2` 2.34.0 => 2.35.0 ([#942](https://github.com/pingidentity/terraform-provider-pingone/issues/942))
* bump `github.com/hashicorp/terraform-plugin-testing` 1.10.0 => 1.11.0 ([#943](https://github.com/pingidentity/terraform-provider-pingone/issues/943))
* bump `golang.org/x/crypto` 0.26.0 => 0.31.0 ([#946](https://github.com/pingidentity/terraform-provider-pingone/issues/946))

BUG FIXES:

* `resource/pingone_mfa_device_policy`: Fixed "unexpected new value: .fido2: was null, but now cty.ObjectVal" when `fido2` is applied, then removed from resource configuration. ([#940](https://github.com/pingidentity/terraform-provider-pingone/issues/940))

## 1.2.0 (18 November 2024)

NOTES:

* Corrected broken documentation bookmark references. ([#931](https://github.com/pingidentity/terraform-provider-pingone/issues/931))
* Upgraded go version to 1.23 to align with the go [release policy](https://go.dev/doc/devel/release#policy). ([#931](https://github.com/pingidentity/terraform-provider-pingone/issues/931))
* bump `github.com/hashicorp/terraform-plugin-framework-jsontypes` 0.1.0 => 0.2.0 ([#937](https://github.com/pingidentity/terraform-provider-pingone/issues/937))
* bump `github.com/hashicorp/terraform-plugin-framework-validators` 0.13.0 => 0.15.0 ([#937](https://github.com/pingidentity/terraform-provider-pingone/issues/937))
* bump `github.com/hashicorp/terraform-plugin-framework` 1.11.0 => 1.13.0 ([#937](https://github.com/pingidentity/terraform-provider-pingone/issues/937))
* bump `github.com/hashicorp/terraform-plugin-go` 0.23.0 => 0.25.0 ([#937](https://github.com/pingidentity/terraform-provider-pingone/issues/937))
* bump `github.com/hashicorp/terraform-plugin-mux` 0.16.0 => 0.17.0 ([#937](https://github.com/pingidentity/terraform-provider-pingone/issues/937))
* bump `github.com/patrickcping/pingone-go-sdk-v2/authorize` 0.6.0 => 0.7.0 ([#932](https://github.com/pingidentity/terraform-provider-pingone/issues/932))
* bump `github.com/patrickcping/pingone-go-sdk-v2/credentials` 0.9.0 => 0.10.0 ([#932](https://github.com/pingidentity/terraform-provider-pingone/issues/932))
* bump `github.com/patrickcping/pingone-go-sdk-v2/management` 0.43.0 => 0.44.0 ([#932](https://github.com/pingidentity/terraform-provider-pingone/issues/932))
* bump `github.com/patrickcping/pingone-go-sdk-v2/mfa` 0.20.0 => 0.21.0 ([#932](https://github.com/pingidentity/terraform-provider-pingone/issues/932))
* bump `github.com/patrickcping/pingone-go-sdk-v2/risk` 0.16.0 => 0.17.0 ([#932](https://github.com/pingidentity/terraform-provider-pingone/issues/932))
* bump `github.com/patrickcping/pingone-go-sdk-v2/verify` 0.7.0 => 0.8.0 ([#932](https://github.com/pingidentity/terraform-provider-pingone/issues/932))
* bump `github.com/patrickcping/pingone-go-sdk-v2` 0.12.3 => 0.12.4 ([#932](https://github.com/pingidentity/terraform-provider-pingone/issues/932))

ENHANCEMENTS:

* `data-source/pingone_application`: Add `session_not_on_or_after_duration` field to SAML applications. ([#934](https://github.com/pingidentity/terraform-provider-pingone/issues/934))
* `resource/pingone_application`: Add `session_not_on_or_after_duration` field to SAML applications. ([#934](https://github.com/pingidentity/terraform-provider-pingone/issues/934))
* `resource/pingone_mfa_device_policy`: Added `[email|sms|voice].otp.otp_length` field to allow admins to specify the length of the OTP displayed to users for SMS, Voice or Email delivery methods. ([#935](https://github.com/pingidentity/terraform-provider-pingone/issues/935))
* `resource/pingone_mfa_device_policy`: Added the `totp.uri_parameters` field to allow custom key:value pairs for authenticators that support `otpauth` URI parameters. ([#936](https://github.com/pingidentity/terraform-provider-pingone/issues/936))
* `resource/pingone_risk_predictor`: Added the `predictor_bot_detection.include_repeated_events_without_sdk` field to choose whether to expand the range of bot activity that PingOne Protect can detect. ([#939](https://github.com/pingidentity/terraform-provider-pingone/issues/939))
* `resource/pingone_risk_predictor`: Added the `predictor_device.should_validate_payload_signature` field to enforce requirement that the Signals SDK payload be provided as a signed JWT for suspicious device predictors. ([#938](https://github.com/pingidentity/terraform-provider-pingone/issues/938))

BUG FIXES:

* Fixed potential "Cannot find .." errors in multiple resources and data sources when many configuration items of the same type exist in an environment (fix paged results). ([#932](https://github.com/pingidentity/terraform-provider-pingone/issues/932))
* Fixed potential missing results in data sources that return multiple configuration items (fix paged results). ([#932](https://github.com/pingidentity/terraform-provider-pingone/issues/932))

## 1.1.1 (28 August 2024)

NOTES:

* `resource/pingone_environment`: Align example HCL with best practice on creating blank/empty DaVinci service environments. ([#907](https://github.com/pingidentity/terraform-provider-pingone/issues/907))
* `resource/pingone_population_default`: Suppress warning on creation where the default population for an environment cannot be found. ([#906](https://github.com/pingidentity/terraform-provider-pingone/issues/906))
* bump `github.com/hashicorp/terraform-plugin-framework-timetypes` 0.4.0 => 0.5.0 ([#908](https://github.com/pingidentity/terraform-provider-pingone/issues/908))
* bump `github.com/hashicorp/terraform-plugin-framework` 1.10.0 => 1.11.0 ([#908](https://github.com/pingidentity/terraform-provider-pingone/issues/908))
* bump `github.com/hashicorp/terraform-plugin-testing` 1.9.0 => 1.10.0 ([#908](https://github.com/pingidentity/terraform-provider-pingone/issues/908))

## 1.1.0 (05 August 2024)

NOTES:

* bump `github.com/hashicorp/terraform-plugin-framework-timetypes` 0.3.0 => 0.4.0 ([#901](https://github.com/pingidentity/terraform-provider-pingone/issues/901))
* bump `github.com/hashicorp/terraform-plugin-framework-validators` 0.12.0 => 0.13.0 ([#901](https://github.com/pingidentity/terraform-provider-pingone/issues/901))
* bump `github.com/hashicorp/terraform-plugin-framework` 1.9.0 => 1.10.0 ([#901](https://github.com/pingidentity/terraform-provider-pingone/issues/901))
* bump `github.com/hashicorp/terraform-plugin-testing` 1.8.0 => 1.9.0 ([#901](https://github.com/pingidentity/terraform-provider-pingone/issues/901))
* bump `github.com/patrickcping/pingone-go-sdk-v2/management` 0.42.0 => 0.43.0 ([#901](https://github.com/pingidentity/terraform-provider-pingone/issues/901))
* bump `github.com/patrickcping/pingone-go-sdk-v2/verify` 0.6.0 => 0.7.0 ([#901](https://github.com/pingidentity/terraform-provider-pingone/issues/901))
* bump `github.com/patrickcping/pingone-go-sdk-v2` 0.12.2 => 0.12.3 ([#901](https://github.com/pingidentity/terraform-provider-pingone/issues/901))

ENHANCEMENTS:

* `data_source/pingone_verify_policy`: Added support for the `fail_expired_id`, `provider_auto`, `provider_manual`, and `retry_attempts` properties. ([#888](https://github.com/pingidentity/terraform-provider-pingone/issues/888))
* `resource/pingone_verify_policy`: Added support for the `fail_expired_id`, `provider_auto`, `provider_manual`, and `retry_attempts` properties. ([#888](https://github.com/pingidentity/terraform-provider-pingone/issues/888))

## 1.0.0 (17 July 2024) :rocket:

BREAKING CHANGES:

* `data_source/pingone_resource`: Removed the `client_secret` attribute. Use the `pingone_resource_secret` resource or data source going forward. ([#819](https://github.com/pingidentity/terraform-provider-pingone/issues/819))
* `data-source/pingone_application`: Changed the `access_control_group_options` attribute data type. ([#682](https://github.com/pingidentity/terraform-provider-pingone/issues/682))
* `data-source/pingone_application`: Changed the `external_link_options` attribute data type. ([#682](https://github.com/pingidentity/terraform-provider-pingone/issues/682))
* `data-source/pingone_application`: Changed the `icon` attribute data type. ([#682](https://github.com/pingidentity/terraform-provider-pingone/issues/682))
* `data-source/pingone_application`: Changed the `oidc_options.certificate_based_authentication` attribute data type. ([#682](https://github.com/pingidentity/terraform-provider-pingone/issues/682))
* `data-source/pingone_application`: Changed the `oidc_options.cors_settings` attribute data type. ([#682](https://github.com/pingidentity/terraform-provider-pingone/issues/682))
* `data-source/pingone_application`: Changed the `oidc_options.mobile_app.integrity_detection.cache_duration` attribute data type. ([#682](https://github.com/pingidentity/terraform-provider-pingone/issues/682))
* `data-source/pingone_application`: Changed the `oidc_options.mobile_app.integrity_detection.google_play` attribute data type. ([#682](https://github.com/pingidentity/terraform-provider-pingone/issues/682))
* `data-source/pingone_application`: Changed the `oidc_options.mobile_app.integrity_detection` attribute data type. ([#682](https://github.com/pingidentity/terraform-provider-pingone/issues/682))
* `data-source/pingone_application`: Changed the `oidc_options.mobile_app` attribute data type. ([#682](https://github.com/pingidentity/terraform-provider-pingone/issues/682))
* `data-source/pingone_application`: Changed the `oidc_options` attribute data type. ([#682](https://github.com/pingidentity/terraform-provider-pingone/issues/682))
* `data-source/pingone_application`: Changed the `saml_options.cors_settings` attribute data type. ([#682](https://github.com/pingidentity/terraform-provider-pingone/issues/682))
* `data-source/pingone_application`: Changed the `saml_options.idp_signing_key` attribute data type. ([#682](https://github.com/pingidentity/terraform-provider-pingone/issues/682))
* `data-source/pingone_application`: Changed the `saml_options.sp_verification` ([#682](https://github.com/pingidentity/terraform-provider-pingone/issues/682))
* `data-source/pingone_application`: Changed the `saml_options` attribute data type. ([#682](https://github.com/pingidentity/terraform-provider-pingone/issues/682))
* `data-source/pingone_application`: Removal of deprecated attribute `saml_options.sp_verification_certificate_ids`.  Use the `saml_options.sp_verification.certificate_ids` attribute going forward. ([#681](https://github.com/pingidentity/terraform-provider-pingone/issues/681))
* `data-source/pingone_application`: Removed `oidc_options.client_secret`.  Use the `pingone_application_secret` resource or data source going forward. ([#781](https://github.com/pingidentity/terraform-provider-pingone/issues/781))
* `data-source/pingone_application`: Renamed `oidc_options.allow_wildcards_in_redirect_uris` to `oidc_options.allow_wildcard_in_redirect_uris` to align with the API. ([#887](https://github.com/pingidentity/terraform-provider-pingone/issues/887))
* `data-source/pingone_application`: Renamed `oidc_options.token_endpoint_authn_method` to `oidc_options.token_endpoint_auth_method` to align with the API. ([#887](https://github.com/pingidentity/terraform-provider-pingone/issues/887))
* `data-source/pingone_environment`: Changed the `service.bookmark` parameter data type and renamed to `services.bookmarks`. ([#665](https://github.com/pingidentity/terraform-provider-pingone/issues/665))
* `data-source/pingone_environment`: Changed the `service` parameter data type and renamed to `services`. ([#665](https://github.com/pingidentity/terraform-provider-pingone/issues/665))
* `data-source/pingone_flow_policies`: Changed the `data_filter` parameter data type and renamed to `data_filters`. ([#664](https://github.com/pingidentity/terraform-provider-pingone/issues/664))
* `data-source/pingone_flow_policy`: Changed the `davinci_application` and `trigger` computed attribute data types from list block to single object types. ([#795](https://github.com/pingidentity/terraform-provider-pingone/issues/795))
* `data-source/pingone_gateway`: Changed the `user_type` data type from a set of objects to a map of objects and renamed to `user_types`. ([#798](https://github.com/pingidentity/terraform-provider-pingone/issues/798))
* `data-source/pingone_gateway`: Renamed `radius_client` to `radius_clients` and changed data type from block set to set of objects. ([#798](https://github.com/pingidentity/terraform-provider-pingone/issues/798))
* `data-source/pingone_gateway`: Renamed `user_type.push_password_changes_to_ldap` to `user_types.allow_password_changes`. ([#798](https://github.com/pingidentity/terraform-provider-pingone/issues/798))
* `data-source/pingone_gateway`: Renamed `user_type.user_migation.attribute_mapping` to `user_types.new_user_lookup.attribute_mappings` and changed data type from block set to set of objects. ([#798](https://github.com/pingidentity/terraform-provider-pingone/issues/798))
* `data-source/pingone_gateway`: Renamed `user_type.user_migation` to `user_types.new_user_lookup` and changed data type from block set to single object. ([#798](https://github.com/pingidentity/terraform-provider-pingone/issues/798))
* `data-source/pingone_gateway`: Renamed `user_type.user_migration.lookup_filter_pattern` to `user_types.new_user_lookup.ldap_filter_pattern`. ([#798](https://github.com/pingidentity/terraform-provider-pingone/issues/798))
* `data-source/pingone_groups`: Changed the `data_filter` parameter data type and renamed to `data_filters`. ([#677](https://github.com/pingidentity/terraform-provider-pingone/issues/677))
* `data-source/pingone_licenses`: Changed the `data_filter` parameter data type and renamed to `data_filters`. ([#730](https://github.com/pingidentity/terraform-provider-pingone/issues/730))
* `data-source/pingone_organization`: Removal of deprecated platform URL computed attributes.  Consider using the [PingOne Utilities module](https://registry.terraform.io/modules/pingidentity/utils/pingone/latest) going forward. ([#628](https://github.com/pingidentity/terraform-provider-pingone/issues/628))
* `data-source/pingone_password_policy`: Moved `password_age.max` to `password_age_max`. ([#801](https://github.com/pingidentity/terraform-provider-pingone/issues/801))
* `data-source/pingone_password_policy`: Moved `password_age.min` to `password_age_min`. ([#801](https://github.com/pingidentity/terraform-provider-pingone/issues/801))
* `data-source/pingone_password_policy`: Removed ineffectual `bypass_policy`. ([#801](https://github.com/pingidentity/terraform-provider-pingone/issues/801))
* `data-source/pingone_password_policy`: Renamed `account_lockout` to `lockout`. ([#801](https://github.com/pingidentity/terraform-provider-pingone/issues/801))
* `data-source/pingone_password_policy`: Renamed `environment_default` to `default`. ([#801](https://github.com/pingidentity/terraform-provider-pingone/issues/801))
* `data-source/pingone_password_policy`: Renamed `exclude_commonly_used_passwords` to `excludes_commonly_used_passwords`. ([#801](https://github.com/pingidentity/terraform-provider-pingone/issues/801))
* `data-source/pingone_password_policy`: Renamed `exclude_profile_data` to `excludes_profile_data`. ([#801](https://github.com/pingidentity/terraform-provider-pingone/issues/801))
* `data-source/pingone_password_policy`: Renamed `password_history` to `history`. ([#801](https://github.com/pingidentity/terraform-provider-pingone/issues/801))
* `data-source/pingone_password_policy`: Renamed `password_length` to `length`. ([#801](https://github.com/pingidentity/terraform-provider-pingone/issues/801))
* `data-source/pingone_populations`: Changed the `data_filter` parameter data type and renamed to `data_filters`. ([#664](https://github.com/pingidentity/terraform-provider-pingone/issues/664))
* `data-source/pingone_resource_scope`: Existing `resource_id` field made read only.  Use `resource_type` and `custom_resource_id` instead. ([#863](https://github.com/pingidentity/terraform-provider-pingone/issues/863))
* `data-source/pingone_resource_scope`: New `resource_type` field is a required field, and new `custom_resource_id` field is an optional field.  The combination of these fields ensure the correct resource is selected without encountering issue. ([#863](https://github.com/pingidentity/terraform-provider-pingone/issues/863))
* `data-source/pingone_trusted_email_domain_dkim`: Removed unnecessary `id` attribute. ([#802](https://github.com/pingidentity/terraform-provider-pingone/issues/802))
* `data-source/pingone_trusted_email_domain_dkim`: Renamed `region.token` to `regions.tokens`. ([#802](https://github.com/pingidentity/terraform-provider-pingone/issues/802))
* `data-source/pingone_trusted_email_domain_dkim`: Renamed `region` to `regions`. ([#802](https://github.com/pingidentity/terraform-provider-pingone/issues/802))
* `data-source/pingone_trusted_email_domain_ownership`: Renamed `region` to `regions`. ([#803](https://github.com/pingidentity/terraform-provider-pingone/issues/803))
* `data-source/pingone_user`: Removed the deprecated `status` attribute.  Use the `enabled` parameter going forward. ([#647](https://github.com/pingidentity/terraform-provider-pingone/issues/647))
* `data-source/pingone_users`: Changed the `data_filter` parameter data type and renamed to `data_filters`. ([#664](https://github.com/pingidentity/terraform-provider-pingone/issues/664))
* `pingone_environment`: The `region` parameter's values now aligns with the API request/response payload values.  See the Upgrade Guide for details. ([#828](https://github.com/pingidentity/terraform-provider-pingone/issues/828))
* `resource/pingone_application_resource_grant`: Changed parameters `resource_id` and `scopes` to be read-only.  Use the `resource_name` and `scope_names` parameters going forward. ([#657](https://github.com/pingidentity/terraform-provider-pingone/issues/657))
* `resource/pingone_application_resource_grant`: Changed parameters `resource_name` and `scope_names` to be required parameters. ([#657](https://github.com/pingidentity/terraform-provider-pingone/issues/657))
* `resource/pingone_application_resource_grant`: Existing `resource_name` field removed.  Use `resource_type` and `custom_resource_id` instead. ([#863](https://github.com/pingidentity/terraform-provider-pingone/issues/863))
* `resource/pingone_application_resource_grant`: Existing `scope_names` field removed and existing `scopes` field is now a required field.  Use `scopes` to define the list of scopes for the grant instead. ([#863](https://github.com/pingidentity/terraform-provider-pingone/issues/863))
* `resource/pingone_application_resource_grant`: New `resource_type` field is a required field, and new `custom_resource_id` field is an optional field.  The combination of these fields ensure the correct resource is selected without encountering issue. ([#863](https://github.com/pingidentity/terraform-provider-pingone/issues/863))
* `resource/pingone_application`: Changed the `access_control_group_options` attribute data type. ([#683](https://github.com/pingidentity/terraform-provider-pingone/issues/683))
* `resource/pingone_application`: Changed the `external_link_options` attribute data type. ([#683](https://github.com/pingidentity/terraform-provider-pingone/issues/683))
* `resource/pingone_application`: Changed the `icon` attribute data type. ([#683](https://github.com/pingidentity/terraform-provider-pingone/issues/683))
* `resource/pingone_application`: Changed the `oidc_options.certificate_based_authentication` attribute data type. ([#683](https://github.com/pingidentity/terraform-provider-pingone/issues/683))
* `resource/pingone_application`: Changed the `oidc_options.cors_settings` attribute data type. ([#683](https://github.com/pingidentity/terraform-provider-pingone/issues/683))
* `resource/pingone_application`: Changed the `oidc_options.mobile_app.integrity_detection.cache_duration` attribute data type. ([#683](https://github.com/pingidentity/terraform-provider-pingone/issues/683))
* `resource/pingone_application`: Changed the `oidc_options.mobile_app.integrity_detection.google_play` attribute data type. ([#683](https://github.com/pingidentity/terraform-provider-pingone/issues/683))
* `resource/pingone_application`: Changed the `oidc_options.mobile_app.integrity_detection` attribute data type. ([#683](https://github.com/pingidentity/terraform-provider-pingone/issues/683))
* `resource/pingone_application`: Changed the `oidc_options.mobile_app` attribute data type. ([#683](https://github.com/pingidentity/terraform-provider-pingone/issues/683))
* `resource/pingone_application`: Changed the `oidc_options` attribute data type. ([#683](https://github.com/pingidentity/terraform-provider-pingone/issues/683))
* `resource/pingone_application`: Changed the `saml_options.cors_settings` attribute data type. ([#683](https://github.com/pingidentity/terraform-provider-pingone/issues/683))
* `resource/pingone_application`: Changed the `saml_options.idp_signing_key` attribute data type and made required in the schema when defining SAML applications. ([#683](https://github.com/pingidentity/terraform-provider-pingone/issues/683))
* `resource/pingone_application`: Changed the `saml_options.sp_verification` ([#683](https://github.com/pingidentity/terraform-provider-pingone/issues/683))
* `resource/pingone_application`: Changed the `saml_options` attribute data type. ([#683](https://github.com/pingidentity/terraform-provider-pingone/issues/683))
* `resource/pingone_application`: Removal of deprecated parameter `saml_options.idp_signing_key_id`.  Use the `saml_options.idp_signing_key` parameter going forward. ([#656](https://github.com/pingidentity/terraform-provider-pingone/issues/656))
* `resource/pingone_application`: Removal of deprecated parameter `saml_options.sp_verification_certificate_ids`.  Use the `saml_options.sp_verification.certificate_ids` parameter going forward. ([#681](https://github.com/pingidentity/terraform-provider-pingone/issues/681))
* `resource/pingone_application`: Removal of deprecated parameters `oidc_options.bundle_id`, `oidc_options.package_name`.  Use the `oidc_options.mobile_app.bundle_id` and `oidc_options.mobile_app.package_name` parameters going forward. ([#656](https://github.com/pingidentity/terraform-provider-pingone/issues/656))
* `resource/pingone_application`: Removed `oidc_options.client_secret`.  Use the `pingone_application_secret` resource or data source going forward. ([#781](https://github.com/pingidentity/terraform-provider-pingone/issues/781))
* `resource/pingone_application`: Renamed `oidc_options.allow_wildcards_in_redirect_uris` to `oidc_options.allow_wildcard_in_redirect_uris` to align with the API. ([#887](https://github.com/pingidentity/terraform-provider-pingone/issues/887))
* `resource/pingone_application`: Renamed `oidc_options.token_endpoint_authn_method` to `oidc_options.token_endpoint_auth_method` to align with the API. ([#887](https://github.com/pingidentity/terraform-provider-pingone/issues/887))
* `resource/pingone_branding_settings`: Changed the `logo_image` parameter data type. ([#661](https://github.com/pingidentity/terraform-provider-pingone/issues/661))
* `resource/pingone_branding_theme`: Changed the `logo` and `background_image` parameters data type. ([#661](https://github.com/pingidentity/terraform-provider-pingone/issues/661))
* `resource/pingone_custom_domain_verify`: Changed the `timeouts` data type. ([#786](https://github.com/pingidentity/terraform-provider-pingone/issues/786))
* `resource/pingone_environment`: Changed the `service.bookmark` parameter data type and renamed to `services.bookmarks`. ([#665](https://github.com/pingidentity/terraform-provider-pingone/issues/665))
* `resource/pingone_environment`: Changed the `service` parameter data type, renamed to `services` and made a required parameter. ([#665](https://github.com/pingidentity/terraform-provider-pingone/issues/665))
* `resource/pingone_environment`: Moved `service.type` to `services.type` and made a required parameter. ([#665](https://github.com/pingidentity/terraform-provider-pingone/issues/665))
* `resource/pingone_environment`: Removal of deprecated parameter `default_population` and computed attribute `default_population_id`. ([#629](https://github.com/pingidentity/terraform-provider-pingone/issues/629))
* `resource/pingone_environment`: Removed `timeouts` parameter block. ([#643](https://github.com/pingidentity/terraform-provider-pingone/issues/643))
* `resource/pingone_gateway`: Changed the `user_type` data type from a set of objects to a map of objects and renamed to `user_types`. ([#798](https://github.com/pingidentity/terraform-provider-pingone/issues/798))
* `resource/pingone_gateway`: Renamed `radius_client` to `radius_clients` and changed data type from block set to set of objects. ([#798](https://github.com/pingidentity/terraform-provider-pingone/issues/798))
* `resource/pingone_gateway`: Renamed `user_type.push_password_changes_to_ldap` to `user_types.allow_password_changes`. ([#798](https://github.com/pingidentity/terraform-provider-pingone/issues/798))
* `resource/pingone_gateway`: Renamed `user_type.user_migation.attribute_mapping` to `user_types.new_user_lookup.attribute_mappings` and changed data type from block set to set of objects. ([#798](https://github.com/pingidentity/terraform-provider-pingone/issues/798))
* `resource/pingone_gateway`: Renamed `user_type.user_migation` to `user_types.new_user_lookup` and changed data type from block set to single object. ([#798](https://github.com/pingidentity/terraform-provider-pingone/issues/798))
* `resource/pingone_gateway`: Renamed `user_type.user_migration.lookup_filter_pattern` to `user_types.new_user_lookup.ldap_filter_pattern`. ([#798](https://github.com/pingidentity/terraform-provider-pingone/issues/798))
* `resource/pingone_identity_provider`: Changed the `facebook`, `google`, `linkedin`, `yahoo`, `amazon`, `twitter`, `apple`, `paypal`, `microsoft`, `github`, `openid_connect` and `saml` parameter data types. ([#662](https://github.com/pingidentity/terraform-provider-pingone/issues/662))
* `resource/pingone_identity_provider`: Changed the `icon` and `login_button_icon` parameters data type. ([#661](https://github.com/pingidentity/terraform-provider-pingone/issues/661))
* `resource/pingone_identity_provider`: Replaced `saml.idp_verification_certificate_ids` with `saml.idp_verification.certificates.*.id`. ([#830](https://github.com/pingidentity/terraform-provider-pingone/issues/830))
* `resource/pingone_identity_provider`: Replaced `saml.sp_signing_key_id` with `saml.sp_signing.key.id`. ([#830](https://github.com/pingidentity/terraform-provider-pingone/issues/830))
* `resource/pingone_mfa_application_push_credential`: Changed the `apns`, `fcm` and `hms` data types. ([#644](https://github.com/pingidentity/terraform-provider-pingone/issues/644))
* `resource/pingone_mfa_application_push_credential`: Removal of deprecated parameter `fcm.key`.  Use the `fcm.google_service_account_credentials` parameter going forward. ([#630](https://github.com/pingidentity/terraform-provider-pingone/issues/630))
* `resource/pingone_mfa_device_policy`: Changed the data types of `email`, `voice`, `sms`, `totp`, `mobile` and `fido2` from list of objects to single object type. ([#809](https://github.com/pingidentity/terraform-provider-pingone/issues/809))
* `resource/pingone_mfa_device_policy`: Moved `{email,sms,voice,mobile,totp}.otp_failure_cooldown_duration` to `{email,sms,voice,mobile,totp}.otp.failure.cool_down.duration`. ([#809](https://github.com/pingidentity/terraform-provider-pingone/issues/809))
* `resource/pingone_mfa_device_policy`: Moved `{email,sms,voice,mobile,totp}.otp_failure_cooldown_timeunit` to `{email,sms,voice,mobile,totp}.otp.failure.cool_down.time_unit`. ([#809](https://github.com/pingidentity/terraform-provider-pingone/issues/809))
* `resource/pingone_mfa_device_policy`: Moved `{email,sms,voice,mobile,totp}.otp_failure_count` to `{email,sms,voice,mobile,totp}.otp.failure.count`. ([#809](https://github.com/pingidentity/terraform-provider-pingone/issues/809))
* `resource/pingone_mfa_device_policy`: Moved `{email,sms,voice}.otp_lifetime_duration` to `{email,sms,voice}.otp.lifetime_duration.duration`. ([#809](https://github.com/pingidentity/terraform-provider-pingone/issues/809))
* `resource/pingone_mfa_device_policy`: Moved `{email,sms,voice}.otp_lifetime_duration` to `{email,sms,voice}.otp.lifetime.duration`. ([#809](https://github.com/pingidentity/terraform-provider-pingone/issues/809))
* `resource/pingone_mfa_device_policy`: Moved `{email,sms,voice}.otp_lifetime_timeunit` to `{email,sms,voice}.otp.lifetime.time_unit`. ([#809](https://github.com/pingidentity/terraform-provider-pingone/issues/809))
* `resource/pingone_mfa_device_policy`: Moved `device_selection` to `authentication.device_selection`. ([#809](https://github.com/pingidentity/terraform-provider-pingone/issues/809))
* `resource/pingone_mfa_device_policy`: Moved `mobile.application.auto_enrollment_enabled` to `mobile.applications.*.auto_enrollment.enabled`. ([#809](https://github.com/pingidentity/terraform-provider-pingone/issues/809))
* `resource/pingone_mfa_device_policy`: Moved `mobile.application.device_authorization_enabled` to `mobile.applications.*.device_authorization.enabled`. ([#809](https://github.com/pingidentity/terraform-provider-pingone/issues/809))
* `resource/pingone_mfa_device_policy`: Moved `mobile.application.device_authorization_extra_verification` to `mobile.applications.*.device_authorization.extra_verification`. ([#809](https://github.com/pingidentity/terraform-provider-pingone/issues/809))
* `resource/pingone_mfa_device_policy`: Moved `mobile.application.otp_enabled` to `mobile.applications.*.otp.enabled`. ([#809](https://github.com/pingidentity/terraform-provider-pingone/issues/809))
* `resource/pingone_mfa_device_policy`: Moved `mobile.application.pairing_key_lifetime_duration` to `mobile.applications.*.pairing_key_lifetime.duration`. ([#809](https://github.com/pingidentity/terraform-provider-pingone/issues/809))
* `resource/pingone_mfa_device_policy`: Moved `mobile.application.pairing_key_lifetime_timeunit` to `mobile.applications.*.pairing_key_lifetime.time_unit`. ([#809](https://github.com/pingidentity/terraform-provider-pingone/issues/809))
* `resource/pingone_mfa_device_policy`: Moved `mobile.application.push_enabled` to `mobile.applications.*.push.enabled`. ([#809](https://github.com/pingidentity/terraform-provider-pingone/issues/809))
* `resource/pingone_mfa_device_policy`: Moved `mobile.application.push_limit_count` to `mobile.applications.*.push_limit.count`. ([#809](https://github.com/pingidentity/terraform-provider-pingone/issues/809))
* `resource/pingone_mfa_device_policy`: Moved `mobile.application.push_limit_lock_duration` to `mobile.applications.*.push_limit.lock_duration.duration`. ([#809](https://github.com/pingidentity/terraform-provider-pingone/issues/809))
* `resource/pingone_mfa_device_policy`: Moved `mobile.application.push_limit_lock_timeunit` to `mobile.applications.*.push_limit.lock_duration.time_unit`. ([#809](https://github.com/pingidentity/terraform-provider-pingone/issues/809))
* `resource/pingone_mfa_device_policy`: Moved `mobile.application.push_limit_time_period_duration` to `mobile.applications.*.push_limit.time_period.duration`. ([#809](https://github.com/pingidentity/terraform-provider-pingone/issues/809))
* `resource/pingone_mfa_device_policy`: Moved `mobile.application.push_limit_time_period_timeunit` to `mobile.applications.*.push_limit.time_period.time_unit`. ([#809](https://github.com/pingidentity/terraform-provider-pingone/issues/809))
* `resource/pingone_mfa_device_policy`: Moved `mobile.application.push_timeout_duration` to `mobile.applications.*.push_timeout.duration`. ([#809](https://github.com/pingidentity/terraform-provider-pingone/issues/809))
* `resource/pingone_mfa_device_policy`: Moved `mobile.application.push_timeout_timeunit` to `mobile.applications.*.push_timeout_time_unit`. ([#809](https://github.com/pingidentity/terraform-provider-pingone/issues/809))
* `resource/pingone_mfa_device_policy`: Renamed `mobile.application` to `mobile.applications` and changed the data type to a map of objects. ([#809](https://github.com/pingidentity/terraform-provider-pingone/issues/809))
* `resource/pingone_mfa_fido_policy`: Removal of deprecated resource. ([#625](https://github.com/pingidentity/terraform-provider-pingone/issues/625))
* `resource/pingone_mfa_policies`: Removal of deprecated resource. ([#626](https://github.com/pingidentity/terraform-provider-pingone/issues/626))
* `resource/pingone_mfa_policy`: Removal of deprecated parameters `platform` and `security_key`.  Use the `fido2` parameter going forward. ([#627](https://github.com/pingidentity/terraform-provider-pingone/issues/627))
* `resource/pingone_mfa_settings`: Changed the `pairing` and `lockout` data types. ([#797](https://github.com/pingidentity/terraform-provider-pingone/issues/797))
* `resource/pingone_mfa_settings`: Removed `phone_extensions_enabled` and moved into nested attribute object.  Use `phone_extensions.enabled` going forward. ([#797](https://github.com/pingidentity/terraform-provider-pingone/issues/797))
* `resource/pingone_mfa_settings`: Removed deprecated parameter `authentication`. ([#645](https://github.com/pingidentity/terraform-provider-pingone/issues/645))
* `resource/pingone_notification_policy`: Changed the `quota` data type. ([#789](https://github.com/pingidentity/terraform-provider-pingone/issues/789))
* `resource/pingone_notification_settings_email`: Changed the `from` and `reply_to` data types. ([#796](https://github.com/pingidentity/terraform-provider-pingone/issues/796))
* `resource/pingone_password_policy`: Moved `password_age.max` to `password_age_max`. ([#801](https://github.com/pingidentity/terraform-provider-pingone/issues/801))
* `resource/pingone_password_policy`: Moved `password_age.min` to `password_age_min`. ([#801](https://github.com/pingidentity/terraform-provider-pingone/issues/801))
* `resource/pingone_password_policy`: Removed ineffectual `bypass_policy`. ([#801](https://github.com/pingidentity/terraform-provider-pingone/issues/801))
* `resource/pingone_password_policy`: Renamed `account_lockout` to `lockout`. ([#801](https://github.com/pingidentity/terraform-provider-pingone/issues/801))
* `resource/pingone_password_policy`: Renamed `environment_default` to `default`. ([#801](https://github.com/pingidentity/terraform-provider-pingone/issues/801))
* `resource/pingone_password_policy`: Renamed `exclude_commonly_used_passwords` to `excludes_commonly_used_passwords`. ([#801](https://github.com/pingidentity/terraform-provider-pingone/issues/801))
* `resource/pingone_password_policy`: Renamed `exclude_profile_data` to `excludes_profile_data`. ([#801](https://github.com/pingidentity/terraform-provider-pingone/issues/801))
* `resource/pingone_password_policy`: Renamed `password_history` to `history`. ([#801](https://github.com/pingidentity/terraform-provider-pingone/issues/801))
* `resource/pingone_password_policy`: Renamed `password_length` to `length`. ([#801](https://github.com/pingidentity/terraform-provider-pingone/issues/801))
* `resource/pingone_resource_attribute`: Changed parameter `resource_id` to be read-only.  Use the `resource_name` parameter going forward. ([#658](https://github.com/pingidentity/terraform-provider-pingone/issues/658))
* `resource/pingone_resource_attribute`: Changed parameter `resource_name` to be a required parameter. ([#658](https://github.com/pingidentity/terraform-provider-pingone/issues/658))
* `resource/pingone_resource_attribute`: Existing `resource_name` field removed.  Use `resource_type` and `custom_resource_id` instead. ([#863](https://github.com/pingidentity/terraform-provider-pingone/issues/863))
* `resource/pingone_resource_attribute`: New `resource_type` field is a required field, and new `custom_resource_id` field is an optional field.  The combination of these fields ensure the correct resource is selected without encountering issue. ([#863](https://github.com/pingidentity/terraform-provider-pingone/issues/863))
* `resource/pingone_resource`: Removed the `client_secret` attribute. Use the `pingone_resource_secret` resource or data source going forward. ([#819](https://github.com/pingidentity/terraform-provider-pingone/issues/819))
* `resource/pingone_schema_attribute`: Changed parameter `schema_id` to be read-only.  Use the optional `schema_name` parameter going forward. ([#660](https://github.com/pingidentity/terraform-provider-pingone/issues/660))
* `resource/pingone_user`: Removed the deprecated `status` parameter.  Use the `enabled` parameter going forward. ([#647](https://github.com/pingidentity/terraform-provider-pingone/issues/647))
* `resource/pingone_webhook`: Changed the `filter_options` parameter data type. ([#663](https://github.com/pingidentity/terraform-provider-pingone/issues/663))
* Removed the provider parameter `force_delete_production_type`.  Use the `global_options.environment.production_type_force_delete` parameter going forward. ([#787](https://github.com/pingidentity/terraform-provider-pingone/issues/787))
* Renamed the `pingone_mfa_policies` data source to `pingone_mfa_device_policies`. ([#788](https://github.com/pingidentity/terraform-provider-pingone/issues/788))
* Renamed the `pingone_mfa_policy` resource to `pingone_mfa_device_policy`. ([#788](https://github.com/pingidentity/terraform-provider-pingone/issues/788))
* Renamed the `pingone_role_assignment_user` resource to `pingone_user_role_assignment`. ([#843](https://github.com/pingidentity/terraform-provider-pingone/issues/843))
* Replaced the `region` parameter (and `PINGONE_REGION` environment variable) with `region_code` (defaulted with the `PINGONE_REGION_CODE` environment variable).  See the Upgrade Guide for details. ([#828](https://github.com/pingidentity/terraform-provider-pingone/issues/828))

NOTES:

* `data-source/pingone_license`: Migrated to plugin framework. ([#894](https://github.com/pingidentity/terraform-provider-pingone/issues/894))
* `data-source/pingone_password_policy`: Migrated to plugin framework. ([#801](https://github.com/pingidentity/terraform-provider-pingone/issues/801))
* `data-source/pingone_trusted_email_domain_dkim`: Migrated to plugin framework. ([#802](https://github.com/pingidentity/terraform-provider-pingone/issues/802))
* `data-source/pingone_trusted_email_domain_ownership`: Migrated to plugin framework. ([#803](https://github.com/pingidentity/terraform-provider-pingone/issues/803))
* `resource/pingone_application`: Changed the `oidc_options.mobile_app.integrity_detection.excluded_platforms` attribute data type. ([#683](https://github.com/pingidentity/terraform-provider-pingone/issues/683))
* `resource/pingone_application`: Migrated to plugin framework. ([#683](https://github.com/pingidentity/terraform-provider-pingone/issues/683))
* `resource/pingone_gateway`: Migrated to plugin framework. ([#798](https://github.com/pingidentity/terraform-provider-pingone/issues/798))
* `resource/pingone_image`: Migrated to plugin framework. ([#509](https://github.com/pingidentity/terraform-provider-pingone/issues/509))
* `resource/pingone_mfa_device_policy`: Migrated to plugin framework. ([#809](https://github.com/pingidentity/terraform-provider-pingone/issues/809))
* `resource/pingone_mfa_settings`: Migrated to plugin framework. ([#797](https://github.com/pingidentity/terraform-provider-pingone/issues/797))
* `resource/pingone_notification_template_content`: Migrated to plugin framework. ([#837](https://github.com/pingidentity/terraform-provider-pingone/issues/837))
* `resource/pingone_password_policy`: Migrated to plugin framework. ([#801](https://github.com/pingidentity/terraform-provider-pingone/issues/801))
* `resource/pingone_resource`: Migrated to plugin framework. ([#819](https://github.com/pingidentity/terraform-provider-pingone/issues/819))
* bump `github.com/patrickcping/pingone-go-sdk-v2/authorize` 0.5.0 => 0.6.0 ([#871](https://github.com/pingidentity/terraform-provider-pingone/issues/871))
* bump `github.com/patrickcping/pingone-go-sdk-v2/credentials` 0.6.2 => 0.7.0 ([#798](https://github.com/pingidentity/terraform-provider-pingone/issues/798))
* bump `github.com/patrickcping/pingone-go-sdk-v2/credentials` 0.8.0 => 0.9.0 ([#871](https://github.com/pingidentity/terraform-provider-pingone/issues/871))
* bump `github.com/patrickcping/pingone-go-sdk-v2/management` 0.38.0 => 0.39.0 ([#798](https://github.com/pingidentity/terraform-provider-pingone/issues/798))
* bump `github.com/patrickcping/pingone-go-sdk-v2/management` 0.41.0 => 0.42.0 ([#871](https://github.com/pingidentity/terraform-provider-pingone/issues/871))
* bump `github.com/patrickcping/pingone-go-sdk-v2/mfa` 0.19.0 => 0.20.0 ([#871](https://github.com/pingidentity/terraform-provider-pingone/issues/871))
* bump `github.com/patrickcping/pingone-go-sdk-v2/risk` 0.15.1 => 0.16.0 ([#871](https://github.com/pingidentity/terraform-provider-pingone/issues/871))
* bump `github.com/patrickcping/pingone-go-sdk-v2/verify` 0.5.0 => 0.6.0 ([#871](https://github.com/pingidentity/terraform-provider-pingone/issues/871))
* bump `github.com/patrickcping/pingone-go-sdk-v2` 0.11.8 => 0.11.9 ([#798](https://github.com/pingidentity/terraform-provider-pingone/issues/798))
* bump `github.com/patrickcping/pingone-go-sdk-v2` 0.12.1 => 0.12.2 ([#871](https://github.com/pingidentity/terraform-provider-pingone/issues/871))
* Changed date-time fields to use custom RFC3339 data type. ([#784](https://github.com/pingidentity/terraform-provider-pingone/issues/784))
* Changed JSON fields to use custom JSON data type. ([#785](https://github.com/pingidentity/terraform-provider-pingone/issues/785))
* Reformat API error responses to be clearer to read. ([#864](https://github.com/pingidentity/terraform-provider-pingone/issues/864))

FEATURES:

* **New Data Source:** `pingone_application_secret` ([#781](https://github.com/pingidentity/terraform-provider-pingone/issues/781))
* **New Data Source:** `pingone_resource_secret` ([#819](https://github.com/pingidentity/terraform-provider-pingone/issues/819))
* **New Resource:** `pingone_alert_channel` ([#848](https://github.com/pingidentity/terraform-provider-pingone/issues/848))
* **New Resource:** `pingone_application_resource_permission` ([#820](https://github.com/pingidentity/terraform-provider-pingone/issues/820))
* **New Resource:** `pingone_application_resource` ([#818](https://github.com/pingidentity/terraform-provider-pingone/issues/818))
* **New Resource:** `pingone_authorize_api_service_operation` ([#825](https://github.com/pingidentity/terraform-provider-pingone/issues/825))
* **New Resource:** `pingone_authorize_api_service` ([#824](https://github.com/pingidentity/terraform-provider-pingone/issues/824))
* **New Resource:** `pingone_authorize_application_role_permission` ([#821](https://github.com/pingidentity/terraform-provider-pingone/issues/821))
* **New Resource:** `pingone_authorize_application_role` ([#817](https://github.com/pingidentity/terraform-provider-pingone/issues/817))
* **New Resource:** `pingone_population_default_identity_provider` ([#831](https://github.com/pingidentity/terraform-provider-pingone/issues/831))
* **New Resource:** `pingone_resource_secret` ([#819](https://github.com/pingidentity/terraform-provider-pingone/issues/819))
* **New Resource:** `pingone_user_application_role_assignment` ([#822](https://github.com/pingidentity/terraform-provider-pingone/issues/822))

ENHANCEMENTS:

* `data_source/pingone_resource`: Added support for the `application_permissions_settings` attribute. ([#819](https://github.com/pingidentity/terraform-provider-pingone/issues/819))
* `data-source/pingone_application`: Added support for the `DEVICE_CODE` grant type for OIDC applications. ([#834](https://github.com/pingidentity/terraform-provider-pingone/issues/834))
* `data-source/pingone_gateway`: Added `follow_referrals` LDAP gateway parameter. ([#798](https://github.com/pingidentity/terraform-provider-pingone/issues/798))
* `data-source/pingone_gateway`: Added `new_user_lookup.update_user_on_successful_authentication` LDAP gateway parameter. ([#798](https://github.com/pingidentity/terraform-provider-pingone/issues/798))
* `data-source/pingone_group`: Added `custom_data` field to be able to read custom JSON data attached to a group. ([#850](https://github.com/pingidentity/terraform-provider-pingone/issues/850))
* `resource/pingone_application_secret`: Support for handling previous secrets for application secret rotation. ([#781](https://github.com/pingidentity/terraform-provider-pingone/issues/781))
* `resource/pingone_application`: Added support for the `DEVICE_CODE` grant type for OIDC applications. ([#834](https://github.com/pingidentity/terraform-provider-pingone/issues/834))
* `resource/pingone_environment`: Support the creation of trials enabled environments. ([#849](https://github.com/pingidentity/terraform-provider-pingone/issues/849))
* `resource/pingone_gateway`: Added `follow_referrals` LDAP gateway parameter. ([#798](https://github.com/pingidentity/terraform-provider-pingone/issues/798))
* `resource/pingone_gateway`: Added `new_user_lookup.update_user_on_successful_authentication` LDAP gateway parameter. ([#798](https://github.com/pingidentity/terraform-provider-pingone/issues/798))
* `resource/pingone_group`: Added `custom_data` field to be able to append custom JSON data to a group. ([#850](https://github.com/pingidentity/terraform-provider-pingone/issues/850))
* `resource/pingone_identity_provider`: Added ability to set the SP signing key algorithm. ([#830](https://github.com/pingidentity/terraform-provider-pingone/issues/830))
* `resource/pingone_identity_provider`: Added the `pkce_method` property for OIDC Identity Providers. ([#829](https://github.com/pingidentity/terraform-provider-pingone/issues/829))
* `resource/pingone_mfa_device_policy`: Added `prompt_for_nickname_on_pairing` for each device method, which provides a prompt for users to provide nicknames for devices on pairing. ([#809](https://github.com/pingidentity/terraform-provider-pingone/issues/809))
* `resource/pingone_mfa_device_policy`: Added the `default` field to track (in state) whether the policy is the default for the environment. ([#844](https://github.com/pingidentity/terraform-provider-pingone/issues/844))
* `resource/pingone_mfa_settings`: Added `users.mfa_enabled` that, when set to `true`, will enable MFA by default for new users. ([#797](https://github.com/pingidentity/terraform-provider-pingone/issues/797))
* `resource/pingone_resource`: Added support for the `application_permissions_settings` property to be able to add permissions to access tokens. ([#819](https://github.com/pingidentity/terraform-provider-pingone/issues/819))
* `resource/pingone_risk_predictor`: Added support for the `ADVERSARY_IN_THE_MIDDLE` and `EMAIL_REPUTATION` predictors. ([#835](https://github.com/pingidentity/terraform-provider-pingone/issues/835))
* `resource/pingone_schema_attribute`: Added data protection validation to mitigate accidental deletion of custom user data. ([#879](https://github.com/pingidentity/terraform-provider-pingone/issues/879))
* Inclusion of a new optional provider parameter `append_user_agent` to append a custom string to the `User-Agent` header when making API requests to the PingOne service. ([#828](https://github.com/pingidentity/terraform-provider-pingone/issues/828))
* Support ability to grant the "Application Owner" role to users, groups of users, connections and admin applications. ([#862](https://github.com/pingidentity/terraform-provider-pingone/issues/862))
* Support the new `AU` tenant region with the `com.au` top level domain. ([#828](https://github.com/pingidentity/terraform-provider-pingone/issues/828))

BUG FIXES:

* `data-source/pingone_agreement_localization`: Correct `locale` validation to add missing ISO country codes. ([#858](https://github.com/pingidentity/terraform-provider-pingone/issues/858))
* `data-source/pingone_language`: Correct `locale` validation to add missing ISO country codes. ([#858](https://github.com/pingidentity/terraform-provider-pingone/issues/858))
* `resource/pingone_agreement_localization_revision`: Fixed inability to retrieve agreement text on import. ([#861](https://github.com/pingidentity/terraform-provider-pingone/issues/861))
* `resource/pingone_agreement_localization_revision`: Fixed intermittent "The revision can not take effect in the past" error when leaving `effective_at` blank. ([#883](https://github.com/pingidentity/terraform-provider-pingone/issues/883))
* `resource/pingone_application_resource_grant`: Fixed broken grants when a resource or scope changes it's ID (scopes and resources are re-created, not triggering a re-creation of the grants) ([#863](https://github.com/pingidentity/terraform-provider-pingone/issues/863))
* `resource/pingone_application_resource_grant`: Fixed issue where the provider produces an inconsistent result after apply when new scopes are added to, or existing scopes removed from, an existing grant. ([#863](https://github.com/pingidentity/terraform-provider-pingone/issues/863))
* `resource/pingone_application_secret`: Fixed state inconsistency issue when retrieving an application's client secret. ([#781](https://github.com/pingidentity/terraform-provider-pingone/issues/781))
* `resource/pingone_application`: Fixed state inconsistency issue when retrieving an application's client secret. ([#781](https://github.com/pingidentity/terraform-provider-pingone/issues/781))
* `resource/pingone_credential_issuance_rule`: Correct `notification.template.locale` validation to add missing ISO country codes. ([#858](https://github.com/pingidentity/terraform-provider-pingone/issues/858))
* `resource/pingone_custom_domain_verify`: Fixed ineffectual `timeouts` configuration. ([#786](https://github.com/pingidentity/terraform-provider-pingone/issues/786))
* `resource/pingone_gateway`: Fixed error when configuring gateways that are generic LDAP v3 compliant directories, or OpenDJ Directory servers. ([#871](https://github.com/pingidentity/terraform-provider-pingone/issues/871))
* `resource/pingone_gateway`: Fixed issue that, when updating a `user_types` object, Terraform re-creates the full `user_types` object instead of updating the object in place. ([#798](https://github.com/pingidentity/terraform-provider-pingone/issues/798))
* `resource/pingone_language_update`: Fixed "The language must be enabled before it is set as the default" error when setting a language as enabled and the environment default. ([#884](https://github.com/pingidentity/terraform-provider-pingone/issues/884))
* `resource/pingone_language`: Correct `locale` validation to add missing ISO country codes. ([#858](https://github.com/pingidentity/terraform-provider-pingone/issues/858))
* `resource/pingone_mfa_device_policy`: Fixed blocking error when attempting to destroy the default MFA device policy for the environment.  This is now a warning instead of an error. ([#845](https://github.com/pingidentity/terraform-provider-pingone/issues/845))
* `resource/pingone_mfa_device_policy`: Resource can now be modified with Terraform if the `default` property is modified to `true` in the console or by API directly. ([#844](https://github.com/pingidentity/terraform-provider-pingone/issues/844))
* `resource/pingone_mfa_fido2_policy`: Fixed blocking error when attempting to destroy the default MFA FIDO2 policy for the environment.  This is now a warning instead of an error. ([#845](https://github.com/pingidentity/terraform-provider-pingone/issues/845))
* `resource/pingone_notification_policy`: Fixed blocking error when attempting to destroy the default notification policy for the environment.  This is now a warning instead of an error. ([#845](https://github.com/pingidentity/terraform-provider-pingone/issues/845))
* `resource/pingone_notification_template_content`: Correct `locale` validation to add missing ISO country codes. ([#858](https://github.com/pingidentity/terraform-provider-pingone/issues/858))
* `resource/pingone_password_policy`: Fixed blocking error when attempting to destroy the default password policy for the environment.  This is now a warning instead of an error. ([#845](https://github.com/pingidentity/terraform-provider-pingone/issues/845))
* `resource/pingone_resource_attribute`: Remove restrictive validation preventing config generation of default OIDC attributes ([#859](https://github.com/pingidentity/terraform-provider-pingone/issues/859))
* `resource/pingone_resource_scope`: Fixed blocking errors that result from removing multiple resource scopes that are already assigned to an application. ([#852](https://github.com/pingidentity/terraform-provider-pingone/issues/852))
* `resource/pingone_resource_scope`: Fixed blocking errors that result from removing multiple resource scopes that are already assigned to an application. ([#854](https://github.com/pingidentity/terraform-provider-pingone/issues/854))
* `resource/pingone_risk_policy`: Fixed blocking error when attempting to destroy the default risk policy for the environment.  This is now a warning instead of an error. ([#845](https://github.com/pingidentity/terraform-provider-pingone/issues/845))
* `resource/pingone_schema_attribute`: Fixed issue where schema attributes have the `required` field set to the incorrect boolean value. ([#879](https://github.com/pingidentity/terraform-provider-pingone/issues/879))
* `resource/pingone_sign_on_policy`: Fixed blocking error when attempting to destroy the default sign on policy for the environment.  This is now a warning instead of an error. ([#845](https://github.com/pingidentity/terraform-provider-pingone/issues/845))
* `resource/pingone_system_application`: Fixed intermittent "Cannot find applications by type" error. ([#885](https://github.com/pingidentity/terraform-provider-pingone/issues/885))
* `resource/pingone_verify_policy`: Fixed blocking error when attempting to destroy the default verify policy for the environment.  This is now a warning instead of an error. ([#845](https://github.com/pingidentity/terraform-provider-pingone/issues/845))
* `resource/pingone_verify_policy`: Resource can now be modified with Terraform if the `default` property is modified to `true` in the console or by API directly. ([#844](https://github.com/pingidentity/terraform-provider-pingone/issues/844))
* `resource/pingone_verify_voice_phrase_content`: Correct `locale` validation to add missing ISO country codes. ([#858](https://github.com/pingidentity/terraform-provider-pingone/issues/858))
