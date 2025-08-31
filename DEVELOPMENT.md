# Development Guide

## Prerequisites

- Go 1.21 or later
- Terraform 1.0 or later
- GoReleaser (for releases)
- golangci-lint (for code quality)

## Development Workflow

### 1. Building the Provider

```bash
go build -v .
```

### 2. Testing

Run unit tests:
```bash
go test ./...
```

Run tests with coverage:
```bash
go test -race -coverprofile=coverage.out -covermode=atomic ./...
```

### 3. Code Quality

Run linting:
```bash
golangci-lint run
```

Check formatting:
```bash
gofmt -s -l .
```

Run static analysis:
```bash
go vet ./...
```

### 4. Generate Documentation

```bash
go generate ./...
```

This will:
- Format example Terraform files
- Generate provider documentation using `terraform-plugin-docs`

### 5. Acceptance Testing

To run acceptance tests against a real BMLT server:

```bash
export BMLT_HOST="https://your-bmlt-server.com/main_server"
export BMLT_USERNAME="your-username"
export BMLT_PASSWORD="your-password"
export TF_ACC=1

go test -v -cover ./internal/provider/
```

## Local Testing

### 1. Build and Install Provider Locally

```bash
# Build the provider
go build -o terraform-provider-bmlt

# Create local provider directory
mkdir -p ~/.terraform.d/plugins/registry.terraform.io/bmlt-enabled/bmlt/99.0.0/darwin_amd64

# Copy the binary
cp terraform-provider-bmlt ~/.terraform.d/plugins/registry.terraform.io/bmlt-enabled/bmlt/99.0.0/darwin_amd64/
```

### 2. Test Configuration

Create a `test.tf` file:

```hcl
terraform {
  required_providers {
    bmlt = {
      source  = "registry.terraform.io/bmlt-enabled/bmlt"
      version = "99.0.0"
    }
  }
}

provider "bmlt" {
  host     = "https://demo.na-bmlt.org/main_server"
  username = "your-username"
  password = "your-password"
}

data "bmlt_formats" "all" {}

output "format_count" {
  value = length(data.bmlt_formats.all.formats)
}
```

Run Terraform:
```bash
terraform init
terraform plan
terraform apply
```

## Release Process

### 1. Prepare Release

1. Update version in any relevant files
2. Ensure all tests pass
3. Update documentation if needed
4. Commit all changes

### 2. Create Release

```bash
# Create and push a tag
git tag v1.0.0
git push origin v1.0.0
```

The GitHub Actions release workflow will automatically:
- Build the provider for all supported platforms
- Sign the release artifacts with GPG
- Create a GitHub release
- Upload all artifacts

### 3. Publish to Terraform Registry

After the GitHub release is created:

1. Go to the [Terraform Registry](https://registry.terraform.io/publish/provider)
2. Sign in with your GitHub account
3. Select your provider repository
4. Follow the publishing process

## GitHub Actions Workflows

### Continuous Integration

- **Test**: Runs on every push and pull request
  - Unit tests on multiple Go versions
  - Acceptance tests
  - Coverage reporting

- **Code Quality**: Ensures code standards
  - Linting with golangci-lint
  - Format checking
  - Static analysis

- **Security**: Scans for vulnerabilities
  - Go vulnerability checking
  - Security issue detection
  - Dependency scanning

- **Documentation**: Keeps docs in sync
  - Generates provider documentation
  - Comments on PRs when docs are out of sync

### Release Management

- **Release**: Triggered by version tags
  - Cross-platform builds
  - GPG signing
  - GitHub release creation

## Required Secrets

For releases, configure these repository secrets:

- `GPG_PRIVATE_KEY`: Your GPG private key
- `PASSPHRASE`: GPG key passphrase

For acceptance tests (optional):
- `BMLT_HOST`: Test BMLT server URL
- `BMLT_USERNAME`: Test username
- `BMLT_PASSWORD`: Test password

## Provider Structure

```
.
├── .github/
│   └── workflows/          # GitHub Actions workflows
├── docs/                   # Generated documentation
├── examples/               # Example Terraform configurations
├── internal/
│   └── provider/          # Provider implementation
├── templates/              # Documentation templates
├── .golangci.yml          # Linting configuration
├── .goreleaser.yml        # Release configuration
└── terraform-registry-manifest.json
```

## Best Practices

1. **Write Tests**: Always add tests for new features
2. **Documentation**: Keep examples and docs updated
3. **Code Quality**: Use the provided linting configuration
4. **Security**: Run security scans regularly
5. **Versioning**: Follow semantic versioning
6. **Backwards Compatibility**: Avoid breaking changes in minor versions

## Troubleshooting

### Build Issues

```bash
# Clean module cache
go clean -modcache

# Download dependencies
go mod download

# Verify dependencies
go mod verify
```

### Test Issues

```bash
# Run tests with verbose output
go test -v ./...

# Run specific test
go test -v -run TestSpecificFunction ./internal/provider
```

### GoReleaser Issues

```bash
# Check configuration
goreleaser check

# Test build locally (without releasing)
goreleaser build --snapshot --rm-dist
```
