resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_application" "my_awesome_native_app" {
  environment_id = pingone_environment.my_environment.id
  name           = "Awesome Native Application"
  enabled        = true

  oidc_options {
    type                        = "NATIVE_APP"
    grant_types                 = ["CLIENT_CREDENTIALS"]
    token_endpoint_authn_method = "CLIENT_SECRET_BASIC"
    bundle_id                   = "com.example.my_ios_app"
    package_name                = "com.example.my_android_app"

    mobile_app {
      bundle_id    = "com.example.my_ios_app"
      package_name = "com.example.my_android_app"
      # ...
    }
  }
}

resource "pingone_digital_wallet_application" "my_digital_wallet_app" {
  environment_id = pingone_environment.my_environment.id
  application_id = pingone_application.my_awesome_native_app.id
  name           = "Awesome Digital Wallet Application"
  app_open_url   = "https://wallet.example.com/appopen"
}