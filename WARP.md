# WARP.md

This file provides guidance to WARP (warp.dev) when working with code in this repository.

## Project Overview

This is a **Terraform Provider for BMLT** (Basic Meeting List Toolbox), written in Go using the Terraform Plugin Framework v1.9.0. The provider manages BMLT server resources including meetings, formats, service bodies, and users through Infrastructure as Code.

**Key technologies:**
- Go 1.23+ (using `github.com/bmlt-enabled/bmlt-server-go-client` for API interactions)
- Terraform Plugin Framework
- OAuth2 password flow authentication
- GoReleaser for multi-platform builds

## Development Commands

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

# Run single test function
go test -v -run TestSpecificFunction ./internal/provider

# Acceptance tests (requires BMLT server credentials)
export BMLT_HOST="https://your-server.com/main_server"
export BMLT_USERNAME="your-username"  
export BMLT_PASSWORD="your-password"
make testacc
```

### Code Quality
```bash
# Lint code (requires golangci-lint)
make lint

# Format code 
make fmt

# Check formatting
make fmtcheck

# Static analysis
go vet ./...
```

### Documentation
```bash
# Generate provider documentation
make docs
# This runs: go generate ./...
```

### Local Development Testing
```bash
# Build and install provider locally for testing
go build -o terraform-provider-bmlt
mkdir -p ~/.terraform.d/plugins/registry.terraform.io/bmlt-enabled/bmlt/99.0.0/darwin_amd64
cp terraform-provider-bmlt ~/.terraform.d/plugins/registry.terraform.io/bmlt-enabled/bmlt/99.0.0/darwin_amd64/

# Test with example configuration in ./t/ directory
cd t && terraform init && terraform plan
```

## Code Architecture

### High-Level Structure

```
.
├── main.go                           # Provider entry point
├── internal/provider/                # Core provider implementation
│   ├── provider.go                   # Provider configuration & client setup
│   ├── *_resource.go                 # Resource implementations (CRUD operations)
│   └── *_data_source.go              # Data source implementations (read-only)
├── examples/                         # Terraform configuration examples
├── .github/workflows/                # CI/CD automation
└── scripts/                          # Development scripts
```

### Provider Architecture Pattern

The provider follows the **Terraform Plugin Framework** pattern:

1. **Provider Configuration** (`provider.go`):
   - Handles OAuth2 authentication using username/password flow
   - Creates authenticated `BMTLClientData` struct containing:
     - `*bmlt.APIClient` - Generated Go client for BMLT API
     - `context.Context` - OAuth2-authenticated context
   - Supports configuration via provider block or environment variables

2. **Resource Pattern** (e.g., `format_resource.go`):
   - Each resource implements `resource.Resource` interface
   - Standard CRUD operations: `Create()`, `Read()`, `Update()`, `Delete()`
   - Schema definition using Plugin Framework types
   - Model structs map Terraform state to Go structs
   - Import support for existing resources

3. **Data Source Pattern** (e.g., `formats_data_source.go`):
   - Read-only access to BMLT resources
   - Implements `datasource.DataSource` interface
   - Used for querying existing resources (formats, meetings, etc.)

### Key Architectural Decisions

- **Authentication**: OAuth2 password flow with automatic token refresh
- **API Client**: Uses generated Go client from `github.com/bmlt-enabled/bmlt-server-go-client`
- **State Management**: Plugin Framework handles Terraform state automatically
- **Error Handling**: Comprehensive HTTP status code handling (401, 403, 404, 422, 500)
- **Versioning**: Git tag-based releases with GoReleaser for cross-platform builds

### Resources Managed

- **bmlt_format**: Meeting formats (e.g., "Open", "Closed") with translations
- **bmlt_meeting**: NA/AA meetings with location, time, format assignments
- **bmlt_service_body**: Organizational units (Areas, Regions) with user assignments
- **bmlt_user**: BMLT server users with different permission levels

### Data Sources Available

- **bmlt_formats**: Query all available formats
- **bmlt_meetings**: Query meetings with optional filtering (service body, day, etc.)  
- **bmlt_service_bodies**: Query all service bodies
- **bmlt_users**: Query all users

## Configuration Requirements

The provider supports two authentication methods:

### Method 1: Username/Password (OAuth2 Flow)
The provider exchanges your credentials for an OAuth2 access token automatically.

1. **Provider block** (development):
```hcl
provider "bmlt" {
  host     = "https://bmlt.example.com/main_server"
  username = "your_username"
  password = "your_password"
}
```

2. **Environment variables** (production):
```bash
export BMLT_HOST="https://bmlt.example.com/main_server"
export BMLT_USERNAME="your_username"
export BMLT_PASSWORD="your_password"
```

### Method 2: Access Token (Direct)
Use a pre-generated OAuth2 access token directly - ideal for CI/CD pipelines.

1. **Provider block**:
```hcl
provider "bmlt" {
  host         = "https://bmlt.example.com/main_server"
  access_token = "your_oauth2_access_token"
}
```

2. **Environment variable**:
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

## Release Process

The project uses **GitHub Actions** for CI/CD:

1. **Push/PR**: Runs tests, linting, security scans
2. **Git tag**: Triggers GoReleaser build for multiple platforms
3. **GitHub Release**: Automatically created with signed artifacts
4. **Terraform Registry**: Manual publishing after GitHub release

Required repository secrets:
- `GPG_PRIVATE_KEY`: For signing release artifacts
- `PASSPHRASE`: GPG key passphrase

## Import Support

All resources support Terraform import using their numeric ID:

```bash
terraform import bmlt_format.example 123
terraform import bmlt_meeting.example 456
terraform import bmlt_service_body.example 789
terraform import bmlt_user.example 101
```
