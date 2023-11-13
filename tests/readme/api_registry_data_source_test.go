package readme_unit_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/liveoaklabs/readme-api-go-client/readme"
	rdmeprovider "github.com/liveoaklabs/terraform-provider-readme/readme"
	"gopkg.in/h2non/gock.v1"
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// testing. The factory function will be invoked for every Terraform CLI command
// executed to create a provider server to which the CLI can reattach.
var testProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"readme": providerserver.NewProtocol6WithError(NewTestProvider(nil)),
}

func NewTestProvider(client *readme.Client) func() tfprotov6.ProviderServer {
	return func() tfprotov6.ProviderServer {
		return rdmeprovider.New("text", client)
	}
}

var providerConfig = `
provider "readme" {
	api_token = "testing"
}
`

func TestAPIRegistryDataSource(t *testing.T) {
	tfConfig := `data "readme_api_registry" "test" { uuid = "somethingUnique" }`

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + tfConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.readme_api_registry.test",
						"uuid",
						"somethingUnique",
					),
					resource.TestCheckResourceAttr("data.readme_api_registry.test", "id", "readme"),
					resource.TestCheckResourceAttr(
						"data.readme_api_registry.test",
						"definition",
						`{"one": "two"}`,
					),
				),
			},
		},
	})
}

func TestAPIRegistryDataSource_GetError(t *testing.T) {
	gock.New(testURL).
		Get("/").
		Persist().
		Reply(401).
		JSON(map[string]string{})
	defer gock.Off()

	expectError, _ := regexp.Compile(
		`Unable to retrieve API registry metadata\.`,
	)

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `data "readme_api_registry" "test" { uuid = "somethingUnique" }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.readme_api_registry.test",
						"uuid",
						"somethingUnique",
					),
				),

				ExpectError: expectError,
			},
		},
	})
}
