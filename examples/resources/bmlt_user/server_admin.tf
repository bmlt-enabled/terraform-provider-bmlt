# Create a server administrator with full privileges
resource "bmlt_user" "server_admin" {
  username     = "server.admin"
  password     = var.server_admin_password
  type         = "serverAdmin"
  display_name = "Server Administrator"
  description  = "Main server administrator account"
  email        = "admin@bmlt.example.com"
}

variable "server_admin_password" {
  description = "Password for server admin user"
  type        = string
  sensitive   = true
}
