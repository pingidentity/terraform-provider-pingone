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
