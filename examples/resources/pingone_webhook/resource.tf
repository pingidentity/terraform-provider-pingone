resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_webhook" "my_webhook" {
  environment_id = pingone_environment.my_environment.id

  name    = "My webhook"
  enabled = true

  http_endpoint_url = "https://audit.bxretail.org/"
  http_endpoint_headers = {
    Authorization = "Basic usernamepassword"
  }

  format = "ACTIVITY"

  filter_options = {
    included_action_types = ["ACCOUNT.LINKED", "ACCOUNT.UNLINKED"]
  }
}