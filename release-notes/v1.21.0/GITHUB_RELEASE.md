### ENHANCEMENTS

[76447246](https://github.com/pingidentity/terraform-provider-pingone/commit/76447246) `resource/pingone_risk_policy`: Added `mitigations`, `fallback`, and `targets` attributes to support mitigation-style risk policies (mutually exclusive with `overrides`) [#1327](https://github.com/pingidentity/terraform-provider-pingone/pull/1327)
[8e3bae60](https://github.com/pingidentity/terraform-provider-pingone/commit/8e3bae60) `resource/pingone_application`: Added support for `oidc_options.mobile_app.passcode_grace_period` [#1341](https://github.com/pingidentity/terraform-provider-pingone/pull/1341)
[8e3bae60](https://github.com/pingidentity/terraform-provider-pingone/commit/8e3bae60) `data-source/pingone_application`: Added support for `oidc_options.mobile_app.passcode_grace_period` [#1341](https://github.com/pingidentity/terraform-provider-pingone/pull/1341)
[1cd4a01f](https://github.com/pingidentity/terraform-provider-pingone/commit/1cd4a01f) `resource/pingone_application`: Added `oidc_options.signing.key_rotation_policy_id` attribute to allow configuration of the signing key rotation policy, fixing an issue where updates reset it to the PingOne default. [#1343](https://github.com/pingidentity/terraform-provider-pingone/pull/1343)
[df22bdaa](https://github.com/pingidentity/terraform-provider-pingone/commit/df22bdaa) `resource/pingone_mfa_device_policy`: Added `policy_type`, `desktop`, `yubikey`, and `oath_token` device method attributes, and the mobile application `biometrics_enabled`, `ip_pairing_configuration`, and `new_request_duration_configuration` sub-fields, for parity with `pingone_mfa_device_policy_default` [#1344](https://github.com/pingidentity/terraform-provider-pingone/pull/1344)
[f3a06383](https://github.com/pingidentity/terraform-provider-pingone/commit/f3a06383) `resource/pingone_davinci_flow`: Added support for the `graph_data.elements.nodes.%.data.outcomes` field [#1345](https://github.com/pingidentity/terraform-provider-pingone/pull/1345)

### BUG FIXES

[e2a137f3](https://github.com/pingidentity/terraform-provider-pingone/commit/e2a137f3) `resource/pingone_davinci_flow`: Fixed a panic when creating a flow with only the name attribute configured. [#1340](https://github.com/pingidentity/terraform-provider-pingone/pull/1340)
[df22bdaa](https://github.com/pingidentity/terraform-provider-pingone/commit/df22bdaa) `resource/pingone_mfa_device_policy_default`: Removed unsupported `yubikey.pairing_key_lifetime` and `oath_token.pairing_key_lifetime` attributes [#1344](https://github.com/pingidentity/terraform-provider-pingone/pull/1344)

### NOTES

[2d9b4a25](https://github.com/pingidentity/terraform-provider-pingone/commit/2d9b4a25) Update Connector Reference Guide (01 June 2026). [#1320](https://github.com/pingidentity/terraform-provider-pingone/pull/1320)

