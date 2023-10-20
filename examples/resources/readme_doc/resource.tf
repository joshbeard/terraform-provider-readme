# Manage docs on ReadMe.
# Create a category
# resource "readme_category" "example" {
#     title = "My example category"
#     type  = "guide"
# }

# Create a doc in the category
resource "readme_doc" "example" {
  # title can be specified as an attribute or in the body front matter.
  body = "---\ntitle: Frontmatter Doc\n---\n\nHello! Welcome to my document!"
  #title = "My Example Doc"
  # body = "Hello! Welcome to my document!"

  # category can be specified as an attribute or in the body front matter.
  # Use the `readme_category` resource to manage categories.
  # category = readme_category.example.id
  category_slug = "example"

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
  #body = chomp(file("mydoc.md"))
}
