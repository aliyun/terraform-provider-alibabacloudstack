---
subcategory: "ASCM"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ascm_service_cluster_by_product"
sidebar_current: "docs-alibabacloudstack-datasource-ascm-service-cluster-by-product"
description: |-
    Provides a list of service clusters to the user.
---
# alibabacloudstack_ascm_service_cluster_by_product

This data source provides a list of service clusters to the user.

## Example Usage

```hcl
data "alibabacloudstack_ascm_regions_by_product" "example" {
  product_name = "ecs"
}

output "regions" {
  value = data.alibabacloudstack_ascm_regions_by_product.example.region_list
}
```

## Argument Reference
The following arguments are supported:

* `ids` - (Optional) A list of region IDs to filter the results.
* `product_name` - (Required) The name of the product for which to retrieve regions.
* `organization` - (Optional) The organization for which to retrieve regions.

## Attributes Reference
The following attributes are exported:

* `ids` - A list of IDs of the regions.
* `region_list` - A list of regions. Each element contains the following attributes:
    * `region_id` - The unique identifier of the region.
    * `region_type` - The type of the region.