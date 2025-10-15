resource "slack_channel_canvas" "example" {
  channel_id = "C07R5HJUBNX"
  content    = "# Hello, World!\n\nThis is a canvas managed by Terraform."
}

