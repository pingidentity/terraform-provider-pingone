resource "pingone_davinci_connector_instance" "akamaiApConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "akamaiApConnector"
  }
  name = "My awesome akamaiApConnector"
  property {
    name  = "akamaiUserRiskHeader"
    type  = "string"
    value = var.akamaiapconnector_property_akamai_user_risk_header
  }
  property {
    name  = "highValueThreshold"
    type  = "string"
    value = var.akamaiapconnector_property_high_value_threshold
  }
  property {
    name  = "mediumValueThreshold"
    type  = "string"
    value = var.akamaiapconnector_property_medium_value_threshold
  }
}
