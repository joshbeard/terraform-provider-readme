package readme_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func Test_Category_DataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					resource "readme_category" "test" {
						title = "Test Category"
						type  = "guide"
					}
					data "readme_category" "test" {
						slug       = "test-category"
						depends_on = [ readme_category.test ]
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.readme_category.test",
						"slug",
						"test-category",
					),
					resource.TestCheckResourceAttr(
						"data.readme_category.test",
						"title",
						"Test Category",
					),
					resource.TestCheckResourceAttr(
						"data.readme_category.test",
						"type",
						"guide",
					),
					resource.TestMatchResourceAttr(
						"data.readme_category.test",
						"id",
						regexp.MustCompile(`^[a-z0-9]{6,}$`),
					),
					resource.TestMatchResourceAttr(
						"data.readme_category.test",
						"version_id",
						regexp.MustCompile(`^[a-z0-9]{6,}$`),
					),
					resource.TestMatchResourceAttr(
						"data.readme_category.test",
						"project",
						regexp.MustCompile(`^[a-z0-9]{6,}$`),
					),
					resource.TestMatchResourceAttr(
						"data.readme_category.test",
						"created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d+Z$`),
					),
					resource.TestCheckResourceAttr(
						"data.readme_category.test",
						"reference",
						"false",
					),
				),
			},
		},
	})
}
