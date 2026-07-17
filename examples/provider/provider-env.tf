terraform {
  required_providers {
    pingone = {
      source  = "pingidentity/pingone"
      version = ">= 1.21, < 1.22"
    }
  }
}

provider "pingone" {
}

resource "pingone_environment" "my_environment" {
  # ...
}
