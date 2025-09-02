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

# Example: Look up a user by username
data "bmlt_user" "by_username" {
  username = "admin"
}

output "user_by_username" {
  value = {
    id           = data.bmlt_user.by_username.id
    username     = data.bmlt_user.by_username.username
    display_name = data.bmlt_user.by_username.display_name
    type         = data.bmlt_user.by_username.type
    email        = data.bmlt_user.by_username.email
    description  = data.bmlt_user.by_username.description
    owner_id     = data.bmlt_user.by_username.owner_id
  }
}
