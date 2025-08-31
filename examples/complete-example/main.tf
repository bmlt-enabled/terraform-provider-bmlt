terraform {
  required_providers {
    bmlt = {
      source = "bmlt-enabled/bmlt"
    }
  }
}

provider "bmlt" {
  host     = var.bmlt_host
  username = var.bmlt_username
  password = var.bmlt_password
}

# Variables
variable "bmlt_host" {
  description = "BMLT server host URL"
  type        = string
}

variable "bmlt_username" {
  description = "BMLT username"
  type        = string
  sensitive   = true
}

variable "bmlt_password" {
  description = "BMLT password"
  type        = string
  sensitive   = true
}

# Data sources to retrieve existing information
data "bmlt_formats" "all" {}

data "bmlt_service_bodies" "all" {}

# Create a custom format
resource "bmlt_format" "custom_format" {
  world_id = "CUSTOM_TF"
  type     = "FC3"

  translations {
    key         = "en"
    name        = "Terraform Managed"
    description = "A format managed by Terraform"
    language    = "en"
  }
}

# Create a user
resource "bmlt_user" "service_admin" {
  username     = "terraform_admin"
  password     = "SecurePassword123!"
  type         = "serviceBodyAdmin"
  display_name = "Terraform Service Admin"
  description  = "User created and managed by Terraform"
  email        = "terraform-admin@example.com"
}

# Create a service body
resource "bmlt_service_body" "terraform_area" {
  name              = "Terraform Managed Area"
  description       = "An area service body managed by Terraform"
  type              = "AS"
  admin_user_id     = bmlt_user.service_admin.id
  assigned_user_ids = [bmlt_user.service_admin.id]
  email             = "area@example.com"
  url               = "https://example.com/area"
}

# Create a meeting
resource "bmlt_meeting" "terraform_meeting" {
  service_body_id = bmlt_service_body.terraform_area.id
  format_ids      = [data.bmlt_formats.all.formats[0].id, bmlt_format.custom_format.id]
  venue_type      = 1 # In-person
  day             = 1 # Monday
  start_time      = "19:00"
  duration        = "01:30"
  time_zone       = "America/New_York"
  latitude        = 40.7128
  longitude       = -74.0060
  published       = true
  name            = "Terraform Managed Meeting"

  # Location details
  location_text         = "Community Center Room A"
  location_street       = "123 Terraform Ave"
  location_municipality = "Tech City"
  location_province     = "NY"
  location_postal_code_1 = "12345"
  location_nation       = "USA"

  # Contact information
  contact_name_1  = "Jane Smith"
  contact_phone_1 = "555-0123"
  contact_email_1 = "jane@example.com"

  # Additional info
  comments = "Managed by Terraform - please coordinate changes"
}

# Outputs
output "created_format_id" {
  description = "ID of the created format"
  value       = bmlt_format.custom_format.id
}

output "created_user_id" {
  description = "ID of the created user"
  value       = bmlt_user.service_admin.id
}

output "created_service_body_id" {
  description = "ID of the created service body"
  value       = bmlt_service_body.terraform_area.id
}

output "created_meeting_id" {
  description = "ID of the created meeting"
  value       = bmlt_meeting.terraform_meeting.id
}

output "total_formats" {
  description = "Total number of formats in the system"
  value       = length(data.bmlt_formats.all.formats)
}

output "total_service_bodies" {
  description = "Total number of service bodies in the system"
  value       = length(data.bmlt_service_bodies.all.service_bodies)
}
