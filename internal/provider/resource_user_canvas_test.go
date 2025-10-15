package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccUserCanvasResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccUserCanvasResourceConfig("test content"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("slack_user_canvas.test", "content", "test content"),
				),
			},
			// Update and Read testing
			{
				Config: testAccUserCanvasResourceConfig("new content"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("slack_user_canvas.test", "content", "new content"),
				),
			},
		},
	})
}

func testAccUserCanvasResourceConfig(content string) string {
	return `
provider "slack" {
  slack_token     = "test-token"
  slack_workspace = "test-workspace"
}

resource "slack_user_canvas" "test" {
	content    = "` + content + `"
}
`
}