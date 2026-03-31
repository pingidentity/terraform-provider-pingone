resource "pingone_davinci_connector_instance" "argyleConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "argyleConnector"
  }
  name = "My awesome argyleConnector"
  property {
    name  = "accountId"
    type  = "string"
    value = var.argyleconnector_property_account_id
  }
  property {
    name  = "apiUrl"
    type  = "string"
    value = var.argyleconnector_property_api_url
  }
  property {
    name  = "backToSearchButtonTitle"
    type  = "string"
    value = var.argyleconnector_property_back_to_search_button_title
  }
  property {
    name  = "cantFindLinkItemTitle"
    type  = "string"
    value = var.argyleconnector_property_cant_find_link_item_title
  }
  property {
    name  = "clientId"
    type  = "string"
    value = var.argyleconnector_property_client_id
  }
  property {
    name  = "clientSecret"
    type  = "string"
    value = var.argyleconnector_property_client_secret
  }
  property {
    name  = "closeOnOutsideClick"
    type  = "string"
    value = var.argyleconnector_property_close_on_outside_click
  }
  property {
    name  = "companyName"
    type  = "string"
    value = var.argyleconnector_property_company_name
  }
  property {
    name  = "excludeCategories"
    type  = "string"
    value = var.argyleconnector_property_exclude_categories
  }
  property {
    name  = "excludeLinkItems"
    type  = "string"
    value = var.argyleconnector_property_exclude_link_items
  }
  property {
    name  = "exitButtonTitle"
    type  = "string"
    value = var.argyleconnector_property_exit_button_title
  }
  property {
    name  = "introSearchPlaceholder"
    type  = "string"
    value = var.argyleconnector_property_intro_search_placeholder
  }
  property {
    name  = "javascriptWebUrl"
    type  = "string"
    value = var.argyleconnector_property_javascript_web_url
  }
  property {
    name  = "linkItems"
    type  = "string"
    value = var.argyleconnector_property_link_items
  }
  property {
    name  = "loadingText"
    type  = "string"
    value = var.argyleconnector_property_loading_text
  }
  property {
    name  = "nextEvent"
    type  = "string"
    value = var.argyleconnector_property_next_event
  }
  property {
    name  = "pluginKey"
    type  = "string"
    value = var.argyleconnector_property_plugin_key
  }
  property {
    name  = "searchScreenSubtitle"
    type  = "string"
    value = var.argyleconnector_property_search_screen_subtitle
  }
  property {
    name  = "searchScreenTitle"
    type  = "string"
    value = var.argyleconnector_property_search_screen_title
  }
  property {
    name  = "showBackToSearchButton"
    type  = "string"
    value = var.argyleconnector_property_show_back_to_search_button
  }
  property {
    name  = "showCategories"
    type  = "string"
    value = var.argyleconnector_property_show_categories
  }
  property {
    name  = "showCloseButton"
    type  = "string"
    value = var.argyleconnector_property_show_close_button
  }
  property {
    name  = "userId"
    type  = "string"
    value = var.argyleconnector_property_user_id
  }
  property {
    name  = "userToken"
    type  = "string"
    value = var.argyleconnector_property_user_token
  }
}
