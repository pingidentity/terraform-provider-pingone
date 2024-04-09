provider "pingone" {
  client_id      = var.client_id
  client_secret  = var.client_secret
  environment_id = var.environment_id
  region         = var.region

  global_options {

    environment {
      // This option should not be used in environments that contain production data.  Data loss may occur.
      production_type_force_delete = true
    }

    population {
      // This option should not be used in environments that contain production data.  Data loss may occur.
      contains_users_force_delete = true
    }

  }
}
