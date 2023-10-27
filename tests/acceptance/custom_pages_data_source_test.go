package readme_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func Test_CustomPages_DataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					resource "readme_custom_page" "test1" {
						title  = "Test Custom Page 1"
						body   = "This is Custom Page 1."
					}
					resource "readme_custom_page" "test2" {
						title  = "Test Custom Page 2"
						body   = "This is Custom Page 2."
					}
					resource "readme_custom_page" "test3" {
						title  = "Test Custom Page 3"
						body   = "This is Custom Page 3."
					}
					data "readme_custom_pages" "test" {
						depends_on = [
							readme_custom_page.test1,
							readme_custom_page.test2,
							readme_custom_page.test3,
						]
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					customPageDataSourceChecks(1),
					customPageDataSourceChecks(2),
					customPageDataSourceChecks(3),
				),
			},
		},
	})
}

func Test_CustomPages_DataSource_NoResults(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					data "readme_custom_pages" "test" {}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckNoResourceAttr(
						"data.readme_custom_pages.test",
						"results.#",
					),
				),
			},
		},
	})
}

func customPageDataSourceChecks(i int) resource.TestCheckFunc {
	resourceTitle := fmt.Sprintf("data.readme_custom_pages.test")
	prefix := fmt.Sprintf("results.%d.", i-1)

	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(
			resourceTitle,
			"results.#",
			"3",
		),
		resource.TestMatchResourceAttr(
			resourceTitle,
			prefix+"created_at",
			regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{2,3}Z$`),
		),
		resource.TestMatchResourceAttr(
			resourceTitle,
			prefix+"updated_at",
			regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{2,3}Z$`),
		),
		resource.TestCheckResourceAttr(
			resourceTitle,
			prefix+"fullscreen",
			"false",
		),
		resource.TestCheckResourceAttr(
			resourceTitle,
			prefix+"html",
			"",
		),
		resource.TestCheckResourceAttr(
			resourceTitle,
			prefix+"htmlmode",
			"false",
		),
		resource.TestMatchResourceAttr(
			resourceTitle,
			prefix+"slug",
			regexp.MustCompile(`^test-custom-page-\d$`),
		),
		resource.TestMatchResourceAttr(
			resourceTitle,
			prefix+"title",
			regexp.MustCompile(`^Test Custom Page \d$`),
		),
		resource.TestMatchResourceAttr(
			resourceTitle,
			prefix+"body",
			regexp.MustCompile(`^This is Custom Page \d\.$`),
		),
		resource.TestCheckResourceAttr(
			resourceTitle,
			prefix+"hidden",
			"true",
		),
		resource.TestCheckResourceAttr(
			resourceTitle,
			prefix+"revision",
			"2",
		),
	)
}
