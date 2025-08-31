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

var _ resource.Resource = &MeetingResource{}
var _ resource.ResourceWithImportState = &MeetingResource{}

func NewMeetingResource() resource.Resource {
	return &MeetingResource{}
}

type MeetingResource struct {
	client *BMTLClientData
}

// MeetingResourceModel describes the resource data model (focusing on key fields)
type MeetingResourceModel struct {
	Id                   types.String  `tfsdk:"id"`
	ServiceBodyId        types.Int64   `tfsdk:"service_body_id"`
	FormatIds            []types.Int64 `tfsdk:"format_ids"`
	VenueType            types.Int64   `tfsdk:"venue_type"`
	TemporarilyVirtual   types.Bool    `tfsdk:"temporarily_virtual"`
	Day                  types.Int64   `tfsdk:"day"`
	StartTime            types.String  `tfsdk:"start_time"`
	Duration             types.String  `tfsdk:"duration"`
	TimeZone             types.String  `tfsdk:"time_zone"`
	Latitude             types.Float64 `tfsdk:"latitude"`
	Longitude            types.Float64 `tfsdk:"longitude"`
	Published            types.Bool    `tfsdk:"published"`
	Email                types.String  `tfsdk:"email"`
	WorldId              types.String  `tfsdk:"world_id"`
	Name                 types.String  `tfsdk:"name"`
	LocationText         types.String  `tfsdk:"location_text"`
	LocationInfo         types.String  `tfsdk:"location_info"`
	LocationStreet       types.String  `tfsdk:"location_street"`
	LocationMunicipality types.String  `tfsdk:"location_municipality"`
	LocationProvince     types.String  `tfsdk:"location_province"`
	LocationPostalCode1  types.String  `tfsdk:"location_postal_code_1"`
	LocationNation       types.String  `tfsdk:"location_nation"`
	VirtualMeetingLink   types.String  `tfsdk:"virtual_meeting_link"`
	ContactName1         types.String  `tfsdk:"contact_name_1"`
	ContactPhone1        types.String  `tfsdk:"contact_phone_1"`
	ContactEmail1        types.String  `tfsdk:"contact_email_1"`
	Comments             types.String  `tfsdk:"comments"`
}

func (r *MeetingResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_meeting"
}

func (r *MeetingResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Meeting resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Meeting identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"service_body_id": schema.Int64Attribute{
				MarkdownDescription: "Service body identifier",
				Required:            true,
			},
			"format_ids": schema.ListAttribute{
				MarkdownDescription: "List of format identifiers",
				Required:            true,
				ElementType:         types.Int64Type,
			},
			"venue_type": schema.Int64Attribute{
				MarkdownDescription: "Venue type (1=in-person, 2=virtual, 3=hybrid)",
				Required:            true,
			},
			"temporarily_virtual": schema.BoolAttribute{
				MarkdownDescription: "Whether the meeting is temporarily virtual",
				Optional:            true,
			},
			"day": schema.Int64Attribute{
				MarkdownDescription: "Day of the week (0=Sunday, 1=Monday, etc.)",
				Required:            true,
			},
			"start_time": schema.StringAttribute{
				MarkdownDescription: "Meeting start time (HH:MM format)",
				Required:            true,
			},
			"duration": schema.StringAttribute{
				MarkdownDescription: "Meeting duration (HH:MM format)",
				Required:            true,
			},
			"time_zone": schema.StringAttribute{
				MarkdownDescription: "Time zone (e.g., America/New_York)",
				Optional:            true,
			},
			"latitude": schema.Float64Attribute{
				MarkdownDescription: "Latitude coordinate",
				Required:            true,
			},
			"longitude": schema.Float64Attribute{
				MarkdownDescription: "Longitude coordinate",
				Required:            true,
			},
			"published": schema.BoolAttribute{
				MarkdownDescription: "Whether the meeting is published",
				Required:            true,
			},
			"email": schema.StringAttribute{
				MarkdownDescription: "Meeting email",
				Optional:            true,
			},
			"world_id": schema.StringAttribute{
				MarkdownDescription: "World identifier",
				Optional:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Meeting name",
				Required:            true,
			},
			"location_text": schema.StringAttribute{
				MarkdownDescription: "Location text",
				Optional:            true,
			},
			"location_info": schema.StringAttribute{
				MarkdownDescription: "Location info",
				Optional:            true,
			},
			"location_street": schema.StringAttribute{
				MarkdownDescription: "Street address",
				Optional:            true,
			},
			"location_municipality": schema.StringAttribute{
				MarkdownDescription: "Municipality",
				Optional:            true,
			},
			"location_province": schema.StringAttribute{
				MarkdownDescription: "Province",
				Optional:            true,
			},
			"location_postal_code_1": schema.StringAttribute{
				MarkdownDescription: "Postal code",
				Optional:            true,
			},
			"location_nation": schema.StringAttribute{
				MarkdownDescription: "Nation",
				Optional:            true,
			},
			"virtual_meeting_link": schema.StringAttribute{
				MarkdownDescription: "Virtual meeting link",
				Optional:            true,
			},
			"contact_name_1": schema.StringAttribute{
				MarkdownDescription: "Primary contact name",
				Optional:            true,
			},
			"contact_phone_1": schema.StringAttribute{
				MarkdownDescription: "Primary contact phone",
				Optional:            true,
			},
			"contact_email_1": schema.StringAttribute{
				MarkdownDescription: "Primary contact email",
				Optional:            true,
			},
			"comments": schema.StringAttribute{
				MarkdownDescription: "Comments",
				Optional:            true,
			},
		},
	}
}

func (r *MeetingResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*BMTLClientData)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *BMTLClientData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *MeetingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *MeetingResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert format IDs
	var formatIds []int32
	for _, id := range data.FormatIds {
		formatIds = append(formatIds, int32(id.ValueInt64()))
	}

	// Convert model to API request
	createRequest := bmlt.MeetingCreate{
		ServiceBodyId:        int32(data.ServiceBodyId.ValueInt64()),
		FormatIds:            formatIds,
		VenueType:            int32(data.VenueType.ValueInt64()),
		TemporarilyVirtual:   data.TemporarilyVirtual.ValueBoolPointer(),
		Day:                  int32(data.Day.ValueInt64()),
		StartTime:            data.StartTime.ValueString(),
		Duration:             data.Duration.ValueString(),
		TimeZone:             data.TimeZone.ValueStringPointer(),
		Latitude:             float32(data.Latitude.ValueFloat64()),
		Longitude:            float32(data.Longitude.ValueFloat64()),
		Published:            data.Published.ValueBool(),
		Email:                data.Email.ValueStringPointer(),
		WorldId:              data.WorldId.ValueStringPointer(),
		Name:                 data.Name.ValueString(),
		LocationText:         data.LocationText.ValueStringPointer(),
		LocationInfo:         data.LocationInfo.ValueStringPointer(),
		LocationStreet:       data.LocationStreet.ValueStringPointer(),
		LocationMunicipality: data.LocationMunicipality.ValueStringPointer(),
		LocationProvince:     data.LocationProvince.ValueStringPointer(),
		LocationPostalCode1:  data.LocationPostalCode1.ValueStringPointer(),
		LocationNation:       data.LocationNation.ValueStringPointer(),
		VirtualMeetingLink:   data.VirtualMeetingLink.ValueStringPointer(),
		ContactName1:         data.ContactName1.ValueStringPointer(),
		ContactPhone1:        data.ContactPhone1.ValueStringPointer(),
		ContactEmail1:        data.ContactEmail1.ValueStringPointer(),
		Comments:             data.Comments.ValueStringPointer(),
	}

	// Create meeting
	meeting, httpResp, err := r.client.Client.RootServerAPI.CreateMeeting(r.client.Context).MeetingCreate(createRequest).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create meeting, got error: %s", err))
		return
	}

	if httpResp.StatusCode != 201 {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("API returned status %d", httpResp.StatusCode))
		return
	}

	// Map response back to model
	data.Id = types.StringValue(strconv.Itoa(int(meeting.Id)))

	// Update all fields from response
	r.updateModelFromMeeting(data, meeting)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *MeetingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *MeetingResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.ParseInt(data.Id.ValueString(), 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse meeting ID: %s", err))
		return
	}

	meeting, httpResp, err := r.client.Client.RootServerAPI.GetMeeting(r.client.Context, id).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read meeting, got error: %s", err))
		return
	}

	if httpResp.StatusCode == 404 {
		resp.State.RemoveResource(ctx)
		return
	}

	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("API returned status %d", httpResp.StatusCode))
		return
	}

	r.updateModelFromMeeting(data, meeting)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *MeetingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *MeetingResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.ParseInt(data.Id.ValueString(), 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse meeting ID: %s", err))
		return
	}

	// Convert format IDs
	var formatIds []int32
	for _, id := range data.FormatIds {
		formatIds = append(formatIds, int32(id.ValueInt64()))
	}

	updateRequest := bmlt.MeetingUpdate{
		ServiceBodyId:        int32(data.ServiceBodyId.ValueInt64()),
		FormatIds:            formatIds,
		VenueType:            int32(data.VenueType.ValueInt64()),
		TemporarilyVirtual:   data.TemporarilyVirtual.ValueBoolPointer(),
		Day:                  int32(data.Day.ValueInt64()),
		StartTime:            data.StartTime.ValueString(),
		Duration:             data.Duration.ValueString(),
		TimeZone:             data.TimeZone.ValueStringPointer(),
		Latitude:             float32(data.Latitude.ValueFloat64()),
		Longitude:            float32(data.Longitude.ValueFloat64()),
		Published:            data.Published.ValueBool(),
		Email:                data.Email.ValueStringPointer(),
		WorldId:              data.WorldId.ValueStringPointer(),
		Name:                 data.Name.ValueString(),
		LocationText:         data.LocationText.ValueStringPointer(),
		LocationInfo:         data.LocationInfo.ValueStringPointer(),
		LocationStreet:       data.LocationStreet.ValueStringPointer(),
		LocationMunicipality: data.LocationMunicipality.ValueStringPointer(),
		LocationProvince:     data.LocationProvince.ValueStringPointer(),
		LocationPostalCode1:  data.LocationPostalCode1.ValueStringPointer(),
		LocationNation:       data.LocationNation.ValueStringPointer(),
		VirtualMeetingLink:   data.VirtualMeetingLink.ValueStringPointer(),
		ContactName1:         data.ContactName1.ValueStringPointer(),
		ContactPhone1:        data.ContactPhone1.ValueStringPointer(),
		ContactEmail1:        data.ContactEmail1.ValueStringPointer(),
		Comments:             data.Comments.ValueStringPointer(),
	}

	httpResp, err := r.client.Client.RootServerAPI.UpdateMeeting(r.client.Context, id).MeetingUpdate(updateRequest).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update meeting, got error: %s", err))
		return
	}

	if httpResp.StatusCode != 204 {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("API returned status %d", httpResp.StatusCode))
		return
	}

	// Re-read the meeting to ensure state is consistent with server
	meetingId, err := strconv.ParseInt(data.Id.ValueString(), 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse meeting ID: %s", err))
		return
	}

	updatedMeeting, httpResp, err := r.client.Client.RootServerAPI.GetMeeting(r.client.Context, meetingId).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read updated meeting, got error: %s", err))
		return
	}

	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("API returned status %d when reading updated meeting", httpResp.StatusCode))
		return
	}

	// Update all fields from the server response
	r.updateModelFromMeeting(data, updatedMeeting)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *MeetingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *MeetingResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.ParseInt(data.Id.ValueString(), 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse meeting ID: %s", err))
		return
	}

	httpResp, err := r.client.Client.RootServerAPI.DeleteMeeting(r.client.Context, id).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete meeting, got error: %s", err))
		return
	}

	if httpResp.StatusCode != 204 {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("API returned status %d", httpResp.StatusCode))
		return
	}
}

func (r *MeetingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}


// Helper function to update model from API response
func (r *MeetingResource) updateModelFromMeeting(data *MeetingResourceModel, meeting *bmlt.Meeting) {
	data.ServiceBodyId = types.Int64Value(int64(meeting.ServiceBodyId))
	data.VenueType = types.Int64Value(int64(meeting.VenueType))
	data.TemporarilyVirtual = types.BoolValue(meeting.TemporarilyVirtual)
	data.Day = types.Int64Value(int64(meeting.Day))
	data.StartTime = types.StringValue(meeting.StartTime)
	data.Duration = types.StringValue(meeting.Duration)
	data.TimeZone = nullableString(meeting.TimeZone)
	data.Latitude = types.Float64Value(float64(meeting.Latitude))
	data.Longitude = types.Float64Value(float64(meeting.Longitude))
	data.Published = types.BoolValue(meeting.Published)
	data.Email = nullableString(meeting.Email)
	data.WorldId = nullableString(meeting.WorldId)
	data.Name = types.StringValue(meeting.Name)
	data.LocationText = types.StringPointerValue(meeting.LocationText)
	data.LocationInfo = types.StringPointerValue(meeting.LocationInfo)
	data.LocationStreet = types.StringPointerValue(meeting.LocationStreet)
	data.LocationMunicipality = types.StringPointerValue(meeting.LocationMunicipality)
	data.LocationProvince = types.StringPointerValue(meeting.LocationProvince)
	data.LocationPostalCode1 = types.StringPointerValue(meeting.LocationPostalCode1)
	data.LocationNation = types.StringPointerValue(meeting.LocationNation)
	data.VirtualMeetingLink = types.StringPointerValue(meeting.VirtualMeetingLink)
	data.ContactName1 = types.StringPointerValue(meeting.ContactName1)
	data.ContactPhone1 = types.StringPointerValue(meeting.ContactPhone1)
	data.ContactEmail1 = types.StringPointerValue(meeting.ContactEmail1)
	data.Comments = types.StringPointerValue(meeting.Comments)

	// Handle format IDs
	var responseFormatIds []types.Int64
	for _, formatId := range meeting.FormatIds {
		responseFormatIds = append(responseFormatIds, types.Int64Value(int64(formatId)))
	}
	data.FormatIds = responseFormatIds
}
