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

var _ resource.Resource = &UserResource{}
var _ resource.ResourceWithImportState = &UserResource{}

func NewUserResource() resource.Resource {
	return &UserResource{}
}

type UserResource struct {
	client *BMTLClientData
}

type UserResourceModel struct {
	Id          types.String `tfsdk:"id"`
	Username    types.String `tfsdk:"username"`
	Password    types.String `tfsdk:"password"`
	Type        types.String `tfsdk:"type"`
	DisplayName types.String `tfsdk:"display_name"`
	Description types.String `tfsdk:"description"`
	Email       types.String `tfsdk:"email"`
	OwnerId     types.Int64  `tfsdk:"owner_id"`
}

func (r *UserResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (r *UserResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "User resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "User identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "Username",
				Required:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "User password",
				Optional:            true,
				Sensitive:           true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "User type",
				Required:            true,
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "Display name",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "User description",
				Optional:            true,
			},
			"email": schema.StringAttribute{
				MarkdownDescription: "User email",
				Optional:            true,
			},
			"owner_id": schema.Int64Attribute{
				MarkdownDescription: "Owner identifier",
				Optional:            true,
			},
		},
	}
}

func (r *UserResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *UserResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *UserResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert model to API request
	createRequest := bmlt.UserCreate{
		Username:    data.Username.ValueString(),
		Password:    data.Password.ValueString(),
		Type:        data.Type.ValueString(),
		DisplayName: data.DisplayName.ValueString(),
		Description: data.Description.ValueStringPointer(),
		Email:       data.Email.ValueStringPointer(),
		OwnerId:     nil, // Handle OwnerId separately
	}

	// Handle optional OwnerId
	if !data.OwnerId.IsNull() {
		createRequest.OwnerId = bmlt.PtrInt32(safeInt64ToInt32(data.OwnerId.ValueInt64()))
	}

	// Create user
	user, httpResp, err := r.client.Client.RootServerAPI.CreateUser(r.client.Context).UserCreate(createRequest).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create user, got error: %s", err))
		return
	}

	if httpResp.StatusCode != HTTPStatusCreated {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("API returned status %d", httpResp.StatusCode))
		return
	}

	// Map response back to model
	data.Id = types.StringValue(strconv.Itoa(int(user.Id)))
	r.updateModelFromUser(data, user)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *UserResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *UserResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.ParseInt(data.Id.ValueString(), 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse user ID: %s", err))
		return
	}

	user, httpResp, err := r.client.Client.RootServerAPI.GetUser(r.client.Context, id).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read user, got error: %s", err))
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

	r.updateModelFromUser(data, user)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *UserResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *UserResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.ParseInt(data.Id.ValueString(), 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse user ID: %s", err))
		return
	}

	updateRequest := bmlt.UserUpdate{
		Username:    data.Username.ValueString(),
		Type:        data.Type.ValueString(),
		DisplayName: data.DisplayName.ValueString(),
		Description: data.Description.ValueStringPointer(),
		Email:       data.Email.ValueStringPointer(),
		Password:    data.Password.ValueStringPointer(),
		OwnerId:     nil, // Handle OwnerId separately
	}

	// Handle optional OwnerId
	if !data.OwnerId.IsNull() {
		updateRequest.OwnerId = bmlt.PtrInt32(safeInt64ToInt32(data.OwnerId.ValueInt64()))
	}

	httpResp, err := r.client.Client.RootServerAPI.UpdateUser(r.client.Context, id).UserUpdate(updateRequest).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update user, got error: %s", err))
		return
	}

	if httpResp.StatusCode != HTTPStatusNoContent {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("API returned status %d", httpResp.StatusCode))
		return
	}

	// Re-read the user to ensure state is consistent with server
	userId, err := strconv.ParseInt(data.Id.ValueString(), 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse user ID: %s", err))
		return
	}

	updatedUser, httpResp, err := r.client.Client.RootServerAPI.GetUser(r.client.Context, userId).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read updated user, got error: %s", err))
		return
	}

	if httpResp.StatusCode != HTTPStatusOK {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("API returned status %d when reading updated user", httpResp.StatusCode))
		return
	}

	// Update all fields from the server response
	r.updateModelFromUser(data, updatedUser)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *UserResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *UserResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.ParseInt(data.Id.ValueString(), 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse user ID: %s", err))
		return
	}

	httpResp, err := r.client.Client.RootServerAPI.DeleteUser(r.client.Context, id).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete user, got error: %s", err))
		return
	}

	if httpResp.StatusCode != HTTPStatusNoContent {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("API returned status %d", httpResp.StatusCode))
		return
	}
}

func (r *UserResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Helper function to update model from API response
func (r *UserResource) updateModelFromUser(data *UserResourceModel, user *bmlt.User) {
	data.Username = types.StringValue(user.Username)
	data.Type = types.StringValue(user.Type)
	data.DisplayName = types.StringValue(user.DisplayName)
	data.Description = nullableString(user.Description)
	data.Email = nullableString(user.Email)
	data.OwnerId = types.Int64Value(int64(user.OwnerId))
	// Note: Password is not returned from API for security reasons
}
