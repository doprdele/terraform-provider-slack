# Terraform Provider for Slack Canvases

This Terraform provider allows you to manage Slack Canvases.

## Example Usage

```hcl
provider "slack" {
  slack_token     = "xoxb-your-token"
  slack_workspace = "your-workspace"
}

resource "slack_user_canvas" "example" {
  content = "# Hello, World!\n\nThis is a canvas managed by Terraform."
}

data "slack_user_canvas" "example" {
  id = slack_user_canvas.example.id
}
```

## Schema

### Optional

- `slack_token` (String, Sensitive) The Slack API token. It can also be set via the `SLACK_TOKEN` environment variable.
- `slack_workspace` (String) The Slack workspace name. It can also be set via the `SLACK_WORKSPACE` environment variable.

## Resources

- [slack_user_canvas](./docs/resources/user_canvas.md)

## Data Sources

- [slack_user_canvas](./docs/data-sources/user_canvas.md)

## Release Process

Releases are automated via the `Test and Release` GitHub Actions workflow. Successful pushes to `main` run tests, build provider binaries for Linux, macOS, and Windows (amd64/arm64), generate the `SHA256SUMS` manifest, and sign it with the OpenPGP key identified by the `GPG_FINGERPRINT` secret. Semantic-release creates a `v<major>.<minor>.<patch>` tag and GitHub release containing the required Terraform/OpenTofu assets.

## Author

- Evan Sarmiento
