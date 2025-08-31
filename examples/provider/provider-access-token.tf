terraform {
  required_providers {
    bmlt = {
      source = "bmlt-enabled/bmlt"
    }
  }
}

# Using access token authentication
provider "bmlt" {
  host         = "https://bmlt.example.com/main_server"
  access_token = var.bmlt_access_token
}

# Or using environment variables:
# export BMLT_HOST="https://bmlt.example.com/main_server"
# export BMLT_ACCESS_TOKEN="your_oauth2_access_token"

variable "bmlt_access_token" {
  description = "BMLT OAuth2 access token"
  type        = string
  sensitive   = true
}
