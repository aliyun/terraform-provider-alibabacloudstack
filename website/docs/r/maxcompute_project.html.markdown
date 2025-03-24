---
subcategory: "MaxCompute"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_maxcompute_project"
sidebar_current: "docs-alibabacloudstack-resource-maxcompute-project"
description: |-
  Provides a Alibabacloudstack maxcompute project resource.
---

# alibabacloudstack_maxcompute_project

The project is the basic unit of operation in maxcompute. 

->**NOTE:** Available in 1.77.0+.

## Example Usage

Basic Usage

```terraform
resource "alibabacloudstack_maxcompute_project" "example" {
  project_name       = "tf_maxcompute_project"
  specification_type = "OdpsStandard"
  order_type         = "PayAsYouGo"
}
```
## Argument Reference

The following arguments are supported:
* `name` - (Required, ForceNew) It has been deprecated from provider version 1.110.0 and `project_name` instead.
* `quota_id` - (Required) The quota ID of the maxcompute project. 
* `disk` - (Required) The disk size of the maxcompute project. 

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the maxcompute project. It is the same as its name.
* `name` - The name of the maxcompute project. 

## Import

MaxCompute project can be imported using the *name* or ID, e.g.

```
$ terraform import alibabacloudstack_maxcompute_project.example tf_maxcompute_project
```