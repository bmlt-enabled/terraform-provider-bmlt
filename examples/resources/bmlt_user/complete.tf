# Complete user configuration with all available fields
resource "bmlt_user" "complete_user" {
  username     = "complete.user"
  password     = var.user_password
  type         = "serviceBodyAdmin"
  display_name = "Complete Example User"
  description  = "A user account demonstrating all available configuration options"
  email        = "complete.user@example.na.org"
  owner_id     = 1 # Owned by user with ID 1
}

variable "user_password" {
  description = "Password for the user account"
  type        = string
  sensitive   = true
}
