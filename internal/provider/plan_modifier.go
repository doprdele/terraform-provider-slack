package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

// useStateForUnknownModifier is a plan modifier that uses the state value for an
// attribute if the plan value is unknown.
func useStateForUnknownModifier() planmodifier.String {
	return &useStateForUnknownModifierImpl{}
}

type useStateForUnknownModifierImpl struct{}

func (m *useStateForUnknownModifierImpl) Description(ctx context.Context) string {
	return "Uses the state value for an attribute if the plan value is unknown."
}

func (m *useStateForUnknownModifierImpl) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m *useStateForUnknownModifierImpl) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if !req.PlanValue.IsUnknown() {
		return
	}
	if req.StateValue.IsNull() {
		return
	}
	resp.PlanValue = req.StateValue
}
