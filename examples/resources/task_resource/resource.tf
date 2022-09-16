resource "cts_task" "example" {
  name        = "foo"
  description = "Dscription"
  module      = "/module"
  providers   = ["consul"]
  condition = {
    kv = {
      path = "consul/path"
    }
  }
}