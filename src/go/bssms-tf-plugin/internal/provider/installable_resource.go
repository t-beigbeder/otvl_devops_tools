package provider

import (
	"bssms/provisioner"
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &installableResource{}
	_ resource.ResourceWithConfigure = &installableResource{}
)

// NewInstallableResource is a helper function to simplify the provider implementation.
func NewInstallableResource() resource.Resource {
	return &installableResource{}
}

// installableResource is the resource implementation.
type installableResource struct {
	configDir string
}

// Metadata returns the resource type name.
func (r *installableResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_installable"
}

// Schema defines the schema for the resource.
func (r *installableResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name":    schema.StringAttribute{Required: true},
			"pri_key": schema.StringAttribute{Computed: true, Sensitive: true},
			"pub_key": schema.StringAttribute{Computed: true, Sensitive: true},
		},
	}
}

func (r *installableResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
}

func (r *installableResource) getInstallHost(name string, diags diag.Diagnostics) *provisioner.InstallHost {
	ihs, err := provisioner.LoadInstallHosts(r.configDir)
	if err != nil {
		diags.AddError(
			"Error retrieving Installables",
			fmt.Sprintf("LoadInstallHosts: %s", err),
		)
		return nil
	}
	for _, ih := range ihs {
		if ih.Name == name {
			return &ih
		}
	}
	diags.AddError(
		"Error retrieving Installable",
		fmt.Sprintf("Fetch %s in %s", name, r.configDir),
	)
	return nil
}

func (r *installableResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan installableResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := provisioner.MergePhase0(r.configDir, plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Installable",
			fmt.Sprintf("RunPhase0: %s", err),
		)
		return
	}
	ih := r.getInstallHost(plan.Name.ValueString(), resp.Diagnostics)
	if ih == nil {
		return
	}
	plan.PriKey = types.StringValue(ih.PrivateKey)
	plan.PubKey = types.StringValue(ih.PubKey)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *installableResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state installableResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	ih := r.getInstallHost(state.Name.ValueString(), resp.Diagnostics)
	if ih == nil {
		return
	}
	state.PriKey = types.StringValue(ih.PrivateKey)
	state.PubKey = types.StringValue(ih.PubKey)
	tflog.Info(ctx, "Reading Installable", map[string]interface{}{"ih": ih})
	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *installableResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan installableResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := provisioner.MergePhase0(r.configDir, plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Installable",
			fmt.Sprintf("RunPhase0: %s", err),
		)
		return
	}
	ih := r.getInstallHost(plan.Name.ValueString(), resp.Diagnostics)
	if ih == nil {
		return
	}
	plan.PriKey = types.StringValue(ih.PrivateKey)
	plan.PubKey = types.StringValue(ih.PubKey)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func (r *installableResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state installableResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	ih := r.getInstallHost(state.Name.ValueString(), resp.Diagnostics)
	if ih == nil {
		return
	}
	// nothing to remove, key-pairs are unique to one cloud-init run
}

// installableResourceModel maps the resource schema data.
type installableResourceModel struct {
	Name   types.String `tfsdk:"name"`
	PriKey types.String `tfsdk:"pri_key"`
	PubKey types.String `tfsdk:"pub_key"`
}
