resource "pingone_davinci_connector_instance" "smtpConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "smtpConnector"
  }
  name = "My awesome smtpConnector"
  properties = jsonencode({
    "hostname" = var.smtpconnector_property_hostname
    "name" = var.smtpconnector_property_name
    "password" = var.smtpconnector_property_password
    "port" = var.smtpconnector_property_port
    "secureFlag" = var.smtpconnector_property_secure_flag
    "username" = var.smtpconnector_property_username
  })
}
