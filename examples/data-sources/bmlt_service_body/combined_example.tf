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

# Example 1: Look up a service body by ID (efficient - direct API call)
data "bmlt_service_body" "camna_by_id" {
  service_body_id = 1046
}

# Example 2: Look up a service body by name (requires fetching all service bodies to filter)
data "bmlt_service_body" "camna_by_name" {
  name = "Central Allegheny Mountain Area"
}

# Example 3: Use a variable to dynamically choose the lookup method
variable "service_body_lookup" {
  description = "Service body to look up"
  type = object({
    id   = optional(number)
    name = optional(string)
  })
  
  # Example values - provide only one:
  default = {
    name = "Central Allegheny Mountain Area"  # Use name lookup
    # id = 1046                               # Or use ID lookup (comment out name)
  }
}

data "bmlt_service_body" "dynamic" {
  service_body_id = var.service_body_lookup.id
  name            = var.service_body_lookup.name
}

# Output the results
output "service_body_details" {
  value = {
    by_id = {
      id          = data.bmlt_service_body.camna_by_id.id
      name        = data.bmlt_service_body.camna_by_id.name
      description = data.bmlt_service_body.camna_by_id.description
      type        = data.bmlt_service_body.camna_by_id.type
      parent_id   = data.bmlt_service_body.camna_by_id.parent_id
    }
    by_name = {
      id          = data.bmlt_service_body.camna_by_name.id
      name        = data.bmlt_service_body.camna_by_name.name
      description = data.bmlt_service_body.camna_by_name.description
      type        = data.bmlt_service_body.camna_by_name.type
      parent_id   = data.bmlt_service_body.camna_by_name.parent_id
    }
    dynamic = {
      id          = data.bmlt_service_body.dynamic.id
      name        = data.bmlt_service_body.dynamic.name
      description = data.bmlt_service_body.dynamic.description
      type        = data.bmlt_service_body.dynamic.type
      parent_id   = data.bmlt_service_body.dynamic.parent_id
    }
  }
  
  description = "Service body details retrieved using different lookup methods"
}

# Verify both methods return the same service body (they should!)
output "service_bodies_match" {
  value = {
    ids_match   = data.bmlt_service_body.camna_by_id.id == data.bmlt_service_body.camna_by_name.id
    names_match = data.bmlt_service_body.camna_by_id.name == data.bmlt_service_body.camna_by_name.name
  }
}
