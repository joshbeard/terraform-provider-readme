package readme

import (
	"context"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/liveoaklabs/readme-api-go-client/readme"
	"github.com/liveoaklabs/terraform-provider-readme/readme/frontmatter"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &changelogResource{}
	_ resource.ResourceWithConfigure   = &changelogResource{}
	_ resource.ResourceWithImportState = &changelogResource{}
)

// changelogResource is the data source implementation.
type changelogResource struct {
	client *readme.Client
}

// changelogResourceModel is the resource model used by the readme_changelog resource.
type changelogResourceModel struct {
	Algolia   types.Object `tfsdk:"algolia"`
	Body      types.String `tfsdk:"body"`
	BodyClean types.String `tfsdk:"body_clean"`
	CreatedAt types.String `tfsdk:"created_at"`
	HTML      types.String `tfsdk:"html"`
	Hidden    types.Bool   `tfsdk:"hidden"`
	ID        types.String `tfsdk:"id"`
	Metadata  types.Object `tfsdk:"metadata"`
	Revision  types.Int64  `tfsdk:"revision"`
	Slug      types.String `tfsdk:"slug"`
	Title     types.String `tfsdk:"title"`
	Type      types.String `tfsdk:"type"`
	UpdatedAt types.String `tfsdk:"updated_at"`
}

// changelogResourceMapToModel maps a readme.Changelog to a changelogResourceModel
// for use in the readme_custom_page resource.
func changelogResourceMapToModel(changelog readme.Changelog, plan changelogResourceModel) changelogResourceModel {
	return changelogResourceModel{
		Algolia:   docModelAlgoliaValue(changelog.Algolia),
		Body:      plan.Body,
		BodyClean: types.StringValue(changelog.Body),
		CreatedAt: types.StringValue(changelog.CreatedAt),
		HTML:      types.StringValue(changelog.HTML),
		Hidden:    types.BoolValue(changelog.Hidden),
		ID:        types.StringValue(changelog.ID),
		Metadata:  docModelMetadataValue(changelog.Metadata),
		Revision:  types.Int64Value(int64(changelog.Revision)),
		Slug:      types.StringValue(changelog.Slug),
		Title:     types.StringValue(changelog.Title),
		Type:      types.StringValue(changelog.Type),
		UpdatedAt: types.StringValue(changelog.UpdatedAt),
	}
}

// NewChangelogResource is a helper function to simplify the provider implementation.
func NewChangelogResource() resource.Resource {
	return &changelogResource{}
}

// Metadata returns the data source type name.
func (r *changelogResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_changelog"
}

// Configure adds the provider configured client to the data source.
func (r *changelogResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*readme.Client)
}

// ModifyPlan is used for modifying the plan before it is applied. In particular,
// this is used to normalize the body attribute and to update dynamic attributes.
func (r *changelogResource) ModifyPlan(
	ctx context.Context,
	req resource.ModifyPlanRequest,
	resp *resource.ModifyPlanResponse,
) {
	plan := &changelogResourceModel{}
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() || plan == nil {
		return
	}

	// Trim leading and trailing whitespace from the body.
	// The ReadMe API normalizes this, but we need to track the original value
	// provided by the user.
	// The 'body_clean' attribute is used to track the normalized value to
	// compare against the API response.
	body := strings.TrimSpace(plan.Body.ValueString())
	plan.BodyClean = types.StringValue(body)

	if plan.Hidden.IsNull() {
		plan.Hidden = types.BoolValue(true)
	}

	diags := resp.Plan.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)

	state := &changelogResourceModel{}
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if state == nil {
		return
	}

	// Several attributes are refreshed whenever the changelog is modified.
	// This may need to be added to if additional attributes are discovered to
	// be dynamic.
	if plan.BodyClean != state.BodyClean ||
		plan.Hidden != state.Hidden {
		tflog.Info(ctx, "Changelog body has changed. Refreshing dynamic attributes.")

		plan.Algolia = types.ObjectUnknown(map[string]attr.Type{
			"record_count":    types.Int64Type,
			"publish_pending": types.BoolType,
			"updated_at":      types.StringType,
		})
		plan.HTML = types.StringUnknown()
		plan.Revision = types.Int64Unknown()
		plan.UpdatedAt = types.StringUnknown()
		plan.Metadata = types.ObjectUnknown(map[string]attr.Type{
			"description": types.StringType,
			"image": types.ListType{
				ElemType: types.StringType,
			},
			"title": types.StringType,
		})
	}

	diags = resp.Plan.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// ValidateConfig is used for validating attribute values.
func (r changelogResource) ValidateConfig(
	ctx context.Context,
	req resource.ValidateConfigRequest,
	resp *resource.ValidateConfigResponse,
) {
	var data changelogResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if data.Title.IsNull() {
		// check front matter for 'title'.
		titleMatter, diag := frontmatter.GetValue(ctx, data.Body.ValueString(), "Title")
		if diag != "" {
			resp.Diagnostics.AddAttributeError(
				path.Root("title"),
				"Error checking front matter during validation.",
				diag,
			)

			return
		}

		// Fail if title is not set in front matter or the attribute.
		if titleMatter == (reflect.Value{}) {
			resp.Diagnostics.AddAttributeError(
				path.Root("title"),
				"Missing required attribute.",
				"'title' must be set using the attribute or in the body front matter.",
			)

			return
		}
	}

	if data.Type.IsNull() {
		// check front matter for 'type'.
		typeMatter, diag := frontmatter.GetValue(ctx, data.Body.ValueString(), "Type")
		if diag != "" {
			resp.Diagnostics.AddAttributeError(
				path.Root("type"),
				"Error checking front matter during validation.",
				diag,
			)

			return
		}

		// Fail if type is not set in front matter or the attribute.
		if typeMatter == (reflect.Value{}) {
			resp.Diagnostics.AddAttributeError(
				path.Root("type"),
				"Missing required attribute.",
				"'type' must be set using the attribute or in the body front matter.",
			)

			return
		}
	}
}

// Create creates the changelog and sets the initial Terraform state.
func (r *changelogResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan changelogResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	hidden := plan.Hidden.ValueBoolPointer()
	if hidden == nil {
		hidden = boolPoint(true)
	}

	params := readme.ChangelogParams{
		Title:  plan.Title.ValueString(),
		Body:   plan.Body.ValueString(),
		Hidden: hidden,
		Type:   plan.Type.ValueString(),
	}

	changelog, _, err := r.client.Changelog.Create(params)
	if err != nil {
		resp.Diagnostics.AddError("Unable to create changelog.", err.Error())

		return
	}

	// Get the changelog
	changelog, _, err = r.client.Changelog.Get(changelog.Slug)
	if err != nil {
		resp.Diagnostics.AddError("Unable to retrieve changelog.", err.Error())

		return
	}

	state := changelogResourceMapToModel(changelog, plan)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

// Read refreshes the Terraform state with the latest data.
func (r *changelogResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state changelogResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	changelog, _, err := r.client.Changelog.Get(state.Slug.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to retrieve changelog.", err.Error())

		return
	}

	state = changelogResourceMapToModel(changelog, state)

	diags := resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the changelog and sets the updated Terraform state on success.
func (r *changelogResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state changelogResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	hidden := plan.Hidden.ValueBoolPointer()
	if hidden == nil {
		hidden = boolPoint(true)
	}

	params := readme.ChangelogParams{
		Title:  plan.Title.ValueString(),
		Body:   plan.Body.ValueString(),
		Hidden: hidden,
		Type:   plan.Type.ValueString(),
	}

	_, _, err := r.client.Changelog.Update(state.Slug.ValueString(), params)
	if err != nil {
		resp.Diagnostics.AddError("Unable to update changelog.", err.Error())

		return
	}

	// Get the changelog
	changelog, _, err := r.client.Changelog.Get(state.Slug.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to retrieve changelog.", err.Error())

		return
	}

	state = changelogResourceMapToModel(changelog, plan)

	diags := resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

// Delete deletes the changelog and removes the Terraform state on success.
func (r *changelogResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state changelogResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, apiResponse, err := r.client.Changelog.Delete(state.Slug.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to delete changelog", clientError(err, apiResponse))
	}
}

// ImportState imports a changelog by its slug.
func (r *changelogResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	resource.ImportStatePassthroughID(ctx, path.Root("slug"), req, resp)
}

// Schema for the readme_changelog resource.
func (r *changelogResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// nolint:goconst
		Description: "Manage changelogs on ReadMe.com\n\n" +
			"Changelogs on ReadMe support setting some attributes using front matter. " +
			"Resource attributes take precedence over front matter attributes in the provider.\n\n" +
			"Refer to <https://docs.readme.com/main/docs/rdme> for more information about using front matter in " +
			"ReadMe docs and changelogs.\n\n" +
			"See <https://docs.readme.com/main/reference/createchangelog> for more information about this API endpoint.",
		Attributes: map[string]schema.Attribute{
			"algolia": schema.SingleNestedAttribute{
				Description: "Metadata about the Algolia search integration. " +
					"See <https://docs.readme.com/main/docs/search> for more information.",
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"publish_pending": schema.BoolAttribute{
						Computed: true,
					},
					"record_count": schema.Int64Attribute{
						Computed: true,
					},
					"updated_at": schema.StringAttribute{
						Computed: true,
					},
				},
			},
			"body": schema.StringAttribute{
				Description: "The body of the changelog. Optionally use front matter to set certain attributes. ",
				Required:    true,
			},
			"body_clean": schema.StringAttribute{
				Description: "The body of the changelog after normalization.",
				Computed:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "The date the changelog was created.",
				Computed:    true,
			},
			"hidden": schema.BoolAttribute{
				Description: "Whether the changelog is hidden. This can alternatively be set using the `hidden` front matter key.",
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(true),
				PlanModifiers: []planmodifier.Bool{
					frontmatter.GetBool("Hidden"),
				},
			},
			"html": schema.StringAttribute{
				Description: "The body source formatted in HTML.",
				Computed:    true,
			},
			"id": schema.StringAttribute{
				Description: "The ID of the changelog.",
				Computed:    true,
			},
			"metadata": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"description": schema.StringAttribute{
						Computed: true,
					},
					"image": schema.ListAttribute{
						Computed:    true,
						ElementType: types.StringType,
					},
					"title": schema.StringAttribute{
						Computed: true,
					},
				},
			},
			"title": schema.StringAttribute{
				Description: "__REQUIRED.__ The title of the changelog. This can alternatively be set using the `title` front matter key.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					frontmatter.GetString("Title"),
				},
			},
			"type": schema.StringAttribute{
				Description: "__REQUIRED.__ The type of changelog. This can alternatively be set using the `type` front matter key. " +
					"Valid values: added, fixed, improved, deprecated, removed",
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.String{
					frontmatter.GetString("Type"),
				},
			},
			"revision": schema.Int64Attribute{
				Description: "The revision of the changelog.",
				Computed:    true,
			},
			"slug": schema.StringAttribute{
				Description: "The slug of the changelog.",
				Computed:    true,
			},
			"updated_at": schema.StringAttribute{
				Description: "The date the changelog was last updated.",
				Computed:    true,
			},
		},
	}
}
