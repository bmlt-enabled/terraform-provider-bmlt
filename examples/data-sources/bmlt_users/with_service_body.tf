# Query existing users to find IDs for service body assignment
data "bmlt_users" "all" {}

# Find specific users by username
locals {
  chair_user = [
    for user in data.bmlt_users.all.users :
    user if user.username == "area.chair"
  ]

  secretary_user = [
    for user in data.bmlt_users.all.users :
    user if user.username == "area.secretary"
  ]
}

# Use the found users in a service body configuration
resource "bmlt_service_body" "area" {
  name        = "Metro Area"
  description = "Metropolitan Area Service Committee"
  type        = "AS"

  # Use the data source to reference existing users
  admin_user_id = length(local.chair_user) > 0 ? local.chair_user[0].id : null
  assigned_user_ids = [
    for user in [local.chair_user, local.secretary_user] :
    user[0].id if length(user) > 0
  ]

  url   = "https://metroarea.na.org"
  email = "info@metroarea.na.org"
}

# Output user information for verification
output "selected_users" {
  description = "Users selected for service body assignment"
  value = {
    chair     = length(local.chair_user) > 0 ? local.chair_user[0] : null
    secretary = length(local.secretary_user) > 0 ? local.secretary_user[0] : null
  }
}
