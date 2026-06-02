---
subcategory: "MaxCompute"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_maxcompute_cu"
sidebar_current: "docs-Alibabacloudstack-resource-maxcompute-cu"
description: |-
  Provides a Alibabacloudstack maxcompute cu resource.
---

# alibabacloudstack_maxcompute_cu

The cu is the basic unit of operation in maxcompute. 

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
* `cu_name` - (Required, ForceNew) The name of the maxcompute cu. Must be between 3 and 27 characters. <!--  AI CREATE  -->
* `cu_num` - (Required) The number of CUs for the maxcompute cu. Must be at least 1. 
* `cluster_name` - (Required, ForceNew) The cluster name of the maxcompute cu.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the maxcompute cu.

## Import

MaxCompute CU can be imported using the CuId, e.g.

```
$ terraform import alibabacloudstack_maxcompute_cu.example <cu_id>
```