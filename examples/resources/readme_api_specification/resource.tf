# # Create an API specification resource.
resource "readme_api_specification" "example" {
  # 'definition' accepts a string of an OpenAPI specification definition JSON.
  #definition = file("${path.module}/../../../tests/testdata/example1.json")
  definition = file("${path.module}/example0.json")

  # When an API specification is created, a category is also created but is
  # not deleted when the API specification is deleted. Set this parameter to
  # true to delete the category when the API specification is deleted.
  delete_category = true
}
#
# # Output the ID of the created resource.
# output "created_spec_id" {
#   value = readme_api_specification.example.id
# }
#
# # Output the specification JSON of the created resource.
# output "created_spec_json" {
#   value = readme_api_specification.example.definition
# }
# The "readme_category" resource manages the lifecycle of a category in ReadMe.
resource "readme_category" "example" {
    title = "My example category"
    type  = "guide"
}

# Manage docs on ReadMe.
resource "readme_doc" "example" {
    # title can be specified as an attribute or in the body front matter.
    title = "Example0"

    # category can be specified as an attribute or in the body front matter.
    # Use the `readme_category` resource to manage categories.
    category = readme_category.example.id

    # category_slug can be specified as an attribute or in the body front matter.
    # category_slug = "foo-bar"

    # hidden can be specified as an attribute or in the body front matter.
    hidden = false

    # order can be specified as an attribute or in the body front matter.
    order = 99

    # type can be specified as an attribute or in the body front matter.
    type = "basic"

    # body can be read from a file using Terraform's `file()` function.
    # For best results, wrap the string with the `chomp()` function to remove
    # trailing newlines. ReadMe's API trims these implicitly.
    body = chomp(file("mydoc.md"))
}
