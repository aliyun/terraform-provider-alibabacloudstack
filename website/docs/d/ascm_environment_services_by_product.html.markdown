---
subcategory: "ASCM"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_ascm_environment_services_by_product"
sidebar_current: "docs-apsarastack-datasource-ascm-environment-services"
description: |-
    Provides a list of environment services to the user.
---

# apsarastack\_ascm_environment_services_by_product

This data source provides the environment services of the current Apsara Stack Cloud user.

## Example Usage

```
data "apsarastack_ascm_environment_services_by_product" "default" {
  output_file = "environment"
}
output "envser" {
  value = data.apsarastack_ascm_environment_services_by_product.default.*
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of environment service IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `result` - A list of environment services. Each element contains the following attributes:  
