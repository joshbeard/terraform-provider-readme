package readme_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func Test_Doc_DataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					resource "readme_category" "test" {
						title = "Test Category"
						type  = "guide"
					}
					resource "readme_doc" "test" {
						title    = "My Test Doc"
						body     = "This is a test body"
						category = readme_category.test.id
						type     = "basic"
					}
					data "readme_doc" "test" {
						slug       = readme_doc.test.slug
						depends_on = [ readme_doc.test ]
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr(
						"data.readme_doc.test",
						"id",
						regexp.MustCompile(`^[a-z0-9]{6,}$`),
					),
					resource.TestCheckResourceAttr(
						"data.readme_doc.test",
						"slug",
						"my-test-doc",
					),
					resource.TestMatchResourceAttr(
						"data.readme_doc.test",
						"category",
						regexp.MustCompile(`^[a-z0-9]{6,}$`),
					),
					resource.TestCheckResourceAttr(
						"data.readme_doc.test",
						"type",
						"basic",
					),
					resource.TestCheckResourceAttr(
						"data.readme_doc.test",
						"body",
						"This is a test body",
					),
				),
			},
		},
	})
}
