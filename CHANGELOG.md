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
