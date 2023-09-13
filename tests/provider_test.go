package readme_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/liveoaklabs/terraform-provider-readme/readme"
	"github.com/stretchr/testify/assert"
)

const (
	// testToken is a dummy token the provider is configured with and used
	// throughout tests.
	testToken = "hunter2"
	// providerConfig is a shared configuration that sets a mock url and token.
	// The URL points to our gock mock server.
	providerConfig = (`
		provider "readme" {
			api_token = "` + testToken + `"
		}
	`)
)

// var testAccProviders map[string]*schema.Provider
// var testAccProviders map[string]func() provider.Provider
// var testAccProvider func() provider.Provider

func TestProvider(t *testing.T) {
	t.Parallel()
	resp := provider.SchemaResponse{}

	prov := readme.New("dev")()
	prov.Schema(context.Background(), provider.SchemaRequest{}, &resp)

	assert.False(t, resp.Diagnostics.HasError())
}

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// testing. The factory function will be invoked for every Terraform CLI command
// executed to create a provider server to which the CLI can reattach.
var testProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"readme": providerserver.NewProtocol6WithError(readme.New("dev")()),
}

// err := providerserver.Serve(context.Background(), readme.New(version), providerserver.ServeOpts{
// 	Address: "registry.terraform.io/liveoaklabs/readme",
// 	// This provider requires Terraform 1.0+
// 	ProtocolVersion: 6,
// })
// if err != nil {
// 	log.Fatal("Error setting up ReadMe provider:", err)
// }

// func init() {
// 	testAccProvider = New("dev")
// 	testAccProviders = map[string]func() provider.Provider{
// 		"readme": testAccProvider,
// 	}
// }

// func TestProvider(t *testing.T) {
// 	if err := New("dev").InternalValidate(); err != nil {
// 		t.Fatalf("err: %s", err)
// 	}
// }
//
// func TestProvider_impl(t *testing.T) {
// 	var _ *schema.Provider = Provider()
// }
