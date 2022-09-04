resource "pingone_certificate" "my_certificate" {
  environment_id = pingone_environment.my_environment.id

  usage_type        = "SSL/TLS"
  pkcs7_file_base64 = filebase64("../path/to/certificate.p7b")
}
