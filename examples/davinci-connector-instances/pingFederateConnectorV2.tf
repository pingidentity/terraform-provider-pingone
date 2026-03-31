resource "pingone_davinci_connector_instance" "pingFederateConnectorV2" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "pingFederateConnectorV2"
  }
  name = "My awesome pingFederateConnectorV2"
  property {
    name  = "authType"
    type  = "string"
    value = var.pingfederateconnectorv2_property_auth_type
  }
  property {
    name  = "button"
    type  = "string"
    value = var.pingfederateconnectorv2_property_button
  }
  property {
    name  = "customCSS"
    type  = "string"
    value = var.pingfederateconnectorv2_property_custom_css
  }
  property {
    name  = "customHTML"
    type  = "string"
    value = var.pingfederateconnectorv2_property_custom_html
  }
  property {
    name  = "customScript"
    type  = "string"
    value = var.pingfederateconnectorv2_property_custom_script
  }
  property {
    name = "openId"
    type = "string"
    value = jsonencode({
      "properties" : {
        "skRedirectUri" : {
          "type" : "string",
          "displayName" : "Redirect URL",
          "info" : "Enter this in your identity provider configuration to allow it to redirect the browser back to DaVinci. If you use a custom PingOne domain, modify the URL accordingly.",
          "preferredControlType" : "textField",
          "disabled" : true,
          "initializeValue" : "SINGULARKEY_REDIRECT_URI",
          "copyToClip" : true
        },
        "clientId" : {
          "type" : "string",
          "displayName" : "Client ID",
          "placeholder" : "",
          "preferredControlType" : "textField",
          "required" : true,
          "value" : "${var.pingfederateconnectorv2_property_client_id}"
        },
        "clientSecret" : {
          "type" : "string",
          "displayName" : "Client Secret",
          "preferredControlType" : "textField",
          "secure" : true,
          "required" : true,
          "value" : "${var.pingfederateconnectorv2_property_client_secret}"
        },
        "scope" : {
          "type" : "string",
          "displayName" : "Scope",
          "preferredControlType" : "textField",
          "requiredValue" : "openid",
          "value" : "${var.pingfederateconnectorv2_property_client_scope}",
          "required" : true
        },
        "issuerUrl" : {
          "type" : "string",
          "displayName" : "Base URL",
          "preferredControlType" : "textField",
          "value" : "${var.pingfederateconnectorv2_property_base_url}",
          "required" : true
        },
        "returnToUrl" : {
          "displayName" : "Application Return To URL",
          "preferredControlType" : "textField",
          "info" : "When using the embedded flow player widget and an IDP/Social Login connector, provide a callback URL to return back to the application.",
          "value" : "${var.pingfederateconnectorv2_property_application_callback}"
        }
      }
    })
  }
  property {
    name  = "showPoweredBy"
    type  = "string"
    value = var.pingfederateconnectorv2_property_show_powered_by
  }
  property {
    name  = "skipButtonPress"
    type  = "string"
    value = var.pingfederateconnectorv2_property_skip_button_press
  }
  property {
    name  = "widgetLogoUrl"
    type  = "string"
    value = var.pingfederateconnectorv2_property_widget_logo_url
  }
  property {
    name  = "widgetUrl"
    type  = "string"
    value = var.pingfederateconnectorv2_property_widget_url
  }
}
