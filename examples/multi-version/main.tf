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
