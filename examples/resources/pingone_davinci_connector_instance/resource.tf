resource "pingone_davinci_connector_instance" "annotation_connector_example" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "annotationConnector"
  }

  name = "myAnnotationConnector"
}

resource "pingone_davinci_connector_instance" "crowdstrike_connector_example" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "crowdStrikeConnector"
  }

  name = "CrowdStrike"

  properties = jsonencode({
    "clientId" : var.crowdstrike_client_id,
    "clientSecret" : var.crowdstrike_client_secret
  })
}
