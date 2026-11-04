### NEW DATA SOURCES

- `pingone_key_rotation_policy`

### ENHANCEMENTS

- `pingone_notification_settings_email` resource now supports OAuth 2.0 authentication (`auth_url`, `grant_type`, `client_id`, `client_secret`, `assertion`, `scopes`) and custom header authentication (`header_name`, `header_value`) for custom email providers.
- `resource/pingone_environment`: Added support for `Recognize` to the `services.type` attribute.
- `data-source/pingone_environment`: Added support for `Recognize` to the `services.type` attribute.
- `resource/pingone_application`: Added `include_typ` and `skip_token_session_validation` attributes to the `oidc_options` block, added support for the `AI_AGENT` application type, added `external_id` attribute to the `external_link_options` block, added an `enable_always_accept_acs_url_in_signed_authn_request` attribute to the `saml_options` block, and added `hybrid_join_device_attribute`, `sts_token_type`, `use_mex13`, `user_attribute_username_match`, `user_attribute_x509_match`, `valid_issuance_certificates`, `assertion_duration`, and `ws_trust_version` attributes to the `wsfed_options` block.
- `data-source/pingone_application`: Added `include_typ` and `skip_token_session_validation` attributes to the `oidc_options` block, added `external_id` attribute to the `external_link_options` block, added an `enable_always_accept_acs_url_in_signed_authn_request` attribute to the `saml_options` block, and added `hybrid_join_device_attribute`, `sts_token_type`, `use_mex13`, `user_attribute_username_match`, `user_attribute_x509_match`, `valid_issuance_certificates`, `assertion_duration`, and `ws_trust_version` attributes to the `wsfed_options` block.
- `resource/pingone_system_application`: Added `logout_type` and `manage_devices_via_my_account_enabled` attributes.
- `data-source/pingone_system_application`: Added `logout_type` and `manage_devices_via_my_account_enabled` attributes.
- `resource/pingone_application`: Added `apple_team_id`, `apple_app_attest_settings`, `ping_id_settings`, the PingID mobile app configuration settings (`allow_authentication_device_locked`, `allow_un_pair_and_change_device_from_app`, `mobile_app_security_pin_require`, `mobile_device_biometrics_setting`, `notify_optional_mobile_version_updates`, `notify_require_mobile_version_updates`, `otp_push_notification_enabled`, `show_authentication_information`), and a `push` block (with `number_matching.type`) to the `oidc_options.mobile_app` block.
- `data-source/pingone_application`: Added `apple_team_id`, `apple_app_attest_settings`, `ping_id_settings`, the PingID mobile app configuration settings (`allow_authentication_device_locked`, `allow_un_pair_and_change_device_from_app`, `mobile_app_security_pin_require`, `mobile_device_biometrics_setting`, `notify_optional_mobile_version_updates`, `notify_require_mobile_version_updates`, `otp_push_notification_enabled`, `show_authentication_information`), and a `push` block (with `number_matching.type`) to the `oidc_options.mobile_app` block.
- `resource/pingone_user_role_assignment`: Added a `scope_group_id` attribute to support group-scoped admin role assignments.
- `resource/pingone_group_role_assignment`: Added a `scope_group_id` attribute to support group-scoped admin role assignments.
- `resource/pingone_role_assignment_user`: Added a `scope_group_id` attribute to support group-scoped admin role assignments.
- `resource/pingone_custom_role`: Added `GROUP` as a valid `applicable_to` value.
- `data-source/pingone_role`: Added `Advanced Help Desk Admin`, `Advanced Identity Cloud Tenant Auditor`, `PingOne Privilege Administrator`, `PingOne Privilege Read-Only Administrator`, and `Promotion Admin` as valid `name` values.

### BUG FIXES

- `resource/pingone_risk_policy`: fixed a `Provider produced inconsistent result after apply` error when configuring `mitigations` with an empty list. An empty list is now rejected at plan time with a clear validation error, matching the behavior of `overrides`.
- `resource/pingone_application`: Removed validation that was disallowing the `signing` attribute for OIDC `SERVICE` applications.
- Fixed an issue where schema validation would reject valid URLs containing percent-encoded characters.

### NOTES

- Several attributes that hold secret values are now marked as sensitive: the `oauth.client_secret` and `api_key.value` attributes of the `pingone_davinci_application` and `pingone_davinci_applications` data sources, the `properties` attribute of the `pingone_davinci_connector_instance` and `pingone_davinci_connector_instances` data sources, and the `http_endpoint_headers` and `connection_details_headers` attributes of the `pingone_webhook` resource.
- The provider's `client_secret` and `api_access_token` attributes are now marked as sensitive.
- bump `github.com/pingidentity/pingone-go-sdk-legacy` v0.17.0 => v0.18.0
- bump `github.com/pingidentity/pingone-go-sdk-legacy/management` v0.79.0 => v0.83.0

