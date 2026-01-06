resource "pingone_davinci_connector_instance" "annotation_connector_example" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "annotationConnector"
  }

  name = "myAnnotationConnector"
}

resource "pingone_davinci_connector_instance" "pingfederate_connector_example" {
  environment_id = var.pingone_environment_id
  connector = {
    id = "pingFederateConnectorV2"
  }
  name = "myPingFederateConnector"

  properties = jsonencode({
    "openId" : {
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
          "value" : var.pingfederate_client_id
        },
        "clientSecret" : {
          "type" : "string",
          "displayName" : "Client Secret",
          "preferredControlType" : "textField",
          "secure" : true,
          "required" : true,
          "value" : var.pingfederate_client_secret
        },
        "scope" : {
          "type" : "string",
          "displayName" : "Scope",
          "preferredControlType" : "textField",
          "requiredValue" : "openid",
          "value" : "openid",
          "required" : true
        },
        "issuerUrl" : {
          "type" : "string",
          "displayName" : "Base URL",
          "preferredControlType" : "textField",
          "value" : "https://example.com",
          "required" : true
        },
        "returnToUrl" : {
          "displayName" : "Application Return To URL",
          "preferredControlType" : "textField",
          "info" : "When using the embedded flow player widget and an IDP/Social Login connector, provide a callback URL to return back to the application.",
          "value" : "https://example.com/callback"
        }
      }
    }
  })
}