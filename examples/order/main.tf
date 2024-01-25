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

resource "customs3_s3_bucket" "example" {
  buckets = [{
    name = "test-bucket-2398756"
    tags = "yourbucket"
  }]
}
