package provider

import (
	"bssms/bssms"
	"bssms/provisioner"
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
			"name":    schema.StringAttribute{Required: true},
			"pub_key": schema.StringAttribute{Computed: true, Sensitive: true},
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

func (r *secretsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan secretsResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	err := provisioner.RunIhs(
		&bssms.ProvisionerConfig{
			BaseConfig:   bssms.BaseConfig{ctx},
			UnsafeTls:    r.pc.UnsafeTls,
			ProxyAddress: r.pc.ProxyAddress,
		},
		[]provisioner.InstallHost{
			{
				Installable: bssms.Installable{
					Name:         plan.Name.ValueString(),
					ServerUuid:   "",
					MacAddress:   "",
					IPExtAddress: "",
					IPIntAddress: "",
					IPAddress:    "",
					EncSecrets:   "",
				},
				PrivateKey: "",
				PubKey:     "",
				Secrets:    nil,
			},
		})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Secrets",
			fmt.Sprintf("RunIhs: %s", err),
		)
		return
	}
}

func (r *secretsResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	//TODO implement me
	panic("implement me")
}

func (r *secretsResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	//TODO implement me
	panic("implement me")
}

func (r *secretsResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	//TODO implement me
	panic("implement me")
}

// secretsResourceModel maps the resource schema data.
type secretsResourceModel struct {
	Name    types.String `tfsdk:"name"`
	PubKey  types.String `tfsdk:"pub_key"`
	Secrets types.Map    `tfsdk:"secrets"`
}
