package readme_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func Test_CustomPage_DataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					resource "readme_custom_page" "test" {
						title  = "Test Custom Page"
						hidden = false
						body   = "This is my custom page."
					}
					data "readme_custom_page" "test" {
						slug       = "test-custom-page"
						depends_on = [ readme_custom_page.test ]
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr(
						"data.readme_custom_page.test",
						"id",
						regexp.MustCompile(`^[a-z0-9]{6,}$`),
					),
					resource.TestMatchResourceAttr(
						"data.readme_custom_page.test",
						"created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{2,3}Z$`),
					),
					resource.TestMatchResourceAttr(
						"data.readme_custom_page.test",
						"updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{2,3}Z$`),
					),
					resource.TestCheckResourceAttr(
						"data.readme_custom_page.test",
						"fullscreen",
						"false",
					),
					resource.TestCheckResourceAttr(
						"data.readme_custom_page.test",
						"html",
						"",
					),
					resource.TestCheckResourceAttr(
						"data.readme_custom_page.test",
						"htmlmode",
						"false",
					),
					resource.TestCheckResourceAttr(
						"data.readme_custom_page.test",
						"slug",
						"test-custom-page",
					),
					resource.TestCheckResourceAttr(
						"data.readme_custom_page.test",
						"title",
						"Test Custom Page",
					),
					resource.TestCheckResourceAttr(
						"data.readme_custom_page.test",
						"body",
						"This is my custom page.",
					),
					resource.TestCheckResourceAttr(
						"data.readme_custom_page.test",
						"hidden",
						"false",
					),
					resource.TestCheckResourceAttr(
						"data.readme_custom_page.test",
						"revision",
						"2",
					),
				),
			},
		},
	})
}

func Test_CustomPage_DataSource_NotFound(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					data "readme_custom_page" "test" {
						slug = "does-not-exist"
					}
				`,
				ExpectError: regexp.MustCompile(`Unable to retrieve custom pages`),
			},
		},
	})
}
