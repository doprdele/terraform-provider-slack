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

## Author

- Evan Sarmiento
