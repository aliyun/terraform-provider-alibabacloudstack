---
subcategory: "ASCM"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ascm_service_cluster_by_product"
sidebar_current: "docs-alibabacloudstack-datasource-ascm-service-clusters-by-product"
description: |-
    Provides a list of service cluster to the user.
---

# alibabacloudstack\_ascm_service_clusters_by_product

This data source provides the service clusters of the current Apsara Stack Cloud user.

## Example Usage

```
data "alibabacloudstack_ascm_service_cluster_by_product" "cluster" {
  output_file = "cluster"
  product_name = "slb"
}

output "cluster" {
  value = data.alibabacloudstack_ascm_service_cluster_by_product.cluster.*
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of instance family IDs.
* `product_name` - (Required) Filter the results by specifying name of the service.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `cluster_list` - A list of instance families. Each element contains the following attributes:
    * `cluster_by_region` - cluster by a region.
