package readme_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func Test_Category_Resource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					resource "readme_category" "test" {
						title = "Test Category"
						type  = "guide"
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr(
						"readme_category.test",
						"created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d+Z$`),
					),
					resource.TestMatchResourceAttr(
						"readme_category.test",
						"id",
						regexp.MustCompile(`^[a-z0-9]{6,}$`),
					),
					resource.TestMatchResourceAttr(
						"readme_category.test",
						"order",
						regexp.MustCompile(`^\d+$`),
					),
					resource.TestMatchResourceAttr(
						"readme_category.test",
						"project",
						regexp.MustCompile(`^[a-z0-9]{6,}$`),
					),
					resource.TestCheckResourceAttr(
						"readme_category.test",
						"reference",
						"false",
					),
					resource.TestCheckResourceAttr(
						"readme_category.test",
						"slug",
						"test-category",
					),
					resource.TestCheckResourceAttr(
						"readme_category.test",
						"title",
						"Test Category",
					),
					resource.TestCheckResourceAttr(
						"readme_category.test",
						"type",
						"guide",
					),
					resource.TestMatchResourceAttr(
						"readme_category.test",
						"version_id",
						regexp.MustCompile(`^[a-z0-9]{6,}$`),
					),
				),
			},
			// Test updating the category.
			{
				Config: `
					resource "readme_category" "test" {
						title = "My Updated Title"
						type  = "guide"
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						"readme_category.test",
						"title",
						"My Updated Title",
					),
				),
			},
			// Test importing.
			{
				ResourceName:      "readme_category.test",
				ImportState:       true,
				ImportStateId:     "test-category",
				ImportStateVerify: true,
			},
		},
	})
}

func Test_Category_Resource_Validation_Error(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					resource "readme_category" "test" {
						title = "Test Category"
						type  = "invalid"
					}
				`,
				ExpectError: regexp.MustCompile(
					"Category type must be 'guide' or 'reference'",
				),
			},
		},
	})
}
