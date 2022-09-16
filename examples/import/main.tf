terraform {
  required_providers {
    cts = {
      versions = ["0.0.1"]
      source = "qzx/hashicups"
    }
  }
}

provider "cts" {}