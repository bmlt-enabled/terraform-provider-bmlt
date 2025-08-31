# First create the users
resource "bmlt_user" "area_chair" {
  username     = "area.chair"
  password     = var.area_chair_password
  type         = "serviceBodyAdmin"
  display_name = "Area Chairperson"
  description  = "Chairperson for the Area Service Committee"
  email        = "chair@area.na.org"
}

resource "bmlt_user" "area_secretary" {
  username     = "area.secretary"
  password     = var.area_secretary_password
  type         = "serviceBodyAdmin"
  display_name = "Area Secretary"
  description  = "Secretary for the Area Service Committee"
  email        = "secretary@area.na.org"
}

# Then create the service body referencing the users
resource "bmlt_service_body" "area" {
  name          = "Central Area"
  description   = "Central Area Service Committee"
  type          = "AS"
  admin_user_id = bmlt_user.area_chair.id
  assigned_user_ids = [
    bmlt_user.area_chair.id,
    bmlt_user.area_secretary.id
  ]

  url   = "https://central.na.org"
  email = "webservant@central.na.org"
}

# Variables for sensitive passwords
variable "area_chair_password" {
  description = "Password for area chair user"
  type        = string
  sensitive   = true
}

variable "area_secretary_password" {
  description = "Password for area secretary user"
  type        = string
  sensitive   = true
}
