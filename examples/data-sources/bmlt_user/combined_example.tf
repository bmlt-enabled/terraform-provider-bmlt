terraform {
  required_providers {
    bmlt = {
      source = "bmlt-enabled/bmlt"
    }
  }
}

provider "bmlt" {
  host     = "https://bmlt.example.com/main_server"
  username = "your_username"
  password = "your_password"
}

# Example 1: Look up a user by ID (efficient - direct API call)
data "bmlt_user" "admin_by_id" {
  user_id = 1
}

# Example 2: Look up a user by username (requires fetching all users to filter)
data "bmlt_user" "admin_by_username" {
  username = "admin"
}

# Example 3: Use a variable to dynamically choose the lookup method
variable "user_lookup" {
  description = "User to look up"
  type = object({
    id       = optional(number)
    username = optional(string)
  })

  # Example values - provide only one:
  default = {
    username = "admin" # Use username lookup
    # id = 1            # Or use ID lookup (comment out username)
  }
}

data "bmlt_user" "dynamic" {
  user_id  = var.user_lookup.id
  username = var.user_lookup.username
}

# Output the results
output "user_details" {
  value = {
    by_id = {
      id           = data.bmlt_user.admin_by_id.id
      username     = data.bmlt_user.admin_by_id.username
      display_name = data.bmlt_user.admin_by_id.display_name
      type         = data.bmlt_user.admin_by_id.type
    }
    by_username = {
      id           = data.bmlt_user.admin_by_username.id
      username     = data.bmlt_user.admin_by_username.username
      display_name = data.bmlt_user.admin_by_username.display_name
      type         = data.bmlt_user.admin_by_username.type
    }
    dynamic = {
      id           = data.bmlt_user.dynamic.id
      username     = data.bmlt_user.dynamic.username
      display_name = data.bmlt_user.dynamic.display_name
      type         = data.bmlt_user.dynamic.type
    }
  }

  description = "User details retrieved using different lookup methods"
}
