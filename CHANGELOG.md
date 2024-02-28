## 0.26.1 (28 February 2024)

NOTES:

* `resource/pingone_credential_issuer_profile`: Added customisable timeout for resource creation, used to tune the polling of a platform bootstrapped credential issuer profile, before one is forcefully created. ([#752](https://github.com/pingidentity/terraform-provider-pingone/issues/752))
* `resource/pingone_notification_policy`: Corrected documentation HCL examples. ([#747](https://github.com/pingidentity/terraform-provider-pingone/issues/747))
* `resource/pingone_sign_on_policy`: Migrated to plugin framework. ([#747](https://github.com/pingidentity/terraform-provider-pingone/issues/747))
* bump `github.com/patrickcping/pingone-go-sdk-v2/management` 0.36.0 => 0.37.0 ([#746](https://github.com/pingidentity/terraform-provider-pingone/issues/746))
* bump `github.com/patrickcping/pingone-go-sdk-v2/risk` 0.12.2 => 0.13.0 ([#746](https://github.com/pingidentity/terraform-provider-pingone/issues/746))
* bump `github.com/patrickcping/pingone-go-sdk-v2` 0.11.5 => 0.11.6 ([#746](https://github.com/pingidentity/terraform-provider-pingone/issues/746))

BUG FIXES:

* `resource/pingone_credential_issuer_profile`: Fixed race condition leading to a "A resource with the specified name already exists" error when creating a credential issuer profile at the same time as creating a new environment. ([#752](https://github.com/pingidentity/terraform-provider-pingone/issues/752))
* `resource/pingone_key`: Fixed unnecessary replacement plan when certain properties are modified. ([#747](https://github.com/pingidentity/terraform-provider-pingone/issues/747))
* `resource/pingone_mfa_fido2_policy`: Resource can now be modified with Terraform if the `default` property is modified to `true` in the console or by API directly. ([#747](https://github.com/pingidentity/terraform-provider-pingone/issues/747))
* `resource/pingone_notification_policy`: The resource can now be modified with Terraform if the `default` property is modified to `true` in the console or by API directly. ([#747](https://github.com/pingidentity/terraform-provider-pingone/issues/747))
* `resource/pingone_password_policy`: Updated the validation rule to allow the user's minimum password length to be set between 8 and 32 characters (inclusive) long. ([#740](https://github.com/pingidentity/terraform-provider-pingone/issues/740))
* `resource/pingone_risk_policy`: Fix for predictors, that are not configured in the policy, get included with a weight/score of 0. ([#751](https://github.com/pingidentity/terraform-provider-pingone/issues/751))
* `resource/pingone_risk_policy`: Fixed "Cannot find risk predictor from compact name" error when applying a policy containing "bot detection", "new device" or "adversary in the middle" predictors. ([#748](https://github.com/pingidentity/terraform-provider-pingone/issues/748))
* `resource/pingone_risk_policy`: Provider now waits for confirmation that, on destroy, the risk policy has been successfully removed in the environment. ([#748](https://github.com/pingidentity/terraform-provider-pingone/issues/748))
* `resource/pingone_risk_policy`: Resource can now be modified with Terraform if the `default` property is modified to `true` in the console or by API directly. ([#747](https://github.com/pingidentity/terraform-provider-pingone/issues/747))
* `resource/pingone_sign_on_policy`: Resource can now be modified with Terraform if the `default` property is modified to `true` in the console or by API directly. ([#747](https://github.com/pingidentity/terraform-provider-pingone/issues/747))
* `resource/pingone_verify_policy`: Resource can now be modified with Terraform if the `default` property is modified to `true` in the console or by API directly. ([#747](https://github.com/pingidentity/terraform-provider-pingone/issues/747))

## 0.26.0 (31 January 2024)

NOTES:

* `data-source/pingone_licenses`: Migrated to plugin framework. ([#728](https://github.com/pingidentity/terraform-provider-pingone/issues/728))
* bump `github.com/google/uuid` 1.5.0 => 1.6.0 ([#729](https://github.com/pingidentity/terraform-provider-pingone/issues/729))
* bump `github.com/hashicorp/terraform-plugin-framework` 1.4.2 => 1.5.0 ([#729](https://github.com/pingidentity/terraform-provider-pingone/issues/729))
* bump `github.com/hashicorp/terraform-plugin-go` 0.20.0 => 0.21.0 ([#729](https://github.com/pingidentity/terraform-provider-pingone/issues/729))
* bump `github.com/patrickcping/pingone-go-sdk-v2/management` 0.35.0 => 0.36.0 ([#729](https://github.com/pingidentity/terraform-provider-pingone/issues/729))
* bump `github.com/patrickcping/pingone-go-sdk-v2` 0.11.4 => 0.11.5 ([#729](https://github.com/pingidentity/terraform-provider-pingone/issues/729))

FEATURES:

* **New Data Source:** `pingone_application_flow_policy_assignments` ([#727](https://github.com/pingidentity/terraform-provider-pingone/issues/727))
* **New Data Source:** `pingone_application_sign_on_policy_assignments` ([#727](https://github.com/pingidentity/terraform-provider-pingone/issues/727))
* **New Resource:** `pingone_identity_propagation_plan` ([#726](https://github.com/pingidentity/terraform-provider-pingone/issues/726))

ENHANCEMENTS:

* `resource/pingone_notification_policy`: Added support for `Email` delivery method quotas. ([#722](https://github.com/pingidentity/terraform-provider-pingone/issues/722))
* `resource/pingone_notification_template_content`: Added support for `credential_verification` and `new_device_paired` notification templates. ([#720](https://github.com/pingidentity/terraform-provider-pingone/issues/720))

BUG FIXES:

* `data-source/pingone_licenses`: Fixed `data_filter.name` defined as an optional parameter.  The `data_filter.name` parameter is now required when defining the `data_filter` block. ([#728](https://github.com/pingidentity/terraform-provider-pingone/issues/728))
* `resource/pingone_sign_on_policy_action`: Corrected the "login" sign on policy action error messages when configuring a gateway with missing configuration. ([#721](https://github.com/pingidentity/terraform-provider-pingone/issues/721))

## 0.25.1 (15 January 2024)

NOTES:

* Added the ability to append custom information (by environment variable) to the default user agent string sent with every API request.  See registry documentation for details. ([#708](https://github.com/pingidentity/terraform-provider-pingone/issues/708))
* Upgrade go to `v1.21`. ([#707](https://github.com/pingidentity/terraform-provider-pingone/issues/707))
* `resource/pingone_environment`: Added ability to override `region` parameter validation with custom values. ([#706](https://github.com/pingidentity/terraform-provider-pingone/issues/706))
* bump `github.com/cloudflare/circl` 1.3.6 => 1.3.7 ([#711](https://github.com/pingidentity/terraform-provider-pingone/issues/711))
* bump `github.com/patrickcping/pingone-go-sdk-v2/agreementmanagement` v0.3.0 => v0.3.1 ([#707](https://github.com/pingidentity/terraform-provider-pingone/issues/707))
* bump `github.com/patrickcping/pingone-go-sdk-v2/authorize` v0.4.0 => v0.4.1 ([#707](https://github.com/pingidentity/terraform-provider-pingone/issues/707))
* bump `github.com/patrickcping/pingone-go-sdk-v2/credentials` v0.6.1 => v0.6.2 ([#707](https://github.com/pingidentity/terraform-provider-pingone/issues/707))
* bump `github.com/patrickcping/pingone-go-sdk-v2/management` v0.34.0 => v0.35.0 ([#707](https://github.com/pingidentity/terraform-provider-pingone/issues/707))
* bump `github.com/patrickcping/pingone-go-sdk-v2/mfa` v0.18.2 => v0.18.3 ([#707](https://github.com/pingidentity/terraform-provider-pingone/issues/707))
* bump `github.com/patrickcping/pingone-go-sdk-v2/risk` v0.12.1 => v0.12.2 ([#707](https://github.com/pingidentity/terraform-provider-pingone/issues/707))
* bump `github.com/patrickcping/pingone-go-sdk-v2/verify` v0.4.0 => v0.4.1 ([#707](https://github.com/pingidentity/terraform-provider-pingone/issues/707))
* bump `github.com/patrickcping/pingone-go-sdk-v2` v0.11.3 => v0.11.4 ([#707](https://github.com/pingidentity/terraform-provider-pingone/issues/707))

BUG FIXES:

* `resource/pingone_form`: Fix error "This attribute contains duplicate values of the same type" when configuring multiple form controls of the same type. ([#714](https://github.com/pingidentity/terraform-provider-pingone/issues/714))
* `resource/pingone_form`: Fixed "Provider produced inconsistent result after apply" error when configuring `PASSWORD` or `PASSWORD_VERIFY` type form controls with the `validation` parameter set. ([#715](https://github.com/pingidentity/terraform-provider-pingone/issues/715))

## 0.25.0 (02 January 2024)

NOTES:

* Add `lifecycle.prevent_destroy` best practice to documentation and examples for data-carrying resources, to mitigate potential accidental data loss. ([#691](https://github.com/pingidentity/terraform-provider-pingone/issues/691))
* To avoid plan inconsistency issues in earlier versions of Terraform, the provider now requires Terraform `v1.3` or later. ([#704](https://github.com/pingidentity/terraform-provider-pingone/issues/704))
* `data-source/pingone_application`: Deprecated the `saml_options.sp_verification_certificate_ids` attribute.  This attribute will be removed in the next major release.  Use the `saml_options.sp_verification.certificate_ids` attribute going forward. ([#680](https://github.com/pingidentity/terraform-provider-pingone/issues/680))
* `resource/pingone_application_attribute_mapping`: Corrected application attribute mapping documentation example when using custom OIDC scopes. ([#684](https://github.com/pingidentity/terraform-provider-pingone/issues/684))
* `resource/pingone_application`: Deprecated the `saml_options.sp_verification_certificate_ids` parameter.  This parameter will be removed in the next major release.  Use the `saml_options.sp_verification.certificate_ids` parameter going forward. ([#680](https://github.com/pingidentity/terraform-provider-pingone/issues/680))
* bump `github.com/google/uuid` v1.4.0 => v1.5.0 ([#701](https://github.com/pingidentity/terraform-provider-pingone/issues/701))
* bump `github.com/hashicorp/terraform-plugin-go` v0.19.1 => v0.20.0 ([#701](https://github.com/pingidentity/terraform-provider-pingone/issues/701))
* bump `github.com/hashicorp/terraform-plugin-mux` v0.12.0 => v0.13.0 ([#701](https://github.com/pingidentity/terraform-provider-pingone/issues/701))
* bump `github.com/hashicorp/terraform-plugin-sdk/v2` v2.30.0 => v2.31.0 ([#701](https://github.com/pingidentity/terraform-provider-pingone/issues/701))
* bump `github.com/hashicorp/terraform-plugin-testing` v1.5.1 => v1.6.0 ([#689](https://github.com/pingidentity/terraform-provider-pingone/issues/689))
* bump `github.com/patrickcping/pingone-go-sdk-v2/credentials` v0.6.0 => v0.6.1 ([#701](https://github.com/pingidentity/terraform-provider-pingone/issues/701))
* bump `github.com/patrickcping/pingone-go-sdk-v2/management` v0.33.0 => v0.34.0 ([#701](https://github.com/pingidentity/terraform-provider-pingone/issues/701))
* bump `github.com/patrickcping/pingone-go-sdk-v2/mfa` v0.18.1 => v0.18.2 ([#701](https://github.com/pingidentity/terraform-provider-pingone/issues/701))
* bump `github.com/patrickcping/pingone-go-sdk-v2/risk` v0.12.0 => v0.12.1 ([#701](https://github.com/pingidentity/terraform-provider-pingone/issues/701))
* bump `github.com/patrickcping/pingone-go-sdk-v2` v0.11.2 => v0.11.3 ([#701](https://github.com/pingidentity/terraform-provider-pingone/issues/701))
* bump `golang.org/x/crypto` v0.16.0 => v0.17.0 ([#699](https://github.com/pingidentity/terraform-provider-pingone/issues/699))

FEATURES:

* **New Resource:** `pingone_application_secret` ([#709](https://github.com/pingidentity/terraform-provider-pingone/issues/709))
* **New Resource:** `pingone_form` ([#655](https://github.com/pingidentity/terraform-provider-pingone/issues/655))
* **New Resource:** `pingone_forms_recaptcha_v2` ([#655](https://github.com/pingidentity/terraform-provider-pingone/issues/655))

ENHANCEMENTS:

* `data-source/pingone_application`: Added the `saml_options.sp_verification.authn_request_signed` attribute to support the "Enforce Signed AuthnRequest" option for SAML applications. ([#680](https://github.com/pingidentity/terraform-provider-pingone/issues/680))
* `resource/pingone_application`: Added the `saml_options.sp_verification.authn_request_signed` parameter to support the "Enforce Signed AuthnRequest" option for SAML applications. ([#680](https://github.com/pingidentity/terraform-provider-pingone/issues/680))
* `resource/pingone_key`: Added the `pkcs12_file_password` parameter to allow import of encrypted PKCS12 keys. ([#678](https://github.com/pingidentity/terraform-provider-pingone/issues/678))
* `resource/pingone_webhook`: Added the `tls_client_auth_key_pair_id` parameter to support outbound mTLS authentication to the endpoint used to post subscription messages to. ([#679](https://github.com/pingidentity/terraform-provider-pingone/issues/679))

BUG FIXES:

* Fix HTTP/HTTPS URL validation on multiple resources. See issue ([#686](https://github.com/pingidentity/terraform-provider-pingone/issues/686)) for details. ([#687](https://github.com/pingidentity/terraform-provider-pingone/issues/687))
* `resource/pingone_user`: Resolve plan inconsistency issues when using Teraform version `v1.3`. ([#704](https://github.com/pingidentity/terraform-provider-pingone/issues/704))

## 0.24.0 (30 November 2023)

NOTES:

* `data-source/pingone_population`: Update schema documentation. ([#670](https://github.com/pingidentity/terraform-provider-pingone/issues/670))
* `resource/pingone_identity_provider`: Migrated to plugin framework. ([#649](https://github.com/pingidentity/terraform-provider-pingone/issues/649))
* `resource/pingone_population`: Update schema documentation. ([#670](https://github.com/pingidentity/terraform-provider-pingone/issues/670))
* `resource/pingone_user`: Corrected documentation HCL example. ([#669](https://github.com/pingidentity/terraform-provider-pingone/issues/669))
* bump `github.com/hashicorp/terraform-plugin-go` v0.19.0 => v0.19.1 ([#674](https://github.com/pingidentity/terraform-provider-pingone/issues/674))
* bump `github.com/patrickcping/pingone-go-sdk-v2/management` v0.32.0 => v0.33.0 ([#674](https://github.com/pingidentity/terraform-provider-pingone/issues/674))
* bump `github.com/patrickcping/pingone-go-sdk-v2/mfa` v0.18.0 => v0.18.1 ([#675](https://github.com/pingidentity/terraform-provider-pingone/issues/675))
* bump `github.com/patrickcping/pingone-go-sdk-v2` v0.11.0 => v0.11.1 ([#674](https://github.com/pingidentity/terraform-provider-pingone/issues/674))
* bump `github.com/patrickcping/pingone-go-sdk-v2` v0.11.1 => v0.11.2 ([#675](https://github.com/pingidentity/terraform-provider-pingone/issues/675))

FEATURES:

* **New Data Source:** `pingone_group` ([#667](https://github.com/pingidentity/terraform-provider-pingone/issues/667))
* **New Data Source:** `pingone_groups` ([#667](https://github.com/pingidentity/terraform-provider-pingone/issues/667))
* **New Resource:** `pingone_user_group_assignment` ([#668](https://github.com/pingidentity/terraform-provider-pingone/issues/668))

ENHANCEMENTS:

* `data-source/pingone_application`: Added `oidc_options.cors_settings` and `saml_options.cors_settings` parameters. ([#673](https://github.com/pingidentity/terraform-provider-pingone/issues/673))
* `resource/pingone_application`: Added `oidc_options.cors_settings` and `saml_options.cors_settings` parameters. ([#673](https://github.com/pingidentity/terraform-provider-pingone/issues/673))

BUG FIXES:

* `resource/pingone_mfa_policy`: Fixed error when creating MFA device policy with Always Display Devices method selection. ([#675](https://github.com/pingidentity/terraform-provider-pingone/issues/675))
* `resource/pingone_user`: Fixed inconsistent result when attempting to create a user in where `account.status` is `LOCKED`. ([#654](https://github.com/pingidentity/terraform-provider-pingone/issues/654))

## 0.23.1 (11 November 2023)

NOTES:

* `resource/pingone_population_default`: Corrected documentation notes. ([#648](https://github.com/pingidentity/terraform-provider-pingone/issues/648))
* bump `github.com/hashicorp/terraform-plugin-sdk/v2` v2.29.0 => v2.30.0 ([#653](https://github.com/pingidentity/terraform-provider-pingone/issues/653))
* bump `github.com/patrickcping/pingone-go-sdk-v2/agreementmanagement` v0.2.2 => v0.3.0 ([#653](https://github.com/pingidentity/terraform-provider-pingone/issues/653))
* bump `github.com/patrickcping/pingone-go-sdk-v2/authorize` v0.3.1 => v0.4.0 ([#653](https://github.com/pingidentity/terraform-provider-pingone/issues/653))
* bump `github.com/patrickcping/pingone-go-sdk-v2/credentials` v0.5.0 => v0.6.0 ([#653](https://github.com/pingidentity/terraform-provider-pingone/issues/653))
* bump `github.com/patrickcping/pingone-go-sdk-v2/management` v0.31.0 => v0.32.0 ([#653](https://github.com/pingidentity/terraform-provider-pingone/issues/653))
* bump `github.com/patrickcping/pingone-go-sdk-v2/mfa` v0.17.0 => v0.18.0 ([#653](https://github.com/pingidentity/terraform-provider-pingone/issues/653))
* bump `github.com/patrickcping/pingone-go-sdk-v2/risk` v0.11.0 => v0.12.0 ([#653](https://github.com/pingidentity/terraform-provider-pingone/issues/653))
* bump `github.com/patrickcping/pingone-go-sdk-v2/verify` v0.3.1 => v0.4.0 ([#653](https://github.com/pingidentity/terraform-provider-pingone/issues/653))
* bump `github.com/patrickcping/pingone-go-sdk-v2` v0.10.9 => v0.11.0 ([#653](https://github.com/pingidentity/terraform-provider-pingone/issues/653))

BUG FIXES:

* `resource/pingone_environment`: Fixed `Value Conversion Error` when attempting to create an environment using calculated (unknown) `service` blocks. ([#652](https://github.com/pingidentity/terraform-provider-pingone/issues/652))
* `resource/pingone_mfa_policy`: Fixed `409 Conflict` error when attempting to destroy an MFA sign-on policy action and MFA device policy, where the MFA policy is assigned to the sign-on policy action. ([#650](https://github.com/pingidentity/terraform-provider-pingone/issues/650))

## 0.23.0 (07 November 2023)

NOTES:

* Updated documentation examples to remove reference to deprecated parameters/attributes. ([#603](https://github.com/pingidentity/terraform-provider-pingone/issues/603))
* `data-source/pingone_role`: Migrated to plugin framework. ([#592](https://github.com/pingidentity/terraform-provider-pingone/issues/592))
* `data-source/pingone_trusted_email_domain`: Corrected documentation descriptions. ([#593](https://github.com/pingidentity/terraform-provider-pingone/issues/593))
* `resource/pingone_application_role_assignment`: Migrated to plugin framework. ([#634](https://github.com/pingidentity/terraform-provider-pingone/issues/634))
* `resource/pingone_environment`: Deprecated the `default_population` block and `default_population_id` attribute in favour of the new `pingone_population_default` resource. ([#485](https://github.com/pingidentity/terraform-provider-pingone/issues/485))
* `resource/pingone_environment`: Removed the ability to import the resource including a default population.  Default populations are now managed with the `pingone_population_default` resource. ([#485](https://github.com/pingidentity/terraform-provider-pingone/issues/485))
* `resource/pingone_gateway_role_assignment`: Migrated to plugin framework. ([#636](https://github.com/pingidentity/terraform-provider-pingone/issues/636))
* `resource/pingone_role_assignment_user`: Migrated to plugin framework. ([#637](https://github.com/pingidentity/terraform-provider-pingone/issues/637))
* bump `github.com/google/uuid` v1.3.1 => v1.4.0 ([#623](https://github.com/pingidentity/terraform-provider-pingone/issues/623))
* bump `github.com/hashicorp/terraform-plugin-framework` v1.4.1 => v1.4.2 ([#623](https://github.com/pingidentity/terraform-provider-pingone/issues/623))
* bump `github.com/patrickcping/pingone-go-sdk-v2/credentials` v0.4.1 => v0.5.0 ([#623](https://github.com/pingidentity/terraform-provider-pingone/issues/623))
* bump `github.com/patrickcping/pingone-go-sdk-v2/management` v0.30.0 => v0.31.0 ([#623](https://github.com/pingidentity/terraform-provider-pingone/issues/623))
* bump `github.com/patrickcping/pingone-go-sdk-v2` v0.10.8 => v0.10.9 ([#623](https://github.com/pingidentity/terraform-provider-pingone/issues/623))

FEATURES:

* **New Data Source:** `pingone_application` ([#598](https://github.com/pingidentity/terraform-provider-pingone/issues/598))
* **New Data Source:** `pingone_gateway` ([#613](https://github.com/pingidentity/terraform-provider-pingone/issues/613))
* **New Data Source:** `pingone_roles` ([#599](https://github.com/pingidentity/terraform-provider-pingone/issues/599))
* **New Resource:** `pingone_group_role_assignment` ([#602](https://github.com/pingidentity/terraform-provider-pingone/issues/602))
* **New Resource:** `pingone_population_default` ([#600](https://github.com/pingidentity/terraform-provider-pingone/issues/600))

ENHANCEMENTS:

* `data-source/pingone_environment`: Added `service.tags` computed attribute that signifies whether the selected PingOne environment was created without example/demo configuration in the DaVinci service. ([#620](https://github.com/pingidentity/terraform-provider-pingone/issues/620))
* `data-source/pingone_population`: Add `default` computed attribute, to show whether the population is the default for the environment. ([#639](https://github.com/pingidentity/terraform-provider-pingone/issues/639))
* `data-source/pingone_role`: Now supports read-only attributes `applicable_to` and `permissions`. ([#599](https://github.com/pingidentity/terraform-provider-pingone/issues/599))
* `data-source/pingone_role`: Now supports the ability to look up role data by ID. ([#599](https://github.com/pingidentity/terraform-provider-pingone/issues/599))
* `resource/pingone_environment`: Added `service.tags` parameter allows for a creation of an environment without example/demo configuration in the DaVinci service. ([#620](https://github.com/pingidentity/terraform-provider-pingone/issues/620))

BUG FIXES:

* `resource/pingone_environment`: Fixed ineffectual `timeouts` block configuration. ([#640](https://github.com/pingidentity/terraform-provider-pingone/issues/640))

## 0.22.0 (17 October 2023)

NOTES:

* `data-source/pingone_population`: Corrected deprecated retry method. ([#574](https://github.com/pingidentity/terraform-provider-pingone/issues/574))
* `resource/pingone_certificate`: Adjusted documentation for PEM certificate import. ([#572](https://github.com/pingidentity/terraform-provider-pingone/issues/572))
* `resource/pingone_environment`: Corrected deprecated retry method. ([#574](https://github.com/pingidentity/terraform-provider-pingone/issues/574))
* `resource/pingone_key`: Migrated to plugin framework. ([#575](https://github.com/pingidentity/terraform-provider-pingone/issues/575))
* `resource/pingone_risk_predictor`: Adjust code to no longer use deprecated API parameter `composition`.  Full support of multiple risk composition policies is planned for a future release. ([#590](https://github.com/pingidentity/terraform-provider-pingone/issues/590))
* `resource/pingone_system_application`: Corrected deprecated retry method. ([#574](https://github.com/pingidentity/terraform-provider-pingone/issues/574))
* bump `github.com/hashicorp/terraform-plugin-framework` v1.4.0 => v1.4.1 ([#588](https://github.com/pingidentity/terraform-provider-pingone/issues/588))
* bump `github.com/patrickcping/pingone-go-sdk-v2/agreementmanagement` v0.2.1 => v0.2.2 ([#588](https://github.com/pingidentity/terraform-provider-pingone/issues/588))
* bump `github.com/patrickcping/pingone-go-sdk-v2/authorize` v0.3.0 => v0.3.1 ([#588](https://github.com/pingidentity/terraform-provider-pingone/issues/588))
* bump `github.com/patrickcping/pingone-go-sdk-v2/credentials` v0.4.0 => v0.4.1 ([#588](https://github.com/pingidentity/terraform-provider-pingone/issues/588))
* bump `github.com/patrickcping/pingone-go-sdk-v2/management` v0.29.0 => v0.30.0 ([#588](https://github.com/pingidentity/terraform-provider-pingone/issues/588))
* bump `github.com/patrickcping/pingone-go-sdk-v2/mfa` v0.16.0 => v0.17.0 ([#588](https://github.com/pingidentity/terraform-provider-pingone/issues/588))
* bump `github.com/patrickcping/pingone-go-sdk-v2/risk` v0.10.0 => v0.11.0 ([#588](https://github.com/pingidentity/terraform-provider-pingone/issues/588))
* bump `github.com/patrickcping/pingone-go-sdk-v2/verify` v0.3.0 => v0.3.1 ([#588](https://github.com/pingidentity/terraform-provider-pingone/issues/588))
* bump `github.com/patrickcping/pingone-go-sdk-v2` v0.10.7 => v0.10.8 ([#588](https://github.com/pingidentity/terraform-provider-pingone/issues/588))

ENHANCEMENTS:

* `data_source/pingone_credential_type`: Added support for image attributes in a verifiable credential configuration. ([#579](https://github.com/pingidentity/terraform-provider-pingone/issues/579))
* `resource/pingone_application`: Added pushed authorization request (PAR) support for OIDC configurations. ([#583](https://github.com/pingidentity/terraform-provider-pingone/issues/583))
* `resource/pingone_credential_type`: Added support for image attributes in a verifiable credential configuration. ([#579](https://github.com/pingidentity/terraform-provider-pingone/issues/579))

BUG FIXES:

* Fixed blocking error on plan on multiple resources when the parent PingOne environment is removed outside of Terraform. ([#578](https://github.com/pingidentity/terraform-provider-pingone/issues/578))
* `resource/pingone_application_resource_grant`: Fixed inability to apply application grants to the Self-Service and Portal system applications. ([#573](https://github.com/pingidentity/terraform-provider-pingone/issues/573))
* `resource/pingone_application_role_assignment`: Fixed blocking error on plan when the parent application is removed outside of Terraform. ([#578](https://github.com/pingidentity/terraform-provider-pingone/issues/578))
* `resource/pingone_gateway`: Fixed incorrect unchanged plan when updating the `user_type` parameter. ([#586](https://github.com/pingidentity/terraform-provider-pingone/issues/586))
* `resource/pingone_resource_attribute`: Fixed blocking error on delete when the parent PingOne resource is removed outside of Terraform. ([#578](https://github.com/pingidentity/terraform-provider-pingone/issues/578))
* `resource/pingone_resource_attribute`: Fixed blocking error on plan when the parent PingOne resource is removed outside of Terraform. ([#578](https://github.com/pingidentity/terraform-provider-pingone/issues/578))
* `resource/pingone_resource_scope_openid`: Fixed blocking error on delete when the parent PingOne resource is removed outside of Terraform. ([#578](https://github.com/pingidentity/terraform-provider-pingone/issues/578))
* `resource/pingone_resource_scope_openid`: Fixed blocking error on plan when the parent PingOne resource is removed outside of Terraform. ([#578](https://github.com/pingidentity/terraform-provider-pingone/issues/578))
* `resource/pingone_resource_scope_openid`: Fixed missing `resource_id` on import. ([#585](https://github.com/pingidentity/terraform-provider-pingone/issues/585))
* `resource/pingone_resource_scope_pingone_api`: Fixed blocking error on delete when the parent PingOne resource is removed outside of Terraform. ([#578](https://github.com/pingidentity/terraform-provider-pingone/issues/578))
* `resource/pingone_resource_scope_pingone_api`: Fixed blocking error on plan when the parent PingOne resource is removed outside of Terraform. ([#578](https://github.com/pingidentity/terraform-provider-pingone/issues/578))
* `resource/pingone_resource_scope_pingone_api`: Fixed missing `resource_id` on import. ([#585](https://github.com/pingidentity/terraform-provider-pingone/issues/585))
* `resource/pingone_resource_scope`: Fixed blocking error on plan when the parent PingOne resource is removed outside of Terraform. ([#578](https://github.com/pingidentity/terraform-provider-pingone/issues/578))
* `resource/pingone_system_application`: Fixed blocking error on plan when the parent environment is removed outside of Terraform. ([#578](https://github.com/pingidentity/terraform-provider-pingone/issues/578))

## 0.21.0 (18 September 2023)

NOTES:

* `data-source/pingone_organization`: Deprecated the calculated attributes `base_url_agreement_management`, `base_url_api`, `base_url_apps`, `base_url_auth`, `base_url_console`, `base_url_orchestrate`.  Please review published modules for the PingOne provider on the Terraform Registry to gain equivalent functionality. ([#564](https://github.com/pingidentity/terraform-provider-pingone/issues/564))
* `data-source/pingone_resource_scope`: Migrated to plugin framework. ([#555](https://github.com/pingidentity/terraform-provider-pingone/issues/555))
* `resource/pingone_application_resource_grant`: Deprecated the `resource_id` parameter in favour of the `resource_name` parameter to avoid dependency on the `pingone_resource` data-source. The `resource_id` parameter will be made read-only in a future release. ([#555](https://github.com/pingidentity/terraform-provider-pingone/issues/555))
* `resource/pingone_application_resource_grant`: Deprecated the `scopes` parameter in favour of the `scope_names` parameter to avoid dependency on the `pingone_resource_scope` data-source. The `scopes` parameter will be made read-only in a future release. ([#555](https://github.com/pingidentity/terraform-provider-pingone/issues/555))
* `resource/pingone_group_nesting`: Migrated to plugin framework. ([#543](https://github.com/pingidentity/terraform-provider-pingone/issues/543))
* `resource/pingone_group`: Migrated to plugin framework. ([#543](https://github.com/pingidentity/terraform-provider-pingone/issues/543))
* `resource/pingone_resource_scope_openid`: Migrated to plugin framework. ([#555](https://github.com/pingidentity/terraform-provider-pingone/issues/555))
* `resource/pingone_resource_scope_pingone_api`: Migrated to plugin framework. ([#555](https://github.com/pingidentity/terraform-provider-pingone/issues/555))
* `resource/pingone_resource_scope`: Migrated to plugin framework. ([#555](https://github.com/pingidentity/terraform-provider-pingone/issues/555))
* bump `github.com/hashicorp/terraform-plugin-framework` v1.3.5 => v1.4.0 ([#568](https://github.com/pingidentity/terraform-provider-pingone/issues/568))
* bump `github.com/hashicorp/terraform-plugin-go` v0.18.0 => v0.19.0 ([#568](https://github.com/pingidentity/terraform-provider-pingone/issues/568))
* bump `github.com/hashicorp/terraform-plugin-mux` v0.11.2 => v0.12.0 ([#568](https://github.com/pingidentity/terraform-provider-pingone/issues/568))
* bump `github.com/hashicorp/terraform-plugin-sdk/v2` v2.28.0 => v2.29.0 ([#568](https://github.com/pingidentity/terraform-provider-pingone/issues/568))
* bump `github.com/patrickcping/pingone-go-sdk-v2/management` v0.27.0 => v0.28.0 ([#556](https://github.com/pingidentity/terraform-provider-pingone/issues/556))
* bump `github.com/patrickcping/pingone-go-sdk-v2/management` v0.28.0 => v0.29.0 ([#568](https://github.com/pingidentity/terraform-provider-pingone/issues/568))
* bump `github.com/patrickcping/pingone-go-sdk-v2/risk` v0.9.0 => v0.10.0 ([#568](https://github.com/pingidentity/terraform-provider-pingone/issues/568))
* bump `github.com/patrickcping/pingone-go-sdk-v2` v0.10.5 => v0.10.6 ([#556](https://github.com/pingidentity/terraform-provider-pingone/issues/556))
* bump `github.com/patrickcping/pingone-go-sdk-v2` v0.10.6 => v0.10.7 ([#568](https://github.com/pingidentity/terraform-provider-pingone/issues/568))

ENHANCEMENTS:

* `resource/pingone_application`: Added support ability to require a signed request object for OIDC applications. ([#559](https://github.com/pingidentity/terraform-provider-pingone/issues/559))
* `resource/pingone_application`: Added support for additional refresh token replay protection configuration on OIDC applications. ([#560](https://github.com/pingidentity/terraform-provider-pingone/issues/560))
* `resource/pingone_application`: Added support for whether `requestedAuthnContext` is taken into account in SAML application policy decision-making. ([#542](https://github.com/pingidentity/terraform-provider-pingone/issues/542))
* `resource/pingone_risk_predictor`: Added support for Bot detection and Suspicious device predictor types. ([#558](https://github.com/pingidentity/terraform-provider-pingone/issues/558))
* `resource/pingone_system_application`: Support the ability to apply active theme configuration to the PingOne Portal and Self-Service applications. ([#541](https://github.com/pingidentity/terraform-provider-pingone/issues/541))

BUG FIXES:

* `resource/pingone_user`: Fixed ineffectual `initial_password` parameter. ([#566](https://github.com/pingidentity/terraform-provider-pingone/issues/566))

## 0.20.1 (05 September 2023)

NOTES:

* bump `github.com/golangci/golangci-lint` v1.54.1 => v1.54.2 ([#538](https://github.com/pingidentity/terraform-provider-pingone/issues/538))
* bump `github.com/google/uuid` v1.3.0 => v1.3.1 ([#538](https://github.com/pingidentity/terraform-provider-pingone/issues/538))
* bump `github.com/hashicorp/terraform-plugin-framework-validators` v0.11.0 => v0.12.0 ([#538](https://github.com/pingidentity/terraform-provider-pingone/issues/538))
* bump `github.com/hashicorp/terraform-plugin-sdk/v2` v2.27.0 => v2.28.0 ([#538](https://github.com/pingidentity/terraform-provider-pingone/issues/538))
* bump `github.com/patrickcping/pingone-go-sdk-v2/management` v0.26.0 => v0.27.0 ([#538](https://github.com/pingidentity/terraform-provider-pingone/issues/538))
* bump `github.com/patrickcping/pingone-go-sdk-v2` v0.10.4 => v0.10.5 ([#538](https://github.com/pingidentity/terraform-provider-pingone/issues/538))
* bump `github.com/terraform-linters/tflint` v0.47.0 => v0.48.0 ([#538](https://github.com/pingidentity/terraform-provider-pingone/issues/538))

BUG FIXES:

* `resource/pingone_environment`: Fixed "Incompatible environment region for the organization tenant" error when the `region` parameter is defaulted from the client connection. ([#535](https://github.com/pingidentity/terraform-provider-pingone/issues/535))

## 0.20.0 (29 August 2023)

NOTES:

* Code optimisation for all resources and data sources to remove redundant code. ([#507](https://github.com/pingidentity/terraform-provider-pingone/issues/507))
* Optimised code and add input validation to import resource state for every resource. ([#520](https://github.com/pingidentity/terraform-provider-pingone/issues/520))
* `resource/pingone_credential_issuer_profile`: Removal of redundant code. ([#518](https://github.com/pingidentity/terraform-provider-pingone/issues/518))
* `resource/pingone_credential_type`: Improved the `credential_type` documentation example. Corrected the placement of `card_design_template` within the example, and clarified the usage of `pingone_image` resource to assign the `background_image` and `logo_image` values. ([#518](https://github.com/pingidentity/terraform-provider-pingone/issues/518))
* `resource/pingone_sign_on_policy_action`: Fix potential "slice out of bounds" issues. ([#525](https://github.com/pingidentity/terraform-provider-pingone/issues/525))
* bump `github.com/hashicorp/terraform-plugin-framework` v1.3.4 => v1.3.5 ([#524](https://github.com/pingidentity/terraform-provider-pingone/issues/524))
* bump `github.com/patrickcping/pingone-go-sdk-v2/credentials` v0.3.1 => v0.4.0 ([#524](https://github.com/pingidentity/terraform-provider-pingone/issues/524))
* bump `github.com/patrickcping/pingone-go-sdk-v2` v0.10.3 => v0.10.4 ([#524](https://github.com/pingidentity/terraform-provider-pingone/issues/524))

FEATURES:

* **New Data Source:** `pingone_verify_voice_phrase` ([#512](https://github.com/pingidentity/terraform-provider-pingone/issues/512))
* **New Data Source:** `pingone_verify_voice_phrase_content` ([#512](https://github.com/pingidentity/terraform-provider-pingone/issues/512))
* **New Data Source:** `pingone_verify_voice_phrase_contents` ([#512](https://github.com/pingidentity/terraform-provider-pingone/issues/512))
* **New Resource:** `pingone_verify_voice_phrase` ([#512](https://github.com/pingidentity/terraform-provider-pingone/issues/512))
* **New Resource:** `pingone_verify_voice_phrase_content` ([#512](https://github.com/pingidentity/terraform-provider-pingone/issues/512))

ENHANCEMENTS:

* `data-source/pingone_credential_type`: Now supports `revoke_on_delete`, `issuer_id`, `created_at`, and `updated_at` in the datasource response. ([#518](https://github.com/pingidentity/terraform-provider-pingone/issues/518))
* `data-source/pingone_user` Enhance the user schema with the full attribute model. ([#467](https://github.com/pingidentity/terraform-provider-pingone/issues/467))
* `resource/pingone_credential_type`: Now supports `revoke_on_delete` configuration option. Read only attributes `issuer_id`, `created_at`, and `updated_at` are available in state. ([#518](https://github.com/pingidentity/terraform-provider-pingone/issues/518))
* `resource/pingone_user` Enhance the user schema with the full attribute model. ([#467](https://github.com/pingidentity/terraform-provider-pingone/issues/467))
* `resource/pingone_verify_policy`: Now supports `voice` configuration option enabling voice verification. ([#512](https://github.com/pingidentity/terraform-provider-pingone/issues/512))

BUG FIXES:

* `resource/pingone_agreement_localization_enable`: Fixed error when attempting to import resource state. ([#520](https://github.com/pingidentity/terraform-provider-pingone/issues/520))
* `resource/pingone_agreement_localization_revision`: Fixed `Cannot import non-existent remote object` error when attempting to import resource state. ([#520](https://github.com/pingidentity/terraform-provider-pingone/issues/520))
* `resource/pingone_agreement_localization`: Fixed missing `language_id` parameter value when attempting to import resource state. ([#520](https://github.com/pingidentity/terraform-provider-pingone/issues/520))
* `resource/pingone_application_attribute_mapping`: Fixed `Provider produced inconsistent result after apply` error when attempting to remove the previously configured `saml_subject_nameformat` parameter value from the `saml_subject` core attribute. ([#520](https://github.com/pingidentity/terraform-provider-pingone/issues/520))
* `resource/pingone_application_attribute_mapping`: Fixed missing `saml_subject_nameformat` parameter value when attempting to import resource state. ([#520](https://github.com/pingidentity/terraform-provider-pingone/issues/520))
* `resource/pingone_application_sign_on_policy_assignment`: Fixed missing `sign_on_policy_id` parameter value when attempting to import resource state. ([#520](https://github.com/pingidentity/terraform-provider-pingone/issues/520))
* `resource/pingone_branding_theme_default`: Fixed error when attempting to import resource state. ([#520](https://github.com/pingidentity/terraform-provider-pingone/issues/520))
* `resource/pingone_credential_issuer_profile`: Corrected `created_at` attribute value. ([#518](https://github.com/pingidentity/terraform-provider-pingone/issues/518))
* `resource/pingone_group_nesting`: Fixed missing `nested_group_id` parameter value when attempting to import resource state. ([#520](https://github.com/pingidentity/terraform-provider-pingone/issues/520))
* `resource/pingone_language_update`: Fixed missing `language_id` parameter value when attempting to import resource state. ([#520](https://github.com/pingidentity/terraform-provider-pingone/issues/520))
* `resource/pingone_mfa_settings`: Fixed `Not Found` error when attempting to import resource state. ([#520](https://github.com/pingidentity/terraform-provider-pingone/issues/520))
* `resource/pingone_notification_template_content`: Fixed missing `locale` parameter value when attempting to import resource state. ([#520](https://github.com/pingidentity/terraform-provider-pingone/issues/520))
* `resource/pingone_phone_delivery_settings`: Fixed `Value Conversion Error` error when attempting to import resource state. ([#520](https://github.com/pingidentity/terraform-provider-pingone/issues/520))
* `resource/pingone_resource_schema_attribute`: Fixed undefined response type error when attempting to import resource state. ([#520](https://github.com/pingidentity/terraform-provider-pingone/issues/520))
* `resource/pingone_resource_scope_openid`: Fixed panic error when attempting to import resource state. ([#520](https://github.com/pingidentity/terraform-provider-pingone/issues/520))
* `resource/pingone_resource_scope_pingone_api`: Fixed panic error when attempting to import resource state. ([#520](https://github.com/pingidentity/terraform-provider-pingone/issues/520))
* `resource/pingone_risk_predictor`: Fixed missing `predictor_composite.composition.condition_json` parameter value when attempting to import resource state. ([#520](https://github.com/pingidentity/terraform-provider-pingone/issues/520))
* `resource/pingone_schema_attribute`: Fixed missing `schema_name` parameter value when attempting to import resource state. ([#520](https://github.com/pingidentity/terraform-provider-pingone/issues/520))
* `resource/pingone_system_application`: Fixed validation error when attempting to import resource state. ([#520](https://github.com/pingidentity/terraform-provider-pingone/issues/520))

## 0.19.1 (16 August 2023)

NOTES:

* Code optimisation for all resources and data sources to remove duplicate service client code. ([#511](https://github.com/pingidentity/terraform-provider-pingone/issues/511))
* `resource/pingone_custom_domain_ssl`: Migrated to plugin framework. ([#506](https://github.com/pingidentity/terraform-provider-pingone/issues/506))
* `resource/pingone_custom_domain_verify`: Migrated to plugin framework. ([#506](https://github.com/pingidentity/terraform-provider-pingone/issues/506))
* `resource/pingone_custom_domain`: Migrated to plugin framework. ([#506](https://github.com/pingidentity/terraform-provider-pingone/issues/506))
* `resource/pingone_trusted_email_domain`: Migrated to plugin framework. ([#508](https://github.com/pingidentity/terraform-provider-pingone/issues/508))
* `resource/pingone_webhook`: Migrated to plugin framework. ([#505](https://github.com/pingidentity/terraform-provider-pingone/issues/505))
* bump `github.com/golangci/golangci-lint` v1.53.3 => v1.54.1 ([#515](https://github.com/pingidentity/terraform-provider-pingone/issues/515))
* bump `github.com/patrickcping/pingone-go-sdk-v2/credentials` v0.3.0 => v0.3.1 ([#515](https://github.com/pingidentity/terraform-provider-pingone/issues/515))
* bump `github.com/patrickcping/pingone-go-sdk-v2/management` v0.25.0 => v0.26.0 ([#515](https://github.com/pingidentity/terraform-provider-pingone/issues/515))
* bump `github.com/patrickcping/pingone-go-sdk-v2` v0.10.2 => v0.10.3 ([#515](https://github.com/pingidentity/terraform-provider-pingone/issues/515))

BUG FIXES:

* `data-source/pingone_credential_types` Fixed panic error when retrieving credential types. ([#515](https://github.com/pingidentity/terraform-provider-pingone/issues/515))
* `resource/pingone_resource`: Fixed blocking error on plan when OpenID resource is removed outside of Terraform. ([#501](https://github.com/pingidentity/terraform-provider-pingone/issues/501))

## 0.19.0 (08 August 2023)

NOTES:

* `data-source/pingone_population`: Optimised finding populations by name for environments with large numbers of populations present. ([#487](https://github.com/pingidentity/terraform-provider-pingone/issues/487))
* `data-source/pingone_resource`: Migrated to plugin framework. ([#493](https://github.com/pingidentity/terraform-provider-pingone/issues/493))
* `data-source/pingone_role`: Add DaVinci roles to the documentation examples. ([#496](https://github.com/pingidentity/terraform-provider-pingone/issues/496))
* `resource/pingone_environment`: Corrected documentation and examples to denote optional nature of the `default_population` block. ([#486](https://github.com/pingidentity/terraform-provider-pingone/issues/486))
* `resource/pingone_resource_attribute`: Deprecated the `resource_id` parameter in favour of the `resource_name` parameter to avoid dependency on the `pingone_resource` data-source.  The `resource_id` parameter will be made read-only in a future release. ([#493](https://github.com/pingidentity/terraform-provider-pingone/issues/493))
* `resource/pingone_schema_attribute`: Deprecated the `schema_id` parameter in favour of the optional `schema_name` parameter to avoid dependency on the `pingone_schema` data-source.  The `schema_id` parameter will be made read-only in a future release. ([#493](https://github.com/pingidentity/terraform-provider-pingone/issues/493))
* bump `github.com/hashicorp/terraform-plugin-framework-validators` v0.10.0 => v0.11.0 ([#504](https://github.com/pingidentity/terraform-provider-pingone/issues/504))
* bump `github.com/hashicorp/terraform-plugin-framework` v1.3.2 => v1.3.4 ([#504](https://github.com/pingidentity/terraform-provider-pingone/issues/504))
* bump `github.com/patrickcping/pingone-go-sdk-v2/authorize` v0.2.1 => v0.3.0 ([#504](https://github.com/pingidentity/terraform-provider-pingone/issues/504))
* bump `github.com/patrickcping/pingone-go-sdk-v2/credentials` v0.2.1 => v0.3.0 ([#504](https://github.com/pingidentity/terraform-provider-pingone/issues/504))
* bump `github.com/patrickcping/pingone-go-sdk-v2/management` v0.24.0 => v0.25.0 ([#504](https://github.com/pingidentity/terraform-provider-pingone/issues/504))
* bump `github.com/patrickcping/pingone-go-sdk-v2/mfa` v0.15.0 => v0.16.0 ([#504](https://github.com/pingidentity/terraform-provider-pingone/issues/504))
* bump `github.com/patrickcping/pingone-go-sdk-v2/risk` v0.8.1 => v0.9.0 ([#504](https://github.com/pingidentity/terraform-provider-pingone/issues/504))
* bump `github.com/patrickcping/pingone-go-sdk-v2/verify` v0.2.1 => v0.3.0 ([#504](https://github.com/pingidentity/terraform-provider-pingone/issues/504))
* bump `github.com/patrickcping/pingone-go-sdk-v2` v0.10.1 => v0.10.2 ([#504](https://github.com/pingidentity/terraform-provider-pingone/issues/504))

ENHANCEMENTS:

* `resource/pingone_mfa_settings`: Now supports `phone_extensions_enabled` option. ([#489](https://github.com/pingidentity/terraform-provider-pingone/issues/489))

BUG FIXES:

* Fixed bug when attempting to override service hostnames with the `service_endpoints` block in provider configuration. ([#498](https://github.com/pingidentity/terraform-provider-pingone/issues/498))
* `resource/pingone_credential_issuance_rule`: Corrected `digital_wallet_application_id` from `REQUIRED` to `OPTIONAL` per API specification. ([#490](https://github.com/pingidentity/terraform-provider-pingone/issues/490))
* `resource/pingone_environment`: Fix errors that occur if the `service` block is left undefined. ([#486](https://github.com/pingidentity/terraform-provider-pingone/issues/486))
* `resource/pingone_environment`: Fix for intermittent error stating the default population couldn't be updated on environment creation. ([#486](https://github.com/pingidentity/terraform-provider-pingone/issues/486))
* `resource/pingone_mfa_policy`: Fixed blocking error on plan when MFA device policy is removed outside of Terraform. ([#500](https://github.com/pingidentity/terraform-provider-pingone/issues/500))

## 0.18.1 (18 July 2023)

NOTES:

* bump `github.com/hashicorp/terraform-plugin-mux` v0.11.1 => v0.11.2 ([#481](https://github.com/pingidentity/terraform-provider-pingone/issues/481))
* bump `github.com/patrickcping/pingone-go-sdk-v2` v0.10.0 => v0.10.1 ([#483](https://github.com/pingidentity/terraform-provider-pingone/issues/483))

BUG FIXES:

* Fixed provider parameter error "Invalid parameter format.  Expected URL with https scheme" when attempting to override service hostnames. ([#483](https://github.com/pingidentity/terraform-provider-pingone/issues/483))

## 0.18.0 (17 July 2023)

NOTES:

* All resources/data sources: SDK response parsing code optimisation. ([#460](https://github.com/pingidentity/terraform-provider-pingone/issues/460))
* Code optimisation with the PingOne client SDK. ([#471](https://github.com/pingidentity/terraform-provider-pingone/issues/471))
* Corrected "Upgrade MFA Policies to use FIDO2 with Passkeys" guide text. ([#455](https://github.com/pingidentity/terraform-provider-pingone/issues/455))
* Now sets a provider-specific UserAgent on the SDK client. ([#474](https://github.com/pingidentity/terraform-provider-pingone/issues/474))
* `data-source/pingone_environment`: Optimised environment filtering by name. ([#469](https://github.com/pingidentity/terraform-provider-pingone/issues/469))
* `resource/pingone_application_resource_grant`: Migrated to plugin framework. ([#456](https://github.com/pingidentity/terraform-provider-pingone/issues/456))
* bump `github.com/hashicorp/terraform-plugin-docs` v0.15.0 => v0.16.0 ([#473](https://github.com/pingidentity/terraform-provider-pingone/issues/473))
* bump `github.com/hashicorp/terraform-plugin-framework-timeouts` v0.4.0 => v0.4.1 ([#473](https://github.com/pingidentity/terraform-provider-pingone/issues/473))
* bump `github.com/patrickcping/pingone-go-sdk-v2/agreementmanagement` v0.2.0 => v0.2.1 ([#473](https://github.com/pingidentity/terraform-provider-pingone/issues/473))
* bump `github.com/patrickcping/pingone-go-sdk-v2/authorize` v0.2.0 => v0.2.1 ([#473](https://github.com/pingidentity/terraform-provider-pingone/issues/473))
* bump `github.com/patrickcping/pingone-go-sdk-v2/credentials` v0.2.0 => v0.2.1 ([#473](https://github.com/pingidentity/terraform-provider-pingone/issues/473))
* bump `github.com/patrickcping/pingone-go-sdk-v2/management` v0.23.0 => v0.24.0 ([#473](https://github.com/pingidentity/terraform-provider-pingone/issues/473))
* bump `github.com/patrickcping/pingone-go-sdk-v2/mfa` v0.14.0 => v0.15.0 ([#473](https://github.com/pingidentity/terraform-provider-pingone/issues/473))
* bump `github.com/patrickcping/pingone-go-sdk-v2/risk` v0.8.0 => v0.8.1 ([#473](https://github.com/pingidentity/terraform-provider-pingone/issues/473))
* bump `github.com/patrickcping/pingone-go-sdk-v2/verify` v0.2.0 => v0.2.1 ([#473](https://github.com/pingidentity/terraform-provider-pingone/issues/473))
* bump `github.com/patrickcping/pingone-go-sdk-v2` v0.9.0 => v0.10.0 ([#473](https://github.com/pingidentity/terraform-provider-pingone/issues/473))

FEATURES:

* **New Data Source:** `pingone_user_role_assignments` ([#479](https://github.com/pingidentity/terraform-provider-pingone/issues/479))
* **New Resource:** `pingone_key_rotation_policy` ([#466](https://github.com/pingidentity/terraform-provider-pingone/issues/466))

ENHANCEMENTS:

* Added provider parameter to connect to the PingOne API service via HTTP proxy. ([#471](https://github.com/pingidentity/terraform-provider-pingone/issues/471))
* `resource/pingone_application_resource_grant`: Add validation to disallow assignment of the `openid` scope from the `openid` resource to avoid error. ([#457](https://github.com/pingidentity/terraform-provider-pingone/issues/457))
* `resource/pingone_application`: Now supports the `saml_option` of `slo_window`. ([#468](https://github.com/pingidentity/terraform-provider-pingone/issues/468))
* `resource/pingone_identity_provider`: Now supports the SAML IdP configuration options: `slo_binding`, `slo_endpoint`, `slo_response_endpoint` and `slo_window`. ([#468](https://github.com/pingidentity/terraform-provider-pingone/issues/468))
* `resource/pingone_mfa_policy`: Now supports `new_device_notification` option. ([#477](https://github.com/pingidentity/terraform-provider-pingone/issues/477))
* `resource/pingone_notification_policy`: Now supports country limit configuration. ([#458](https://github.com/pingidentity/terraform-provider-pingone/issues/458))
* `resource/pingone_webhook`: Now supports `ip_address_exposed` and `useragent_exposed` in filter options configuration. ([#470](https://github.com/pingidentity/terraform-provider-pingone/issues/470))

BUG FIXES:

* `resource/pingone_application`: Fixed bug where the `pkce_enforcement` parameter wasn't being configured correctly on OIDC worker apps. ([#475](https://github.com/pingidentity/terraform-provider-pingone/issues/475))

## 0.17.1 (05 July 2023)

NOTES:

* `data-source/pingone_user`: Migrated to plugin framework. ([#453](https://github.com/pingidentity/terraform-provider-pingone/issues/453))
* `data-source/pingone_users`: Migrated to plugin framework. ([#453](https://github.com/pingidentity/terraform-provider-pingone/issues/453))
* `resource/pingone_mfa_fido_policy`: Corrected MFA policy upgrade guide link in registry documentation. ([#454](https://github.com/pingidentity/terraform-provider-pingone/issues/454))
* `resource/pingone_mfa_policy`: Corrected MFA policy upgrade guide link in registry documentation. ([#454](https://github.com/pingidentity/terraform-provider-pingone/issues/454))
* `resource/pingone_user`: Migrated to plugin framework. ([#453](https://github.com/pingidentity/terraform-provider-pingone/issues/453))

## 0.17.0 (05 July 2023)

BREAKING CHANGES:

* `resource/pingone_mfa_fido_policy`: This resource is deprecated, please use the `pingone_mfa_fido2_policy` resource going forward.  This resource is no longer configurable for environments created after 19th June 2023, nor environments that have been upgraded to use the latest FIDO2 policies. Existing environments that were created before 19th June 2023 and have not been upgraded can continue to use this resource to facilitate migration. ([#441](https://github.com/pingidentity/terraform-provider-pingone/issues/441))
* `resource/pingone_mfa_policy`: The `platform` and `security_key` FIDO device types are deprecated and need to be replaced with the `fido2` device type.  `platform` and `security_key` are no longer configurable for newly created environments, or existing environments that have not had their environment upgraded to use the latest FIDO2 policies.  Existing environments that have not been upgraded to use the latest FIDO2 policies can continue to use the factors to facilitate migration. ([#437](https://github.com/pingidentity/terraform-provider-pingone/issues/437))

NOTES:

* Changed "Risk" documentation subcategory to "Protect", for the release of PingOne Protect. ([#452](https://github.com/pingidentity/terraform-provider-pingone/issues/452))
* Code optimisations in each resource/data source to remove the need to override the region on each operation. ([#439](https://github.com/pingidentity/terraform-provider-pingone/issues/439))
* Minor adjustments of multiple documentation HCL examples. ([#389](https://github.com/pingidentity/terraform-provider-pingone/issues/389))
* `resource/pingone_mfa_application_push_credential`: Corrected documentation HCL example. ([#389](https://github.com/pingidentity/terraform-provider-pingone/issues/389))
* `resource/pingone_mfa_application_push_credential`: Migrated to plugin framework. **IMPORTANT**: The resource will show drift and will need to be re-applied to ensure consistency in the stored resource state. ([#426](https://github.com/pingidentity/terraform-provider-pingone/issues/426))
* bump `github.com/hashicorp/terraform-plugin-framework-timeouts` v0.3.1 => v0.4.0 ([#449](https://github.com/pingidentity/terraform-provider-pingone/issues/449))
* bump `github.com/hashicorp/terraform-plugin-framework` v1.3.1 => v1.3.2 ([#449](https://github.com/pingidentity/terraform-provider-pingone/issues/449))
* bump `github.com/hashicorp/terraform-plugin-go` v0.15.0 => v0.17.0 ([#449](https://github.com/pingidentity/terraform-provider-pingone/issues/449))
* bump `github.com/hashicorp/terraform-plugin-go` v0.17.0 => v0.18.0 ([#450](https://github.com/pingidentity/terraform-provider-pingone/issues/450))
* bump `github.com/hashicorp/terraform-plugin-mux` v0.10.0 => v0.11.1 ([#449](https://github.com/pingidentity/terraform-provider-pingone/issues/449))
* bump `github.com/hashicorp/terraform-plugin-sdk/v2` v2.26.1 => v2.27.0 ([#449](https://github.com/pingidentity/terraform-provider-pingone/issues/449))
* bump `github.com/patrickcping/pingone-go-sdk-v2/agreementmanagement` v0.1.4 => v0.2.0 ([#449](https://github.com/pingidentity/terraform-provider-pingone/issues/449))
* bump `github.com/patrickcping/pingone-go-sdk-v2/authorize` v0.1.7 => v0.2.0 ([#449](https://github.com/pingidentity/terraform-provider-pingone/issues/449))
* bump `github.com/patrickcping/pingone-go-sdk-v2/credentials` v0.1.0 => v0.2.0 ([#449](https://github.com/pingidentity/terraform-provider-pingone/issues/449))
* bump `github.com/patrickcping/pingone-go-sdk-v2/management` v0.22.0 => v0.23.0 ([#449](https://github.com/pingidentity/terraform-provider-pingone/issues/449))
* bump `github.com/patrickcping/pingone-go-sdk-v2/mfa` v0.12.0 => v0.13.0 ([#449](https://github.com/pingidentity/terraform-provider-pingone/issues/449))
* bump `github.com/patrickcping/pingone-go-sdk-v2/mfa` v0.13.0 => v0.14.0 ([#437](https://github.com/pingidentity/terraform-provider-pingone/issues/437))
* bump `github.com/patrickcping/pingone-go-sdk-v2/risk` v0.7.1 => v0.8.0 ([#449](https://github.com/pingidentity/terraform-provider-pingone/issues/449))
* bump `github.com/patrickcping/pingone-go-sdk-v2/verify` v0.1.0 => v0.2.0 ([#449](https://github.com/pingidentity/terraform-provider-pingone/issues/449))
* bump `github.com/patrickcping/pingone-go-sdk-v2` v0.8.0 => v0.9.0 ([#449](https://github.com/pingidentity/terraform-provider-pingone/issues/449))
* bump `github.com/patrickcping/pingone-go-sdk-v2` v0.9.0 => v0.9.1 ([#437](https://github.com/pingidentity/terraform-provider-pingone/issues/437))

FEATURES:

* **New Data Source:** `pingone_mfa_policies` ([#437](https://github.com/pingidentity/terraform-provider-pingone/issues/437))
* **New Data Source:** `pingone_phone_delivery_settings_list` ([#419](https://github.com/pingidentity/terraform-provider-pingone/issues/419))
* **New Guide:** `Upgrade MFA Policies to use FIDO2 with Passkeys` ([#437](https://github.com/pingidentity/terraform-provider-pingone/issues/437))
* **New Resource:** `pingone_mfa_fido2_policy` ([#441](https://github.com/pingidentity/terraform-provider-pingone/issues/441))
* **New Resource:** `pingone_mfa_policies` ([#437](https://github.com/pingidentity/terraform-provider-pingone/issues/437))
* **New Resource:** `pingone_notification_settings` ([#419](https://github.com/pingidentity/terraform-provider-pingone/issues/419))
* **New Resource:** `pingone_phone_delivery_settings` ([#419](https://github.com/pingidentity/terraform-provider-pingone/issues/419))

ENHANCEMENTS:

* Add provider configuration parameters to be able to override the PingOne service URL hostnames. ([#439](https://github.com/pingidentity/terraform-provider-pingone/issues/439))
* `resource/pingone_mfa_application_push_credential`: PingOne MFA has moved to Firebase Cloud Messaging for sending push messages.  `fcm.key` has now been deprecated, `fcm.google_service_account_credentials` should be used going forward. ([#426](https://github.com/pingidentity/terraform-provider-pingone/issues/426))
* `resource/pingone_mfa_policy`: Add support for the new `fido2` MFA device type to enable support for passkeys.  The `fido2` device type is only configurable for newly created environments, or existing environments that have been upgraded to use the latest FIDO2 policies. ([#437](https://github.com/pingidentity/terraform-provider-pingone/issues/437))
* `resource/pingone_mfa_policy`: Support the ability to phase out MFA devices using new `pairing_disabled` parameters for each device type in the policy. ([#437](https://github.com/pingidentity/terraform-provider-pingone/issues/437))
* `resource/pingone_mfa_policy`: Supports the ability to define the pairing key lifetime and push limit for mobile applications. ([#437](https://github.com/pingidentity/terraform-provider-pingone/issues/437))
* `resource/pingone_notification_template_content`: Add support for P1Verify and P1Credentials notification templates: `email_phone_verification`, `id_verification`, `credential_issued`, `credential_updated`, `digital_wallet_pairing`, `credential_revoked`. ([#428](https://github.com/pingidentity/terraform-provider-pingone/issues/428))

BUG FIXES:

* Global fix for "*value* is not a valid Enum*Object*" errors. ([#449](https://github.com/pingidentity/terraform-provider-pingone/issues/449))
* `resource/pingone_branding_settings`: Added missing schema validators to prevent misconfiguration of the single-object `logo_image` block. ([#440](https://github.com/pingidentity/terraform-provider-pingone/issues/440))
* `resource/pingone_branding_theme`: Added missing schema validators to prevent misconfiguration of the single-object `background_image` and `logo` blocks. ([#440](https://github.com/pingidentity/terraform-provider-pingone/issues/440))
* `resource/pingone_risk_policy`: Fix for "BOT is not a valid EnumPredictorType" error. ([#449](https://github.com/pingidentity/terraform-provider-pingone/issues/449))
* `resource/pingone_risk_predictor`: Fix for "BOT is not a valid EnumPredictorType" error. ([#449](https://github.com/pingidentity/terraform-provider-pingone/issues/449))

## 0.16.0 (19 June 2023)

NOTES:

* Adjusted documentation for multiple resources/datasources to clarify PingOne ID attribute validation and plan replacement on change. ([#404](https://github.com/pingidentity/terraform-provider-pingone/issues/404))
* `data-source/pingone_credential_issuance_rule`: Code optimisation of ENUM values to TF state. ([#417](https://github.com/pingidentity/terraform-provider-pingone/issues/417))
* `data-source/pingone_credential_type`: Code optimisation of ENUM values to TF state. ([#417](https://github.com/pingidentity/terraform-provider-pingone/issues/417))
* `data-source/pingone_environment`: Code optimisation of ENUM values to TF state. ([#417](https://github.com/pingidentity/terraform-provider-pingone/issues/417))
* `data_source/pingone_credential_type`: Corrected typo in data source description. ([#406](https://github.com/pingidentity/terraform-provider-pingone/issues/406))
* `resource/pingone_agreement_localization_revision`: Code optimisation of ENUM values to TF state. ([#417](https://github.com/pingidentity/terraform-provider-pingone/issues/417))
* `resource/pingone_branding_settings`: Adjusted schema such that the icon `id` required parameter no longer triggers a replacement plan on change. ([#404](https://github.com/pingidentity/terraform-provider-pingone/issues/404))
* `resource/pingone_credential_issuance_rule`: Code optimisation of ENUM values to TF state. ([#417](https://github.com/pingidentity/terraform-provider-pingone/issues/417))
* `resource/pingone_credential_type`: Code optimisation of ENUM values to TF state. ([#417](https://github.com/pingidentity/terraform-provider-pingone/issues/417))
* `resource/pingone_environment`: Code optimisation of ENUM values to TF state. ([#417](https://github.com/pingidentity/terraform-provider-pingone/issues/417))
* `resource/pingone_flow_policy_assignment`: Adjusted schema such that the `flow_policy_id` required parameter no longer triggers a replacement plan on change. ([#404](https://github.com/pingidentity/terraform-provider-pingone/issues/404))
* `resource/pingone_schema_attribute`: Migrated to plugin framework. ([#414](https://github.com/pingidentity/terraform-provider-pingone/issues/414))
* `resource/pingone_sign_on_policy_action`: Adjust documentation to clarify the where conditions can and cannot be used in a policy action. ([#412](https://github.com/pingidentity/terraform-provider-pingone/issues/412))
* bump `github.com/golangci/golangci-lint` v1.53.2 => v1.53.3 ([#424](https://github.com/pingidentity/terraform-provider-pingone/issues/424))
* bump `github.com/hashicorp/terraform-plugin-docs` v0.14.1 => v0.15.0 ([#424](https://github.com/pingidentity/terraform-provider-pingone/issues/424))
* bump `github.com/hashicorp/terraform-plugin-framework` v1.2.0 => v1.3.1 ([#424](https://github.com/pingidentity/terraform-provider-pingone/issues/424))
* bump `github.com/hashicorp/terraform-plugin-log` v0.8.0 => v0.9.0 ([#424](https://github.com/pingidentity/terraform-provider-pingone/issues/424))
* bump `github.com/patrickcping/pingone-go-sdk-v2/mfa` v0.11.0 => v0.12.0 ([#424](https://github.com/pingidentity/terraform-provider-pingone/issues/424))
* bump `github.com/patrickcping/pingone-go-sdk-v2` v0.7.3 => v0.8.0 ([#424](https://github.com/pingidentity/terraform-provider-pingone/issues/424))
* bump `github.com/terraform-linters/tflint` v0.46.1 => v0.47.0 ([#424](https://github.com/pingidentity/terraform-provider-pingone/issues/424))

FEATURES:

* **New Data Source:** `pingone_verify_policies` ([#403](https://github.com/pingidentity/terraform-provider-pingone/issues/403))
* **New Data Source:** `pingone_verify_policy` ([#403](https://github.com/pingidentity/terraform-provider-pingone/issues/403))
* **New Resource:** `pingone_verify_policy` ([#403](https://github.com/pingidentity/terraform-provider-pingone/issues/403))

ENHANCEMENTS:

* `resource/pingone_schema_attribute`: Supports properties for enumerated values and regular expression validation. ([#414](https://github.com/pingidentity/terraform-provider-pingone/issues/414))
* `resource/pingone_sign_on_policy_action`: Support user provisioning gateway configuration on the "Login" sign-on policy action ([#407](https://github.com/pingidentity/terraform-provider-pingone/issues/407))

BUG FIXES:

* `data_source/pingone_credential_issuer_profile`: Fixed mismatched `created_at` and `updated_at` mapping. ([#406](https://github.com/pingidentity/terraform-provider-pingone/issues/406))
* `resource/pingone_branding_theme`: Fixed change to theme icon/background ID triggers a replacement plan on change, leading to removal failures if the theme is set as default. ([#404](https://github.com/pingidentity/terraform-provider-pingone/issues/404))
* `resource/pingone_credential_issuance_rule`: Fixed incorrect replacement of resource when the `digital_wallet_id` value was changed. ([#406](https://github.com/pingidentity/terraform-provider-pingone/issues/406))
* `resource/pingone_credential_issuer_profile`: Fixed mismatched `created_at` and `updated_at` mapping. ([#406](https://github.com/pingidentity/terraform-provider-pingone/issues/406))
* `resource/pingone_digital_wallet_application`: Fixed incorrect replacement of resource when the `application_id` value was changed. ([#406](https://github.com/pingidentity/terraform-provider-pingone/issues/406))
* `resource/pingone_sign_on_policy_action`: Fix panic crash when defining `conditions.user_attribute_equals` and/or `conditions.user_is_member_of_any_population_id` in a sign-on policy action that is priority 1. ([#412](https://github.com/pingidentity/terraform-provider-pingone/issues/412))

## 0.15.1 (07 June 2023)

NOTES:

* `resource/pingone_population`: Migrated to plugin framework. ([#400](https://github.com/pingidentity/terraform-provider-pingone/issues/400))

BUG FIXES:

* `resource/pingone_population`: Fixed panic on plan when population is removed outside of Terraform. ([#400](https://github.com/pingidentity/terraform-provider-pingone/issues/400))

## 0.15.0 (06 June 2023)

NOTES:

* `resource/pingone_application_attribute_mapping`: Updated documentation to reflect support for only administrator defined applications. ([#395](https://github.com/pingidentity/terraform-provider-pingone/issues/395))
* `resource/pingone_application_flow_policy_assignment`: Updated documentation to reflect support for Portal and Self-Service built-in system applications. ([#395](https://github.com/pingidentity/terraform-provider-pingone/issues/395))
* `resource/pingone_application_resource_grant`: Updated documentation to reflect support for Portal and Self-Service built-in system applications. ([#395](https://github.com/pingidentity/terraform-provider-pingone/issues/395))
* `resource/pingone_application_role_assignment`: Updated documentation to reflect support for only administrator defined applications. ([#395](https://github.com/pingidentity/terraform-provider-pingone/issues/395))
* `resource/pingone_application_sign_on_policy_assignment`: Updated documentation to reflect support for Portal and Self-Service built-in system applications. ([#395](https://github.com/pingidentity/terraform-provider-pingone/issues/395))
* `resource/pingone_application`: Updated documentation to reflect support for only administrator defined applications. ([#395](https://github.com/pingidentity/terraform-provider-pingone/issues/395))
* `resource/pingone_branding_settings`: Migrated to plugin framework. ([#374](https://github.com/pingidentity/terraform-provider-pingone/issues/374))
* `resource/pingone_branding_theme`: Corrected documentation example syntax. ([#374](https://github.com/pingidentity/terraform-provider-pingone/issues/374))
* `resource/pingone_branding_theme`: Migrated to plugin framework. ([#374](https://github.com/pingidentity/terraform-provider-pingone/issues/374))
* `resource/pingone_credential_type`: Improved the documentation for the `title` and `description` attributes and explained their correlation to fields in the `card_design_template`. ([#377](https://github.com/pingidentity/terraform-provider-pingone/issues/377))
* `resource/pingone_risk_predictor`: Corrected custom map type documentation example. ([#378](https://github.com/pingidentity/terraform-provider-pingone/issues/378))
* `resource/pingone_risk_predictor`: Corrected documentation examples. ([#387](https://github.com/pingidentity/terraform-provider-pingone/issues/387))
* bump `github.com/golangci/golangci-lint` v1.52.2 => v1.53.2 ([#396](https://github.com/pingidentity/terraform-provider-pingone/issues/396))
* bump `github.com/patrickcping/pingone-go-sdk-v2/management` v0.20.0 => v0.22.0 ([#396](https://github.com/pingidentity/terraform-provider-pingone/issues/396))
* bump `github.com/patrickcping/pingone-go-sdk-v2/mfa` v0.10.0 => v0.11.0 ([#396](https://github.com/pingidentity/terraform-provider-pingone/issues/396))
* bump `github.com/patrickcping/pingone-go-sdk-v2/risk` v0.6.0 => v0.7.1 ([#396](https://github.com/pingidentity/terraform-provider-pingone/issues/396))
* bump `github.com/patrickcping/pingone-go-sdk-v2` v0.7.0 => v0.7.3 ([#396](https://github.com/pingidentity/terraform-provider-pingone/issues/396))

FEATURES:

* **New Resource:** `pingone_branding_theme_default` ([#375](https://github.com/pingidentity/terraform-provider-pingone/issues/375))
* **New Resource:** `pingone_risk_policy` ([#381](https://github.com/pingidentity/terraform-provider-pingone/issues/381))
* **New Resource:** `pingone_system_application` ([#395](https://github.com/pingidentity/terraform-provider-pingone/issues/395))

ENHANCEMENTS:

* `resource/pingone_application_role_assignment`: Added validation logic to ensure that only valid applications are accepted for role assignment. ([#395](https://github.com/pingidentity/terraform-provider-pingone/issues/395))

BUG FIXES:

* `resource/pingone_credential_type`: Fix the validation rules for `title` and `description`. The rules incorrectly compared the `metadata.name` and `metadata.description` attributes to the ${cardTitle} and ${cardSubTitle} fields in the `card_design_template`. The rules are now correctly applied to `title` and `description`. ([#377](https://github.com/pingidentity/terraform-provider-pingone/issues/377))

## 0.14.0 (23 May 2023)

BREAKING CHANGES:

* `resource/pingone_application`: Signature algorithms `SHA224withRSA` and `SHA224withECDSA` removed as they are no longer supported by the platform. (P14C-50332) ([#358](https://github.com/pingidentity/terraform-provider-pingone/issues/358))
* `resource/pingone_key`: Signature algorithms `SHA224withRSA` and `SHA224withECDSA` removed as they are no longer supported by the platform. (P14C-50332) ([#358](https://github.com/pingidentity/terraform-provider-pingone/issues/358))

NOTES:

* Upgraded the provider protocol version from 5 to 6.  Use of the provider requires Terraform CLI 1.1 or later. ([#354](https://github.com/pingidentity/terraform-provider-pingone/issues/354))
* `pingone_application`: Deprecated `oidc_options.bundle_id` and `oidc_options.package_name` from the schema.  Customers should use `oidc_options.mobile_app.bundle_id` and `oidc_options.mobile_app.package_name` going forward. ([#363](https://github.com/pingidentity/terraform-provider-pingone/issues/363))
* `pingone_application`: Updated Native mobile example in registry documentation to remove deprecated attributes. ([#365](https://github.com/pingidentity/terraform-provider-pingone/issues/365))
* `resource/pingone_application_flow_policy_assignment`: Update documentation example to select from multiple DaVinci application flow policies. ([#360](https://github.com/pingidentity/terraform-provider-pingone/issues/360))
* `resource/pingone_mfa_application_push_credential`: Documentation: Remove deprecated attributes from example HCL. ([#372](https://github.com/pingidentity/terraform-provider-pingone/issues/372))
* `resource/pingone_mfa_policy`: Documentation: Remove deprecated attributes from example HCL. ([#372](https://github.com/pingidentity/terraform-provider-pingone/issues/372))
* bump `github.com/hashicorp/terraform-plugin-mux` from v0.9.0 => v0.10.0 ([#354](https://github.com/pingidentity/terraform-provider-pingone/issues/354))
* bump `github.com/patrickcping/pingone-go-sdk-v2/agreementmanagement` v0.1.3 => v0.1.4 ([#361](https://github.com/pingidentity/terraform-provider-pingone/issues/361))
* bump `github.com/patrickcping/pingone-go-sdk-v2/authorize` v0.1.6 => v0.1.7 ([#361](https://github.com/pingidentity/terraform-provider-pingone/issues/361))
* bump `github.com/patrickcping/pingone-go-sdk-v2/management` v0.19.1 => v0.20.0 ([#361](https://github.com/pingidentity/terraform-provider-pingone/issues/361))
* bump `github.com/patrickcping/pingone-go-sdk-v2/mfa` v0.9.3 => v0.10.0 ([#361](https://github.com/pingidentity/terraform-provider-pingone/issues/361))
* bump `github.com/patrickcping/pingone-go-sdk-v2` v0.6.4 => v0.7.0 ([#361](https://github.com/pingidentity/terraform-provider-pingone/issues/361))

FEATURES:

* **New Data Source:** `pingone_credential_issuance_rule` ([#359](https://github.com/pingidentity/terraform-provider-pingone/issues/359))
* **New Data Source:** `pingone_credential_issuer_profile` ([#359](https://github.com/pingidentity/terraform-provider-pingone/issues/359))
* **New Data Source:** `pingone_credential_type` ([#359](https://github.com/pingidentity/terraform-provider-pingone/issues/359))
* **New Data Source:** `pingone_credential_types` ([#359](https://github.com/pingidentity/terraform-provider-pingone/issues/359))
* **New Data Source:** `pingone_digital_wallet_application` ([#359](https://github.com/pingidentity/terraform-provider-pingone/issues/359))
* **New Data Source:** `pingone_digital_wallet_applications` ([#359](https://github.com/pingidentity/terraform-provider-pingone/issues/359))
* **New Resource:** `pingone_credential_issuance_rule` ([#359](https://github.com/pingidentity/terraform-provider-pingone/issues/359))
* **New Resource:** `pingone_credential_issuer_profile` ([#359](https://github.com/pingidentity/terraform-provider-pingone/issues/359))
* **New Resource:** `pingone_credential_type` ([#359](https://github.com/pingidentity/terraform-provider-pingone/issues/359))
* **New Resource:** `pingone_digital_wallet_application` ([#359](https://github.com/pingidentity/terraform-provider-pingone/issues/359))
* **New Resource:** `pingone_risk_predictor` ([#350](https://github.com/pingidentity/terraform-provider-pingone/issues/350))

BUG FIXES:

* `pingone_application`: Fix unexpected replacement plan when specifying a mobile native application without configuring deprecated attributes `oidc_options.bundle_id` or `oidc_options.package_name`. ([#365](https://github.com/pingidentity/terraform-provider-pingone/issues/365))

## 0.13.1 (02 May 2023)

NOTES:

* Simplified SDK request retry code for all resources/datasources. ([#348](https://github.com/pingidentity/terraform-provider-pingone/issues/348))
* bump `github.com/patrickcping/pingone-go-sdk-v2/agreementmanagement` v0.1.2 => v0.1.3 ([#351](https://github.com/pingidentity/terraform-provider-pingone/issues/351))
* bump `github.com/patrickcping/pingone-go-sdk-v2/authorize` v0.1.5 => v0.1.6 ([#351](https://github.com/pingidentity/terraform-provider-pingone/issues/351))
* bump `github.com/patrickcping/pingone-go-sdk-v2/management` v0.19.0 => v0.19.1 ([#351](https://github.com/pingidentity/terraform-provider-pingone/issues/351))
* bump `github.com/patrickcping/pingone-go-sdk-v2/mfa` v0.9.2 => v0.9.3 ([#351](https://github.com/pingidentity/terraform-provider-pingone/issues/351))
* bump `github.com/patrickcping/pingone-go-sdk-v2` v0.6.3 => v0.6.4 ([#351](https://github.com/pingidentity/terraform-provider-pingone/issues/351))

BUG FIXES:

* `resource/pingone_application`: Fix for "There was an unexpected error with the service" error when attempting to create an application immediately after creation of the parent environment. (2) ([#351](https://github.com/pingidentity/terraform-provider-pingone/issues/351))
* `resource/pingone_notification_template_content`: Fix issue where the notification template content with no variant is configured correctly in PingOne but the template content is not effective. ([#349](https://github.com/pingidentity/terraform-provider-pingone/issues/349))

## 0.13.0 (25 April 2023)

BREAKING CHANGES:

* `resource/pingone_application`: Moved from SafetyNet Attestation API to Google Play Integration API for Android integrity detection (P14C-37640).  Customers wanting to enable Android/Google integrity detection for mobile apps will need to upgrade to the latest provider version as `oidc_options.mobile_app.integrity_detection` now requires the `google_play` block to be defined. ([#344](https://github.com/pingidentity/terraform-provider-pingone/issues/344))

NOTES:

* Updated the external documentation site link to `terraform.pingidentity.com` for the getting started guide on the index docs page. ([#340](https://github.com/pingidentity/terraform-provider-pingone/issues/340))
* `resource/pingone_application`: Expanded the native application documentation example for mobile app use case. ([#344](https://github.com/pingidentity/terraform-provider-pingone/issues/344))
* bump `github.com/bflad/tfproviderlint` v0.28.1 => v0.29.0 ([#347](https://github.com/pingidentity/terraform-provider-pingone/issues/347))
* bump `github.com/patrickcping/pingone-go-sdk-v2/agreementmanagement` v0.1.1 => v0.1.2 ([#345](https://github.com/pingidentity/terraform-provider-pingone/issues/345))
* bump `github.com/patrickcping/pingone-go-sdk-v2/authorize` v0.1.4 => v0.1.5 ([#345](https://github.com/pingidentity/terraform-provider-pingone/issues/345))
* bump `github.com/patrickcping/pingone-go-sdk-v2/management` v0.18.0 => v0.19.0 ([#345](https://github.com/pingidentity/terraform-provider-pingone/issues/345))
* bump `github.com/patrickcping/pingone-go-sdk-v2/mfa` v0.9.1 => v0.9.2 ([#345](https://github.com/pingidentity/terraform-provider-pingone/issues/345))
* bump `github.com/patrickcping/pingone-go-sdk-v2` v0.6.2 => v0.6.3 ([#345](https://github.com/pingidentity/terraform-provider-pingone/issues/345))
* bump `github.com/terraform-linters/tflint` v0.46.0 => v0.46.1 ([#347](https://github.com/pingidentity/terraform-provider-pingone/issues/347))

BUG FIXES:

* `resource/pingone_application`: Added a default value to optional `oidc_options.mobile_app.passcode_refresh_seconds` attribute. ([#344](https://github.com/pingidentity/terraform-provider-pingone/issues/344))
* `resource/pingone_application`: Fix for "There was an unexpected error with the service" error when attempting to create an application immediately after creation of the parent environment. ([#345](https://github.com/pingidentity/terraform-provider-pingone/issues/345))

## 0.12.0 (18 April 2023)

NOTES:

* Bump `github.com/golangci/golangci-lint` from 1.52.0 to 1.52.2 ([#334](https://github.com/pingidentity/terraform-provider-pingone/issues/334))
* Bump `github.com/hashicorp/terraform-plugin-framework` from 1.1.1 to 1.2.0 ([#334](https://github.com/pingidentity/terraform-provider-pingone/issues/334))
* Bump `github.com/hashicorp/terraform-plugin-go` from 0.14.3 to 0.15.0 ([#334](https://github.com/pingidentity/terraform-provider-pingone/issues/334))
* Bump `github.com/hashicorp/terraform-plugin-sdk/v2` from 2.25.0 to 2.26.1 ([#334](https://github.com/pingidentity/terraform-provider-pingone/issues/334))
* Bump `github.com/terraform-linters/tflint` from 0.45.0 to 0.46.0 ([#334](https://github.com/pingidentity/terraform-provider-pingone/issues/334))
* Updated the index documentation to refer to the more detailed getting started guide at [pingidentity.github.io/terraform-docs/](https://pingidentity.github.io/terraform-docs/) ([#309](https://github.com/pingidentity/terraform-provider-pingone/issues/309))
* `data-source/pingone_schema`: Migrated to plugin framework. ([#306](https://github.com/pingidentity/terraform-provider-pingone/issues/306))
* `resource/pingone_application_attribute_mapping`: Migrated to plugin framework. ([#329](https://github.com/pingidentity/terraform-provider-pingone/issues/329))
* `resource/pingone_application`: Update the documentation example for external link. ([#333](https://github.com/pingidentity/terraform-provider-pingone/issues/333))
* `resource/pingone_environment`: Code optimisations for default computed schema values. ([#335](https://github.com/pingidentity/terraform-provider-pingone/issues/335))
* `resource/pingone_identity_provider_attribute`: Migrated to plugin framework. ([#329](https://github.com/pingidentity/terraform-provider-pingone/issues/329))
* `resource/pingone_identity_provider_attribute`: Reformatted the social provider and external identity provider attribute reference documentation. ([#329](https://github.com/pingidentity/terraform-provider-pingone/issues/329))
* `resource/pingone_resource_attribute`: Migrated to plugin framework. ([#329](https://github.com/pingidentity/terraform-provider-pingone/issues/329))
* `resource/pingone_webhook`: Added link to "Subscription Action Types" api reference for a full list of configurable action types. ([#332](https://github.com/pingidentity/terraform-provider-pingone/issues/332))
* bump `github.com/patrickcping/pingone-go-sdk-v2/management` v0.17.1 => v0.18.0 ([#336](https://github.com/pingidentity/terraform-provider-pingone/issues/336))
* bump `github.com/patrickcping/pingone-go-sdk-v2` v0.6.1 => v0.6.2 ([#336](https://github.com/pingidentity/terraform-provider-pingone/issues/336))

ENHANCEMENTS:

* `resource/pingone_application_attribute_mapping`: Support the ability to override the `sub` and `saml_subject` core attributes. ([#329](https://github.com/pingidentity/terraform-provider-pingone/issues/329))
* `resource/pingone_application_attribute_mapping`: Supports the ability to set attribute level scopes and enabled/disabled status in the ID token and on the userinfo endpoint for OIDC applications. ([#329](https://github.com/pingidentity/terraform-provider-pingone/issues/329))
* `resource/pingone_application`: Provide support for certificate based authentication. ([#311](https://github.com/pingidentity/terraform-provider-pingone/issues/311))
* `resource/pingone_identity_provider_attribute`: Support the ability to override the `username` core attribute. ([#329](https://github.com/pingidentity/terraform-provider-pingone/issues/329))
* `resource/pingone_key`: Support custom CRL for keys with type `ISSUANCE`. ([#312](https://github.com/pingidentity/terraform-provider-pingone/issues/312))
* `resource/pingone_resource_attribute`: Support the ability to override the `sub` core attribute for custom resources. ([#329](https://github.com/pingidentity/terraform-provider-pingone/issues/329))

BUG FIXES:

* `resource/pingone_resource_attribute`: Fix error when deleting predefined OpenID Connect resource attribute.  Now resets the value back to the environment default. ([#329](https://github.com/pingidentity/terraform-provider-pingone/issues/329))

## 0.11.1 (20 March 2023)

NOTES:

* `resource/pingone_application_flow_policy_assignment`: Expanded the HCL example in the registry documentation. ([#304](https://github.com/pingidentity/terraform-provider-pingone/issues/304))

BUG FIXES:

* `resource/pingone_environment`: Fix for "inconsistent result" error when using implicitly defined default value (from the SDK client) for `pingone_environment` `region` attribute. ([#305](https://github.com/pingidentity/terraform-provider-pingone/issues/305))

## 0.11.0 (20 March 2023)

NOTES:

* `data-source/pingone_environment`: Migrated to plugin framework. ([#292](https://github.com/pingidentity/terraform-provider-pingone/issues/292))
* `resource/pingone_application`: Changed the `idp_signing_key_id` attribute in SAML apps to expect a computed value from the platform (P14C-47055) ([#296](https://github.com/pingidentity/terraform-provider-pingone/issues/296))
* `resource/pingone_application`: Deprecates the `idp_signing_key_id` attribute for new `idp_signing_key` block in SAML apps. ([#296](https://github.com/pingidentity/terraform-provider-pingone/issues/296))
* `resource/pingone_environment`: Migrated to plugin framework. ([#292](https://github.com/pingidentity/terraform-provider-pingone/issues/292))
* bump `github.com/golangci/golangci-lint` v1.51.2 => v1.52.0 ([#300](https://github.com/pingidentity/terraform-provider-pingone/issues/300))
* bump `github.com/patrickcping/pingone-go-sdk-v2/agreementmanagement` v0.1.0 => v0.1.1 ([#300](https://github.com/pingidentity/terraform-provider-pingone/issues/300))
* bump `github.com/patrickcping/pingone-go-sdk-v2/authorize` v0.1.3 => v0.1.4 ([#300](https://github.com/pingidentity/terraform-provider-pingone/issues/300))
* bump `github.com/patrickcping/pingone-go-sdk-v2/management` v0.16.0 => v0.17.1 ([#300](https://github.com/pingidentity/terraform-provider-pingone/issues/300))
* bump `github.com/patrickcping/pingone-go-sdk-v2/mfa` v0.9.0 => v0.9.1 ([#300](https://github.com/pingidentity/terraform-provider-pingone/issues/300))
* bump `github.com/patrickcping/pingone-go-sdk-v2` v0.6.0 => v0.6.1 ([#300](https://github.com/pingidentity/terraform-provider-pingone/issues/300))

FEATURES:

* **New Data Source:** `pingone_flow_policies` ([#277](https://github.com/pingidentity/terraform-provider-pingone/issues/277))
* **New Data Source:** `pingone_flow_policy` ([#277](https://github.com/pingidentity/terraform-provider-pingone/issues/277))
* **New Resource:** `pingone_application_flow_policy_assignment` ([#277](https://github.com/pingidentity/terraform-provider-pingone/issues/277))

ENHANCEMENTS:

* `resource/pingone_application`: Adds support for defining the signing algorithm to apply to assertion/response signing in SAML apps. ([#296](https://github.com/pingidentity/terraform-provider-pingone/issues/296))
* `resource/pingone_environment`: The `default_population` parameter and `default_population_id` attributes, when an environment is created from new, now align correctly with the platform's own Default population. ([#292](https://github.com/pingidentity/terraform-provider-pingone/issues/292))

## 0.10.0 (13 March 2023)

NOTES:

* Update `pingone_sign_on_policy` and `pingone_sign_on_policy_action` documentation example for the MFA action. ([#275](https://github.com/pingidentity/terraform-provider-pingone/issues/275))
* `data-source/pingone_population`: Use common `environment_id` link ID schema definition. ([#287](https://github.com/pingidentity/terraform-provider-pingone/issues/287))
* `data-source/pingone_populations`: Use common `environment_id` link ID schema definition. ([#287](https://github.com/pingidentity/terraform-provider-pingone/issues/287))
* `data-source/pingone_trusted_email_domain`: Use common `environment_id` link ID schema definition. ([#287](https://github.com/pingidentity/terraform-provider-pingone/issues/287))
* `resource/pingone_group_nesting`: Corrected the schema documentation. ([#276](https://github.com/pingidentity/terraform-provider-pingone/issues/276))
* `resource/pingone_notification_policy`: Use common `environment_id` link ID schema definition. ([#287](https://github.com/pingidentity/terraform-provider-pingone/issues/287))
* `resource/pingone_notification_settings_email`: Use common `environment_id` link ID schema definition. ([#287](https://github.com/pingidentity/terraform-provider-pingone/issues/287))
* `resource/pingone_trusted_email_address`: Use common `environment_id` link ID schema definition. ([#287](https://github.com/pingidentity/terraform-provider-pingone/issues/287))
* bump `github.com/hashicorp/terraform-plugin-docs` v0.13.0 => v0.14.1 ([#285](https://github.com/pingidentity/terraform-provider-pingone/issues/285))
* bump `github.com/patrickcping/pingone-go-sdk-v2/authorize` v0.1.2 => v0.1.3 ([#285](https://github.com/pingidentity/terraform-provider-pingone/issues/285))
* bump `github.com/patrickcping/pingone-go-sdk-v2/management` v0.15.0 => v0.16.0 ([#285](https://github.com/pingidentity/terraform-provider-pingone/issues/285))
* bump `github.com/patrickcping/pingone-go-sdk-v2/mfa` v0.8.0 => v0.9.0 ([#285](https://github.com/pingidentity/terraform-provider-pingone/issues/285))
* bump `github.com/patrickcping/pingone-go-sdk-v2` v0.5.3 => v0.6.0 ([#285](https://github.com/pingidentity/terraform-provider-pingone/issues/285))

FEATURES:

* **New Data Source:** `pingone_agreement` ([#278](https://github.com/pingidentity/terraform-provider-pingone/issues/278))
* **New Data Source:** `pingone_agreement_localization` ([#278](https://github.com/pingidentity/terraform-provider-pingone/issues/278))
* **New Data Source:** `pingone_environments` ([#284](https://github.com/pingidentity/terraform-provider-pingone/issues/284))
* **New Data Source:** `pingone_organization` ([#283](https://github.com/pingidentity/terraform-provider-pingone/issues/283))
* **New Resource:** `pingone_agreement` ([#278](https://github.com/pingidentity/terraform-provider-pingone/issues/278))
* **New Resource:** `pingone_agreement_enable` ([#278](https://github.com/pingidentity/terraform-provider-pingone/issues/278))
* **New Resource:** `pingone_agreement_localization` ([#278](https://github.com/pingidentity/terraform-provider-pingone/issues/278))
* **New Resource:** `pingone_agreement_localization_enable` ([#278](https://github.com/pingidentity/terraform-provider-pingone/issues/278))
* **New Resource:** `pingone_agreement_revision` ([#278](https://github.com/pingidentity/terraform-provider-pingone/issues/278))

BUG FIXES:

* `resource/pingone_application`: Fix a bug where `pkce_enforcement` couldn't be set on native application types. ([#282](https://github.com/pingidentity/terraform-provider-pingone/issues/282))
* `resource/pingone_application`: Fix input validation for mobile native uri values on `post_logout_redirect_uris`, `redirect_uris` and `target_link_uri` parameters. ([#282](https://github.com/pingidentity/terraform-provider-pingone/issues/282))

## 0.9.0 (23 February 2023)

NOTES:

* Added plugin mux factory and plugin framework (v6 protocol) provider to facilitate migration from SDKv2 (v5 protocol) ([#252](https://github.com/pingidentity/terraform-provider-pingone/issues/252))
* bump `github.com/golangci/golangci-lint` v1.51.1 => v1.51.2 ([#270](https://github.com/pingidentity/terraform-provider-pingone/issues/270))
* bump `github.com/hashicorp/go-getter` v1.6.2 => v1.7.0 ([#256](https://github.com/pingidentity/terraform-provider-pingone/issues/256))
* bump `github.com/hashicorp/terraform-plugin-mux` v0.8.0 => v0.9.0 ([#270](https://github.com/pingidentity/terraform-provider-pingone/issues/270))
* bump `github.com/hashicorp/terraform-plugin-sdk/v2` v2.24.1 => v2.25.0 ([#270](https://github.com/pingidentity/terraform-provider-pingone/issues/270))
* bump `github.com/patrickcping/pingone-go-sdk-v2/authorize` v0.1.1 => v0.1.2 ([#270](https://github.com/pingidentity/terraform-provider-pingone/issues/270))
* bump `github.com/patrickcping/pingone-go-sdk-v2/management` v0.14.0 => v0.15.0 ([#266](https://github.com/pingidentity/terraform-provider-pingone/issues/266))
* bump `github.com/patrickcping/pingone-go-sdk-v2/mfa` v0.7.2 => v0.8.0 ([#264](https://github.com/pingidentity/terraform-provider-pingone/issues/264))
* bump `github.com/patrickcping/pingone-go-sdk-v2` v0.5.2 => v0.5.3 ([#270](https://github.com/pingidentity/terraform-provider-pingone/issues/270))
* bump `golang.org/x/net` v0.5.0 => v0.7.0 ([#257](https://github.com/pingidentity/terraform-provider-pingone/issues/257))

FEATURES:

* **New Data Source:** `pingone_population` ([#255](https://github.com/pingidentity/terraform-provider-pingone/issues/255))
* **New Data Source:** `pingone_populations` ([#255](https://github.com/pingidentity/terraform-provider-pingone/issues/255))
* **New Data Source:** `pingone_trusted_email_domain` ([#253](https://github.com/pingidentity/terraform-provider-pingone/issues/253))
* **New Resource:** `pingone_notification_policy` ([#268](https://github.com/pingidentity/terraform-provider-pingone/issues/268))
* **New Resource:** `pingone_notification_settings_email` ([#269](https://github.com/pingidentity/terraform-provider-pingone/issues/269))
* **New Resource:** `pingone_trusted_email_address` ([#253](https://github.com/pingidentity/terraform-provider-pingone/issues/253))

ENHANCEMENTS:

* `resource/pingone_application`: Support for Huawei HMS push notification configuration. ([#264](https://github.com/pingidentity/terraform-provider-pingone/issues/264))
* `resource/pingone_gateway`: Now supports RADIUS gateways. ([#266](https://github.com/pingidentity/terraform-provider-pingone/issues/266))
* `resource/pingone_mfa_application_push_credential`: Support for Huawei HMS push notification configuration. ([#264](https://github.com/pingidentity/terraform-provider-pingone/issues/264))

## 0.8.1 (14 February 2023)

NOTES:

* `resource/pingone_application_role_assignment`: Updated documentation to add more examples and clarify the schema requirements. ([#247](https://github.com/pingidentity/terraform-provider-pingone/issues/247))
* `resource/pingone_gateway_role_assignment`: Updated documentation to add more examples and clarify the schema requirements. ([#247](https://github.com/pingidentity/terraform-provider-pingone/issues/247))
* `resource/pingone_role_assignment_user`: Updated documentation to add more examples and clarify the schema requirements. ([#247](https://github.com/pingidentity/terraform-provider-pingone/issues/247))
* bump `github.com/golangci/golangci-lint` v1.50.1 => v1.51.1 ([#242](https://github.com/pingidentity/terraform-provider-pingone/issues/242))
* bump `github.com/hashicorp/terraform-plugin-log` v0.7.0 => v0.8.0 ([#245](https://github.com/pingidentity/terraform-provider-pingone/issues/245))
* bump `github.com/terraform-linters/tflint` v0.44.1 => v0.45.0 ([#246](https://github.com/pingidentity/terraform-provider-pingone/issues/246))

ENHANCEMENTS:

* `resource/pingone_application`: Changed input validation to add support for localhost `http://` endpoints in OIDC applications. ([#244](https://github.com/pingidentity/terraform-provider-pingone/issues/244))

BUG FIXES:

* `resource/pingone_sign_on_policy`: Corrected input validation regex for the sign on policy `name` attribute. ([#248](https://github.com/pingidentity/terraform-provider-pingone/issues/248))

## 0.8.0 (12 January 2023)

NOTES:

* bump `github.com/patrickcping/pingone-go-sdk-v2/management` v0.13.0 => v0.14.0 ([#230](https://github.com/pingidentity/terraform-provider-pingone/issues/230))
* bump `github.com/patrickcping/pingone-go-sdk-v2/mfa` v0.7.1 => v0.7.2 ([#230](https://github.com/pingidentity/terraform-provider-pingone/issues/230))
* bump `github.com/patrickcping/pingone-go-sdk-v2` v0.5.1 => v0.5.2 ([#230](https://github.com/pingidentity/terraform-provider-pingone/issues/230))
* resource/pingone_mfa_settings: Deprecate attribute block `authentication` and attribute `authentication.device_selection` as device selection settings have moved to the `pingone_mfa_policy` resource. ([#230](https://github.com/pingidentity/terraform-provider-pingone/issues/230))

FEATURES:

* **New Resource:** `pingone_notification_template_content` ([#229](https://github.com/pingidentity/terraform-provider-pingone/issues/229))

ENHANCEMENTS:

* resource/pingone_mfa_policy: Support per application push notification timeout by adding optional parameter `push_timeout_duration` to the `mobile.application` block. ([#231](https://github.com/pingidentity/terraform-provider-pingone/issues/231))
* resource/pingone_mfa_policy: Support per policy device selection settings by adding the optional parameter `device_selection`, previously found on the `pingone_mfa_settings` resource. ([#230](https://github.com/pingidentity/terraform-provider-pingone/issues/230))

## 0.7.1 (09 January 2023)

NOTES:

* Removed documented reference to the Fraud service due to capability merge with Risk ([#224](https://github.com/pingidentity/terraform-provider-pingone/issues/224))
* bump `github.com/hashicorp/terraform-plugin-sdk/v2` v2.24.0 => v2.24.1 ([#214](https://github.com/pingidentity/terraform-provider-pingone/issues/214))
* bump `github.com/patrickcping/pingone-go-sdk-v2/authorize` v0.1.0 => v0.1.1 ([#228](https://github.com/pingidentity/terraform-provider-pingone/issues/228))
* bump `github.com/patrickcping/pingone-go-sdk-v2/management` v0.12.0 => v0.13.0 ([#228](https://github.com/pingidentity/terraform-provider-pingone/issues/228))
* bump `github.com/patrickcping/pingone-go-sdk-v2/mfa` v0.7.0 => v0.7.1 ([#228](https://github.com/pingidentity/terraform-provider-pingone/issues/228))
* bump `github.com/patrickcping/pingone-go-sdk-v2` v0.5.0 => v0.5.1 ([#228](https://github.com/pingidentity/terraform-provider-pingone/issues/228))
* bump `github.com/terraform-linters/tflint` v0.42.2 => v0.44.1 ([#221](https://github.com/pingidentity/terraform-provider-pingone/issues/221))
* resource/pingone_application: Removed redundant support for `tags` parameter on SAML type applications. ([#228](https://github.com/pingidentity/terraform-provider-pingone/issues/228))

ENHANCEMENTS:

* resource/pingone_application: Optional parameter `home_page_url` added to the SAML application options. ([#228](https://github.com/pingidentity/terraform-provider-pingone/issues/228))
* resource/pingone_application: Support better resiliency of rolling refresh tokens by adding the optional parameter `refresh_token_rolling_grace_period_duration` to the OIDC application options.  This is useful in the case of network errors on the client. ([#228](https://github.com/pingidentity/terraform-provider-pingone/issues/228))
* resource/pingone_application: Support options for post login redirect by adding the optional parameter `target_link_uri` to the OIDC application options. ([#228](https://github.com/pingidentity/terraform-provider-pingone/issues/228))
* resource/pingone_application: Support the ability to hide an application from the Application Portal through the new optional parameter `hidden_from_app_portal`. ([#228](https://github.com/pingidentity/terraform-provider-pingone/issues/228))
* resource/pingone_application: Support third party initiated login by adding the optional parameter `initiate_login_uri` to the OIDC application options. ([#228](https://github.com/pingidentity/terraform-provider-pingone/issues/228))
* resource/pingone_application: Support wildcards use in redirect URIs by adding the optional parameter `allow_wildcards_in_redirect_uris` to the OIDC application options. ([#228](https://github.com/pingidentity/terraform-provider-pingone/issues/228))
* resource/pingone_sign_on_policy_action: Added `last_sign_on_older_than_seconds_mfa` condition that can only be set to an MFA Sign on policy action. ([#225](https://github.com/pingidentity/terraform-provider-pingone/issues/225))

BUG FIXES:

* data-source/pingone_user: Fixed provider panic crash when the user cannot be found. ([#227](https://github.com/pingidentity/terraform-provider-pingone/issues/227))
* resource/pingone_sign_on_policy_action: Added `value_boolean` to the `user_attribute_equals` condition block as the existing `value` property didn't correctly interpret boolean values. ([#225](https://github.com/pingidentity/terraform-provider-pingone/issues/225))
* resource/pingone_sign_on_policy_action: Fixed bug where the `last_sign_on_older_than_seconds` condition, when set to an MFA Sign on policy action that was then changed in the console lead to a provider crash on next replan. ([#225](https://github.com/pingidentity/terraform-provider-pingone/issues/225))

## 0.7.0 (07 November 2022)

NOTES:

* bump `github.com/golangci/golangci-lint` v1.50.0 => v1.50.1 ([#207](https://github.com/pingidentity/terraform-provider-pingone/issues/207))
* bump `github.com/hashicorp/terraform-plugin-sdk/v2` v2.23.0 => v2.24.0 ([#198](https://github.com/pingidentity/terraform-provider-pingone/issues/198))
* bump `github.com/patrickcping/pingone-go-sdk-v2/management` v0.11.2 => v0.12.0 ([#207](https://github.com/pingidentity/terraform-provider-pingone/issues/207))
* bump `github.com/patrickcping/pingone-go-sdk-v2/mfa` v0.6.1 => v0.7.0 ([#207](https://github.com/pingidentity/terraform-provider-pingone/issues/207))
* bump `github.com/patrickcping/pingone-go-sdk-v2` v0.4.2 => v0.4.3 ([#207](https://github.com/pingidentity/terraform-provider-pingone/issues/207))
* bump `github.com/patrickcping/pingone-go-sdk-v2` v0.4.3 => v0.5.0 ([#208](https://github.com/pingidentity/terraform-provider-pingone/issues/208))
* bump `github.com/terraform-linters/tflint` v0.41.0 => v0.42.2 ([#204](https://github.com/pingidentity/terraform-provider-pingone/issues/204))

FEATURES:

* **New Data Source:** `pingone_resource_attribute` ([#205](https://github.com/pingidentity/terraform-provider-pingone/issues/205))
* **New Resource:** `pingone_branding_settings` ([#195](https://github.com/pingidentity/terraform-provider-pingone/issues/195))
* **New Resource:** `pingone_branding_theme` ([#195](https://github.com/pingidentity/terraform-provider-pingone/issues/195))
* **New Resource:** `pingone_image` ([#186](https://github.com/pingidentity/terraform-provider-pingone/issues/186))
* **New Resource:** `pingone_mfa_fido_policy` ([#194](https://github.com/pingidentity/terraform-provider-pingone/issues/194))
* **New Resource:** `pingone_resource_attribute` ([#205](https://github.com/pingidentity/terraform-provider-pingone/issues/205))
* **New Resource:** `pingone_resource_scope_openid` ([#205](https://github.com/pingidentity/terraform-provider-pingone/issues/205))
* **New Resource:** `pingone_resource_scope_pingone_api` ([#205](https://github.com/pingidentity/terraform-provider-pingone/issues/205))

ENHANCEMENTS:

* Optional parameter `api_access_token` added to the provider configuration, to allow use of a PingOne API access token obtained prior to execution. ([#208](https://github.com/pingidentity/terraform-provider-pingone/issues/208))
* data-source/pingone_resource: Added read only support for the `introspect_endpoint_auth_method` and `client_secret` attributes. ([#205](https://github.com/pingidentity/terraform-provider-pingone/issues/205))
* data-source/pingone_resource_scope: Added read only support for the `mapped_claims` attribute. ([#205](https://github.com/pingidentity/terraform-provider-pingone/issues/205))
* resource/pingone_environment: No longer forces re-creation of the environment resource if the license ID is changed. ([#206](https://github.com/pingidentity/terraform-provider-pingone/issues/206))
* resource/pingone_resource: Added support for the optional `introspect_endpoint_auth_method` and the computed `client_secret` attributes. ([#205](https://github.com/pingidentity/terraform-provider-pingone/issues/205))

BUG FIXES:

* data-source/pingone_licenses: Remove the value restriction on the license `package` field on when filtering.  Package values are not fixed and can change over time. ([#206](https://github.com/pingidentity/terraform-provider-pingone/issues/206))
* resource/pingone_resource: Removed the potential for defective management of PingOne API and OpenID Connect resources. ([#205](https://github.com/pingidentity/terraform-provider-pingone/issues/205))
* resource/pingone_resource_scope: Removed the potential for defective management of PingOne API and OpenID Connect resource scopes. ([#205](https://github.com/pingidentity/terraform-provider-pingone/issues/205))

## 0.6.1 (15 October 2022)

NOTES:

* bump `github.com/patrickcping/pingone-go-sdk-v2/management` v0.11.0 => v0.11.1 ([#181](https://github.com/pingidentity/terraform-provider-pingone/issues/181))
* bump `github.com/patrickcping/pingone-go-sdk-v2/management` v0.11.1 => v0.11.2 ([#187](https://github.com/pingidentity/terraform-provider-pingone/issues/187))
* bump `github.com/patrickcping/pingone-go-sdk-v2/mfa` v0.6.0 => v0.6.1 ([#181](https://github.com/pingidentity/terraform-provider-pingone/issues/181))
* bump `github.com/patrickcping/pingone-go-sdk-v2` v0.4.0 => v0.4.1 ([#181](https://github.com/pingidentity/terraform-provider-pingone/issues/181))
* bump `github.com/patrickcping/pingone-go-sdk-v2` v0.4.1 => v0.4.2 ([#187](https://github.com/pingidentity/terraform-provider-pingone/issues/187))

ENHANCEMENTS:

* resource/pingone_application: Add support for `universal_app_link` attribute for native mobile applications. ([#185](https://github.com/pingidentity/terraform-provider-pingone/issues/185))
* resource/pingone_application: Add support for integrity detection `excluded_platforms` attribute for native mobile applications. ([#185](https://github.com/pingidentity/terraform-provider-pingone/issues/185))

BUG FIXES:

* resource/pingone_mfa_settings: Made the `lockout` configuration block and `lockout.duration_seconds` optional in the schema. ([#181](https://github.com/pingidentity/terraform-provider-pingone/issues/181))
* resource/pingone_trusted_email_domain_dkim: Corrected documentation example. ([#184](https://github.com/pingidentity/terraform-provider-pingone/issues/184))
* resource/pingone_trusted_email_domain_spf: Corrected documentation example. ([#184](https://github.com/pingidentity/terraform-provider-pingone/issues/184))

## 0.6.0 (10 October 2022)

NOTES:

* Documentation: Organised registry documentation into subcategories ([#169](https://github.com/pingidentity/terraform-provider-pingone/issues/169))
* bump `github.com/golangci/golangci-lint` v1.49.0 => v1.50.0 ([#176](https://github.com/pingidentity/terraform-provider-pingone/issues/176))
* bump `github.com/patrickcping/pingone-go-sdk-v2/management` v0.10.0 => v0.11.0 ([#170](https://github.com/pingidentity/terraform-provider-pingone/issues/170))
* bump `github.com/patrickcping/pingone-go-sdk-v2/mfa` v0.5.1 => v0.6.0 ([#170](https://github.com/pingidentity/terraform-provider-pingone/issues/170))
* bump `github.com/patrickcping/pingone-go-sdk-v2` v0.3.8 => v0.4.0 ([#170](https://github.com/pingidentity/terraform-provider-pingone/issues/170))

FEATURES:

* **New Data Source:** `pingone_language` ([#162](https://github.com/pingidentity/terraform-provider-pingone/issues/162))
* **New Data Source:** `pingone_license` ([#164](https://github.com/pingidentity/terraform-provider-pingone/issues/164))
* **New Data Source:** `pingone_licenses` ([#164](https://github.com/pingidentity/terraform-provider-pingone/issues/164))
* **New Data Source:** `pingone_user` ([#168](https://github.com/pingidentity/terraform-provider-pingone/issues/168))
* **New Data Source:** `pingone_users` ([#168](https://github.com/pingidentity/terraform-provider-pingone/issues/168))
* **New Resource:** `pingone_authorize_decision_endpoint` ([#160](https://github.com/pingidentity/terraform-provider-pingone/issues/160))
* **New Resource:** `pingone_language` ([#162](https://github.com/pingidentity/terraform-provider-pingone/issues/162))
* **New Resource:** `pingone_language_update` ([#162](https://github.com/pingidentity/terraform-provider-pingone/issues/162))
* **New Resource:** `pingone_mfa_application_push_credential` ([#170](https://github.com/pingidentity/terraform-provider-pingone/issues/170))
* **New Resource:** `pingone_mfa_policy` ([#170](https://github.com/pingidentity/terraform-provider-pingone/issues/170))

BUG FIXES:

* resource/pingone_role_assignment_user: Corrected import command on registry documentation. ([#172](https://github.com/pingidentity/terraform-provider-pingone/issues/172))

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
