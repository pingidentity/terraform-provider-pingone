### BREAKING CHANGES

[04dd9f04](https://github.com/pingidentity/terraform-provider-pingone/commit/04dd9f04) `resource/pingone_verify_policy`: Removed the `government_id.aadhaar.otp` attribute which is no longer supported by the PingOne API. [#1213](https://github.com/pingidentity/terraform-provider-pingone/pull/1213)
[04dd9f04](https://github.com/pingidentity/terraform-provider-pingone/commit/04dd9f04) `data-source/pingone_verify_policy`: Removed the `government_id.aadhaar.otp` attribute which is no longer supported by the PingOne API. [#1213](https://github.com/pingidentity/terraform-provider-pingone/pull/1213)

### NEW DATA SOURCES

[3de133af](https://github.com/pingidentity/terraform-provider-pingone/commit/3de133af) `pingone_password_policies` [#1212](https://github.com/pingidentity/terraform-provider-pingone/pull/1212)
[99a0c28f](https://github.com/pingidentity/terraform-provider-pingone/commit/99a0c28f) `pingone_system_application` [#1214](https://github.com/pingidentity/terraform-provider-pingone/pull/1214)
[75c45bae](https://github.com/pingidentity/terraform-provider-pingone/commit/75c45bae) `pingone_system_applications` [#1218](https://github.com/pingidentity/terraform-provider-pingone/pull/1218)
[8a19b47f](https://github.com/pingidentity/terraform-provider-pingone/commit/8a19b47f) `pingone_application_role_assignments` [#1219](https://github.com/pingidentity/terraform-provider-pingone/pull/1219)
[f0fce53a](https://github.com/pingidentity/terraform-provider-pingone/commit/f0fce53a) `pingone_group_role_assignments` [#1222](https://github.com/pingidentity/terraform-provider-pingone/pull/1222)
[6ed32689](https://github.com/pingidentity/terraform-provider-pingone/commit/6ed32689) `pingone_schema_attribute` [#1225](https://github.com/pingidentity/terraform-provider-pingone/pull/1225)

### ENHANCEMENTS

[6ed32689](https://github.com/pingidentity/terraform-provider-pingone/commit/6ed32689) `resource/pingone_schema_attribute`: Add support for managing imported `STANDARD` schema attributes [#1225](https://github.com/pingidentity/terraform-provider-pingone/pull/1225)
[6ed32689](https://github.com/pingidentity/terraform-provider-pingone/commit/6ed32689) `resource/pingone_schema_attribute`: Remove data loss protection from `unique`, allowing in-place updates (no resource replacement required) [#1225](https://github.com/pingidentity/terraform-provider-pingone/pull/1225)

### BUG FIXES

[af4463b9](https://github.com/pingidentity/terraform-provider-pingone/commit/af4463b9) Fixed a potential panic that could occur when the `PINGONE_ENVIRONMENT_ID` environment variable was set to a different value than the `environment_id` value defined in the provider configuration block. This bug only affected beta DaVinci resources. [#1220](https://github.com/pingidentity/terraform-provider-pingone/pull/1220)

### NOTES

[0c8bfa05](https://github.com/pingidentity/terraform-provider-pingone/commit/0c8bfa05) bump `github.com/hashicorp/terraform-plugin-sdk/v2` v2.38.1 => v2.38.2 [#1206](https://github.com/pingidentity/terraform-provider-pingone/pull/1206)
[0c8bfa05](https://github.com/pingidentity/terraform-provider-pingone/commit/0c8bfa05) bump `github.com/go-git/go-git/v5` v5.14.0 => v5.16.5 [#1209](https://github.com/pingidentity/terraform-provider-pingone/pull/1209)

