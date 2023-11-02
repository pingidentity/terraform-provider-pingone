resource "pingone_certificate" "my_certificate" {
  environment_id = pingone_environment.my_environment.id

  usage_type = "SSL/TLS"
  pem_file   = file("../path/to/certificate.pem")
}
