package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCanvasResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccCanvasResourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("slack_canvas.test", "channel_id", "C07R5HJUBNX"),
					resource.TestCheckResourceAttr("slack_canvas.test", "content", "test content"),
				),
			},
		},
	})
}

func testAccCanvasResourceConfig() string {
	return `
terraform {
  required_providers {
    slack = {
      source = "doprdele/slack"
    }
  }
}

provider "slack" {
  slack_token     = "test-token"
  slack_workspace = "test-workspace"
}

resource "slack_canvas" "test" {
	content    = "test content"
	channel_id = "C07R5HJUBNX"
}
`
}
