terraform {
  required_providers {
    cts = {
      versions = ["0.0.1"]
      source = "qzx/hashicups"
    }
  }
}

resource "cts_task" "this" {
  name = "cts-task-1"
  description = "Tasks description"
  module = "/module"
  providers = ["consul", "vault"]
  enabled = true
  condition = {
    kv = {
      path = "cts/consul/path"
      recurse = true
      use_as_module_input = true
    }
  }
}

data "cts_task" "this" {
  name = "cts-task-1"
}