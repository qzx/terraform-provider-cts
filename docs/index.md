---
page_title: "Provider: CTS"
subcategory: ""
description: |-
  Terraform provider for interacting with CTS API.
---

# CTS Provider

Use the navigation to the left to read about the available resources.

## Example Usage

Do not keep your authentication password in HCL for production environments, use Terraform environment variables.

```terraform
provider "cts" {}
```

## Schema

### Optional

- **host** (String, Optional) HashiCups API address (defaults to `localhost:8558`)