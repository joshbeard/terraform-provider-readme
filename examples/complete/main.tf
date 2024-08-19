terraform {
  required_providers {
    readme = {
      source = "liveoaklabs/readme"
    }
  }
}

variable "api_token" {
  type        = string
  description = "The API token for the ReadMe provider"
}

provider "readme" {
  api_token = var.api_token
}


resource "readme_version" "example" {
  from          = "1.0.0"
  version       = "1.0.1"
  is_hidden     = false
  is_beta       = false
  is_deprecated = false
  is_stable     = false
  codename      = "next"
}

resource "readme_category" "example" {
  title   = "My Example Docs"
  type    = "guide"
  version = readme_version.example.version_clean
}

resource "readme_custom_page" "this" {
  title  = "My custom page"
  body   = "This is a test custom page"
  hidden = false

  # html      = "<h1>My custom page</h1>"
  # html_mode = true
}

resource "readme_doc" "example" {
  title   = "Test Doc"
  version = readme_version.example.version_clean

  category_slug = readme_category.example.slug
  # category = <category_id>

  hidden = false
  order  = 1
  type   = "basic"

  body = "This is a test doc"

  # error = {}

  # parent_doc      = "<parent doc id>"
  # parent_doc_slug = "<parent doc slug>"
  # verify_parent_doc = true

  # use_slug = "existing-doc-slug"
}

resource "readme_api_specification" "example" {
  definition      = file("${path.module}/petstore.json")
  semver          = readme_version.example.version_clean
  delete_category = true
}

resource "readme_changelog" "example" {
  body   = file("${path.module}/changelog/2024-08-19.md")
  hidden = false

  # The title is set in frontmatter
  # title = "2024-08-19"

  # type = "added"
}

# Note: Images aren't truly stateful. If the source changes, a new image will
# be created with a new URL. The API does not support deleting existing images.
resource "readme_image" "example" {
  source = "${path.module}/../../.github/readme/lob-logo.png"
}
