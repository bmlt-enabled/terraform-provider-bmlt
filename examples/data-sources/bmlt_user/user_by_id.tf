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

# Example: Look up a user by ID
data "bmlt_user" "by_id" {
  user_id = 123
}

output "user_by_id" {
  value = {
    id           = data.bmlt_user.by_id.id
    username     = data.bmlt_user.by_id.username
    display_name = data.bmlt_user.by_id.display_name
    type         = data.bmlt_user.by_id.type
    email        = data.bmlt_user.by_id.email
    description  = data.bmlt_user.by_id.description
    owner_id     = data.bmlt_user.by_id.owner_id
  }
}
