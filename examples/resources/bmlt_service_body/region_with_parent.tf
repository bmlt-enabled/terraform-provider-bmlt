# Create a parent area first
resource "bmlt_service_body" "metro_area" {
  name              = "Metro Area"
  description       = "Metropolitan Area Service Committee"
  type              = "AS" # Area Service
  admin_user_id     = 1
  assigned_user_ids = [1]
}

# Create a region under the area
resource "bmlt_service_body" "north_region" {
  name              = "North Region"
  description       = "North Regional Service Committee"
  type              = "RS" # Regional Service
  parent_id         = bmlt_service_body.metro_area.id
  admin_user_id     = 2
  assigned_user_ids = [2, 3]
}
