package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &FormatsDataSource{}

func NewFormatsDataSource() datasource.DataSource {
	return &FormatsDataSource{}
}

// FormatsDataSource defines the data source implementation.
type FormatsDataSource struct {
	client *BMTLClientData
}

// FormatsDataSourceModel describes the data source data model.
type FormatsDataSourceModel struct {
	Formats []FormatModel `tfsdk:"formats"`
	Id      types.String  `tfsdk:"id"`
}

type FormatModel struct {
	Id           types.Int64              `tfsdk:"id"`
	WorldId      types.String             `tfsdk:"world_id"`
	Type         types.String             `tfsdk:"type"`
	Translations []FormatTranslationModel `tfsdk:"translations"`
}

type FormatTranslationModel struct {
	Key         types.String `tfsdk:"key"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Language    types.String `tfsdk:"language"`
}

func (d *FormatsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest,
	resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_formats"
}

func (d *FormatsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Formats data source allows you to retrieve information about available meeting formats.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Placeholder identifier for the data source.",
				Computed:            true,
			},
			"formats": schema.ListNestedAttribute{
				MarkdownDescription: "List of formats",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							MarkdownDescription: "Format identifier",
							Computed:            true,
						},
						"world_id": schema.StringAttribute{
							MarkdownDescription: "World identifier for the format",
							Computed:            true,
						},
						"type": schema.StringAttribute{
							MarkdownDescription: "Format type",
							Computed:            true,
						},
						"translations": schema.ListNestedAttribute{
							MarkdownDescription: "Format translations",
							Computed:            true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"key": schema.StringAttribute{
										MarkdownDescription: "Translation key",
										Computed:            true,
									},
									"name": schema.StringAttribute{
										MarkdownDescription: "Format name",
										Computed:            true,
									},
									"description": schema.StringAttribute{
										MarkdownDescription: "Format description",
										Computed:            true,
									},
									"language": schema.StringAttribute{
										MarkdownDescription: "Language code",
										Computed:            true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (d *FormatsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest,
	resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*BMTLClientData)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			clientTypeError(req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *FormatsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data FormatsDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get formats from the API
	formats, httpResp, err := d.client.Client.RootServerAPI.GetFormats(d.client.Context).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read formats, got error: %s", err))
		return
	}

	if httpResp.StatusCode != HTTPStatusOK {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("API returned status %d", httpResp.StatusCode))
		return
	}

	// Map response body to model
	for _, format := range formats {
		formatModel := FormatModel{
			Id:      types.Int64Value(int64(format.Id)),
			WorldId: types.StringValue(format.WorldId),
			Type:    types.StringValue(format.Type),
		}

		// Handle translations
		var translations []FormatTranslationModel
		if format.Translations != nil {
			for _, translation := range format.Translations {
				translationModel := FormatTranslationModel{
					Key:         types.StringValue(translation.Key),
					Name:        types.StringValue(translation.Name),
					Description: types.StringValue(translation.Description),
					Language:    types.StringValue(translation.Language),
				}
				translations = append(translations, translationModel)
			}
		}
		formatModel.Translations = translations

		data.Formats = append(data.Formats, formatModel)
	}

	// Set ID for the data source
	data.Id = types.StringValue("placeholder")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
