# Service body with multiple users assigned
resource "bmlt_service_body" "district" {
  name              = "District 42"
  description       = "District 42 Service Committee"
  type              = "AS"
  admin_user_id     = 5
  assigned_user_ids = [5, 6, 7, 8, 9] # Multiple users have access

  url      = "https://district42.org"
  helpline = "+1-555-HELP"
  email    = "contact@district42.org"
}
