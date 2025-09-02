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

# Example: Look up a service body by name
data "bmlt_service_body" "by_name" {
  name = "Central Allegheny Mountain Area"
}

output "service_body_by_name" {
  value = {
    id               = data.bmlt_service_body.by_name.id
    name             = data.bmlt_service_body.by_name.name
    description      = data.bmlt_service_body.by_name.description
    type             = data.bmlt_service_body.by_name.type
    parent_id        = data.bmlt_service_body.by_name.parent_id
    admin_user_id    = data.bmlt_service_body.by_name.admin_user_id
    assigned_user_ids = data.bmlt_service_body.by_name.assigned_user_ids
    url              = data.bmlt_service_body.by_name.url
    helpline         = data.bmlt_service_body.by_name.helpline
    email            = data.bmlt_service_body.by_name.email
    world_id         = data.bmlt_service_body.by_name.world_id
  }
}
