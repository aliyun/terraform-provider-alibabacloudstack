---
subcategory: "MaxCompute"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_maxcompute_cus"
sidebar_current: "docs-alibabacloudstack-datasource-maxcompute-cus"
description: |-
  Provides a datasource of Max Compute Cus
---

# alibabacloudstack_maxcompute_cus

This data source provides Max Compute Cus


## Example Usage

```hcl
data "alibabacloudstack_maxcompute_cus" "example" {
  name_regex = "example-cu"
}

output "cus" {
  value = data.alibabacloudstack_maxcompute_cus.example.cus
}
```

## Argument Reference
The following arguments are supported:

* `ids` - (Optional) A list of CU IDs to filter the results.
* `name_regex` - (Optional) A regex pattern to filter CUs by name.
* `cluster_name` - (Optional) The name of the cluster to filter CUs.

## Attributes Reference
The following attributes are exported:

* `ids` - A list of IDs of the cluster.
* `cus` - A list of CUs. Each element contains the following attributes:
    * `id` - The unique identifier of the cluster.
    * `cu_name` - The name of the CU.
    * `cu_num` - The number of CUs.
    * `cluster_name` - The name of the cluster.
