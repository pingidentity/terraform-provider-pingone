resource "pingone_davinci_connector_instance" "pingOneSnodeRiskConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "pingOneSnodeRiskConnector"
  }
  name = "My awesome pingOneSnodeRiskConnector"
  property {
    name  = "cookie"
    type  = "string"
    value = var.pingonesnoderiskconnector_property_cookie
  }
  property {
    name  = "customAttributes"
    type  = "string"
    value = var.pingonesnoderiskconnector_property_custom_attributes
  }
  property {
    name  = "envId"
    type  = "string"
    value = var.pingonesnoderiskconnector_property_env_id
  }
  property {
    name  = "externalId"
    type  = "string"
    value = var.pingonesnoderiskconnector_property_external_id
  }
  property {
    name  = "flowType"
    type  = "string"
    value = var.pingonesnoderiskconnector_property_flow_type
  }
  property {
    name  = "outcomes"
    type  = "string"
    value = var.pingonesnoderiskconnector_property_outcomes
  }
  property {
    name  = "riskId"
    type  = "string"
    value = var.pingonesnoderiskconnector_property_risk_id
  }
  property {
    name  = "riskPolicySetId"
    type  = "string"
    value = var.pingonesnoderiskconnector_property_risk_policy_set_id
  }
  property {
    name  = "sessionId"
    type  = "string"
    value = var.pingonesnoderiskconnector_property_session_id
  }
  property {
    name  = "showPoweredBy"
    type  = "string"
    value = var.pingonesnoderiskconnector_property_show_powered_by
  }
  property {
    name  = "skRiskFP"
    type  = "string"
    value = var.pingonesnoderiskconnector_property_sk_risk_fp
  }
  property {
    name  = "skipButtonPress"
    type  = "string"
    value = var.pingonesnoderiskconnector_property_skip_button_press
  }
  property {
    name  = "subtype"
    type  = "string"
    value = var.pingonesnoderiskconnector_property_subtype
  }
  property {
    name  = "targetResourceId"
    type  = "string"
    value = var.pingonesnoderiskconnector_property_target_resource_id
  }
  property {
    name  = "targetResourceName"
    type  = "string"
    value = var.pingonesnoderiskconnector_property_target_resource_name
  }
  property {
    name  = "targetedPolicy"
    type  = "boolean"
    value = var.pingonesnoderiskconnector_property_targeted_policy
  }
  property {
    name  = "userGroups"
    type  = "string"
    value = var.pingonesnoderiskconnector_property_user_groups
  }
  property {
    name  = "userId"
    type  = "string"
    value = var.pingonesnoderiskconnector_property_user_id
  }
  property {
    name  = "userName"
    type  = "string"
    value = var.pingonesnoderiskconnector_property_user_name
  }
  property {
    name  = "userType"
    type  = "string"
    value = var.pingonesnoderiskconnector_property_user_type
  }
}
