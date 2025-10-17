package provider

import (
	"context"
	"fmt"

	"github.com/doprdele/terraform-provider-slack/internal/slack"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource = &userCanvasResource{}
)

type userCanvasResource struct {
	client *slack.Client
}

type userCanvasResourceModel struct {
	ID        types.String `tfsdk:"id"`
	Content   types.String `tfsdk:"content"`
	ChannelID types.String `tfsdk:"channel_id"`
	Private   types.Bool   `tfsdk:"private"`
	UserIDs   []string     `tfsdk:"user_ids"`
}

func (r *userCanvasResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user_canvas"
}

func (r *userCanvasResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"content": schema.StringAttribute{
				Required: true,
			},
			"channel_id": schema.StringAttribute{
				Optional: true,
			},
			"private": schema.BoolAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"user_ids": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
			},
		},
	}
}

func (r *userCanvasResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data userCanvasResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	id, err := r.client.CreateUserCanvas(data.Content.ValueString(), data.ChannelID.ValueString(), data.Private.ValueBool(), data.UserIDs)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create canvas, got error: %s", err))
		return
	}

	data.ID = types.StringValue(id)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *userCanvasResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data userCanvasResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// There is no API to read the content of a user canvas. We can only verify it exists.
	_, err := r.client.ReadCanvas(data.ID.ValueString())
	if err != nil {
		// If the canvas is not found, we remove it from the state.
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *userCanvasResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data userCanvasResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.UpdateUserCanvas(data.ID.ValueString(), data.Content.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update canvas, got error: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *userCanvasResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data userCanvasResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteUserCanvas(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete canvas, got error: %s", err))
		return
	}
}

func (r *userCanvasResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
