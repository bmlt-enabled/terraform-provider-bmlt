# Create a basic Area Service Body
resource "bmlt_service_body" "metro_area" {
  name              = "Metro Area"
  description       = "Metropolitan Area Service Committee"
  type              = "AS" # Area Service
  admin_user_id     = 1
  assigned_user_ids = [1, 2]
}
