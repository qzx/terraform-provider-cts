---
page_title: "Task Resource - terraform-provider-cts"
subcategory: ""
description: |-
  The task data source allows you to retrieve information all available CTS tasks.
---

# Data Source `tasks`

The task data source allows you to retrieve information all available CTS tasks.

## Example Usage

```terraform
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

```

## Argument Reference

The following attributes are allowed.

- `name` - The task name.
- `description` - The task description.
- `module` - The task module.
- `providers` - The task providers.
- `condition` - The task condition.
- `enabled` - If the task is enabled. 