package provider

import (
	"bssms/bssms"
	"bssms/provisioner"
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"net"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &secretsResource{}
	_ resource.ResourceWithConfigure = &secretsResource{}
)

// NewSecretsResource is a helper function to simplify the provider implementation.
func NewSecretsResource() resource.Resource {
	return &secretsResource{}
}

// secretsResource is the resource implementation.
type secretsResource provisionerConfig

// Metadata returns the resource type name.
func (r *secretsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_secrets"
}

// Schema defines the schema for the resource.
func (r *secretsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name":        schema.StringAttribute{Required: true},
			"pub_key":     schema.StringAttribute{Required: true, Sensitive: true},
			"server_uuid": schema.StringAttribute{Required: true},
			"ip_ext_addresses": schema.ListAttribute{
				ElementType: types.StringType,
				Required:    true,
			},
			"mac_ext_address": schema.StringAttribute{Required: true},
			"ip_int_addresses": schema.ListAttribute{
				ElementType: types.StringType,
				Required:    true,
			},
			"mac_int_address": schema.StringAttribute{Required: true},
			"secrets": schema.MapAttribute{
				ElementType: types.StringType,
				Required:    true,
				Sensitive:   true,
			},
		},
	}
}

func (r *secretsResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	pc, ok := req.ProviderData.(provisionerConfig)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected provisionerConfig, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}
	r.configDir = pc.configDir
	r.pc = pc.pc
}

func getIPV4(addresses []types.String) (string, error) {
	for _, address := range addresses {
		ip := net.ParseIP(address.ValueString())
		if ip != nil && ip.To4() != nil {
			return ip.String(), nil
		}
	}
	return "", errors.New("no IPv4 address found")
}

func makeInstallHost(model secretsResourceModel) (ih provisioner.InstallHost, err error) {
	ih = provisioner.InstallHost{
		Installable: bssms.Installable{
			Name:          model.Name.ValueString(),
			ServerUuid:    model.ServerUuid.ValueString(),
			MacExtAddress: model.MacExtAddress.ValueString(),
			MacIntAddress: model.MacIntAddress.ValueString(),
		},
		PubKey:  model.PubKey.ValueString(),
		Secrets: map[string]string{},
	}
	if ih.IPExtAddress, err = getIPV4(model.IPExtAddresses); err != nil {
		return
	}
	if ih.IPIntAddress, err = getIPV4(model.IPIntAddresses); err != nil {
		return
	}
	for k, v := range model.Secrets.Elements() {
		ih.Secrets[k] = v.String()
	}
	return
}

func (r *secretsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan secretsResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	ih, err := makeInstallHost(plan)
	if err != nil {
		resp.Diagnostics.AddError("Error creating InstallHost", err.Error())
		return
	}
	err = provisioner.RunIhs(
		&bssms.ProvisionerConfig{
			BaseConfig:   bssms.BaseConfig{ctx},
			UnsafeTls:    r.pc.UnsafeTls,
			ProxyAddress: r.pc.ProxyAddress,
		},
		[]provisioner.InstallHost{ih})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Secrets",
			fmt.Sprintf("RunIhs: %s", err),
		)
		return
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *secretsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state secretsResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	//Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *secretsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan secretsResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	ih, err := makeInstallHost(plan)
	if err != nil {
		resp.Diagnostics.AddError("Error updating InstallHost", err.Error())
		return
	}
	err = provisioner.RunIhs(
		&bssms.ProvisionerConfig{
			BaseConfig:   bssms.BaseConfig{ctx},
			UnsafeTls:    r.pc.UnsafeTls,
			ProxyAddress: r.pc.ProxyAddress,
		},
		[]provisioner.InstallHost{ih})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Secrets",
			fmt.Sprintf("RunIhs: %s", err),
		)
		return
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *secretsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state secretsResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// secretsResourceModel maps the resource schema data.
type secretsResourceModel struct {
	Name           types.String   `tfsdk:"name"`
	PubKey         types.String   `tfsdk:"pub_key"`
	ServerUuid     types.String   `tfsdk:"server_uuid"`
	IPExtAddresses []types.String `tfsdk:"ip_ext_addresses"`
	MacExtAddress  types.String   `tfsdk:"mac_ext_address"`
	IPIntAddresses []types.String `tfsdk:"ip_int_addresses"`
	MacIntAddress  types.String   `tfsdk:"mac_int_address"`
	Secrets        types.Map      `tfsdk:"secrets"`
}
