package provider

import (
	"context"
	"fmt"

	"github.com/doprdele/terraform-provider-slack-canvas/internal/slack"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource = &userCanvasDataSource{}
)

type userCanvasDataSource struct {
	client *slack.Client
}

type userCanvasDataSourceModel struct {
	ID        types.String `tfsdk:"id"`
	Content   types.String `tfsdk:"content"`
	ChannelID types.String `tfsdk:"channel_id"`
	Private   types.Bool   `tfsdk:"private"`
	UserIDs   []string     `tfsdk:"user_ids"`
}

func (d *userCanvasDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user_canvas"
}

func (d *userCanvasDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required: true,
			},
			"content": schema.StringAttribute{
				Computed: true,
			},
			"channel_id": schema.StringAttribute{
				Computed: true,
			},
			"private": schema.BoolAttribute{
				Computed: true,
			},
			"user_ids": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
			},
		},
	}
}

func (d *userCanvasDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data userCanvasDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	fileInfo, err := d.client.GetFileInfo(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read file info, got error: %s", err))
		return
	}

	data.Content = types.StringValue("") // Content is not available from files.info
	data.ChannelID = types.StringValue(fileInfo.Group)
	data.Private = types.BoolValue(fileInfo.IsPublic)
	data.UserIDs = fileInfo.Shares.Private[fileInfo.Group]

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}