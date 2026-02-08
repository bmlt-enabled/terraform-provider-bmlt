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

var _ resource.Resource = &ServiceBodyResource{}
var _ resource.ResourceWithImportState = &ServiceBodyResource{}

func NewServiceBodyResource() resource.Resource {
	return &ServiceBodyResource{}
}

type ServiceBodyResource struct {
	client *BMTLClientData
}

type ServiceBodyResourceModel struct {
	Id              types.String  `tfsdk:"id"`
	ParentId        types.Int64   `tfsdk:"parent_id"`
	Name            types.String  `tfsdk:"name"`
	Description     types.String  `tfsdk:"description"`
	Type            types.String  `tfsdk:"type"`
	AdminUserId     types.Int64   `tfsdk:"admin_user_id"`
	AssignedUserIds []types.Int64 `tfsdk:"assigned_user_ids"`
	Url             types.String  `tfsdk:"url"`
	Helpline        types.String  `tfsdk:"helpline"`
	Email           types.String  `tfsdk:"email"`
	WorldId         types.String  `tfsdk:"world_id"`
	ForceDelete     types.Bool    `tfsdk:"force_delete"`
}

func (r *ServiceBodyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service_body"
}

func (r *ServiceBodyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Service body resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Service body identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"parent_id": schema.Int64Attribute{
				MarkdownDescription: "Parent service body identifier",
				Optional:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Service body name",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Service body description",
				Required:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Service body type",
				Required:            true,
			},
			"admin_user_id": schema.Int64Attribute{
				MarkdownDescription: "Admin user identifier",
				Required:            true,
			},
			"assigned_user_ids": schema.ListAttribute{
				MarkdownDescription: "List of assigned user identifiers",
				Required:            true,
				ElementType:         types.Int64Type,
			},
			"url": schema.StringAttribute{
				MarkdownDescription: "Service body URL",
				Optional:            true,
			},
			"helpline": schema.StringAttribute{
				MarkdownDescription: "Service body helpline",
				Optional:            true,
			},
			"email": schema.StringAttribute{
				MarkdownDescription: "Service body email",
				Optional:            true,
			},
			"world_id": schema.StringAttribute{
				MarkdownDescription: "World identifier",
				Optional:            true,
			},
			"force_delete": schema.BoolAttribute{
				MarkdownDescription: "Force delete the service body even if it has associated meetings",
				Optional:            true,
			},
		},
	}
}

func (r *ServiceBodyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ServiceBodyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *ServiceBodyResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert assigned user IDs
	var assignedUserIds []int32
	for _, id := range data.AssignedUserIds {
		assignedUserIds = append(assignedUserIds, safeInt64ToInt32(id.ValueInt64()))
	}

	// Create NullableInt32 for ParentId
	var parentId bmlt.NullableInt32
	if !data.ParentId.IsNull() {
		parentId.Set(bmlt.PtrInt32(safeInt64ToInt32(data.ParentId.ValueInt64())))
	}

	// Convert model to API request
	createRequest := bmlt.ServiceBodyCreate{
		ParentId:        parentId,
		Name:            data.Name.ValueString(),
		Description:     data.Description.ValueString(),
		Type:            data.Type.ValueString(),
		AdminUserId:     safeInt64ToInt32(data.AdminUserId.ValueInt64()),
		AssignedUserIds: assignedUserIds,
		Url:             data.Url.ValueStringPointer(),
		Helpline:        data.Helpline.ValueStringPointer(),
		Email:           data.Email.ValueStringPointer(),
		WorldId:         data.WorldId.ValueStringPointer(),
	}

	// Create service body
	serviceBody, httpResp, err := r.client.Client.RootServerAPI.CreateServiceBody(r.client.Context).ServiceBodyCreate(createRequest).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create service body, got error: %s", err))
		return
	}

	if httpResp.StatusCode != HTTPStatusCreated {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("API returned status %d", httpResp.StatusCode))
		return
	}

	// Map response back to model
	data.Id = types.StringValue(strconv.Itoa(int(serviceBody.Id)))
	r.updateModelFromServiceBody(data, serviceBody)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ServiceBodyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *ServiceBodyResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.ParseInt(data.Id.ValueString(), 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse service body ID: %s", err))
		return
	}

	serviceBody, httpResp, err := r.client.Client.RootServerAPI.GetServiceBody(r.client.Context, id).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read service body, got error: %s", err))
		return
	}

	if httpResp.StatusCode == HTTPStatusNotFound {
		resp.State.RemoveResource(ctx)
		return
	}

	if httpResp.StatusCode != HTTPStatusOK {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("API returned status %d", httpResp.StatusCode))
		return
	}

	r.updateModelFromServiceBody(data, serviceBody)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ServiceBodyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *ServiceBodyResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.ParseInt(data.Id.ValueString(), 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse service body ID: %s", err))
		return
	}

	// Convert assigned user IDs
	var assignedUserIds []int32
	for _, id := range data.AssignedUserIds {
		assignedUserIds = append(assignedUserIds, safeInt64ToInt32(id.ValueInt64()))
	}

	// Create NullableInt32 for ParentId
	var parentId bmlt.NullableInt32
	if !data.ParentId.IsNull() {
		parentId.Set(bmlt.PtrInt32(safeInt64ToInt32(data.ParentId.ValueInt64())))
	}

	updateRequest := bmlt.ServiceBodyUpdate{
		ParentId:        parentId,
		Name:            data.Name.ValueString(),
		Description:     data.Description.ValueString(),
		Type:            data.Type.ValueString(),
		AdminUserId:     safeInt64ToInt32(data.AdminUserId.ValueInt64()),
		AssignedUserIds: assignedUserIds,
		Url:             data.Url.ValueStringPointer(),
		Helpline:        data.Helpline.ValueStringPointer(),
		Email:           data.Email.ValueStringPointer(),
		WorldId:         data.WorldId.ValueStringPointer(),
	}

	httpResp, err := r.client.Client.RootServerAPI.UpdateServiceBody(r.client.Context, id).ServiceBodyUpdate(updateRequest).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update service body, got error: %s", err))
		return
	}

	if httpResp.StatusCode != HTTPStatusNoContent {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("API returned status %d", httpResp.StatusCode))
		return
	}

	// Re-read the service body to ensure state is consistent with server
	serviceBodyId, err := strconv.ParseInt(data.Id.ValueString(), 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse service body ID: %s", err))
		return
	}

	updatedServiceBody, httpResp, err := r.client.Client.RootServerAPI.GetServiceBody(r.client.Context, serviceBodyId).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read updated service body, got error: %s", err))
		return
	}

	if httpResp.StatusCode != HTTPStatusOK {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("API returned status %d when reading updated service body", httpResp.StatusCode))
		return
	}

	// Update all fields from the server response
	r.updateModelFromServiceBody(data, updatedServiceBody)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ServiceBodyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *ServiceBodyResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.ParseInt(data.Id.ValueString(), 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse service body ID: %s", err))
		return
	}

	deleteReq := r.client.Client.RootServerAPI.DeleteServiceBody(r.client.Context, id)
	if data.ForceDelete.ValueBool() {
		deleteReq = deleteReq.Force("true")
	}
	httpResp, err := deleteReq.Execute()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete service body, got error: %s", err))
		return
	}

	if httpResp.StatusCode != HTTPStatusNoContent {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("API returned status %d", httpResp.StatusCode))
		return
	}
}

func (r *ServiceBodyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Helper function to update model from API response
func (r *ServiceBodyResource) updateModelFromServiceBody(data *ServiceBodyResourceModel, serviceBody *bmlt.ServiceBody) {
	// Handle nullable ParentId
	if serviceBody.ParentId.IsSet() && serviceBody.ParentId.Get() != nil {
		data.ParentId = types.Int64Value(int64(*serviceBody.ParentId.Get()))
	} else {
		data.ParentId = types.Int64Null()
	}

	data.Name = types.StringValue(serviceBody.Name)
	data.Description = types.StringValue(serviceBody.Description)
	data.Type = types.StringValue(serviceBody.Type)
	data.AdminUserId = types.Int64Value(int64(serviceBody.AdminUserId))
	data.Url = nullableString(serviceBody.Url)
	data.Helpline = nullableString(serviceBody.Helpline)
	data.Email = nullableString(serviceBody.Email)
	data.WorldId = nullableString(serviceBody.WorldId)

	// Handle assigned user IDs
	var assignedUserIds []types.Int64
	for _, userId := range serviceBody.AssignedUserIds {
		assignedUserIds = append(assignedUserIds, types.Int64Value(int64(userId)))
	}
	data.AssignedUserIds = assignedUserIds
}
