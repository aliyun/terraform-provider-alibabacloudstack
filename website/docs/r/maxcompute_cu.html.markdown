---
subcategory: "MaxCompute"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_maxcompute_cu"
sidebar_current: "docs-apsarastack-resource-maxcompute-cu"
description: |-
  Provides a Apsarastack maxcompute cu resource.
---

# apsarastack\_maxcompute\_project

The cu is the basic unit of operation in maxcompute. 

->**NOTE:** Available in 1.77.0+.

## Example Usage

Basic Usage

```terraform
resource "apsarastack_maxcompute_cu" "example" {
   cu_name      = "tf_testAccApsaraStack7898"
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
$ terraform import apsarastack_maxcompute_project.example tf_maxcompute_project
```
