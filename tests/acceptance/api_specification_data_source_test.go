package readme_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func Test_APISpecification_DataSource(t *testing.T) {
	createResource := `
		resource "readme_api_specification" "test" {
			definition      = file("testdata/example1.json")
			delete_category = true
		}
	`
	t.Run("matching id", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: createResource + `
						data "readme_api_specification" "test" {
							id = readme_api_specification.test.id
						}
					`,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.readme_api_specification.test", "title", "Test API Spec 1"),
					),
				},
			},
		})
	})

	t.Run("matching title", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: createResource + `
						data "readme_api_specification" "test" {
							title = readme_api_specification.test.title
						}
					`,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.readme_api_specification.test", "title", "Test API Spec 1"),
					),
				},
			},
		})
	})

	t.Run("matching id and title", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: createResource + `
						data "readme_api_specification" "test" {
							title = readme_api_specification.test.title
							id    = readme_api_specification.test.id
						}
					`,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.readme_api_specification.test", "title", "Test API Spec 1"),
					),
				},
			},
		})
	})

	t.Run("matching title and category id", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: createResource + `
						data "readme_api_specification" "test" {
							title  = readme_api_specification.test.title
							filter = { category_id = readme_api_specification.test.category.id }
						}
					`,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.readme_api_specification.test", "title", "Test API Spec 1"),
					),
				},
			},
		})
	})

	t.Run("matching title and category title", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: createResource + `
						data "readme_api_specification" "test" {
							title  = readme_api_specification.test.title
							filter = { category_title = readme_api_specification.test.category.title }
						}
					`,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.readme_api_specification.test", "title", "Test API Spec 1"),
					),
				},
			},
		})
	})

	t.Run("matching title and category slug", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: createResource + `
						data "readme_api_specification" "test" {
							title  = readme_api_specification.test.title
							filter = { category_slug = readme_api_specification.test.category.slug }
						}
					`,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.readme_api_specification.test", "title", "Test API Spec 1"),
					),
				},
			},
		})
	})

	t.Run("matching title, category id, and category title", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: createResource + `
						data "readme_api_specification" "test" {
							title  = readme_api_specification.test.title
							filter = {
								category_id    = readme_api_specification.test.category.id
								category_title = readme_api_specification.test.category.title
							}
						}
					`,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.readme_api_specification.test", "title", "Test API Spec 1"),
					),
				},
			},
		})
	})

	t.Run("matching title, category id, and category slug", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: createResource + `
						data "readme_api_specification" "test" {
							title  = readme_api_specification.test.title
							filter = {
								category_id   = readme_api_specification.test.category.id
								category_slug = readme_api_specification.test.category.slug
							}
						}
					`,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.readme_api_specification.test", "title", "Test API Spec 1"),
					),
				},
			},
		})
	})

	// 	{
	// 		name: "it should return an API spec when a matching id and category_id is found",
	// 		config: `
	// 			data "readme_api_specification" "$TEST$" {
	// 				id     = readme_api_specification.test.id
	// 				filter = { category_id = readme_api_specification.test.category.id }
	// 		}`,
	// 	},
	// 	{
	// 		name: "it should return an API spec when a matching id and category_title is found",
	// 		config: `
	// 			data "readme_api_specification" "$TEST$" {
	// 				id     = readme_api_specification.test.id
	// 				filter = { category_title = readme_api_specification.test.category.title }
	// 		}`,
	// 	},
	// 	{
	// 		name: "it should return an API spec when a matching id and category_slug is found",
	// 		config: `
	// 			data "readme_api_specification" "$TEST$" {
	// 				id     = readme_api_specification.test.id
	// 				filter = { category_slug = readme_api_specification.test.category.slug }
	// 		}`,
	// 	},
	//

	// Negative tests
	t.Run("error when no ID or title is provided", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config:      createResource + `data "readme_api_specification" "test" {} `,
					ExpectError: regexp.MustCompile("An ID or title must be specified to retrieve an API specification"),
				},
			},
		})
	})

	t.Run("error when no matching title", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: createResource + `
						data "readme_api_specification" "test" {
							title = "This doesn't exist"
						}
					`,
					ExpectError: regexp.MustCompile("Unable to find API specification with title: This doesn't exist"),
				},
			},
		})
	})

	t.Run("error when no matching id", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: createResource + `
						data "readme_api_specification" "test" {
							id = "000000000000000"
						}
					`,
					ExpectError: regexp.MustCompile("API specification not found"),
				},
			},
		})
	})

	t.Run("error when no matching ID or category slug", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: createResource + `
						data "readme_api_specification" "test" {
							id     = readme_api_specification.test.id
							filter = { category_slug = "does-not-exist" }
						}
					`,
					ExpectError: regexp.MustCompile("Unable to find API specification with the specified criteria"),
				},
			},
		})
	})

	t.Run("error when no matching title or category slug", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: createResource + `
						data "readme_api_specification" "test" {
							title  = readme_api_specification.test.title
							filter = { category_slug = "does-not-exist" }
						}
					`,
					ExpectError: regexp.MustCompile("Unable to find API specification with title: Test API Spec 1"),
				},
			},
		})
	})

	t.Run("error when no matching title or category id", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: createResource + `
						data "readme_api_specification" "test" {
							title  = readme_api_specification.test.title
							filter = { category_id = "does-not-exist" }
						}
					`,
					ExpectError: regexp.MustCompile("Unable to find API specification with title: Test API Spec 1"),
				},
			},
		})
	})

	t.Run("error when no matching title or category title", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: createResource + `
						data "readme_api_specification" "test" {
							title  = readme_api_specification.test.title
							filter = { category_title = "does-not-exist" }
						}
					`,
					ExpectError: regexp.MustCompile("Unable to find API specification with title: Test API Spec 1"),
				},
			},
		})
	})

	t.Run("error when no matching title or has category", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: createResource + `
						data "readme_api_specification" "test" {
							title  = "does not exist"
							filter = { has_category = false }
						}
					`,
					ExpectError: regexp.MustCompile("Unable to find API specification with title: does not exist"),
				},
			},
		})
	})
}
