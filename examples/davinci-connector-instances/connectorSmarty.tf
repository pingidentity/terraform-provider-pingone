resource "pingone_davinci_connector_instance" "connectorSmarty" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorSmarty"
  }
  name = "My awesome connectorSmarty"
  property {
    name  = "address1"
    type  = "string"
    value = var.connectorsmarty_property_address1
  }
  property {
    name  = "address2"
    type  = "string"
    value = var.connectorsmarty_property_address2
  }
  property {
    name  = "administrativeArea"
    type  = "string"
    value = var.connectorsmarty_property_administrative_area
  }
  property {
    name  = "authId"
    type  = "string"
    value = var.connectorsmarty_property_auth_id
  }
  property {
    name  = "authToken"
    type  = "string"
    value = var.connectorsmarty_property_auth_token
  }
  property {
    name  = "candidates"
    type  = "string"
    value = var.connectorsmarty_property_candidates
  }
  property {
    name  = "city"
    type  = "string"
    value = var.connectorsmarty_property_city
  }
  property {
    name  = "country"
    type  = "string"
    value = var.connectorsmarty_property_country
  }
  property {
    name  = "license"
    type  = "string"
    value = var.connectorsmarty_property_license
  }
  property {
    name  = "locality"
    type  = "string"
    value = var.connectorsmarty_property_locality
  }
  property {
    name  = "match"
    type  = "string"
    value = var.connectorsmarty_property_match
  }
  property {
    name  = "postalCode"
    type  = "string"
    value = var.connectorsmarty_property_postal_code
  }
  property {
    name  = "state"
    type  = "string"
    value = var.connectorsmarty_property_state
  }
  property {
    name  = "street"
    type  = "string"
    value = var.connectorsmarty_property_street
  }
}
