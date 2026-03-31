resource "pingone_davinci_connector_instance" "pingOneFormsConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "pingOneFormsConnector"
  }
  name = "My awesome pingOneFormsConnector"
  property {
    name  = "agreement"
    type  = "string"
    value = var.pingoneformsconnector_property_agreement
  }
  property {
    name  = "agreementId"
    type  = "string"
    value = var.pingoneformsconnector_property_agreement_id
  }
  property {
    name  = "agreementSectionLabel"
    type  = "string"
    value = var.pingoneformsconnector_property_agreement_section_label
  }
  property {
    name  = "authenticationMethodList"
    type  = "string"
    value = var.pingoneformsconnector_property_authentication_method_list
  }
  property {
    name  = "authenticationMethodSource"
    type  = "string"
    value = var.pingoneformsconnector_property_authentication_method_source
  }
  property {
    name  = "buttonText"
    type  = "string"
    value = var.pingoneformsconnector_property_button_text
  }
  property {
    name  = "challenge"
    type  = "string"
    value = var.pingoneformsconnector_property_challenge
  }
  property {
    name  = "collectBehavioralData"
    type  = "string"
    value = var.pingoneformsconnector_property_collect_behavioral_data
  }
  property {
    name  = "componentVisibility"
    type  = "string"
    value = var.pingoneformsconnector_property_component_visibility
  }
  property {
    name  = "deviceProfilingSectionLabel"
    type  = "string"
    value = var.pingoneformsconnector_property_device_profiling_section_label
  }
  property {
    name  = "dynamicText"
    type  = "string"
    value = var.pingoneformsconnector_property_dynamic_text
  }
  property {
    name  = "enableMagicLinkAuthentication"
    type  = "string"
    value = var.pingoneformsconnector_property_enable_magic_link_authentication
  }
  property {
    name  = "enablePolling"
    type  = "string"
    value = var.pingoneformsconnector_property_enable_polling
  }
  property {
    name  = "enableRisk"
    type  = "string"
    value = var.pingoneformsconnector_property_enable_risk
  }
  property {
    name  = "form"
    type  = "string"
    value = var.pingoneformsconnector_property_form
  }
  property {
    name  = "formData"
    type  = "string"
    value = var.pingoneformsconnector_property_form_data
  }
  property {
    name  = "innerCss"
    type  = "string"
    value = var.pingoneformsconnector_property_inner_css
  }
  property {
    name  = "innerHtml"
    type  = "string"
    value = var.pingoneformsconnector_property_inner_html
  }
  property {
    name  = "isIAFDetectionEnabled"
    type  = "string"
    value = var.pingoneformsconnector_property_is_iafdetection_enabled
  }
  property {
    name  = "linkBrandingThemesUrl"
    type  = "string"
    value = var.pingoneformsconnector_property_link_branding_themes_url
  }
  property {
    name  = "linkFormsUrl"
    type  = "string"
    value = var.pingoneformsconnector_property_link_forms_url
  }
  property {
    name  = "linkMFAPolicies"
    type  = "string"
    value = var.pingoneformsconnector_property_link_mfapolicies
  }
  property {
    name  = "linkWithP1User"
    type  = "string"
    value = var.pingoneformsconnector_property_link_with_p1_user
  }
  property {
    name  = "message"
    type  = "string"
    value = var.pingoneformsconnector_property_message
  }
  property {
    name  = "mfaPolicyId"
    type  = "string"
    value = var.pingoneformsconnector_property_mfa_policy_id
  }
  property {
    name  = "nextEvent"
    type  = "string"
    value = var.pingoneformsconnector_property_next_event
  }
  property {
    name  = "pingidAgent"
    type  = "string"
    value = var.pingoneformsconnector_property_pingid_agent
  }
  property {
    name  = "pingidAgentPort"
    type  = "string"
    value = var.pingoneformsconnector_property_pingid_agent_port
  }
  property {
    name  = "pingidAgentTimeout"
    type  = "string"
    value = var.pingoneformsconnector_property_pingid_agent_timeout
  }
  property {
    name  = "pollChallengeStatus"
    type  = "string"
    value = var.pingoneformsconnector_property_poll_challenge_status
  }
  property {
    name  = "pollInterval"
    type  = "string"
    value = var.pingoneformsconnector_property_poll_interval
  }
  property {
    name  = "pollRetries"
    type  = "string"
    value = var.pingoneformsconnector_property_poll_retries
  }
  property {
    name  = "pollingSectionLabel"
    type  = "string"
    value = var.pingoneformsconnector_property_polling_section_label
  }
  property {
    name  = "population"
    type  = "string"
    value = var.pingoneformsconnector_property_population
  }
  property {
    name  = "populationId"
    type  = "string"
    value = var.pingoneformsconnector_property_population_id
  }
  property {
    name  = "publicKeyCredentialCreationOptions"
    type  = "string"
    value = var.pingoneformsconnector_property_public_key_credential_creation_options
  }
  property {
    name  = "publicKeyCredentialRequestOptions"
    type  = "string"
    value = var.pingoneformsconnector_property_public_key_credential_request_options
  }
  property {
    name  = "qrCodeContents"
    type  = "string"
    value = var.pingoneformsconnector_property_qr_code_contents
  }
  property {
    name  = "qrCodeSectionLabel"
    type  = "string"
    value = var.pingoneformsconnector_property_qr_code_section_label
  }
  property {
    name  = "returnUrl"
    type  = "string"
    value = var.pingoneformsconnector_property_return_url
  }
  property {
    name  = "returnUrlLabel"
    type  = "string"
    value = var.pingoneformsconnector_property_return_url_label
  }
  property {
    name  = "sectionLabelFido2"
    type  = "string"
    value = var.pingoneformsconnector_property_section_label_fido2
  }
  property {
    name  = "showContinueButton"
    type  = "string"
    value = var.pingoneformsconnector_property_show_continue_button
  }
  property {
    name  = "theme"
    type  = "string"
    value = var.pingoneformsconnector_property_theme
  }
  property {
    name  = "themeId"
    type  = "string"
    value = var.pingoneformsconnector_property_theme_id
  }
  property {
    name  = "universalDeviceIdentification"
    type  = "string"
    value = var.pingoneformsconnector_property_universal_device_identification
  }
  property {
    name  = "userId"
    type  = "string"
    value = var.pingoneformsconnector_property_user_id
  }
}
