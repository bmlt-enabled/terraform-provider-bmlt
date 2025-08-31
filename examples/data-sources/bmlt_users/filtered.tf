# Query all users and filter locally
data "bmlt_users" "all" {}

# Find server administrators
locals {
  server_admins = [
    for user in data.bmlt_users.all.users :
    user if user.type == "serverAdmin"
  ]

  service_body_admins = [
    for user in data.bmlt_users.all.users :
    user if user.type == "serviceBodyAdmin"
  ]

  # Find specific user by username
  webservant = [
    for user in data.bmlt_users.all.users :
    user if user.username == "webservant"
  ]
}

# Output filtered results
output "admin_users" {
  description = "Server administrator users"
  value       = local.server_admins
}

output "service_body_admins" {
  description = "Service body administrator users"
  value       = local.service_body_admins
}

output "webservant_user" {
  description = "Webservant user details"
  value       = length(local.webservant) > 0 ? local.webservant[0] : null
}
