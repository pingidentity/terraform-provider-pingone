### BREAKING CHANGES

- `resource/pingone_verify_policy`: Removed the `voice` attribute which is no longer supported by the PingOne API.
- `data-source/pingone_verify_policy`: Removed the `voice` attribute which is no longer supported by the PingOne API.
- `resource/pingone_verify_voice_phrase`: The previously deprecated `pingone_verify_voice_phrase` resource has been removed as the underlying API is no longer supported by the PingOne API.
- `resource/pingone_verify_voice_phrase_content`: The previously deprecated `pingone_verify_voice_phrase_content` resource has been removed as the underlying API is no longer supported by the PingOne API.
- `data-source/pingone_verify_voice_phrase`: The previously deprecated `pingone_verify_voice_phrase` data source has been removed as the underlying API is no longer supported by the PingOne API.
- `data-source/pingone_verify_voice_phrase_content`: The previously deprecated `pingone_verify_voice_phrase_content` data source has been removed as the underlying API is no longer supported by the PingOne API.
- `data-source/pingone_verify_voice_phrase_contents`: The previously deprecated `pingone_verify_voice_phrase_contents` data source has been removed as the underlying API is no longer supported by the PingOne API.

### ENHANCEMENTS

- `resource/pingone_mfa_fido2_policy`: Added a `public_key_credential_hints` attribute to provide public key credential hints to the browser during FIDO registration.
- `resource/pingone_gateway`: Added the `blast_radius_mitigation` object to the `radius_clients` blocks, supporting the `require_msg_auth` and `limit_proxy_state` RADIUS gateway API fields.
- `data-source/pingone_gateway`: Added the `blast_radius_mitigation` object to the `radius_clients` blocks, supporting the `require_msg_auth` and `limit_proxy_state` RADIUS gateway API fields.

### BUG FIXES

- `resource/pingone_language_update`: Fix updates to an already-default language, allow state-only removal of default English, and provide an actionable error when unsetting default English.
- `resource/pingone_notification_template_content`: Return a clear error when importing or managing the read-only default template content.
- `resource/pingone_notification_settings_email`: Fixed an issue where the `requests.body` attribute rejected the documented `${to}` and `${message}` variable formats with a transform suffix of `json`, `html`, `url`, or `base64` (for example, `${to|url}`).

