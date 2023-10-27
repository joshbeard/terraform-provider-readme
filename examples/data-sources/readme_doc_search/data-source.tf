<<<<<<< HEAD
# Search for docs on ReadMe.
data "readme_doc_search" "example" {
  query = "*"
}

output "example_doc_search" {
  value = data.readme_doc_search.example
=======
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
  type     = "basic"
}

resource "readme_doc" "test2" {
  title    = "Testing Colorado"
  body     = "Tiger Pear Giraffe Grape Lemon"
  category = readme_category.test1.id
  type     = "basic"
}

resource "readme_doc" "test3" {
  title    = "Testing Virginia"
  body     = "Koala Pineapple Zebra Plum Avocado"
  category = readme_category.test2.id
  type     = "basic"
>>>>>>> 2c212af (wip)
}

