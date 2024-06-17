provider "pingone" {
  client_id      = var.client_id
  client_secret  = var.client_secret
  environment_id = var.environment_id
  region_code    = var.region_code

  global_options {

    population {
      // This option should not be used in environments that contain production data.  Data loss may occur.
      contains_users_force_delete = true
    }

  }
}
