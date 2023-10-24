package readme_test

// -----------------------------------------------------------------------------
// Tests for the readme_doc resource's frontmatter handling.
// -----------------------------------------------------------------------------

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// frontmatterTest is a struct that contains the configuration for a test case
// for the readme_doc resource's frontmatter handling.
type frontmatterTest struct {
	AttrName   string // Name of the attribute to test
	AttrConfig string // Config for the attribute preceding the frontmatter.
	AttrValue  string // Value to expect when Frontmatter is ignored.
	FMConfig   string // Config for the when the frontmatter is set without the attribute.
	FMValue    string // Value to expect when Frontmatter is used.
}

// generateFrontmatterTest generates a test configuration and checks for the
// given frontmatterTest. A category is created implicitly.
// This tests that the resource attributes take precedence over the
// frontmatter. When an attribute is not set, the frontmatter is used.
func generateFrontmatterTest(fmTest frontmatterTest) func(t *testing.T) {
	baseConfig := `
		resource "readme_category" "test" {
			title = "Test Category"
			type  = "guide"
		}
	`

	return func(t *testing.T) {
		t.Run(fmt.Sprintf("Frontmatter for %s attribute", fmTest.AttrName), func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: testProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					// === Frontmatter attribute is used when set and attribute is not set.
					{
						Config: baseConfig + fmTest.AttrConfig,
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestMatchResourceAttr(
								"readme_doc.test",
								fmTest.AttrName,
								regexp.MustCompile(fmt.Sprintf(`^%s$`, fmTest.AttrValue)),
							),
						),
					},
					// === Frontmatter attribute is ignored when attribute is set.
					{
						Config: baseConfig + fmTest.FMConfig,
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestMatchResourceAttr(
								"readme_doc.test",
								fmTest.AttrName,
								regexp.MustCompile(fmt.Sprintf(`^%s$`, fmTest.FMValue)),
							),
						),
					},
				},
			})
		})
	}
}

// TestDocResourceFrontMatter tests that the resource attributes take precedence
// over the frontmatter. When an attribute is not set, the frontmatter is used.
func TestDocResourceFrontMatter(t *testing.T) {
	testCases := []frontmatterTest{
		{
			AttrName: "title",
			AttrConfig: `
				resource "readme_doc" "test" {
					title    = "My Test Doc"
					body     = "---\ntitle: Frontmatter Title\n---\nThis is a test body"
					category = readme_category.test.id
					type     = "basic"
				}`,
			AttrValue: "My Test Doc",
			FMConfig: `
				resource "readme_doc" "test" {
					body     = "---\ntitle: Frontmatter Title\n---\nThis is a test body"
					category = readme_category.test.id
					type     = "basic"
				}`,
			FMValue: "Frontmatter Title",
		},
		{
			AttrName: "category_slug",
			AttrConfig: `
				resource "readme_doc" "test" {
					title         = "My Test Doc"
					body          = "---\ncategorySlug: ignored\n---\nThis is a test body"
					category_slug = readme_category.test.slug
					type          = "basic"
				}`,
			AttrValue: "test-category",
			FMConfig: `
				resource "readme_doc" "test" {
					title    = "My Test Doc"
					body     = "---\ncategorySlug: test-category\n---\nThis is a test body"
					type     = "basic"
				}`,
			FMValue: "test-category",
		},
	}

	for _, test := range testCases {
		t.Run(test.AttrName, generateFrontmatterTest(test))
	}
}

// TestDocResourceFrontMatterInvalid tests that the resource fails when the
// frontmatter is invalid YAML.
func TestDocResourceFrontMatterInvalid(t *testing.T) {
	// === negative test: frontmatter is invalid/cannot be unmarshalled.
	t.Run("Frontmatter can't be unmarshalled", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: `
						resource "readme_category" "test" {
							title = "test category"
							type  = "guide"
						}
						resource "readme_doc" "test" {
							title    = "my test doc"
							body     = "---\ntitle:i need space\n---\nthis is a test body"
							category = readme_category.test.id
							type     = "basic"
						}`,
					ExpectError: regexp.MustCompile(`yaml: unmarshal errors:`),
				},
			},
		})
	})
}
