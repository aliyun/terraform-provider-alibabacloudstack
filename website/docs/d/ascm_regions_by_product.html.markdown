---
subcategory: "ASCM"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_ascm_regions_by_product"
sidebar_current: "docs-apsarastack-datasource-ascm-regions-by-product"
description: |-
    Provides a list of regions to the user.
---

# apsarastack\_ascm_regions_by_product

This data source provides the regions of the current Apsara Stack Cloud user.

## Example Usage

```
data "apsarastack_ascm_regions_by_product" "regions" {
  output_file = "product_regions"
  product_name = "ecs"
}
output "regions" {
  value = data.apsarastack_ascm_regions_by_product.regions.*
}
```

## Argument Reference

The following arguments are supported:

* `product_name` - (Required) Filter the results by specified The name of the service.
* `organization` - (Optional) Filter the results by the specified name of the organization.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `region_list` - A list of regions. Each element contains the following attributes:
    * `region_id` - ID of the region.
    * `region_type` - type of region.
