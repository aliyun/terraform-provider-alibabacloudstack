---
subcategory: "MaxCompute"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_maxcompute_cu"
sidebar_current: "docs-alibabacloudstack-resource-maxcompute-cu"
description: |-
  Provides a Alibabacloudstack maxcompute cu resource.
---

# alibabacloudstack\_maxcompute\_project

The cu is the basic unit of operation in maxcompute. 

->**NOTE:** Available in 1.77.0+.

## Example Usage

Basic Usage

```terraform
resource "alibabacloudstack_maxcompute_cu" "example" {
   cu_name      = "tf_testAccAlibabacloudStack7898"
   cu_num       = "1"
   cluster_name = "HYBRIDODPSCLUSTER-A-20210520-07B0"
}
```
## Argument Reference

The following arguments are supported:
* `id` - (Required, ForceNew) The id of the maxcompute cu.
* `cu_name` - (Required, ForceNew, Available in 1.110.0+) The name of the maxcompute cu.
* `cu_num` - (Required, ForceNew) The name of the maxcompute cu.
* `cluster_name` - (Required, ForceNew) The cluster name of the maxcompute cu.


## Import

MaxCompute project can be imported using the *name* or ID, e.g.

```
$ terraform import alibabacloudstack_maxcompute_project.example tf_maxcompute_project
```
