### NEW RESOURCES

- pingone_application_metadata

### NEW DATA SOURCES

- `pingone_credential_type_versions`
- `pingone_credential_type_version`
- `pingone_phone_delivery_settings`
- `pingone_risk_policy`

### ENHANCEMENTS

- `resource/pingone_branding_theme`: Added configuration attributes to align with the API: `application_background_color`, `background_outline_color`, `body_text_size`, `body_text_weight`, `button_border_color`, `button_border_width`, `button_corner_radius`, `button_hover_state_border_color`, `button_hover_state_fill_color`, `button_hover_state_text_color`, `button_text_size`, `button_text_weight`, `card_border_color`, `card_border_width`, `card_corner_radius`, `card_horizontal_alignment`, `card_logo_alignment`, `card_shadow`, `card_vertical_alignment`, `focus_rectangle_color`, `footer_localized`, `foreground_highlight_color`, `foreground_main_color`, `global_font`, `header`, `header_background_color`, `header_localized`, `input_border_width`, `input_box_border_color`, `input_corner_radius`, `input_label_position`, `input_label_text_color`, `input_label_text_size`, `input_label_text_weight`, `input_value_text_color`, `input_value_text_size`, `input_value_text_weight`, `link_text_hover_color`, `link_text_size`, `link_text_weight`, `logo_height`, `sub_title_text_color`, `sub_title_text_size`, `sub_title_text_weight`, `title_text_color`, `title_text_size`, `title_text_weight`
- `resource/pingone_credential_type`: Added computed `version` attribute (`id`, `number`, `uri`)
- `data-source/pingone_credential_type`: Added computed `version` attribute (`id`, `number`, `uri`)
- `resource/pingone_phone_delivery_settings`: The `provider_custom.authentication` block now supports the `OAUTH2` and `CUSTOM_HEADER` authentication methods, with new parameters `auth_url`, `grant_type`, `assertion`, `client_id`, `client_secret`, `scopes`, `header_name`, and `header_value`, and a read-only `client_authentication_method` parameter.
- Added the `mapped_claims` and `enable_mapped_claims` attributes to the `pingone_resource_scope` resource, and the `enable_mapped_claims` attribute to the `pingone_resource_scope_openid` resource, to manage mapped claims on resource scopes.
- Added the `enable_mapped_claims` attribute to the `pingone_resource_scope` data source.
- `resource/pingone_environment`: Added support for `Privilege` and `AdvancedIdentityCloud` to the `services.type` attribute. Also update the `tags` attribute to support the `AUTHENTICATION_MODE_AGENT` and `AUTHENTICATION_MODE_AGENTLESS` tags used by `Privilege`.
- `data-source/pingone_environment`: Added support for `Privilege` and `AdvancedIdentityCloud` to the `services.type` attribute. Also update the `tags` attribute to support the `AUTHENTICATION_MODE_AGENT` and `AUTHENTICATION_MODE_AGENTLESS` tags used by `Privilege`.
- `resource/pingone_notification_template_content`: Updated enum for `template_name` attribute for template names now supported by the API: `account_created`, `account_updated`, `password_change_admin`, `password_change_user`, `password_recovery_new`, and `verification_code_new`.
- `resource/pingone_mfa_device_policy`: Added `block_disabled_users` and `block_users_with_disabled_mfa` attributes to align with API.
- `resource/pingone_mfa_device_policy_default`: Added `block_disabled_users` and `block_users_with_disabled_mfa` attributes to align with API.

### BUG FIXES

- `resource/pingone_sign_on_policy_action`: The `identifier_first.discovery_rule` block is now an ordered list. Previously the block was an unordered set, so rule precedence could not be controlled from the Terraform configuration. Existing configurations with multiple discovery rules may see a one-time in-place update on their next apply, changing the effective rule evaluation order in PingOne. Unlike a set, a list does not deduplicate identical `discovery_rule` blocks; duplicate rules are sent to PingOne as configured.
- Fixed an issue where some valid ISO 3166 country codes were rejected by provider validation.

### NOTES

- bump `github.com/patrickcping/pingone-go-sdk-v2/management` v0.70.0 => v0.71.0
- bump `github.com/patrickcping/pingone-go-sdk-v2/credentials` v0.12.1 => v0.13.0
- bump `github.com/patrickcping/pingone-go-sdk-v2/management` v0.72.0 => v0.73.0
- bump `github.com/patrickcping/pingone-go-sdk-v2` v0.14.14 => v0.15.0
- bump `github.com/patrickcping/pingone-go-sdk-v2/management` v0.74.0 => v0.76.0
- bump `github.com/patrickcping/pingone-go-sdk-v2/mfa` v0.25.1 => v0.26.0
- Update Connector Reference Guide (01 October 2026).

