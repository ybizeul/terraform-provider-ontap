package ontap

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	ontap "github.com/ybizeul/terraform-provider-ontap/ontap_client_go"
)

// Ensure ONTAPProvider satisfies various provider interfaces.
var _ provider.Provider = &ONTAPProvider{}
var _ provider.ProviderWithMetadata = &ONTAPProvider{}

// ONTAPProvider defines the provider implementation.
type ONTAPProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// ONTAPProviderModel describes the provider data model.
type ONTAPProviderModel struct {
	Host            types.String `tfsdk:"hostname"`
	Username        types.String `tfsdk:"username"`
	Password        types.String `tfsdk:"password"`
	IgnoreSSLErrors types.Bool   `tfsdk:"ignore_ssl_errors"`
}

func (p *ONTAPProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "ontap"
	resp.Version = p.version
}

func (p *ONTAPProvider) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"hostname": {
				MarkdownDescription: "ONTAP Management Hostname",
				Type:                types.StringType,
				Required:            true,
			},
			"username": {
				MarkdownDescription: "ONTAP Username",
				Type:                types.StringType,
				Required:            true,
				Sensitive:           true,
			},
			"password": {
				MarkdownDescription: "ONTAP Password",
				Type:                types.StringType,
				Required:            true,
				Sensitive:           true,
			},
			"ignore_ssl_errors": {
				MarkdownDescription: "Ignore SSL Errors",
				Type:                types.BoolType,
				Optional:            true,
			},
		},
	}, nil
}

func (p *ONTAPProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data ONTAPProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Configuration values are now available.
	// if data.Endpoint.IsNull() { /* ... */ }

	// Example client configuration for data sources and resources
	if data.IgnoreSSLErrors.Null {
		data.IgnoreSSLErrors.Value = false
	}

	client, _ := ontap.NewClient(&data.Host.Value, &data.Username.Value, &data.Password.Value, data.IgnoreSSLErrors.Value)

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *ONTAPProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewQtreeResource,
		NewSVMResource,
	}
}

func (p *ONTAPProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewQtreeDataSource,
		NewSVMDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &ONTAPProvider{
			version: version,
		}
	}
}
