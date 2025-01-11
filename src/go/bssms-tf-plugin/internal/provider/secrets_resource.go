package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/resource"
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

// orderResource is the resource implementation.
type secretsResource struct{}

// Metadata returns the resource type name.
func (r *secretsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_secrets"
}

// Schema defines the schema for the resource.
func (r *secretsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {

}

func (r *secretsResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	//TODO implement me
	panic("implement me")
}

func (r *secretsResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	//TODO implement me
	panic("implement me")
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
