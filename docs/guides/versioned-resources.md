---
page_title: Versioned Resources
---

API Specifications, Categories, and Docs are versioned in ReadMe.

ReadMe creates a `1.0.0` automatically and sets it as the default version.
Any resources deployed that don't specify a version will use the default
version.

## Creating New Versions

When a new version is created in ReadMe, it is created _from_ another source
version. This copies all of the source version's resources to the new version,
making them unmanaged by Terraform.

To manage resources entirely with Terraform, consider creating a hidden version
to use as a template for new versions. This template version should be
completely empty - no categories, docs, or API specifications.

```terraform
resource "readme_version" "template" {
  from          = "1.0.0"
  version       = "0"
  is_hidden     = true
  codename      = "Template Version"
}
```

Ensure the `v0` version is hidden and empty.

Managed versions should refer to `0` as the source version.

```terraform
resource "readme_version" "v1_1_0" {
  from          = readme_version.template.version_clean # or "0"
  version       = "1.1.0"
  is_hidden     = false
  is_beta       = false
  is_deprecated = false
  is_stable     = false
  codename      = "current"
}
```

## Versioned Resources

```terraform
resource "readme_category" "example" {
  title   = "My Example Docs"
  type    = "guide"
  version = readme_version.v1_1_0.version_clean
}

resource "readme_doc" "example" {
  title         = "Test Doc"
  version       = readme_version.v1_1_0.version_clean
  category_slug = readme_category.example.slug


  hidden = false
  order  = 1
  type   = "basic"

  body = "This is a test doc"
}

resource "readme_api_specification" "example" {
  definition      = file("${path.module}/petstore.json")
  semver          = readme_version.v1_1_0.version_clean
  delete_category = true
}
```

## Organizing Versioned Content

### Wrapper Module

Consider creating a wrapper module for managing versioned resources.

```terraform
module "v1_1_0" {
  source = "./modules/version"

  version = "1.1.0"
  guides  = fileset("${path.module}/versions/v1.1/guides", "*.md")
  specs   = fileset("${path.module}/versions/v1.1/specs", "*.json")
}
```

### Version-Specific Modules

Another pat
