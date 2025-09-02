# Terraform Provider for BMLT

[![Go Report Card](https://goreportcard.com/badge/github.com/bmlt-enabled/terraform-provider-bmlt)](https://goreportcard.com/report/github.com/bmlt-enabled/terraform-provider-bmlt)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![GitHub release](https://img.shields.io/github/release/bmlt-enabled/terraform-provider-bmlt.svg)](https://github.com/bmlt-enabled/terraform-provider-bmlt/releases)

A Terraform provider for managing [BMLT (Basic Meeting List Toolbox)](https://bmlt.app) resources through Infrastructure as Code.

## Features

This provider allows you to manage BMLT server resources including:

- **Meetings** - Create, update, and manage 12 step meetings with location, time, and format assignments
- **Formats** - Manage meeting formats (e.g., "Open", "Closed") with multi-language translations
- **Service Bodies** - Organize areas, regions, and other organizational units with user assignments
- **Users** - Manage BMLT server users with different permission levels

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- Go >= 1.24.4 (for development)
- Access to a BMLT server with API credentials

## Quick Start

### 1. Configure the Provider

Add the provider to your Terraform configuration:

```hcl
terraform {
  required_providers {
    bmlt = {
      source  = "bmlt-enabled/bmlt"
      version = "~> 1.0"
    }
  }
}

provider "bmlt" {
  host     = "https://your-bmlt-server.com/main_server"
  username = var.bmlt_username
  password = var.bmlt_password
}
```

### 2. Authentication Methods

#### Method 1: Username/Password (OAuth2 Flow)
The provider exchanges your credentials for an OAuth2 access token automatically.

**Using provider block:**
```hcl
provider "bmlt" {
  host     = "https://bmlt.example.com/main_server"
  username = "your_username"
  password = "your_password"
}
```

**Using environment variables (recommended for production):**
```bash
export BMLT_HOST="https://bmlt.example.com/main_server"
export BMLT_USERNAME="your_username"
export BMLT_PASSWORD="your_password"
```

#### Method 2: Access Token (Direct)
Use a pre-generated OAuth2 access token directly - ideal for CI/CD pipelines.

**Using provider block:**
```hcl
provider "bmlt" {
  host         = "https://bmlt.example.com/main_server"
  access_token = "your_oauth2_access_token"
}
```

**Using environment variables:**
```bash
export BMLT_HOST="https://bmlt.example.com/main_server"
export BMLT_ACCESS_TOKEN="your_oauth2_access_token"
```

**Generate Access Token:**
```bash
curl -X POST "https://your-server.com/main_server/api/v1/auth/token" \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "grant_type=password&username=your_username&password=your_password"
```

### 3. Basic Usage Example

```hcl
# Query existing formats with language filtering for easy access
data "bmlt_formats" "english" {
  language = "en"  # This populates the formats_by_key map
}

# Create a custom format
resource "bmlt_format" "custom_format" {
  world_id = "CUSTOM"
  type     = "FC3"

  translations {
    key         = "en"
    name        = "Custom Format"
    description = "A custom meeting format"
    language    = "en"
  }
}

# Create a user
resource "bmlt_user" "service_admin" {
  username     = "terraform_admin"
  password     = "SecurePassword123!"
  type         = "serviceBodyAdmin"
  display_name = "Terraform Admin"
  email        = "admin@example.com"
}

# Create a service body
resource "bmlt_service_body" "area" {
  name              = "Terraform Area"
  description       = "An area managed by Terraform"
  type              = "AS"
  admin_user_id     = bmlt_user.service_admin.id
  assigned_user_ids = [bmlt_user.service_admin.id]
}

# Create a meeting using format keys (much simpler than before!)
resource "bmlt_meeting" "example" {
  service_body_id = bmlt_service_body.area.id
  format_ids = [
    data.bmlt_formats.english.formats_by_key["O"].id,   # Open
    data.bmlt_formats.english.formats_by_key["D"].id,   # Discussion  
    data.bmlt_formats.english.formats_by_key["BT"].id,  # Basic Text
  ]
  venue_type      = 1 # In-person
  day             = 1 # Monday
  start_time      = "19:00"
  duration        = "01:30"
  latitude        = 40.7128
  longitude       = -74.0060
  published       = true
  name            = "Example Meeting"
  
  location_text         = "Community Center"
  location_street       = "123 Main St"
  location_municipality = "Anytown"
  location_province     = "NY"
  location_postal_code_1 = "12345"
}
```

## Resources

### `bmlt_format`
Manages meeting formats with multi-language translations.

```hcl
resource "bmlt_format" "example" {
  world_id = "EXAMPLE"
  type     = "FC3"

  translations {
    key         = "en"
    name        = "Example Format"
    description = "An example meeting format"
    language    = "en"
  }

  translations {
    key         = "es"
    name        = "Formato de Ejemplo"
    description = "Un formato de reunión de ejemplo"
    language    = "es"
  }
}
```

### `bmlt_user`
Manages BMLT server users with different permission levels.

```hcl
resource "bmlt_user" "example" {
  username     = "example_user"
  password     = "SecurePassword123!"
  type         = "serviceBodyAdmin"
  display_name = "Example User"
  description  = "An example user account"
  email        = "user@example.com"
  owner_id     = 1 # Optional: ID of the user who owns this account
}
```

### `bmlt_service_body`
Manages organizational units like areas and regions.

```hcl
resource "bmlt_service_body" "example" {
  name              = "Example Area"
  description       = "An example area service body"
  type              = "AS" # Area Service
  admin_user_id     = bmlt_user.example.id
  assigned_user_ids = [bmlt_user.example.id]
  email             = "area@example.com"
  url               = "https://example.com/area"
  helpline          = "555-HELP"
  parent_id         = 1 # Optional: Parent service body ID
}
```

### `bmlt_meeting`
Manages NA/AA meetings with comprehensive location and contact information.

```hcl
resource "bmlt_meeting" "example" {
  service_body_id        = bmlt_service_body.example.id
  format_ids             = [1, 2, 3]
  venue_type             = 1 # 1=in-person, 2=virtual, 3=hybrid
  temporarily_virtual    = false
  day                    = 1 # 0=Sunday, 1=Monday, etc.
  start_time             = "19:00"
  duration               = "01:30"
  time_zone              = "America/New_York"
  latitude               = 40.7128
  longitude              = -74.0060
  published              = true
  name                   = "Example Meeting"
  
  # Location details
  location_text          = "Community Center Room A"
  location_info          = "Enter through main entrance"
  location_street        = "123 Main Street"
  location_municipality  = "Anytown"
  location_province      = "NY"
  location_postal_code_1 = "12345"
  location_nation        = "USA"
  
  # Virtual meeting details (if applicable)
  virtual_meeting_link   = "https://zoom.us/j/123456789"
  
  # Contact information
  contact_name_1         = "John Doe"
  contact_phone_1        = "555-0123"
  contact_email_1        = "john@example.com"
  
  # Additional information
  comments               = "Wheelchair accessible"
}
```

## Data Sources

### `bmlt_formats`
Retrieve information about all available meeting formats.

**Basic usage (all formats):**
```hcl
data "bmlt_formats" "all" {}

output "format_count" {
  value = length(data.bmlt_formats.all.formats)
}
```

**Language-filtered formats (recommended):**
```hcl
# Get formats with language filtering and convenient key-based access
data "bmlt_formats" "english" {
  language = "en"  # Populates formats_by_key map
}

# Easy access to specific formats by their key
resource "bmlt_meeting" "example" {
  format_ids = [
    data.bmlt_formats.english.formats_by_key["O"].id,   # Open
    data.bmlt_formats.english.formats_by_key["D"].id,   # Discussion
    data.bmlt_formats.english.formats_by_key["WC"].id,  # Wheelchair Accessible
  ]
  # ... other meeting attributes
}

# Get all available format keys
output "available_format_keys" {
  value = keys(data.bmlt_formats.english.formats_by_key)
}
```

### `bmlt_meetings`
Query meetings with optional filtering.

```hcl
data "bmlt_meetings" "monday_meetings" {
  days = "1" # Monday only
}

data "bmlt_meetings" "area_meetings" {
  service_body_ids = "5,6,7"
}
```

### `bmlt_service_bodies`
Retrieve information about all service bodies.

```hcl
data "bmlt_service_bodies" "all" {}
```

### `bmlt_users`
Retrieve information about all users.

```hcl
data "bmlt_users" "all" {}
```

## Import Support

All resources support Terraform import using their numeric ID:

```bash
terraform import bmlt_format.example 123
terraform import bmlt_meeting.example 456
terraform import bmlt_service_body.example 789
terraform import bmlt_user.example 101
```

## Development

### Prerequisites

- Go 1.24.4+
- Make
- golangci-lint (for linting)

### Building

```bash
# Quick build
go build -v .

# Build with version info and install locally
make build

# Install for local testing
make install
```

### Testing

```bash
# Unit tests
go test ./...

# Unit tests with coverage
go test -race -coverprofile=coverage.out -covermode=atomic ./...

# Acceptance tests (requires BMLT server credentials)
export BMLT_HOST="https://your-server.com/main_server"
export BMLT_USERNAME="your-username"
export BMLT_PASSWORD="your-password"
make testacc
```

### Code Quality

```bash
# Lint code
make lint

# Format code
make fmt

# Generate documentation
make docs
```

### Local Development Testing

```bash
# Build and install provider locally
go build -o terraform-provider-bmlt
mkdir -p ~/.terraform.d/plugins/registry.terraform.io/bmlt-enabled/bmlt/99.0.0/darwin_amd64
cp terraform-provider-bmlt ~/.terraform.d/plugins/registry.terraform.io/bmlt-enabled/bmlt/99.0.0/darwin_amd64/

# Test with example configuration
cd examples/complete-example
terraform init
terraform plan
```

## Documentation

Complete documentation is available in the [`docs/`](./docs/) directory:

- [Provider Configuration](./docs/index.md)
- Resources:
  - [`bmlt_format`](./docs/resources/format.md)
  - [`bmlt_meeting`](./docs/resources/meeting.md) 
  - [`bmlt_service_body`](./docs/resources/service_body.md)
  - [`bmlt_user`](./docs/resources/user.md)
- Data Sources:
  - [`bmlt_formats`](./docs/data-sources/formats.md)
  - [`bmlt_meetings`](./docs/data-sources/meetings.md)
  - [`bmlt_service_bodies`](./docs/data-sources/service_bodies.md)
  - [`bmlt_users`](./docs/data-sources/users.md)

## Contributing

Contributions are welcome! Please see our [development guide](./DEVELOPMENT.md) for details on:

- Setting up your development environment
- Running tests
- Code style and standards
- Submitting pull requests

## Release Process

This project uses GitHub Actions for CI/CD:

1. **Push/PR**: Runs tests, linting, security scans
2. **Git tag**: Triggers GoReleaser build for multiple platforms
3. **GitHub Release**: Automatically created with signed artifacts
4. **Terraform Registry**: Published after GitHub release

## Support

- **Issues**: Report bugs or request features via [GitHub Issues](https://github.com/bmlt-enabled/terraform-provider-bmlt/issues)
- **Discussions**: Join the conversation on [GitHub Discussions](https://github.com/bmlt-enabled/terraform-provider-bmlt/discussions)
- **BMLT Community**: Visit [bmlt.app](https://bmlt.app) for more information about BMLT

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Security

If you discover a security vulnerability, please email security@bmlt.app instead of creating a public issue.

---

*This provider is developed and maintained by the BMLT community. It is not affiliated with Terraform or HashiCorp.*
