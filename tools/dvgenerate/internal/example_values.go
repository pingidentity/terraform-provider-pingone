package internal

type ExampleValue struct {
	OverridingType *string
	Value          string
}

var (
	ExampleValues = map[string]map[string]ExampleValue{
		"connector1Kosmos": {
			"customAuth": ExampleValue{Value: `jsonencode({})`},
		},

		"connectorAWSLambda": {
			// "accessKeyId":     ExampleValue{Value: "var.aws_access_key_id"},
			"region": ExampleValue{Value: "\"eu-west-1\""},
			// "secretAccessKey": ExampleValue{Value: "var.aws_secret_access_key"},
		},

		"awsIdpConnector": {
			"openId": ExampleValue{Value: `jsonencode({})`},
		},

		"connectorAmazonAwsSecretsManager": {
			// "accessKeyId":     ExampleValue{Value: "var.aws_access_key_id"},
			"region": ExampleValue{Value: "\"eu-west-1\""},
			// "secretAccessKey": ExampleValue{Value: "var.aws_secret_access_key"},
		},

		"connectorAbuseipdb": {
			// "apiKey": ExampleValue{Value: "var.api_key"},
		},

		"connectorAcuant": {
			"customAuth": ExampleValue{Value: `jsonencode({})`},
		},

		"adobemarketoConnector": {
			// "clientId":     ExampleValue{Value: "var.client_id"},
			// "clientSecret": ExampleValue{Value: "var.client_secret"},
			// "endpoint":     ExampleValue{Value: "# property value"},
		},

		"connectorAllthenticate": {
			"customAuth": ExampleValue{Value: `jsonencode({})`},
		},

		"connectorAmazonDynamoDB": {
			// "awsAccessKey":    ExampleValue{Value: "var.aws_access_key_id"},
			// "awsAccessSecret": ExampleValue{Value: "var.aws_secret_access_key"},
			"awsRegion": ExampleValue{Value: "\"eu-west-1\""},
		},

		"amazonSimpleEmailConnector": {
			// "awsAccessKey":    ExampleValue{Value: "var.aws_access_key_id"},
			// "awsAccessSecret": ExampleValue{Value: "var.aws_secret_access_key"},
			"awsRegion": ExampleValue{Value: "\"eu-west-1\""},
			"from":      ExampleValue{Value: "\"support@bxretail.org\""},
		},

		"appleConnector": {
			"customAuth": ExampleValue{Value: `jsonencode({
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
			  })`},
		},

		"argyleConnector": {
			// "apiUrl":           ExampleValue{Value: "# property value"},
			// "clientId":     ExampleValue{Value: "var.client_id"},
			// "clientSecret": ExampleValue{Value: "var.client_secret"},
			// "javascriptWebUrl": ExampleValue{Value: "# property value"},
			// "pluginKey":        ExampleValue{Value: "# property value"},
		},

		"connectorAsignio": {
			"customAuth": ExampleValue{Value: `jsonencode({})`},
		},

		"connectorAuthid": {
			"customAuth": ExampleValue{Value: `jsonencode({})`},
		},

		"authenticIdConnector": {
			// "accountAccessKey":     ExampleValue{Value: "# property value"},
			"androidSDKLicenseKey": ExampleValue{Value: "var.authenticidconnector_property_android_sdk_license_key"},
			// "apiUrl":               ExampleValue{Value: "# property value"},
			"baseURL": ExampleValue{Value: "var.base_url"},
			// "clientCertificate":    ExampleValue{Value: "# property value"},
			// "clientKey":            ExampleValue{Value: "# property value"},
			"iOSSDKLicenseKey": ExampleValue{Value: "var.authenticidconnector_property_ios_sdk_license_key"},
			// "passphrase":           ExampleValue{Value: "# property value"},
			// "secretToken":          ExampleValue{Value: "# property value"},
		},

		"connectorAuthomize": {
			// "apiKey": ExampleValue{Value: "var.api_key"},
		},

		"azureUserManagementConnector": {
			"baseURL": ExampleValue{Value: "var.base_url"},
			// "customApiUrl": ExampleValue{Value: "# property value"},
			"customAuth": ExampleValue{Value: `jsonencode({})`},
		},

		"connectorBadge": {
			"customAuth": ExampleValue{Value: `jsonencode({})`},
		},

		"bambooConnector": {
			// "apiKey":  ExampleValue{Value: "var.api_key"},
			"baseURL": ExampleValue{Value: "var.base_url"},
			// "companySubDomain": ExampleValue{Value: "# property value"},
			// "flowId":           ExampleValue{Value: "# property value"},
			// "skWebhookUri":     ExampleValue{Value: "# property value"},
			// "webhookToken":     ExampleValue{Value: "# property value"},
		},

		"connectorBerbix": {
			// "domainName": ExampleValue{Value: "# property value"},
			// "path":       ExampleValue{Value: "# property value"},
			// "username": ExampleValue{Value: "var.username"},
		},

		"connectorBeyondIdentity": {
			"openId": ExampleValue{Value: `jsonencode({})`},
		},

		"connectorBTps": {
			// "apiKey": ExampleValue{Value: "var.api_key"},
			// "apiUser": ExampleValue{Value: "# property value"},
			// "domain":  ExampleValue{Value: "# property value"},
		},

		"connectorBTpra": {
			// "clientID":     ExampleValue{Value: "var.client_id"},
			// "clientSecret": ExampleValue{Value: "var.client_secret"},
			"praAPIurl": ExampleValue{Value: "var.pra_api_url"},
		},

		"connectorBTrs": {
			// "clientID":     ExampleValue{Value: "var.client_id"},
			// "clientSecret": ExampleValue{Value: "var.client_secret"},
			"rsAPIurl": ExampleValue{Value: "var.rs_api_url"},
		},

		"biocatchConnector": {
			// "apiUrl":           ExampleValue{Value: "# property value"},
			// "customerId":       ExampleValue{Value: "# property value"},
			// "javascriptCdnUrl": ExampleValue{Value: "# property value"},
			// "sdkToken":         ExampleValue{Value: "# property value"},
			// "truthApiKey":      ExampleValue{Value: "# property value"},
			// "truthApiUrl":      ExampleValue{Value: "# property value"},
		},

		"bitbucketIdpConnector": {
			"oauth2": ExampleValue{Value: `jsonencode({})`},
		},

		"castleConnector": {
			// "apiSecret": ExampleValue{Value: "# property value"},
		},

		"connectorCircleAccess": {
			// "appKey":      ExampleValue{Value: "# property value"},
			"customAuth": ExampleValue{Value: `jsonencode({})`},
			// "loginUrl":    ExampleValue{Value: "# property value"},
			// "readKey":     ExampleValue{Value: "# property value"},
			// "returnToUrl": ExampleValue{Value: "# property value"},
			// "writeKey":    ExampleValue{Value: "# property value"},
		},

		"connectorClearbit": {
			// "apiKey": ExampleValue{Value: "var.api_key"},
			// "riskApiVersion": ExampleValue{Value: "# property value"},
			// "version":        ExampleValue{Value: "# property value"},
		},

		"connectorCloudflare": {
			// "accountId": ExampleValue{Value: "# property value"},
			// "apiToken":  ExampleValue{Value: "# property value"},
		},

		"codeSnippetConnector": {
			// "code":         ExampleValue{Value: "# property value"},
			// "inputSchema":  ExampleValue{Value: "# property value"},
			// "outputSchema": ExampleValue{Value: "# property value"},
		},

		"complyAdvatangeConnector": {
			// "apiKey":  ExampleValue{Value: "var.api_key"},
			"baseURL": ExampleValue{Value: "var.base_url"},
		},

		"connectIdConnector": {
			"customAuth": ExampleValue{Value: `jsonencode({})`},
		},

		"cookieConnector": {
			// "hmacSigningKey": ExampleValue{Value: "# property value"},
		},

		"credovaConnector": {
			"baseURL": ExampleValue{Value: "var.base_url"},
			// "password": ExampleValue{Value: "var.password"},
			// "username": ExampleValue{Value: "var.username"},
		},

		"crowdStrikeConnector": {
			"baseURL": ExampleValue{Value: "var.base_url"},
			// "clientId":     ExampleValue{Value: "var.client_id"},
			// "clientSecret": ExampleValue{Value: "var.client_secret"},
		},

		"connectorDaonidv": {
			"openId": ExampleValue{Value: `jsonencode({})`},
		},

		"daonConnector": {
			// "apiUrl":   ExampleValue{Value: "# property value"},
			// "password": ExampleValue{Value: "var.password"},
			// "username": ExampleValue{Value: "var.username"},
		},

		"dataZooConnector": {
			// "password": ExampleValue{Value: "var.password"},
			// "username": ExampleValue{Value: "var.username"},
		},

		"connector-oai-datadogapi": {
			// "authApiKey":         ExampleValue{Value: "# property value"},
			// "authApplicationKey": ExampleValue{Value: "# property value"},
			// "basePath":           ExampleValue{Value: "# property value"},
		},

		"connectorDeBounce": {
			// "apiKey": ExampleValue{Value: "var.api_key"},
		},

		"digilockerConnector": {
			"oauth2": ExampleValue{Value: `jsonencode({})`},
		},

		"digidentityConnector": {
			"oauth2": ExampleValue{Value: `jsonencode({})`},
		},

		"duoConnector": {
			"customAuth": ExampleValue{Value: `jsonencode({})`},
		},

		"entrustConnector": {
			// "applicationId": ExampleValue{Value: "# property value"},
			// "serviceDomain": ExampleValue{Value: "# property value"},
		},

		"equifaxConnector": {
			"baseURL": ExampleValue{Value: "var.base_url"},
			// "clientId":     ExampleValue{Value: "var.client_id"},
			// "clientSecret": ExampleValue{Value: "var.client_secret"},
			// "equifaxSoapApiEnvironment": ExampleValue{Value: "# property value"},
			// "memberNumber":              ExampleValue{Value: "# property value"},
			// "password": ExampleValue{Value: "var.password"},
			// "username": ExampleValue{Value: "var.username"},
		},

		"facebookIdpConnector": {
			"oauth2": ExampleValue{Value: `jsonencode({
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
			  })`},
		},

		"fingerprintjsConnector": {
			// "apiToken":         ExampleValue{Value: "# property value"},
			// "javascriptCdnUrl": ExampleValue{Value: "# property value"},
			// "token":            ExampleValue{Value: "# property value"},
		},

		"finicityConnector": {
			// "appKey":        ExampleValue{Value: "# property value"},
			"baseURL": ExampleValue{Value: "var.base_url"},
			// "partnerId":     ExampleValue{Value: "# property value"},
			// "partnerSecret": ExampleValue{Value: "# property value"},
		},

		"flowConnector": {
			// "enforcedSignedToken": ExampleValue{Value: "# property value"},
			// "inputSchema":         ExampleValue{Value: "# property value"},
			// "pemPublicKey":        ExampleValue{Value: "# property value"},
		},

		"connectorFreshdesk": {
			// "apiKey":  ExampleValue{Value: "var.api_key"},
			"baseURL": ExampleValue{Value: "var.base_url"},
			// "version": ExampleValue{Value: "# property value"},
		},

		"connectorFreshservice": {
			// "apiKey": ExampleValue{Value: "var.api_key"},
			// "domain": ExampleValue{Value: "# property value"},
		},

		"gbgConnector": {
			// "password": ExampleValue{Value: "var.password"},
			// "requestUrl": ExampleValue{Value: "# property value"},
			// "soapAction": ExampleValue{Value: "# property value"},
			// "username": ExampleValue{Value: "var.username"},
		},

		"githubIdpConnector": {
			"oauth2": ExampleValue{Value: `jsonencode({})`},
		},

		"connectorGoogleanalyticsUA": {
			"trackingID": ExampleValue{Value: "var.tracking_id"},
			// "version":    ExampleValue{Value: "# property value"},
		},

		"connectorGoogleChromeEnterprise": {
			"customAuth": ExampleValue{Value: `jsonencode({})`},
		},

		"googleConnector": {
			"openId": ExampleValue{Value: `jsonencode({})`},
		},

		"googleWorkSpaceAdminConnector": {
			// "iss":        ExampleValue{Value: "# property value"},
			// "privateKey": ExampleValue{Value: "# property value"},
			// "sub":        ExampleValue{Value: "# property value"},
		},

		"httpConnector": {
			// "connectionId":       ExampleValue{Value: "# property value"},
			// "recaptchaSecretKey": ExampleValue{Value: "# property value"},
			// "recaptchaSiteKey":   ExampleValue{Value: "# property value"},
		},

		"connectorHuman": {
			// "humanAuthenticationToken": ExampleValue{Value: "# property value"},
			"humanCustomerID": ExampleValue{Value: "var.human_customer_id"},
			// "humanPolicyName":          ExampleValue{Value: "# property value"},
		},

		"humanCompromisedConnector": {
			// "appId":     ExampleValue{Value: "# property value"},
			// "authToken": ExampleValue{Value: "# property value"},
		},

		"connectorHyprAdapt": {
			// "accessToken": ExampleValue{Value: "# property value"},
		},

		"hyprConnector": {
			"customAuth": ExampleValue{Value: `jsonencode({})`},
		},

		"haveIBeenPwnedConnector": {
			// "apiKey": ExampleValue{Value: "var.api_key"},
			// "apiUrl":    ExampleValue{Value: "# property value"},
			// "userAgent": ExampleValue{Value: "# property value"},
		},

		"connectorHello": {
			"customAuth": ExampleValue{Value: `jsonencode({})`},
		},

		"connectorHubspot": {
			// "bearerToken": ExampleValue{Value: "# property value"},
		},

		"idDatawebConnector": {
			"customAuth": ExampleValue{Value: `jsonencode({})`},
		},

		"idranddConnector": {
			// "apiKey": ExampleValue{Value: "var.api_key"},
			// "apiUrl": ExampleValue{Value: "# property value"},
		},

		"connectorIdMeIdentity": {
			"openId": ExampleValue{Value: `jsonencode({})`},
		},

		"idMeConnector": {
			"oauth2": ExampleValue{Value: `jsonencode({})`},
		},

		"idemiaConnector": {
			// "apiKey":  ExampleValue{Value: "var.api_key"},
			"baseURL": ExampleValue{Value: "var.base_url"},
		},

		"skPeopleIntelligenceConnector": {
			// "authUrl":      ExampleValue{Value: "# property value"},
			// "clientId":     ExampleValue{Value: "var.client_id"},
			// "clientSecret": ExampleValue{Value: "var.client_secret"},
			// "dppa":         ExampleValue{Value: "# property value"},
			// "glba":         ExampleValue{Value: "# property value"},
			// "searchUrl":    ExampleValue{Value: "# property value"},
		},

		"connectorIdmelon": {
			"customAuth": ExampleValue{Value: `jsonencode({})`},
		},

		"idmissionConnector": {
			// "authDescription":  ExampleValue{Value: "# property value"},
			// "connectorName":    ExampleValue{Value: "# property value"},
			// "description":      ExampleValue{Value: "# property value"},
			// "details1":         ExampleValue{Value: "# property value"},
			// "details2":         ExampleValue{Value: "# property value"},
			// "iconUrl":          ExampleValue{Value: "# property value"},
			// "iconUrlPng":       ExampleValue{Value: "# property value"},
			// "loginId":          ExampleValue{Value: "# property value"},
			// "merchantId":       ExampleValue{Value: "# property value"},
			// "password": ExampleValue{Value: "var.password"},
			// "productId":        ExampleValue{Value: "# property value"},
			// "productName":      ExampleValue{Value: "# property value"},
			// "showCredAddedOn":  ExampleValue{Value: "# property value"},
			// "showCredAddedVia": ExampleValue{Value: "# property value"},
			// "title":            ExampleValue{Value: "# property value"},
			// "toolTip":          ExampleValue{Value: "# property value"},
			// "url":              ExampleValue{Value: "# property value"},
		},

		"idrampOidcConnector": {
			"customAuth": ExampleValue{Value: `jsonencode({})`},
		},

		"incodeConnector": {
			"customAuth": ExampleValue{Value: `jsonencode({})`},
		},

		"connectorInfinipoint": {
			"customAuth": ExampleValue{Value: `jsonencode({})`},
		},

		"connectorJamf": {
			// "jamfPassword": ExampleValue{Value: "# property value"},
			// "jamfUsername": ExampleValue{Value: "# property value"},
			// "serverName":   ExampleValue{Value: "# property value"},
		},

		"connectorJiraServiceDesk": {
			"JIRAServiceDeskAuth":       ExampleValue{Value: "var.jira_service_desk_auth"},
			"JIRAServiceDeskCreateData": ExampleValue{Value: "var.jira_service_desk_create_data"},
			"JIRAServiceDeskURL":        ExampleValue{Value: "var.jira_service_desk_url"},
			"JIRAServiceDeskUpdateData": ExampleValue{Value: "var.jira_service_desk_update_data"},
			// "method":                    ExampleValue{Value: "# property value"},
		},

		"jiraConnector": {
			// "apiKey": ExampleValue{Value: "var.api_key"},
			// "apiUrl": ExampleValue{Value: "# property value"},
			// "email":  ExampleValue{Value: "# property value"},
		},

		"jumioConnector": {
			// "apiKey": ExampleValue{Value: "var.api_key"},
			// "authDescription":            ExampleValue{Value: "# property value"},
			// "authUrl":                    ExampleValue{Value: "# property value"},
			// "authorizationTokenLifetime": ExampleValue{Value: "# property value"},
			// "baseColor":                  ExampleValue{Value: "# property value"},
			// "bgColor":                    ExampleValue{Value: "# property value"},
			// "callbackUrl":                ExampleValue{Value: "# property value"},
			// "clientSecret": ExampleValue{Value: "var.client_secret"},
			// "connectorName":              ExampleValue{Value: "# property value"},
			// "description":                ExampleValue{Value: "# property value"},
			// "details1":                   ExampleValue{Value: "# property value"},
			// "details2":                   ExampleValue{Value: "# property value"},
			// "doNotShowInIframe":          ExampleValue{Value: "# property value"},
			// "docVerificationUrl":         ExampleValue{Value: "# property value"},
			// "headerImageUrl":             ExampleValue{Value: "# property value"},
			// "iconUrl":                    ExampleValue{Value: "# property value"},
			// "iconUrlPng":                 ExampleValue{Value: "# property value"},
			// "locale":                     ExampleValue{Value: "# property value"},
			// "showCredAddedOn":            ExampleValue{Value: "# property value"},
			// "showCredAddedVia":           ExampleValue{Value: "# property value"},
			// "title":                      ExampleValue{Value: "# property value"},
			// "toolTip":                    ExampleValue{Value: "# property value"},
		},

		"kbaConnector": {
			// "authDescription":  ExampleValue{Value: "# property value"},
			// "connectorName":    ExampleValue{Value: "# property value"},
			// "description":      ExampleValue{Value: "# property value"},
			// "details1":         ExampleValue{Value: "# property value"},
			// "details2":         ExampleValue{Value: "# property value"},
			// "formFieldsList":   ExampleValue{Value: "# property value"},
			// "iconUrl":          ExampleValue{Value: "# property value"},
			// "iconUrlPng":       ExampleValue{Value: "# property value"},
			// "showCredAddedOn":  ExampleValue{Value: "# property value"},
			// "showCredAddedVia": ExampleValue{Value: "# property value"},
			// "title":            ExampleValue{Value: "# property value"},
			// "toolTip":          ExampleValue{Value: "# property value"},
		},

		"kaizenVoizConnector": {
			// "apiUrl":           ExampleValue{Value: "# property value"},
			// "applicationName":  ExampleValue{Value: "# property value"},
			// "authDescription":  ExampleValue{Value: "# property value"},
			// "connectorName":    ExampleValue{Value: "# property value"},
			// "description":      ExampleValue{Value: "# property value"},
			// "details1":         ExampleValue{Value: "# property value"},
			// "details2":         ExampleValue{Value: "# property value"},
			// "iconUrl":          ExampleValue{Value: "# property value"},
			// "iconUrlPng":       ExampleValue{Value: "# property value"},
			// "showCredAddedOn":  ExampleValue{Value: "# property value"},
			// "showCredAddedVia": ExampleValue{Value: "# property value"},
			// "title":            ExampleValue{Value: "# property value"},
			// "toolTip":          ExampleValue{Value: "# property value"},
		},

		"connectorKeyless": {
			"customAuth": ExampleValue{Value: `jsonencode({})`},
		},

		"pingOneLDAPConnector": {
			// "clientId":     ExampleValue{Value: "var.client_id"},
			// "clientSecret": ExampleValue{Value: "var.client_secret"},
			// "envId":        ExampleValue{Value: "# property value"},
			// "gatewayId":    ExampleValue{Value: "# property value"},
			// "region":       ExampleValue{Value: "# property value"},
		},

		"lexisnexisV2Connector": {
			// "apiKey": ExampleValue{Value: "var.api_key"},
			// "apiUrl":          ExampleValue{Value: "# property value"},
			// "orgId":           ExampleValue{Value: "# property value"},
			"useCustomApiURL": ExampleValue{Value: "var.use_custom_api_url"},
		},

		"linkedInConnector": {
			"oauth2": ExampleValue{Value: `jsonencode({})`},
		},

		"connectorMailchimp": {
			// "transactionalApiKey":     ExampleValue{Value: "# property value"},
			// "transactionalApiVersion": ExampleValue{Value: "# property value"},
		},

		"connectorMailgun": {
			// "apiKey": ExampleValue{Value: "var.api_key"},
			// "apiVersion":    ExampleValue{Value: "# property value"},
			// "mailgunDomain": ExampleValue{Value: "# property value"},
		},

		"melissaConnector": {
			// "apiKey": ExampleValue{Value: "var.api_key"},
		},

		"connectorMicrosoftIntune": {
			// "clientId":     ExampleValue{Value: "var.client_id"},
			// "clientSecret": ExampleValue{Value: "var.client_secret"},
			// "domainName":   ExampleValue{Value: "# property value"},
			// "grantType":    ExampleValue{Value: "# property value"},
			// "scope":        ExampleValue{Value: "# property value"},
			// "tenant":       ExampleValue{Value: "# property value"},
		},

		"microsoftIdpConnector": {
			"openId": ExampleValue{Value: `jsonencode({})`},
		},

		"microsoftTeamsConnector": {
			"customAuth": ExampleValue{Value: `jsonencode({})`},
		},

		"nuanceConnector": {
			// "authDescription":  ExampleValue{Value: "# property value"},
			// "configSetName":    ExampleValue{Value: "# property value"},
			// "connectorName":    ExampleValue{Value: "# property value"},
			// "description":      ExampleValue{Value: "# property value"},
			// "details1":         ExampleValue{Value: "# property value"},
			// "details2":         ExampleValue{Value: "# property value"},
			// "iconUrl":          ExampleValue{Value: "# property value"},
			// "iconUrlPng":       ExampleValue{Value: "# property value"},
			// "passphrase1":      ExampleValue{Value: "# property value"},
			// "passphrase2":      ExampleValue{Value: "# property value"},
			// "passphrase3":      ExampleValue{Value: "# property value"},
			// "passphrase4":      ExampleValue{Value: "# property value"},
			// "passphrase5":      ExampleValue{Value: "# property value"},
			// "showCredAddedOn":  ExampleValue{Value: "# property value"},
			// "showCredAddedVia": ExampleValue{Value: "# property value"},
			// "title":            ExampleValue{Value: "# property value"},
			// "toolTip":          ExampleValue{Value: "# property value"},
		},

		"genericConnector": {
			"customAuth": ExampleValue{Value: `jsonencode({})`},
		},

		"connectorOpswat": {
			// "clientId":     ExampleValue{Value: "var.client_id"},
			// "clientSecret": ExampleValue{Value: "var.client_secret"},
			// "crossDomainApiPort": ExampleValue{Value: "# property value"},
			// "maDomain":           ExampleValue{Value: "# property value"},
		},

		"oneTrustConnector": {
			// "clientId":     ExampleValue{Value: "var.client_id"},
			// "clientSecret": ExampleValue{Value: "var.client_secret"},
		},

		"onfidoConnector": {
			// "androidPackageName":        ExampleValue{Value: "# property value"},
			// "apiKey": ExampleValue{Value: "var.api_key"},
			// "authDescription":           ExampleValue{Value: "# property value"},
			"baseURL": ExampleValue{Value: "var.base_url"},
			// "connectorName":             ExampleValue{Value: "# property value"},
			// "customizeSteps":            ExampleValue{Value: "# property value"},
			// "description":               ExampleValue{Value: "# property value"},
			// "details1":                  ExampleValue{Value: "# property value"},
			// "details2":                  ExampleValue{Value: "# property value"},
			// "iOSBundleId":               ExampleValue{Value: "# property value"},
			// "iconUrl":                   ExampleValue{Value: "# property value"},
			// "iconUrlPng":                ExampleValue{Value: "# property value"},
			"javascriptCSSUrl": ExampleValue{Value: "var.javascript_css_url"},
			// "javascriptCdnUrl":          ExampleValue{Value: "# property value"},
			// "language":                  ExampleValue{Value: "# property value"},
			// "referenceStepsList":        ExampleValue{Value: "# property value"},
			// "referrerUrl":               ExampleValue{Value: "# property value"},
			// "retrieveReports":           ExampleValue{Value: "# property value"},
			// "shouldCloseOnOverlayClick": ExampleValue{Value: "# property value"},
			// "showCredAddedOn":           ExampleValue{Value: "# property value"},
			// "showCredAddedVia":          ExampleValue{Value: "# property value"},
			// "stepsList":                 ExampleValue{Value: "# property value"},
			// "title":                     ExampleValue{Value: "# property value"},
			// "toolTip":                   ExampleValue{Value: "# property value"},
			// "useLanguage":               ExampleValue{Value: "# property value"},
			// "useModal":                  ExampleValue{Value: "# property value"},
			// "viewDescriptions":          ExampleValue{Value: "# property value"},
			// "viewTitle":                 ExampleValue{Value: "# property value"},
		},

		"connectorPaloAltoPrisma": {
			"baseURL": ExampleValue{Value: "var.base_url"},
			// "prismaPassword": ExampleValue{Value: "# property value"},
			// "prismaUsername": ExampleValue{Value: "# property value"},
		},

		"connector-oai-pingaccessadministrativeapi": {
			// "authPassword":    ExampleValue{Value: "# property value"},
			// "authUsername":    ExampleValue{Value: "# property value"},
			// "basePath":        ExampleValue{Value: "# property value"},
			// "sslVerification": ExampleValue{Value: "# property value"},
		},

		"connector-oai-pfadminapi": {
			// "authPassword":    ExampleValue{Value: "# property value"},
			// "authUsername":    ExampleValue{Value: "# property value"},
			// "basePath":        ExampleValue{Value: "# property value"},
			// "sslVerification": ExampleValue{Value: "# property value"},
		},

		"pingFederateConnectorV2": {
			"openId": ExampleValue{Value: `jsonencode({
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
			  })`},
		},

		"pingIdConnector": {
			"customAuth": ExampleValue{Value: `jsonencode({
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
			  })`},
		},

		"pingOneAuthorizeConnector": {
			// "clientId":     ExampleValue{Value: "var.client_id"},
			// "clientSecret": ExampleValue{Value: "var.client_secret"},
			"endpointURL": ExampleValue{Value: "var.endpoint_url"},
		},

		"pingOneCredentialsConnector": {
			"clientId":     ExampleValue{Value: "var.pingone_worker_app_client_id"},
			"clientSecret": ExampleValue{Value: "var.pingone_worker_app_client_secret"},
			// "digitalWalletApplicationId": ExampleValue{Value: "# property value"},
			"envId": ExampleValue{Value: "var.pingone_worker_app_environment_id"},
			// "region":                     ExampleValue{Value: "# property value"},
		},

		"pingOneMfaConnector": {
			"clientId":     ExampleValue{Value: "var.pingone_worker_app_client_id"},
			"clientSecret": ExampleValue{Value: "var.pingone_worker_app_client_secret"},
			"envId":        ExampleValue{Value: "var.pingone_worker_app_environment_id"},
			// "policyId":     ExampleValue{Value: "# property value"},
			// "region":       ExampleValue{Value: "# property value"},
		},

		"notificationsConnector": {
			// "clientId":     ExampleValue{Value: "var.client_id"},
			// "clientSecret": ExampleValue{Value: "var.client_secret"},
			// "envId":                ExampleValue{Value: "# property value"},
			// "notificationPolicyId": ExampleValue{Value: "# property value"},
			// "region":               ExampleValue{Value: "# property value"},
		},

		"pingOneRiskConnector": {
			"clientId":     ExampleValue{Value: "var.pingone_worker_app_client_id"},
			"clientSecret": ExampleValue{Value: "var.pingone_worker_app_client_secret"},
			"envId":        ExampleValue{Value: "var.pingone_worker_app_environment_id"},
			// "region":       ExampleValue{Value: "# property value"},
		},

		"pingOneScopeConsentConnector": {
			"clientId":     ExampleValue{Value: "var.pingone_worker_app_client_id"},
			"clientSecret": ExampleValue{Value: "var.pingone_worker_app_client_secret"},
			"envId":        ExampleValue{Value: "var.pingone_worker_app_environment_id"},
			// "region":       ExampleValue{Value: "# property value"},
		},

		"pingOneVerifyConnector": {
			"clientId":     ExampleValue{Value: "var.pingone_worker_app_client_id"},
			"clientSecret": ExampleValue{Value: "var.pingone_worker_app_client_secret"},
			"envId":        ExampleValue{Value: "var.pingone_worker_app_environment_id"},
			// "region":       ExampleValue{Value: "# property value"},
		},

		"pingOneSSOConnector": {
			"clientId":     ExampleValue{Value: "var.pingone_worker_app_client_id"},
			"clientSecret": ExampleValue{Value: "var.pingone_worker_app_client_secret"},
			"envId":        ExampleValue{Value: "var.pingone_worker_app_environment_id"},
			// "region":       ExampleValue{Value: "# property value"},
		},

		"proveConnector": {
			"baseURL": ExampleValue{Value: "var.base_url"},
			// "clientId": ExampleValue{Value: "var.client_id"},
			// "grantType": ExampleValue{Value: "# property value"},
			// "password": ExampleValue{Value: "var.password"},
			// "username": ExampleValue{Value: "var.username"},
		},

		"payfoneConnector": {
			// "appClientId":          ExampleValue{Value: "# property value"},
			"baseURL": ExampleValue{Value: "var.base_url"},
			// "clientId": ExampleValue{Value: "var.client_id"},
			// "password": ExampleValue{Value: "var.password"},
			// "simulatorMode":        ExampleValue{Value: "# property value"},
			// "simulatorPhoneNumber": ExampleValue{Value: "# property value"},
			// "skCallbackBaseUrl":    ExampleValue{Value: "# property value"},
			// "username": ExampleValue{Value: "var.username"},
		},

		"rsaConnector": {
			// "accessId":  ExampleValue{Value: "# property value"},
			// "accessKey": ExampleValue{Value: "# property value"},
			"baseURL": ExampleValue{Value: "var.base_url"},
		},

		"inveridConnector": {
			// "getApiKey":    ExampleValue{Value: "# property value"},
			// "host":         ExampleValue{Value: "# property value"},
			// "postApiKey":   ExampleValue{Value: "# property value"},
			// "skWebhookUri": ExampleValue{Value: "# property value"},
			// "timeToLive":   ExampleValue{Value: "# property value"},
		},

		"connectorIdiVERIFIED": {
			// "apiSecret":  ExampleValue{Value: "# property value"},
			// "companyKey": ExampleValue{Value: "# property value"},
			// "idiEnv":     ExampleValue{Value: "# property value"},
			// "siteKey":    ExampleValue{Value: "# property value"},
			// "uniqueUrl":  ExampleValue{Value: "# property value"},
		},

		"samlIdpConnector": {
			"saml": ExampleValue{Value: `jsonencode({})`},
		},

		"seonConnector": {
			"baseURL": ExampleValue{Value: "var.base_url"},
			// "licenseKey": ExampleValue{Value: "# property value"},
		},

		"smtpConnector": {
			// "hostname":   ExampleValue{Value: "# property value"},
			// "name":       ExampleValue{Value: "# property value"},
			// "password": ExampleValue{Value: "var.password"},
			// "port":       ExampleValue{Value: "# property value"},
			// "secureFlag": ExampleValue{Value: "# property value"},
			// "username": ExampleValue{Value: "var.username"},
		},

		"connectorIdentityNow": {
			// "clientId":     ExampleValue{Value: "var.client_id"},
			// "clientSecret": ExampleValue{Value: "var.client_secret"},
			// "tenant":       ExampleValue{Value: "# property value"},
		},

		"connectorSalesforceMarketingCloud": {
			"SalesforceMarketingCloudURL": ExampleValue{Value: "var.salesforce_marketing_cloud_url"},
			// "accountId":                   ExampleValue{Value: "# property value"},
			// "clientId":     ExampleValue{Value: "var.client_id"},
			// "clientSecret": ExampleValue{Value: "var.client_secret"},
			// "scope":                       ExampleValue{Value: "# property value"},
		},

		"salesforceConnector": {
			// "adminUsername": ExampleValue{Value: "# property value"},
			// "consumerKey":   ExampleValue{Value: "# property value"},
			// "domainName":    ExampleValue{Value: "# property value"},
			// "environment":   ExampleValue{Value: "# property value"},
			// "privateKey":    ExampleValue{Value: "# property value"},
		},

		"connectorSaviyntFlow": {
			// "domainName":      ExampleValue{Value: "# property value"},
			// "path":            ExampleValue{Value: "# property value"},
			// "saviyntPassword": ExampleValue{Value: "# property value"},
			// "saviyntUserName": ExampleValue{Value: "# property value"},
		},

		"securIdConnector": {
			// "apiUrl":    ExampleValue{Value: "# property value"},
			// "clientKey": ExampleValue{Value: "# property value"},
		},

		"connectorSecuronix": {
			// "domainName": ExampleValue{Value: "# property value"},
			// "token":      ExampleValue{Value: "# property value"},
		},

		"connectorSegment": {
			// "version":  ExampleValue{Value: "# property value"},
			// "writeKey": ExampleValue{Value: "# property value"},
		},

		"sentilinkConnector": {
			// "account":          ExampleValue{Value: "# property value"},
			// "apiUrl":           ExampleValue{Value: "# property value"},
			// "javascriptCdnUrl": ExampleValue{Value: "# property value"},
			// "token":            ExampleValue{Value: "# property value"},
		},

		"servicenowConnector": {
			// "adminUsername": ExampleValue{Value: "# property value"},
			// "apiUrl":        ExampleValue{Value: "# property value"},
			// "password": ExampleValue{Value: "var.password"},
		},

		"connectorShopify": {
			// "accessToken":          ExampleValue{Value: "# property value"},
			// "apiVersion":           ExampleValue{Value: "# property value"},
			// "multipassSecret":      ExampleValue{Value: "# property value"},
			// "multipassStoreDomain": ExampleValue{Value: "# property value"},
			// "yourStoreName":        ExampleValue{Value: "# property value"},
		},

		"connectorSignicat": {
			"customAuth": ExampleValue{Value: `jsonencode({})`},
		},

		"singpassLoginConnector": {
			"customAuth": ExampleValue{Value: `jsonencode({})`},
		},

		"slackConnector": {
			"oauth2": ExampleValue{Value: `jsonencode({})`},
		},

		"connectorSmarty": {
			// "authId":    ExampleValue{Value: "# property value"},
			// "authToken": ExampleValue{Value: "# property value"},
			// "license":   ExampleValue{Value: "# property value"},
		},

		"socureConnector": {
			// "apiKey":  ExampleValue{Value: "var.api_key"},
			"baseURL": ExampleValue{Value: "var.base_url"},
			// "customApiUrl": ExampleValue{Value: "# property value"},
			// "sdkKey":       ExampleValue{Value: "# property value"},
			// "skWebhookUri": ExampleValue{Value: "# property value"},
		},

		"splunkConnector": {
			// "apiUrl": ExampleValue{Value: "# property value"},
			// "port":   ExampleValue{Value: "# property value"},
			// "token":  ExampleValue{Value: "# property value"},
		},

		"connectorSpotify": {
			"oauth2": ExampleValue{Value: `jsonencode({})`},
		},

		"connectorSpycloud": {
			// "apiKey": ExampleValue{Value: "var.api_key"},
		},

		"connectorSvipe": {
			"customAuth": ExampleValue{Value: `jsonencode({})`},
		},

		"tmtConnector": {
			// "apiKey": ExampleValue{Value: "var.api_key"},
			// "apiSecret": ExampleValue{Value: "# property value"},
			// "apiUrl":    ExampleValue{Value: "# property value"},
		},

		"connectorTableau": {
			// "addFlowPermissionsRequestBody": ExampleValue{Value: "# property value"},
			// "addUsertoSiteRequestBody":      ExampleValue{Value: "# property value"},
			// "apiVersion":                    ExampleValue{Value: "# property value"},
			// "authId":                        ExampleValue{Value: "# property value"},
			// "createScheduleBody":            ExampleValue{Value: "# property value"},
			// "datasourceId":                  ExampleValue{Value: "# property value"},
			// "flowId":                        ExampleValue{Value: "# property value"},
			// "groupId":                       ExampleValue{Value: "# property value"},
			// "jobId":                         ExampleValue{Value: "# property value"},
			// "scheduleId":                    ExampleValue{Value: "# property value"},
			// "serverUrl":                     ExampleValue{Value: "# property value"},
			// "siteId":                        ExampleValue{Value: "# property value"},
			// "taskId":                        ExampleValue{Value: "# property value"},
			// "updateScheduleRequestBody":     ExampleValue{Value: "# property value"},
			// "updateUserRequestBody":         ExampleValue{Value: "# property value"},
			// "userId":                        ExampleValue{Value: "# property value"},
			// "workbookId":                    ExampleValue{Value: "# property value"},
		},

		"telesignConnector": {
			// "authDescription":  ExampleValue{Value: "# property value"},
			// "connectorName":    ExampleValue{Value: "# property value"},
			// "description":      ExampleValue{Value: "# property value"},
			// "details1":         ExampleValue{Value: "# property value"},
			// "details2":         ExampleValue{Value: "# property value"},
			// "iconUrl":          ExampleValue{Value: "# property value"},
			// "iconUrlPng":       ExampleValue{Value: "# property value"},
			// "password": ExampleValue{Value: "var.password"},
			// "providerName":     ExampleValue{Value: "# property value"},
			// "showCredAddedOn":  ExampleValue{Value: "# property value"},
			// "showCredAddedVia": ExampleValue{Value: "# property value"},
			// "title":            ExampleValue{Value: "# property value"},
			// "toolTip":          ExampleValue{Value: "# property value"},
			// "username": ExampleValue{Value: "var.username"},
		},

		"tutloxpConnector": {
			// "apiUrl":   ExampleValue{Value: "# property value"},
			// "dppaCode": ExampleValue{Value: "# property value"},
			// "glbCode":  ExampleValue{Value: "# property value"},
			// "password": ExampleValue{Value: "var.password"},
			// "username": ExampleValue{Value: "var.username"},
		},

		"transunionConnector": {
			// "apiUrl":                   ExampleValue{Value: "# property value"},
			// "docVerificationPassword":  ExampleValue{Value: "# property value"},
			// "docVerificationPublicKey": ExampleValue{Value: "# property value"},
			// "docVerificationSecret":    ExampleValue{Value: "# property value"},
			// "docVerificationSiteId":    ExampleValue{Value: "# property value"},
			// "docVerificationUsername":  ExampleValue{Value: "# property value"},
			// "idVerificationPassword":   ExampleValue{Value: "# property value"},
			// "idVerificationPublicKey":  ExampleValue{Value: "# property value"},
			// "idVerificationSecret":     ExampleValue{Value: "# property value"},
			// "idVerificationSiteId":     ExampleValue{Value: "# property value"},
			// "idVerificationUsername":   ExampleValue{Value: "# property value"},
			// "kbaPassword":              ExampleValue{Value: "# property value"},
			// "kbaPublicKey":             ExampleValue{Value: "# property value"},
			// "kbaSecret":                ExampleValue{Value: "# property value"},
			// "kbaSiteId":                ExampleValue{Value: "# property value"},
			// "kbaUsername":              ExampleValue{Value: "# property value"},
			// "otpPassword":              ExampleValue{Value: "# property value"},
			// "otpPublicKey":             ExampleValue{Value: "# property value"},
			// "otpSecret":                ExampleValue{Value: "# property value"},
			// "otpSiteId":                ExampleValue{Value: "# property value"},
			// "otpUsername":              ExampleValue{Value: "# property value"},
		},

		"connectorTrulioo": {
			// "clientID":     ExampleValue{Value: "var.client_id"},
			// "clientSecret": ExampleValue{Value: "var.client_secret"},
		},

		"twilioConnector": {
			// "accountSid":              ExampleValue{Value: "# property value"},
			// "authDescription":         ExampleValue{Value: "# property value"},
			// "authMessageTemplate":     ExampleValue{Value: "# property value"},
			// "authToken":               ExampleValue{Value: "# property value"},
			// "connectorName":           ExampleValue{Value: "# property value"},
			// "description":             ExampleValue{Value: "# property value"},
			// "details1":                ExampleValue{Value: "# property value"},
			// "details2":                ExampleValue{Value: "# property value"},
			// "iconUrl":                 ExampleValue{Value: "# property value"},
			// "iconUrlPng":              ExampleValue{Value: "# property value"},
			// "registerMessageTemplate": ExampleValue{Value: "# property value"},
			// "senderPhoneNumber":       ExampleValue{Value: "# property value"},
			// "showCredAddedOn":         ExampleValue{Value: "# property value"},
			// "showCredAddedVia":        ExampleValue{Value: "# property value"},
			// "title":                   ExampleValue{Value: "# property value"},
			// "toolTip":                 ExampleValue{Value: "# property value"},
		},

		"unifyIdConnector": {
			// "accountId":        ExampleValue{Value: "# property value"},
			// "apiKey": ExampleValue{Value: "var.api_key"},
			// "connectorName":    ExampleValue{Value: "# property value"},
			// "details1":         ExampleValue{Value: "# property value"},
			// "details2":         ExampleValue{Value: "# property value"},
			// "iconUrl":          ExampleValue{Value: "# property value"},
			// "iconUrlPng":       ExampleValue{Value: "# property value"},
			// "sdkToken":         ExampleValue{Value: "# property value"},
			// "showCredAddedOn":  ExampleValue{Value: "# property value"},
			// "showCredAddedVia": ExampleValue{Value: "# property value"},
			// "toolTip":          ExampleValue{Value: "# property value"},
		},

		"userPolicyConnector": {
			// "passwordExpiryInDays":          ExampleValue{Value: "# property value"},
			// "passwordExpiryNotification":    ExampleValue{Value: "# property value"},
			// "passwordLengthMax":             ExampleValue{Value: "# property value"},
			// "passwordLengthMin":             ExampleValue{Value: "# property value"},
			// "passwordLockoutAttempts":       ExampleValue{Value: "# property value"},
			// "passwordPreviousXPasswords":    ExampleValue{Value: "# property value"},
			// "passwordRequireLowercase":      ExampleValue{Value: "# property value"},
			// "passwordRequireNumbers":        ExampleValue{Value: "# property value"},
			// "passwordRequireSpecial":        ExampleValue{Value: "# property value"},
			// "passwordRequireUppercase":      ExampleValue{Value: "# property value"},
			// "passwordSpacesOk":              ExampleValue{Value: "# property value"},
			// "passwordsEnabled":              ExampleValue{Value: "# property value"},
			// "temporaryPasswordExpiryInDays": ExampleValue{Value: "# property value"},
		},

		"skUserPool": {
			"customAttributes": ExampleValue{Value: `jsonencode({
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
			})`},
		},

		"connectorValidsoft": {
			"customAuth": ExampleValue{Value: `jsonencode({})`},
		},

		"connectorVericlouds": {
			// "apiSecret": ExampleValue{Value: "# property value"},
			// "apiKey": ExampleValue{Value: "var.api_key"},
		},

		"veriffConnector": {
			// "access_token": ExampleValue{Value: "var.access_token"},
			"baseURL": ExampleValue{Value: "var.base_url"},
			// "password":     ExampleValue{Value: "var.password"},
		},

		"connector443id": {
			// "apiKey": ExampleValue{Value: "var.api_key"},
		},

		"webhookConnector": {
			// "urls": ExampleValue{Value: "# property value"},
		},

		"connectorWhatsAppBusiness": {
			// "accessToken":  ExampleValue{Value: "# property value"},
			// "appSecret":    ExampleValue{Value: "# property value"},
			// "skWebhookUri": ExampleValue{Value: "# property value"},
			// "verifyToken":  ExampleValue{Value: "# property value"},
			// "version":      ExampleValue{Value: "# property value"},
		},

		"connectorWinmagic": {
			"openId": ExampleValue{Value: `jsonencode({})`},
		},

		"wireWheelConnector": {
			"baseURL": ExampleValue{Value: "var.base_url"},
			// "clientId":     ExampleValue{Value: "var.client_id"},
			// "clientSecret": ExampleValue{Value: "var.client_secret"},
			// "issuerId":     ExampleValue{Value: "# property value"},
		},

		"twitterIdpConnector": {
			"customAuth": ExampleValue{Value: `jsonencode({})`},
		},

		"yotiConnector": {
			"customAuth": ExampleValue{Value: `jsonencode({})`},
		},

		"connectorZendesk": {
			// "apiToken":      ExampleValue{Value: "# property value"},
			// "emailUsername": ExampleValue{Value: "# property value"},
			// "subdomain":     ExampleValue{Value: "# property value"},
		},

		"zoopConnector": {
			// "agencyId": ExampleValue{Value: "# property value"},
			// "apiKey": ExampleValue{Value: "var.api_key"},
			// "apiUrl":   ExampleValue{Value: "# property value"},
		},

		"connectorZscaler": {
			// "basePath":        ExampleValue{Value: "# property value"},
			"baseURL":       ExampleValue{Value: "var.base_url"},
			"zscalerAPIkey": ExampleValue{Value: "var.zscaler_api_key"},
			// "zscalerPassword": ExampleValue{Value: "# property value"},
			// "zscalerUsername": ExampleValue{Value: "# property value"},
		},

		"iproovConnector": {
			// "allowLandscape":       ExampleValue{Value: "# property value"},
			// "apiKey": ExampleValue{Value: "var.api_key"},
			// "authDescription":      ExampleValue{Value: "# property value"},
			"baseURL": ExampleValue{Value: "var.base_url"},
			// "color1":               ExampleValue{Value: "# property value"},
			// "color2":               ExampleValue{Value: "# property value"},
			// "color3":               ExampleValue{Value: "# property value"},
			// "color4":               ExampleValue{Value: "# property value"},
			// "connectorName":        ExampleValue{Value: "# property value"},
			// "customTitle":          ExampleValue{Value: "# property value"},
			// "description":          ExampleValue{Value: "# property value"},
			// "details1":             ExampleValue{Value: "# property value"},
			// "details2":             ExampleValue{Value: "# property value"},
			// "enableCameraSelector": ExampleValue{Value: "# property value"},
			// "iconUrl":              ExampleValue{Value: "# property value"},
			// "iconUrlPng":           ExampleValue{Value: "# property value"},
			"javascriptCSSUrl": ExampleValue{Value: "var.javascript_css_url"},
			// "javascriptCdnUrl":     ExampleValue{Value: "# property value"},
			// "kioskMode":            ExampleValue{Value: "# property value"},
			// "logo":                 ExampleValue{Value: "# property value"},
			// "password": ExampleValue{Value: "var.password"},
			// "secret":               ExampleValue{Value: "# property value"},
			// "showCountdown":        ExampleValue{Value: "# property value"},
			// "showCredAddedOn":      ExampleValue{Value: "# property value"},
			// "showCredAddedVia":     ExampleValue{Value: "# property value"},
			// "startScreenTitle":     ExampleValue{Value: "# property value"},
			// "title":                ExampleValue{Value: "# property value"},
			// "toolTip":              ExampleValue{Value: "# property value"},
			// "username": ExampleValue{Value: "var.username"},
		},

		"iovationConnector": {
			// "apiUrl":             ExampleValue{Value: "# property value"},
			// "javascriptCdnUrl":   ExampleValue{Value: "# property value"},
			// "subKey":             ExampleValue{Value: "# property value"},
			// "subscriberAccount":  ExampleValue{Value: "# property value"},
			// "subscriberId":       ExampleValue{Value: "# property value"},
			// "subscriberPasscode": ExampleValue{Value: "# property value"},
			// "version":            ExampleValue{Value: "# property value"},
		},

		"connectorIPGeolocationio": {
			// "apiKey": ExampleValue{Value: "var.api_key"},
		},

		"connectorIPregistry": {
			// "apiKey": ExampleValue{Value: "var.api_key"},
		},

		"connectorIPStack": {
			"allowInsecureIPStackConnection": ExampleValue{Value: "var.allow_insecure_ip_stack_connection"},
			// "apiKey":                         ExampleValue{Value: "var.api_key"},
		},

		"neoeyedConnector": {
			// "appKey":           ExampleValue{Value: "# property value"},
			// "javascriptCdnUrl": ExampleValue{Value: "# property value"},
		},

		"connectorTruid": {
			"customAuth": ExampleValue{Value: `jsonencode({})`},
		},
	}
)
