resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_notification_policy" "user" {
  environment_id = pingone_environment.my_environment.id

  name = "User Quota SMS, Voice and Email"

  quota {
    type             = "USER"
    delivery_methods = ["SMS", "Voice"]
    total            = 30
  }

  quota {
    type             = "USER"
    delivery_methods = ["Email"]
    total            = 30
  }
}