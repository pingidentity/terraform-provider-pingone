resource "pingone_davinci_connector_instance" "connectorIdiVERIFIED" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorIdiVERIFIED"
  }
  name = "My awesome connectorIdiVERIFIED"
  property {
    name  = "apiSecret"
    type  = "string"
    value = var.connectoridiverified_property_api_secret
  }
  property {
    name  = "companyKey"
    type  = "string"
    value = var.connectoridiverified_property_company_key
  }
  property {
    name  = "dppa"
    type  = "string"
    value = var.connectoridiverified_property_dppa
  }
  property {
    name  = "glba"
    type  = "string"
    value = var.connectoridiverified_property_glba
  }
  property {
    name  = "idiEnv"
    type  = "string"
    value = var.connectoridiverified_property_idi_env
  }
  property {
    name  = "reqBody"
    type  = "string"
    value = var.connectoridiverified_property_req_body
  }
  property {
    name  = "siteKey"
    type  = "string"
    value = var.connectoridiverified_property_site_key
  }
  property {
    name  = "uniqueUrl"
    type  = "string"
    value = var.connectoridiverified_property_unique_url
  }
}
