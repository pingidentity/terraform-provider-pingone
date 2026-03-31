resource "pingone_davinci_connector_instance" "locationPolicyConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "locationPolicyConnector"
  }
  name = "My awesome locationPolicyConnector"
  property {
    name  = "confidenceThreshold"
    type  = "string"
    value = var.locationpolicyconnector_property_confidence_threshold
  }
  property {
    name  = "geolocationPermission"
    type  = "string"
    value = var.locationpolicyconnector_property_geolocation_permission
  }
  property {
    name  = "ipAddressList"
    type  = "string"
    value = var.locationpolicyconnector_property_ip_address_list
  }
  property {
    name  = "regions"
    type  = "string"
    value = var.locationpolicyconnector_property_regions
  }
}
