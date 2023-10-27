package readme_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func Test_Categories_DataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					resource "readme_category" "test" {
						title = "Test Category"
						type  = "guide"
					}
					data "readme_categories" "test" {
						depends_on = [ readme_category.test ]
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.readme_categories.test",
						"categories.#",
						"1",
					),
					resource.TestCheckResourceAttr(
						"data.readme_categories.test",
						"categories.0.slug",
						"test-category",
					),
					resource.TestCheckResourceAttr(
						"data.readme_categories.test",
						"categories.0.title",
						"Test Category",
					),
					resource.TestCheckResourceAttr(
						"data.readme_categories.test",
						"categories.0.type",
						"guide",
					),
					resource.TestMatchResourceAttr(
						"data.readme_categories.test",
						"categories.0.id",
						regexp.MustCompile(`^[a-z0-9]{6,}$`),
					),
					resource.TestMatchResourceAttr(
						"data.readme_categories.test",
						"categories.0.version_id",
						regexp.MustCompile(`^[a-z0-9]{6,}$`),
					),
					resource.TestMatchResourceAttr(
						"data.readme_categories.test",
						"categories.0.project",
						regexp.MustCompile(`^[a-z0-9]{6,}$`),
					),
					resource.TestMatchResourceAttr(
						"data.readme_categories.test",
						"categories.0.created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d+Z$`),
					),
					resource.TestCheckResourceAttr(
						"data.readme_categories.test",
						"categories.0.reference",
						"false",
					),
				),
			},
		},
	})
}
