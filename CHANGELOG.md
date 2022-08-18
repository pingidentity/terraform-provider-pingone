## 0.2.1 (Unreleased)

NOTES:

* Added a regexp validation to any schema attribute that represents a PingOne ID.  Applies to all resources and data sources. ([#72](https://github.com/pingidentity/terraform-provider-pingone/issues/72))

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
