package readme

import (
	"context"
	"fmt"
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// otherStringChanged is a plan modifier that plans a change for an
// attribute if another specified *string* attribute is changed.
type otherStringChanged struct {
	otherAttribute   path.Path
	otherField       string
	checkFrontmatter bool
}

// Custom plan modifier to flag a *string* attribute for change if another
// specified *string* attribute changes.
func stringChangeIfOtherString(
	attribute path.Path,
	otherField string,
	checkFrontmatter bool,
) planmodifier.String {
	return otherStringChanged{
		otherAttribute:   attribute,
		otherField:       otherField,
		checkFrontmatter: checkFrontmatter,
	}
}

// Custom plan modifier to flag an *int64* attribute for change if another
// specified *string* attribute changes.
func int64ChangeIfOtherString(
	attribute path.Path,
	otherField string,
	checkFrontmatter bool,
) planmodifier.Int64 {
	return otherStringChanged{
		otherAttribute:   attribute,
		otherField:       otherField,
		checkFrontmatter: checkFrontmatter,
	}
}

// Description returns a plain text description of the modifier's behavior.
func (m otherStringChanged) Description(ctx context.Context) string {
	return "If another attribute is changed, this attribute will be changed."
}

// MarkdownDescription returns a markdown formatted description of the
// modifier's behavior.
func (m otherStringChanged) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

// PlanModifyString implements a modifier for planning a change for an
// attribute if another specified *string* attribute changes.
func (m otherStringChanged) PlanModifyString(
	ctx context.Context,
	req planmodifier.StringRequest,
	resp *planmodifier.StringResponse,
) {
	var isChanged bool
	var otherConfigValue, otherPlanValue, otherStateValue types.String

	// The config is loaded to know whether to check frontmatter if the
	// attribute is not set.
	resp.Diagnostics.Append(
		req.Config.GetAttribute(ctx, m.otherAttribute, &otherConfigValue)...)

	// Load state to compare to config value (from attribute or frontmatter).
	resp.Diagnostics.Append(
		req.State.GetAttribute(ctx, m.otherAttribute, &otherStateValue)...)

	// Load plan to compare to config value (from attribute or frontmatter).
	resp.Diagnostics.Append(
		req.Plan.GetAttribute(ctx, m.otherAttribute, &otherPlanValue)...)

	var bodyPlanValue types.String
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("body"), &bodyPlanValue)...)

	// If the attribute isn't set, check the body front matter.
	if m.checkFrontmatter && otherConfigValue.IsNull() {
		value, diag := valueFromFrontMatter(ctx, bodyPlanValue.ValueString(), m.otherField)
		if diag != "" {
			resp.Diagnostics.AddError("Error parsing front matter.", diag)

			return
		}

		// If the value from frontmatter is not empty, compare it to the
		// current state.
		if value != (reflect.Value{}) {
			tflog.Debug(ctx, fmt.Sprintf(
				"%s was found in frontmatter with value %s",
				m.otherAttribute, value))

			// If the value from frontmatter is different from the current
			// plan, mark this attribute as changed.
			isChanged = value.Interface().(string) != otherPlanValue.ValueString()
		} else {
			tflog.Debug(ctx, fmt.Sprintf(
				"value for %s was not found in frontmatter",
				m.otherAttribute))
		}
	} else {
		// If the attribute is set, compare it to the current state and ignore
		// the frontmatter.
		tflog.Debug(ctx, "otherStringChanged: not checking front matter")
		isChanged = otherConfigValue != otherStateValue && !otherStateValue.IsNull()
	}

	// If the other attribute is changed, mark this attribute as unknown to
	// trigger a change.
	if isChanged {
		tflog.Debug(ctx, fmt.Sprintf(
			"otherStringChanged: %s otherConfigValue (%s) != otherStateValue (%s)",
			m.otherAttribute, otherConfigValue, otherStateValue))
		resp.PlanValue = types.StringUnknown()

		return
	}

	// If the other attribute is not changed, set this attribute to the
	// current plan value.
	resp.PlanValue = req.PlanValue
}

// PlanModifyInt64 implements a modifier for planning a change for an
// *int64* attribute if another specified *string* attribute changes.
// The string attribute is optionally cheked for a value in the frontmatter.
func (m otherStringChanged) PlanModifyInt64(
	ctx context.Context,
	req planmodifier.Int64Request,
	resp *planmodifier.Int64Response,
) {
	var isChanged bool
	var otherConfigValue, otherPlanValue, otherStateValue types.String

	// The config is loaded to know whether to check frontmatter if the
	// attribute is not set.
	resp.Diagnostics.Append(
		req.Config.GetAttribute(ctx, m.otherAttribute, &otherConfigValue)...)

	// Load state to compare to config value (from attribute or frontmatter).
	resp.Diagnostics.Append(
		req.State.GetAttribute(ctx, m.otherAttribute, &otherStateValue)...)

	// Load plan to compare to config value (from attribute or frontmatter).
	resp.Diagnostics.Append(
		req.Plan.GetAttribute(ctx, m.otherAttribute, &otherPlanValue)...)

	var bodyPlanValue types.String
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("body"), &bodyPlanValue)...)

	// If the attribute isn't set, check the body front matter.
	if m.checkFrontmatter && otherConfigValue.IsNull() {
		value, diag := valueFromFrontMatter(ctx, bodyPlanValue.ValueString(), m.otherField)
		if diag != "" {
			resp.Diagnostics.AddError("Error parsing front matter.", diag)

			return
		}

		// If the value from frontmatter is not empty, compare it to the
		// current state.
		if value != (reflect.Value{}) {
			tflog.Debug(ctx, fmt.Sprintf(
				"%s was found in frontmatter with value %s",
				m.otherAttribute, value))

			// If the value from frontmatter is different from the current
			// plan, mark this attribute as changed.
			isChanged = value.Interface().(string) != otherPlanValue.ValueString()
		} else {
			tflog.Debug(ctx, fmt.Sprintf(
				"value for %s was not found in frontmatter",
				m.otherAttribute))
		}
	} else {
		// If the attribute is set, compare it to the current state and ignore
		// the frontmatter.
		tflog.Debug(ctx, "otherStringChanged: not checking front matter")
		isChanged = otherConfigValue != otherStateValue && !otherStateValue.IsNull()
	}

	// If the other attribute is changed, mark this attribute as unknown to
	// trigger a change.
	if isChanged {
		tflog.Debug(ctx, fmt.Sprintf(
			"otherStringChanged: %s otherConfigValue (%s) != otherStateValue (%s)",
			m.otherAttribute, otherConfigValue, otherStateValue))
		resp.PlanValue = types.Int64Unknown()

		return
	}

	// If the other attribute is not changed, set this attribute to the
	// current plan value.
	resp.PlanValue = req.PlanValue
}

// -----------------------------------------------------------------------------
// otherInt64Changed is a plan modifier that plans a change for an
// attribute if another specified *int64* attribute is changed.
type otherInt64Changed struct {
	otherAttribute path.Path
}

// stringChangeIfOtherInt64 is a custom plan modifier to flag a *string*
// attribute for change if another specified *int64* attribute changes.
func stringChangeIfOtherInt64(attribute path.Path) planmodifier.String {
	return otherInt64Changed{
		otherAttribute: attribute,
	}
}

// int64ChangeIfOtherInt64 is a custom plan modifier to flag an *int64*
// attribute for change if another specified *int64* attribute changes.
func int64ChangeIfOtherInt64(attribute path.Path) planmodifier.Int64 {
	return otherInt64Changed{
		otherAttribute: attribute,
	}
}

// Description returns a plain text description of the modifier's behavior.
func (m otherInt64Changed) Description(ctx context.Context) string {
	return "If another attribute is changed, this attribute will be changed."
}

// MarkdownDescription returns a markdown formatted description of the
// modifier's behavior.
func (m otherInt64Changed) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m otherInt64Changed) PlanModifyString(
	ctx context.Context,
	req planmodifier.StringRequest,
	resp *planmodifier.StringResponse,
) {
	var otherConfigValue, otherStateValue types.Int64

	resp.Diagnostics.Append(
		req.Config.GetAttribute(ctx, m.otherAttribute, &otherConfigValue)...)

	resp.Diagnostics.Append(
		req.State.GetAttribute(ctx, m.otherAttribute, &otherStateValue)...)

	// If the other attribute is changed, mark this attribute as unknown.
	if otherConfigValue != otherStateValue && !otherConfigValue.IsNull() {
		tflog.Debug(ctx, fmt.Sprintf(
			"otherStringChanged: %s otherPlanValue (%d) != otherStateValue (%d)",
			m.otherAttribute, otherConfigValue, otherStateValue))
		resp.PlanValue = types.StringUnknown()
	}
}

func (m otherInt64Changed) PlanModifyInt64(
	ctx context.Context,
	req planmodifier.Int64Request,
	resp *planmodifier.Int64Response,
) {
	var otherConfigValue, otherStateValue types.Int64

	resp.Diagnostics.Append(
		req.Config.GetAttribute(ctx, m.otherAttribute, &otherConfigValue)...)

	resp.Diagnostics.Append(
		req.State.GetAttribute(ctx, m.otherAttribute, &otherStateValue)...)

	// If the other attribute is changed, mark this attribute as unknown.
	if otherConfigValue != otherStateValue && !otherConfigValue.IsNull() {
		tflog.Debug(ctx, fmt.Sprintf(
			"otherStringChanged: %s otherPlanValue (%d) != otherStateValue (%d)",
			m.otherAttribute, otherConfigValue, otherStateValue))
		resp.PlanValue = types.Int64Unknown()
	}
}

// -----------------------------------------------------------------------------
// otherBoolChanged is a plan modifier that plans a change for an
// attribute if another specified *int64* attribute is changed.
type otherBoolChanged struct {
	otherAttribute path.Path
}

// stringChangeIfOtherInt64 is a custom plan modifier to flag a *string*
// attribute for change if another specified *int64* attribute changes.
func stringChangeIfOtherBool(attribute path.Path) planmodifier.String {
	return otherBoolChanged{
		otherAttribute: attribute,
	}
}

// int64ChangeIfOtherInt64 is a custom plan modifier to flag an *int64*
// attribute for change if another specified *int64* attribute changes.
func int64ChangeIfOtherBool(attribute path.Path) planmodifier.Int64 {
	return otherBoolChanged{
		otherAttribute: attribute,
	}
}

// Description returns a plain text description of the modifier's behavior.
func (m otherBoolChanged) Description(ctx context.Context) string {
	return "If another bool attribute is changed, this attribute will be changed."
}

// MarkdownDescription returns a markdown formatted description of the
// modifier's behavior.
func (m otherBoolChanged) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m otherBoolChanged) PlanModifyString(
	ctx context.Context,
	req planmodifier.StringRequest,
	resp *planmodifier.StringResponse,
) {
	var otherConfigValue, otherStateValue types.Bool

	resp.Diagnostics.Append(
		req.Config.GetAttribute(ctx, m.otherAttribute, &otherConfigValue)...)

	resp.Diagnostics.Append(
		req.State.GetAttribute(ctx, m.otherAttribute, &otherStateValue)...)

	// If the other attribute is changed, mark this attribute as unknown.
	if otherConfigValue != otherStateValue && !otherConfigValue.IsNull() {
		tflog.Debug(ctx, fmt.Sprintf(
			"otherStringChanged: %s otherPlanValue (%v) != otherStateValue (%v)",
			m.otherAttribute, otherConfigValue, otherStateValue))
		resp.PlanValue = types.StringUnknown()
	}
}

func (m otherBoolChanged) PlanModifyInt64(
	ctx context.Context,
	req planmodifier.Int64Request,
	resp *planmodifier.Int64Response,
) {
	var otherConfigValue, otherStateValue types.Bool

	resp.Diagnostics.Append(
		req.Config.GetAttribute(ctx, m.otherAttribute, &otherConfigValue)...)

	resp.Diagnostics.Append(
		req.State.GetAttribute(ctx, m.otherAttribute, &otherStateValue)...)

	// If the other attribute is changed, mark this attribute as unknown.
	if otherConfigValue != otherStateValue && !otherConfigValue.IsNull() {
		tflog.Debug(ctx, fmt.Sprintf(
			"otherStringChanged: %s otherPlanValue (%v) != otherStateValue (%v)",
			m.otherAttribute, otherConfigValue, otherStateValue))
		resp.PlanValue = types.Int64Unknown()
	}
}
