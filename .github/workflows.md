# GitHub Actions Workflows

This directory contains GitHub Actions workflows for the BMLT Terraform Provider.

## Workflows

### 🧪 `test.yml` - Testing
- **Triggers**: Push to `main`, Pull Requests
- **Purpose**: Runs unit tests and acceptance tests
- **Go Versions**: Tests against Go 1.21.x and 1.22.x
- **Features**:
  - Unit tests with race detection and coverage
  - Acceptance tests for Terraform provider
  - Coverage reporting to Codecov

### 🚀 `release.yml` - Release Management
- **Triggers**: Push tags matching `v*`
- **Purpose**: Builds and releases the provider using GoReleaser
- **Features**:
  - Cross-platform builds (Linux, macOS, Windows, FreeBSD)
  - Multiple architectures (amd64, 386, arm, arm64)
  - GPG signing of release artifacts
  - Automatic GitHub release creation

### 🔍 `code-quality.yml` - Code Quality
- **Triggers**: Push to `main`, Pull Requests
- **Purpose**: Ensures code quality and formatting
- **Features**:
  - golangci-lint for comprehensive linting
  - Format checking with gofmt
  - Static analysis with go vet

### 📚 `docs.yml` - Documentation
- **Triggers**: Pull Requests affecting provider code or documentation
- **Purpose**: Ensures provider documentation is up-to-date
- **Features**:
  - Automatically generates documentation
  - Comments on PRs when docs are out of sync
  - Fails CI if documentation needs updating

### 🔒 `security.yml` - Security Scanning
- **Triggers**: Push to `main`, Pull Requests, Weekly schedule
- **Purpose**: Scans for security vulnerabilities
- **Features**:
  - govulncheck for Go vulnerability scanning
  - gosec for security issue detection
  - Nancy for dependency vulnerability scanning

## Setup Requirements

### For Releases
To use the release workflow, you need to set up the following repository secrets:

1. **GPG_PRIVATE_KEY**: Your GPG private key for signing releases
2. **PASSPHRASE**: The passphrase for your GPG key

#### Setting up GPG for Releases

1. Generate a GPG key if you don't have one:
   ```bash
   gpg --full-generate-key
   ```

2. Export your private key:
   ```bash
   gpg --armor --export-secret-keys YOUR_KEY_ID
   ```

3. Add the exported key to GitHub repository secrets as `GPG_PRIVATE_KEY`
4. Add your passphrase to GitHub repository secrets as `PASSPHRASE`

### For Testing
If you want to run acceptance tests against a real BMLT server, add these secrets:
- `BMLT_HOST`: Your test BMLT server URL
- `BMLT_USERNAME`: Username for authentication
- `BMLT_PASSWORD`: Password for authentication

## Configuration Files

- **`.golangci.yml`**: Configuration for golangci-lint with comprehensive linting rules
- **`.goreleaser.yml`**: GoReleaser configuration for building and releasing
- **`terraform-registry-manifest.json`**: Terraform Registry manifest for provider metadata

## Usage

These workflows will automatically run based on their triggers. For releases:

1. Create and push a git tag:
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```

2. The release workflow will automatically:
   - Build the provider for all supported platforms
   - Sign the release artifacts
   - Create a GitHub release
   - Upload all artifacts to the release

The provider will then be available for download and can be published to the Terraform Registry if desired.
