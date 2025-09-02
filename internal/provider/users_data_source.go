package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &UsersDataSource{}

func NewUsersDataSource() datasource.DataSource {
	return &UsersDataSource{}
}

type UsersDataSource struct {
	client *BMTLClientData
}

type UsersDataSourceModel struct {
	Users []UserModel  `tfsdk:"users"`
	Id    types.String `tfsdk:"id"`
}

type UserModel struct {
	Id          types.Int64  `tfsdk:"id"`
	Username    types.String `tfsdk:"username"`
	Type        types.String `tfsdk:"type"`
	DisplayName types.String `tfsdk:"display_name"`
	Description types.String `tfsdk:"description"`
	Email       types.String `tfsdk:"email"`
	OwnerId     types.Int64  `tfsdk:"owner_id"`
}

func (d *UsersDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_users"
}

func (d *UsersDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Users data source allows you to retrieve information about users.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Placeholder identifier for the data source.",
				Computed:            true,
			},
			"users": schema.ListNestedAttribute{
				MarkdownDescription: "List of users",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							MarkdownDescription: "User identifier",
							Computed:            true,
						},
						"username": schema.StringAttribute{
							MarkdownDescription: "Username",
							Computed:            true,
						},
						"type": schema.StringAttribute{
							MarkdownDescription: "User type",
							Computed:            true,
						},
						"display_name": schema.StringAttribute{
							MarkdownDescription: "Display name",
							Computed:            true,
						},
						"description": schema.StringAttribute{
							MarkdownDescription: "User description",
							Computed:            true,
						},
						"email": schema.StringAttribute{
							MarkdownDescription: "User email",
							Computed:            true,
						},
						"owner_id": schema.Int64Attribute{
							MarkdownDescription: "Owner identifier",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *UsersDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *UsersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data UsersDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get users from the API
	users, httpResp, err := d.client.Client.RootServerAPI.GetUsers(d.client.Context).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read users, got error: %s", err))
		return
	}

	if httpResp.StatusCode != HTTPStatusOK {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("API returned status %d", httpResp.StatusCode))
		return
	}

	// Map response body to model
	for _, user := range users {
		userModel := UserModel{
			Id:          types.Int64Value(int64(user.Id)),
			Username:    types.StringValue(user.Username),
			Type:        types.StringValue(user.Type),
			DisplayName: types.StringValue(user.DisplayName),
			Description: types.StringValue(user.Description),
			Email:       types.StringValue(user.Email),
			OwnerId:     types.Int64Value(int64(user.OwnerId)),
		}

		data.Users = append(data.Users, userModel)
	}

	data.Id = types.StringValue("placeholder")
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
