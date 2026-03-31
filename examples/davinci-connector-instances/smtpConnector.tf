resource "pingone_davinci_connector_instance" "smtpConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "smtpConnector"
  }
  name = "My awesome smtpConnector"
  property {
    name  = "email"
    type  = "string"
    value = var.smtpconnector_property_email
  }
  property {
    name  = "from"
    type  = "string"
    value = var.smtpconnector_property_from
  }
  property {
    name  = "hostname"
    type  = "string"
    value = var.smtpconnector_property_hostname
  }
  property {
    name  = "htmlMessage"
    type  = "string"
    value = var.smtpconnector_property_html_message
  }
  property {
    name  = "message"
    type  = "string"
    value = var.smtpconnector_property_message
  }
  property {
    name  = "name"
    type  = "string"
    value = var.smtpconnector_property_name
  }
  property {
    name  = "password"
    type  = "string"
    value = var.smtpconnector_property_password
  }
  property {
    name  = "port"
    type  = "string"
    value = var.smtpconnector_property_port
  }
  property {
    name  = "secureFlag"
    type  = "string"
    value = var.smtpconnector_property_secure_flag
  }
  property {
    name  = "sendTestEmail"
    type  = "string"
    value = var.smtpconnector_property_send_test_email
  }
  property {
    name  = "subject"
    type  = "string"
    value = var.smtpconnector_property_subject
  }
  property {
    name  = "username"
    type  = "string"
    value = var.smtpconnector_property_username
  }
}
