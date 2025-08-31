# Complete service body configuration with all optional fields
resource "bmlt_service_body" "complete_area" {
  name              = "Greater Metro Area"
  description       = "Greater Metropolitan Area Service Committee covering multiple districts"
  type              = "AS"
  parent_id         = 1 # Optional parent service body
  admin_user_id     = 10
  assigned_user_ids = [10, 11, 12, 13]

  # Contact information
  url      = "https://greatermarea.na.org"
  helpline = "+1-800-266-2262"
  email    = "webservant@greatermarea.na.org"

  # World service identifier for integration
  world_id = "GMArea2023"
}
