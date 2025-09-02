package provider

import (
	"context"
	"fmt"

	"github.com/bmlt-enabled/bmlt-server-go-client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &ServiceBodyDataSource{}

func NewServiceBodyDataSource() datasource.DataSource {
	return &ServiceBodyDataSource{}
}

type ServiceBodyDataSource struct {
	client *BMTLClientData
}

type ServiceBodyDataSourceModel struct {
	Id              types.Int64   `tfsdk:"id"`
	ServiceBodyId   types.Int64   `tfsdk:"service_body_id"`
	Name            types.String  `tfsdk:"name"`
	ParentId        types.Int64   `tfsdk:"parent_id"`
	Description     types.String  `tfsdk:"description"`
	Type            types.String  `tfsdk:"type"`
	AdminUserId     types.Int64   `tfsdk:"admin_user_id"`
	AssignedUserIds []types.Int64 `tfsdk:"assigned_user_ids"`
	Url             types.String  `tfsdk:"url"`
	Helpline        types.String  `tfsdk:"helpline"`
	Email           types.String  `tfsdk:"email"`
	WorldId         types.String  `tfsdk:"world_id"`
}

func (d *ServiceBodyDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service_body"
}

func (d *ServiceBodyDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Service body data source allows you to retrieve information about a specific service body by either service body ID or name.",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "Service body identifier (computed from service_body_id or resolved from name)",
				Computed:            true,
			},
			"service_body_id": schema.Int64Attribute{
				MarkdownDescription: "Service body identifier to look up (mutually exclusive with name)",
				Optional:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Service body name to look up (mutually exclusive with service_body_id)",
				Optional:            true,
			},
			"parent_id": schema.Int64Attribute{
				MarkdownDescription: "Parent service body identifier",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Service body description",
				Computed:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Service body type",
				Computed:            true,
			},
			"admin_user_id": schema.Int64Attribute{
				MarkdownDescription: "Admin user identifier",
				Computed:            true,
			},
			"assigned_user_ids": schema.ListAttribute{
				MarkdownDescription: "List of assigned user identifiers",
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"url": schema.StringAttribute{
				MarkdownDescription: "Service body URL",
				Computed:            true,
			},
			"helpline": schema.StringAttribute{
				MarkdownDescription: "Service body helpline",
				Computed:            true,
			},
			"email": schema.StringAttribute{
				MarkdownDescription: "Service body email",
				Computed:            true,
			},
			"world_id": schema.StringAttribute{
				MarkdownDescription: "World identifier",
				Computed:            true,
			},
		},
	}
}

func (d *ServiceBodyDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ServiceBodyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ServiceBodyDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate that exactly one of service_body_id or name is provided
	hasServiceBodyId := !data.ServiceBodyId.IsNull()
	hasName := !data.Name.IsNull()

	if !hasServiceBodyId && !hasName {
		resp.Diagnostics.AddError(
			"Missing Required Argument",
			"Either 'service_body_id' or 'name' must be provided.",
		)
		return
	}

	if hasServiceBodyId && hasName {
		resp.Diagnostics.AddError(
			"Conflicting Arguments",
			"Cannot specify both 'service_body_id' and 'name'. Please provide only one.",
		)
		return
	}

	// If service_body_id is provided, fetch service body directly
	if hasServiceBodyId {
		serviceBodyId := data.ServiceBodyId.ValueInt64()
		serviceBody, httpResp, err := d.client.Client.RootServerAPI.GetServiceBody(d.client.Context, serviceBodyId).Execute()
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read service body, got error: %s", err))
			return
		}

		if httpResp.StatusCode == HTTPStatusNotFound {
			resp.Diagnostics.AddError("Service Body Not Found", fmt.Sprintf("Service body with ID %d not found", serviceBodyId))
			return
		}

		if httpResp.StatusCode != HTTPStatusOK {
			resp.Diagnostics.AddError("API Error", fmt.Sprintf("API returned status %d", httpResp.StatusCode))
			return
		}

		// Map response to model
		d.mapServiceBodyToModel(&data, serviceBody)
	} else {
		// If name is provided, fetch all service bodies and filter
		targetName := data.Name.ValueString()
		serviceBodies, httpResp, err := d.client.Client.RootServerAPI.GetServiceBodies(d.client.Context).Execute()
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read service bodies, got error: %s", err))
			return
		}

		if httpResp.StatusCode != HTTPStatusOK {
			resp.Diagnostics.AddError("API Error", fmt.Sprintf("API returned status %d", httpResp.StatusCode))
			return
		}

		// Find service body by name
		var foundServiceBody *bmlt.ServiceBody
		for _, serviceBody := range serviceBodies {
			if serviceBody.Name == targetName {
				foundServiceBody = &serviceBody
				break
			}
		}

		if foundServiceBody == nil {
			resp.Diagnostics.AddError("Service Body Not Found", fmt.Sprintf("Service body with name '%s' not found", targetName))
			return
		}

		// Map response to model
		d.mapServiceBodyToModel(&data, foundServiceBody)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Helper function to map service body API response to model
func (d *ServiceBodyDataSource) mapServiceBodyToModel(data *ServiceBodyDataSourceModel, serviceBody *bmlt.ServiceBody) {
	data.Id = types.Int64Value(int64(serviceBody.Id))
	data.ServiceBodyId = types.Int64Value(int64(serviceBody.Id))
	data.Name = types.StringValue(serviceBody.Name)
	data.Description = types.StringValue(serviceBody.Description)
	data.Type = types.StringValue(serviceBody.Type)
	data.AdminUserId = types.Int64Value(int64(serviceBody.AdminUserId))
	data.Url = types.StringValue(serviceBody.Url)
	data.Helpline = types.StringValue(serviceBody.Helpline)
	data.Email = types.StringValue(serviceBody.Email)
	data.WorldId = types.StringValue(serviceBody.WorldId)

	// Handle nullable ParentId
	if serviceBody.ParentId.IsSet() && serviceBody.ParentId.Get() != nil {
		data.ParentId = types.Int64Value(int64(*serviceBody.ParentId.Get()))
	} else {
		data.ParentId = types.Int64Null()
	}

	// Handle assigned user IDs
	var assignedUserIds []types.Int64
	for _, userId := range serviceBody.AssignedUserIds {
		assignedUserIds = append(assignedUserIds, types.Int64Value(int64(userId)))
	}
	data.AssignedUserIds = assignedUserIds
}
