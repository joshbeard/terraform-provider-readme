package readme_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func Test_Version_Resource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					resource "readme_version" "test" {
					  version   = "1.1.1"
					  from      = "1.0.0"
					  is_stable = false
					  is_hidden = false
					  codename  = "test1"
					}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr(
						"readme_version.test",
						"created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d+Z$`),
					),
					resource.TestCheckResourceAttr(
						"readme_version.test",
						"version",
						"1.1.1",
					),
					resource.TestCheckResourceAttr(
						"readme_version.test",
						"version_clean",
						"1.1.1",
					),
					resource.TestCheckResourceAttr(
						"readme_version.test",
						"codename",
						"test1",
					),
					resource.TestMatchResourceAttr(
						"readme_version.test",
						"forked_from",
						regexp.MustCompile(`^[a-z0-9]{24}$`),
					),
					resource.TestCheckResourceAttr(
						"readme_version.test",
						"is_beta",
						"false",
					),
					resource.TestCheckResourceAttr(
						"readme_version.test",
						"is_deprecated",
						"false",
					),
					resource.TestCheckResourceAttr(
						"readme_version.test",
						"is_hidden",
						"false",
					),
					resource.TestCheckResourceAttr(
						"readme_version.test",
						"is_stable",
						"false",
					),
				),
			},
			// Test updating a version.
			{
				Config: `
					resource "readme_version" "test" {
					  version   = "1.1.1"
					  from      = "1.0.0"
					  is_stable = false
					  is_hidden = false
					  codename  = "test2"
					}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						"readme_version.test",
						"codename",
						"test2",
					),
				),
			},
			// === Test importing a version.
			{
				ResourceName:  "readme_version.test",
				ImportState:   true,
				ImportStateId: "1.1.1",
			},
		},
	})
}

func Test_Version_Resource_Validation(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					resource "readme_version" "test" {
					  version   = "1.1.1"
					  from      = "1.0.0"
					  is_stable = true
					  is_hidden = true
					  codename  = "test"
					}`,
				ExpectError: regexp.MustCompile("A stable version cannot be hidden"),
			},
			{
				Config: `
					resource "readme_version" "test" {
					  version       = "1.1.1"
					  from          = "1.0.0"
					  is_stable     = true
					  is_hidden     = false
					  is_deprecated = true
					  codename      = "test"
					}`,
				ExpectError: regexp.MustCompile("A stable version cannot be deprecated"),
			},
		},
	})
}
