package readme_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var ReadmeProjectConfig = struct {
	ProjectName string
	Subdomain   string
	Plan        string
}{
	ProjectName: "jbeard-202309",
	Subdomain:   "jbeard-202309",
	Plan:        "business2018",
}

func Test_Project_DataSource(t *testing.T) {
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
						ReadmeProjectConfig.ProjectName,
					),
					resource.TestCheckResourceAttr(
						"data.readme_project.test",
						"subdomain",
						ReadmeProjectConfig.Subdomain,
					),
					resource.TestCheckResourceAttr(
						"data.readme_project.test",
						"base_url",
						fmt.Sprintf("https://%s.readme.io", ReadmeProjectConfig.Subdomain),
					),
					resource.TestCheckResourceAttr(
						"data.readme_project.test",
						"plan",
						ReadmeProjectConfig.Plan,
					),
					// Verify placeholder id attribute.
					// See https://developer.hashicorp.com/terraform/plugin/framework/acctests#implement-id-attribute
					resource.TestCheckResourceAttr("data.readme_project.test", "id", "readme"),
				),
			},
		},
	})
}
