package provider

import (
	"context"
	"os"
	"strings"

	"github.com/bmlt-enabled/bmlt-server-go-client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"golang.org/x/oauth2"
)

// Ensure BMTProvider satisfies various provider interfaces.
var _ provider.Provider = &BMTProvider{}

// BMTProvider defines the provider implementation.
type BMTProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// BMTProviderModel describes the provider data model.
type BMTProviderModel struct {
	Host        types.String `tfsdk:"host"`
	Username    types.String `tfsdk:"username"`
	Password    types.String `tfsdk:"password"`
	AccessToken types.String `tfsdk:"access_token"`
}

func (p *BMTProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "bmlt"
	resp.Version = p.version
}

func (p *BMTProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				MarkdownDescription: "BMLT server host URL (e.g., https://example.com/main_server)",
				Optional:            true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "Username for BMLT server authentication",
				Optional:            true,
				Sensitive:           true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "Password for BMLT server authentication",
				Optional:            true,
				Sensitive:           true,
			},
			"access_token": schema.StringAttribute{
				MarkdownDescription: "OAuth2 access token for BMLT server authentication (alternative to username/password)",
				Optional:            true,
				Sensitive:           true,
			},
		},
	}
}

func (p *BMTProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data BMTProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Check for unknown values
	if data.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown BMLT API Host",
			"The provider cannot create the BMLT API client as there is an unknown configuration value for the BMLT API host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the BMLT_HOST environment variable.",
		)
	}

	if data.Username.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Unknown BMLT API Username",
			"The provider cannot create the BMLT API client as there is an unknown configuration value for the BMLT API username.",
		)
	}

	if data.Password.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Unknown BMLT API Password",
			"The provider cannot create the BMLT API client as there is an unknown configuration value for the BMLT API password.",
		)
	}

	if data.AccessToken.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("access_token"),
			"Unknown BMLT API Access Token",
			"The provider cannot create the BMLT API client as there is an unknown configuration value for the BMLT API access token.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Get values from config or environment variables
	host := os.Getenv("BMLT_HOST")
	username := os.Getenv("BMLT_USERNAME")
	password := os.Getenv("BMLT_PASSWORD")
	accessToken := os.Getenv("BMLT_ACCESS_TOKEN")

	if !data.Host.IsNull() {
		host = data.Host.ValueString()
	}

	if !data.Username.IsNull() {
		username = data.Username.ValueString()
	}

	if !data.Password.IsNull() {
		password = data.Password.ValueString()
	}

	if !data.AccessToken.IsNull() {
		accessToken = data.AccessToken.ValueString()
	}

	// Validate required host
	if host == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Missing BMLT API Host",
			"The provider requires a BMLT server host URL. Set the host value in the configuration or use the BMLT_HOST environment variable.",
		)
		return
	}

	// Validate authentication method - either username/password OR access_token
	hasUsernamePassword := username != "" && password != ""
	hasAccessToken := accessToken != ""

	if !hasUsernamePassword && !hasAccessToken {
		resp.Diagnostics.AddError(
			"Missing Authentication Configuration",
			"The provider requires authentication. Provide either:\n"+
				"1. Both username and password (via configuration or BMLT_USERNAME/BMLT_PASSWORD environment variables)\n"+
				"2. An access_token (via configuration or BMLT_ACCESS_TOKEN environment variable)",
		)
		return
	}

	if hasUsernamePassword && hasAccessToken {
		resp.Diagnostics.AddError(
			"Conflicting Authentication Configuration",
			"The provider cannot use both username/password and access_token authentication methods simultaneously. "+
				"Please provide either username/password OR access_token, not both.",
		)
		return
	}

	if hasUsernamePassword && (username == "" || password == "") {
		if username == "" {
			resp.Diagnostics.AddAttributeError(
				path.Root("username"),
				"Missing BMLT API Username",
				"Username is required when using username/password authentication. Set the username value in the configuration or use the BMLT_USERNAME environment variable.",
			)
		}
		if password == "" {
			resp.Diagnostics.AddAttributeError(
				path.Root("password"),
				"Missing BMLT API Password",
				"Password is required when using username/password authentication. Set the password value in the configuration or use the BMLT_PASSWORD environment variable.",
			)
		}
		return
	}

	// Parse the host URL to extract scheme, host, and path separately
	var scheme, hostOnly, basePath string
	if len(host) >= 7 && host[:7] == "http://" {
		scheme = "http"
		remaining := host[7:]
		if idx := strings.Index(remaining, "/"); idx != -1 {
			hostOnly = remaining[:idx]
			basePath = remaining[idx:]
		} else {
			hostOnly = remaining
			basePath = ""
		}
	} else if len(host) >= 8 && host[:8] == "https://" {
		scheme = "https"
		remaining := host[8:]
		if idx := strings.Index(remaining, "/"); idx != -1 {
			hostOnly = remaining[:idx]
			basePath = remaining[idx:]
		} else {
			hostOnly = remaining
			basePath = ""
		}
	} else {
		// Default to https if no scheme provided
		scheme = "https"
		if idx := strings.Index(host, "/"); idx != -1 {
			hostOnly = host[:idx]
			basePath = host[idx:]
		} else {
			hostOnly = host
			basePath = ""
		}
	}

	// Normalize basePath
	if basePath != "" {
		if !strings.HasPrefix(basePath, "/") {
			basePath = "/" + basePath
		}
		basePath = strings.TrimSuffix(basePath, "/")
	}

	// Create BMLT client configuration
	cfg := bmlt.NewConfiguration()
	cfg.Scheme = scheme
	cfg.Host = hostOnly

	// Add base path to servers configuration
	if basePath != "" {
		cfg.Servers = bmlt.ServerConfigurations{
			{
				URL:         scheme + "://" + hostOnly + basePath,
				Description: "BMLT server with custom path",
			},
		}
	}

	client := bmlt.NewAPIClient(cfg)

	// Set up authentication based on the provided method
	var authCtx context.Context

	if hasAccessToken {
		// Use provided access token directly
		token := &oauth2.Token{
			AccessToken: accessToken,
			TokenType:   "bearer",
		}
		tokenSource := oauth2.StaticTokenSource(token)
		authCtx = context.WithValue(context.Background(), bmlt.ContextOAuth2, tokenSource)
	} else {
		// Use username/password to obtain token
		oauthConfig := &oauth2.Config{
			Endpoint: oauth2.Endpoint{
				TokenURL: scheme + "://" + hostOnly + basePath + "/api/v1/auth/token",
			},
		}

		// Use background context for OAuth2 authentication to avoid timeouts
		token, err := oauthConfig.PasswordCredentialsToken(context.Background(), username, password)
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to Create BMLT API Client",
				"An unexpected error occurred when creating the BMLT API client using username/password: "+err.Error(),
			)
			return
		}

		// Create authenticated context using background context
		tokenSource := oauthConfig.TokenSource(context.Background(), token)
		authCtx = context.WithValue(context.Background(), bmlt.ContextOAuth2, tokenSource)
	}

	// Create a client data structure to pass to resources and data sources
	clientData := &BMTLClientData{
		Client:  client,
		Context: authCtx,
	}

	resp.DataSourceData = clientData
	resp.ResourceData = clientData
}

// BMTLClientData contains the authenticated client and context for use by resources and data sources
type BMTLClientData struct {
	Client  *bmlt.APIClient
	Context context.Context
}

func (p *BMTProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewFormatResource,
		NewMeetingResource,
		NewServiceBodyResource,
		NewSettingsResource,
		NewUserResource,
	}
}

func (p *BMTProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewFormatsDataSource,
		NewMeetingsDataSource,
		NewServiceBodiesDataSource,
		NewServiceBodyDataSource,
		NewSettingsDataSource,
		NewUserDataSource,
		NewUsersDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &BMTProvider{
			version: version,
		}
	}
}
