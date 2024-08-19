package readme

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/liveoaklabs/readme-api-go-client/readme"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &apiSpecificationResource{}
	_ resource.ResourceWithConfigure   = &apiSpecificationResource{}
	_ resource.ResourceWithImportState = &apiSpecificationResource{}
)

// apiSpecificationResource is the data source implementation.
type apiSpecificationResource struct {
	client *readme.Client
}

// apiSpecificationResourceModel maps the struct from the ReadMe client library to Terraform attributes.
type apiSpecificationResourceModel struct {
	ID             types.String `tfsdk:"id"`
	Category       types.Object `tfsdk:"category"`
	DeleteCategory types.Bool   `tfsdk:"delete_category"`
	UUID           types.String `tfsdk:"uuid"`
	Definition     types.String `tfsdk:"definition"`
	DefinitionNorm types.String `tfsdk:"definition_normalized"`
	LastSynced     types.String `tfsdk:"last_synced"`
	Semver         types.String `tfsdk:"semver"`
	Source         types.String `tfsdk:"source"`
	Title          types.String `tfsdk:"title"`
	Type           types.String `tfsdk:"type"`
	Version        types.String `tfsdk:"version"`
}

// NewAPISpecificationResource is a helper function to simplify the provider implementation.
func NewAPISpecificationResource() resource.Resource {
	return &apiSpecificationResource{}
}

// Metadata returns the data source type name.
func (r *apiSpecificationResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_api_specification"
}

// Configure adds the provider configured client to the data source.
func (r *apiSpecificationResource) Configure(
	_ context.Context,
	req resource.ConfigureRequest,
	_ *resource.ConfigureResponse,
) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*readme.Client)
}

// jsonMatch compares two JSON strings for equality.
func jsonMatch(one, two string) (bool, error) {
	var oneRaw, twoRaw json.RawMessage
	if err := json.Unmarshal([]byte(one), &oneRaw); err != nil {
		return false, fmt.Errorf("error unmarshalling first item: %w", err)
	}
	if err := json.Unmarshal([]byte(two), &twoRaw); err != nil {
		return false, fmt.Errorf("error unmarshalling second item: %w", err)
	}

	return reflect.DeepEqual(oneRaw, twoRaw), nil
}

// specCategoryObject maps a readme.CategorySummary type to a generic ObjectValue and returns the ObjectValue for use
// with the Terraform resource schema.
func specCategoryObject(spec readme.APISpecification) basetypes.ObjectValue {
	object, _ := types.ObjectValue(
		map[string]attr.Type{
			"id":    types.StringType,
			"title": types.StringType,
			"slug":  types.StringType,
			"order": types.Int64Type,
			"type":  types.StringType,
		},
		map[string]attr.Value{
			"id":    types.StringValue(spec.Category.ID),
			"title": types.StringValue(spec.Category.Title),
			"slug":  types.StringValue(spec.Category.Slug),
			"order": types.Int64Value(int64(spec.Category.Order)),
			"type":  types.StringValue(spec.Category.Type),
		})

	return object
}

// Schema defines the API Specification resource attributes.
func (r *apiSpecificationResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Manages an API specification on ReadMe.com\n\n" +
			"The provider creates and updates API specifications by first uploading the definition to the " +
			"API registry and then creating or updating the API specification using the UUID returned from the " +
			"API registry. This is necessary for associating an API specification with its definition. Ensuring " +
			"the definition is created in the API registry is necessary for retrieving the " +
			"remote definition. This behavior is undocumented in the ReadMe API documentation but works the same way " +
			"the official ReadMe `rdme` CLI tool works.\n\n" +
			"## External Changes\n\n" +
			"External changes made to an API specification managed by Terraform will not be detected due to the way " +
			"the API registry works. When a specification definition is updated, the registry UUID changes and is " +
			"only available from the response when the definition is published to the registry. When Terraform runs " +
			"after an external update, there's no way of programmatically retrieving the current state without the " +
			"current UUID. Forcing a Terraform update (e.g. tainting or a manual change) will get things " +
			"synchronized again.\n\n" +
			"## Importing Existing Specifications\n\n" +
			"Importing API specifications is limited due to the behavior of the API registry and associating a " +
			"specification with its definition. When importing, Terraform will replace the remote definition on its " +
			"next run, regardless if it differs from the local definition. This will associate a registry UUID " +
			"with the specification.\n\n" +
			"## Managing API Specification Docs\n\n" +
			"API Specifications created in ReadMe can have a documentation page associated with them. This is " +
			"automatically created by ReadMe when a specification is created. The documentation page is not " +
			"implicitly managed by Terraform. To manage the documentation page, use the `readme_doc` resource " +
			"with the `use_slug` attribute set to the API specification tag slug.\n\n" +
			"See <https://docs.readme.com/main/reference/uploadapispecification> for more information about this API " +
			"endpoint.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the API specification.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"category": schema.ObjectAttribute{
				Description: "Category metadata for the API specification.",
				Computed:    true,
				AttributeTypes: map[string]attr.Type{
					"id":    types.StringType,
					"slug":  types.StringType,
					"order": types.Int64Type,
					"title": types.StringType,
					"type":  types.StringType,
				},
			},
			"definition": schema.StringAttribute{
				Description: "The raw API specification definition JSON.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"definition_normalized": schema.StringAttribute{
				Description: "The normalized API specification definition JSON. This attribute is computed and " +
					"read-only. It is used to compare the definition with the remote definition.",
				Computed: true,
			},
			"delete_category": schema.BoolAttribute{
				Description: "Delete the category associated with the API specification when the resource is deleted.",
				Optional:    true,
			},
			"last_synced": schema.StringAttribute{
				Description: "Timestamp of last synchronization.",
				Computed:    true,
			},
			"uuid": schema.StringAttribute{
				Description: "The API registry UUID associated with the specification.",
				Computed:    true,
			},
			"source": schema.StringAttribute{
				Description: "The creation source of the API specification.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"title": schema.StringAttribute{
				Description: "The title of the API specification derived from the specification JSON.",
				Computed:    true,
			},
			"type": schema.StringAttribute{
				Description: "The type of the API specification.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"version": schema.StringAttribute{
				Description: "The version ID the API specification is associated with.",
				Computed:    true,
			},
			"semver": schema.StringAttribute{
				Description: "The semver(-ish) of the API specification. This value may also be set in the " +
					"definition JSON `info:version` key, but will be ignored if this attribute is set. Changing the " +
					"version of a created resource will replace the API specification. Use unique resources to use " +
					"the same specification across multiple versions.\n\n" +
					"Learn more about document versioning at <https://docs.readme.com/main/docs/versions>.",
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
			},
		},
	}
}

// Create creates the API Specification and sets the initial Terraform state.
func (r *apiSpecificationResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	// Retrieve values from the plan.
	var plan apiSpecificationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save the specification and update the state.
	if updatedPlan, err := r.save(ctx, saveActionCreate, "", plan); err != nil {
		resp.Diagnostics.AddError("Unable to create API specification.", err.Error())
	} else {
		resp.Diagnostics.Append(resp.State.Set(ctx, updatedPlan)...)
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *apiSpecificationResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	// Retrieve the current state.
	var state apiSpecificationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	version, err := r.resolveVersion(state)
	if err != nil {
		resp.Diagnostics.AddError("Unable to resolve version.", err.Error())

		return
	}

	if state.Definition.ValueString() == "" {
		resp.Diagnostics.AddWarning("No definition provided. Skipping read.", "")
	}

	// Fetch the specification if a UUID is available.
	if state.UUID.ValueString() != "" {
		if updatedState, err := r.makePlan(ctx, state.ID.ValueString(), state.DefinitionNorm, state.UUID.ValueString(), version); err != nil {
			if strings.Contains(err.Error(), "API specification not found") {
				tflog.Warn(ctx, fmt.Sprintf("API specification %s not found. Removing from state.", state.ID.ValueString()))
				resp.State.RemoveResource(ctx)
			} else {
				resp.Diagnostics.AddError("Unable to read API specification.", err.Error())
			}
		} else {
			updatedState.DeleteCategory = state.DeleteCategory
			match, _ := jsonMatch(state.Definition.ValueString(), updatedState.Definition.ValueString())
			if !match {
				state.Definition = updatedState.Definition
			} else {
				updatedState.Definition = state.Definition
			}
			resp.Diagnostics.Append(resp.State.Set(ctx, updatedState)...)
		}
	}
}

// Update updates the API Specification and sets the updated Terraform state on success.
func (r *apiSpecificationResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	// Retrieve values from the plan and current state.
	var plan, state apiSpecificationResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update the specification and refresh the state.
	if updatedPlan, err := r.save(ctx, saveActionUpdate, state.ID.ValueString(), plan); err != nil {
		resp.Diagnostics.AddError("Unable to update API specification.", err.Error())
	} else {
		updatedPlan.DeleteCategory = plan.DeleteCategory
		resp.Diagnostics.Append(resp.State.Set(ctx, updatedPlan)...)
	}
}

// Delete deletes the API Specification and removes the Terraform state on success.
func (r *apiSpecificationResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	// Retrieve values from state.
	var state apiSpecificationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, apiResponse, err := r.client.APISpecification.Delete(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to delete API specification.",
			clientError(err, apiResponse),
		)

		return
	}

	// Remove the category if delete_category is set to true.
	// When deleting a specification, its category is not deleted by the API.
	if state.DeleteCategory.ValueBool() {
		category := state.Category.Attributes()
		catSlug := category["slug"].String()
		// Remove double quotes
		catSlug = strings.ReplaceAll(catSlug, "\"", "")

		// Categories are versioned. Get the version ID from the state.

		version, err := r.resolveVersion(state)
		if err != nil {
			resp.Diagnostics.AddError("Unable to resolve version.", err.Error())

			return
		}

		opts := readme.RequestOptions{Version: version}
		_, apiResponse, err := r.client.Category.Delete(catSlug, opts)
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to delete category.",
				clientError(err, apiResponse),
			)

			return
		}
	}
}

// ImportState imports an API Specification by ID.
func (r *apiSpecificationResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	// Use the "id" attribute for importing.
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *apiSpecificationResource) save(
	ctx context.Context,
	action saveAction,
	specID string,
	plan apiSpecificationResourceModel,
) (apiSpecificationResourceModel, error) {
	version, err := r.resolveVersion(plan)
	if err != nil {
		return apiSpecificationResourceModel{}, fmt.Errorf("error resolving version: %w", err)
	}

	definition := plan.DefinitionNorm.ValueString()

	// Upload to the API registry.
	registry, err := r.createRegistry(definition, version)
	if err != nil {
		return apiSpecificationResourceModel{}, fmt.Errorf("error creating registry: %w", err)
	}

	// Perform the create or update action.
	response, err := r.performSaveAction(action, specID, registry.RegistryUUID, version)
	if err != nil {
		return apiSpecificationResourceModel{}, fmt.Errorf("error performing save action: %w", err)
	}

	// Create the final plan.
	m, err := r.makePlan(ctx, response.ID, plan.DefinitionNorm, registry.RegistryUUID, version)
	if err != nil {
		return apiSpecificationResourceModel{}, fmt.Errorf("error creating plan: %w", err)
	}

	m.DeleteCategory = plan.DeleteCategory
	m.Definition = plan.Definition

	return m, nil
}

func (r *apiSpecificationResource) resolveVersion(
	plan apiSpecificationResourceModel,
) (string, error) {
	// If a specific version is provided in the plan, use it.
	if plan.Semver.ValueString() != "" {
		return plan.Semver.ValueString(), nil
	}

	if plan.Version.ValueString() == "" {
		return "", fmt.Errorf("no version or semver provided in the plan")
	}

	versionInfo, _, err := r.client.Version.Get(IDPrefix + plan.Version.ValueString())
	if err != nil {
		return "", fmt.Errorf("error resolving version: %w", err)
	}

	return versionInfo.VersionClean, nil
}

func (r *apiSpecificationResource) performSaveAction(
	action saveAction,
	specID, registryUUID, version string,
) (readme.APISpecificationSaved, error) {
	if action == saveActionUpdate {

		resp, _, err := r.client.APISpecification.Update(specID, UUIDPrefix+registryUUID)

		return resp, err
	}
	requestOptions := readme.RequestOptions{Version: version}
	resp, _, err := r.client.APISpecification.Create(UUIDPrefix+registryUUID, requestOptions)

	return resp, err
}

// normalizeDefinition is a helper function that normalizes the definition JSON
// by setting the `info:version` key to the parameter attribute version when set.
func normalizeDefinition(version, definition string) (string, error) {
	if version == "" {
		return definition, nil
	}

	// Update the definition's version key to avoid churn.
	definitionJSON := map[string]interface{}{}
	err := json.Unmarshal([]byte(definition), &definitionJSON)
	if err != nil {
		return "", fmt.Errorf("unable to unmarshal definition: %w", err)
	}

	info, ok := definitionJSON["info"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("unable to get info from definition")
	}

	info["version"] = version
	definitionJSON["info"] = info

	// Marshal back to string.
	definitionBytes, err := json.Marshal(definitionJSON)
	if err != nil {
		return "", fmt.Errorf("unable to marshal definition: %w", err)
	}

	return string(definitionBytes), nil
}

func (r *apiSpecificationResource) ModifyPlan(
	ctx context.Context,
	req resource.ModifyPlanRequest,
	resp *resource.ModifyPlanResponse,
) {
	var plan, state *apiSpecificationResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() || plan == nil {
		return
	}

	// Only modify the plan if the definition is changing
	if plan.Definition.ValueString() == state.Definition.ValueString() {
		resp.Diagnostics.AddWarning("No changes detected in definition. Skipping plan modification.", "")
		plan.DefinitionNorm = state.DefinitionNorm
		plan.Definition = state.Definition

		return
	}

	// If the semver is set in the plan, ensure it is set in the definition.
	definition, err := normalizeDefinition(plan.Semver.ValueString(), plan.Definition.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to set info version in definition.", err.Error())

		return
	}

	plan.DefinitionNorm = types.StringValue(definition)

	diags := resp.Plan.Set(ctx, plan)

	resp.Diagnostics.Append(diags...)
}

func (r *apiSpecificationResource) makePlan(
	ctx context.Context,
	specID string,
	definition types.String,
	registryUUID, version string,
) (apiSpecificationResourceModel, error) {
	spec, err := r.get(ctx, specID, version)
	if err != nil {
		return apiSpecificationResourceModel{}, fmt.Errorf("error getting specification: %w", err)
	}

	// Map the plan to the resource struct.
	return apiSpecificationResourceModel{
		Category:       specCategoryObject(spec),
		DefinitionNorm: definition,
		ID:             types.StringValue(spec.ID),
		LastSynced:     types.StringValue(spec.LastSynced),
		Semver:         types.StringValue(version),
		Source:         types.StringValue(spec.Source),
		Title:          types.StringValue(spec.Title),
		Type:           types.StringValue(spec.Type),
		UUID:           types.StringValue(registryUUID),
		Version:        types.StringValue(spec.Version),
	}, nil
}

// get is a helper function that retrieves a specification by ID and returns a readme.APISpecification struct.
func (r *apiSpecificationResource) get(ctx context.Context, specID, version string) (readme.APISpecification, error) {
	requestOptions := readme.RequestOptions{Version: version}
	specification, _, err := r.client.APISpecification.Get(specID, requestOptions)
	if err != nil {
		tflog.Error(ctx, fmt.Sprintf("Unable to get specification: %+v", err))

		return specification, fmt.Errorf("unable to get specification id %s (version %s): %w. request options: %+v", specID, version, err, requestOptions)
	}

	if specification.ID == "" {
		return specification, fmt.Errorf("specification response is empty for specification ID %s", specID)
	}

	return specification, nil
}

// createRegistry is a helper function that creates an API registry definition in ReadMe. This is done before any create
// or update of an API specification.
func (r *apiSpecificationResource) createRegistry(
	definition, version string,
) (readme.APIRegistrySaved, error) {
	registry, apiResponse, err := r.client.APIRegistry.Create(definition, version)
	if err != nil {
		return readme.APIRegistrySaved{}, errors.New(clientError(err, apiResponse))
	}

	return registry, nil
}
