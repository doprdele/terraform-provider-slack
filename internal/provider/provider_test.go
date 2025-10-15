package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server instance.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"slack": providerserver.NewProtocol6WithError(New()),
}

func testAccPreCheck(t *testing.T) {
	// You can add code here to run prior to any test case execution, for example assertions
	// about the appropriate environment variables being set.
	if v := os.Getenv("SLACK_TOKEN"); v == "" {
		t.Fatal("SLACK_TOKEN must be set for acceptance tests")
	}
	if v := os.Getenv("SLACK_WORKSPACE"); v == "" {
		t.Fatal("SLACK_WORKSPACE must be set for acceptance tests")
	}
}
