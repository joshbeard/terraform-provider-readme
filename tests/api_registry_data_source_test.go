package readme_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func Test_APIRegistry_DataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					resource "readme_api_specification" "test" {
						definition      = file("testdata/example1.json")
						delete_category = true
					}
					data "readme_api_registry" "test" {
						uuid = readme_api_specification.test.uuid
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.readme_api_registry.test", "id", "readme"),
				),
			},
		},
	})
}

func Test_APIRegistry_DataSource_Error(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      `data "readme_api_registry" "test" { uuid = "doesntexist" }`,
				ExpectError: regexp.MustCompile(`Unable to retrieve API registry metadata\.`),
			},
		},
	})
}
