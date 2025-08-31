terraform {
  required_providers {
    bmlt = {
      source = "bmlt-enabled/bmlt"
    }
  }
}

provider "bmlt" {
  host     = "https://bmlt.example.com/main_server"
  username = var.bmlt_username
  password = var.bmlt_password
}

# Or using environment variables:
# export BMLT_HOST="https://bmlt.example.com/main_server"
# export BMLT_USERNAME="your_username"
# export BMLT_PASSWORD="your_password"

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
