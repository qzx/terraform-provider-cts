terraform {
  required_providers {
    cts = {
      versions = ["0.0.1"]
      source = "qzx/cts"
    }
  }
}

provider "cts" {}

resource "cts_task" "this" {
  name = "cts-task-1"
  description = "Tasks description"
  module = "/module"
  providers = ["consul", "vault"]
  enabled = true
  condition_kv_path = "meta/v2/customers"
  condition_kv_recurse = true
  condition_kv_use_as_modele_input = true
}

data "cts_task" "this" {
  name = "cts-task-1"
}