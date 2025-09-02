package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &ServiceBodiesDataSource{}

func NewServiceBodiesDataSource() datasource.DataSource {
	return &ServiceBodiesDataSource{}
}

type ServiceBodiesDataSource struct {
	client *BMTLClientData
}

type ServiceBodiesDataSourceModel struct {
	ServiceBodies []ServiceBodyModel `tfsdk:"service_bodies"`
	Id            types.String       `tfsdk:"id"`
}

type ServiceBodyModel struct {
	Id              types.Int64   `tfsdk:"id"`
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
}

func (d *ServiceBodiesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service_bodies"
}

func (d *ServiceBodiesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Service bodies data source allows you to retrieve information about service bodies.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Placeholder identifier for the data source.",
				Computed:            true,
			},
			"service_bodies": schema.ListNestedAttribute{
				MarkdownDescription: "List of service bodies",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							MarkdownDescription: "Service body identifier",
							Computed:            true,
						},
						"parent_id": schema.Int64Attribute{
							MarkdownDescription: "Parent service body identifier",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Service body name",
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
				},
			},
		},
	}
}

func (d *ServiceBodiesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ServiceBodiesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ServiceBodiesDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get service bodies from the API
	serviceBodies, httpResp, err := d.client.Client.RootServerAPI.GetServiceBodies(d.client.Context).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read service bodies, got error: %s", err))
		return
	}

	if httpResp.StatusCode != HTTPStatusOK {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("API returned status %d", httpResp.StatusCode))
		return
	}

	// Map response body to model
	for _, serviceBody := range serviceBodies {
		serviceBodyModel := ServiceBodyModel{
			Id:          types.Int64Value(int64(serviceBody.Id)),
			Name:        types.StringValue(serviceBody.Name),
			Description: types.StringValue(serviceBody.Description),
			Type:        types.StringValue(serviceBody.Type),
			AdminUserId: types.Int64Value(int64(serviceBody.AdminUserId)),
			Url:         types.StringValue(serviceBody.Url),
			Helpline:    types.StringValue(serviceBody.Helpline),
			Email:       types.StringValue(serviceBody.Email),
			WorldId:     types.StringValue(serviceBody.WorldId),
		}

		// Handle nullable ParentId
		if serviceBody.ParentId.IsSet() && serviceBody.ParentId.Get() != nil {
			serviceBodyModel.ParentId = types.Int64Value(int64(*serviceBody.ParentId.Get()))
		} else {
			serviceBodyModel.ParentId = types.Int64Null()
		}

		// Handle assigned user IDs
		var assignedUserIds []types.Int64
		for _, userId := range serviceBody.AssignedUserIds {
			assignedUserIds = append(assignedUserIds, types.Int64Value(int64(userId)))
		}
		serviceBodyModel.AssignedUserIds = assignedUserIds

		data.ServiceBodies = append(data.ServiceBodies, serviceBodyModel)
	}

	data.Id = types.StringValue("placeholder")
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
