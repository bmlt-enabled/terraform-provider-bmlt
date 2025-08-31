# User Management Workflow Example
# Query existing users first to understand current state

data "bmlt_users" "existing" {}

# Analyze current user structure
locals {
  current_users = {
    total_count    = length(data.bmlt_users.existing.users)
    server_admins  = length([for u in data.bmlt_users.existing.users : u if u.type == "serverAdmin"])
    service_admins = length([for u in data.bmlt_users.existing.users : u if u.type == "serviceBodyAdmin"])
    observers      = length([for u in data.bmlt_users.existing.users : u if u.type == "observer"])
    deactivated    = length([for u in data.bmlt_users.existing.users : u if u.type == "deactivated"])
  }

  # Check if specific users already exist
  required_users     = ["webservant", "backup.admin", "observer.user"]
  existing_usernames = [for u in data.bmlt_users.existing.users : u.username]
  missing_users = [
    for username in local.required_users :
    username if !contains(local.existing_usernames, username)
  ]
}

# Create missing users if needed
resource "bmlt_user" "webservant" {
  count = contains(local.missing_users, "webservant") ? 1 : 0

  username     = "webservant"
  password     = var.webservant_password
  type         = "serviceBodyAdmin"
  display_name = "Web Servant"
  description  = "Primary web servant account"
  email        = "webservant@example.na.org"
}

resource "bmlt_user" "backup_admin" {
  count = contains(local.missing_users, "backup.admin") ? 1 : 0

  username     = "backup.admin"
  password     = var.backup_admin_password
  type         = "serverAdmin"
  display_name = "Backup Administrator"
  description  = "Backup server administrator"
  email        = "backup@example.na.org"
}

resource "bmlt_user" "observer" {
  count = contains(local.missing_users, "observer.user") ? 1 : 0

  username     = "observer.user"
  password     = var.observer_password
  type         = "observer"
  display_name = "Observer User"
  description  = "Read-only observer account"
  email        = "observer@example.na.org"
}

# Output current state analysis
output "user_analysis" {
  description = "Current user state analysis"
  value = {
    current_state = local.current_users
    missing_users = local.missing_users
    users_to_create = {
      webservant   = contains(local.missing_users, "webservant")
      backup_admin = contains(local.missing_users, "backup.admin")
      observer     = contains(local.missing_users, "observer.user")
    }
  }
}

# Variables for new user passwords
variable "webservant_password" {
  description = "Password for webservant user"
  type        = string
  sensitive   = true
  default     = ""
}

variable "backup_admin_password" {
  description = "Password for backup admin user"
  type        = string
  sensitive   = true
  default     = ""
}

variable "observer_password" {
  description = "Password for observer user"
  type        = string
  sensitive   = true
  default     = ""
}
