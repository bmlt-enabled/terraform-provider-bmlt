package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &SettingsDataSource{}

func NewSettingsDataSource() datasource.DataSource {
	return &SettingsDataSource{}
}

// SettingsDataSource defines the data source implementation.
type SettingsDataSource struct {
	client *BMTLClientData
}

// SettingsDataSourceModel describes the data source data model.
type SettingsDataSourceModel struct {
	Id                                types.String   `tfsdk:"id"`
	GoogleApiKey                      types.String   `tfsdk:"google_api_key"`
	ChangeDepthForMeetings            types.Int64    `tfsdk:"change_depth_for_meetings"`
	DefaultSortKey                    types.String   `tfsdk:"default_sort_key"`
	Language                          types.String   `tfsdk:"language"`
	DefaultDurationTime               types.String   `tfsdk:"default_duration_time"`
	RegionBias                        types.String   `tfsdk:"region_bias"`
	DistanceUnits                     types.String   `tfsdk:"distance_units"`
	MeetingStatesAndProvinces         []types.String `tfsdk:"meeting_states_and_provinces"`
	MeetingCountiesAndSubProvinces    []types.String `tfsdk:"meeting_counties_and_sub_provinces"`
	SearchSpecMapCenterLongitude      types.Float64  `tfsdk:"search_spec_map_center_longitude"`
	SearchSpecMapCenterLatitude       types.Float64  `tfsdk:"search_spec_map_center_latitude"`
	SearchSpecMapCenterZoom           types.Int64    `tfsdk:"search_spec_map_center_zoom"`
	NumberOfMeetingsForAuto           types.Int64    `tfsdk:"number_of_meetings_for_auto"`
	AutoGeocodingEnabled              types.Bool     `tfsdk:"auto_geocoding_enabled"`
	CountyAutoGeocodingEnabled        types.Bool     `tfsdk:"county_auto_geocoding_enabled"`
	ZipAutoGeocodingEnabled           types.Bool     `tfsdk:"zip_auto_geocoding_enabled"`
	DefaultClosedStatus               types.Bool     `tfsdk:"default_closed_status"`
	EnableLanguageSelector            types.Bool     `tfsdk:"enable_language_selector"`
	IncludeServiceBodyEmailInSemantic types.Bool     `tfsdk:"include_service_body_email_in_semantic"`
	BmltTitle                         types.String   `tfsdk:"bmlt_title"`
	BmltNotice                        types.String   `tfsdk:"bmlt_notice"`
}

func (d *SettingsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_settings"
}

func (d *SettingsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Settings data source. Retrieves all BMLT server settings.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Settings identifier (always 'settings' for this singleton resource)",
			},
			"google_api_key": schema.StringAttribute{
				MarkdownDescription: "Google API key for geocoding",
				Computed:            true,
				Sensitive:           true,
			},
			"change_depth_for_meetings": schema.Int64Attribute{
				MarkdownDescription: "Change depth for meetings",
				Computed:            true,
			},
			"default_sort_key": schema.StringAttribute{
				MarkdownDescription: "Default sort key for meetings",
				Computed:            true,
			},
			"language": schema.StringAttribute{
				MarkdownDescription: "Default language for the server",
				Computed:            true,
			},
			"default_duration_time": schema.StringAttribute{
				MarkdownDescription: "Default duration time for meetings",
				Computed:            true,
			},
			"region_bias": schema.StringAttribute{
				MarkdownDescription: "Region bias for geocoding",
				Computed:            true,
			},
			"distance_units": schema.StringAttribute{
				MarkdownDescription: "Distance units (e.g., 'mi' for miles, 'km' for kilometers)",
				Computed:            true,
			},
			"meeting_states_and_provinces": schema.ListAttribute{
				MarkdownDescription: "List of meeting states and provinces",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"meeting_counties_and_sub_provinces": schema.ListAttribute{
				MarkdownDescription: "List of meeting counties and sub-provinces",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"search_spec_map_center_longitude": schema.Float64Attribute{
				MarkdownDescription: "Search specification map center longitude",
				Computed:            true,
			},
			"search_spec_map_center_latitude": schema.Float64Attribute{
				MarkdownDescription: "Search specification map center latitude",
				Computed:            true,
			},
			"search_spec_map_center_zoom": schema.Int64Attribute{
				MarkdownDescription: "Search specification map center zoom level",
				Computed:            true,
			},
			"number_of_meetings_for_auto": schema.Int64Attribute{
				MarkdownDescription: "Number of meetings for auto geocoding",
				Computed:            true,
			},
			"auto_geocoding_enabled": schema.BoolAttribute{
				MarkdownDescription: "Whether auto geocoding is enabled",
				Computed:            true,
			},
			"county_auto_geocoding_enabled": schema.BoolAttribute{
				MarkdownDescription: "Whether county auto geocoding is enabled",
				Computed:            true,
			},
			"zip_auto_geocoding_enabled": schema.BoolAttribute{
				MarkdownDescription: "Whether ZIP code auto geocoding is enabled",
				Computed:            true,
			},
			"default_closed_status": schema.BoolAttribute{
				MarkdownDescription: "Default closed status for meetings",
				Computed:            true,
			},
			"enable_language_selector": schema.BoolAttribute{
				MarkdownDescription: "Whether to enable the language selector",
				Computed:            true,
			},
			"include_service_body_email_in_semantic": schema.BoolAttribute{
				MarkdownDescription: "Whether to include service body email in semantic output",
				Computed:            true,
			},
			"bmlt_title": schema.StringAttribute{
				MarkdownDescription: "BMLT server title",
				Computed:            true,
			},
			"bmlt_notice": schema.StringAttribute{
				MarkdownDescription: "BMLT server notice",
				Computed:            true,
			},
		},
	}
}

func (d *SettingsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *SettingsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data SettingsDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get settings from API
	settings, httpResp, err := d.client.Client.RootServerAPI.GetSettings(d.client.Context).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read settings, got error: %s", err))
		return
	}

	if httpResp.StatusCode != HTTPStatusOK {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("API returned status %d", httpResp.StatusCode))
		return
	}

	// Set a constant ID for this singleton data source
	data.Id = types.StringValue("settings")

	// Map settings to model
	if settings.GoogleApiKey != nil {
		data.GoogleApiKey = types.StringValue(*settings.GoogleApiKey)
	} else {
		data.GoogleApiKey = types.StringNull()
	}

	if settings.ChangeDepthForMeetings != nil {
		data.ChangeDepthForMeetings = types.Int64Value(int64(*settings.ChangeDepthForMeetings))
	} else {
		data.ChangeDepthForMeetings = types.Int64Null()
	}

	if settings.DefaultSortKey.IsSet() && settings.DefaultSortKey.Get() != nil {
		data.DefaultSortKey = types.StringValue(*settings.DefaultSortKey.Get())
	} else {
		data.DefaultSortKey = types.StringNull()
	}

	if settings.Language != nil {
		data.Language = types.StringValue(*settings.Language)
	} else {
		data.Language = types.StringNull()
	}

	if settings.DefaultDurationTime != nil {
		data.DefaultDurationTime = types.StringValue(*settings.DefaultDurationTime)
	} else {
		data.DefaultDurationTime = types.StringNull()
	}

	if settings.RegionBias != nil {
		data.RegionBias = types.StringValue(*settings.RegionBias)
	} else {
		data.RegionBias = types.StringNull()
	}

	if settings.DistanceUnits != nil {
		data.DistanceUnits = types.StringValue(*settings.DistanceUnits)
	} else {
		data.DistanceUnits = types.StringNull()
	}

	if settings.MeetingStatesAndProvinces != nil {
		var states []types.String
		for _, s := range settings.MeetingStatesAndProvinces {
			states = append(states, types.StringValue(s))
		}
		data.MeetingStatesAndProvinces = states
	} else {
		data.MeetingStatesAndProvinces = nil
	}

	if settings.MeetingCountiesAndSubProvinces != nil {
		var counties []types.String
		for _, c := range settings.MeetingCountiesAndSubProvinces {
			counties = append(counties, types.StringValue(c))
		}
		data.MeetingCountiesAndSubProvinces = counties
	} else {
		data.MeetingCountiesAndSubProvinces = nil
	}

	if settings.SearchSpecMapCenterLongitude != nil {
		data.SearchSpecMapCenterLongitude = types.Float64Value(float64(*settings.SearchSpecMapCenterLongitude))
	} else {
		data.SearchSpecMapCenterLongitude = types.Float64Null()
	}

	if settings.SearchSpecMapCenterLatitude != nil {
		data.SearchSpecMapCenterLatitude = types.Float64Value(float64(*settings.SearchSpecMapCenterLatitude))
	} else {
		data.SearchSpecMapCenterLatitude = types.Float64Null()
	}

	if settings.SearchSpecMapCenterZoom != nil {
		data.SearchSpecMapCenterZoom = types.Int64Value(int64(*settings.SearchSpecMapCenterZoom))
	} else {
		data.SearchSpecMapCenterZoom = types.Int64Null()
	}

	if settings.NumberOfMeetingsForAuto != nil {
		data.NumberOfMeetingsForAuto = types.Int64Value(int64(*settings.NumberOfMeetingsForAuto))
	} else {
		data.NumberOfMeetingsForAuto = types.Int64Null()
	}

	if settings.AutoGeocodingEnabled != nil {
		data.AutoGeocodingEnabled = types.BoolValue(*settings.AutoGeocodingEnabled)
	} else {
		data.AutoGeocodingEnabled = types.BoolNull()
	}

	if settings.CountyAutoGeocodingEnabled != nil {
		data.CountyAutoGeocodingEnabled = types.BoolValue(*settings.CountyAutoGeocodingEnabled)
	} else {
		data.CountyAutoGeocodingEnabled = types.BoolNull()
	}

	if settings.ZipAutoGeocodingEnabled != nil {
		data.ZipAutoGeocodingEnabled = types.BoolValue(*settings.ZipAutoGeocodingEnabled)
	} else {
		data.ZipAutoGeocodingEnabled = types.BoolNull()
	}

	if settings.DefaultClosedStatus != nil {
		data.DefaultClosedStatus = types.BoolValue(*settings.DefaultClosedStatus)
	} else {
		data.DefaultClosedStatus = types.BoolNull()
	}

	if settings.EnableLanguageSelector != nil {
		data.EnableLanguageSelector = types.BoolValue(*settings.EnableLanguageSelector)
	} else {
		data.EnableLanguageSelector = types.BoolNull()
	}

	if settings.IncludeServiceBodyEmailInSemantic != nil {
		data.IncludeServiceBodyEmailInSemantic = types.BoolValue(*settings.IncludeServiceBodyEmailInSemantic)
	} else {
		data.IncludeServiceBodyEmailInSemantic = types.BoolNull()
	}

	if settings.BmltTitle != nil {
		data.BmltTitle = types.StringValue(*settings.BmltTitle)
	} else {
		data.BmltTitle = types.StringNull()
	}

	if settings.BmltNotice != nil {
		data.BmltNotice = types.StringValue(*settings.BmltNotice)
	} else {
		data.BmltNotice = types.StringNull()
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
