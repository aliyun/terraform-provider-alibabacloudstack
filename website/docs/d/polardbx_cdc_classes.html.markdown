---
subcategory: "Cloud-Native Distributed Database PolarDB-X 2.0"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_polardbx_cdc_classes"
sidebar_current: "docs-Alibabacloudstack-datasource-polardbx-cdc-classes"
description: |-
  Query the available CDC (Change Data Capture) classes for a PolarDB-X instance.
---

# alibabacloudstack_polardbx_cdc_classes

Queries the available CDC (Change Data Capture) classes for a specified PolarDB-X instance. CDC classes define the CPU and memory configurations that can be used for CDC nodes.

## Example Usage

```hcl
data "alibabacloudstack_polardbx_cdc_classes" "example" {
  instance_id = "pxc-example-instance"
  cpu         = 4
  memory      = 16
}

output "cdc_classes" {
  value = data.alibabacloudstack_polardbx_cdc_classes.example.cdc_classes
}
```

## Argument Reference
The following arguments are supported:

* `instance_id` - (Required) The ID of the PolarDB-X instance.
* `ids` - (Optional) A list of CDC class IDs to filter the results.
* `cpu` - (Optional) The number of CPU cores to filter CDC classes.
* `memory` - (Optional) The memory size (in GB) to filter CDC classes.
* `sorted_by` - (Optional) The sorting method for the returned CDC classes. Valid values: `CPU`, `Memory`.

## Attributes Reference
In addition to the arguments listed above, the following attributes are exported:

* `ids` - A list of CDC class IDs.
* `cdc_classes` - A list of CDC class configurations, containing the following attributes:
  * `id` - The CDC class unique identifier (ClassCode).
  * `cpu` - The number of CPU cores.
  * `memory` - The memory size (in GB).
