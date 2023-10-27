package readme_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func Test_APISpecifications_DataSource(t *testing.T) {
	createResource := `
		resource "readme_api_specification" "test1" {
			definition      = file("testdata/example1.json")
			delete_category = true
		}
		resource "readme_api_specification" "test2" {
			definition      = file("testdata/example2.json")
			delete_category = true
		}
		resource "readme_api_specification" "test3" {
			definition      = file("testdata/example3.json")
			delete_category = true
		}
	`

	t.Run("get all specs", func(t *testing.T) {
		rsc := "data.readme_api_specifications.no_filters"
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: createResource + `
						data "readme_api_specifications" "no_filters" {
							depends_on = [
								readme_api_specification.test1,
								readme_api_specification.test2,
								readme_api_specification.test3,
							]
					}`,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(rsc, "specs.#", "3"),
					),
				},
			},
		})
	})

	t.Run("sorted by title", func(t *testing.T) {
		rsc := "data.readme_api_specifications.sort_by_title"
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: createResource + `
						data "readme_api_specifications" "sort_by_title" {
							depends_on = [
								readme_api_specification.test1,
								readme_api_specification.test2,
								readme_api_specification.test3,
							]
							sort_by = "title"
					}`,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(rsc, "specs.0.title", "Test API Spec 1"),
						resource.TestCheckResourceAttr(rsc, "specs.1.title", "Test API Spec 2"),
						resource.TestCheckResourceAttr(rsc, "specs.2.title", "Test API Spec 3"),
					),
				},
			},
		})
	})

	t.Run("sorted by last_synced", func(t *testing.T) {
		rsc := "data.readme_api_specifications.sort_by_last_synced"
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: createResource + `
						data "readme_api_specifications" "sort_by_last_synced" {
							depends_on = [
								readme_api_specification.test1,
								readme_api_specification.test2,
								readme_api_specification.test3,
							]
							sort_by = "last_synced"
					}`,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(rsc, "specs.#", "3"),
					),
				},
				{
					Config: `
						resource "readme_api_specification" "test1" {
							definition      = file("testdata/example1.json")
							delete_category = true
						}
						resource "readme_api_specification" "test2" {
							definition      = file("testdata/example2-alt.json")
							delete_category = true
						}
						resource "readme_api_specification" "test3" {
							definition      = file("testdata/example3.json")
							delete_category = true
						}
						data "readme_api_specifications" "sort_by_last_synced" {
							depends_on = [
								readme_api_specification.test1,
								readme_api_specification.test2,
								readme_api_specification.test3,
							]
							sort_by = "last_synced"
					}`,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(rsc, "specs.0.title", "Test API Spec 2 Updated"),
					),
				},
			},
		})
	})

	t.Run("filter by title", func(t *testing.T) {
		rsc := "data.readme_api_specifications.title_filter"
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: createResource + `
						data "readme_api_specifications" "title_filter" {
							depends_on = [
								readme_api_specification.test1,
								readme_api_specification.test2,
								readme_api_specification.test3,
							]
							filter = {
								title = ["Test API Spec 1", "Test API Spec 3"]
							}
					}`,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(rsc, "specs.#", "2"),
					),
				},
			},
		})
	})

	t.Run("filter by version", func(t *testing.T) {
		t.Skip("SKIPPED: Implement this test")

		rsc := "data.readme_api_specifications.version_filter"
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: createResource + `
						resource "readme_version" "v1_1_0" {
							version = "1.1.0"
							from    = "1.0.0"
						}
						resource "readme_api_specification" "test1_1_1_0" {
							definition      = file("testdata/example1_1.1.0.json")
							delete_category = true
							depends_on      = [readme_version.v1_1_0]
						}
						data "readme_api_specifications" "version_filter" {
							depends_on = [
								readme_api_specification.test1,
								readme_api_specification.test2,
								readme_api_specification.test3,
								readme_api_specification.test1_1_1_0,
							]
							filter = {
								version = ["1.1.0"]
							}
					}`,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(rsc, "specs.0.title", "Test API Spec 1 v1.1.0"),
						resource.TestCheckResourceAttr(rsc, "specs.#", "1"),
					),
				},
			},
		})
	})

	t.Run("no filter matches", func(t *testing.T) {
		rsc := "data.readme_api_specifications.no_matches"
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: createResource + `
						data "readme_api_specifications" "no_matches" {
							depends_on = [
								readme_api_specification.test1,
								readme_api_specification.test2,
								readme_api_specification.test3,
							]
							filter = {
								has_category   = true
								category_id    = ["000000011111122222223333"]
								category_slug  = ["this-doesnt-exist"]
								category_title = ["This doesn't exist"]
								title          = ["This doesn't exist either"]
							}
					}`,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(rsc, "specs.#", "0"),
					),
				},
			},
		})
	})
}

// t.Run("it should return API specs that match the provided has_category filter", func(t *testing.T) {
// 	r := "data.readme_api_specifications.version_filter"
// 	resource.Test(t, resource.TestCase{
// 		ProtoV6ProviderFactories: testProtoV6ProviderFactories,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: createResource + `
// 					resource "readme_api_specification" "test1_1_1_0" {
// 						definition      = file("testdata/example1_1.1.0.json")
// 						delete_category = true
// 					}
// 					data "readme_api_specifications" "version_filter" {
// 						depends_on = [
// 							readme_api_specification.test1,
// 							readme_api_specification.test2,
// 							readme_api_specification.test3,
// 							readme_api_specification.test1_1_1_0,
// 						]
// 						filter = {
// 							has_category = true
// 						}
// 				}`,
// 				Check: resource.ComposeAggregateTestCheckFunc(
// 					resource.TestCheckResourceAttr(r, "specs.0.title", "Test API Spec 1 v1.1.0"),
// 					resource.TestCheckResourceAttr(r, "specs.#.title", "1"),
// 				),
// 			},
// 		},
// 	})
// 	// 	config: `data "readme_api_specifications" "test" {
// 	// 		filter = { has_category = true }
// 	// 	}`,
// 	// 	response: []readme.APISpecification{
// 	// 		testdata.APISpecifications[0],
// 	// 		testdata.APISpecifications[1],
// 	// 	},
// 	// 	excluded: []readme.APISpecification{testdata.APISpecifications[2]},
// })
// {
// 	name: "it should return only API specs that have a category when has_category is true and used with another filter",
// 	config: fmt.Sprintf(`data "readme_api_specifications" "test" {
// 		filter = {
// 			has_category  = true
// 			category_slug = ["%s", "%s", "test-api-spec-without-category"]
// 		}
// 	}`,
// 		testdata.APISpecifications[0].Category.Slug,
// 		testdata.APISpecifications[1].Category.Slug,
// 	),
// 	response: []readme.APISpecification{
// 		testdata.APISpecifications[0],
// 		testdata.APISpecifications[1],
// 	},
// 	excluded: []readme.APISpecification{testdata.APISpecifications[2]},
// },
// {
// 	name: "it should return an empty list when has_category is true and used with another filter that doesn't match any API specs",
// 	config: `data "readme_api_specifications" "test" {
// 		filter = {
// 			has_category  = true
// 			category_slug = ["test-api-spec-without-category"]
// 		}
// 	}`,
// 	response: []readme.APISpecification{},
// },
//
// {
