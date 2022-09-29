resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_authorize_decision_endpoint" "my_awesome_decision_endpoint" {
  environment_id = pingone_environment.my_environment.id
  name           = "Awesome Decision Endpoint"

  record_recent_requests = false
}
