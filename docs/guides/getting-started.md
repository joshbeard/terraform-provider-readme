---
page_title: Getting Started
---

## Configure the Provider

```terraform
terraform {
  required_providers {
    readme = {
      source = "liveoaklabs/readme"
    }
  }
}

provider "readme" {
  api_token = "<API_TOKEN>"
}
```

## Manage Resources

```terraform
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

resource "readme_doc" "example" {
  title         = "Test Doc"
  version       = readme_version.example.version_clean
  category_slug = readme_category.example.slug


  hidden = false
  order  = 1
  type   = "basic"

  body = "This is a test doc"
}

resource "readme_api_specification" "example" {
  definition      = file("${path.module}/petstore.json")
  semver          = readme_version.example.version_clean
  delete_category = true
}
```
