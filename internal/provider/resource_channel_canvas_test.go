package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccChannelCanvasResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccChannelCanvasResourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("slack_channel_canvas.test", "channel_id", "C07R5HJUBNX"),
					resource.TestCheckResourceAttr("slack_channel_canvas.test", "content", "test content"),
				),
			},
		},
	})
}

func testAccChannelCanvasResourceConfig() string {
	return `


provider "slack" {
  slack_token     = "test-token"
  slack_workspace = "test-workspace"
}

resource "slack_channel_canvas" "test" {
	content    = "test content"
	channel_id = "C07R5HJUBNX"
}
`
}


