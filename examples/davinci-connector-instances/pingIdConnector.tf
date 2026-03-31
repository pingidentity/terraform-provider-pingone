resource "pingone_davinci_connector_instance" "pingIdConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "pingIdConnector"
  }
  name = "My awesome pingIdConnector"
  property {
    name  = "appIconUrl"
    type  = "string"
    value = var.pingidconnector_property_app_icon_url
  }
  property {
    name  = "appName"
    type  = "string"
    value = var.pingidconnector_property_app_name
  }
  property {
    name  = "authType"
    type  = "string"
    value = var.pingidconnector_property_auth_type
  }
  property {
    name  = "button"
    type  = "string"
    value = var.pingidconnector_property_button
  }
  property {
    name  = "claimsNameValuePairs"
    type  = "string"
    value = var.pingidconnector_property_claims_name_value_pairs
  }
  property {
    name = "customAuth"
    type = "string"
    value = jsonencode({
      "properties" : {
        "pingIdProperties" : {
          "displayName" : "PingID properties file",
          "preferredControlType" : "secureTextArea",
          "hashedVisibility" : true,
          "required" : true,
          "info" : "Paste the contents of the PingID properties file into this field.",
          "value" : "${file(var.pingidconnector_property_pingid_properties_file_path)}"
        },
        "returnToUrl" : {
          "displayName" : "Application Return To URL",
          "preferredControlType" : "textField",
          "info" : "When using the embedded flow player widget and an IDP/Social Login connector, provide a callback URL to return back to the application."
        }
      }
    })
  }
  property {
    name  = "fname"
    type  = "string"
    value = var.pingidconnector_property_fname
  }
  property {
    name  = "group"
    type  = "string"
    value = var.pingidconnector_property_group
  }
  property {
    name  = "lname"
    type  = "string"
    value = var.pingidconnector_property_lname
  }
  property {
    name  = "passwordlessContext"
    type  = "string"
    value = var.pingidconnector_property_passwordless_context
  }
  property {
    name  = "phone"
    type  = "string"
    value = var.pingidconnector_property_phone
  }
  property {
    name  = "pingidIp"
    type  = "string"
    value = var.pingidconnector_property_pingid_ip
  }
  property {
    name  = "saasid"
    type  = "string"
    value = var.pingidconnector_property_saasid
  }
  property {
    name  = "showPoweredBy"
    type  = "string"
    value = var.pingidconnector_property_show_powered_by
  }
  property {
    name  = "skipButtonPress"
    type  = "string"
    value = var.pingidconnector_property_skip_button_press
  }
  property {
    name  = "sub"
    type  = "string"
    value = var.pingidconnector_property_sub
  }
  property {
    name  = "voiceNumber"
    type  = "string"
    value = var.pingidconnector_property_voice_number
  }
}
