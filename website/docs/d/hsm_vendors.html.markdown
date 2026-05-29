---
subcategory: "Key Management Service"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_hsm_vendors"
sidebar_current: "docs-alibabacloudstack-datasource-hsm-vendors"
description: |-
  Query Alibaba Cloud Hardware Security Module (HSM) vendors and their products.
---

# alibabacloudstack_hsm_vendors

Query Alibaba Cloud Hardware Security Module (HSM) vendors and their products.

> `NOTE` This data source is used to query HSM vendors and does not support create, modify, or delete operations.

## Example Usage

```hcl
data "alibabacloudstack_hsm_vendors" "default" {
}

output "first_vendor_code" {
  value = data.alibabacloudstack_hsm_vendors.default.vendors.0.code
}

output "first_vendor_name" {
  value = data.alibabacloudstack_hsm_vendors.default.vendors.0.name
}

output "first_vendor_first_product_code" {
  value = data.alibabacloudstack_hsm_vendors.default.vendors.0.products.0.code
}

output "first_vendor_first_product_name" {
  value = data.alibabacloudstack_hsm_vendors.default.vendors.0.products.0.name
}
```

## Argument Reference

The following arguments are supported:

* `ids` (List, Optional) - A list of vendor codes used to filter query results.

## Attributes Reference

The following attributes are exported:

* `id` (String) - The resource ID.
* `vendors` (List) - A list of HSM vendors. Each element contains the following attributes:
  * `code` (String) - The code of the HSM vendor.
  * `name` (String) - The name of the HSM vendor.
  * `products` (List) - A list of HSM products. Each element contains the following attributes:
    * `code` (String) - The code of the HSM vendor product.
    * `name` (String) - The name of the HSM vendor product.
