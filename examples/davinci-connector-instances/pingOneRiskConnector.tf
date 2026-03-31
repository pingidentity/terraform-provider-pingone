resource "pingone_davinci_connector_instance" "pingOneRiskConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "pingOneRiskConnector"
  }
  name = "My awesome pingOneRiskConnector"
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
    name  = "completionStatus"
    type  = "string"
    value = var.pingoneriskconnector_property_completion_status
  }
  property {
    name  = "cookie"
    type  = "string"
    value = var.pingoneriskconnector_property_cookie
  }
  property {
    name  = "createdAt"
    type  = "string"
    value = var.pingoneriskconnector_property_created_at
  }
  property {
    name  = "customAttributes"
    type  = "string"
    value = var.pingoneriskconnector_property_custom_attributes
  }
  property {
    name  = "envId"
    type  = "string"
    value = var.pingone_worker_app_environment_id
  }
  property {
    name  = "externalId"
    type  = "string"
    value = var.pingoneriskconnector_property_external_id
  }
  property {
    name  = "feedbackCategory"
    type  = "string"
    value = var.pingoneriskconnector_property_feedback_category
  }
  property {
    name  = "flowType"
    type  = "string"
    value = var.pingoneriskconnector_property_flow_type
  }
  property {
    name  = "ipAddress"
    type  = "string"
    value = var.pingoneriskconnector_property_ip_address
  }
  property {
    name  = "password"
    type  = "string"
    value = var.pingoneriskconnector_property_password
  }
  property {
    name  = "passwordAlgorithm"
    type  = "string"
    value = var.pingoneriskconnector_property_password_algorithm
  }
  property {
    name  = "reason"
    type  = "string"
    value = var.pingoneriskconnector_property_reason
  }
  property {
    name  = "reasonAutomatedAttack"
    type  = "string"
    value = var.pingoneriskconnector_property_reason_automated_attack
  }
  property {
    name  = "reasonCompromisedAccount"
    type  = "string"
    value = var.pingoneriskconnector_property_reason_compromised_account
  }
  property {
    name  = "reasonFalseHighRisk"
    type  = "string"
    value = var.pingoneriskconnector_property_reason_false_high_risk
  }
  property {
    name  = "reasonFriendlyBot"
    type  = "string"
    value = var.pingoneriskconnector_property_reason_friendly_bot
  }
  property {
    name  = "reasonNewAccountFraud"
    type  = "string"
    value = var.pingoneriskconnector_property_reason_new_account_fraud
  }
  property {
    name  = "region"
    type  = "string"
    value = var.pingoneriskconnector_property_region
  }
  property {
    name  = "riskEvaluationId"
    type  = "string"
    value = var.pingoneriskconnector_property_risk_evaluation_id
  }
  property {
    name  = "riskId"
    type  = "string"
    value = var.pingoneriskconnector_property_risk_id
  }
  property {
    name  = "riskPolicySetId"
    type  = "string"
    value = var.pingoneriskconnector_property_risk_policy_set_id
  }
  property {
    name  = "sessionId"
    type  = "string"
    value = var.pingoneriskconnector_property_session_id
  }
  property {
    name  = "showPoweredBy"
    type  = "string"
    value = var.pingoneriskconnector_property_show_powered_by
  }
  property {
    name  = "skRiskFP"
    type  = "string"
    value = var.pingoneriskconnector_property_sk_risk_fp
  }
  property {
    name  = "skipButtonPress"
    type  = "string"
    value = var.pingoneriskconnector_property_skip_button_press
  }
  property {
    name  = "subtype"
    type  = "string"
    value = var.pingoneriskconnector_property_subtype
  }
  property {
    name  = "targetResourceId"
    type  = "string"
    value = var.pingoneriskconnector_property_target_resource_id
  }
  property {
    name  = "targetResourceName"
    type  = "string"
    value = var.pingoneriskconnector_property_target_resource_name
  }
  property {
    name  = "targetedPolicy"
    type  = "boolean"
    value = var.pingoneriskconnector_property_targeted_policy
  }
  property {
    name  = "userAgent"
    type  = "string"
    value = var.pingoneriskconnector_property_user_agent
  }
  property {
    name  = "userGroups"
    type  = "string"
    value = var.pingoneriskconnector_property_user_groups
  }
  property {
    name  = "userId"
    type  = "string"
    value = var.pingoneriskconnector_property_user_id
  }
  property {
    name  = "userName"
    type  = "string"
    value = var.pingoneriskconnector_property_user_name
  }
  property {
    name  = "userType"
    type  = "string"
    value = var.pingoneriskconnector_property_user_type
  }
}
