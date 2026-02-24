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
