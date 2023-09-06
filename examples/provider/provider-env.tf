terraform {
  required_providers {
    pingone = {
      source  = "pingidentity/pingone"
      version = "~> 0.21"
    }
  }
}

provider "pingone" {
}

resource "pingone_environment" "my_environment" {
  # ...
}