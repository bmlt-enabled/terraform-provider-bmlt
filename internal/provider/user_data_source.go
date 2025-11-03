package provider

import (
	"context"
	"fmt"

	"github.com/bmlt-enabled/bmlt-server-go-client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &UserDataSource{}

func NewUserDataSource() datasource.DataSource {
	return &UserDataSource{}
}

type UserDataSource struct {
	client *BMTLClientData
}

type UserDataSourceModel struct {
	Id          types.Int64  `tfsdk:"id"`
	UserId      types.Int64  `tfsdk:"user_id"`
	Username    types.String `tfsdk:"username"`
	Type        types.String `tfsdk:"type"`
	DisplayName types.String `tfsdk:"display_name"`
	Description types.String `tfsdk:"description"`
	Email       types.String `tfsdk:"email"`
	OwnerId     types.Int64  `tfsdk:"owner_id"`
	LastLoginAt types.String `tfsdk:"last_login_at"`
}

func (d *UserDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (d *UserDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "User data source allows you to retrieve information about a specific user by either user ID or username.",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "User identifier (computed from user_id or resolved from username)",
				Computed:            true,
			},
			"user_id": schema.Int64Attribute{
				MarkdownDescription: "User identifier to look up (mutually exclusive with username)",
				Optional:            true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "Username to look up (mutually exclusive with user_id)",
				Optional:            true,
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
			"last_login_at": schema.StringAttribute{
				MarkdownDescription: "Last login timestamp (computed from last token generation)",
				Computed:            true,
			},
		},
	}
}

func (d *UserDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *UserDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data UserDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate that exactly one of user_id or username is provided
	hasUserId := !data.UserId.IsNull()
	hasUsername := !data.Username.IsNull()

	if !hasUserId && !hasUsername {
		resp.Diagnostics.AddError(
			"Missing Required Argument",
			"Either 'user_id' or 'username' must be provided.",
		)
		return
	}

	if hasUserId && hasUsername {
		resp.Diagnostics.AddError(
			"Conflicting Arguments",
			"Cannot specify both 'user_id' and 'username'. Please provide only one.",
		)
		return
	}

	// If user_id is provided, fetch user directly
	if hasUserId {
		userId := data.UserId.ValueInt64()
		user, httpResp, err := d.client.Client.RootServerAPI.GetUser(d.client.Context, userId).Execute()
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read user, got error: %s", err))
			return
		}

		if httpResp.StatusCode == HTTPStatusNotFound {
			resp.Diagnostics.AddError("User Not Found", fmt.Sprintf("User with ID %d not found", userId))
			return
		}

		if httpResp.StatusCode != HTTPStatusOK {
			resp.Diagnostics.AddError("API Error", fmt.Sprintf("API returned status %d", httpResp.StatusCode))
			return
		}

		// Map response to model
		data.Id = types.Int64Value(int64(user.Id))
		data.UserId = types.Int64Value(int64(user.Id))
		data.Username = types.StringValue(user.Username)
		data.Type = types.StringValue(user.Type)
		data.DisplayName = types.StringValue(user.DisplayName)
		data.Description = types.StringValue(user.Description)
		data.Email = types.StringValue(user.Email)
		data.OwnerId = types.Int64Value(int64(user.OwnerId))
		data.LastLoginAt = nullableString(user.LastLoginAt)
	} else {
		// If username is provided, fetch all users and filter
		targetUsername := data.Username.ValueString()
		users, httpResp, err := d.client.Client.RootServerAPI.GetUsers(d.client.Context).Execute()
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read users, got error: %s", err))
			return
		}

		if httpResp.StatusCode != HTTPStatusOK {
			resp.Diagnostics.AddError("API Error", fmt.Sprintf("API returned status %d", httpResp.StatusCode))
			return
		}

		// Find user by username
		var foundUser *bmlt.User
		for _, user := range users {
			if user.Username == targetUsername {
				foundUser = &user
				break
			}
		}

		if foundUser == nil {
			resp.Diagnostics.AddError("User Not Found", fmt.Sprintf("User with username '%s' not found", targetUsername))
			return
		}

		// Map response to model
		data.Id = types.Int64Value(int64(foundUser.Id))
		data.UserId = types.Int64Value(int64(foundUser.Id))
		data.Username = types.StringValue(foundUser.Username)
		data.Type = types.StringValue(foundUser.Type)
		data.DisplayName = types.StringValue(foundUser.DisplayName)
		data.Description = types.StringValue(foundUser.Description)
		data.Email = types.StringValue(foundUser.Email)
		data.OwnerId = types.Int64Value(int64(foundUser.OwnerId))
		data.LastLoginAt = nullableString(foundUser.LastLoginAt)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
