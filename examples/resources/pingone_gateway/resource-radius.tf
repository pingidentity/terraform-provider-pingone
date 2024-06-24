resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_gateway" "my_radius_gateway" {
  environment_id = pingone_environment.my_environment.id
  name           = "My RADIUS Gateway"
  enabled        = true
  type           = "RADIUS"

  radius_default_shared_secret = var.radius_default_shared_secret
  radius_davinci_policy_id     = var.radius_davinci_policy_id

  radius_clients = [
    {
      ip = "127.0.0.1"
    }
  ]

}