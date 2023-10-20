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
func stringChangeIfOtherString(attribute path.Path) planmodifier.String {
	return otherStringChanged{
		otherAttribute: attribute,
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
	var otherConfigValue, otherStateValue types.String

	resp.Diagnostics.Append(
		req.Config.GetAttribute(ctx, m.otherAttribute, &otherConfigValue)...)

	resp.Diagnostics.Append(
		req.State.GetAttribute(ctx, m.otherAttribute, &otherStateValue)...)

	// If the other attribute is changed, mark this attribute as unknown.
	if otherConfigValue != otherStateValue && !otherConfigValue.IsNull() {
		tflog.Debug(ctx, fmt.Sprintf(
			"otherStringChanged: %s otherPlanValue (%s) != otherStateValue (%s)",
			m.otherAttribute, otherConfigValue, otherStateValue))
		resp.PlanValue = types.StringUnknown()
	}
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

	resp.Diagnostics.Append(
		req.Config.GetAttribute(ctx, m.otherAttribute, &otherConfigValue)...)

	resp.Diagnostics.Append(
		req.State.GetAttribute(ctx, m.otherAttribute, &otherStateValue)...)

	resp.Diagnostics.Append(
		req.Plan.GetAttribute(ctx, m.otherAttribute, &otherPlanValue)...)

	var bodyPlanValue types.String
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("body"), &bodyPlanValue)...)

	// var bodyPlanValue types.String
	// resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("body"), &bodyPlanValue)...)

	// If the attribute isn't set, check the body front matter.
	if m.checkFrontmatter && otherConfigValue.IsNull() {
		tflog.Info(ctx, fmt.Sprintf("otherStringChanged: checking front matter for %s", m.otherField))
		value, diag := valueFromFrontMatter(ctx, bodyPlanValue.ValueString(), m.otherField)
		if diag != "" {
			resp.Diagnostics.AddError("Error parsing front matter.", diag)

			return
		}

		tflog.Info(ctx, bodyPlanValue.ValueString())

		tflog.Info(ctx, fmt.Sprintf("%s value: %v", m.otherAttribute, value))

		// If the value from frontmatter is not empty, compare it to the
		// current state.
		if value != (reflect.Value{}) {
			tflog.Info(ctx, fmt.Sprintf("%s was found in frontmatter with value %s", m.otherAttribute, value))
			isChanged = value.Interface().(string) != otherPlanValue.ValueString()
		} else {
			tflog.Info(ctx, fmt.Sprintf("%s was not found in frontmatter", m.otherAttribute))
		}

		// if value != (reflect.Value{}) {
		// 	tflog.Debug(ctx, fmt.Sprintf("%s: setting value from front matter", req.Path))
		// 	// resp.PlanValue = types.StringValue(value.Interface().(string))
		// 	resp.PlanValue = types.Int64Value(value.Interface().(int64))
		//
		// 	return
		// }
	} else {
		tflog.Info(ctx, "otherStringChanged: not checking front matter")
		isChanged = otherConfigValue != otherStateValue && !otherConfigValue.IsNull()
	}

	tflog.Info(ctx, fmt.Sprintf(
		"otherStringChanged: isChanged %v %s otherPlanValue (%s) otherConfigValue (%s) otherStateValue (%s)",
		isChanged, m.otherAttribute, otherPlanValue, otherConfigValue, otherStateValue))

	// If the other attribute is changed, mark this attribute as unknown.
	// if otherPlanValue == otherStateValue && !otherConfigValue.IsNull() { //&& !otherConfigValue.IsNull() {
	if isChanged {
		tflog.Info(ctx, fmt.Sprintf(
			"otherStringChanged: %s otherPlanValue (%s) != otherStateValue (%s)",
			m.otherAttribute, otherConfigValue, otherStateValue))
		resp.PlanValue = types.Int64Unknown()

		return
	}
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
