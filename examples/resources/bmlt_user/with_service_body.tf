# Create users first
resource "bmlt_user" "chair" {
  username     = "district.chair"
  password     = var.chair_password
  type         = "serviceBodyAdmin"
  display_name = "District Chair"
  description  = "District Service Committee Chairperson"
  email        = "chair@district.na.org"
}

resource "bmlt_user" "secretary" {
  username     = "district.secretary"
  password     = var.secretary_password
  type         = "serviceBodyAdmin"
  display_name = "District Secretary"
  description  = "District Service Committee Secretary"
  email        = "secretary@district.na.org"
}

resource "bmlt_user" "webservant" {
  username     = "district.web"
  password     = var.web_password
  type         = "serviceBodyAdmin"
  display_name = "District Web Servant"
  description  = "District Web Servant"
  email        = "web@district.na.org"
}

# Then create service body using the users
resource "bmlt_service_body" "district" {
  name          = "District 50"
  description   = "District 50 Service Committee"
  type          = "AS"
  admin_user_id = bmlt_user.chair.id
  assigned_user_ids = [
    bmlt_user.chair.id,
    bmlt_user.secretary.id,
    bmlt_user.webservant.id
  ]

  url      = "https://district50.na.org"
  helpline = "+1-555-RECOVER"
  email    = "info@district50.na.org"
}

variable "chair_password" {
  description = "Password for chair user"
  type        = string
  sensitive   = true
}

variable "secretary_password" {
  description = "Password for secretary user"
  type        = string
  sensitive   = true
}

variable "web_password" {
  description = "Password for web servant user"
  type        = string
  sensitive   = true
}
