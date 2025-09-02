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

# Example: Look up a service body by ID
data "bmlt_service_body" "by_id" {
  service_body_id = 1046
}

output "service_body_by_id" {
  value = {
    id               = data.bmlt_service_body.by_id.id
    name             = data.bmlt_service_body.by_id.name
    description      = data.bmlt_service_body.by_id.description
    type             = data.bmlt_service_body.by_id.type
    parent_id        = data.bmlt_service_body.by_id.parent_id
    admin_user_id    = data.bmlt_service_body.by_id.admin_user_id
    assigned_user_ids = data.bmlt_service_body.by_id.assigned_user_ids
    url              = data.bmlt_service_body.by_id.url
    helpline         = data.bmlt_service_body.by_id.helpline
    email            = data.bmlt_service_body.by_id.email
    world_id         = data.bmlt_service_body.by_id.world_id
  }
}
