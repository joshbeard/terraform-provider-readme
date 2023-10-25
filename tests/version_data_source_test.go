package readme_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func Test_Version_DataSource(t *testing.T) {
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
					data "readme_version" "example" {
						version_clean = "1.1.0"
						depends_on    = [ readme_version.example ]
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr(
						"data.readme_version.example",
						"id",
						regexp.MustCompile(`^[a-z0-9]{24}$`),
					),
					resource.TestCheckResourceAttr(
						"data.readme_version.example",
						"version",
						"1.1.0",
					),
					resource.TestCheckResourceAttr(
						"data.readme_version.example",
						"version_clean",
						"1.1.0",
					),
					resource.TestCheckResourceAttr(
						"data.readme_version.example",
						"codename",
						"test",
					),
					resource.TestMatchResourceAttr(
						"data.readme_version.example",
						"forked_from",
						regexp.MustCompile(`^[a-z0-9]{24}$`),
					),
					resource.TestCheckResourceAttr(
						"data.readme_version.example",
						"is_beta",
						"false",
					),
					resource.TestCheckResourceAttr(
						"data.readme_version.example",
						"is_deprecated",
						"false",
					),
					resource.TestCheckResourceAttr(
						"data.readme_version.example",
						"is_hidden",
						"false",
					),
					resource.TestCheckResourceAttr(
						"data.readme_version.example",
						"is_stable",
						"false",
					),
					resource.TestMatchResourceAttr(
						"data.readme_version.example",
						"created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d+Z$`),
					),
				),
			},
		},
	})
}
