resource "pingone_davinci_connector_instance" "appleConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "appleConnector"
  }
  name = "My awesome appleConnector"
  properties = jsonencode({
    "customAuth" = jsonencode({
				"properties": {
				  "providerName": {
					"displayName": "Provider Name",
					"preferredControlType": "textField",
					"value": "${var.appleconnector_property_provider_name}"
				  },
				  "skRedirectUri": {
					"displayName": "DaVinci Redirect URL",
					"info": "Your DaVinci redirect URL. This allows an identity provider to redirect the browser back to DaVinci.",
					"preferredControlType": "textField",
					"disabled": true,
					"initializeValue": "SINGULARKEY_REDIRECT_URI",
					"copyToClip": true
				  },
				  "iss": {
					"displayName": "Issuer",
					"info": "The issuer registered claim identifies the principal that issued the client secret. Since the client secret was generated for your developer team, use your 10-character Team ID associated with your developer account.",
					"preferredControlType": "textField",
					"required": true,
					"value": "${var.appleconnector_property_issuer}"
				  },
				  "kid": {
					"displayName": "Key ID",
					"info": "A 10-character key identifier generated for the Sign in with Apple private key associated with your developer account.",
					"preferredControlType": "textField",
					"required": true,
					"value": "${var.appleconnector_property_key_id}"
				  },
				  "issuerUrl": {
					"displayName": "Issuer URL",
					"preferredControlType": "textField",
					"required": true,
					"value": "${var.appleconnector_property_issuer_url}"
				  },
				  "authorizationEndpoint": {
					"preferredControlType": "textField",
					"displayName": "Authorization Endpoint",
					"required": true,
					"value": "${var.appleconnector_property_authorization_endpoint}"
				  },
				  "tokenEndpoint": {
					"preferredControlType": "textField",
					"displayName": "Token Endpoint",
					"required": true,
					"value": "${var.appleconnector_property_token_endpoint}"
				  },
				  "clientId": {
					"displayName": "Client ID",
					"preferredControlType": "textField",
					"required": true,
					"value": "${var.appleconnector_property_client_id}"
				  },
				  "clientSecret": {
					"displayName": "Private Key",
					"info": "Content of your 'Sign in with Apple' private key associated with your developer account.",
					"preferredControlType": "textArea",
					"secure": true,
					"required": true,
					"value": "${var.appleconnector_property_private_key}"
				  },
				  "scope": {
					"displayName": "Scope",
					"preferredControlType": "textField",
					"requiredValue": "email",
					"required": true,
					"value": "${var.appleconnector_property_scope}"
				  },
				  "userConnectorAttributeMapping": {
					"type": "object",
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
						  "value1": "sub"
						},
						"name": {
						  "value1": "email"
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
						"name": "sub",
						"description": "Sub",
						"type": "string",
						"value": null,
						"minLength": "1",
						"maxLength": "300",
						"required": true,
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
				  "disableCreateUser": {
					"displayName": "Disable Shadow User Creation",
					"preferredControlType": "toggleSwitch",
					"value": false,
					"info": "A shadow user is implicitly created, unless disabled."
				  },
				  "returnToUrl": {
					"displayName": "Application Return To URL",
					"preferredControlType": "textField",
					"info": "When using the embedded flow player widget and an IDP/Social Login connector, provide a callback URL to return back to the application."
				  }
				}
			  })
  })
}
