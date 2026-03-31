resource "pingone_davinci_connector_instance" "pingOneCredentialsConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "pingOneCredentialsConnector"
  }
  name = "My awesome pingOneCredentialsConnector"
  property {
    name  = "applicationInstance"
    type  = "string"
    value = var.pingonecredentialsconnector_property_application_instance
  }
  property {
    name  = "applyIssue"
    type  = "string"
    value = var.pingonecredentialsconnector_property_apply_issue
  }
  property {
    name  = "applyRevoke"
    type  = "string"
    value = var.pingonecredentialsconnector_property_apply_revoke
  }
  property {
    name  = "applyUpdate"
    type  = "string"
    value = var.pingonecredentialsconnector_property_apply_update
  }
  property {
    name  = "clientId"
    type  = "string"
    value = var.pingone_worker_app_client_id
  }
  property {
    name  = "clientSecret"
    type  = "string"
    value = var.pingone_worker_app_client_secret
  }
  property {
    name  = "credentialId"
    type  = "string"
    value = var.pingonecredentialsconnector_property_credential_id
  }
  property {
    name  = "credentialTypeId"
    type  = "string"
    value = var.pingonecredentialsconnector_property_credential_type_id
  }
  property {
    name  = "credentialsVerificationId"
    type  = "string"
    value = var.pingonecredentialsconnector_property_credentials_verification_id
  }
  property {
    name  = "didMethod"
    type  = "string"
    value = var.pingonecredentialsconnector_property_did_method
  }
  property {
    name  = "digitalWalletApplicationId"
    type  = "string"
    value = var.pingonecredentialsconnector_property_digital_wallet_application_id
  }
  property {
    name  = "digitalWalletId"
    type  = "string"
    value = var.pingonecredentialsconnector_property_digital_wallet_id
  }
  property {
    name  = "digitalWalletStatus"
    type  = "string"
    value = var.pingonecredentialsconnector_property_digital_wallet_status
  }
  property {
    name  = "envId"
    type  = "string"
    value = var.pingone_worker_app_environment_id
  }
  property {
    name  = "expirationDate"
    type  = "string"
    value = var.pingonecredentialsconnector_property_expiration_date
  }
  property {
    name  = "issuanceRuleId"
    type  = "string"
    value = var.pingonecredentialsconnector_property_issuance_rule_id
  }
  property {
    name  = "issuerFilterDids"
    type  = "string"
    value = var.pingonecredentialsconnector_property_issuer_filter_dids
  }
  property {
    name  = "issuerFilterEnvIds"
    type  = "string"
    value = var.pingonecredentialsconnector_property_issuer_filter_env_ids
  }
  property {
    name  = "message"
    type  = "string"
    value = var.pingonecredentialsconnector_property_message
  }
  property {
    name  = "notificationMethods"
    type  = "string"
    value = var.pingonecredentialsconnector_property_notification_methods
  }
  property {
    name  = "protocol"
    type  = "string"
    value = var.pingonecredentialsconnector_property_protocol
  }
  property {
    name  = "region"
    type  = "string"
    value = var.pingonecredentialsconnector_property_region
  }
  property {
    name  = "reqCredType"
    type  = "string"
    value = var.pingonecredentialsconnector_property_req_cred_type
  }
  property {
    name  = "requestedCredKeys"
    type  = "string"
    value = var.pingonecredentialsconnector_property_requested_cred_keys
  }
  property {
    name  = "requestedCredentials"
    type  = "string"
    value = var.pingonecredentialsconnector_property_requested_credentials
  }
  property {
    name  = "showPoweredBy"
    type  = "string"
    value = var.pingonecredentialsconnector_property_show_powered_by
  }
  property {
    name  = "skipButtonPress"
    type  = "string"
    value = var.pingonecredentialsconnector_property_skip_button_press
  }
  property {
    name  = "templateLocale"
    type  = "string"
    value = var.pingonecredentialsconnector_property_template_locale
  }
  property {
    name  = "templateVariables"
    type  = "string"
    value = var.pingonecredentialsconnector_property_template_variables
  }
  property {
    name  = "templateVariant"
    type  = "string"
    value = var.pingonecredentialsconnector_property_template_variant
  }
  property {
    name  = "userId"
    type  = "string"
    value = var.pingonecredentialsconnector_property_user_id
  }
}
