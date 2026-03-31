resource "pingone_davinci_connector_instance" "daonConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "daonConnector"
  }
  name = "My awesome daonConnector"
  property {
    name  = "apiUrl"
    type  = "string"
    value = var.daonconnector_property_api_url
  }
  property {
    name  = "authId"
    type  = "string"
    value = var.daonconnector_property_auth_id
  }
  property {
    name  = "description"
    type  = "string"
    value = var.daonconnector_property_description
  }
  property {
    name  = "password"
    type  = "string"
    value = var.daonconnector_property_password
  }
  property {
    name  = "policyUrl"
    type  = "string"
    value = var.daonconnector_property_policy_url
  }
  property {
    name  = "pushNotificationType"
    type  = "string"
    value = var.daonconnector_property_push_notification_type
  }
  property {
    name  = "secureImageTransactionContent"
    type  = "string"
    value = var.daonconnector_property_secure_image_transaction_content
  }
  property {
    name  = "secureTextTransactionContent"
    type  = "string"
    value = var.daonconnector_property_secure_text_transaction_content
  }
  property {
    name  = "secureTransactionContentType"
    type  = "string"
    value = var.daonconnector_property_secure_transaction_content_type
  }
  property {
    name  = "type"
    type  = "string"
    value = var.daonconnector_property_type
  }
  property {
    name  = "userId"
    type  = "string"
    value = var.daonconnector_property_user_id
  }
  property {
    name  = "userLogin"
    type  = "string"
    value = var.daonconnector_property_user_login
  }
  property {
    name  = "username"
    type  = "string"
    value = var.daonconnector_property_username
  }
}
