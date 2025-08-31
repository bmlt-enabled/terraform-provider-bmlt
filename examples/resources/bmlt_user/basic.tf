# Create a basic service body admin user
resource "bmlt_user" "webservant" {
  username     = "webservant"
  password     = var.webservant_password
  type         = "serviceBodyAdmin"
  display_name = "Web Servant"
}

variable "webservant_password" {
  description = "Password for webservant user"
  type        = string
  sensitive   = true
}
