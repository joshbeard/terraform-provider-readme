package readme_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func Test_Versions_DataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					resource "readme_version" "example" {
					  version   = "1.1.0"
					  from      = "1.0.0"
					  is_stable = false
					  is_hidden = false
					  codename  = "test"
					}
					data "readme_versions" "example" {
						depends_on = [ readme_version.example ]
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.readme_versions.example",
						"versions.#",
						"2",
					),
					resource.TestMatchResourceAttr(
						"data.readme_versions.example",
						"versions.0.id",
						regexp.MustCompile(`^[a-z0-9]{24}$`),
					),
					resource.TestCheckResourceAttr(
						"data.readme_versions.example",
						"versions.1.version",
						"1.1.0",
					),
					resource.TestCheckResourceAttr(
						"data.readme_versions.example",
						"versions.1.version_clean",
						"1.1.0",
					),
					resource.TestCheckResourceAttr(
						"data.readme_versions.example",
						"versions.1.codename",
						"test",
					),
					resource.TestMatchResourceAttr(
						"data.readme_versions.example",
						"versions.1.forked_from",
						regexp.MustCompile(`^[a-z0-9]{24}$`),
					),
					resource.TestCheckResourceAttr(
						"data.readme_versions.example",
						"versions.1.is_beta",
						"false",
					),
					resource.TestCheckResourceAttr(
						"data.readme_versions.example",
						"versions.1.is_deprecated",
						"false",
					),
					resource.TestCheckResourceAttr(
						"data.readme_versions.example",
						"versions.1.is_hidden",
						"false",
					),
					resource.TestCheckResourceAttr(
						"data.readme_versions.example",
						"versions.1.is_stable",
						"false",
					),
				),
			},
		},
	})
}
