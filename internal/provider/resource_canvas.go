package provider

import (
	"context"
	"fmt"

	"github.com/doprdele/terraform-provider-slack-canvas/internal/slack"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource = &canvasResource{}
)



type canvasResource struct{
	client *slack.Client
}

type canvasResourceModel struct {
	ID        types.String `tfsdk:"id"`
	Content   types.String `tfsdk:"content"`
	ChannelID types.String `tfsdk:"channel_id"`
}

func (r *canvasResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_canvas"
}

func (r *canvasResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"content": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					useStateForUnknownModifier(),
				},
			},
			"channel_id": schema.StringAttribute{
				Required: true,
			},
		},
	}
}

func (r *canvasResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data canvasResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create the canvas.
	id, err := r.client.CreateCanvas(data.Content.ValueString(), data.ChannelID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create canvas, got error: %s", err))
		return
	}

	data.ID = types.StringValue(id)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *canvasResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data canvasResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// There is no API to read the content of a channel canvas. We can only verify it exists.
	// We assume that if the ReadCanvas call doesn't return an error, the canvas still exists.
	_, err := r.client.ReadCanvas(data.ID.ValueString())
	if err != nil {
		// If the canvas is not found, we remove it from the state.
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *canvasResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Not implemented.
}

func (r *canvasResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Not implemented.
}





func (r *canvasResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
