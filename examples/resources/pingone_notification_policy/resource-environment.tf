resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_notification_policy" "environment" {
  environment_id = pingone_environment.my_environment.id

  name = "Environment Quota SMS Voice and Email"

  quota = [
    {
      type             = "ENVIRONMENT"
      delivery_methods = ["SMS", "Voice"]
      total            = 100
    },
    {
      type             = "ENVIRONMENT"
      delivery_methods = ["Email"]
      total            = 100
    }
  ]
}