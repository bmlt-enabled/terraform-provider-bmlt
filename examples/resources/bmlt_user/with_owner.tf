# Create a parent admin user
resource "bmlt_user" "regional_admin" {
  username     = "regional.admin"
  password     = var.regional_admin_password
  type         = "serviceBodyAdmin"
  display_name = "Regional Administrator"
  description  = "Regional Service Office Administrator"
  email        = "rso@region.na.org"
}

# Create a child user owned by the regional admin
resource "bmlt_user" "area_admin" {
  username     = "area.admin"
  password     = var.area_admin_password
  type         = "serviceBodyAdmin"
  display_name = "Area Administrator"
  description  = "Area Service Committee Administrator"
  email        = "asc@area.na.org"
  owner_id     = bmlt_user.regional_admin.id # Owned by regional admin
}

variable "regional_admin_password" {
  description = "Password for regional admin user"
  type        = string
  sensitive   = true
}

variable "area_admin_password" {
  description = "Password for area admin user"
  type        = string
  sensitive   = true
}
