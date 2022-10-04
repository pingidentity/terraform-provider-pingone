## 0.6.0 (Unreleased)

## 0.5.2 (04 October 2022)

ENHANCEMENTS:

* data-source/pingone_environment: Add support for `organization_id` computed attribute. ([#166](https://github.com/pingidentity/terraform-provider-pingone/issues/166))
* resource/pingone_environment: Add support for `organization_id` computed attribute. ([#166](https://github.com/pingidentity/terraform-provider-pingone/issues/166))

## 0.5.1 (30 September 2022)

NOTES:

* bump `github.com/terraform-linters/tflint` v0.40.1 => v0.41.0 ([#157](https://github.com/pingidentity/terraform-provider-pingone/issues/157))
* pingone_application: Clarified documentation for fixed enum fields. ([#161](https://github.com/pingidentity/terraform-provider-pingone/issues/161))

ENHANCEMENTS:

* pingone_application: Add support for "Service" type applications. ([#161](https://github.com/pingidentity/terraform-provider-pingone/issues/161))

BUG FIXES:

* pingone_application: Correct the `type` parameter validation for "Custom" application types. ([#161](https://github.com/pingidentity/terraform-provider-pingone/issues/161))

## 0.5.0 (25 September 2022)

NOTES:

* bump `github.com/hashicorp/terraform-plugin-sdk/v2` v2.21.0 => v2.23.0 ([#146](https://github.com/pingidentity/terraform-provider-pingone/issues/146))
* bump `github.com/patrickcping/pingone-go-sdk-v2/management` v0.9.0 => v0.10.0 ([#145](https://github.com/pingidentity/terraform-provider-pingone/issues/145))
* bump `github.com/patrickcping/pingone-go-sdk-v2` v0.3.7 => v0.3.8 ([#145](https://github.com/pingidentity/terraform-provider-pingone/issues/145))
* bump `github.com/terraform-linters/tflint` v0.39.3 => v0.40.1 ([#147](https://github.com/pingidentity/terraform-provider-pingone/issues/147))
* pingone_application: Changed `tags` from `List` type to `Set` type. ([#149](https://github.com/pingidentity/terraform-provider-pingone/issues/149))
* pingone_gateway_credential: Corrected documentation example HCL. ([#153](https://github.com/pingidentity/terraform-provider-pingone/issues/153))

FEATURES:

* **New Data Source:** `pingone_trusted_email_domain_dkim` ([#134](https://github.com/pingidentity/terraform-provider-pingone/issues/134))
* **New Data Source:** `pingone_trusted_email_domain_ownership` ([#134](https://github.com/pingidentity/terraform-provider-pingone/issues/134))
* **New Data Source:** `pingone_trusted_email_domain_spf` ([#134](https://github.com/pingidentity/terraform-provider-pingone/issues/134))
* **New Resource:** `pingone_group_nesting` ([#144](https://github.com/pingidentity/terraform-provider-pingone/issues/144))
* **New Resource:** `pingone_mfa_settings` ([#140](https://github.com/pingidentity/terraform-provider-pingone/issues/140))
* **New Resource:** `pingone_trusted_email_domain` ([#134](https://github.com/pingidentity/terraform-provider-pingone/issues/134))
* **New Resource:** `pingone_webhook` ([#143](https://github.com/pingidentity/terraform-provider-pingone/issues/143))

ENHANCEMENTS:

* data-source/pingone_environment: Support for the `solution` environment context attribute (`CUSTOMER`, `WORKFORCE` and custom) ([#137](https://github.com/pingidentity/terraform-provider-pingone/issues/137))
* pingone_application: Add support for "External Link" type applications by adding `external_link_options` configuration block. ([#155](https://github.com/pingidentity/terraform-provider-pingone/issues/155))
* resource/pingone_environment: Support for the `solution` environment context attribute (`CUSTOMER` and custom) ([#137](https://github.com/pingidentity/terraform-provider-pingone/issues/137))
* resource/pingone_sign_on_policy_action: Added the *PingID* and *PingID Windows Login Passwordless* sign-on policy actions for workforce environments ([#141](https://github.com/pingidentity/terraform-provider-pingone/issues/141))

BUG FIXES:

* Fix panic error when a HTTP level error is returned from API after retry ([#136](https://github.com/pingidentity/terraform-provider-pingone/issues/136))
* pingone_application: Fix for `access_control_group_options.groups` showing changes when the values are the same but in different order. ([#149](https://github.com/pingidentity/terraform-provider-pingone/issues/149))
* pingone_application: Fix for `oidc_options.grant_types` showing changes when the values are the same but in different order. ([#149](https://github.com/pingidentity/terraform-provider-pingone/issues/149))
* pingone_application: Fix for `oidc_options.post_logout_redirect_uris` showing changes when the values are the same but in different order. ([#149](https://github.com/pingidentity/terraform-provider-pingone/issues/149))
* pingone_application: Fix for `oidc_options.redirect_uris` showing changes when the values are the same but in different order. ([#149](https://github.com/pingidentity/terraform-provider-pingone/issues/149))
* pingone_application: Fix for `oidc_options.response_types` showing changes when the values are the same but in different order. ([#149](https://github.com/pingidentity/terraform-provider-pingone/issues/149))
* pingone_application: Fix for `saml_options.acs_urls` showing changes when the values are the same but in different order. ([#149](https://github.com/pingidentity/terraform-provider-pingone/issues/149))
* pingone_application: Fix for `saml_options.sp_verification_certificate_ids` showing changes when the values are the same but in different order. ([#149](https://github.com/pingidentity/terraform-provider-pingone/issues/149))
* pingone_environment: Fix for `services` showing changes when the values are the same but in different order. ([#149](https://github.com/pingidentity/terraform-provider-pingone/issues/149))
* pingone_sign_on_policy_action: Fix for `registration_confirm_user_attributes` on `login` type sign-on policy action has no effect, causing change on replan. ([#152](https://github.com/pingidentity/terraform-provider-pingone/issues/152))

## 0.4.0 (11 September 2022)

NOTES:

* Bump `goreleaser/goreleaser-action` from 3.0.0 to 3.1.0 ([#87](https://github.com/pingidentity/terraform-provider-pingone/issues/87))
* Change default API call retry timeout from 30s to 10m ([#126](https://github.com/pingidentity/terraform-provider-pingone/issues/126))
* Documentation: Updates and corrections to examples ([#107](https://github.com/pingidentity/terraform-provider-pingone/issues/107))
* bump `github.com/katbyte/terrafmt` from 0.4.0 to 0.5.2 ([#65](https://github.com/pingidentity/terraform-provider-pingone/issues/65))
* bump `github.com/patrickcping/pingone-go-sdk-v2/management` v0.7.0 => v0.8.0 ([#110](https://github.com/pingidentity/terraform-provider-pingone/issues/110))
* bump `github.com/patrickcping/pingone-go-sdk-v2/management` v0.8.0 => v0.9.0 ([#128](https://github.com/pingidentity/terraform-provider-pingone/issues/128))
* bump `github.com/patrickcping/pingone-go-sdk-v2` v0.3.5 => v0.3.6 ([#110](https://github.com/pingidentity/terraform-provider-pingone/issues/110))
* bump `github.com/patrickcping/pingone-go-sdk-v2` v0.3.6 => v0.3.7 ([#128](https://github.com/pingidentity/terraform-provider-pingone/issues/128))

FEATURES:

* **New Data Source:** `pingone_certificate` ([#110](https://github.com/pingidentity/terraform-provider-pingone/issues/110))
* **New Data Source:** `pingone_certificate_export` ([#110](https://github.com/pingidentity/terraform-provider-pingone/issues/110))
* **New Data Source:** `pingone_certificate_signing_request` ([#110](https://github.com/pingidentity/terraform-provider-pingone/issues/110))
* **New Resource:** `pingone_certificate` ([#110](https://github.com/pingidentity/terraform-provider-pingone/issues/110))
* **New Resource:** `pingone_certificate_signing_response` ([#110](https://github.com/pingidentity/terraform-provider-pingone/issues/110))
* **New Resource:** `pingone_custom_domain` ([#126](https://github.com/pingidentity/terraform-provider-pingone/issues/126))
* **New Resource:** `pingone_custom_domain_ssl` ([#126](https://github.com/pingidentity/terraform-provider-pingone/issues/126))
* **New Resource:** `pingone_custom_domain_verify` ([#126](https://github.com/pingidentity/terraform-provider-pingone/issues/126))
* **New Resource:** `pingone_gateway` ([#101](https://github.com/pingidentity/terraform-provider-pingone/issues/101))
* **New Resource:** `pingone_gateway_credential` ([#101](https://github.com/pingidentity/terraform-provider-pingone/issues/101))
* **New Resource:** `pingone_gateway_role_assignment` ([#101](https://github.com/pingidentity/terraform-provider-pingone/issues/101))
* **New Resource:** `pingone_key` ([#96](https://github.com/pingidentity/terraform-provider-pingone/issues/96))

## 0.3.1 (02 September 2022)

NOTES:

* Added structure to automatically retry OAuth token calls where returned errors are retryable. ([#105](https://github.com/pingidentity/terraform-provider-pingone/issues/105))
* All resources/datasources: Certain HTTP level API errors become more readable and show better detail. ([#105](https://github.com/pingidentity/terraform-provider-pingone/issues/105))
* bump `github.com/patrickcping/pingone-go-sdk-v2/management` v0.6.0 => v0.7.0 ([#98](https://github.com/pingidentity/terraform-provider-pingone/issues/98))
* bump `github.com/patrickcping/pingone-go-sdk-v2` v0.3.4 => v0.3.5 ([#98](https://github.com/pingidentity/terraform-provider-pingone/issues/98))
* resource/pingone_sign_on_policy_action: Added schema conflict advice for `enforce_lockout_for_identity_providers` when using the `identity_provider` typed sign on policy action ([#103](https://github.com/pingidentity/terraform-provider-pingone/issues/103))
* resource/pingone_sign_on_policy_action: Added schema conflict advice for `social_provider_ids` when using the `identity_provider` typed sign on policy action ([#103](https://github.com/pingidentity/terraform-provider-pingone/issues/103))

BUG FIXES:

* resource/pingone_application: Fix `idp_signing_key_id` on `pingone_application` resource has no effect ([#98](https://github.com/pingidentity/terraform-provider-pingone/issues/98))
* resource/pingone_application_attribute_mapping: Import issue: `invalid id..` when ID is correctly specified ([#104](https://github.com/pingidentity/terraform-provider-pingone/issues/104))
* resource/pingone_application_resource_grant: Import issue: `invalid id..` when ID is correctly specified ([#104](https://github.com/pingidentity/terraform-provider-pingone/issues/104))
* resource/pingone_application_sign_on_policy_assignment: Import issue: `invalid id..` when ID is correctly specified ([#104](https://github.com/pingidentity/terraform-provider-pingone/issues/104))
* resource/pingone_identity_provider_attribute: Import issue: `invalid id..` when ID is correctly specified ([#104](https://github.com/pingidentity/terraform-provider-pingone/issues/104))
* resource/pingone_resource_scope: Import issue: `invalid id..` when ID is correctly specified ([#104](https://github.com/pingidentity/terraform-provider-pingone/issues/104))
* resource/pingone_schema_attribute: Import issue: `invalid id..` when ID is correctly specified ([#104](https://github.com/pingidentity/terraform-provider-pingone/issues/104))
* resource/pingone_sign_on_policy_action: Fix `registration_confirm_user_attributes` (`identity_provider` typed sign on policy action) has no effect ([#103](https://github.com/pingidentity/terraform-provider-pingone/issues/103))
* resource/pingone_sign_on_policy_action: Fix for `social_provider_ids` showing changes when the values are the same but in different order. ([#103](https://github.com/pingidentity/terraform-provider-pingone/issues/103))
* resource/pingone_sign_on_policy_action: Import issue: `invalid id..` when ID is correctly specified ([#104](https://github.com/pingidentity/terraform-provider-pingone/issues/104))

## 0.3.0 (30 August 2022)

NOTES:

* All resources/datasources: API errors become more readable and show better detail. ([#84](https://github.com/pingidentity/terraform-provider-pingone/issues/84))
* All resources/datasources: Added structure to automatically retry API calls where returned errors are retryable. ([#84](https://github.com/pingidentity/terraform-provider-pingone/issues/84))
* All resources/datasources: Fix `Cannot decode error response` warning on some API errors. ([#84](https://github.com/pingidentity/terraform-provider-pingone/issues/84))
* bump `github.com/patrickcping/pingone-go-sdk-v2/management` v0.5.0 => v0.6.0 ([#79](https://github.com/pingidentity/terraform-provider-pingone/issues/79))
* bump `github.com/patrickcping/pingone-go-sdk-v2` v0.3.3 => v0.3.4 ([#79](https://github.com/pingidentity/terraform-provider-pingone/issues/79))
* resource/pingone_application: Correction to documentation text ([#81](https://github.com/pingidentity/terraform-provider-pingone/issues/81))

FEATURES:

* **New Data Source:** `pingone_resource` ([#71](https://github.com/pingidentity/terraform-provider-pingone/issues/71))
* **New Data Source:** `pingone_resource_scope` ([#71](https://github.com/pingidentity/terraform-provider-pingone/issues/71))
* **New Resource:** `pingone_identity_provider` ([#79](https://github.com/pingidentity/terraform-provider-pingone/issues/79))
* **New Resource:** `pingone_identity_provider_attribute` ([#79](https://github.com/pingidentity/terraform-provider-pingone/issues/79))
* **New Resource:** `pingone_resource` ([#71](https://github.com/pingidentity/terraform-provider-pingone/issues/71))
* **New Resource:** `pingone_resource_scope` ([#71](https://github.com/pingidentity/terraform-provider-pingone/issues/71))

BUG FIXES:

* resource/pingone_application: Fix `pingone_application` error: `Once specified, refreshTokenDuration cannot be nullified` ([#88](https://github.com/pingidentity/terraform-provider-pingone/issues/88))
* resource/pingone_application_attribute_mapping: Fix import ID parsing error. ([#84](https://github.com/pingidentity/terraform-provider-pingone/issues/84))
* resource/pingone_application_resource_grant: Fix import ID parsing error. ([#84](https://github.com/pingidentity/terraform-provider-pingone/issues/84))
* resource/pingone_application_sign_on_policy_assignment: Fix import ID parsing error. ([#84](https://github.com/pingidentity/terraform-provider-pingone/issues/84))
* resource/pingone_environment: Fix index out of range panic on environment creation error. ([#84](https://github.com/pingidentity/terraform-provider-pingone/issues/84))
* resource/pingone_sign_on_policy: Fix import ID parsing error. ([#84](https://github.com/pingidentity/terraform-provider-pingone/issues/84))

## 0.2.1 (18 August 2022)

NOTES:

* Added a regexp validation to any schema attribute that represents a PingOne ID.  Applies to all resources and data sources. ([#72](https://github.com/pingidentity/terraform-provider-pingone/issues/72))

BUG FIXES:

* resource/pingone_application_resource_grant: Fix for `pingone_application_resource_grant` `scopes` showing changes when the values are the same but in different order. ([#74](https://github.com/pingidentity/terraform-provider-pingone/issues/74))

## 0.2.0 (18 August 2022)

NOTES:

* bump `github.com/golangci/golangci-lint` v1.47.2 => v1.48.0 ([#53](https://github.com/pingidentity/terraform-provider-pingone/issues/53))
* bump `github.com/hashicorp/terraform-plugin-log` v0.6.0 => v0.7.0 ([#44](https://github.com/pingidentity/terraform-provider-pingone/issues/44))
* bump `github.com/hashicorp/terraform-plugin-sdk/v2` v2.19.0 => v2.20.0 ([#46](https://github.com/pingidentity/terraform-provider-pingone/issues/46))
* bump `github.com/hashicorp/terraform-plugin-sdk/v2` v2.20.0 => v2.21.0 ([#62](https://github.com/pingidentity/terraform-provider-pingone/issues/62))
* bump `github.com/patrickcping/pingone-go-sdk-v2/management` v0.2.0 => v0.3.0 ([#48](https://github.com/pingidentity/terraform-provider-pingone/issues/48))
* bump `github.com/patrickcping/pingone-go-sdk-v2/management` v0.3.0 => v0.4.0 ([#56](https://github.com/pingidentity/terraform-provider-pingone/issues/56))
* bump `github.com/patrickcping/pingone-go-sdk-v2/management` v0.4.0 => v0.5.0 ([#42](https://github.com/pingidentity/terraform-provider-pingone/issues/42))
* bump `github.com/patrickcping/pingone-go-sdk-v2` v0.3.0 => v0.3.1 ([#48](https://github.com/pingidentity/terraform-provider-pingone/issues/48))
* bump `github.com/patrickcping/pingone-go-sdk-v2` v0.3.1 => v0.3.2 ([#56](https://github.com/pingidentity/terraform-provider-pingone/issues/56))
* bump `github.com/patrickcping/pingone-go-sdk-v2` v0.3.2 => v0.3.3 ([#42](https://github.com/pingidentity/terraform-provider-pingone/issues/42))
* bump `github.com/terraform-linters/tflint` v0.38.1 => v0.39.1 ([#45](https://github.com/pingidentity/terraform-provider-pingone/issues/45))
* bump `github.com/terraform-linters/tflint` v0.39.1 => v0.39.3 ([#64](https://github.com/pingidentity/terraform-provider-pingone/issues/64))
* resource/pingone_environment: Now sets the services (bill of materials) at the point of environment creation ([#57](https://github.com/pingidentity/terraform-provider-pingone/issues/57))

FEATURES:

* **New Data Source:** `pingone_password_policy` ([#41](https://github.com/pingidentity/terraform-provider-pingone/issues/41))
* **New Resource:** `pingone_application` ([#50](https://github.com/pingidentity/terraform-provider-pingone/issues/50))
* **New Resource:** `pingone_application_attribute_mapping` ([#50](https://github.com/pingidentity/terraform-provider-pingone/issues/50))
* **New Resource:** `pingone_application_resource_grant` ([#50](https://github.com/pingidentity/terraform-provider-pingone/issues/50))
* **New Resource:** `pingone_application_role_assignment` ([#50](https://github.com/pingidentity/terraform-provider-pingone/issues/50))
* **New Resource:** `pingone_application_sign_on_policy_assignment` ([#50](https://github.com/pingidentity/terraform-provider-pingone/issues/50))
* **New Resource:** `pingone_password_policy` ([#41](https://github.com/pingidentity/terraform-provider-pingone/issues/41))
* **New Resource:** `pingone_sign_on_policy` ([#42](https://github.com/pingidentity/terraform-provider-pingone/issues/42))
* **New Resource:** `pingone_sign_on_policy_action` ([#42](https://github.com/pingidentity/terraform-provider-pingone/issues/42))

BUG FIXES:

* resource/pingone_environment: Fix `region` attribute nil value on replan causing resource re-creation ([#51](https://github.com/pingidentity/terraform-provider-pingone/issues/51))

## 0.1.1 (10 August 2022)

NOTES:

* bump `github.com/patrickcping/pingone-go-sdk-v2/management` v0.3.0 => v0.4.0 ([#55](https://github.com/pingidentity/terraform-provider-pingone/issues/55))
* bump `github.com/patrickcping/pingone-go-sdk-v2` v0.3.1 => v0.3.2 ([#55](https://github.com/pingidentity/terraform-provider-pingone/issues/55))

BUG FIXES:

* resource/pingone_environment: Fix error `PING_ONE_DAVINCI is not a valid EnumProductType` ([#55](https://github.com/pingidentity/terraform-provider-pingone/issues/55))

## 0.1.0 (23 July 2022)

:fire: Initial provider release :fire:

FEATURES:

* **New Resource:** `pingone_environment`
* **New Resource:** `pingone_group`
* **New Resource:** `pingone_population`
* **New Resource:** `pingone_role_assignment_user`
* **New Resource:** `pingone_schema_attribute`
* **New Resource:** `pingone_user`
* **New Data Source:** `pingone_environment`
* **New Data Source:** `pingone_role`
* **New Data Source:** `pingone_schema`
