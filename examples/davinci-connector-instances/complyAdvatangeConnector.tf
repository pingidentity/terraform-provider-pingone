resource "pingone_davinci_connector_instance" "complyAdvatangeConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "complyAdvatangeConnector"
  }
  name = "My awesome complyAdvatangeConnector"
  property {
    name  = "apiKey"
    type  = "string"
    value = var.complyadvatangeconnector_property_api_key
  }
  property {
    name  = "baseUrl"
    type  = "string"
    value = var.complyadvatangeconnector_property_base_url
  }
  property {
    name  = "countryCodes"
    type  = "string"
    value = var.complyadvatangeconnector_property_country_codes
  }
  property {
    name  = "filterEntityTypes"
    type  = "string"
    value = var.complyadvatangeconnector_property_filter_entity_types
  }
  property {
    name  = "filterTypes"
    type  = "string"
    value = var.complyadvatangeconnector_property_filter_types
  }
  property {
    name  = "fuzziness"
    type  = "string"
    value = var.complyadvatangeconnector_property_fuzziness
  }
  property {
    name  = "removeDeceased"
    type  = "string"
    value = var.complyadvatangeconnector_property_remove_deceased
  }
  property {
    name  = "searchTerm"
    type  = "string"
    value = var.complyadvatangeconnector_property_search_term
  }
  property {
    name  = "yearOfBirth"
    type  = "string"
    value = var.complyadvatangeconnector_property_year_of_birth
  }
}
