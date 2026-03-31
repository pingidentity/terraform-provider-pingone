resource "pingone_davinci_connector_instance" "securIdConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "securIdConnector"
  }
  name = "My awesome securIdConnector"
  property {
    name  = "apiUrl"
    type  = "string"
    value = var.securidconnector_property_api_url
  }
  property {
    name  = "assurancePolicyId"
    type  = "string"
    value = var.securidconnector_property_assurance_policy_id
  }
  property {
    name  = "authnAttemptTimeout"
    type  = "string"
    value = var.securidconnector_property_authn_attempt_timeout
  }
  property {
    name  = "clientKey"
    type  = "string"
    value = var.securidconnector_property_client_key
  }
  property {
    name  = "htmlConfig0"
    type  = "string"
    value = var.securidconnector_property_html_config0
  }
  property {
    name  = "htmlConfig1"
    type  = "string"
    value = var.securidconnector_property_html_config1
  }
  property {
    name  = "htmlConfig2"
    type  = "string"
    value = var.securidconnector_property_html_config2
  }
  property {
    name  = "keepAttempt"
    type  = "string"
    value = var.securidconnector_property_keep_attempt
  }
  property {
    name  = "subjectName"
    type  = "string"
    value = var.securidconnector_property_subject_name
  }
}
