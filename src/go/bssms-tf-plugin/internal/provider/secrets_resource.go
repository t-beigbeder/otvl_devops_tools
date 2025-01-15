package provider

import (
	"bssms/bssms"
	"bssms/provisioner"
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/lrstanley/go-bogon"
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
			"ip_v4_addresses": schema.ListAttribute{
				ElementType: types.StringType,
				Required:    true,
			},
			"mac_addresses": schema.ListAttribute{
				ElementType: types.StringType,
				Required:    true,
			},
			"secrets": schema.MapAttribute{
				ElementType: types.StringType,
				Required:    true,
				Sensitive:   true,
			},
			"tofu_running": schema.BoolAttribute{Computed: true},
			"installed":    schema.BoolAttribute{Computed: true},
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

func ipIsExt(sip string) bool {
	tbogon, _ := bogon.Is(sip)
	return !tbogon
}

func makeInstallHost(model secretsResourceModel) (ih provisioner.InstallHost, err error) {
	ih = provisioner.InstallHost{
		Installable: bssms.Installable{
			Name:       model.Name.ValueString(),
			ServerUuid: model.ServerUuid.ValueString(),
		},
		PubKey:  model.PubKey.ValueString(),
		Secrets: map[string]string{},
	}
	if len(model.IPV4Addresses) != len(model.MacAddresses) {
		return ih, fmt.Errorf("IPV4Addresses and MacAddresses are not the same length")
	}
	for i, address := range model.IPV4Addresses {
		if ipIsExt(address.ValueString()) {
			ih.IPExtAddress = address.ValueString()
			ih.MacExtAddress = model.MacAddresses[i].ValueString()
		} else {
			ih.IPIntAddress = address.ValueString()
			ih.MacIntAddress = model.MacAddresses[i].ValueString()
		}
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
	tflog.Info(ctx, "Create: Installing Host", map[string]interface{}{"ih": ih})
	err = provisioner.TofuRun(
		&bssms.ProvisionerConfig{
			BaseConfig:   bssms.BaseConfig{ctx},
			UnsafeTls:    r.pc.UnsafeTls,
			ProxyAddress: r.pc.ProxyAddress,
		},
		r.configDir, ih)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Secrets",
			fmt.Sprintf("TofuRun: %s", err),
		)
		return
	}
	sih, err := provisioner.ReadInstallHost(r.configDir, ih.Name)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Secrets",
			fmt.Sprintf("ReadInstallHost: %s", err),
		)
		return
	}
	if !sih.Installed {
		resp.Diagnostics.AddError(
			"Error creating Secrets",
			fmt.Sprintf("ReadInstallHost: %s not installed", sih.Name),
		)
		return
	}

	plan.Installed = types.BoolValue(true)
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
	ih, err := provisioner.ReadInstallHost(r.configDir, state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Secrets",
			fmt.Sprintf("ReadInstallHost: %s", err),
		)
		return
	}
	tflog.Info(ctx, "Read: Installing Host", map[string]interface{}{"ih": ih})
	state.TofuRunning = types.BoolValue(ih.TofuRunning)
	state.Installed = types.BoolValue(ih.Installed)

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
		resp.Diagnostics.AddError("Error creating InstallHost", err.Error())
		return
	}
	tflog.Info(ctx, "Update: Installing Host", map[string]interface{}{"ih": ih})
	err = provisioner.TofuRun(
		&bssms.ProvisionerConfig{
			BaseConfig:   bssms.BaseConfig{ctx},
			UnsafeTls:    r.pc.UnsafeTls,
			ProxyAddress: r.pc.ProxyAddress,
		},
		r.configDir, ih)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Secrets",
			fmt.Sprintf("TofuRun: %s", err),
		)
		return
	}
	sih, err := provisioner.ReadInstallHost(r.configDir, ih.Name)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Secrets",
			fmt.Sprintf("ReadInstallHost: %s", err),
		)
		return
	}
	if !sih.Installed {
		resp.Diagnostics.AddError(
			"Error updating Secrets",
			fmt.Sprintf("ReadInstallHost: %s not installed", sih.Name),
		)
		return
	}

	plan.Installed = types.BoolValue(true)
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
	Name          types.String   `tfsdk:"name"`
	PubKey        types.String   `tfsdk:"pub_key"`
	ServerUuid    types.String   `tfsdk:"server_uuid"`
	IPV4Addresses []types.String `tfsdk:"ip_v4_addresses"`
	MacAddresses  []types.String `tfsdk:"mac_addresses"`
	Secrets       types.Map      `tfsdk:"secrets"`
	TofuRunning   types.Bool     `tfsdk:"tofu_running"`
	Installed     types.Bool     `tfsdk:"installed"`
}
