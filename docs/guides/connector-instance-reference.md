---
layout: ""
page_title: "Connector Instance Parameter Reference"
description: |-
  The guide describes the connection parameters for all connectors in the DaVinci platform, with examples of how to define within Terraform using the `pingone_davinci_connector_instance` resource.
---

# DaVinci Connection Instance Definitions

Below is a list of all available DaVinci Connections available for use in the `pingone_davinci_connector_instance` resource. 
If the `value` type of a property is not defined it must be inferred.


## 1Kosmos connector

Connector ID (`connector.id` in the resource): `connector1Kosmos`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `customAuth` (json):  Console display name: "Custom Parameters".


Example:
```terraform
resource "pingone_davinci_connector_instance" "connector1Kosmos" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connector1Kosmos"
  }
  name = "My awesome connector1Kosmos"
  properties = jsonencode({
    "customAuth" = jsonencode({})
  })
}
```


## AWS Lambda

Connector ID (`connector.id` in the resource): `connectorAWSLambda`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `accessKeyId` (string): Access key to connect to AWS Environment. Console display name: "Access Key Id".
* `region` (string): AWS Region where the Lambda function is created. Console display name: "AWS Region".
* `secretAccessKey` (string): Secret Key to access the AWS. Console display name: "AWS Secret Key".


Example:
```terraform
resource "pingone_davinci_connector_instance" "connectorAWSLambda" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorAWSLambda"
  }
  name = "My awesome connectorAWSLambda"
  properties = jsonencode({
    "accessKeyId" = var.connectorawslambda_property_access_key_id
    "region" = "eu-west-1"
    "secretAccessKey" = var.connectorawslambda_property_secret_access_key
  })
}
```


## AWS Login

Connector ID (`connector.id` in the resource): `awsIdpConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `openId` (json):  Console display name: "OpenId Parameters".


Example:
```terraform
resource "pingone_davinci_connector_instance" "awsIdpConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "awsIdpConnector"
  }
  name = "My awesome awsIdpConnector"
  properties = jsonencode({
    "openId" = jsonencode({})
  })
}
```


## AWS Secrets Manager

Connector ID (`connector.id` in the resource): `connectorAmazonAwsSecretsManager`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `accessKeyId` (string): The AWS Access Key. Console display name: "AWS Access Key".
* `region` (string): The AWS Region. Console display name: "AWS Region".
* `secretAccessKey` (string): The AWS Access Secret. Console display name: "AWS Access Secret".


Example:
```terraform
resource "pingone_davinci_connector_instance" "connectorAmazonAwsSecretsManager" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorAmazonAwsSecretsManager"
  }
  name = "My awesome connectorAmazonAwsSecretsManager"
  properties = jsonencode({
    "accessKeyId" = var.connectoramazonawssecretsmanager_property_access_key_id
    "region" = "eu-west-1"
    "secretAccessKey" = var.connectoramazonawssecretsmanager_property_secret_access_key
  })
}
```


## AbuseIPDB

Connector ID (`connector.id` in the resource): `connectorAbuseipdb`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `apiKey` (string): API Key gathered from AbuseIPDB tenant. Console display name: "API Key".


Example:
```terraform
resource "pingone_davinci_connector_instance" "connectorAbuseipdb" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorAbuseipdb"
  }
  name = "My awesome connectorAbuseipdb"
  properties = jsonencode({
    "apiKey" = var.connectorabuseipdb_property_api_key
  })
}
```


## ActiveCampaign API

Connector ID (`connector.id` in the resource): `connector-oai-activecampaignapi`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `authApiKey` (string): The authentication key to the ActiveCampaign API. Console display name: "API Key".
* `authApiVersion` (string): The version of the ActiveCampaign API. Console display name: "API Version".
* `basePath` (string): The base URL for contacting the API. Console display name: "Base Path".


Example:
```terraform
resource "pingone_davinci_connector_instance" "connector-oai-activecampaignapi" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connector-oai-activecampaignapi"
  }
  name = "My awesome connector-oai-activecampaignapi"
  properties = jsonencode({
    "authApiKey" = var.connector-oai-activecampaignapi_property_auth_api_key
    "authApiVersion" = var.connector-oai-activecampaignapi_property_auth_api_version
    "basePath" = var.connector-oai-activecampaignapi_property_base_path
  })
}
```


## Acuant

Connector ID (`connector.id` in the resource): `connectorAcuant`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `customAuth` (json):  Console display name: "Custom Parameters".


Example:
```terraform
resource "pingone_davinci_connector_instance" "connectorAcuant" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorAcuant"
  }
  name = "My awesome connectorAcuant"
  properties = jsonencode({
    "customAuth" = jsonencode({})
  })
}
```


## Adobe Marketo

Connector ID (`connector.id` in the resource): `adobemarketoConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `clientId` (string): Your Adobe Marketo client ID. Console display name: "Client ID".
* `clientSecret` (string): Your Adobe Marketo client secret. Console display name: "Client Secret".
* `endpoint` (string): The API endpoint for your Adobe Marketo instance, such as "abc123.mktorest.com/rest". Console display name: "API URL".


Example:
```terraform
resource "pingone_davinci_connector_instance" "adobemarketoConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "adobemarketoConnector"
  }
  name = "My awesome adobemarketoConnector"
  properties = jsonencode({
    "clientId" = var.adobemarketoconnector_property_client_id
    "clientSecret" = var.adobemarketoconnector_property_client_secret
    "endpoint" = var.adobemarketoconnector_property_endpoint
  })
}
```


## Akamai MFA

Connector ID (`connector.id` in the resource): `akamaiConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `customAuth` (json):  Console display name: "Custom Parameters".


Example:
```terraform
resource "pingone_davinci_connector_instance" "akamaiConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "akamaiConnector"
  }
  name = "My awesome akamaiConnector"
  properties = jsonencode({
    "customAuth" = var.akamaiconnector_property_custom_auth
  })
}
```


## Allthenticate

Connector ID (`connector.id` in the resource): `connectorAllthenticate`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `customAuth` (json):  Console display name: "Custom Parameters".


Example:
```terraform
resource "pingone_davinci_connector_instance" "connectorAllthenticate" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorAllthenticate"
  }
  name = "My awesome connectorAllthenticate"
  properties = jsonencode({
    "customAuth" = jsonencode({})
  })
}
```


## Amazon DynamoDB

Connector ID (`connector.id` in the resource): `connectorAmazonDynamoDB`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `awsAccessKey` (string): Your AWS Access Key. Console display name: "AWS Access Key".
* `awsAccessSecret` (string): Access Secret corresponding with Access Key found in Your Security Credentials. Console display name: "AWS Access Secret".
* `awsRegion` (string): The AWS Region you are using the connector for. Console display name: "AWS Region".


Example:
```terraform
resource "pingone_davinci_connector_instance" "connectorAmazonDynamoDB" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorAmazonDynamoDB"
  }
  name = "My awesome connectorAmazonDynamoDB"
  properties = jsonencode({
    "awsAccessKey" = var.connectoramazondynamodb_property_aws_access_key
    "awsAccessSecret" = var.connectoramazondynamodb_property_aws_access_secret
    "awsRegion" = "eu-west-1"
  })
}
```


## Amazon Simple Email Service

Connector ID (`connector.id` in the resource): `amazonSimpleEmailConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `awsAccessKey` (string):  Console display name: "AWS Access Key".
* `awsAccessSecret` (string):  Console display name: "AWS Access Secret".
* `awsRegion` (string):  Console display name: "AWS Region".
* `from` (string): The email address that the message appears to originate from, as registered with your AWS account, such as "support@mycompany.com". Console display name: "From (Default) *".


Example:
```terraform
resource "pingone_davinci_connector_instance" "amazonSimpleEmailConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "amazonSimpleEmailConnector"
  }
  name = "My awesome amazonSimpleEmailConnector"
  properties = jsonencode({
    "awsAccessKey" = var.amazonsimpleemailconnector_property_aws_access_key
    "awsAccessSecret" = var.amazonsimpleemailconnector_property_aws_access_secret
    "awsRegion" = "eu-west-1"
    "from" = "support@bxretail.org"
  })
}
```


## Annotation

Connector ID (`connector.id` in the resource): `annotationConnector`

*No properties*


Example:
```terraform
resource "pingone_davinci_connector_instance" "annotationConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "annotationConnector"
  }
  name = "My awesome annotationConnector"
}
```


## Apple Login

Connector ID (`connector.id` in the resource): `appleConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `customAuth` (json):  Console display name: "Custom Parameters".


Example:
```terraform
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
```


## Argyle

Connector ID (`connector.id` in the resource): `argyleConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `apiUrl` (string):  Console display name: "API Server URL".
* `clientId` (string):  Console display name: "Client ID".
* `clientSecret` (string):  Console display name: "Client Secret".
* `javascriptWebUrl` (string): Argyle loader javascript web URL. Console display name: "Argyle Loader Javascript Web URL".
* `pluginKey` (string):  Console display name: "Plugin Key".


Example:
```terraform
resource "pingone_davinci_connector_instance" "argyleConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "argyleConnector"
  }
  name = "My awesome argyleConnector"
  properties = jsonencode({
    "apiUrl" = var.argyleconnector_property_api_url
    "clientId" = var.argyleconnector_property_client_id
    "clientSecret" = var.argyleconnector_property_client_secret
    "javascriptWebUrl" = var.argyleconnector_property_javascript_web_url
    "pluginKey" = var.argyleconnector_property_plugin_key
  })
}
```


## Asignio

Connector ID (`connector.id` in the resource): `connectorAsignio`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `customAuth` (json):  Console display name: "Custom Parameters".


Example:
```terraform
resource "pingone_davinci_connector_instance" "connectorAsignio" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorAsignio"
  }
  name = "My awesome connectorAsignio"
  properties = jsonencode({
    "customAuth" = jsonencode({})
  })
}
```


## AuthID

Connector ID (`connector.id` in the resource): `connectorAuthid`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `customAuth` (json):  Console display name: "Custom Parameters".


Example:
```terraform
resource "pingone_davinci_connector_instance" "connectorAuthid" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorAuthid"
  }
  name = "My awesome connectorAuthid"
  properties = jsonencode({
    "customAuth" = jsonencode({})
  })
}
```


## AuthenticID

Connector ID (`connector.id` in the resource): `authenticIdConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `accountAccessKey` (string): Your Account Access Key provided by AuthenticID . Console display name: "Account Access Key".
* `androidSDKLicenseKey` (string): License key is whitelisted for specific package name. Console display name: "Android SDK Licence Key".
* `apiUrl` (string): AuthenticID REST API URL for sandbox/production environments. Console display name: "REST API URL".
* `baseUrl` (string): AuthenticID API URL for sandbox/production environments. Console display name: "Base URL".
* `clientCertificate` (string): Your Client Certificate provided by AuthenticID. Console display name: "Client Certificate".
* `clientKey` (string): Your Client Key provided by AuthenticID. Console display name: "Client Key".
* `iOSSDKLicenseKey` (string): License key is whitelisted for specific bundle id. Console display name: "iOS SDK Licence Key".
* `passphrase` (string): Your Certificate Passphrase provided by AuthenticID. Console display name: "Certificate Passphrase".
* `secretToken` (string): Your Secret Token provided by AuthenticID. Console display name: "Secret Token".


Example:
```terraform
resource "pingone_davinci_connector_instance" "authenticIdConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "authenticIdConnector"
  }
  name = "My awesome authenticIdConnector"
  properties = jsonencode({
    "accountAccessKey" = var.authenticidconnector_property_account_access_key
    "androidSDKLicenseKey" = var.authenticidconnector_property_android_sdk_license_key
    "apiUrl" = var.authenticidconnector_property_api_url
    "baseUrl" = var.authenticidconnector_property_base_url
    "clientCertificate" = var.authenticidconnector_property_client_certificate
    "clientKey" = var.authenticidconnector_property_client_key
    "iOSSDKLicenseKey" = var.authenticidconnector_property_ios_sdk_license_key
    "passphrase" = var.authenticidconnector_property_passphrase
    "secretToken" = var.authenticidconnector_property_secret_token
  })
}
```


## Authomize API

Connector ID (`connector.id` in the resource): `connector-oai-authomizeapireference`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `authApiKey` (string): Your Authomize API key. Console display name: "API Key".
* `basePath` (string): The base URL for the Authomize API, such as "https://api.authomize.com". Console display name: "Base URL".


Example:
```terraform
resource "pingone_davinci_connector_instance" "connector-oai-authomizeapireference" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connector-oai-authomizeapireference"
  }
  name = "My awesome connector-oai-authomizeapireference"
  properties = jsonencode({
    "authApiKey" = var.connector-oai-authomizeapireference_property_auth_api_key
    "basePath" = var.connector-oai-authomizeapireference_property_base_path
  })
}
```


## Authomize Incident Connector

Connector ID (`connector.id` in the resource): `connectorAuthomize`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `apiKey` (string): The API Key from the Authomize API Tokens creation page. Console display name: "Authomize API Key".


Example:
```terraform
resource "pingone_davinci_connector_instance" "connectorAuthomize" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorAuthomize"
  }
  name = "My awesome connectorAuthomize"
  properties = jsonencode({
    "apiKey" = var.connectorauthomize_property_api_key
  })
}
```


## Azure AD User Management

Connector ID (`connector.id` in the resource): `azureUserManagementConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `baseUrl` (string): The Microsoft API URL to target. For a custom value, select Use Custom API URL and enter a value in the Custom API URL field. Console display name: "API URL".
* `customApiUrl` (string): The URL for the Microsoft Graph API, such as "https://graph.microsoft.com/v1.0". Console display name: "Custom API URL".
* `customAuth` (json):  Console display name: "Custom Parameters".


Example:
```terraform
resource "pingone_davinci_connector_instance" "azureUserManagementConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "azureUserManagementConnector"
  }
  name = "My awesome azureUserManagementConnector"
  properties = jsonencode({
    "baseUrl" = var.azureusermanagementconnector_property_base_url
    "customApiUrl" = var.azureusermanagementconnector_property_custom_api_url
    "customAuth" = jsonencode({})
  })
}
```


## Badge

Connector ID (`connector.id` in the resource): `connectorBadge`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `customAuth` (json):  Console display name: "Custom Parameters".


Example:
```terraform
resource "pingone_davinci_connector_instance" "connectorBadge" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorBadge"
  }
  name = "My awesome connectorBadge"
  properties = jsonencode({
    "customAuth" = jsonencode({})
  })
}
```


## BambooHR

Connector ID (`connector.id` in the resource): `bambooConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `apiKey` (string):  Console display name: "API Key".
* `baseUrl` (string):  BambooHR Base URL. Console display name: "Base URL".
* `companySubDomain` (string):  Your BambooHR subdomain. Console display name: "Company Sub Domain".
* `flowId` (string): Select ID of the flow to execute when BambooHR sends a webhook. Console display name: "Flow ID".
* `skWebhookUri` (string): Use this url as the Webhook URL in the Third Party Integration's configuration. Console display name: "DaVinci Webhook URL".
* `webhookToken` (string): Create a webhook token and configure it in the bambooHR webhook url. Console display name: "Webhook Token".


Example:
```terraform
resource "pingone_davinci_connector_instance" "bambooConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "bambooConnector"
  }
  name = "My awesome bambooConnector"
  properties = jsonencode({
    "apiKey" = var.bambooconnector_property_api_key
    "baseUrl" = var.bambooconnector_property_base_url
    "companySubDomain" = var.bambooconnector_property_company_sub_domain
    "flowId" = var.bambooconnector_property_flow_id
    "skWebhookUri" = var.bambooconnector_property_sk_webhook_uri
    "webhookToken" = var.bambooconnector_property_webhook_token
  })
}
```


## Berbix

Connector ID (`connector.id` in the resource): `connectorBerbix`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `domainName` (string): Provide Berbix domain name. Console display name: "Domain Name".
* `path` (string): Provide path of the API. Console display name: "Path".
* `username` (string): Provide your Berbix user name. Console display name: "User Name".


Example:
```terraform
resource "pingone_davinci_connector_instance" "connectorBerbix" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorBerbix"
  }
  name = "My awesome connectorBerbix"
  properties = jsonencode({
    "domainName" = var.connectorberbix_property_domain_name
    "path" = var.connectorberbix_property_path
    "username" = var.connectorberbix_property_username
  })
}
```


## Beyond Identity

Connector ID (`connector.id` in the resource): `connectorBeyondIdentity`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `openId` (json):  Console display name: "OpenId Parameters".


Example:
```terraform
resource "pingone_davinci_connector_instance" "connectorBeyondIdentity" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorBeyondIdentity"
  }
  name = "My awesome connectorBeyondIdentity"
  properties = jsonencode({
    "openId" = jsonencode({})
  })
}
```


## BeyondTrust - Password Safe

Connector ID (`connector.id` in the resource): `connectorBTps`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `apiKey` (string): API Key from your Password Safe environment. Console display name: "API Key".
* `apiUser` (string): API User from your Password Safe environment. Console display name: "API User".
* `domain` (string): Domain of your Password Safe environment. Console display name: "PasswordSafe Hostname".


Example:
```terraform
resource "pingone_davinci_connector_instance" "connectorBTps" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorBTps"
  }
  name = "My awesome connectorBTps"
  properties = jsonencode({
    "apiKey" = var.connectorbtps_property_api_key
    "apiUser" = var.connectorbtps_property_api_user
    "domain" = var.connectorbtps_property_domain
  })
}
```


## BeyondTrust - Privileged Remote Access

Connector ID (`connector.id` in the resource): `connectorBTpra`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `clientID` (string): PRA API Client ID. Console display name: "Client ID".
* `clientSecret` (string): PRA API Client Secret. Console display name: "Client Secret".
* `praAPIurl` (string): URL of PRA Appliance. Console display name: "PRA Web API Address".


Example:
```terraform
resource "pingone_davinci_connector_instance" "connectorBTpra" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorBTpra"
  }
  name = "My awesome connectorBTpra"
  properties = jsonencode({
    "clientID" = var.connectorbtpra_property_client_i_d
    "clientSecret" = var.connectorbtpra_property_client_secret
    "praAPIurl" = var.pra_api_url
  })
}
```


## BeyondTrust - Remote Support

Connector ID (`connector.id` in the resource): `connectorBTrs`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `clientID` (string): RS API Client ID. Console display name: "Client ID".
* `clientSecret` (string): RS API Client Secret. Console display name: "Client Secret".
* `rsAPIurl` (string): URL of RS Appliance. Console display name: "RS Web API Address".


Example:
```terraform
resource "pingone_davinci_connector_instance" "connectorBTrs" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorBTrs"
  }
  name = "My awesome connectorBTrs"
  properties = jsonencode({
    "clientID" = var.connectorbtrs_property_client_i_d
    "clientSecret" = var.connectorbtrs_property_client_secret
    "rsAPIurl" = var.rs_api_url
  })
}
```


## BioCatch

Connector ID (`connector.id` in the resource): `biocatchConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `apiUrl` (string):  Console display name: "API Server URL".
* `customerId` (string):  Console display name: "Customer ID".
* `javascriptCdnUrl` (string):  Console display name: "Javascript CDN URL".
* `sdkToken` (string):  Console display name: "SDK Token".
* `truthApiKey` (string): Fraudulent/Genuine Session Reporting API Key. Console display name: "Truth-mapping API Key".
* `truthApiUrl` (string): Fraudulent/Genuine Session Reporting. Console display name: "Truth-mapping API URL".


Example:
```terraform
resource "pingone_davinci_connector_instance" "biocatchConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "biocatchConnector"
  }
  name = "My awesome biocatchConnector"
  properties = jsonencode({
    "apiUrl" = var.biocatchconnector_property_api_url
    "customerId" = var.biocatchconnector_property_customer_id
    "javascriptCdnUrl" = var.biocatchconnector_property_javascript_cdn_url
    "sdkToken" = var.biocatchconnector_property_sdk_token
    "truthApiKey" = var.biocatchconnector_property_truth_api_key
    "truthApiUrl" = var.biocatchconnector_property_truth_api_url
  })
}
```


## Bitbucket Login

Connector ID (`connector.id` in the resource): `bitbucketIdpConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `oauth2` (json):  Console display name: "Oauth2 Parameters".


Example:
```terraform
resource "pingone_davinci_connector_instance" "bitbucketIdpConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "bitbucketIdpConnector"
  }
  name = "My awesome bitbucketIdpConnector"
  properties = jsonencode({
    "oauth2" = jsonencode({})
  })
}
```


## CASTLE

Connector ID (`connector.id` in the resource): `castleConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `apiSecret` (string): Your 32-character Castle API secret, such as “Olc…QBF”. Console display name: "API Secret".


Example:
```terraform
resource "pingone_davinci_connector_instance" "castleConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "castleConnector"
  }
  name = "My awesome castleConnector"
  properties = jsonencode({
    "apiSecret" = var.castleconnector_property_api_secret
  })
}
```


## CLEAR

Connector ID (`connector.id` in the resource): `connectorClear`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `customAuth` (json):  Console display name: "Custom Parameters".


Example:
```terraform
resource "pingone_davinci_connector_instance" "connectorClear" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorClear"
  }
  name = "My awesome connectorClear"
  properties = jsonencode({
    "customAuth" = var.connectorclear_property_custom_auth
  })
}
```


## Challenge

Connector ID (`connector.id` in the resource): `challengeConnector`

*No properties*


Example:
```terraform
resource "pingone_davinci_connector_instance" "challengeConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "challengeConnector"
  }
  name = "My awesome challengeConnector"
}
```


## Circle Access

Connector ID (`connector.id` in the resource): `connectorCircleAccess`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `appKey` (string): App Key. Console display name: "App Key".
* `customAuth` (json):  Console display name: "Custom Parameters".
* `loginUrl` (string): The URL of your Circle Access login. Console display name: "Login Url".
* `readKey` (string): Read Key. Console display name: "Read Key".
* `returnToUrl` (string): When using the embedded flow player widget and an IDP/Social Login connector, provide a callback URL to return back to the application. Console display name: "Application Return To URL".
* `writeKey` (string): Write key. Console display name: "Write Key".


Example:
```terraform
resource "pingone_davinci_connector_instance" "connectorCircleAccess" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorCircleAccess"
  }
  name = "My awesome connectorCircleAccess"
  properties = jsonencode({
    "appKey" = var.connectorcircleaccess_property_app_key
    "customAuth" = jsonencode({})
    "loginUrl" = var.connectorcircleaccess_property_login_url
    "readKey" = var.connectorcircleaccess_property_read_key
    "returnToUrl" = var.connectorcircleaccess_property_return_to_url
    "writeKey" = var.connectorcircleaccess_property_write_key
  })
}
```


## Clearbit

Connector ID (`connector.id` in the resource): `connectorClearbit`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `apiKey` (string): Clearbit API Key. Console display name: "API Key".
* `riskApiVersion` (string): Clearbit - Risk API Version. Console display name: "Risk API Version".
* `version` (string): Clearbit - Person API Version. Console display name: "Person API Version".


Example:
```terraform
resource "pingone_davinci_connector_instance" "connectorClearbit" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorClearbit"
  }
  name = "My awesome connectorClearbit"
  properties = jsonencode({
    "apiKey" = var.connectorclearbit_property_api_key
    "riskApiVersion" = var.connectorclearbit_property_risk_api_version
    "version" = var.connectorclearbit_property_version
  })
}
```


## Cloudflare

Connector ID (`connector.id` in the resource): `connectorCloudflare`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `accountId` (string): Cloudflare Account ID. Console display name: "Account ID".
* `apiToken` (string): Cloudflare API Token. Console display name: "API Token".


Example:
```terraform
resource "pingone_davinci_connector_instance" "connectorCloudflare" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorCloudflare"
  }
  name = "My awesome connectorCloudflare"
  properties = jsonencode({
    "accountId" = var.connectorcloudflare_property_account_id
    "apiToken" = var.connectorcloudflare_property_api_token
  })
}
```


## Code Snippet

Connector ID (`connector.id` in the resource): `codeSnippetConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `code` (string): Follow example for code. Caution: Custom code is for advanced users only. Before using custom code, review the security risks in the DaVinci documentation by searching for "Using custom code safely". Console display name: "Code Snippet".
* `inputSchema` (string): Follow example for JSON schema. Console display name: "Input Schema".
* `outputSchema` (string): Follow example for JSON schema. Console display name: "Output Schema".


Example:
```terraform
resource "pingone_davinci_connector_instance" "codeSnippetConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "codeSnippetConnector"
  }
  name = "My awesome codeSnippetConnector"
  properties = jsonencode({
    "code" = var.codesnippetconnector_property_code
    "inputSchema" = var.codesnippetconnector_property_input_schema
    "outputSchema" = var.codesnippetconnector_property_output_schema
  })
}
```


## Comply Advantage

Connector ID (`connector.id` in the resource): `complyAdvatangeConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `apiKey` (string): API Key is the API key that you can retrieve from Comply Advantage Admin Portal. Console display name: "API Key".
* `baseUrl` (string): Comply Advantage API URL for sandbox/production environments. Console display name: "Base URL".


Example:
```terraform
resource "pingone_davinci_connector_instance" "complyAdvatangeConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "complyAdvatangeConnector"
  }
  name = "My awesome complyAdvatangeConnector"
  properties = jsonencode({
    "apiKey" = var.complyadvatangeconnector_property_api_key
    "baseUrl" = var.complyadvatangeconnector_property_base_url
  })
}
```


## ConnectID

Connector ID (`connector.id` in the resource): `connectIdConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `customAuth` (json):  Console display name: "Custom Parameters".


Example:
```terraform
resource "pingone_davinci_connector_instance" "connectIdConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectIdConnector"
  }
  name = "My awesome connectIdConnector"
  properties = jsonencode({
    "customAuth" = jsonencode({})
  })
}
```


## Cookie

Connector ID (`connector.id` in the resource): `cookieConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `hmacSigningKey` (string): Base64 encoded 256 bit key. Console display name: "HMAC Signing Key".


Example:
```terraform
resource "pingone_davinci_connector_instance" "cookieConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "cookieConnector"
  }
  name = "My awesome cookieConnector"
  properties = jsonencode({
    "hmacSigningKey" = var.cookieconnector_property_hmac_signing_key
  })
}
```


## Copper API

Connector ID (`connector.id` in the resource): `connector-oai-copperdeveloperapi`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `basePath` (string): The base URL for contacting the API. Console display name: "Base Path".
* `contentType` (string): Content type. Console display name: "Content-Type".
* `xPWAccessToken` (string): API Key. Console display name: "X-PW-AccessToken".
* `xPWApplication` (string): Application. Console display name: "X-PW-Application".
* `xPWUserEmail` (string): Email address of token owner. Console display name: "X-PW-UserEmail".


Example:
```terraform
resource "pingone_davinci_connector_instance" "connector-oai-copperdeveloperapi" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connector-oai-copperdeveloperapi"
  }
  name = "My awesome connector-oai-copperdeveloperapi"
  properties = jsonencode({
    "basePath" = var.connector-oai-copperdeveloperapi_property_base_path
    "contentType" = var.connector-oai-copperdeveloperapi_property_content_type
    "xPWAccessToken" = var.connector-oai-copperdeveloperapi_property_x_p_w_access_token
    "xPWApplication" = var.connector-oai-copperdeveloperapi_property_x_p_w_application
    "xPWUserEmail" = var.connector-oai-copperdeveloperapi_property_x_p_w_user_email
  })
}
```


## Credova

Connector ID (`connector.id` in the resource): `credovaConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `baseUrl` (string): Base URL for Credova API. Console display name: "Base URL".
* `password` (string): Password for the Credova Developer Portal. Console display name: "Credova Password".
* `username` (string): Username for the Credova Developer Portal. Console display name: "Credova Username".


Example:
```terraform
resource "pingone_davinci_connector_instance" "credovaConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "credovaConnector"
  }
  name = "My awesome credovaConnector"
  properties = jsonencode({
    "baseUrl" = var.credovaconnector_property_base_url
    "password" = var.credovaconnector_property_password
    "username" = var.credovaconnector_property_username
  })
}
```


## CrowdStrike

Connector ID (`connector.id` in the resource): `crowdStrikeConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `baseURL` (string): The base URL of the CrowdStrike environment. Console display name: "CrowdStrike Base URL".
* `clientId` (string): The Client ID of the application in CrowdStrike. Console display name: "Client ID".
* `clientSecret` (string): The Client Secret provided by CrowdStrike. Console display name: "Client Secret".


Example:
```terraform
resource "pingone_davinci_connector_instance" "crowdStrikeConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "crowdStrikeConnector"
  }
  name = "My awesome crowdStrikeConnector"
  properties = jsonencode({
    "baseURL" = var.base_url
    "clientId" = var.crowdstrikeconnector_property_client_id
    "clientSecret" = var.crowdstrikeconnector_property_client_secret
  })
}
```


## Daon IDV

Connector ID (`connector.id` in the resource): `connectorDaonidv`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `openId` (json):  Console display name: "OpenId Parameters".


Example:
```terraform
resource "pingone_davinci_connector_instance" "connectorDaonidv" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorDaonidv"
  }
  name = "My awesome connectorDaonidv"
  properties = jsonencode({
    "openId" = jsonencode({})
  })
}
```


## Daon IdentityX

Connector ID (`connector.id` in the resource): `daonConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `apiUrl` (string): The protocol, host and base path to the IdX API. E.g. https://api.identityx-cloud.com/tenant1/IdentityXServices/rest/v1. Console display name: "API Base URL".
* `password` (string): The password of the user to authenticate API calls. Console display name: "Admin Password".
* `username` (string): The userId to authenticate API calls. Console display name: "Admin Username".


Example:
```terraform
resource "pingone_davinci_connector_instance" "daonConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "daonConnector"
  }
  name = "My awesome daonConnector"
  properties = jsonencode({
    "apiUrl" = var.daonconnector_property_api_url
    "password" = var.daonconnector_property_password
    "username" = var.daonconnector_property_username
  })
}
```


## Data Zoo

Connector ID (`connector.id` in the resource): `dataZooConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `password` (string):  Console display name: "Data Zoo Password".
* `username` (string):  Console display name: "Data Zoo Username".


Example:
```terraform
resource "pingone_davinci_connector_instance" "dataZooConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "dataZooConnector"
  }
  name = "My awesome dataZooConnector"
  properties = jsonencode({
    "password" = var.datazooconnector_property_password
    "username" = var.datazooconnector_property_username
  })
}
```


## Datadog API

Connector ID (`connector.id` in the resource): `connector-oai-datadogapi`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `authApiKey` (string): The API key for an account that has access to the Datadog API. Console display name: "Authentication API Key".
* `authApplicationKey` (string): The Application key for an account that has access to the Datadog API. Console display name: "Authentication Application Key".
* `basePath` (string): The base URL for contacting the Datadog API, such as "https://api.us3.datadoghq.com". Console display name: "API URL".


Example:
```terraform
resource "pingone_davinci_connector_instance" "connector-oai-datadogapi" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connector-oai-datadogapi"
  }
  name = "My awesome connector-oai-datadogapi"
  properties = jsonencode({
    "authApiKey" = var.connector-oai-datadogapi_property_auth_api_key
    "authApplicationKey" = var.connector-oai-datadogapi_property_auth_application_key
    "basePath" = var.connector-oai-datadogapi_property_base_path
  })
}
```


## DeBounce

Connector ID (`connector.id` in the resource): `connectorDeBounce`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `apiKey` (string): A DeBounce API Key is physically a token/code of 13 random alphanumeric characters. If you need to create an API key, please log in to your DeBounce account and then navigate to the API section. Console display name: "API Key".


Example:
```terraform
resource "pingone_davinci_connector_instance" "connectorDeBounce" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorDeBounce"
  }
  name = "My awesome connectorDeBounce"
  properties = jsonencode({
    "apiKey" = var.connectordebounce_property_api_key
  })
}
```


## Device Policy

Connector ID (`connector.id` in the resource): `devicePolicyConnector`

*No properties*


Example:
```terraform
resource "pingone_davinci_connector_instance" "devicePolicyConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "devicePolicyConnector"
  }
  name = "My awesome devicePolicyConnector"
}
```


## DigiLocker

Connector ID (`connector.id` in the resource): `digilockerConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `oauth2` (json):  Console display name: "Oauth2 Parameters".


Example:
```terraform
resource "pingone_davinci_connector_instance" "digilockerConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "digilockerConnector"
  }
  name = "My awesome digilockerConnector"
  properties = jsonencode({
    "oauth2" = jsonencode({})
  })
}
```


## Digidentity

Connector ID (`connector.id` in the resource): `digidentityConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `oauth2` (json):  Console display name: "Oauth2 Parameters".


Example:
```terraform
resource "pingone_davinci_connector_instance" "digidentityConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "digidentityConnector"
  }
  name = "My awesome digidentityConnector"
  properties = jsonencode({
    "oauth2" = jsonencode({})
  })
}
```


## Druva inSync Cloud API

Connector ID (`connector.id` in the resource): `connector-oai-druvainsynccloud`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `authClientId` (string): The Client ID of the authenticating application. Console display name: "Client ID".
* `authClientSecret` (string): The Secret Key for the authenticating application. Console display name: "Secret Key".
* `authTokenUrl` (string): The URL used to obtain an access token. Console display name: "Token URL".
* `basePath` (string): The base URL for contacting the API. Console display name: "Base Path".


Example:
```terraform
resource "pingone_davinci_connector_instance" "connector-oai-druvainsynccloud" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connector-oai-druvainsynccloud"
  }
  name = "My awesome connector-oai-druvainsynccloud"
  properties = jsonencode({
    "authClientId" = var.connector-oai-druvainsynccloud_property_auth_client_id
    "authClientSecret" = var.connector-oai-druvainsynccloud_property_auth_client_secret
    "authTokenUrl" = var.connector-oai-druvainsynccloud_property_auth_token_url
    "basePath" = var.connector-oai-druvainsynccloud_property_base_path
  })
}
```


## Duo

Connector ID (`connector.id` in the resource): `duoConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `customAuth` (json):  Console display name: "Custom Parameters".


Example:
```terraform
resource "pingone_davinci_connector_instance" "duoConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "duoConnector"
  }
  name = "My awesome duoConnector"
  properties = jsonencode({
    "customAuth" = jsonencode({})
  })
}
```


## Entrust

Connector ID (`connector.id` in the resource): `entrustConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `applicationId` (string): The application ID for the Identity as a Service application. Console display name: "Application ID".
* `serviceDomain` (string): The domain of the Entrust service. Format is '<customer>.<region>.trustedauth.com'. For example, 'mycompany.us.trustedauth.com'. Console display name: "Service Domain".


Example:
```terraform
resource "pingone_davinci_connector_instance" "entrustConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "entrustConnector"
  }
  name = "My awesome entrustConnector"
  properties = jsonencode({
    "applicationId" = var.entrustconnector_property_application_id
    "serviceDomain" = var.entrustconnector_property_service_domain
  })
}
```


## Equifax

Connector ID (`connector.id` in the resource): `equifaxConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `baseUrl` (string): Base URL for Equifax API. Console display name: "Base URL".
* `clientId` (string): When you Create a New App, Equifax will assign a Client ID per environment for the API Product. Console display name: "Client ID".
* `clientSecret` (string): When you Create a New App, Equifax will assign a Client Secret per environment for the API Product. Console display name: "Client Secret".
* `equifaxSoapApiEnvironment` (string): SOAP API WSDL Environment. Console display name: "SOAP API Environment".
* `memberNumber` (string): Unique Identifier of Customer. Please contact Equifax Sales Representative during client onboarding for this value. Console display name: "Member Number".
* `password` (string): Password provided by Equifax for SOAP API. Console display name: "Password for SOAP API".
* `username` (string): Username provided by Equifax for SOAP API. Console display name: "Username for SOAP API".


Example:
```terraform
resource "pingone_davinci_connector_instance" "equifaxConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "equifaxConnector"
  }
  name = "My awesome equifaxConnector"
  properties = jsonencode({
    "baseUrl" = var.equifaxconnector_property_base_url
    "clientId" = var.equifaxconnector_property_client_id
    "clientSecret" = var.equifaxconnector_property_client_secret
    "equifaxSoapApiEnvironment" = var.equifaxconnector_property_equifax_soap_api_environment
    "memberNumber" = var.equifaxconnector_property_member_number
    "password" = var.equifaxconnector_property_password
    "username" = var.equifaxconnector_property_username
  })
}
```


## Error Message

Connector ID (`connector.id` in the resource): `errorConnector`

*No properties*


Example:
```terraform
resource "pingone_davinci_connector_instance" "errorConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "errorConnector"
  }
  name = "My awesome errorConnector"
}
```


## Facebook Login

Connector ID (`connector.id` in the resource): `facebookIdpConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `oauth2` (json):  Console display name: "Oauth2 Parameters".


Example:
```terraform
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
```


## Fingerprint JS

Connector ID (`connector.id` in the resource): `fingerprintjsConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `apiToken` (string):  Console display name: "Fingerprint Subscription API Token".
* `javascriptCdnUrl` (string):  Console display name: "Javascript CDN URL".
* `token` (string):  Console display name: "Fingerprint Subscription Browser Token".


Example:
```terraform
resource "pingone_davinci_connector_instance" "fingerprintjsConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "fingerprintjsConnector"
  }
  name = "My awesome fingerprintjsConnector"
  properties = jsonencode({
    "apiToken" = var.fingerprintjsconnector_property_api_token
    "javascriptCdnUrl" = var.fingerprintjsconnector_property_javascript_cdn_url
    "token" = var.fingerprintjsconnector_property_token
  })
}
```


## Finicity

Connector ID (`connector.id` in the resource): `finicityConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `appKey` (string): Finicity App Key from Finicity Developer Portal. Console display name: "Finicity App Key".
* `baseUrl` (string): Base URL for Finicity API. Console display name: "Base URL".
* `partnerId` (string): The partner id you can obtain from your Finicity developer dashboard. Console display name: "Partner ID".
* `partnerSecret` (string): Partner Secret from Finicity Developer Portal. Console display name: "Partner Secret".


Example:
```terraform
resource "pingone_davinci_connector_instance" "finicityConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "finicityConnector"
  }
  name = "My awesome finicityConnector"
  properties = jsonencode({
    "appKey" = var.finicityconnector_property_app_key
    "baseUrl" = var.finicityconnector_property_base_url
    "partnerId" = var.finicityconnector_property_partner_id
    "partnerSecret" = var.finicityconnector_property_partner_secret
  })
}
```


## Flow Analytics

Connector ID (`connector.id` in the resource): `analyticsConnector`

*No properties*


Example:
```terraform
resource "pingone_davinci_connector_instance" "analyticsConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "analyticsConnector"
  }
  name = "My awesome analyticsConnector"
}
```


## Flow Conductor

Connector ID (`connector.id` in the resource): `flowConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `enforcedSignedToken` (boolean):  Console display name: "Enforce Signed Token".
* `inputSchema` (string): Follow example for JSON schema. Console display name: "Input Schema".
* `pemPublicKey` (string): pem public key. Console display name: "Public Key".


Example:
```terraform
resource "pingone_davinci_connector_instance" "flowConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "flowConnector"
  }
  name = "My awesome flowConnector"
  properties = jsonencode({
    "enforcedSignedToken" = var.flowconnector_property_enforced_signed_token
    "inputSchema" = var.flowconnector_property_input_schema
    "pemPublicKey" = var.flowconnector_property_pem_public_key
  })
}
```


## Forter

Connector ID (`connector.id` in the resource): `forterConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `apiVersion` (string): API Version. Console display name: " Forter API Version".
* `secretKey` (string): Secret Key from Forter tenant. Console display name: "Forter Secret Key".
* `siteId` (string): Site ID from Forter tenant. Console display name: "Forter SiteID".


Example:
```terraform
resource "pingone_davinci_connector_instance" "forterConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "forterConnector"
  }
  name = "My awesome forterConnector"
  properties = jsonencode({
    "apiVersion" = var.forterconnector_property_api_version
    "secretKey" = var.forterconnector_property_secret_key
    "siteId" = var.forterconnector_property_site_id
  })
}
```


## Freshdesk

Connector ID (`connector.id` in the resource): `connectorFreshdesk`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `apiKey` (string): Make sure that the "APIkey:X" is Base64-encoded before pasting into the text field. Console display name: "Freshdesk API Key".
* `baseURL` (string): The <tenant>.freshdesk.com URL or custom domain. Console display name: "Freshdesk Base URL (or Domain)".
* `version` (string): The current Freshdesk API Version. Console display name: "Freshdesk API Version".


Example:
```terraform
resource "pingone_davinci_connector_instance" "connectorFreshdesk" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorFreshdesk"
  }
  name = "My awesome connectorFreshdesk"
  properties = jsonencode({
    "apiKey" = var.connectorfreshdesk_property_api_key
    "baseURL" = var.base_url
    "version" = var.connectorfreshdesk_property_version
  })
}
```


## Freshservice

Connector ID (`connector.id` in the resource): `connectorFreshservice`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `apiKey` (string): Your Freshservice API key. Console display name: "API Key".
* `domain` (string): Your Freshservice domain. Example: https://domain.freshservice.com/. Console display name: "Domain".


Example:
```terraform
resource "pingone_davinci_connector_instance" "connectorFreshservice" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorFreshservice"
  }
  name = "My awesome connectorFreshservice"
  properties = jsonencode({
    "apiKey" = var.connectorfreshservice_property_api_key
    "domain" = var.connectorfreshservice_property_domain
  })
}
```


## Functions

Connector ID (`connector.id` in the resource): `functionsConnector`

*No properties*


Example:
```terraform
resource "pingone_davinci_connector_instance" "functionsConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "functionsConnector"
  }
  name = "My awesome functionsConnector"
}
```


## GBG

Connector ID (`connector.id` in the resource): `gbgConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `password` (string):  Console display name: "GBG Password".
* `requestUrl` (string):  Console display name: "Request URL".
* `soapAction` (string): SOAP Action is a header required for the soap request. Console display name: "Soap Action URL".
* `username` (string):  Console display name: "GBG Username".


Example:
```terraform
resource "pingone_davinci_connector_instance" "gbgConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "gbgConnector"
  }
  name = "My awesome gbgConnector"
  properties = jsonencode({
    "password" = var.gbgconnector_property_password
    "requestUrl" = var.gbgconnector_property_request_url
    "soapAction" = var.gbgconnector_property_soap_action
    "username" = var.gbgconnector_property_username
  })
}
```


## GitHub API

Connector ID (`connector.id` in the resource): `connector-oai-github`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `apiVersion` (string): The GitHub v3 REST API version, such as "2022-11-28". Console display name: "API Version".
* `authBearerToken` (string): The authentication bearer token that has access to GitHub v3 REST API. Console display name: "Authentication Bearer Token".
* `basePath` (string): The base URL for the GitHub API, such as "https://api.github.com". Console display name: "API URL".


Example:
```terraform
resource "pingone_davinci_connector_instance" "connector-oai-github" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connector-oai-github"
  }
  name = "My awesome connector-oai-github"
  properties = jsonencode({
    "apiVersion" = var.connector-oai-github_property_api_version
    "authBearerToken" = var.connector-oai-github_property_auth_bearer_token
    "basePath" = var.connector-oai-github_property_base_path
  })
}
```


## GitHub Login

Connector ID (`connector.id` in the resource): `githubIdpConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `oauth2` (json):  Console display name: "Oauth2 Parameters".


Example:
```terraform
resource "pingone_davinci_connector_instance" "githubIdpConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "githubIdpConnector"
  }
  name = "My awesome githubIdpConnector"
  properties = jsonencode({
    "oauth2" = jsonencode({})
  })
}
```


## Google Analytics (Universal Analytics)

Connector ID (`connector.id` in the resource): `connectorGoogleanalyticsUA`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `trackingID` (string): The tracking ID / web property ID. The format is UA-XXXX-Y. All collected data is associated by this ID. Console display name: "Tracking ID".
* `version` (string): The Protocol version. The current value is '1'. This will only change when there are changes made that are not backwards compatible. Console display name: "Version".


Example:
```terraform
resource "pingone_davinci_connector_instance" "connectorGoogleanalyticsUA" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorGoogleanalyticsUA"
  }
  name = "My awesome connectorGoogleanalyticsUA"
  properties = jsonencode({
    "trackingID" = var.tracking_id
    "version" = var.connectorgoogleanalyticsua_property_version
  })
}
```


## Google Chrome Enterprise Device Trust

Connector ID (`connector.id` in the resource): `connectorGoogleChromeEnterprise`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `customAuth` (json):  Console display name: "Custom Parameters".


Example:
```terraform
resource "pingone_davinci_connector_instance" "connectorGoogleChromeEnterprise" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorGoogleChromeEnterprise"
  }
  name = "My awesome connectorGoogleChromeEnterprise"
  properties = jsonencode({
    "customAuth" = jsonencode({})
  })
}
```


## Google Login

Connector ID (`connector.id` in the resource): `googleConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `openId` (json):  Console display name: "OpenId Parameters".


Example:
```terraform
resource "pingone_davinci_connector_instance" "googleConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "googleConnector"
  }
  name = "My awesome googleConnector"
  properties = jsonencode({
    "openId" = jsonencode({})
  })
}
```


## Google Workspace Admin

Connector ID (`connector.id` in the resource): `googleWorkSpaceAdminConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `iss` (string): The email address associated with the Google Workspace service, such as "google-workspace-admin@xenon-set-123456.iam.gserviceaccount.com". Console display name: "Service Account Email Address".
* `privateKey` (string): The private key associated with the public key that you added to the Google Workspace service. Console display name: "Private Key".
* `sub` (string): The administrator's email address. Console display name: "Admin Email Address".


Example:
```terraform
resource "pingone_davinci_connector_instance" "googleWorkSpaceAdminConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "googleWorkSpaceAdminConnector"
  }
  name = "My awesome googleWorkSpaceAdminConnector"
  properties = jsonencode({
    "iss" = var.googleworkspaceadminconnector_property_iss
    "privateKey" = var.googleworkspaceadminconnector_property_private_key
    "sub" = var.googleworkspaceadminconnector_property_sub
  })
}
```


## HTTP

Connector ID (`connector.id` in the resource): `httpConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `connectionId` (string):  Console display name: "Select an OpenID token management connection for signed HTTP responses.".
* `recaptchaSecretKey` (string): The Secret Key from reCAPTCHA Admin dashboard. Console display name: "reCAPTCHA v2 Secret Key".
* `recaptchaSiteKey` (string): The Site Key from reCAPTCHA Admin dashboard. Console display name: "reCAPTCHA v2 Site Key".
* `whiteList` (string): Enter the hostname for the trusted sites that host your HTML. Note: Ensure that the content hosted on these sites can be trusted and that publishing safeguards are in place to prevent unexpected issues. Console display name: "Trusted Sites".


Example:
```terraform
resource "pingone_davinci_connector_instance" "httpConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "httpConnector"
  }
  name = "My awesome httpConnector"
  properties = jsonencode({
    "connectionId" = var.httpconnector_property_connection_id
    "recaptchaSecretKey" = var.httpconnector_property_recaptcha_secret_key
    "recaptchaSiteKey" = var.httpconnector_property_recaptcha_site_key
    "whiteList" = var.httpconnector_property_white_list
  })
}
```


## HUMAN

Connector ID (`connector.id` in the resource): `connectorHuman`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `humanAuthenticationToken` (string): Bearer Token from HUMAN. Console display name: "HUMAN Authentication Token".
* `humanCustomerID` (string): Customer ID from HUMAN. Console display name: "HUMAN Customer ID".
* `humanPolicyName` (string): HUMAN mitigation policy name. Console display name: "HUMAN Policy Name".


Example:
```terraform
resource "pingone_davinci_connector_instance" "connectorHuman" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorHuman"
  }
  name = "My awesome connectorHuman"
  properties = jsonencode({
    "humanAuthenticationToken" = var.connectorhuman_property_human_authentication_token
    "humanCustomerID" = var.human_customer_id
    "humanPolicyName" = var.connectorhuman_property_human_policy_name
  })
}
```


## HUMAN

Connector ID (`connector.id` in the resource): `humanCompromisedConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `appId` (string): App ID from your HUMAN Tenant. Console display name: "HUMAN App ID".
* `authToken` (string): Auth Token from your HUMAN Tenant. Console display name: "HUMAN Auth Token".


Example:
```terraform
resource "pingone_davinci_connector_instance" "humanCompromisedConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "humanCompromisedConnector"
  }
  name = "My awesome humanCompromisedConnector"
  properties = jsonencode({
    "appId" = var.humancompromisedconnector_property_app_id
    "authToken" = var.humancompromisedconnector_property_auth_token
  })
}
```


## HYPR

Connector ID (`connector.id` in the resource): `hyprConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `customAuth` (json):  Console display name: "Custom Parameters".


Example:
```terraform
resource "pingone_davinci_connector_instance" "hyprConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "hyprConnector"
  }
  name = "My awesome hyprConnector"
  properties = jsonencode({
    "customAuth" = jsonencode({})
  })
}
```


## HYPR Adapt

Connector ID (`connector.id` in the resource): `connectorHyprAdapt`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `accessToken` (string): Access Token. Console display name: "HYPR Adapt Access Token".


Example:
```terraform
resource "pingone_davinci_connector_instance" "connectorHyprAdapt" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorHyprAdapt"
  }
  name = "My awesome connectorHyprAdapt"
  properties = jsonencode({
    "accessToken" = var.connectorhypradapt_property_access_token
  })
}
```


## Have I Been Pwned

Connector ID (`connector.id` in the resource): `haveIBeenPwnedConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `apiKey` (string):  Console display name: "Have I Been Pwned API Key".
* `apiUrl` (string):  Console display name: "API Server URL".
* `userAgent` (string):  


Example:
```terraform
resource "pingone_davinci_connector_instance" "haveIBeenPwnedConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "haveIBeenPwnedConnector"
  }
  name = "My awesome haveIBeenPwnedConnector"
  properties = jsonencode({
    "apiKey" = var.haveibeenpwnedconnector_property_api_key
    "apiUrl" = var.haveibeenpwnedconnector_property_api_url
    "userAgent" = var.haveibeenpwnedconnector_property_user_agent
  })
}
```


## Hellō Connector

Connector ID (`connector.id` in the resource): `connectorHello`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `customAuth` (json):  Console display name: "Custom Parameters".


Example:
```terraform
resource "pingone_davinci_connector_instance" "connectorHello" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorHello"
  }
  name = "My awesome connectorHello"
  properties = jsonencode({
    "customAuth" = jsonencode({})
  })
}
```


## HubSpot Companies API

Connector ID (`connector.id` in the resource): `connector-oai-hubspotcompanies`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `authBearerToken` (string): The authenticating token. Console display name: "Bearer token".
* `basePath` (string): The base URL for contacting the API. Console display name: "Base Path".


Example:
```terraform
resource "pingone_davinci_connector_instance" "connector-oai-hubspotcompanies" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connector-oai-hubspotcompanies"
  }
  name = "My awesome connector-oai-hubspotcompanies"
  properties = jsonencode({
    "authBearerToken" = var.connector-oai-hubspotcompanies_property_auth_bearer_token
    "basePath" = var.connector-oai-hubspotcompanies_property_base_path
  })
}
```


## Hubspot

Connector ID (`connector.id` in the resource): `connectorHubspot`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `bearerToken` (string): Your unique API key. Console display name: "API Key".


Example:
```terraform
resource "pingone_davinci_connector_instance" "connectorHubspot" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorHubspot"
  }
  name = "My awesome connectorHubspot"
  properties = jsonencode({
    "bearerToken" = var.connectorhubspot_property_bearer_token
  })
}
```


## ID DataWeb

Connector ID (`connector.id` in the resource): `idDatawebConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `customAuth` (json):  Console display name: "Custom Parameters".


Example:
```terraform
resource "pingone_davinci_connector_instance" "idDatawebConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "idDatawebConnector"
  }
  name = "My awesome idDatawebConnector"
  properties = jsonencode({
    "customAuth" = jsonencode({})
  })
}
```


## ID R&D

Connector ID (`connector.id` in the resource): `idranddConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `apiKey` (string):  Console display name: "API Key".
* `apiUrl` (string):  Console display name: "API Server URL".


Example:
```terraform
resource "pingone_davinci_connector_instance" "idranddConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "idranddConnector"
  }
  name = "My awesome idranddConnector"
  properties = jsonencode({
    "apiKey" = var.idranddconnector_property_api_key
    "apiUrl" = var.idranddconnector_property_api_url
  })
}
```


## ID.me

Connector ID (`connector.id` in the resource): `idMeConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `oauth2` (json):  Console display name: "Oauth2 Parameters".


Example:
```terraform
resource "pingone_davinci_connector_instance" "idMeConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "idMeConnector"
  }
  name = "My awesome idMeConnector"
  properties = jsonencode({
    "oauth2" = jsonencode({})
  })
}
```


## ID.me - Community Verification

Connector ID (`connector.id` in the resource): `idmecommunityConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `openId` (json):  Console display name: "OpenId Parameters".


Example:
```terraform
resource "pingone_davinci_connector_instance" "idmecommunityConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "idmecommunityConnector"
  }
  name = "My awesome idmecommunityConnector"
  properties = jsonencode({
    "openId" = var.idmecommunityconnector_property_open_id
  })
}
```


## ID.me - Identity Verification

Connector ID (`connector.id` in the resource): `connectorIdMeIdentity`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `openId` (json):  Console display name: "OpenId Parameters".


Example:
```terraform
resource "pingone_davinci_connector_instance" "connectorIdMeIdentity" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorIdMeIdentity"
  }
  name = "My awesome connectorIdMeIdentity"
  properties = jsonencode({
    "openId" = jsonencode({})
  })
}
```


## IDEMIA

Connector ID (`connector.id` in the resource): `idemiaConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `apikey` (string):  Console display name: "API Key".
* `baseUrl` (string): Base Url for IDEMIA API. Can be found in the dashboard documents. Console display name: "IDEMIA API base URL".


Example:
```terraform
resource "pingone_davinci_connector_instance" "idemiaConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "idemiaConnector"
  }
  name = "My awesome idemiaConnector"
  properties = jsonencode({
    "apikey" = var.idemiaconnector_property_apikey
    "baseUrl" = var.idemiaconnector_property_base_url
  })
}
```


## IDI Data

Connector ID (`connector.id` in the resource): `skPeopleIntelligenceConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `authUrl` (string):  Console display name: "Authorization URL".
* `clientId` (string):  Console display name: "Client ID".
* `clientSecret` (string):  Console display name: "Client Secret".
* `dppa` (string):  Console display name: "DPPA".
* `glba` (string):  Console display name: "GLBA".
* `searchUrl` (string):  Console display name: "Search URL".


Example:
```terraform
resource "pingone_davinci_connector_instance" "skPeopleIntelligenceConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "skPeopleIntelligenceConnector"
  }
  name = "My awesome skPeopleIntelligenceConnector"
  properties = jsonencode({
    "authUrl" = var.skpeopleintelligenceconnector_property_auth_url
    "clientId" = var.skpeopleintelligenceconnector_property_client_id
    "clientSecret" = var.skpeopleintelligenceconnector_property_client_secret
    "dppa" = var.skpeopleintelligenceconnector_property_dppa
    "glba" = var.skpeopleintelligenceconnector_property_glba
    "searchUrl" = var.skpeopleintelligenceconnector_property_search_url
  })
}
```


## IDI coreIDENTITY

Connector ID (`connector.id` in the resource): `connectorIdiVERIFIED`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `apiSecret` (string): Please enter your API secret that IDI coreIDENTITY has provided you. Console display name: "API Secret".
* `companyKey` (string): Please enter the company key that IDI coreIDENTITY has assigned. Console display name: "Company Key".
* `idiEnv` (string): Please choose which coreIDENTITY environment you would like to query . Console display name: "Environment".
* `siteKey` (string): Please enter your site key that IDI coreIDENTITY has provided you. Console display name: "Site Key".
* `uniqueUrl` (string): Please enter your unique URL that IDI coreIDENTITY has provided you. Console display name: "Unique URL".


Example:
```terraform
resource "pingone_davinci_connector_instance" "connectorIdiVERIFIED" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorIdiVERIFIED"
  }
  name = "My awesome connectorIdiVERIFIED"
  properties = jsonencode({
    "apiSecret" = var.connectoridiverified_property_api_secret
    "companyKey" = var.connectoridiverified_property_company_key
    "idiEnv" = var.connectoridiverified_property_idi_env
    "siteKey" = var.connectoridiverified_property_site_key
    "uniqueUrl" = var.connectoridiverified_property_unique_url
  })
}
```


## IDmelon

Connector ID (`connector.id` in the resource): `connectorIdmelon`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `customAuth` (json):  Console display name: "Custom Parameters".


Example:
```terraform
resource "pingone_davinci_connector_instance" "connectorIdmelon" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorIdmelon"
  }
  name = "My awesome connectorIdmelon"
  properties = jsonencode({
    "customAuth" = jsonencode({})
  })
}
```


## IDmission

Connector ID (`connector.id` in the resource): `idmissionConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `authDescription` (string):  Console display name: "Authentication Description".
* `connectorName` (string):  Console display name: "Connector Name".
* `description` (string):  Console display name: "Description".
* `details1` (string):  Console display name: "Credentials Details 1".
* `details2` (string):  Console display name: "Credentials Details 2".
* `iconUrl` (string):  Console display name: "Icon URL".
* `iconUrlPng` (string):  Console display name: "Icon URL in PNG".
* `loginId` (string):  Console display name: "Sign On ID".
* `merchantId` (string):  Console display name: "Merchant ID".
* `password` (string):  Console display name: "Password".
* `productId` (string):  Console display name: "Product ID".
* `productName` (string):  Console display name: "Product Name".
* `showCredAddedOn` (boolean):  Console display name: "Show Credentials Added On?".
* `showCredAddedVia` (boolean):  Console display name: "Show Credentials Added through ?".
* `title` (string):  Console display name: "Title".
* `toolTip` (string):  Console display name: "Tooltip".
* `url` (string):  Console display name: "IDmission Server URL".


Example:
```terraform
resource "pingone_davinci_connector_instance" "idmissionConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "idmissionConnector"
  }
  name = "My awesome idmissionConnector"
  properties = jsonencode({
    "authDescription" = var.idmissionconnector_property_auth_description
    "connectorName" = var.idmissionconnector_property_connector_name
    "description" = var.idmissionconnector_property_description
    "details1" = var.idmissionconnector_property_details1
    "details2" = var.idmissionconnector_property_details2
    "iconUrl" = var.idmissionconnector_property_icon_url
    "iconUrlPng" = var.idmissionconnector_property_icon_url_png
    "loginId" = var.idmissionconnector_property_login_id
    "merchantId" = var.idmissionconnector_property_merchant_id
    "password" = var.idmissionconnector_property_password
    "productId" = var.idmissionconnector_property_product_id
    "productName" = var.idmissionconnector_property_product_name
    "showCredAddedOn" = var.idmissionconnector_property_show_cred_added_on
    "showCredAddedVia" = var.idmissionconnector_property_show_cred_added_via
    "title" = var.idmissionconnector_property_title
    "toolTip" = var.idmissionconnector_property_tool_tip
    "url" = var.idmissionconnector_property_url
  })
}
```


## IDmission - OIDC

Connector ID (`connector.id` in the resource): `idmissionOidcConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `customAuth` (json):  Console display name: "Custom Parameters".


Example:
```terraform
resource "pingone_davinci_connector_instance" "idmissionOidcConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "idmissionOidcConnector"
  }
  name = "My awesome idmissionOidcConnector"
  properties = jsonencode({
    "customAuth" = var.idmissionoidcconnector_property_custom_auth
  })
}
```


## IdRamp

Connector ID (`connector.id` in the resource): `idrampOidcConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `customAuth` (json):  Console display name: "Custom Parameters".


Example:
```terraform
resource "pingone_davinci_connector_instance" "idrampOidcConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "idrampOidcConnector"
  }
  name = "My awesome idrampOidcConnector"
  properties = jsonencode({
    "customAuth" = jsonencode({})
  })
}
```


## Image

Connector ID (`connector.id` in the resource): `imageConnector`

*No properties*


Example:
```terraform
resource "pingone_davinci_connector_instance" "imageConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "imageConnector"
  }
  name = "My awesome imageConnector"
}
```


## Incode

Connector ID (`connector.id` in the resource): `incodeConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `customAuth` (json):  Console display name: "Custom Parameters".


Example:
```terraform
resource "pingone_davinci_connector_instance" "incodeConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "incodeConnector"
  }
  name = "My awesome incodeConnector"
  properties = jsonencode({
    "customAuth" = jsonencode({})
  })
}
```


## Infinipoint

Connector ID (`connector.id` in the resource): `connectorInfinipoint`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `customAuth` (json):  Console display name: "Custom Parameters".


Example:
```terraform
resource "pingone_davinci_connector_instance" "connectorInfinipoint" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorInfinipoint"
  }
  name = "My awesome connectorInfinipoint"
  properties = jsonencode({
    "customAuth" = jsonencode({})
  })
}
```


## Intellicheck

Connector ID (`connector.id` in the resource): `intellicheckConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `apiKey` (string): API Key from your Intellicheck tenant. Console display name: "API Key".
* `baseUrl` (string): Base URL from your Intellicheck tenant (Including protocol - https://). Console display name: "Base URL".
* `customerId` (string): Customer ID from your Intellicheck tenant. Console display name: "Customer ID".


Example:
```terraform
resource "pingone_davinci_connector_instance" "intellicheckConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "intellicheckConnector"
  }
  name = "My awesome intellicheckConnector"
  properties = jsonencode({
    "apiKey" = var.intellicheckconnector_property_api_key
    "baseUrl" = var.intellicheckconnector_property_base_url
    "customerId" = var.intellicheckconnector_property_customer_id
  })
}
```


## Jamf

Connector ID (`connector.id` in the resource): `connectorJamf`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `jamfPassword` (string): Enter Password for token. Console display name: "JAMF Password".
* `jamfUsername` (string): Enter Username for token. Console display name: "JAMF Username".
* `serverName` (string): Enter Server Name for Base URL. Console display name: "Server Name".


Example:
```terraform
resource "pingone_davinci_connector_instance" "connectorJamf" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorJamf"
  }
  name = "My awesome connectorJamf"
  properties = jsonencode({
    "jamfPassword" = var.connectorjamf_property_jamf_password
    "jamfUsername" = var.connectorjamf_property_jamf_username
    "serverName" = var.connectorjamf_property_server_name
  })
}
```


## Jira

Connector ID (`connector.id` in the resource): `jiraConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `apiKey` (string): You may need to create a token from Jira with your credentials, if you haven't created one. Console display name: "Jira API token".
* `apiUrl` (string): Base URL of the Jira instance. Console display name: "Base Url".
* `email` (string): Email used for your Jira account. Console display name: "Email Address".


Example:
```terraform
resource "pingone_davinci_connector_instance" "jiraConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "jiraConnector"
  }
  name = "My awesome jiraConnector"
  properties = jsonencode({
    "apiKey" = var.jiraconnector_property_api_key
    "apiUrl" = var.jiraconnector_property_api_url
    "email" = var.jiraconnector_property_email
  })
}
```


## Jira Service Desk

Connector ID (`connector.id` in the resource): `connectorJiraServiceDesk`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `JIRAServiceDeskAuth` (string): Bearer Authorization Token for JIRA Service Desk. Console display name: "Bearer Authorization Token for JIRA Service Desk".
* `JIRAServiceDeskCreateData` (string): Raw JSON body to create new JIRA service desk request. Example: {   "requestParticipants": ["qm:a713c8ea-1075-4e30-9d96-891a7d181739:5ad6d69abfa3980ce712caae"   ],   "serviceDeskId": "10",   "requestTypeId": "25",   "requestFieldValues": {     "summary": "Request JSD help via REST",     "description": "I need a new *mouse* for my Mac"   } }. Console display name: "Raw JSON for creating new JIRA service desk request".
* `JIRAServiceDeskURL` (string): URL for JIRA Service Desk. Example: your-domain.atlassian.net. Console display name: "JIRA Service Desk URL".
* `JIRAServiceDeskUpdateData` (string): Raw JSON body to update JIRA service desk request. Example: {"id": "1","additionalComment": {"body": "I have fixed the problem."}}. Console display name: "Raw JSON for updating JIRA service desk".
* `method` (string): The HTTP Method. Console display name: "Method".


Example:
```terraform
resource "pingone_davinci_connector_instance" "connectorJiraServiceDesk" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorJiraServiceDesk"
  }
  name = "My awesome connectorJiraServiceDesk"
  properties = jsonencode({
    "JIRAServiceDeskAuth" = var.jira_service_desk_auth
    "JIRAServiceDeskCreateData" = var.jira_service_desk_create_data
    "JIRAServiceDeskURL" = var.jira_service_desk_url
    "JIRAServiceDeskUpdateData" = var.jira_service_desk_update_data
    "method" = var.connectorjiraservicedesk_property_method
  })
}
```


## Jumio

Connector ID (`connector.id` in the resource): `jumioConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `apiKey` (string):  Console display name: "API Key".
* `authDescription` (string):  Console display name: "Authentication Description".
* `authUrl` (string):  Console display name: "Base URL for Authentication".
* `authorizationTokenLifetime` (number): default: 1800 (30 minutes). maximum: 5184000 (60 days). Console display name: "Time Transaction URL Valid (seconds)".
* `baseColor` (string): Must be passed with bgColor. Console display name: "HEX Main Color".
* `bgColor` (string): Must be passed with baseColor. Console display name: "HEX Background Color.".
* `callbackUrl` (string):  Console display name: "Callback URL".
* `clientSecret` (string):  Console display name: "API Secret".
* `connectorName` (string):  Console display name: "Connector Name".
* `description` (string):  Console display name: "Description".
* `details1` (string):  Console display name: "Credentials Details 1".
* `details2` (string):  Console display name: "Credentials Details 2".
* `doNotShowInIframe` (boolean): If this is true, user will be redirected to the verification url and then redirected back when complete. Console display name: "Do not show in iFrame".
* `docVerificationUrl` (string):  Console display name: "Document Verification Url".
* `headerImageUrl` (string): Logo must be: landscape (16:9 or 4:3), min. height of 192 pixels, size 8-64 KB. Console display name: "Custom Header Logo URL".
* `iconUrl` (string):  Console display name: "Icon URL".
* `iconUrlPng` (string):  Console display name: "Icon URL in PNG".
* `locale` (string): Renders content in the specified language. Console display name: "Locale".
* `showCredAddedOn` (boolean):  Console display name: "Show Credentials Added On?".
* `showCredAddedVia` (boolean):  Console display name: "Show Credentials Added through ?".
* `title` (string):  Console display name: "Title".
* `toolTip` (string):  Console display name: "Tooltip".


Example:
```terraform
resource "pingone_davinci_connector_instance" "jumioConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "jumioConnector"
  }
  name = "My awesome jumioConnector"
  properties = jsonencode({
    "apiKey" = var.jumioconnector_property_api_key
    "authDescription" = var.jumioconnector_property_auth_description
    "authUrl" = var.jumioconnector_property_auth_url
    "authorizationTokenLifetime" = var.jumioconnector_property_authorization_token_lifetime
    "baseColor" = var.jumioconnector_property_base_color
    "bgColor" = var.jumioconnector_property_bg_color
    "callbackUrl" = var.jumioconnector_property_callback_url
    "clientSecret" = var.jumioconnector_property_client_secret
    "connectorName" = var.jumioconnector_property_connector_name
    "description" = var.jumioconnector_property_description
    "details1" = var.jumioconnector_property_details1
    "details2" = var.jumioconnector_property_details2
    "doNotShowInIframe" = var.jumioconnector_property_do_not_show_in_iframe
    "docVerificationUrl" = var.jumioconnector_property_doc_verification_url
    "headerImageUrl" = var.jumioconnector_property_header_image_url
    "iconUrl" = var.jumioconnector_property_icon_url
    "iconUrlPng" = var.jumioconnector_property_icon_url_png
    "locale" = var.jumioconnector_property_locale
    "showCredAddedOn" = var.jumioconnector_property_show_cred_added_on
    "showCredAddedVia" = var.jumioconnector_property_show_cred_added_via
    "title" = var.jumioconnector_property_title
    "toolTip" = var.jumioconnector_property_tool_tip
  })
}
```


## KBA

Connector ID (`connector.id` in the resource): `kbaConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `authDescription` (string):  Console display name: "Authentication Description".
* `connectorName` (string):  Console display name: "Connector Name".
* `description` (string):  Console display name: "Description".
* `details1` (string):  Console display name: "Credentials Details 1".
* `details2` (string):  Console display name: "Credentials Details 2".
* `formFieldsList` (json):  Console display name: "Fields List".
* `iconUrl` (string):  Console display name: "Icon URL".
* `iconUrlPng` (string):  Console display name: "Icon URL in PNG".
* `showCredAddedOn` (boolean):  Console display name: "Show Credentials Added On?".
* `showCredAddedVia` (boolean):  Console display name: "Show Credentials Added through ?".
* `title` (string):  Console display name: "Title".
* `toolTip` (string):  Console display name: "Tooltip".


Example:
```terraform
resource "pingone_davinci_connector_instance" "kbaConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "kbaConnector"
  }
  name = "My awesome kbaConnector"
  properties = jsonencode({
    "authDescription" = var.kbaconnector_property_auth_description
    "connectorName" = var.kbaconnector_property_connector_name
    "description" = var.kbaconnector_property_description
    "details1" = var.kbaconnector_property_details1
    "details2" = var.kbaconnector_property_details2
    "formFieldsList" = var.kbaconnector_property_form_fields_list
    "iconUrl" = var.kbaconnector_property_icon_url
    "iconUrlPng" = var.kbaconnector_property_icon_url_png
    "showCredAddedOn" = var.kbaconnector_property_show_cred_added_on
    "showCredAddedVia" = var.kbaconnector_property_show_cred_added_via
    "title" = var.kbaconnector_property_title
    "toolTip" = var.kbaconnector_property_tool_tip
  })
}
```


## KYXStart

Connector ID (`connector.id` in the resource): `kyxstartConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `clientId` (string): KYXStart Client ID. Console display name: "Client ID".
* `clientSecret` (string): KYXStart Client Secret. Console display name: "Client Secret".
* `tenantName` (string): Tenant Name from KYXStart Account. Console display name: "Tenant Name".


Example:
```terraform
resource "pingone_davinci_connector_instance" "kyxstartConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "kyxstartConnector"
  }
  name = "My awesome kyxstartConnector"
  properties = jsonencode({
    "clientId" = var.kyxstartconnector_property_client_id
    "clientSecret" = var.kyxstartconnector_property_client_secret
    "tenantName" = var.kyxstartconnector_property_tenant_name
  })
}
```


## Kaizen Secure Voiz

Connector ID (`connector.id` in the resource): `kaizenVoizConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `apiUrl` (string): example: http://<server_root>/ksvvoiceservice/rest/service. Console display name: "API Server URL".
* `applicationName` (string):  Console display name: "Application Name".
* `authDescription` (string):  Console display name: "Authentication Description".
* `connectorName` (string):  Console display name: "Connector Name".
* `description` (string):  Console display name: "Description".
* `details1` (string):  Console display name: "Credentials Details 1".
* `details2` (string):  Console display name: "Credentials Details 2".
* `iconUrl` (string):  Console display name: "Icon URL".
* `iconUrlPng` (string):  Console display name: "Icon URL in PNG".
* `showCredAddedOn` (boolean):  Console display name: "Show Credentials Added On?".
* `showCredAddedVia` (boolean):  Console display name: "Show Credentials Added through ?".
* `title` (string):  Console display name: "Title".
* `toolTip` (string):  Console display name: "Tooltip".


Example:
```terraform
resource "pingone_davinci_connector_instance" "kaizenVoizConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "kaizenVoizConnector"
  }
  name = "My awesome kaizenVoizConnector"
  properties = jsonencode({
    "apiUrl" = var.kaizenvoizconnector_property_api_url
    "applicationName" = var.kaizenvoizconnector_property_application_name
    "authDescription" = var.kaizenvoizconnector_property_auth_description
    "connectorName" = var.kaizenvoizconnector_property_connector_name
    "description" = var.kaizenvoizconnector_property_description
    "details1" = var.kaizenvoizconnector_property_details1
    "details2" = var.kaizenvoizconnector_property_details2
    "iconUrl" = var.kaizenvoizconnector_property_icon_url
    "iconUrlPng" = var.kaizenvoizconnector_property_icon_url_png
    "showCredAddedOn" = var.kaizenvoizconnector_property_show_cred_added_on
    "showCredAddedVia" = var.kaizenvoizconnector_property_show_cred_added_via
    "title" = var.kaizenvoizconnector_property_title
    "toolTip" = var.kaizenvoizconnector_property_tool_tip
  })
}
```


## Keyless

Connector ID (`connector.id` in the resource): `connectorKeyless`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `customAuth` (json):  Console display name: "Custom Parameters".


Example:
```terraform
resource "pingone_davinci_connector_instance" "connectorKeyless" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorKeyless"
  }
  name = "My awesome connectorKeyless"
  properties = jsonencode({
    "customAuth" = jsonencode({})
  })
}
```


## Keyri QR Login

Connector ID (`connector.id` in the resource): `connectorKeyri`

*No properties*


Example:
```terraform
resource "pingone_davinci_connector_instance" "connectorKeyri" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorKeyri"
  }
  name = "My awesome connectorKeyri"
}
```


## LDAP

Connector ID (`connector.id` in the resource): `pingOneLDAPConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `clientId` (string): The Client ID of your PingOne Worker application. Console display name: "Client ID".
* `clientSecret` (string): The Client Secret of your PingOne Worker application. Console display name: "Client Secret".
* `envId` (string): Your PingOne environment ID. Console display name: "Environment ID".
* `gatewayId` (string): Your PingOne LDAP gateway ID. Console display name: "Gateway ID".
* `region` (string): The region in which your PingOne environment exists. Console display name: "Region".


Example:
```terraform
resource "pingone_davinci_connector_instance" "pingOneLDAPConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "pingOneLDAPConnector"
  }
  name = "My awesome pingOneLDAPConnector"
  properties = jsonencode({
    "clientId" = var.pingoneldapconnector_property_client_id
    "clientSecret" = var.pingoneldapconnector_property_client_secret
    "envId" = var.pingoneldapconnector_property_env_id
    "gatewayId" = var.pingoneldapconnector_property_gateway_id
    "region" = var.pingoneldapconnector_property_region
  })
}
```


## LaunchDarkly API

Connector ID (`connector.id` in the resource): `connector-oai-launchdarklyrestapi`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `authApiKey` (string): The authentication key to the LaunchDarkly REST API. Console display name: "API Key".
* `basePath` (string): The base URL for contacting the API. Console display name: "Base Path".


Example:
```terraform
resource "pingone_davinci_connector_instance" "connector-oai-launchdarklyrestapi" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connector-oai-launchdarklyrestapi"
  }
  name = "My awesome connector-oai-launchdarklyrestapi"
  properties = jsonencode({
    "authApiKey" = var.connector-oai-launchdarklyrestapi_property_auth_api_key
    "basePath" = var.connector-oai-launchdarklyrestapi_property_base_path
  })
}
```


## LexisNexis

Connector ID (`connector.id` in the resource): `lexisnexisV2Connector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `apiKey` (string): Your LexisNexis API key, such as “o3x9ywfs26rm1zvl”. Console display name: "API Key".
* `apiUrl` (string): The API URL to target. For a custom value, select Use Custom API URL and enter a value in the Custom API URL field. Console display name: "API URL".
* `orgId` (string): Your LexisNexis organization ID, such as “4en6ll2s”. Console display name: "Organization ID".
* `useCustomApiURL` (string): The API URL to target, such as “https://h.online-metrix.net”. Console display name: "Custom API URL".


Example:
```terraform
resource "pingone_davinci_connector_instance" "lexisnexisV2Connector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "lexisnexisV2Connector"
  }
  name = "My awesome lexisnexisV2Connector"
  properties = jsonencode({
    "apiKey" = var.lexisnexisv2connector_property_api_key
    "apiUrl" = var.lexisnexisv2connector_property_api_url
    "orgId" = var.lexisnexisv2connector_property_org_id
    "useCustomApiURL" = var.use_custom_api_url
  })
}
```


## LinkedIn Login

Connector ID (`connector.id` in the resource): `linkedInConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `oauth2` (json):  Console display name: "Oauth2 Parameters".


Example:
```terraform
resource "pingone_davinci_connector_instance" "linkedInConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "linkedInConnector"
  }
  name = "My awesome linkedInConnector"
  properties = jsonencode({
    "oauth2" = jsonencode({})
  })
}
```


## Location Policy

Connector ID (`connector.id` in the resource): `locationPolicyConnector`

*No properties*


Example:
```terraform
resource "pingone_davinci_connector_instance" "locationPolicyConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "locationPolicyConnector"
  }
  name = "My awesome locationPolicyConnector"
}
```


## Mailchimp

Connector ID (`connector.id` in the resource): `connectorMailchimp`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `transactionalApiKey` (string): The Transactional API Key is used to send data to the transactional API. Console display name: "Transactional API Key".
* `transactionalApiVersion` (string): Mailchimp - Transactional API Version. Console display name: "Transactional API Version".


Example:
```terraform
resource "pingone_davinci_connector_instance" "connectorMailchimp" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorMailchimp"
  }
  name = "My awesome connectorMailchimp"
  properties = jsonencode({
    "transactionalApiKey" = var.connectormailchimp_property_transactional_api_key
    "transactionalApiVersion" = var.connectormailchimp_property_transactional_api_version
  })
}
```


## Mailgun

Connector ID (`connector.id` in the resource): `connectorMailgun`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `apiKey` (string): Mailgun API Key. Console display name: "API Key".
* `apiVersion` (string): Mailgun API Version. Console display name: "API Version".
* `mailgunDomain` (string): Name of the desired domain (e.g. mail.mycompany.com). Console display name: "Domain".


Example:
```terraform
resource "pingone_davinci_connector_instance" "connectorMailgun" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorMailgun"
  }
  name = "My awesome connectorMailgun"
  properties = jsonencode({
    "apiKey" = var.connectormailgun_property_api_key
    "apiVersion" = var.connectormailgun_property_api_version
    "mailgunDomain" = var.connectormailgun_property_mailgun_domain
  })
}
```


## Mailjet API

Connector ID (`connector.id` in the resource): `connector-oai-mailjetapi`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `authPassword` (string): API Secret Key. Console display name: "API Secret Key".
* `authUsername` (string): API Key. Console display name: "API Key".
* `basePath` (string): The base URL for contacting the API. Console display name: "Base Path".


Example:
```terraform
resource "pingone_davinci_connector_instance" "connector-oai-mailjetapi" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connector-oai-mailjetapi"
  }
  name = "My awesome connector-oai-mailjetapi"
  properties = jsonencode({
    "authPassword" = var.connector-oai-mailjetapi_property_auth_password
    "authUsername" = var.connector-oai-mailjetapi_property_auth_username
    "basePath" = var.connector-oai-mailjetapi_property_base_path
  })
}
```


## Melissa Global Address

Connector ID (`connector.id` in the resource): `melissaConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `apiKey` (string): License Key is the API key that you can retrieve from Melissa Admin Portal. Console display name: "License Key".


Example:
```terraform
resource "pingone_davinci_connector_instance" "melissaConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "melissaConnector"
  }
  name = "My awesome melissaConnector"
  properties = jsonencode({
    "apiKey" = var.melissaconnector_property_api_key
  })
}
```


## Microsoft Dynamics - Customer Insights

Connector ID (`connector.id` in the resource): `microsoftDynamicsCustomerInsightsConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `baseURL` (string): Base URL. Console display name: "Base URL".
* `clientId` (string): Client ID. Console display name: "Client ID".
* `clientSecret` (string): Client Secret. Console display name: "Client Secret".
* `environmentName` (string): Environment Name. Console display name: "Environment Name".
* `grantType` (string): Grant Type. Console display name: "Grant Type".
* `tenant` (string): Tenant. Console display name: "Tenant".
* `version` (string): Web API Version. Console display name: "Version".


Example:
```terraform
resource "pingone_davinci_connector_instance" "microsoftDynamicsCustomerInsightsConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "microsoftDynamicsCustomerInsightsConnector"
  }
  name = "My awesome microsoftDynamicsCustomerInsightsConnector"
  properties = jsonencode({
    "baseURL" = var.microsoftdynamicscustomerinsightsconnector_property_base_u_r_l
    "clientId" = var.microsoftdynamicscustomerinsightsconnector_property_client_id
    "clientSecret" = var.microsoftdynamicscustomerinsightsconnector_property_client_secret
    "environmentName" = var.microsoftdynamicscustomerinsightsconnector_property_environment_name
    "grantType" = var.microsoftdynamicscustomerinsightsconnector_property_grant_type
    "tenant" = var.microsoftdynamicscustomerinsightsconnector_property_tenant
    "version" = var.microsoftdynamicscustomerinsightsconnector_property_version
  })
}
```


## Microsoft Edge for Business

Connector ID (`connector.id` in the resource): `connectorMicrosoftEdge`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `customAuth` (json):  Console display name: "Custom Parameters".


Example:
```terraform
resource "pingone_davinci_connector_instance" "connectorMicrosoftEdge" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorMicrosoftEdge"
  }
  name = "My awesome connectorMicrosoftEdge"
  properties = jsonencode({
    "customAuth" = var.connectormicrosoftedge_property_custom_auth
  })
}
```


## Microsoft Intune

Connector ID (`connector.id` in the resource): `connectorMicrosoftIntune`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `clientId` (string): Client ID. Console display name: "Client ID".
* `clientSecret` (string): Client Secret. Console display name: "Client Secret".
* `grantType` (string): Grant Type. Console display name: "Grant Type".
* `scope` (string): Scope. Console display name: "Scope".
* `tenant` (string): Tenant. Console display name: "Tenant".


Example:
```terraform
resource "pingone_davinci_connector_instance" "connectorMicrosoftIntune" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorMicrosoftIntune"
  }
  name = "My awesome connectorMicrosoftIntune"
  properties = jsonencode({
    "clientId" = var.connectormicrosoftintune_property_client_id
    "clientSecret" = var.connectormicrosoftintune_property_client_secret
    "grantType" = var.connectormicrosoftintune_property_grant_type
    "scope" = var.connectormicrosoftintune_property_scope
    "tenant" = var.connectormicrosoftintune_property_tenant
  })
}
```


## Microsoft Login

Connector ID (`connector.id` in the resource): `microsoftIdpConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `openId` (json):  Console display name: "OpenId Parameters".


Example:
```terraform
resource "pingone_davinci_connector_instance" "microsoftIdpConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "microsoftIdpConnector"
  }
  name = "My awesome microsoftIdpConnector"
  properties = jsonencode({
    "openId" = jsonencode({})
  })
}
```


## Microsoft Teams

Connector ID (`connector.id` in the resource): `microsoftTeamsConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `customAuth` (json):  Console display name: "Custom Parameters".


Example:
```terraform
resource "pingone_davinci_connector_instance" "microsoftTeamsConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "microsoftTeamsConnector"
  }
  name = "My awesome microsoftTeamsConnector"
  properties = jsonencode({
    "customAuth" = jsonencode({})
  })
}
```


## NuData Security

Connector ID (`connector.id` in the resource): `nudataConnector`

*No properties*


Example:
```terraform
resource "pingone_davinci_connector_instance" "nudataConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "nudataConnector"
  }
  name = "My awesome nudataConnector"
}
```


## Nuance

Connector ID (`connector.id` in the resource): `nuanceConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `authDescription` (string):  Console display name: "Authentication Description".
* `configSetName` (string): The Config Set Name for accessing Nuance API. Console display name: "Config Set Name".
* `connectorName` (string):  Console display name: "Connector Name".
* `description` (string):  Console display name: "Description".
* `details1` (string):  Console display name: "Credentials Details 1".
* `details2` (string):  Console display name: "Credentials Details 2".
* `iconUrl` (string):  Console display name: "Icon URL".
* `iconUrlPng` (string):  Console display name: "Icon URL in PNG".
* `passphrase1` (string): Passphrase that the user will need to speak for voice sample. Console display name: "Passphrase One".
* `passphrase2` (string): Passphrase that the user will need to speak for voice sample. Console display name: "Passphrase Two".
* `passphrase3` (string): Passphrase that the user will need to speak for voice sample. Console display name: "Passphrase Three".
* `passphrase4` (string): Passphrase that the user will need to speak for voice sample. Console display name: "Passphrase Four".
* `passphrase5` (string): Passphrase that the user will need to speak for voice sample. Console display name: "Passphrase Five".
* `showCredAddedOn` (boolean):  Console display name: "Show Credentials Added On?".
* `showCredAddedVia` (boolean):  Console display name: "Show Credentials Added through ?".
* `title` (string):  Console display name: "Title".
* `toolTip` (string):  Console display name: "Tooltip".


Example:
```terraform
resource "pingone_davinci_connector_instance" "nuanceConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "nuanceConnector"
  }
  name = "My awesome nuanceConnector"
  properties = jsonencode({
    "authDescription" = var.nuanceconnector_property_auth_description
    "configSetName" = var.nuanceconnector_property_config_set_name
    "connectorName" = var.nuanceconnector_property_connector_name
    "description" = var.nuanceconnector_property_description
    "details1" = var.nuanceconnector_property_details1
    "details2" = var.nuanceconnector_property_details2
    "iconUrl" = var.nuanceconnector_property_icon_url
    "iconUrlPng" = var.nuanceconnector_property_icon_url_png
    "passphrase1" = var.nuanceconnector_property_passphrase1
    "passphrase2" = var.nuanceconnector_property_passphrase2
    "passphrase3" = var.nuanceconnector_property_passphrase3
    "passphrase4" = var.nuanceconnector_property_passphrase4
    "passphrase5" = var.nuanceconnector_property_passphrase5
    "showCredAddedOn" = var.nuanceconnector_property_show_cred_added_on
    "showCredAddedVia" = var.nuanceconnector_property_show_cred_added_via
    "title" = var.nuanceconnector_property_title
    "toolTip" = var.nuanceconnector_property_tool_tip
  })
}
```


## OIDC & OAuth IdP

Connector ID (`connector.id` in the resource): `genericConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `customAuth` (json):  Console display name: "Custom Parameters".


Example:
```terraform
resource "pingone_davinci_connector_instance" "genericConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "genericConnector"
  }
  name = "My awesome genericConnector"
  properties = jsonencode({
    "customAuth" = jsonencode({})
  })
}
```


## OPSWAT MetaAccess

Connector ID (`connector.id` in the resource): `connectorOpswat`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `clientID` (string): Oauth client key for authenticating API calls with MetaAccess. Console display name: "Oauth Client Key".
* `clientSecret` (string): Oauth client secret for authenticating API calls with MetaAccess. Console display name: "Oauth Client Secret".
* `crossDomainApiPort` (string): MetaAccess Cross-Domain API integration port. Console display name: "Cross-Domain API Port".
* `maDomain` (string): MetaAccess domain for your environment. Console display name: "MetaAccess Domain".


Example:
```terraform
resource "pingone_davinci_connector_instance" "connectorOpswat" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorOpswat"
  }
  name = "My awesome connectorOpswat"
  properties = jsonencode({
    "clientID" = var.connectoropswat_property_client_i_d
    "clientSecret" = var.connectoropswat_property_client_secret
    "crossDomainApiPort" = var.connectoropswat_property_cross_domain_api_port
    "maDomain" = var.connectoropswat_property_ma_domain
  })
}
```


## OneTrust

Connector ID (`connector.id` in the resource): `oneTrustConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `clientId` (string): Your OneTrust application client ID. Console display name: "Client ID".
* `clientSecret` (string): Your OneTrust application client secret. Console display name: "Client Secret".


Example:
```terraform
resource "pingone_davinci_connector_instance" "oneTrustConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "oneTrustConnector"
  }
  name = "My awesome oneTrustConnector"
  properties = jsonencode({
    "clientId" = var.onetrustconnector_property_client_id
    "clientSecret" = var.onetrustconnector_property_client_secret
  })
}
```


## Onfido

Connector ID (`connector.id` in the resource): `onfidoConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `androidPackageName` (string): Your Android Application's Package Name. Console display name: "Android Application Package Name".
* `apiKey` (string):  Console display name: "API Key".
* `authDescription` (string):  Console display name: "Authentication Description".
* `baseUrl` (string):  Console display name: "Base URL".
* `connectorName` (string):  Console display name: "Connector Name".
* `customizeSteps` (boolean):  Console display name: "Customize Steps".
* `description` (string):  Console display name: "Description".
* `details1` (string):  Console display name: "Credentials Details 1".
* `details2` (string):  Console display name: "Credentials Details 2".
* `iOSBundleId` (string): Your iOS Application's Bundle ID. Console display name: "iOS Application Bundle ID".
* `iconUrl` (string):  Console display name: "Icon URL".
* `iconUrlPng` (string):  Console display name: "Icon URL in PNG".
* `javascriptCSSUrl` (string):  Console display name: "CSS URL".
* `javascriptCdnUrl` (string):  Console display name: "Javascript CDN URL".
* `language` (string):  Console display name: "Language".
* `referenceStepsList` (json):  
* `referrerUrl` (string):  Console display name: "Referrer URL".
* `retrieveReports` (boolean):  Console display name: "Retrieve Reports".
* `shouldCloseOnOverlayClick` (boolean):  Console display name: "Close on Overlay Click".
* `showCredAddedOn` (boolean):  Console display name: "Show Credentials Added On?".
* `showCredAddedVia` (boolean):  Console display name: "Show Credentials Added through ?".
* `stepsList` (boolean): The Proof of Address document capture is currently a BETA feature, and it cannot be used in conjunction with the document and face steps as part of a single SDK flow. Console display name: "ID Verification Steps".
* `title` (string):  Console display name: "Title".
* `toolTip` (string):  Console display name: "Tooltip".
* `useLanguage` (boolean):  Console display name: "Customize Language".
* `useModal` (boolean):  Console display name: "Modal".
* `viewDescriptions` (string):  Console display name: "OnFido Description".
* `viewTitle` (string):  Console display name: "OnFido Title".


Example:
```terraform
resource "pingone_davinci_connector_instance" "onfidoConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "onfidoConnector"
  }
  name = "My awesome onfidoConnector"
  properties = jsonencode({
    "androidPackageName" = var.onfidoconnector_property_android_package_name
    "apiKey" = var.onfidoconnector_property_api_key
    "authDescription" = var.onfidoconnector_property_auth_description
    "baseUrl" = var.onfidoconnector_property_base_url
    "connectorName" = var.onfidoconnector_property_connector_name
    "customizeSteps" = var.onfidoconnector_property_customize_steps
    "description" = var.onfidoconnector_property_description
    "details1" = var.onfidoconnector_property_details1
    "details2" = var.onfidoconnector_property_details2
    "iOSBundleId" = var.onfidoconnector_property_i_o_s_bundle_id
    "iconUrl" = var.onfidoconnector_property_icon_url
    "iconUrlPng" = var.onfidoconnector_property_icon_url_png
    "javascriptCSSUrl" = var.javascript_css_url
    "javascriptCdnUrl" = var.onfidoconnector_property_javascript_cdn_url
    "language" = var.onfidoconnector_property_language
    "referenceStepsList" = var.onfidoconnector_property_reference_steps_list
    "referrerUrl" = var.onfidoconnector_property_referrer_url
    "retrieveReports" = var.onfidoconnector_property_retrieve_reports
    "shouldCloseOnOverlayClick" = var.onfidoconnector_property_should_close_on_overlay_click
    "showCredAddedOn" = var.onfidoconnector_property_show_cred_added_on
    "showCredAddedVia" = var.onfidoconnector_property_show_cred_added_via
    "stepsList" = var.onfidoconnector_property_steps_list
    "title" = var.onfidoconnector_property_title
    "toolTip" = var.onfidoconnector_property_tool_tip
    "useLanguage" = var.onfidoconnector_property_use_language
    "useModal" = var.onfidoconnector_property_use_modal
    "viewDescriptions" = var.onfidoconnector_property_view_descriptions
    "viewTitle" = var.onfidoconnector_property_view_title
  })
}
```


## PaloAlto Prisma Connector

Connector ID (`connector.id` in the resource): `connectorPaloAltoPrisma`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `baseURL` (string): Prisma Base URL. Console display name: "Prisma Base URL".
* `prismaPassword` (string): Secret Key. Console display name: "Prisma - Secret Key".
* `prismaUsername` (string): Access Key. Console display name: "Prisma - Access Key".


Example:
```terraform
resource "pingone_davinci_connector_instance" "connectorPaloAltoPrisma" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorPaloAltoPrisma"
  }
  name = "My awesome connectorPaloAltoPrisma"
  properties = jsonencode({
    "baseURL" = var.base_url
    "prismaPassword" = var.connectorpaloaltoprisma_property_prisma_password
    "prismaUsername" = var.connectorpaloaltoprisma_property_prisma_username
  })
}
```


## PingAccess Administration

Connector ID (`connector.id` in the resource): `connector-oai-pingaccessadministrativeapi`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `authPassword` (string): The password for an account that has access to the PingAccess administrative API. Console display name: "Authenticating Password".
* `authUsername` (string): The username for an account that has access to the PingAccess administrative API. Console display name: "Authenticating Username".
* `basePath` (string): The base URL for the PingAccess Administrative API, such as "https://localhost:9000/pa-admin-api/v3". Console display name: "API URL".
* `sslVerification` (string): When enabled, DaVinci verifies the PingAccess SSL certificate and uses encrypted communication. Console display name: "Use SSL Verification".


Example:
```terraform
resource "pingone_davinci_connector_instance" "connector-oai-pingaccessadministrativeapi" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connector-oai-pingaccessadministrativeapi"
  }
  name = "My awesome connector-oai-pingaccessadministrativeapi"
  properties = jsonencode({
    "authPassword" = var.connector-oai-pingaccessadministrativeapi_property_auth_password
    "authUsername" = var.connector-oai-pingaccessadministrativeapi_property_auth_username
    "basePath" = var.connector-oai-pingaccessadministrativeapi_property_base_path
    "sslVerification" = var.connector-oai-pingaccessadministrativeapi_property_ssl_verification
  })
}
```


## PingFederate

Connector ID (`connector.id` in the resource): `pingFederateConnectorV2`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `openId` (json):  Console display name: "OpenId Parameters".


Example:
```terraform
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
```


## PingFederate Administration

Connector ID (`connector.id` in the resource): `connector-oai-pfadminapi`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `authPassword` (string): The password for an account that has access to the PingFederate administrative API. Console display name: "Authenticating Password".
* `authUsername` (string): The username for an account that has access to the PingFederate administrative API. Console display name: "Authenticating Username".
* `basePath` (string): The base URL for the PingFederate administrative API, such as "https://8.8.4.4:9999/pf-admin-api/v1". Console display name: "API URL".
* `sslVerification` (string): When enabled, DaVinci verifies the PingFederate SSL certificate and uses encrypted communication. Console display name: "Use SSL Verification".


Example:
```terraform
resource "pingone_davinci_connector_instance" "connector-oai-pfadminapi" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connector-oai-pfadminapi"
  }
  name = "My awesome connector-oai-pfadminapi"
  properties = jsonencode({
    "authPassword" = var.connector-oai-pfadminapi_property_auth_password
    "authUsername" = var.connector-oai-pfadminapi_property_auth_username
    "basePath" = var.connector-oai-pfadminapi_property_base_path
    "sslVerification" = var.connector-oai-pfadminapi_property_ssl_verification
  })
}
```


## PingID

Connector ID (`connector.id` in the resource): `pingIdConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `customAuth` (json):  Console display name: "Custom Parameters".


Example:
```terraform
resource "pingone_davinci_connector_instance" "pingIdConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "pingIdConnector"
  }
  name = "My awesome pingIdConnector"
  properties = jsonencode({
    "customAuth" = jsonencode({
				"properties": {
				  "pingIdProperties": {
					"displayName": "PingID properties file",
					"preferredControlType": "secureTextArea",
					"hashedVisibility": true,
					"required": true,
					"info": "Paste the contents of the PingID properties file into this field.",
					"value": "${file(var.pingidconnector_property_pingid_properties_file_path)}"
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
```


## PingOne

Connector ID (`connector.id` in the resource): `pingOneSSOConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `clientId` (string): The Client ID of your PingOne Worker application. Console display name: "Client ID".
* `clientSecret` (string): The Client Secret of your PingOne Worker application. Console display name: "Client Secret".
* `envId` (string): Your PingOne environment ID. Console display name: "Environment ID".
* `envRegionInfo` (string): If you want to connect with a different PingOne environment, enter the environment and credential information below. Console display name: "The default PingOne environment is configured automatically.".
* `region` (string): The region in which your PingOne environment exists. Console display name: "Region".


Example:
```terraform
resource "pingone_davinci_connector_instance" "pingOneSSOConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "pingOneSSOConnector"
  }
  name = "My awesome pingOneSSOConnector"
  properties = jsonencode({
    "clientId" = var.pingone_worker_app_client_id
    "clientSecret" = var.pingone_worker_app_client_secret
    "envId" = var.pingone_worker_app_environment_id
    "envRegionInfo" = var.pingonessoconnector_property_env_region_info
    "region" = var.pingonessoconnector_property_region
  })
}
```


## PingOne Advanced Identity Cloud Access Request

Connector ID (`connector.id` in the resource): `accessRequestConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `baseURL` (string): The API URL to target. Console display name: "Identity Cloud Base URL".
* `endUserClientId` (string): The Client ID from the end user account. Console display name: "End User Client ID".
* `endUserClientPrivateKey` (string): The Client Private Key from the end user account. Console display name: "End User Client Private Key".
* `realm` (string): The Realm configured in Identity Cloud. Console display name: "Realm".
* `serviceAccountId` (string): The account ID for your Identity Cloud service account. You can find this ID under the account settings of your service account. Console display name: "Service Account ID".
* `serviceAccountPrivateKey` (string): The private key for your Identity Cloud service account. You can find this private key under the account settings of your service account. Console display name: "Service Account Private Key".


Example:
```terraform
resource "pingone_davinci_connector_instance" "accessRequestConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "accessRequestConnector"
  }
  name = "My awesome accessRequestConnector"
  properties = jsonencode({
    "baseURL" = var.accessrequestconnector_property_base_u_r_l
    "endUserClientId" = var.accessrequestconnector_property_end_user_client_id
    "endUserClientPrivateKey" = var.accessrequestconnector_property_end_user_client_private_key
    "realm" = var.accessrequestconnector_property_realm
    "serviceAccountId" = var.accessrequestconnector_property_service_account_id
    "serviceAccountPrivateKey" = var.accessrequestconnector_property_service_account_private_key
  })
}
```


## PingOne Advanced Identity Cloud Login Connector

Connector ID (`connector.id` in the resource): `pingoneAdvancedIdentityCloudLoginConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `openId` (json):  Console display name: "OpenId Parameters".


Example:
```terraform
resource "pingone_davinci_connector_instance" "pingoneAdvancedIdentityCloudLoginConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "pingoneAdvancedIdentityCloudLoginConnector"
  }
  name = "My awesome pingoneAdvancedIdentityCloudLoginConnector"
  properties = jsonencode({
    "openId" = var.pingoneadvancedidentitycloudloginconnector_property_open_id
  })
}
```


## PingOne Authentication

Connector ID (`connector.id` in the resource): `pingOneAuthenticationConnector`

*No properties*


Example:
```terraform
resource "pingone_davinci_connector_instance" "pingOneAuthenticationConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "pingOneAuthenticationConnector"
  }
  name = "My awesome pingOneAuthenticationConnector"
}
```


## PingOne Authorize

Connector ID (`connector.id` in the resource): `pingOneAuthorizeConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `clientId` (string): The Client ID of the PingOne worker application. Console display name: "Client ID".
* `clientSecret` (string): The Client Secret of the PingOne worker application. Console display name: "Client Secret".
* `endpointURL` (string): The PingOne Authorize decision endpoint or ID to which the connector submits decision requests. Console display name: "Endpoint".


Example:
```terraform
resource "pingone_davinci_connector_instance" "pingOneAuthorizeConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "pingOneAuthorizeConnector"
  }
  name = "My awesome pingOneAuthorizeConnector"
  properties = jsonencode({
    "clientId" = var.pingoneauthorizeconnector_property_client_id
    "clientSecret" = var.pingoneauthorizeconnector_property_client_secret
    "endpointURL" = var.endpoint_url
  })
}
```


## PingOne Authorize - API Access Management

Connector ID (`connector.id` in the resource): `pingauthadapter`

*No properties*


Example:
```terraform
resource "pingone_davinci_connector_instance" "pingauthadapter" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "pingauthadapter"
  }
  name = "My awesome pingauthadapter"
}
```


## PingOne Credentials

Connector ID (`connector.id` in the resource): `pingOneCredentialsConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `clientId` (string): The Client ID of your PingOne Worker application. Console display name: "Client ID".
* `clientSecret` (string): The Client Secret of your PingOne Worker application. Console display name: "Client Secret".
* `digitalWalletApplicationId` (string): Identifier (UUID) associated with the credential digital wallet app. Console display name: "Digital Wallet Application ID".
* `envId` (string): Your PingOne Environment ID. Console display name: "Environment ID".
* `region` (string): The region your PingOne environment is in. Console display name: "Region".


Example:
```terraform
resource "pingone_davinci_connector_instance" "pingOneCredentialsConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "pingOneCredentialsConnector"
  }
  name = "My awesome pingOneCredentialsConnector"
  properties = jsonencode({
    "clientId" = var.pingone_worker_app_client_id
    "clientSecret" = var.pingone_worker_app_client_secret
    "digitalWalletApplicationId" = var.pingonecredentialsconnector_property_digital_wallet_application_id
    "envId" = var.pingone_worker_app_environment_id
    "region" = var.pingonecredentialsconnector_property_region
  })
}
```


## PingOne Forms

Connector ID (`connector.id` in the resource): `pingOneFormsConnector`

*No properties*


Example:
```terraform
resource "pingone_davinci_connector_instance" "pingOneFormsConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "pingOneFormsConnector"
  }
  name = "My awesome pingOneFormsConnector"
}
```


## PingOne MFA

Connector ID (`connector.id` in the resource): `pingOneMfaConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `clientId` (string): The Client ID of your PingOne Worker application. Console display name: "Client ID".
* `clientSecret` (string): The Client Secret of your PingOne Worker application. Console display name: "Client Secret".
* `envId` (string): Your PingOne Environment ID. Console display name: "Environment ID".
* `policyId` (string): The ID of your PingOne MFA device authentication policy. Console display name: "Policy ID".
* `region` (string): The region in which your PingOne environment exists. Console display name: "Region".


Example:
```terraform
resource "pingone_davinci_connector_instance" "pingOneMfaConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "pingOneMfaConnector"
  }
  name = "My awesome pingOneMfaConnector"
  properties = jsonencode({
    "clientId" = var.pingone_worker_app_client_id
    "clientSecret" = var.pingone_worker_app_client_secret
    "envId" = var.pingone_worker_app_environment_id
    "policyId" = var.pingonemfaconnector_property_policy_id
    "region" = var.pingonemfaconnector_property_region
  })
}
```


## PingOne Notifications

Connector ID (`connector.id` in the resource): `notificationsConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `clientId` (string): The Client ID of your PingOne Worker application. Console display name: "Client ID".
* `clientSecret` (string): The Client Secret of your PingOne Worker application. Console display name: "Client Secret".
* `envId` (string): Your PingOne Environment ID. Console display name: "Environment ID".
* `notificationPolicyId` (string): A unique identifier for the policy. Console display name: "Notification Policy ID".
* `region` (string): The region in which your PingOne environment exists. Console display name: "Region".


Example:
```terraform
resource "pingone_davinci_connector_instance" "notificationsConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "notificationsConnector"
  }
  name = "My awesome notificationsConnector"
  properties = jsonencode({
    "clientId" = var.notificationsconnector_property_client_id
    "clientSecret" = var.notificationsconnector_property_client_secret
    "envId" = var.notificationsconnector_property_env_id
    "notificationPolicyId" = var.notificationsconnector_property_notification_policy_id
    "region" = var.notificationsconnector_property_region
  })
}
```


## PingOne Protect

Connector ID (`connector.id` in the resource): `pingOneRiskConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `clientId` (string): The id for your Application found in Ping's Dashboard. Console display name: "Client ID".
* `clientSecret` (string): Client Secret from your App in Ping's Dashboard. Console display name: "Client Secret".
* `envId` (string): Your Environment ID provided by Ping. Console display name: "Environment ID".
* `region` (string): The region your PingOne environment is in. Console display name: "Region".


Example:
```terraform
resource "pingone_davinci_connector_instance" "pingOneRiskConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "pingOneRiskConnector"
  }
  name = "My awesome pingOneRiskConnector"
  properties = jsonencode({
    "clientId" = var.pingone_worker_app_client_id
    "clientSecret" = var.pingone_worker_app_client_secret
    "envId" = var.pingone_worker_app_environment_id
    "region" = var.pingoneriskconnector_property_region
  })
}
```


## PingOne RADIUS Gateway

Connector ID (`connector.id` in the resource): `pingOneIntegrationsConnector`

*No properties*


Example:
```terraform
resource "pingone_davinci_connector_instance" "pingOneIntegrationsConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "pingOneIntegrationsConnector"
  }
  name = "My awesome pingOneIntegrationsConnector"
}
```


## PingOne Scope Consent

Connector ID (`connector.id` in the resource): `pingOneScopeConsentConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `clientId` (string): The Client ID of your PingOne Worker application. Console display name: "Client ID".
* `clientSecret` (string): The Client Secret of your PingOne Worker application. Console display name: "Client Secret".
* `envId` (string): Your PingOne Environment ID. Console display name: "Environment ID".
* `region` (string): The region in which your PingOne environment exists. Console display name: "Region".


Example:
```terraform
resource "pingone_davinci_connector_instance" "pingOneScopeConsentConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "pingOneScopeConsentConnector"
  }
  name = "My awesome pingOneScopeConsentConnector"
  properties = jsonencode({
    "clientId" = var.pingone_worker_app_client_id
    "clientSecret" = var.pingone_worker_app_client_secret
    "envId" = var.pingone_worker_app_environment_id
    "region" = var.pingonescopeconsentconnector_property_region
  })
}
```


## PingOne Verify

Connector ID (`connector.id` in the resource): `pingOneVerifyConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `clientId` (string): The Client ID of your PingOne Worker application. Console display name: "Client ID".
* `clientSecret` (string): The Client Secret of your PingOne Worker application. Console display name: "Client Secret".
* `envId` (string): Your PingOne Environment ID. Console display name: "Environment ID".
* `region` (string): The region your PingOne environment is in. Console display name: "Region".


Example:
```terraform
resource "pingone_davinci_connector_instance" "pingOneVerifyConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "pingOneVerifyConnector"
  }
  name = "My awesome pingOneVerifyConnector"
  properties = jsonencode({
    "clientId" = var.pingone_worker_app_client_id
    "clientSecret" = var.pingone_worker_app_client_secret
    "envId" = var.pingone_worker_app_environment_id
    "region" = var.pingoneverifyconnector_property_region
  })
}
```


## Private ID

Connector ID (`connector.id` in the resource): `privateidConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `customAuth` (json):  Console display name: "Custom Parameters".


Example:
```terraform
resource "pingone_davinci_connector_instance" "privateidConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "privateidConnector"
  }
  name = "My awesome privateidConnector"
  properties = jsonencode({
    "customAuth" = var.privateidconnector_property_custom_auth
  })
}
```


## Prove

Connector ID (`connector.id` in the resource): `payfoneConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `appClientId` (string):  Console display name: "App Client ID".
* `baseUrl` (string):  Console display name: "Prove Base URL".
* `clientId` (string):  Console display name: "Client ID".
* `password` (string):  Console display name: "Password".
* `simulatorMode` (boolean):  Console display name: "Simulator Mode?".
* `simulatorPhoneNumber` (string):  Console display name: "Simulator Phone Number".
* `skCallbackBaseUrl` (string): Use this url as the callback base URL. Console display name: "Callback Base URL".
* `username` (string):  Console display name: "Username".


Example:
```terraform
resource "pingone_davinci_connector_instance" "payfoneConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "payfoneConnector"
  }
  name = "My awesome payfoneConnector"
  properties = jsonencode({
    "appClientId" = var.payfoneconnector_property_app_client_id
    "baseUrl" = var.payfoneconnector_property_base_url
    "clientId" = var.payfoneconnector_property_client_id
    "password" = var.payfoneconnector_property_password
    "simulatorMode" = var.payfoneconnector_property_simulator_mode
    "simulatorPhoneNumber" = var.payfoneconnector_property_simulator_phone_number
    "skCallbackBaseUrl" = var.payfoneconnector_property_sk_callback_base_url
    "username" = var.payfoneconnector_property_username
  })
}
```


## Prove International

Connector ID (`connector.id` in the resource): `proveConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `baseUrl` (string):  Console display name: "Prove Base URL".
* `clientId` (string):  Console display name: "Prove Client ID".
* `grantType` (string):  Console display name: "Prove Grant Type".
* `password` (string):  Console display name: "Prove Password".
* `username` (string):  Console display name: "Prove Username".


Example:
```terraform
resource "pingone_davinci_connector_instance" "proveConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "proveConnector"
  }
  name = "My awesome proveConnector"
  properties = jsonencode({
    "baseUrl" = var.proveconnector_property_base_url
    "clientId" = var.proveconnector_property_client_id
    "grantType" = var.proveconnector_property_grant_type
    "password" = var.proveconnector_property_password
    "username" = var.proveconnector_property_username
  })
}
```


## RSA

Connector ID (`connector.id` in the resource): `rsaConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `accessId` (string): RSA Access ID from Administration API key file. Console display name: "Access ID".
* `accessKey` (string): RSA Access Key from Administration API key file. Console display name: "Access Key".
* `baseUrl` (string): Base URL for RSA API that is provided in Administration API key file. Console display name: "Base URL".


Example:
```terraform
resource "pingone_davinci_connector_instance" "rsaConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "rsaConnector"
  }
  name = "My awesome rsaConnector"
  properties = jsonencode({
    "accessId" = var.rsaconnector_property_access_id
    "accessKey" = var.rsaconnector_property_access_key
    "baseUrl" = var.rsaconnector_property_base_url
  })
}
```


## ReadID by Inverid

Connector ID (`connector.id` in the resource): `inveridConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `getApiKey` (string): Viewer API Key provided to you by Inverid. Console display name: "ReadID Viewer API Key".
* `host` (string): Hostname provided to you by Inverid. Console display name: "ReadID Hostname".
* `postApiKey` (string): Submitter API Key provided to you by Inverid. Console display name: "ReadID Submitter API Key".
* `skWebhookUri` (string): Use this url as the Webhook URL in the Third Party Integration's configuration. Console display name: "Redirect Webhook URI".
* `timeToLive` (string): Specify the duration (in minutes) a users session should stay active. Value must be between 30 and 72000. Console display name: "Time to live for ReadySession".


Example:
```terraform
resource "pingone_davinci_connector_instance" "inveridConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "inveridConnector"
  }
  name = "My awesome inveridConnector"
  properties = jsonencode({
    "getApiKey" = var.inveridconnector_property_get_api_key
    "host" = var.inveridconnector_property_host
    "postApiKey" = var.inveridconnector_property_post_api_key
    "skWebhookUri" = var.inveridconnector_property_sk_webhook_uri
    "timeToLive" = var.inveridconnector_property_time_to_live
  })
}
```


## SAML

Connector ID (`connector.id` in the resource): `samlConnector`

*No properties*


Example:
```terraform
resource "pingone_davinci_connector_instance" "samlConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "samlConnector"
  }
  name = "My awesome samlConnector"
}
```


## SAML IdP

Connector ID (`connector.id` in the resource): `samlIdpConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `saml` (json):  Console display name: "SAML Parameters".


Example:
```terraform
resource "pingone_davinci_connector_instance" "samlIdpConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "samlIdpConnector"
  }
  name = "My awesome samlIdpConnector"
  properties = jsonencode({
    "saml" = jsonencode({})
  })
}
```


## SAP Identity API

Connector ID (`connector.id` in the resource): `connector-oai-sapidentityapis`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `authApiKey` (string): The authentication key to the SAP Identity APIs. Console display name: "API Key".
* `basePath` (string): The base URL for contacting the API. Console display name: "Base Path".


Example:
```terraform
resource "pingone_davinci_connector_instance" "connector-oai-sapidentityapis" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connector-oai-sapidentityapis"
  }
  name = "My awesome connector-oai-sapidentityapis"
  properties = jsonencode({
    "authApiKey" = var.connector-oai-sapidentityapis_property_auth_api_key
    "basePath" = var.connector-oai-sapidentityapis_property_base_path
  })
}
```


## SEON

Connector ID (`connector.id` in the resource): `seonConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `baseURL` (string): The API URL to target. Console display name: "API Base URL".
* `licenseKey` (string): Your SEON license key. For help, see the SEON REST API documentation. Console display name: "License Key".


Example:
```terraform
resource "pingone_davinci_connector_instance" "seonConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "seonConnector"
  }
  name = "My awesome seonConnector"
  properties = jsonencode({
    "baseURL" = var.base_url
    "licenseKey" = var.seonconnector_property_license_key
  })
}
```


## SMTP Client

Connector ID (`connector.id` in the resource): `smtpConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `hostname` (string): Example: smtp-relay.gmail.com. Console display name: "SMTP Server/Host".
* `name` (string): Optional hostname of the client, used for identifying to the server, defaults to hostname of the machine. Console display name: "Client Name".
* `password` (string):  Console display name: "Password".
* `port` (number): Example: 25. Console display name: "SMTP Port".
* `secureFlag` (boolean):  Console display name: "Secure Flag?".
* `username` (string):  Console display name: "Username".


Example:
```terraform
resource "pingone_davinci_connector_instance" "smtpConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "smtpConnector"
  }
  name = "My awesome smtpConnector"
  properties = jsonencode({
    "hostname" = var.smtpconnector_property_hostname
    "name" = var.smtpconnector_property_name
    "password" = var.smtpconnector_property_password
    "port" = var.smtpconnector_property_port
    "secureFlag" = var.smtpconnector_property_secure_flag
    "username" = var.smtpconnector_property_username
  })
}
```


## SailPoint IdentityNow

Connector ID (`connector.id` in the resource): `connectorIdentityNow`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `clientId` (string): Client Id for your client found in IdentityNow's Dashboard. Console display name: "Client ID".
* `clientSecret` (string): Client Secret from your client in IdentityNow's Dashboard. Console display name: "Client Secret".
* `tenant` (string): The org name is displayed within the Org Details section of the dashboard. Console display name: "IdentityNow Tenant".


Example:
```terraform
resource "pingone_davinci_connector_instance" "connectorIdentityNow" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorIdentityNow"
  }
  name = "My awesome connectorIdentityNow"
  properties = jsonencode({
    "clientId" = var.connectoridentitynow_property_client_id
    "clientSecret" = var.connectoridentitynow_property_client_secret
    "tenant" = var.connectoridentitynow_property_tenant
  })
}
```


## Salesforce

Connector ID (`connector.id` in the resource): `salesforceConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `adminUsername` (string): The username of your Salesforce administrator account. Console display name: "Username".
* `consumerKey` (string): The consumer key shown on your Salesforce connected app. Console display name: "Consumer Key".
* `domainName` (string): Your Salesforce domain name, such as "mycompany-dev-ed". Console display name: "Domain Name".
* `environment` (string): If the environment you specify in the Domain Name field is part of a sandbox organization, select Sandbox. Otherwise, select Production. Console display name: "Environment".
* `privateKey` (string): The private key that corresponds to the X.509 certificate you added to your Salesforce connected app. Console display name: "Private Key".


Example:
```terraform
resource "pingone_davinci_connector_instance" "salesforceConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "salesforceConnector"
  }
  name = "My awesome salesforceConnector"
  properties = jsonencode({
    "adminUsername" = var.salesforceconnector_property_admin_username
    "consumerKey" = var.salesforceconnector_property_consumer_key
    "domainName" = var.salesforceconnector_property_domain_name
    "environment" = var.salesforceconnector_property_environment
    "privateKey" = var.salesforceconnector_property_private_key
  })
}
```


## Salesforce Marketing Cloud (BETA)

Connector ID (`connector.id` in the resource): `connectorSalesforceMarketingCloud`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `SalesforceMarketingCloudURL` (string): URL for Salesforce Marketing Cloud. Example: https://YOUR_SUBDOMAIN.rest.marketingcloudapis.com. Console display name: "Salesforce Marketing Cloud URL".
* `accountId` (string): Account identifier, or MID, of the target business unit. Use to switch between business units. If you don’t specify account_id, the returned access token is in the context of the business unit that created the integration. Console display name: "Account ID".
* `clientId` (string): Client ID issued when you create the API integration in Installed Packages. Console display name: "Client ID".
* `clientSecret` (string): Client secret issued when you create the API integration in Installed Packages. Console display name: "Client Secret".
* `scope` (string): Space-separated list of data-access permissions for your application. Console display name: "Scope".


Example:
```terraform
resource "pingone_davinci_connector_instance" "connectorSalesforceMarketingCloud" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorSalesforceMarketingCloud"
  }
  name = "My awesome connectorSalesforceMarketingCloud"
  properties = jsonencode({
    "SalesforceMarketingCloudURL" = var.salesforce_marketing_cloud_url
    "accountId" = var.connectorsalesforcemarketingcloud_property_account_id
    "clientId" = var.connectorsalesforcemarketingcloud_property_client_id
    "clientSecret" = var.connectorsalesforcemarketingcloud_property_client_secret
    "scope" = var.connectorsalesforcemarketingcloud_property_scope
  })
}
```


## Saviynt Connector Flows

Connector ID (`connector.id` in the resource): `connectorSaviyntFlow`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `domainName` (string): Provide your Saviynt domain name. Console display name: "Saviynt Domain Name".
* `path` (string): Provide your Saviynt path name. Console display name: "Saviynt Path Name".
* `saviyntPassword` (string): Provide your Saviynt password. Console display name: "Saviynt Password".
* `saviyntUserName` (string): Provide your Saviynt user name. Console display name: "Saviynt User Name".


Example:
```terraform
resource "pingone_davinci_connector_instance" "connectorSaviyntFlow" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorSaviyntFlow"
  }
  name = "My awesome connectorSaviyntFlow"
  properties = jsonencode({
    "domainName" = var.connectorsaviyntflow_property_domain_name
    "path" = var.connectorsaviyntflow_property_path
    "saviyntPassword" = var.connectorsaviyntflow_property_saviynt_password
    "saviyntUserName" = var.connectorsaviyntflow_property_saviynt_user_name
  })
}
```


## Screen

Connector ID (`connector.id` in the resource): `screenConnector`

*No properties*


Example:
```terraform
resource "pingone_davinci_connector_instance" "screenConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "screenConnector"
  }
  name = "My awesome screenConnector"
}
```


## SecurID

Connector ID (`connector.id` in the resource): `securIdConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `apiUrl` (string): The URL of your SecurID authentication API, such as "https://company.auth.securid.com". Console display name: "SecurID Authentication API REST URL".
* `clientKey` (string): Your SecurID authentication client key, such as "vowc450ahs6nry66vok0pvaizwnfr43ewsqcm7tz". Console display name: "Client Key".


Example:
```terraform
resource "pingone_davinci_connector_instance" "securIdConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "securIdConnector"
  }
  name = "My awesome securIdConnector"
  properties = jsonencode({
    "apiUrl" = var.securidconnector_property_api_url
    "clientKey" = var.securidconnector_property_client_key
  })
}
```


## Securonix

Connector ID (`connector.id` in the resource): `connectorSecuronix`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `domainName` (string): Domain Name. Console display name: "Domain Name".
* `token` (string): Token for authentication. Console display name: "Token".


Example:
```terraform
resource "pingone_davinci_connector_instance" "connectorSecuronix" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorSecuronix"
  }
  name = "My awesome connectorSecuronix"
  properties = jsonencode({
    "domainName" = var.connectorsecuronix_property_domain_name
    "token" = var.connectorsecuronix_property_token
  })
}
```


## Segment

Connector ID (`connector.id` in the resource): `connectorSegment`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `version` (string): Segment - HTTP Tracking API Version. Console display name: "HTTP Tracking API Version".
* `writeKey` (string): The Write Key is used to send data to a specific workplace. Console display name: "Write Key".


Example:
```terraform
resource "pingone_davinci_connector_instance" "connectorSegment" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorSegment"
  }
  name = "My awesome connectorSegment"
  properties = jsonencode({
    "version" = var.connectorsegment_property_version
    "writeKey" = var.connectorsegment_property_write_key
  })
}
```


## SentiLink

Connector ID (`connector.id` in the resource): `sentilinkConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `account` (string): Account ID of SentiLink. Console display name: "Account ID".
* `apiUrl` (string):  Console display name: "API URL".
* `javascriptCdnUrl` (string):  Console display name: "Javascript CDN URL".
* `token` (string): Token ID for SentiLink account. Console display name: "Token ID".


Example:
```terraform
resource "pingone_davinci_connector_instance" "sentilinkConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "sentilinkConnector"
  }
  name = "My awesome sentilinkConnector"
  properties = jsonencode({
    "account" = var.sentilinkconnector_property_account
    "apiUrl" = var.sentilinkconnector_property_api_url
    "javascriptCdnUrl" = var.sentilinkconnector_property_javascript_cdn_url
    "token" = var.sentilinkconnector_property_token
  })
}
```


## ServiceNow

Connector ID (`connector.id` in the resource): `servicenowConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `adminUsername` (string): Your ServiceNow administrator username. Console display name: "Username".
* `apiUrl` (string): The API URL to target, such as "https://mycompany.service-now.com". Console display name: "API URL".
* `password` (string): Your ServiceNow administrator password. Console display name: "Password".


Example:
```terraform
resource "pingone_davinci_connector_instance" "servicenowConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "servicenowConnector"
  }
  name = "My awesome servicenowConnector"
  properties = jsonencode({
    "adminUsername" = var.servicenowconnector_property_admin_username
    "apiUrl" = var.servicenowconnector_property_api_url
    "password" = var.servicenowconnector_property_password
  })
}
```


## Shopify Connector

Connector ID (`connector.id` in the resource): `connectorShopify`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `accessToken` (string): Your store's unique Admin API Access Token that goes into the X-Shopify-Access-Token property. Required scopes when generating Admin API Access Token: 'read_customers' and 'write_customers'. Note any Custom Shopify API calls you intend to use with this connector via Make Custom API Call capability, will have to be added as well. Console display name: "Admin API Access Token".
* `apiVersion` (string): The Shopify version name ( ex. 2022-04 ). Console display name: "API Version Name".
* `multipassSecret` (string): Shopify Multipass Secret. Console display name: "Multipass Secret".
* `multipassStoreDomain` (string): Shopify Multipass Store Domain (yourstorename.myshopify.com). Console display name: "Multipass Store Domain".
* `yourStoreName` (string): The name of your store as Shopify identifies you ( first text that comes after HTTPS:// ). Console display name: "Store Name".


Example:
```terraform
resource "pingone_davinci_connector_instance" "connectorShopify" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorShopify"
  }
  name = "My awesome connectorShopify"
  properties = jsonencode({
    "accessToken" = var.connectorshopify_property_access_token
    "apiVersion" = var.connectorshopify_property_api_version
    "multipassSecret" = var.connectorshopify_property_multipass_secret
    "multipassStoreDomain" = var.connectorshopify_property_multipass_store_domain
    "yourStoreName" = var.connectorshopify_property_your_store_name
  })
}
```


## Sift

Connector ID (`connector.id` in the resource): `siftConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `apiKey` (string): API Key from Sift Tenant. Console display name: "Sift API Key".


Example:
```terraform
resource "pingone_davinci_connector_instance" "siftConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "siftConnector"
  }
  name = "My awesome siftConnector"
  properties = jsonencode({
    "apiKey" = var.siftconnector_property_api_key
  })
}
```


## Signicat

Connector ID (`connector.id` in the resource): `connectorSignicat`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `customAuth` (json):  Console display name: "Custom Parameters".


Example:
```terraform
resource "pingone_davinci_connector_instance" "connectorSignicat" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorSignicat"
  }
  name = "My awesome connectorSignicat"
  properties = jsonencode({
    "customAuth" = jsonencode({})
  })
}
```


## Silverfort

Connector ID (`connector.id` in the resource): `silverfortConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `apiKey` (string): Silverfort Risk API Key. Console display name: "Risk API Key".
* `appUserSecret` (string): Silverfort App User Secret. Console display name: "App User Secret".
* `consoleApi` (string): Silverfort App User ID. Console display name: "App User ID".


Example:
```terraform
resource "pingone_davinci_connector_instance" "silverfortConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "silverfortConnector"
  }
  name = "My awesome silverfortConnector"
  properties = jsonencode({
    "apiKey" = var.silverfortconnector_property_api_key
    "appUserSecret" = var.silverfortconnector_property_app_user_secret
    "consoleApi" = var.silverfortconnector_property_console_api
  })
}
```


## Sinch

Connector ID (`connector.id` in the resource): `sinchConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `acceptLanguage` (string): Language of SMS sent, if using Sinch provided templates will be chosen based on Accept-Language header. Examples include, but are not limited to pl-PL, no-NO, en-US. Console display name: "Language".
* `applicationKey` (string): Verification Application Key from your Sinch Account. Console display name: "Sinch Application Key".
* `secretKey` (string): Verification Secret Key from your Sinch Account. Console display name: "Sinch Secret Key".


Example:
```terraform
resource "pingone_davinci_connector_instance" "sinchConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "sinchConnector"
  }
  name = "My awesome sinchConnector"
  properties = jsonencode({
    "acceptLanguage" = var.sinchconnector_property_accept_language
    "applicationKey" = var.sinchconnector_property_application_key
    "secretKey" = var.sinchconnector_property_secret_key
  })
}
```


## Singpass Login

Connector ID (`connector.id` in the resource): `singpassLoginConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `customAuth` (json):  Console display name: "Custom Parameters".


Example:
```terraform
resource "pingone_davinci_connector_instance" "singpassLoginConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "singpassLoginConnector"
  }
  name = "My awesome singpassLoginConnector"
  properties = jsonencode({
    "customAuth" = jsonencode({})
  })
}
```


## Slack Login

Connector ID (`connector.id` in the resource): `slackConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `oauth2` (json):  Console display name: "Oauth2 Parameters".


Example:
```terraform
resource "pingone_davinci_connector_instance" "slackConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "slackConnector"
  }
  name = "My awesome slackConnector"
  properties = jsonencode({
    "oauth2" = jsonencode({})
  })
}
```


## Smarty Address Validator

Connector ID (`connector.id` in the resource): `connectorSmarty`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `authId` (string): Smarty Authentication ID (Found on 'API Keys' tab in Smarty tenant). Console display name: "Auth ID".
* `authToken` (string): Smarty Authentication Token (Found on 'API Keys' tab in Smarty tenant). Console display name: "Auth Token".
* `license` (string): Smarty License Value (Found on 'Subscriptions' tab in Smarty tenant). Console display name: "License".


Example:
```terraform
resource "pingone_davinci_connector_instance" "connectorSmarty" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorSmarty"
  }
  name = "My awesome connectorSmarty"
  properties = jsonencode({
    "authId" = var.connectorsmarty_property_auth_id
    "authToken" = var.connectorsmarty_property_auth_token
    "license" = var.connectorsmarty_property_license
  })
}
```


## Socure

Connector ID (`connector.id` in the resource): `socureConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `apiKey` (string): ID+ Key is the API key that you can retrieve from Socure Admin Portal. Console display name: "ID+ Key".
* `baseUrl` (string): The Socure API URL to target. For a custom value, select Use Custom API URL and enter a value in the Custom API URL field. Console display name: "API URL".
* `customApiUrl` (string): The URL for the Socure API, such as "https://example.socure.com". Console display name: "Custom API URL".
* `sdkKey` (string): SDK Key that you can retrieve from Socure Admin Portal. Console display name: "SDK Key".
* `skWebhookUri` (string): Use this url as the Webhook URL in the Third Party Integration's configuration. Console display name: "Webhook URL".


Example:
```terraform
resource "pingone_davinci_connector_instance" "socureConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "socureConnector"
  }
  name = "My awesome socureConnector"
  properties = jsonencode({
    "apiKey" = var.socureconnector_property_api_key
    "baseUrl" = var.socureconnector_property_base_url
    "customApiUrl" = var.socureconnector_property_custom_api_url
    "sdkKey" = var.socureconnector_property_sdk_key
    "skWebhookUri" = var.socureconnector_property_sk_webhook_uri
  })
}
```


## Splunk

Connector ID (`connector.id` in the resource): `splunkConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `apiUrl` (string): The Base API URL for Splunk. Console display name: "Base URL".
* `port` (number): API Server Port. Console display name: "Port".
* `token` (string): Splunk Token to make API requests. Console display name: "Token".


Example:
```terraform
resource "pingone_davinci_connector_instance" "splunkConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "splunkConnector"
  }
  name = "My awesome splunkConnector"
  properties = jsonencode({
    "apiUrl" = var.splunkconnector_property_api_url
    "port" = var.splunkconnector_property_port
    "token" = var.splunkconnector_property_token
  })
}
```


## Spotify

Connector ID (`connector.id` in the resource): `connectorSpotify`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `oauth2` (json):  Console display name: "Oauth2 Parameters".


Example:
```terraform
resource "pingone_davinci_connector_instance" "connectorSpotify" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorSpotify"
  }
  name = "My awesome connectorSpotify"
  properties = jsonencode({
    "oauth2" = jsonencode({})
  })
}
```


## SpyCloud Enterprise Protection

Connector ID (`connector.id` in the resource): `connectorSpycloud`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `apiKey` (string): Contact SpyCloud to acquire an Employee ATO Prevention API Key that will work with DaVinci. Console display name: "SpyCloud Employee ATO Prevention API Key".


Example:
```terraform
resource "pingone_davinci_connector_instance" "connectorSpycloud" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorSpycloud"
  }
  name = "My awesome connectorSpycloud"
  properties = jsonencode({
    "apiKey" = var.connectorspycloud_property_api_key
  })
}
```


## String

Connector ID (`connector.id` in the resource): `stringsConnector`

*No properties*


Example:
```terraform
resource "pingone_davinci_connector_instance" "stringsConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "stringsConnector"
  }
  name = "My awesome stringsConnector"
}
```


## TMT Analysis

Connector ID (`connector.id` in the resource): `tmtConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `apiKey` (string): API Key for TMT Analysis. Console display name: "API Key".
* `apiSecret` (string): API Secret for TMT Analysis. Console display name: "API Secret".
* `apiUrl` (string): The Base API URL for TMT Analysis. Console display name: "Base URL".


Example:
```terraform
resource "pingone_davinci_connector_instance" "tmtConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "tmtConnector"
  }
  name = "My awesome tmtConnector"
  properties = jsonencode({
    "apiKey" = var.tmtconnector_property_api_key
    "apiSecret" = var.tmtconnector_property_api_secret
    "apiUrl" = var.tmtconnector_property_api_url
  })
}
```


## Tableau

Connector ID (`connector.id` in the resource): `connectorTableau`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `addFlowPermissionsRequestBody` (string): Add Flow Permissions Request Body in XML Format. Example: <tsRequest><task><flowRun><flow id="flow-id"/><flowRunSpec><flowParameterSpecs><flowParameterSpec parameterId="parameter-id" overrideValue= "overrideValue"/><flowParameterSpecs><flowRunSpec></flowRun></task></tsRequest>. Console display name: "Add Flow Permissions Request Body in XML Format.".
* `addUsertoSiteRequestBody` (string): Add User to Site Request Body in XML Format. Example: <tsRequest><user name="user-name" siteRole="site-role" authSetting="auth-setting" /></tsRequest>. Console display name: "Add User to Site Request Body in XML Format.".
* `apiVersion` (string): The version of the API to use, such as 3.16. Console display name: "api-version".
* `authId` (string): The Tableau-Auth sent along with every request. Console display name: "auth-ID".
* `createScheduleBody` (string): This should contain the entire XML. Eg: <tsRequest><schedule name="schedule-name"priority="schedule-priority"type="schedule-type"frequency="schedule-frequency"executionOrder="schedule-execution-order"><frequencyDetails start="start-time" end="end-time"><intervals><interval interval-expression /></intervals></frequencyDetails></schedule></tsRequest>. Console display name: "XML file format to be used for creating schedule".
* `datasourceId` (string): The ID of the flow. Console display name: "datasource-id".
* `flowId` (string): The flow-id value for the flow you want to add permissions to. Console display name: "flow-id".
* `groupId` (string): The ID of the group. Console display name: "group-id".
* `jobId` (string): The ID of the job. Console display name: "job-id".
* `scheduleId` (string): The ID of the schedule that you are associating with the data source. Console display name: "schedule-id".
* `serverUrl` (string): The tableau server URL Example: https://www.tableau.com:8030. Console display name: "server-url".
* `siteId` (string): The ID of the site that contains the view. Console display name: "site-id".
* `taskId` (string): The ID of the extract refresh task. Console display name: "task-id".
* `updateScheduleRequestBody` (string): This should contain the entire XML. Eg: <tsRequest><schedule name="hourly-schedule-1" priority="50" type="Extract" frequency="Hourly" executionOrder="Parallel"><frequencyDetails start="18:30:00" end="23:00:00"><intervals><interval hours="2" /></intervals></frequencyDetails></schedule></tsRequest>. Console display name: "XML file format to be used for updating schedule".
* `updateUserRequestBody` (string): Update User Request Body in XML Format. <tsRequest><user fullName="new-full-name" email="new-email" password="new-password" siteRole="new-site-role" authSetting="new-auth-setting" /></tsRequest>. Console display name: "Update User Request Body in XML Format.".
* `userId` (string): The ID of the user to get/give information for. Console display name: "user-id".
* `workbookId` (string): The ID of the workbook to add to the schedule. Console display name: "workbook-id".


Example:
```terraform
resource "pingone_davinci_connector_instance" "connectorTableau" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorTableau"
  }
  name = "My awesome connectorTableau"
  properties = jsonencode({
    "addFlowPermissionsRequestBody" = var.connectortableau_property_add_flow_permissions_request_body
    "addUsertoSiteRequestBody" = var.connectortableau_property_add_userto_site_request_body
    "apiVersion" = var.connectortableau_property_api_version
    "authId" = var.connectortableau_property_auth_id
    "createScheduleBody" = var.connectortableau_property_create_schedule_body
    "datasourceId" = var.connectortableau_property_datasource_id
    "flowId" = var.connectortableau_property_flow_id
    "groupId" = var.connectortableau_property_group_id
    "jobId" = var.connectortableau_property_job_id
    "scheduleId" = var.connectortableau_property_schedule_id
    "serverUrl" = var.connectortableau_property_server_url
    "siteId" = var.connectortableau_property_site_id
    "taskId" = var.connectortableau_property_task_id
    "updateScheduleRequestBody" = var.connectortableau_property_update_schedule_request_body
    "updateUserRequestBody" = var.connectortableau_property_update_user_request_body
    "userId" = var.connectortableau_property_user_id
    "workbookId" = var.connectortableau_property_workbook_id
  })
}
```


## Talend Identities Management API

Connector ID (`connector.id` in the resource): `connector-oai-talendim`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `authBearerToken` (string): The authenticating token. Console display name: "Bearer Token".
* `basePath` (string): The base URL for contacting the API. Console display name: "Base Path".


Example:
```terraform
resource "pingone_davinci_connector_instance" "connector-oai-talendim" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connector-oai-talendim"
  }
  name = "My awesome connector-oai-talendim"
  properties = jsonencode({
    "authBearerToken" = var.connector-oai-talendim_property_auth_bearer_token
    "basePath" = var.connector-oai-talendim_property_base_path
  })
}
```


## Talend SCIM API

Connector ID (`connector.id` in the resource): `connector-oai-talendscim`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `authBearerToken` (string): The authenticating token. Console display name: "Bearer Token".
* `basePath` (string): The base URL for contacting the API. Console display name: "Base Path".


Example:
```terraform
resource "pingone_davinci_connector_instance" "connector-oai-talendscim" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connector-oai-talendscim"
  }
  name = "My awesome connector-oai-talendscim"
  properties = jsonencode({
    "authBearerToken" = var.connector-oai-talendscim_property_auth_bearer_token
    "basePath" = var.connector-oai-talendscim_property_base_path
  })
}
```


## Teleport

Connector ID (`connector.id` in the resource): `nodeConnector`

*No properties*


Example:
```terraform
resource "pingone_davinci_connector_instance" "nodeConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "nodeConnector"
  }
  name = "My awesome nodeConnector"
}
```


## Telesign

Connector ID (`connector.id` in the resource): `telesignConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `authDescription` (string):  Console display name: "Authentication Description".
* `connectorName` (string):  Console display name: "Connector Name".
* `description` (string):  Console display name: "Description".
* `details1` (string):  Console display name: "Credentials Details 1".
* `details2` (string):  Console display name: "Credentials Details 2".
* `iconUrl` (string):  Console display name: "Icon URL".
* `iconUrlPng` (string):  Console display name: "Icon URL in PNG".
* `password` (string):  Console display name: "Password".
* `providerName` (string):  Console display name: "Provider Name".
* `showCredAddedOn` (boolean):  Console display name: "Show Credentials Added On?".
* `showCredAddedVia` (boolean):  Console display name: "Show Credentials Added through ?".
* `title` (string):  Console display name: "Title".
* `toolTip` (string):  Console display name: "Tooltip".
* `username` (string):  Console display name: "Username".


Example:
```terraform
resource "pingone_davinci_connector_instance" "telesignConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "telesignConnector"
  }
  name = "My awesome telesignConnector"
  properties = jsonencode({
    "authDescription" = var.telesignconnector_property_auth_description
    "connectorName" = var.telesignconnector_property_connector_name
    "description" = var.telesignconnector_property_description
    "details1" = var.telesignconnector_property_details1
    "details2" = var.telesignconnector_property_details2
    "iconUrl" = var.telesignconnector_property_icon_url
    "iconUrlPng" = var.telesignconnector_property_icon_url_png
    "password" = var.telesignconnector_property_password
    "providerName" = var.telesignconnector_property_provider_name
    "showCredAddedOn" = var.telesignconnector_property_show_cred_added_on
    "showCredAddedVia" = var.telesignconnector_property_show_cred_added_via
    "title" = var.telesignconnector_property_title
    "toolTip" = var.telesignconnector_property_tool_tip
    "username" = var.telesignconnector_property_username
  })
}
```


## Token Management

Connector ID (`connector.id` in the resource): `skOpenIdConnector`

*No properties*


Example:
```terraform
resource "pingone_davinci_connector_instance" "skOpenIdConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "skOpenIdConnector"
  }
  name = "My awesome skOpenIdConnector"
}
```


## TransUnion TLOxp

Connector ID (`connector.id` in the resource): `tutloxpConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `apiUrl` (string): The URL for your TransUnion API. Unnecessary to change unless you're testing against a demo tenant. Console display name: "API URL".
* `dppaCode` (string): The DPPA code that determines the level of data access in the API. Console display name: "DPPA Purpose Code".
* `glbCode` (string): The GLB code that determines the level of data access in the API. Console display name: "GLB Purpose Code".
* `password` (string): The password for your API User. Console display name: "Password".
* `username` (string): The username for your API user. Console display name: "Username".


Example:
```terraform
resource "pingone_davinci_connector_instance" "tutloxpConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "tutloxpConnector"
  }
  name = "My awesome tutloxpConnector"
  properties = jsonencode({
    "apiUrl" = var.tutloxpconnector_property_api_url
    "dppaCode" = var.tutloxpconnector_property_dppa_code
    "glbCode" = var.tutloxpconnector_property_glb_code
    "password" = var.tutloxpconnector_property_password
    "username" = var.tutloxpconnector_property_username
  })
}
```


## TransUnion TruValidate

Connector ID (`connector.id` in the resource): `transunionConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `apiUrl` (string): The Base API URL for TransUnion. Console display name: "Base URL".
* `docVerificationPassword` (string): Password for Document Verification, provided by TransUnion. Console display name: "Password".
* `docVerificationPublicKey` (string): Public Key for Document Verification, provided by TransUnion. Console display name: "Public Key".
* `docVerificationSecret` (string): Secret for Document Verification, provided by TransUnion. Console display name: "Secret".
* `docVerificationSiteId` (string): Site ID for Document Verification, provided by TransUnion. Console display name: "Site ID".
* `docVerificationUsername` (string): Username for Document Verification, provided by TransUnion. Console display name: "Username".
* `idVerificationPassword` (string): Password for ID Verification, provided by TransUnion. Console display name: "Password".
* `idVerificationPublicKey` (string): Public Key for ID Verification, provided by TransUnion. Console display name: "Public Key".
* `idVerificationSecret` (string): Secret for ID Verification, provided by TransUnion. Console display name: "Secret".
* `idVerificationSiteId` (string): Site ID for ID Verification, provided by TransUnion. Console display name: "Site ID".
* `idVerificationUsername` (string): Username for ID Verification, provided by TransUnion. Console display name: "Username".
* `kbaPassword` (string): Password for KBA, provided by TransUnion. Console display name: "Password".
* `kbaPublicKey` (string): Public Key for KBA, provided by TransUnion. Console display name: "Public Key".
* `kbaSecret` (string): Secret for KBA, provided by TransUnion. Console display name: "Secret".
* `kbaSiteId` (string): Site ID for KBA, provided by TransUnion. Console display name: "Site ID".
* `kbaUsername` (string): Username for KBA, provided by TransUnion. Console display name: "Username".
* `otpPassword` (string): Password for otp Verification, provided by TransUnion. Console display name: "Password".
* `otpPublicKey` (string): Public Key for otp Verification, provided by TransUnion. Console display name: "Public Key".
* `otpSecret` (string): Secret for otp Verification, provided by TransUnion. Console display name: "Secret".
* `otpSiteId` (string): Site ID for otp Verification, provided by TransUnion. Console display name: "Site ID".
* `otpUsername` (string): Username for otp Verification, provided by TransUnion. Console display name: "Username".


Example:
```terraform
resource "pingone_davinci_connector_instance" "transunionConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "transunionConnector"
  }
  name = "My awesome transunionConnector"
  properties = jsonencode({
    "apiUrl" = var.transunionconnector_property_api_url
    "docVerificationPassword" = var.transunionconnector_property_doc_verification_password
    "docVerificationPublicKey" = var.transunionconnector_property_doc_verification_public_key
    "docVerificationSecret" = var.transunionconnector_property_doc_verification_secret
    "docVerificationSiteId" = var.transunionconnector_property_doc_verification_site_id
    "docVerificationUsername" = var.transunionconnector_property_doc_verification_username
    "idVerificationPassword" = var.transunionconnector_property_id_verification_password
    "idVerificationPublicKey" = var.transunionconnector_property_id_verification_public_key
    "idVerificationSecret" = var.transunionconnector_property_id_verification_secret
    "idVerificationSiteId" = var.transunionconnector_property_id_verification_site_id
    "idVerificationUsername" = var.transunionconnector_property_id_verification_username
    "kbaPassword" = var.transunionconnector_property_kba_password
    "kbaPublicKey" = var.transunionconnector_property_kba_public_key
    "kbaSecret" = var.transunionconnector_property_kba_secret
    "kbaSiteId" = var.transunionconnector_property_kba_site_id
    "kbaUsername" = var.transunionconnector_property_kba_username
    "otpPassword" = var.transunionconnector_property_otp_password
    "otpPublicKey" = var.transunionconnector_property_otp_public_key
    "otpSecret" = var.transunionconnector_property_otp_secret
    "otpSiteId" = var.transunionconnector_property_otp_site_id
    "otpUsername" = var.transunionconnector_property_otp_username
  })
}
```


## Treasure Data

Connector ID (`connector.id` in the resource): `treasureDataConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `apiKey` (string): Treasure Data API Key. Console display name: "API Key".


Example:
```terraform
resource "pingone_davinci_connector_instance" "treasureDataConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "treasureDataConnector"
  }
  name = "My awesome treasureDataConnector"
  properties = jsonencode({
    "apiKey" = var.treasuredataconnector_property_api_key
  })
}
```


## Trulioo

Connector ID (`connector.id` in the resource): `connectorTrulioo`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `clientID` (string): Trulioo Client ID. Console display name: "Client ID".
* `clientSecret` (string): Trulioo Client Secret. Console display name: "Client Secret".


Example:
```terraform
resource "pingone_davinci_connector_instance" "connectorTrulioo" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorTrulioo"
  }
  name = "My awesome connectorTrulioo"
  properties = jsonencode({
    "clientID" = var.connectortrulioo_property_client_i_d
    "clientSecret" = var.connectortrulioo_property_client_secret
  })
}
```


## Twilio

Connector ID (`connector.id` in the resource): `twilioConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `accountSid` (string):  Console display name: "Account Sid".
* `authDescription` (string):  Console display name: "Authentication Description".
* `authMessageTemplate` (string):  Console display name: "Text Message Template (Authentication)".
* `authToken` (string):  Console display name: "Auth Token".
* `connectorName` (string):  Console display name: "Connector Name".
* `description` (string):  Console display name: "Description".
* `details1` (string):  Console display name: "Credentials Details 1".
* `details2` (string):  Console display name: "Credentials Details 2".
* `iconUrl` (string):  Console display name: "Icon URL".
* `iconUrlPng` (string):  Console display name: "Icon URL in PNG".
* `registerMessageTemplate` (string):  Console display name: "Text Message Template (Registration)".
* `senderPhoneNumber` (string):  Console display name: "Sender Phone Number".
* `showCredAddedOn` (boolean):  Console display name: "Show Credentials Added On?".
* `showCredAddedVia` (boolean):  Console display name: "Show Credentials Added through ?".
* `title` (string):  Console display name: "Title".
* `toolTip` (string):  Console display name: "Tooltip".


Example:
```terraform
resource "pingone_davinci_connector_instance" "twilioConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "twilioConnector"
  }
  name = "My awesome twilioConnector"
  properties = jsonencode({
    "accountSid" = var.twilioconnector_property_account_sid
    "authDescription" = var.twilioconnector_property_auth_description
    "authMessageTemplate" = var.twilioconnector_property_auth_message_template
    "authToken" = var.twilioconnector_property_auth_token
    "connectorName" = var.twilioconnector_property_connector_name
    "description" = var.twilioconnector_property_description
    "details1" = var.twilioconnector_property_details1
    "details2" = var.twilioconnector_property_details2
    "iconUrl" = var.twilioconnector_property_icon_url
    "iconUrlPng" = var.twilioconnector_property_icon_url_png
    "registerMessageTemplate" = var.twilioconnector_property_register_message_template
    "senderPhoneNumber" = var.twilioconnector_property_sender_phone_number
    "showCredAddedOn" = var.twilioconnector_property_show_cred_added_on
    "showCredAddedVia" = var.twilioconnector_property_show_cred_added_via
    "title" = var.twilioconnector_property_title
    "toolTip" = var.twilioconnector_property_tool_tip
  })
}
```


## TypingDNA

Connector ID (`connector.id` in the resource): `typingdnaConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `customAuth` (json):  Console display name: "Custom Parameters".


Example:
```terraform
resource "pingone_davinci_connector_instance" "typingdnaConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "typingdnaConnector"
  }
  name = "My awesome typingdnaConnector"
  properties = jsonencode({
    "customAuth" = var.typingdnaconnector_property_custom_auth
  })
}
```


## UnifyID

Connector ID (`connector.id` in the resource): `unifyIdConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `accountId` (string):  Console display name: "Account ID".
* `apiKey` (string):  Console display name: "API Key".
* `connectorName` (string):  Console display name: "Connector Name".
* `details1` (string):  Console display name: "Credentials Details 1".
* `details2` (string):  Console display name: "Credentials Details 2".
* `iconUrl` (string):  Console display name: "Icon URL".
* `iconUrlPng` (string):  Console display name: "Icon URL in PNG".
* `sdkToken` (string):  Console display name: "SDK Token".
* `showCredAddedOn` (boolean):  Console display name: "Show Credentials Added On?".
* `showCredAddedVia` (boolean):  Console display name: "Show Credentials Added through ?".
* `toolTip` (string):  Console display name: "Tooltip".


Example:
```terraform
resource "pingone_davinci_connector_instance" "unifyIdConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "unifyIdConnector"
  }
  name = "My awesome unifyIdConnector"
  properties = jsonencode({
    "accountId" = var.unifyidconnector_property_account_id
    "apiKey" = var.unifyidconnector_property_api_key
    "connectorName" = var.unifyidconnector_property_connector_name
    "details1" = var.unifyidconnector_property_details1
    "details2" = var.unifyidconnector_property_details2
    "iconUrl" = var.unifyidconnector_property_icon_url
    "iconUrlPng" = var.unifyidconnector_property_icon_url_png
    "sdkToken" = var.unifyidconnector_property_sdk_token
    "showCredAddedOn" = var.unifyidconnector_property_show_cred_added_on
    "showCredAddedVia" = var.unifyidconnector_property_show_cred_added_via
    "toolTip" = var.unifyidconnector_property_tool_tip
  })
}
```


## User Policy

Connector ID (`connector.id` in the resource): `userPolicyConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `passwordExpiryInDays` (number): Choose 0 for never expire. Console display name: "Expires in the specified number of days".
* `passwordExpiryNotification` (boolean):  Console display name: "Notify user before password expires".
* `passwordLengthMax` (number):  Console display name: "Maximum Password Length".
* `passwordLengthMin` (number):  Console display name: "Minimum Password Length".
* `passwordLockoutAttempts` (number):  Console display name: "Number of failed login attempts before account is locked".
* `passwordPreviousXPasswords` (number): Choose 0 if any previous passwords are allowed. This is not recommended. Console display name: "Number of unique user passwords associated with a user".
* `passwordRequireLowercase` (boolean): Should the password contain lowercase characters?. Console display name: "Require Lowercase Characters".
* `passwordRequireNumbers` (boolean): Should the password contain numbers?. Console display name: "Require Numbers".
* `passwordRequireSpecial` (boolean): Should the password contain special character?. Console display name: "Require Special Characters".
* `passwordRequireUppercase` (boolean): Should the password contain uppercase characters?. Console display name: "Require Uppercase Characters".
* `passwordSpacesOk` (boolean): Are spaces allowed in the password?. Console display name: "Spaces Accepted".
* `passwordsEnabled` (boolean):  Console display name: "Passwords Feature Enabled?".
* `temporaryPasswordExpiryInDays` (number): If an administrator sets a temporary password, choose how long before it expires. Console display name: "Temporary password expires in the specified number of days".


Example:
```terraform
resource "pingone_davinci_connector_instance" "userPolicyConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "userPolicyConnector"
  }
  name = "My awesome userPolicyConnector"
  properties = jsonencode({
    "passwordExpiryInDays" = var.userpolicyconnector_property_password_expiry_in_days
    "passwordExpiryNotification" = var.userpolicyconnector_property_password_expiry_notification
    "passwordLengthMax" = var.userpolicyconnector_property_password_length_max
    "passwordLengthMin" = var.userpolicyconnector_property_password_length_min
    "passwordLockoutAttempts" = var.userpolicyconnector_property_password_lockout_attempts
    "passwordPreviousXPasswords" = var.userpolicyconnector_property_password_previous_x_passwords
    "passwordRequireLowercase" = var.userpolicyconnector_property_password_require_lowercase
    "passwordRequireNumbers" = var.userpolicyconnector_property_password_require_numbers
    "passwordRequireSpecial" = var.userpolicyconnector_property_password_require_special
    "passwordRequireUppercase" = var.userpolicyconnector_property_password_require_uppercase
    "passwordSpacesOk" = var.userpolicyconnector_property_password_spaces_ok
    "passwordsEnabled" = var.userpolicyconnector_property_passwords_enabled
    "temporaryPasswordExpiryInDays" = var.userpolicyconnector_property_temporary_password_expiry_in_days
  })
}
```


## User Pool

Connector ID (`connector.id` in the resource): `skUserPool`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `customAttributes` (json):  


Example:
```terraform
resource "pingone_davinci_connector_instance" "skUserPool" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "skUserPool"
  }
  name = "My awesome skUserPool"
  properties = jsonencode({
    "customAttributes" = jsonencode({
				"type" : "array",
				"preferredControlType" : "tableViewAttributes",
				"sections" : [
				  "connectorAttributes"
				],
				"value" : [
				  {
					"name" : "username",
					"description" : "Username",
					"type" : "string",
					"value" : null,
					"minLength" : "1",
					"maxLength" : "300",
					"required" : true,
					"attributeType" : "sk"
				  },
				  {
					"name" : "firstName",
					"description" : "First Name",
					"type" : "string",
					"value" : null,
					"minLength" : "1",
					"maxLength" : "100",
					"required" : false,
					"attributeType" : "sk"
				  },
				  {
					"name" : "lastName",
					"description" : "Last Name",
					"type" : "string",
					"value" : null,
					"minLength" : "1",
					"maxLength" : "100",
					"required" : false,
					"attributeType" : "sk"
				  },
				  {
					"name" : "name",
					"description" : "Display Name",
					"type" : "string",
					"value" : null,
					"minLength" : "1",
					"maxLength" : "250",
					"required" : false,
					"attributeType" : "sk"
				  },
				  {
					"name" : "email",
					"description" : "Email",
					"type" : "string",
					"value" : null,
					"minLength" : "1",
					"maxLength" : "250",
					"required" : false,
					"attributeType" : "sk"
					}
				]
			})
  })
}
```


## ValidSoft

Connector ID (`connector.id` in the resource): `connectorValidsoft`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `customAuth` (json):  Console display name: "Custom Parameters".


Example:
```terraform
resource "pingone_davinci_connector_instance" "connectorValidsoft" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorValidsoft"
  }
  name = "My awesome connectorValidsoft"
  properties = jsonencode({
    "customAuth" = jsonencode({})
  })
}
```


## Variable

Connector ID (`connector.id` in the resource): `variablesConnector`

*No properties*


Example:
```terraform
resource "pingone_davinci_connector_instance" "variablesConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "variablesConnector"
  }
  name = "My awesome variablesConnector"
}
```


## Venafi Account Service API

Connector ID (`connector.id` in the resource): `connector-oai-venafi`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `authApiKey` (string): The authentication key to the Venafi as a Service API for Account Service Operations. Console display name: "API Key".
* `basePath` (string): The base URL for contacting the API. Console display name: "Base Path".


Example:
```terraform
resource "pingone_davinci_connector_instance" "connector-oai-venafi" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connector-oai-venafi"
  }
  name = "My awesome connector-oai-venafi"
  properties = jsonencode({
    "authApiKey" = var.connector-oai-venafi_property_auth_api_key
    "basePath" = var.connector-oai-venafi_property_base_path
  })
}
```


## Vericlouds

Connector ID (`connector.id` in the resource): `connectorVericlouds`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `apiSecret` (string): The API secret assigned by VeriClouds to the customer. The secret is also used for decrypting sensitive data such as leaked passwords. It is important to never share the secret with any 3rd party. Console display name: "apiSecret".
* `apikey` (string): The API key assigned by VeriClouds to the customer. Console display name: "apiKey".


Example:
```terraform
resource "pingone_davinci_connector_instance" "connectorVericlouds" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorVericlouds"
  }
  name = "My awesome connectorVericlouds"
  properties = jsonencode({
    "apiSecret" = var.connectorvericlouds_property_api_secret
    "apikey" = var.connectorvericlouds_property_apikey
  })
}
```


## Veriff

Connector ID (`connector.id` in the resource): `veriffConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `access_token` (string): The API Key provided by Veriff, such as "323aa031-b4af-4e12-b354-de0da91a2ab0". Console display name: "API Key".
* `baseUrl` (string): The API URL to target, such as “https://stationapi.veriff.com/”. Console display name: "Base URL".
* `password` (string): The Share Secret Key from Veriff to create HMAC signature, such as "20bf4sf0-fbg7-488c-b4f1-d9594lf707bk". Console display name: "Shared Secret Key".


Example:
```terraform
resource "pingone_davinci_connector_instance" "veriffConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "veriffConnector"
  }
  name = "My awesome veriffConnector"
  properties = jsonencode({
    "access_token" = var.veriffconnector_property_access_token
    "baseUrl" = var.veriffconnector_property_base_url
    "password" = var.veriffconnector_property_password
  })
}
```


## Verosint

Connector ID (`connector.id` in the resource): `connector443id`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `apiKey` (string): This is the API key from your Verosint account. Remember, Your API KEY is like a serial number for your policy. If you want to utilize more than one policy, you can generate another API KEY and tailor that to a custom policy. Console display name: "API Key".


Example:
```terraform
resource "pingone_davinci_connector_instance" "connector443id" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connector443id"
  }
  name = "My awesome connector443id"
  properties = jsonencode({
    "apiKey" = var.connector443id_property_api_key
  })
}
```


## Vidos

Connector ID (`connector.id` in the resource): `mailchainConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `apiKey` (string): Enter your Vidos API Key obtained from the Vidos Dashboard with appropriate resolver or verifier permissions (visit https://dashboard.vidos.id/iam/api-keys). Console display name: "Vidos API Key".
* `version` (string): The verification API specification version. Console display name: "Verifier Version".


Example:
```terraform
resource "pingone_davinci_connector_instance" "mailchainConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "mailchainConnector"
  }
  name = "My awesome mailchainConnector"
  properties = jsonencode({
    "apiKey" = var.mailchainconnector_property_api_key
    "version" = var.mailchainconnector_property_version
  })
}
```


## Webhook

Connector ID (`connector.id` in the resource): `webhookConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `urls` (string): POST requests will be made to these registered url as selected later. Console display name: "Register URLs".


Example:
```terraform
resource "pingone_davinci_connector_instance" "webhookConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "webhookConnector"
  }
  name = "My awesome webhookConnector"
  properties = jsonencode({
    "urls" = var.webhookconnector_property_urls
  })
}
```


## WhatsApp for Business

Connector ID (`connector.id` in the resource): `connectorWhatsAppBusiness`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `accessToken` (string): WhatsApp Access Token. Console display name: "Access Token".
* `appSecret` (string): WhatsApp App Secret for the application, it is used to verify the webhook signatures. Console display name: "App Secret".
* `skWebhookUri` (string): Use this url as the Webhook URL in the Third Party Integration's configuration. Console display name: "Redirect Webhook URI".
* `verifyToken` (string): Meta webhook verify token. Console display name: "Webhook Verify Token".
* `version` (string): WhatsApp Graph API Version. Console display name: "Version".


Example:
```terraform
resource "pingone_davinci_connector_instance" "connectorWhatsAppBusiness" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorWhatsAppBusiness"
  }
  name = "My awesome connectorWhatsAppBusiness"
  properties = jsonencode({
    "accessToken" = var.connectorwhatsappbusiness_property_access_token
    "appSecret" = var.connectorwhatsappbusiness_property_app_secret
    "skWebhookUri" = var.connectorwhatsappbusiness_property_sk_webhook_uri
    "verifyToken" = var.connectorwhatsappbusiness_property_verify_token
    "version" = var.connectorwhatsappbusiness_property_version
  })
}
```


## WinMagic

Connector ID (`connector.id` in the resource): `connectorWinmagic`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `openId` (json):  Console display name: "OpenId Parameters".


Example:
```terraform
resource "pingone_davinci_connector_instance" "connectorWinmagic" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorWinmagic"
  }
  name = "My awesome connectorWinmagic"
  properties = jsonencode({
    "openId" = jsonencode({})
  })
}
```


## WireWheel

Connector ID (`connector.id` in the resource): `wireWheelConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `baseURL` (string): The base API URL of the WireWheel environment. Console display name: "WireWheel Base API URL".
* `clientId` (string): Client ID from WireWheel Channel settings. Console display name: "Client ID".
* `clientSecret` (string): Client Secret from WireWheel Channel settings. Console display name: "Client Secret".
* `issuerId` (string): Issuer URL from WireWheel Channel settings. Console display name: "Issuer URL".


Example:
```terraform
resource "pingone_davinci_connector_instance" "wireWheelConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "wireWheelConnector"
  }
  name = "My awesome wireWheelConnector"
  properties = jsonencode({
    "baseURL" = var.base_url
    "clientId" = var.wirewheelconnector_property_client_id
    "clientSecret" = var.wirewheelconnector_property_client_secret
    "issuerId" = var.wirewheelconnector_property_issuer_id
  })
}
```


## X Login

Connector ID (`connector.id` in the resource): `twitterIdpConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `customAuth` (json):  Console display name: "Custom Parameters".


Example:
```terraform
resource "pingone_davinci_connector_instance" "twitterIdpConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "twitterIdpConnector"
  }
  name = "My awesome twitterIdpConnector"
  properties = jsonencode({
    "customAuth" = jsonencode({})
  })
}
```


## Yoti

Connector ID (`connector.id` in the resource): `yotiConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `customAuth` (json):  Console display name: "Custom Parameters".


Example:
```terraform
resource "pingone_davinci_connector_instance" "yotiConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "yotiConnector"
  }
  name = "My awesome yotiConnector"
  properties = jsonencode({
    "customAuth" = jsonencode({})
  })
}
```


## Zendesk

Connector ID (`connector.id` in the resource): `connectorZendesk`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `apiToken` (string): An Active Zendesk API Token (admin center->Apps&Integrations->Zendesk API). Console display name: "Zendesk API Token".
* `emailUsername` (string): Email used as 'username' for your Zendesk account. Console display name: "Email of User (username)".
* `subdomain` (string): Your Zendesk subdomain (ex. {subdomain}.zendesk.com/api/v2/...). Console display name: "Subdomain".


Example:
```terraform
resource "pingone_davinci_connector_instance" "connectorZendesk" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorZendesk"
  }
  name = "My awesome connectorZendesk"
  properties = jsonencode({
    "apiToken" = var.connectorzendesk_property_api_token
    "emailUsername" = var.connectorzendesk_property_email_username
    "subdomain" = var.connectorzendesk_property_subdomain
  })
}
```


## Zoop.one

Connector ID (`connector.id` in the resource): `zoopConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `agencyId` (string):  Console display name: "Zoop Agency ID".
* `apiKey` (string):  Console display name: "Zoop API Key".
* `apiUrl` (string):  Console display name: "Zoop API URL".


Example:
```terraform
resource "pingone_davinci_connector_instance" "zoopConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "zoopConnector"
  }
  name = "My awesome zoopConnector"
  properties = jsonencode({
    "agencyId" = var.zoopconnector_property_agency_id
    "apiKey" = var.zoopconnector_property_api_key
    "apiUrl" = var.zoopconnector_property_api_url
  })
}
```


## Zscaler ZIA

Connector ID (`connector.id` in the resource): `connectorZscaler`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `basePath` (string): basePath. Console display name: "Base Path".
* `baseURL` (string): baseURL. Console display name: "Base URL".
* `zscalerAPIkey` (string): Zscaler APIkey. Console display name: "Zscaler APIkey".
* `zscalerPassword` (string): Zscaler Domain Password. Console display name: "Zscaler Password".
* `zscalerUsername` (string): Zscaler Domain Username. Console display name: "Zscaler Username".


Example:
```terraform
resource "pingone_davinci_connector_instance" "connectorZscaler" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorZscaler"
  }
  name = "My awesome connectorZscaler"
  properties = jsonencode({
    "basePath" = var.connectorzscaler_property_base_path
    "baseURL" = var.base_url
    "zscalerAPIkey" = var.zscaler_api_key
    "zscalerPassword" = var.connectorzscaler_property_zscaler_password
    "zscalerUsername" = var.connectorzscaler_property_zscaler_username
  })
}
```


## iProov

Connector ID (`connector.id` in the resource): `iproovConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `allowLandscape` (boolean):  Console display name: "Allow Landscape".
* `apiKey` (string):  Console display name: "API Key".
* `authDescription` (string):  Console display name: "Authentication Description".
* `baseUrl` (string):  Console display name: "Base URL".
* `color1` (string): Ex. #000000. Console display name: "Loading Tint Color".
* `color2` (string): Ex. #000000. Console display name: "Not Ready Tint Color".
* `color3` (string): Ex. #000000. Console display name: "Ready Tint Color".
* `color4` (string): Ex. #000000. Console display name: "Liveness Tint Color".
* `connectorName` (string):  Console display name: "Connector Name".
* `customTitle` (string): Specify a custom title to be shown. Defaults to show an iProov-generated message. Set to empty string "" to hide the message entirely.  Console display name: "Custom Title".
* `description` (string):  Console display name: "Description".
* `details1` (string):  Console display name: "Credentials Details 1".
* `details2` (string):  Console display name: "Credentials Details 2".
* `enableCameraSelector` (boolean):  Console display name: "Enable Camera Selector".
* `iconUrl` (string):  Console display name: "Icon URL".
* `iconUrlPng` (string):  Console display name: "Icon URL in PNG".
* `javascriptCSSUrl` (string):  Console display name: "CSS URL".
* `javascriptCdnUrl` (string):  Console display name: "Javascript CDN URL".
* `kioskMode` (boolean):  Console display name: "Kiosk Mode".
* `logo` (string): You can use a custom logo by simply passing a relative link, absolute path or data URI to your logo. If you do not want a logo to show pass the logo attribute as null. Console display name: "Logo".
* `password` (string):  Console display name: "Password".
* `secret` (string):  Console display name: "Secret".
* `showCountdown` (boolean):  Console display name: "Show Countdown".
* `showCredAddedOn` (boolean):  Console display name: "Show Credentials Added On?".
* `showCredAddedVia` (boolean):  Console display name: "Show Credentials Added through ?".
* `startScreenTitle` (string):  Console display name: "Start Screen Title".
* `title` (string):  Console display name: "Title".
* `toolTip` (string):  Console display name: "Tooltip".
* `username` (string):  Console display name: "Username".


Example:
```terraform
resource "pingone_davinci_connector_instance" "iproovConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "iproovConnector"
  }
  name = "My awesome iproovConnector"
  properties = jsonencode({
    "allowLandscape" = var.iproovconnector_property_allow_landscape
    "apiKey" = var.iproovconnector_property_api_key
    "authDescription" = var.iproovconnector_property_auth_description
    "baseUrl" = var.iproovconnector_property_base_url
    "color1" = var.iproovconnector_property_color1
    "color2" = var.iproovconnector_property_color2
    "color3" = var.iproovconnector_property_color3
    "color4" = var.iproovconnector_property_color4
    "connectorName" = var.iproovconnector_property_connector_name
    "customTitle" = var.iproovconnector_property_custom_title
    "description" = var.iproovconnector_property_description
    "details1" = var.iproovconnector_property_details1
    "details2" = var.iproovconnector_property_details2
    "enableCameraSelector" = var.iproovconnector_property_enable_camera_selector
    "iconUrl" = var.iproovconnector_property_icon_url
    "iconUrlPng" = var.iproovconnector_property_icon_url_png
    "javascriptCSSUrl" = var.javascript_css_url
    "javascriptCdnUrl" = var.iproovconnector_property_javascript_cdn_url
    "kioskMode" = var.iproovconnector_property_kiosk_mode
    "logo" = var.iproovconnector_property_logo
    "password" = var.iproovconnector_property_password
    "secret" = var.iproovconnector_property_secret
    "showCountdown" = var.iproovconnector_property_show_countdown
    "showCredAddedOn" = var.iproovconnector_property_show_cred_added_on
    "showCredAddedVia" = var.iproovconnector_property_show_cred_added_via
    "startScreenTitle" = var.iproovconnector_property_start_screen_title
    "title" = var.iproovconnector_property_title
    "toolTip" = var.iproovconnector_property_tool_tip
    "username" = var.iproovconnector_property_username
  })
}
```


## iProov API

Connector ID (`connector.id` in the resource): `iproovV2Connector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `apiKey` (string): Your iProov Service Provider API key. This can be obtained from your iPortal account. Please contact support@iproov.com for more information. Console display name: "iProov API Key".
* `secret` (string): Your iProov Service Provider Secret. This can be obtained from your iPortal account. Please contact support@iproov.com for more information. Console display name: "iProov Secret".
* `tenant` (string): The iProov tenant URL (do not include https://). This can be obtained from your iPortal account. Please contact support@iproov.com for more information. Console display name: "iProov Tenant".


Example:
```terraform
resource "pingone_davinci_connector_instance" "iproovV2Connector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "iproovV2Connector"
  }
  name = "My awesome iproovV2Connector"
  properties = jsonencode({
    "apiKey" = var.iproovv2connector_property_api_key
    "secret" = var.iproovv2connector_property_secret
    "tenant" = var.iproovv2connector_property_tenant
  })
}
```


## iProov OIDC

Connector ID (`connector.id` in the resource): `connectorSvipe`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `customAuth` (json):  Console display name: "Custom Parameters".


Example:
```terraform
resource "pingone_davinci_connector_instance" "connectorSvipe" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorSvipe"
  }
  name = "My awesome connectorSvipe"
  properties = jsonencode({
    "customAuth" = jsonencode({})
  })
}
```


## iovation

Connector ID (`connector.id` in the resource): `iovationConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `apiUrl` (string):  Console display name: "API Server URL".
* `javascriptCdnUrl` (string): iovation loader javascript CDN. Console display name: "iovation loader Javascript CDN URL".
* `subKey` (string): This will be an iovation assigned value that tracks requests from your site. This is primarily used for debugging and troubleshooting purposes. Console display name: "Sub Key".
* `subscriberAccount` (string):  Console display name: "Subscriber Account".
* `subscriberId` (string):  Console display name: "Subscriber ID".
* `subscriberPasscode` (string):  Console display name: "Subscriber Passcode".
* `version` (string): This is the version of the script to load. The value should either correspond to a specific version you wish to use, or one of the following aliases to get the latest version of the code: general5 - the latest stable version of the javascript, early5 - the latest available version of the javascript. Console display name: "Version".


Example:
```terraform
resource "pingone_davinci_connector_instance" "iovationConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "iovationConnector"
  }
  name = "My awesome iovationConnector"
  properties = jsonencode({
    "apiUrl" = var.iovationconnector_property_api_url
    "javascriptCdnUrl" = var.iovationconnector_property_javascript_cdn_url
    "subKey" = var.iovationconnector_property_sub_key
    "subscriberAccount" = var.iovationconnector_property_subscriber_account
    "subscriberId" = var.iovationconnector_property_subscriber_id
    "subscriberPasscode" = var.iovationconnector_property_subscriber_passcode
    "version" = var.iovationconnector_property_version
  })
}
```


## ipgeolocation.io

Connector ID (`connector.id` in the resource): `connectorIPGeolocationio`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `apiKey` (string): Developer subscription API key. Console display name: "API key".


Example:
```terraform
resource "pingone_davinci_connector_instance" "connectorIPGeolocationio" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorIPGeolocationio"
  }
  name = "My awesome connectorIPGeolocationio"
  properties = jsonencode({
    "apiKey" = var.connectoripgeolocationio_property_api_key
  })
}
```


## ipregistry

Connector ID (`connector.id` in the resource): `connectorIPregistry`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `apiKey` (string): API Key used to authenticate to the ipregistry.co API. Console display name: "API Key".


Example:
```terraform
resource "pingone_davinci_connector_instance" "connectorIPregistry" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorIPregistry"
  }
  name = "My awesome connectorIPregistry"
  properties = jsonencode({
    "apiKey" = var.connectoripregistry_property_api_key
  })
}
```


## ipstack

Connector ID (`connector.id` in the resource): `connectorIPStack`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `allowInsecureIPStackConnection` (string): The Free IPStack Subscription Plan does not support HTTPS connections. For more information refer to https://ipstack.com/plan. Console display name: "Allow Insecure ipstack Connection?".
* `apiKey` (string): The ipstack API key to use the service. Console display name: "API key".


Example:
```terraform
resource "pingone_davinci_connector_instance" "connectorIPStack" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorIPStack"
  }
  name = "My awesome connectorIPStack"
  properties = jsonencode({
    "allowInsecureIPStackConnection" = var.allow_insecure_ip_stack_connection
    "apiKey" = var.connectoripstack_property_api_key
  })
}
```


## mParticle

Connector ID (`connector.id` in the resource): `mparticleConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `clientID` (string): Client ID from mParticle tenant. Console display name: "Client ID".
* `clientSecret` (string): Client Secret from mParticle tenant. Console display name: "Client Secret".
* `pod` (string): Pod from mParticle tenant. Only required for 'Upload an event batch' capability. Console display name: "Pod".


Example:
```terraform
resource "pingone_davinci_connector_instance" "mparticleConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "mparticleConnector"
  }
  name = "My awesome mparticleConnector"
  properties = jsonencode({
    "clientID" = var.mparticleconnector_property_client_i_d
    "clientSecret" = var.mparticleconnector_property_client_secret
    "pod" = var.mparticleconnector_property_pod
  })
}
```


## neoEYED

Connector ID (`connector.id` in the resource): `neoeyedConnector`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `appKey` (string): Unique key for the application. Console display name: "Application Key".
* `javascriptCdnUrl` (string): URL of javascript CDN of neoEYED. Console display name: "Javascript CDN URL".


Example:
```terraform
resource "pingone_davinci_connector_instance" "neoeyedConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "neoeyedConnector"
  }
  name = "My awesome neoeyedConnector"
  properties = jsonencode({
    "appKey" = var.neoeyedconnector_property_app_key
    "javascriptCdnUrl" = var.neoeyedconnector_property_javascript_cdn_url
  })
}
```


## randomuser.me

Connector ID (`connector.id` in the resource): `connectorRandomUserMe`

*No properties*


Example:
```terraform
resource "pingone_davinci_connector_instance" "connectorRandomUserMe" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorRandomUserMe"
  }
  name = "My awesome connectorRandomUserMe"
}
```


## tru.ID

Connector ID (`connector.id` in the resource): `connectorTruid`

Properties (used under the `properties` block in the resource as a key in the JSON object):

* `customAuth` (json):  Console display name: "Custom Parameters".


Example:
```terraform
resource "pingone_davinci_connector_instance" "connectorTruid" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorTruid"
  }
  name = "My awesome connectorTruid"
  properties = jsonencode({
    "customAuth" = jsonencode({})
  })
}
```

