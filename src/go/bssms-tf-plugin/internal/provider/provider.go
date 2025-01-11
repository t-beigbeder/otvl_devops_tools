package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure BssmsProvider satisfies various provider interfaces.
var _ provider.Provider = &BssmsProvider{}
var _ provider.ProviderWithFunctions = &BssmsProvider{}

// BssmsProvider defines the provider implementation.
type BssmsProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// bssmsProviderModel describes the provider data model.
type bssmsProviderModel struct {
	ConfigDir types.String `tfsdk:"config_dir"`
}

func (p *BssmsProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "bssms"
	resp.Version = p.version
}

func (p *BssmsProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"config_dir": schema.StringAttribute{
				Optional: true,
			},
		},
	}
}

func (p *BssmsProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring Bssms ConfigDir")

	var config bssmsProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.ConfigDir.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("config_dir"),
			"Unknown Bssms ConfigDir",
			"set the value statically in the configuration",
		)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	resp.ResourceData = config.ConfigDir.String()
	tflog.Info(ctx, "Configured Bssms ConfigDir", map[string]interface{}{"config_dir": config.ConfigDir})
}

func (p *BssmsProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewInstallableResource,
		NewSecretsResource,
	}
}

func (p *BssmsProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return nil
}

func (p *BssmsProvider) Functions(_ context.Context) []func() function.Function {
	return nil
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &BssmsProvider{
			version: version,
		}
	}
}
