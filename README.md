# Terraform Provider for BMLT

This Terraform provider allows you to manage BMLT (Basic Meeting List Toolbox) resources using Infrastructure as Code.

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.23

## Building The Provider

1. Clone the repository
2. Enter the provider directory: `cd terraform-provider-bmlt`
3. Build the provider using the Go `install` command:

```shell
go install
```

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```shell
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

## Using the Provider

First, build and install the provider.

```shell
make install
```

Then, run the following command to initialize the workspace and apply the sample configuration.

```shell
terraform init && terraform apply
```

## Provider Configuration

The provider can be configured in three ways:

1. **Provider block configuration** (recommended for development):
```hcl
provider "bmlt" {
  host     = "https://bmlt.example.com/main_server"
  username = "your_username"
  password = "your_password"
}
```

2. **Environment variables** (recommended for production):
```bash
export BMLT_HOST="https://bmlt.example.com/main_server"
export BMLT_USERNAME="your_username"
export BMLT_PASSWORD="your_password"
```

3. **Terraform variables**:
```hcl
provider "bmlt" {
  host     = var.bmlt_host
  username = var.bmlt_username
  password = var.bmlt_password
}
```

## Available Resources

### bmlt_format
Manages meeting formats.

```hcl
resource "bmlt_format" "example" {
  world_id = "CUSTOM_FORMAT"
  type     = "FC3"

  translations {
    key         = "en"
    name        = "Custom Format"
    description = "This is a custom format for our region"
    language    = "en"
  }
}
```

### bmlt_meeting
Manages meetings.

```hcl
resource "bmlt_meeting" "example" {
  service_body_id = 1
  format_ids      = [1, 2, 3]
  venue_type      = 1 # In-person
  day             = 1 # Monday
  start_time      = "19:00"
  duration        = "01:30"
  latitude        = 40.7128
  longitude       = -74.0060
  published       = true
  name            = "Monday Night Group"
}
```

### bmlt_service_body
Manages service bodies.

```hcl
resource "bmlt_service_body" "example" {
  name            = "Example Area"
  description     = "Example Area Service Body"
  type            = "AS"
  admin_user_id   = 1
  assigned_user_ids = [1, 2]
}
```

### bmlt_user
Manages users.

```hcl
resource "bmlt_user" "example" {
  username     = "example_user"
  password     = "secure_password"
  type         = "serviceBodyAdmin"
  display_name = "Example User"
  email        = "user@example.com"
}
```

## Available Data Sources

### bmlt_formats
Retrieves all available formats.

```hcl
data "bmlt_formats" "all" {}

output "format_count" {
  value = length(data.bmlt_formats.all.formats)
}
```

### bmlt_meetings
Retrieves meetings with optional filtering.

```hcl
# Get all meetings
data "bmlt_meetings" "all" {}

# Get meetings for specific service bodies
data "bmlt_meetings" "service_body_meetings" {
  service_body_ids = "1,2,3"
}

# Get weekend meetings
data "bmlt_meetings" "weekend_meetings" {
  days = "0,6" # Sunday and Saturday
}
```

### bmlt_service_bodies
Retrieves all service bodies.

```hcl
data "bmlt_service_bodies" "all" {}
```

### bmlt_users
Retrieves all users.

```hcl
data "bmlt_users" "all" {}
```

## Authentication

This provider uses OAuth2 password flow authentication. The provider will automatically:

1. Exchange your username/password for an access token
2. Handle token refresh automatically
3. Include the token in all API requests

## Error Handling

The provider includes comprehensive error handling for:

- Authentication errors (401)
- Authorization errors (403)
- Resource not found errors (404)
- Validation errors (422)
- Server errors (500)

## Import Support

All resources support Terraform import using their ID:

```shell
terraform import bmlt_format.example 123
terraform import bmlt_meeting.example 456
terraform import bmlt_service_body.example 789
terraform import bmlt_user.example 101
```

## Development

### Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.23

### Building

```shell
go build -o terraform-provider-bmlt
```

### Testing

```shell
make test
```

### Generating Documentation

```shell
go generate ./...
```

This will update the docs directory with the latest schema information.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for your changes
5. Run the test suite
6. Submit a pull request

## License

This provider is released under the MIT License. See LICENSE for details.
