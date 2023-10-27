package readme_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func Test_DocSearch_DataSource(t *testing.T) {
	// FIXME: Skip this test for now.
	t.Skip("DocSearch data source tests are skipped because the search " +
		"indexing doesn't work on free accounts")

	tfConfig := `
		resource "readme_category" "test1" {
			title = "Test Category 1"
			type  = "guide"
		}
		resource "readme_category" "test2" {
			title = "Test Category 2"
			type  = "guide"
		}
		resource "readme_doc" "test1" {
			title    = "Testing North Carolina"
			body     = "Monkey Banana Elephant Apple Orange"
			category = readme_category.test1.id
			hidden   = false
			type     = "basic"
		}
		resource "readme_doc" "test2" {
			title    = "Testing Colorado"
			body     = "Tiger Pear Giraffe Grape Lemon"
			category = readme_category.test1.id
			hidden   = false
			type     = "basic"
		}
		resource "readme_doc" "test3" {
			title    = "Testing Virginia"
			body     = "Koala Pineapple Zebra Plum Avocado"
			category = readme_category.test2.id
			hidden   = false
			type     = "basic"
		}
	`
	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: tfConfig + `
					data "readme_doc_search" "test" {
						query = "*"
						depends_on = [
							readme_doc.test1,
							readme_doc.test2,
							readme_doc.test3,
						]
					}
				`,
				ProtoV6ProviderFactories: testProtoV6ProviderFactories,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr(
						"data.readme_doc_search.test",
						"results.#",
						regexp.MustCompile(`^\d+$`),
					),
				),
			},
			{
				Config: tfConfig + `
					data "readme_doc_search" "test" {
						query = "banana"
						depends_on = [
							readme_doc.test1,
							readme_doc.test2,
							readme_doc.test3,
						]
					}
				`,
				PreConfig: func() {
					// Sleep for 30 seconds to allow the search index to update.
					t.Log("Sleeping for 30 seconds to allow the search index to update...")
					// time.Sleep(30 * time.Second)
				},
				ProtoV6ProviderFactories: testProtoV6ProviderFactories,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.readme_doc_search.test",
						"results.#",
						"1",
					),
				),
			},
		},
	})
}
