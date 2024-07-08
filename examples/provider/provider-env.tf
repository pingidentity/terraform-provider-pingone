terraform {
  required_providers {
    pingone = {
      source  = "pingidentity/pingone"
      version = "1.0.0-rc2"
    }
  }
}

provider "pingone" {
}

resource "pingone_environment" "my_environment" {
  # ...
}