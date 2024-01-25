terraform {
  required_providers {
    customs3 = {
      source = "hashicorp.com/edu/custom-s3"
    }
  }
}

provider "customs3" {
  region     = "<region>"
  access_key = "<access_key>"
  secret_key = "<secret_key>"
}

data "customs3_buckets" "example" {}

output "all_buckets" {
  value = data.customs3_buckets.example
}