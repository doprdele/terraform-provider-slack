package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccUserCanvasDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccUserCanvasDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.slack_user_canvas.test", "content", ""),
				),
			},
		},
	})
}

func testAccUserCanvasDataSourceConfig() string {
	return `
resource "slack_user_canvas" "test" {
	content = "test"
}

data "slack_user_canvas" "test" {
	id = slack_user_canvas.test.id
}
`
}
