package readme_test

import (
	"context"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	client "github.com/liveoaklabs/readme-api-go-client/readme"
	"github.com/liveoaklabs/terraform-provider-readme/readme"
	"github.com/stretchr/testify/assert"
)

func TestProvider(t *testing.T) {
	resp := provider.SchemaResponse{}

	readmeClient, err := client.NewClient("dev", "hunter2")
	if err != nil {
		t.Fatalf("Error creating client: %s", err)
	}

	prov := readme.New("dev", readmeClient)()
	prov.Schema(context.Background(), provider.SchemaRequest{}, &resp)

	assert.False(t, resp.Diagnostics.HasError())
}

func setupTestProvider() provider.Provider {
	token := os.Getenv("README_API_TOKEN")
	client, err := client.NewClient(token)
	if err != nil {
		panic(err)
	}
	p := readme.New("dev", client)
	if err != nil {
		panic(err)
	}

	return p()
}

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// testing. The factory function will be invoked for every Terraform CLI command
// executed to create a provider server to which the CLI can reattach.
var testProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"readme": providerserver.NewProtocol6WithError(setupTestProvider()),
}

func Test_Provider_MissingAPIToken(t *testing.T) {
	// Ensure the README_API_TOKEN environment variable is unset.
	// This is necessary because the provider will use the environment variable
	// if the api_token field is not set or empty.
	orig := os.Getenv("README_API_TOKEN")
	os.Unsetenv("README_API_TOKEN")

	defer func() {
		os.Setenv("README_API_TOKEN", orig)
	}()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					provider "readme" {
						api_token = ""
					}
					data "readme_project" "test" {}
				`,
				ExpectError: regexp.MustCompile(`Missing ReadMe API Token`),
			},
		},
	})
}

func Test_Provider_EmptyAPIURL(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					provider "readme" {
						api_token = "hunter2"
						api_url   = ""
					}
					data "readme_project" "test" {}
				`,
				ExpectError: regexp.MustCompile(`Missing ReadMe API URL`),
			},
		},
	})
}
