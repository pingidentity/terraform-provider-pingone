resource "pingone_davinci_connector_instance" "babelStreetConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "babelStreetConnector"
  }
  name = "My awesome babelStreetConnector"
  property {
    name  = "address1City"
    type  = "string"
    value = var.babelstreetconnector_property_address1_city
  }
  property {
    name  = "address1Country"
    type  = "string"
    value = var.babelstreetconnector_property_address1_country
  }
  property {
    name  = "address1HouseNo"
    type  = "string"
    value = var.babelstreetconnector_property_address1_house_no
  }
  property {
    name  = "address1PostCode"
    type  = "string"
    value = var.babelstreetconnector_property_address1_post_code
  }
  property {
    name  = "address1Road"
    type  = "string"
    value = var.babelstreetconnector_property_address1_road
  }
  property {
    name  = "address1State"
    type  = "string"
    value = var.babelstreetconnector_property_address1_state
  }
  property {
    name  = "address1Unit"
    type  = "string"
    value = var.babelstreetconnector_property_address1_unit
  }
  property {
    name  = "address2City"
    type  = "string"
    value = var.babelstreetconnector_property_address2_city
  }
  property {
    name  = "address2Country"
    type  = "string"
    value = var.babelstreetconnector_property_address2_country
  }
  property {
    name  = "address2HouseNo"
    type  = "string"
    value = var.babelstreetconnector_property_address2_house_no
  }
  property {
    name  = "address2PostCode"
    type  = "string"
    value = var.babelstreetconnector_property_address2_post_code
  }
  property {
    name  = "address2Road"
    type  = "string"
    value = var.babelstreetconnector_property_address2_road
  }
  property {
    name  = "address2State"
    type  = "string"
    value = var.babelstreetconnector_property_address2_state
  }
  property {
    name  = "address2Unit"
    type  = "string"
    value = var.babelstreetconnector_property_address2_unit
  }
  property {
    name  = "apiKey"
    type  = "string"
    value = var.babelstreetconnector_property_api_key
  }
  property {
    name  = "content"
    type  = "string"
    value = var.babelstreetconnector_property_content
  }
  property {
    name  = "entityType"
    type  = "string"
    value = var.babelstreetconnector_property_entity_type
  }
  property {
    name  = "entityTypeName1"
    type  = "string"
    value = var.babelstreetconnector_property_entity_type_name1
  }
  property {
    name  = "entityTypeName2"
    type  = "string"
    value = var.babelstreetconnector_property_entity_type_name2
  }
  property {
    name  = "genderName1"
    type  = "string"
    value = var.babelstreetconnector_property_gender_name1
  }
  property {
    name  = "genderName2"
    type  = "string"
    value = var.babelstreetconnector_property_gender_name2
  }
  property {
    name  = "languageName1"
    type  = "string"
    value = var.babelstreetconnector_property_language_name1
  }
  property {
    name  = "languageName2"
    type  = "string"
    value = var.babelstreetconnector_property_language_name2
  }
  property {
    name  = "name"
    type  = "string"
    value = var.babelstreetconnector_property_name
  }
  property {
    name  = "sourceLanguageOfOrigin"
    type  = "string"
    value = var.babelstreetconnector_property_source_language_of_origin
  }
  property {
    name  = "targetLanguage"
    type  = "string"
    value = var.babelstreetconnector_property_target_language
  }
  property {
    name  = "textName1"
    type  = "string"
    value = var.babelstreetconnector_property_text_name1
  }
  property {
    name  = "textName2"
    type  = "string"
    value = var.babelstreetconnector_property_text_name2
  }
}
