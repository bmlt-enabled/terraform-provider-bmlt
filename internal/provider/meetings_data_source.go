package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &MeetingsDataSource{}

func NewMeetingsDataSource() datasource.DataSource {
	return &MeetingsDataSource{}
}

// MeetingsDataSource defines the data source implementation.
type MeetingsDataSource struct {
	client *BMTLClientData
}

// MeetingsDataSourceModel describes the data source data model.
type MeetingsDataSourceModel struct {
	Meetings       []MeetingModel `tfsdk:"meetings"`
	Id             types.String   `tfsdk:"id"`
	MeetingIds     types.String   `tfsdk:"meeting_ids"`
	Days           types.String   `tfsdk:"days"`
	ServiceBodyIds types.String   `tfsdk:"service_body_ids"`
	SearchString   types.String   `tfsdk:"search_string"`
}

type MeetingModel struct {
	Id                           types.Int64   `tfsdk:"id"`
	ServiceBodyId                types.Int64   `tfsdk:"service_body_id"`
	FormatIds                    []types.Int64 `tfsdk:"format_ids"`
	VenueType                    types.Int64   `tfsdk:"venue_type"`
	TemporarilyVirtual           types.Bool    `tfsdk:"temporarily_virtual"`
	Day                          types.Int64   `tfsdk:"day"`
	StartTime                    types.String  `tfsdk:"start_time"`
	Duration                     types.String  `tfsdk:"duration"`
	TimeZone                     types.String  `tfsdk:"time_zone"`
	Latitude                     types.Float64 `tfsdk:"latitude"`
	Longitude                    types.Float64 `tfsdk:"longitude"`
	Published                    types.Bool    `tfsdk:"published"`
	Email                        types.String  `tfsdk:"email"`
	WorldId                      types.String  `tfsdk:"world_id"`
	Name                         types.String  `tfsdk:"name"`
	LocationText                 types.String  `tfsdk:"location_text"`
	LocationInfo                 types.String  `tfsdk:"location_info"`
	LocationStreet               types.String  `tfsdk:"location_street"`
	LocationNeighborhood         types.String  `tfsdk:"location_neighborhood"`
	LocationCitySubsection       types.String  `tfsdk:"location_city_subsection"`
	LocationMunicipality         types.String  `tfsdk:"location_municipality"`
	LocationSubProvince          types.String  `tfsdk:"location_sub_province"`
	LocationProvince             types.String  `tfsdk:"location_province"`
	LocationPostalCode1          types.String  `tfsdk:"location_postal_code_1"`
	LocationNation               types.String  `tfsdk:"location_nation"`
	PhoneMeetingNumber           types.String  `tfsdk:"phone_meeting_number"`
	VirtualMeetingLink           types.String  `tfsdk:"virtual_meeting_link"`
	VirtualMeetingAdditionalInfo types.String  `tfsdk:"virtual_meeting_additional_info"`
	ContactName1                 types.String  `tfsdk:"contact_name_1"`
	ContactName2                 types.String  `tfsdk:"contact_name_2"`
	ContactPhone1                types.String  `tfsdk:"contact_phone_1"`
	ContactPhone2                types.String  `tfsdk:"contact_phone_2"`
	ContactEmail1                types.String  `tfsdk:"contact_email_1"`
	ContactEmail2                types.String  `tfsdk:"contact_email_2"`
	BusLines                     types.String  `tfsdk:"bus_lines"`
	TrainLines                   types.String  `tfsdk:"train_lines"`
	Comments                     types.String  `tfsdk:"comments"`
}

func (d *MeetingsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_meetings"
}

func (d *MeetingsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Meetings data source allows you to retrieve information about meetings.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Placeholder identifier for the data source.",
				Computed:            true,
			},
			"meeting_ids": schema.StringAttribute{
				MarkdownDescription: "Comma delimited meeting ids to filter by",
				Optional:            true,
			},
			"days": schema.StringAttribute{
				MarkdownDescription: "Comma delimited day ids between 0-6 to filter by",
				Optional:            true,
			},
			"service_body_ids": schema.StringAttribute{
				MarkdownDescription: "Comma delimited service body ids to filter by",
				Optional:            true,
			},
			"search_string": schema.StringAttribute{
				MarkdownDescription: "Search string to filter meetings",
				Optional:            true,
			},
			"meetings": schema.ListNestedAttribute{
				MarkdownDescription: "List of meetings",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							MarkdownDescription: "Meeting identifier",
							Computed:            true,
						},
						"service_body_id": schema.Int64Attribute{
							MarkdownDescription: "Service body identifier",
							Computed:            true,
						},
						"format_ids": schema.ListAttribute{
							MarkdownDescription: "List of format identifiers",
							Computed:            true,
							ElementType:         types.Int64Type,
						},
						"venue_type": schema.Int64Attribute{
							MarkdownDescription: "Venue type (1=in-person, 2=virtual, 3=hybrid)",
							Computed:            true,
						},
						"temporarily_virtual": schema.BoolAttribute{
							MarkdownDescription: "Whether the meeting is temporarily virtual",
							Computed:            true,
						},
						"day": schema.Int64Attribute{
							MarkdownDescription: "Day of the week (0=Sunday, 1=Monday, etc.)",
							Computed:            true,
						},
						"start_time": schema.StringAttribute{
							MarkdownDescription: "Meeting start time",
							Computed:            true,
						},
						"duration": schema.StringAttribute{
							MarkdownDescription: "Meeting duration",
							Computed:            true,
						},
						"time_zone": schema.StringAttribute{
							MarkdownDescription: "Time zone",
							Computed:            true,
						},
						"latitude": schema.Float64Attribute{
							MarkdownDescription: "Latitude coordinate",
							Computed:            true,
						},
						"longitude": schema.Float64Attribute{
							MarkdownDescription: "Longitude coordinate",
							Computed:            true,
						},
						"published": schema.BoolAttribute{
							MarkdownDescription: "Whether the meeting is published",
							Computed:            true,
						},
						"email": schema.StringAttribute{
							MarkdownDescription: "Meeting email",
							Computed:            true,
						},
						"world_id": schema.StringAttribute{
							MarkdownDescription: "World identifier",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Meeting name",
							Computed:            true,
						},
						"location_text": schema.StringAttribute{
							MarkdownDescription: "Location text",
							Computed:            true,
						},
						"location_info": schema.StringAttribute{
							MarkdownDescription: "Location info",
							Computed:            true,
						},
						"location_street": schema.StringAttribute{
							MarkdownDescription: "Street address",
							Computed:            true,
						},
						"location_neighborhood": schema.StringAttribute{
							MarkdownDescription: "Neighborhood",
							Computed:            true,
						},
						"location_city_subsection": schema.StringAttribute{
							MarkdownDescription: "City subsection",
							Computed:            true,
						},
						"location_municipality": schema.StringAttribute{
							MarkdownDescription: "Municipality",
							Computed:            true,
						},
						"location_sub_province": schema.StringAttribute{
							MarkdownDescription: "Sub province",
							Computed:            true,
						},
						"location_province": schema.StringAttribute{
							MarkdownDescription: "Province",
							Computed:            true,
						},
						"location_postal_code_1": schema.StringAttribute{
							MarkdownDescription: "Postal code",
							Computed:            true,
						},
						"location_nation": schema.StringAttribute{
							MarkdownDescription: "Nation",
							Computed:            true,
						},
						"phone_meeting_number": schema.StringAttribute{
							MarkdownDescription: "Phone meeting number",
							Computed:            true,
						},
						"virtual_meeting_link": schema.StringAttribute{
							MarkdownDescription: "Virtual meeting link",
							Computed:            true,
						},
						"virtual_meeting_additional_info": schema.StringAttribute{
							MarkdownDescription: "Additional virtual meeting info",
							Computed:            true,
						},
						"contact_name_1": schema.StringAttribute{
							MarkdownDescription: "Primary contact name",
							Computed:            true,
						},
						"contact_name_2": schema.StringAttribute{
							MarkdownDescription: "Secondary contact name",
							Computed:            true,
						},
						"contact_phone_1": schema.StringAttribute{
							MarkdownDescription: "Primary contact phone",
							Computed:            true,
						},
						"contact_phone_2": schema.StringAttribute{
							MarkdownDescription: "Secondary contact phone",
							Computed:            true,
						},
						"contact_email_1": schema.StringAttribute{
							MarkdownDescription: "Primary contact email",
							Computed:            true,
						},
						"contact_email_2": schema.StringAttribute{
							MarkdownDescription: "Secondary contact email",
							Computed:            true,
						},
						"bus_lines": schema.StringAttribute{
							MarkdownDescription: "Bus lines",
							Computed:            true,
						},
						"train_lines": schema.StringAttribute{
							MarkdownDescription: "Train lines",
							Computed:            true,
						},
						"comments": schema.StringAttribute{
							MarkdownDescription: "Comments",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *MeetingsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*BMTLClientData)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *BMTLClientData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = client
}

func (d *MeetingsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data MeetingsDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the API request with optional parameters
	apiReq := d.client.Client.RootServerAPI.GetMeetings(d.client.Context)

	if !data.MeetingIds.IsNull() && data.MeetingIds.ValueString() != "" {
		apiReq = apiReq.MeetingIds(data.MeetingIds.ValueString())
	}
	if !data.Days.IsNull() && data.Days.ValueString() != "" {
		apiReq = apiReq.Days(data.Days.ValueString())
	}
	if !data.ServiceBodyIds.IsNull() && data.ServiceBodyIds.ValueString() != "" {
		apiReq = apiReq.ServiceBodyIds(data.ServiceBodyIds.ValueString())
	}
	if !data.SearchString.IsNull() && data.SearchString.ValueString() != "" {
		apiReq = apiReq.SearchString(data.SearchString.ValueString())
	}

	// Execute the request
	meetings, httpResp, err := apiReq.Execute()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read meetings, got error: %s", err))
		return
	}

	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("API returned status %d", httpResp.StatusCode))
		return
	}

	// Map response body to model
	for _, meeting := range meetings {
		meetingModel := MeetingModel{
			Id:                           types.Int64Value(int64(meeting.Id)),
			ServiceBodyId:                types.Int64Value(int64(meeting.ServiceBodyId)),
			VenueType:                    types.Int64Value(int64(meeting.VenueType)),
			TemporarilyVirtual:           types.BoolValue(meeting.TemporarilyVirtual),
			Day:                          types.Int64Value(int64(meeting.Day)),
			StartTime:                    types.StringValue(meeting.StartTime),
			Duration:                     types.StringValue(meeting.Duration),
			TimeZone:                     types.StringValue(meeting.TimeZone),
			Latitude:                     types.Float64Value(float64(meeting.Latitude)),
			Longitude:                    types.Float64Value(float64(meeting.Longitude)),
			Published:                    types.BoolValue(meeting.Published),
			Email:                        types.StringValue(meeting.Email),
			WorldId:                      types.StringValue(meeting.WorldId),
			Name:                         types.StringValue(meeting.Name),
			LocationText:                 types.StringPointerValue(meeting.LocationText),
			LocationInfo:                 types.StringPointerValue(meeting.LocationInfo),
			LocationStreet:               types.StringPointerValue(meeting.LocationStreet),
			LocationNeighborhood:         types.StringPointerValue(meeting.LocationNeighborhood),
			LocationCitySubsection:       types.StringPointerValue(meeting.LocationCitySubsection),
			LocationMunicipality:         types.StringPointerValue(meeting.LocationMunicipality),
			LocationSubProvince:          types.StringPointerValue(meeting.LocationSubProvince),
			LocationProvince:             types.StringPointerValue(meeting.LocationProvince),
			LocationPostalCode1:          types.StringPointerValue(meeting.LocationPostalCode1),
			LocationNation:               types.StringPointerValue(meeting.LocationNation),
			PhoneMeetingNumber:           types.StringPointerValue(meeting.PhoneMeetingNumber),
			VirtualMeetingLink:           types.StringPointerValue(meeting.VirtualMeetingLink),
			VirtualMeetingAdditionalInfo: types.StringPointerValue(meeting.VirtualMeetingAdditionalInfo),
			ContactName1:                 types.StringPointerValue(meeting.ContactName1),
			ContactName2:                 types.StringPointerValue(meeting.ContactName2),
			ContactPhone1:                types.StringPointerValue(meeting.ContactPhone1),
			ContactPhone2:                types.StringPointerValue(meeting.ContactPhone2),
			ContactEmail1:                types.StringPointerValue(meeting.ContactEmail1),
			ContactEmail2:                types.StringPointerValue(meeting.ContactEmail2),
			BusLines:                     types.StringPointerValue(meeting.BusLines),
			TrainLines:                   types.StringPointerValue(meeting.TrainLines),
			Comments:                     types.StringPointerValue(meeting.Comments),
		}

		// Handle format IDs
		var formatIds []types.Int64
		for _, formatId := range meeting.FormatIds {
			formatIds = append(formatIds, types.Int64Value(int64(formatId)))
		}
		meetingModel.FormatIds = formatIds

		data.Meetings = append(data.Meetings, meetingModel)
	}

	// Set ID for the data source
	var idParts []string
	if !data.MeetingIds.IsNull() {
		idParts = append(idParts, "meeting_ids="+data.MeetingIds.ValueString())
	}
	if !data.Days.IsNull() {
		idParts = append(idParts, "days="+data.Days.ValueString())
	}
	if !data.ServiceBodyIds.IsNull() {
		idParts = append(idParts, "service_body_ids="+data.ServiceBodyIds.ValueString())
	}
	if !data.SearchString.IsNull() {
		idParts = append(idParts, "search_string="+data.SearchString.ValueString())
	}

	if len(idParts) == 0 {
		data.Id = types.StringValue("all")
	} else {
		data.Id = types.StringValue(strings.Join(idParts, "&"))
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
