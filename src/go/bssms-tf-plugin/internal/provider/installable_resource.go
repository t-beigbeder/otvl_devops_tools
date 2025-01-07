package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/resource"
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

// orderResource is the resource implementation.
type installableResource struct{}

// Metadata returns the resource type name.
func (r *installableResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_installable"
}

// Schema defines the schema for the resource.
func (r *installableResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {

}

func (r *installableResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	//TODO implement me
	panic("implement me")
}

func (r *installableResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	//TODO implement me
	panic("implement me")
}

func (r *installableResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	//TODO implement me
	panic("implement me")
}

func (r *installableResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	//TODO implement me
	panic("implement me")
}

func (r *installableResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	//TODO implement me
	panic("implement me")
}
