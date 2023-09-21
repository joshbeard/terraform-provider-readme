resource "readme_category" "test" {
    title = "Test Category0"
        type  = "guide"
}
resource "readme_doc" "test1" {
    title    = "Test Child Doc"
        hidden   = false
        category = readme_category.test.id
}
resource "readme_doc" "test2" {
    title    = "Test Another Child Doc"
        hidden   = false
        category = readme_category.test.id
}
data "readme_category_docs" "test" {
    slug = "test-category0"
    depends_on = [
        readme_category.test,
        readme_doc.test1,
        readme_doc.test2,
    ]
}

# The "readme_category_docs" data source retrieves a list of docs for a category.
# data "readme_category_docs" "example" {
#     slug = "example"
# }

output "category_docs" {
    value = data.readme_category_docs.test
}
