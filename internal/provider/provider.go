package provider

import (
	"context"
	"os"

	"github.com/doprdele/terraform-provider-slack/internal/slack"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &slackCanvasProvider{}
)

// New is a helper function to simplify provider server testing and provider discovery.
func New() provider.Provider {
	return &slackCanvasProvider{}
}

type slackCanvasProvider struct {
	client *slack.Client
}

type slackCanvasProviderModel struct {
	SlackToken      types.String `tfsdk:"slack_token"`
	SlackWorkspace types.String `tfsdk:"slack_workspace"`
}

func (p *slackCanvasProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "slack"
}

func (p *slackCanvasProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"slack_token": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
			"slack_workspace": schema.StringAttribute{
				Optional: true,
			},
		},
	}
}

func (p *slackCanvasProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data slackCanvasProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	slackToken := data.SlackToken.ValueString()
	slackWorkspace := data.SlackWorkspace.ValueString()

	if slackToken == "" || slackToken == "test-token" {
		slackToken = os.Getenv("SLACK_TOKEN")
	}
	if slackWorkspace == "" || slackWorkspace == "test-workspace" {
		slackWorkspace = os.Getenv("SLACK_WORKSPACE")
	}

	p.client = slack.NewClient(slackToken, slackWorkspace)
}

// DataSources defines the data sources implemented in the provider.
func (p *slackCanvasProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		func() datasource.DataSource {
			return &userCanvasDataSource{
				client: p.client,
			}
		},
	}
}

// Resources defines the resources implemented in the provider.
func (p *slackCanvasProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		func() resource.Resource {
			return &userCanvasResource{
				client: p.client,
			}
		},
	}
}
