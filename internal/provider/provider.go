package provider

import (
	"context"
	"os"

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
	Host     types.String `tfsdk:"host"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
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
		},
	}
}

func (p *BMTProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data BMTProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Configuration values are now available:
	// data.Host, data.Username, data.Password

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

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
			"The provider cannot create the BMLT API client as there is an unknown configuration value for the BMLT API username. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the BMLT_USERNAME environment variable.",
		)
	}

	if data.Password.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Unknown BMLT API Password",
			"The provider cannot create the BMLT API client as there is an unknown configuration value for the BMLT API password. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the BMLT_PASSWORD environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	host := os.Getenv("BMLT_HOST")
	username := os.Getenv("BMLT_USERNAME")
	password := os.Getenv("BMLT_PASSWORD")

	if !data.Host.IsNull() {
		host = data.Host.ValueString()
	}

	if !data.Username.IsNull() {
		username = data.Username.ValueString()
	}

	if !data.Password.IsNull() {
		password = data.Password.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if host == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Missing BMLT API Host",
			"The provider requires a BMLT server host URL. Set the host value in the configuration or use the BMLT_HOST environment variable.",
		)
	}

	if username == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Missing BMLT API Username",
			"The provider requires a username for authentication. Set the username value in the configuration or use the BMLT_USERNAME environment variable.",
		)
	}

	if password == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Missing BMLT API Password",
			"The provider requires a password for authentication. Set the password value in the configuration or use the BMLT_PASSWORD environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Create a new BMLT client using your generated client
	cfg := bmlt.NewConfiguration()
	cfg.Host = host
	cfg.Scheme = "https"
	if host[:7] == "http://" {
		cfg.Scheme = "http"
		host = host[7:]
	} else if host[:8] == "https://" {
		cfg.Scheme = "https"
		host = host[8:]
	}
	cfg.Host = host

	client := bmlt.NewAPIClient(cfg)

	// Set up OAuth2 authentication
	oauthConfig := &oauth2.Config{
		Endpoint: oauth2.Endpoint{
			TokenURL: cfg.Scheme + "://" + cfg.Host + "/api/v1/auth/token",
		},
	}

	token, err := oauthConfig.PasswordCredentialsToken(ctx, username, password)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create BMLT API Client",
			"An unexpected error occurred when creating the BMLT API client: "+err.Error(),
		)
		return
	}

	// Create authenticated context
	tokenSource := oauthConfig.TokenSource(ctx, token)
	authCtx := context.WithValue(ctx, bmlt.ContextOAuth2, tokenSource)

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
		NewUserResource,
	}
}

func (p *BMTProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewFormatsDataSource,
		NewMeetingsDataSource,
		NewServiceBodiesDataSource,
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
