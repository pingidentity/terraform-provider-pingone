resource "pingone_davinci_connector_instance" "pingFederateConnectorV2" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "pingFederateConnectorV2"
  }
  name = "My awesome pingFederateConnectorV2"
  properties = jsonencode({
    "openId" = jsonencode({
				"properties": {
				  "skRedirectUri": {
					"type": "string",
					"displayName": "Redirect URL",
					"info": "Enter this in your identity provider configuration to allow it to redirect the browser back to DaVinci. If you use a custom PingOne domain, modify the URL accordingly.",
					"preferredControlType": "textField",
					"disabled": true,
					"initializeValue": "SINGULARKEY_REDIRECT_URI",
					"copyToClip": true
				  },
				  "clientId": {
					"type": "string",
					"displayName": "Client ID",
					"placeholder": "",
					"preferredControlType": "textField",
					"required": true,
					"value": "${var.pingfederateconnectorv2_property_client_id}"
				  },
				  "clientSecret": {
					"type": "string",
					"displayName": "Client Secret",
					"preferredControlType": "textField",
					"secure": true,
					"required": true,
					"value": "${var.pingfederateconnectorv2_property_client_secret}"
				  },
				  "scope": {
					"type": "string",
					"displayName": "Scope",
					"preferredControlType": "textField",
					"requiredValue": "openid",
					"value": "${var.pingfederateconnectorv2_property_client_scope}",
					"required": true
				  },
				  "issuerUrl": {
					"type": "string",
					"displayName": "Base URL",
					"preferredControlType": "textField",
					"value": "${var.pingfederateconnectorv2_property_base_url}",
					"required": true
				  },
				  "returnToUrl": {
					"displayName": "Application Return To URL",
					"preferredControlType": "textField",
					"info": "When using the embedded flow player widget and an IDP/Social Login connector, provide a callback URL to return back to the application.",
					"value": "${var.pingfederateconnectorv2_property_application_callback}"
				  }
				}
			  })
  })
}
