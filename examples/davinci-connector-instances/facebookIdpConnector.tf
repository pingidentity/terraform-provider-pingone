resource "pingone_davinci_connector_instance" "facebookIdpConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "facebookIdpConnector"
  }
  name = "My awesome facebookIdpConnector"
  properties = jsonencode({
    "oauth2" = jsonencode({
				"properties": {
				  "providerName": {
					"type": "string",
					"displayName": "Provider Name",
					"preferredControlType": "textField",
					"value": "Login with Facebook"
				  },
				  "skRedirectUri": {
					"type": "string",
					"displayName": "DaVinci Redirect URL",
					"info": "Enter this in your identity provider configuration to allow it to redirect the browser back to DaVinci. If you use a custom PingOne domain, modify the URL accordingly.",
					"preferredControlType": "textField",
					"disabled": true,
					"initializeValue": "SINGULARKEY_REDIRECT_URI",
					"copyToClip": true
				  },
				  "clientId": {
					"type": "string",
					"displayName": "Application ID",
					"preferredControlType": "textField",
					"required": true,
					"value": "${var.facebookidpconnector_property_application_id}"
				  },
				  "clientSecret": {
					"type": "string",
					"displayName": "Client Secret",
					"preferredControlType": "textField",
					"secure": true,
					"required": true,
					"value": "${var.facebookidpconnector_property_client_secret}"
				  },
				  "scope": {
					"type": "string",
					"displayName": "Scope",
					"preferredControlType": "textField",
					"requiredValue": "email",
					"required": true,
					"value": "${var.facebookidpconnector_property_scope}"
				  },
				  "disableCreateUser": {
					"displayName": "Disable Shadow User",
					"preferredControlType": "toggleSwitch",
					"value": true,
					"info": "A shadow user is implicitly created, unless disabled."
				  },
				  "userConnectorAttributeMapping": {
					"type": "object",
					"displayName": null,
					"preferredControlType": "userConnectorAttributeMapping",
					"newMappingAllowed": true,
					"title1": null,
					"title2": null,
					"sections": [
					  "attributeMapping"
					],
					"value": {
					  "userPoolConnectionId": "defaultUserPool",
					  "mapping": {
						"username": {
						  "value1": "id"
						},
						"name": {
						  "value1": "name"
						},
						"email": {
						  "value1": "email"
						}
					  }
					}
				  },
				  "customAttributes": {
					"type": "array",
					"displayName": "Connector Attributes",
					"preferredControlType": "tableViewAttributes",
					"info": "These attributes will be available in User Connector Attribute Mapping.",
					"sections": [
					  "connectorAttributes"
					],
					"value": [
					  {
						"name": "id",
						"description": "ID",
						"type": "string",
						"value": null,
						"minLength": "1",
						"maxLength": "300",
						"required": true,
						"attributeType": "sk"
					  },
					  {
						"name": "name",
						"description": "Display Name",
						"type": "string",
						"value": null,
						"minLength": "1",
						"maxLength": "250",
						"required": false,
						"attributeType": "sk"
					  },
					  {
						"name": "email",
						"description": "Email",
						"type": "string",
						"value": null,
						"minLength": "1",
						"maxLength": "250",
						"required": false,
						"attributeType": "sk"
					  }
					]
				  },
				  "state": {
					"displayName": "Send state with request",
					"value": true,
					"preferredControlType": "toggleSwitch",
					"info": "Send unique state value with every request"
				  },
				  "returnToUrl": {
					"displayName": "Application Return To URL",
					"preferredControlType": "textField",
					"info": "When using the embedded flow player widget and an IDP/Social Login connector, provide a callback URL to return back to the application.",
					"value": "${var.facebookidpconnector_property_callback_url}"
				  }
				}
			  })
  })
}
