resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_notification_template_content" "email" {
  environment_id = pingone_environment.my_environment.id
  template_name  = "strong_authentication"
  locale         = "en"

  email = {
    body    = "Please approve this transaction with passcode $${otp}."
    subject = "BX Retail Transaction Request"

    from = {
      name    = "BX Retail"
      address = "noreply@bxretail.org"
    }
  }
}

resource "pingone_notification_template_content" "push" {
  environment_id = pingone_environment.my_environment.id
  template_name  = "strong_authentication"
  locale         = "en"

  push = {
    body  = "Please approve this transaction."
    title = "BX Retail Transaction Request"
  }
}

resource "pingone_notification_template_content" "sms" {
  environment_id = pingone_environment.my_environment.id
  template_name  = "strong_authentication"
  locale         = "en"

  sms = {
    content = "Please approve this transaction with passcode $${otp}."
    sender  = "BX Retail"
  }
}

resource "pingone_notification_template_content" "voice" {
  environment_id = pingone_environment.my_environment.id
  template_name  = "strong_authentication"
  locale         = "en"

  voice = {
    content = "Hello <pause1sec> your authentication code is <sayCharValue>$${otp}</sayCharValue><pause1sec><pause1sec><repeatMessage val=2>I repeat <pause1sec>your code is <sayCharValue>$${otp}</sayCharValue></repeatMessage>"
    type    = "Alice"
  }
}
