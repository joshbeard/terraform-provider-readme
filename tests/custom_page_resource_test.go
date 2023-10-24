package readme_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestCustomPageResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					resource "readme_custom_page" "test" {
						title = "Test Page"
						body  = "This is a test body"
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr(
						"readme_custom_page.test",
						"title",
						regexp.MustCompile(`^Test Page$`),
					),
					resource.TestMatchResourceAttr(
						"readme_custom_page.test",
						"body",
						regexp.MustCompile(`^This is a test body$`),
					),
					resource.TestMatchResourceAttr(
						"readme_custom_page.test",
						"id",
						regexp.MustCompile(`^[a-z0-9]{6,}$`),
					),
					resource.TestCheckResourceAttr(
						"readme_custom_page.test",
						"slug",
						"test-page",
					),
					resource.TestMatchResourceAttr(
						"readme_custom_page.test",
						"created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{2,3}Z$`),
					),
					resource.TestMatchResourceAttr(
						"readme_custom_page.test",
						"revision",
						regexp.MustCompile(`^\d+$`),
					),
					resource.TestCheckResourceAttr(
						"readme_custom_page.test",
						"hidden",
						"true",
					),
				),
			},
			{
				Config: `
					resource "readme_custom_page" "test" {
						title  = "Test Page"
						body   = "This is a test body"
						hidden = false
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						"readme_custom_page.test",
						"hidden",
						"false",
					),
				),
			},
			// Test Importing.
			{
				ResourceName:      "readme_custom_page.test",
				ImportState:       true,
				ImportStateId:     "test-page",
				ImportStateVerify: true,
			},
		},
	})
}

func TestCustomPageResource_Errors(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					resource "readme_custom_page" "test" {
						body  = "test body"
					}
				`,
				ExpectError: regexp.MustCompile(
					"'title' must be set using the attribute or in the " +
						"body front matter"),
			},
		},
	})
}

func TestCustomPageResource_FrontMatter(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					resource "readme_custom_page" "test" {
						body   = "---\ntitle: frontmatter title\n---\nThis is a test body"
						hidden = true
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						"readme_custom_page.test",
						"title",
						"frontmatter title",
					),
					resource.TestCheckResourceAttr(
						"readme_custom_page.test",
						"hidden",
						"true",
					),
				),
			},
			{
				Config: `
					resource "readme_custom_page" "test" {
						body  = "---\ntitle: frontmatter title\nhidden: false\n---\nThis is a test body"
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						"readme_custom_page.test",
						"hidden",
						"false",
					),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}
