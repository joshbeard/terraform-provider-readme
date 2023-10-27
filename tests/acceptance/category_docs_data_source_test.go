package readme_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func Test_CategoryDocs_DataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					resource "readme_category" "test" {
						title = "Test Category0"
						type  = "guide"
					}
					resource "readme_doc" "parent" {
						title    = "Test Parent Doc"
						hidden   = false
						category = readme_category.test.id
						body     = "parent doc"
					}
					resource "readme_doc" "child" {
						title           = "Test Child Doc"
						hidden          = false
						category        = readme_category.test.id
						body            = "child doc"
						parent_doc_slug = readme_doc.parent.slug
					}
					data "readme_category_docs" "test" {
						slug       = "test-category0"
						depends_on = [
							readme_doc.parent,
							readme_doc.child,
						]
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr(
						"data.readme_category_docs.test",
						"docs.0.id",
						regexp.MustCompile(`^[a-z0-9]{6,}$`),
					),
					resource.TestCheckResourceAttr(
						"data.readme_category_docs.test",
						"docs.0.title",
						"Test Parent Doc",
					),
					resource.TestCheckResourceAttr(
						"data.readme_category_docs.test",
						"docs.0.slug",
						"test-parent-doc",
					),
					resource.TestMatchResourceAttr(
						"data.readme_category_docs.test",
						"docs.0.order",
						regexp.MustCompile(`^\d+$`),
					),
					resource.TestCheckResourceAttr(
						"data.readme_category_docs.test",
						"docs.0.hidden",
						"false",
					),
					resource.TestMatchResourceAttr(
						"data.readme_category_docs.test",
						"docs.0.children.0.id",
						regexp.MustCompile(`^[a-z0-9]{6,}$`),
					),
					resource.TestCheckResourceAttr(
						"data.readme_category_docs.test",
						"docs.0.children.0.slug",
						"test-child-doc",
					),
					resource.TestCheckResourceAttr(
						"data.readme_category_docs.test",
						"docs.0.children.0.title",
						"Test Child Doc",
					),
					resource.TestMatchResourceAttr(
						"data.readme_category_docs.test",
						"docs.0.children.0.order",
						regexp.MustCompile(`^\d+$`),
					),
					resource.TestCheckResourceAttr(
						"data.readme_category_docs.test",
						"docs.0.children.0.hidden",
						"false",
					),
				),
			},
		},
	})
}
