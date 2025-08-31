terraform {
  required_providers {
    bmlt = {
      source = "bmlt-enabled/bmlt"
    }
  }
}

# Method 1: Username/Password authentication
provider "bmlt" {
  host     = "https://bmlt.example.com/main_server"
  username = var.bmlt_username
  password = var.bmlt_password
}

# Method 2: Access Token authentication
# provider "bmlt" {
#   host         = "https://bmlt.example.com/main_server"
#   access_token = var.bmlt_access_token
# }

# Environment variable options:
# For username/password:
# export BMLT_HOST="https://bmlt.example.com/main_server"
# export BMLT_USERNAME="your_username"
# export BMLT_PASSWORD="your_password"

# For access token:
# export BMLT_HOST="https://bmlt.example.com/main_server"
# export BMLT_ACCESS_TOKEN="your_oauth2_access_token"

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

variable "bmlt_access_token" {
  description = "BMLT OAuth2 access token"
  type        = string
  sensitive   = true
}
