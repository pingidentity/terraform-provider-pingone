resource "pingone_davinci_connector_instance" "melissaConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "melissaConnector"
  }
  name = "My awesome melissaConnector"
  property {
    name  = "a1"
    type  = "string"
    value = var.melissaconnector_property_a1
  }
  property {
    name  = "a2"
    type  = "string"
    value = var.melissaconnector_property_a2
  }
  property {
    name  = "a3"
    type  = "string"
    value = var.melissaconnector_property_a3
  }
  property {
    name  = "a4"
    type  = "string"
    value = var.melissaconnector_property_a4
  }
  property {
    name  = "a5"
    type  = "string"
    value = var.melissaconnector_property_a5
  }
  property {
    name  = "a6"
    type  = "string"
    value = var.melissaconnector_property_a6
  }
  property {
    name  = "a7"
    type  = "string"
    value = var.melissaconnector_property_a7
  }
  property {
    name  = "a8"
    type  = "string"
    value = var.melissaconnector_property_a8
  }
  property {
    name  = "admarea"
    type  = "string"
    value = var.melissaconnector_property_admarea
  }
  property {
    name  = "apiKey"
    type  = "string"
    value = var.melissaconnector_property_api_key
  }
  property {
    name  = "ctry"
    type  = "string"
    value = var.melissaconnector_property_ctry
  }
  property {
    name  = "ddeplo"
    type  = "string"
    value = var.melissaconnector_property_ddeplo
  }
  property {
    name  = "deplo"
    type  = "string"
    value = var.melissaconnector_property_deplo
  }
  property {
    name  = "loc"
    type  = "string"
    value = var.melissaconnector_property_loc
  }
  property {
    name  = "opt"
    type  = "string"
    value = var.melissaconnector_property_opt
  }
  property {
    name  = "org"
    type  = "string"
    value = var.melissaconnector_property_org
  }
  property {
    name  = "postal"
    type  = "string"
    value = var.melissaconnector_property_postal
  }
  property {
    name  = "subadmarea"
    type  = "string"
    value = var.melissaconnector_property_subadmarea
  }
  property {
    name  = "subnatarea"
    type  = "string"
    value = var.melissaconnector_property_subnatarea
  }
}
