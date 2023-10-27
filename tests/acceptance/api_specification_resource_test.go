package readme_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func Test_APISpecification_Resource_Create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					resource "readme_api_specification" "test2" {
						definition      = file("testdata/example2.json")
						delete_category = true
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("readme_api_specification.test2", "title", "Test API Spec 2"),
				),
			},
			{
				Config: `
					resource "readme_api_specification" "test2" {
						definition      = file("testdata/example2-alt.json")
						delete_category = true
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("readme_api_specification.test2", "title", "Test API Spec 2 Updated"),
				),
			},
		},
	})
}

func Test_APISpecification_Resource_Create_Invalid(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					resource "readme_api_specification" "test" {
						definition      = file("testdata/invalid.json")
						delete_category = true
				}`,
				ExpectError: regexp.MustCompile("The spec you uploaded has validation errors"),
			},
		},
	})
}
