package readme_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestProjectDataSource(t *testing.T) {
	tfConfig := `data "readme_project" "test" {}`

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config:                   tfConfig,
				ProtoV6ProviderFactories: testProtoV6ProviderFactories,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify the project response is returned with the data source.
					resource.TestCheckResourceAttr(
						"data.readme_project.test",
						"name",
						"jbeard-202309",
					),
					resource.TestCheckResourceAttr(
						"data.readme_project.test",
						"subdomain",
						"jbeard-202309",
					),
					resource.TestCheckResourceAttr(
						"data.readme_project.test",
						"base_url",
						"https://jbeard-202309.readme.io",
					),
					resource.TestCheckResourceAttr(
						"data.readme_project.test",
						"plan",
						"business2018",
					),
					// Verify placeholder id attribute.
					// See https://developer.hashicorp.com/terraform/plugin/framework/acctests#implement-id-attribute
					resource.TestCheckResourceAttr("data.readme_project.test", "id", "readme"),
				),
			},
		},
	})
}
