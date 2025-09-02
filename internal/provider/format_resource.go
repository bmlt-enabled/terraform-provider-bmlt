package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/bmlt-enabled/bmlt-server-go-client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &FormatResource{}
var _ resource.ResourceWithImportState = &FormatResource{}

func NewFormatResource() resource.Resource {
	return &FormatResource{}
}

// FormatResource defines the resource implementation.
type FormatResource struct {
	client *BMTLClientData
}

// FormatResourceModel describes the resource data model.
type FormatResourceModel struct {
	Id           types.String             `tfsdk:"id"`
	WorldId      types.String             `tfsdk:"world_id"`
	Type         types.String             `tfsdk:"type"`
	Translations []FormatTranslationModel `tfsdk:"translations"`
}

func (r *FormatResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_format"
}

func (r *FormatResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Format resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Format identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"world_id": schema.StringAttribute{
				MarkdownDescription: "World identifier for the format",
				Optional:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Format type",
				Optional:            true,
			},
			"translations": schema.ListNestedAttribute{
				MarkdownDescription: "Format translations",
				Required:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"key": schema.StringAttribute{
							MarkdownDescription: "Translation key",
							Required:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Format name",
							Required:            true,
						},
						"description": schema.StringAttribute{
							MarkdownDescription: "Format description",
							Required:            true,
						},
						"language": schema.StringAttribute{
							MarkdownDescription: "Language code",
							Required:            true,
						},
					},
				},
			},
		},
	}
}

func (r *FormatResource) Configure(ctx context.Context, req resource.ConfigureRequest,
	resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*BMTLClientData)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			clientTypeError(req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *FormatResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *FormatResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Convert model to API request
	createRequest := bmlt.FormatCreate{
		WorldId: data.WorldId.ValueStringPointer(),
		Type:    data.Type.ValueStringPointer(),
	}

	// Convert translations
	var translations []bmlt.FormatTranslation
	for _, t := range data.Translations {
		translation := bmlt.FormatTranslation{
			Key:         t.Key.ValueString(),
			Name:        t.Name.ValueString(),
			Description: t.Description.ValueString(),
			Language:    t.Language.ValueString(),
		}
		translations = append(translations, translation)
	}
	createRequest.Translations = translations

	// Create format
	format, httpResp, err := r.client.Client.RootServerAPI.CreateFormat(r.client.Context).
		FormatCreate(createRequest).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create format, got error: %s", err))
		return
	}

	if httpResp.StatusCode != HTTPStatusCreated {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("API returned status %d", httpResp.StatusCode))
		return
	}

	// Map response back to model
	data.Id = types.StringValue(strconv.Itoa(int(format.Id)))
	data.WorldId = nullableString(format.WorldId)
	data.Type = nullableString(format.Type)

	// Update translations from response
	var responseTranslations []FormatTranslationModel
	if format.Translations != nil {
		for _, t := range format.Translations {
			translationModel := FormatTranslationModel{
				Key:         types.StringValue(t.Key),
				Name:        types.StringValue(t.Name),
				Description: types.StringValue(t.Description),
				Language:    types.StringValue(t.Language),
			}
			responseTranslations = append(responseTranslations, translationModel)
		}
	}
	data.Translations = responseTranslations

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FormatResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *FormatResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Convert ID to int64
	id, err := strconv.ParseInt(data.Id.ValueString(), 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse format ID: %s", err))
		return
	}

	// Get format from API
	format, httpResp, err := r.client.Client.RootServerAPI.GetFormat(r.client.Context, id).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read format, got error: %s", err))
		return
	}

	if httpResp.StatusCode == HTTPStatusNotFound {
		// Format was deleted outside of Terraform
		resp.State.RemoveResource(ctx)
		return
	}

	if httpResp.StatusCode != HTTPStatusOK {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("API returned status %d", httpResp.StatusCode))
		return
	}

	// Map response to model
	data.WorldId = nullableString(format.WorldId)
	data.Type = nullableString(format.Type)

	// Update translations
	var translations []FormatTranslationModel
	if format.Translations != nil {
		for _, t := range format.Translations {
			translationModel := FormatTranslationModel{
				Key:         types.StringValue(t.Key),
				Name:        types.StringValue(t.Name),
				Description: types.StringValue(t.Description),
				Language:    types.StringValue(t.Language),
			}
			translations = append(translations, translationModel)
		}
	}
	data.Translations = translations

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FormatResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *FormatResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Convert ID to int64
	id, err := strconv.ParseInt(data.Id.ValueString(), 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse format ID: %s", err))
		return
	}

	// Convert model to API request
	updateRequest := bmlt.FormatUpdate{
		WorldId: data.WorldId.ValueStringPointer(),
		Type:    data.Type.ValueStringPointer(),
	}

	// Convert translations
	var translations []bmlt.FormatTranslation
	for _, t := range data.Translations {
		translation := bmlt.FormatTranslation{
			Key:         t.Key.ValueString(),
			Name:        t.Name.ValueString(),
			Description: t.Description.ValueString(),
			Language:    t.Language.ValueString(),
		}
		translations = append(translations, translation)
	}
	updateRequest.Translations = translations

	// Update format
	httpResp, err := r.client.Client.RootServerAPI.UpdateFormat(r.client.Context, id).FormatUpdate(updateRequest).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update format, got error: %s", err))
		return
	}

	if httpResp.StatusCode != HTTPStatusNoContent {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("API returned status %d", httpResp.StatusCode))
		return
	}

	// Re-read the format to ensure state is consistent with server
	formatId, err := strconv.ParseInt(data.Id.ValueString(), 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse format ID: %s", err))
		return
	}

	updatedFormat, httpResp, err := r.client.Client.RootServerAPI.GetFormat(r.client.Context, formatId).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read updated format, got error: %s", err))
		return
	}

	if httpResp.StatusCode != HTTPStatusOK {
		resp.Diagnostics.AddError("API Error",
			fmt.Sprintf("API returned status %d when reading updated format", httpResp.StatusCode))
		return
	}

	// Update all fields from the server response
	data.WorldId = nullableString(updatedFormat.WorldId)
	data.Type = nullableString(updatedFormat.Type)

	// Update translations from response
	var responseTranslations []FormatTranslationModel
	if updatedFormat.Translations != nil {
		for _, t := range updatedFormat.Translations {
			translationModel := FormatTranslationModel{
				Key:         types.StringValue(t.Key),
				Name:        types.StringValue(t.Name),
				Description: types.StringValue(t.Description),
				Language:    types.StringValue(t.Language),
			}
			responseTranslations = append(responseTranslations, translationModel)
		}
	}
	data.Translations = responseTranslations

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FormatResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *FormatResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Convert ID to int64
	id, err := strconv.ParseInt(data.Id.ValueString(), 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse format ID: %s", err))
		return
	}

	// Delete format
	httpResp, err := r.client.Client.RootServerAPI.DeleteFormat(r.client.Context, id).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete format, got error: %s", err))
		return
	}

	if httpResp.StatusCode != HTTPStatusNoContent {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("API returned status %d", httpResp.StatusCode))
		return
	}
}

func (r *FormatResource) ImportState(ctx context.Context, req resource.ImportStateRequest,
	resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
