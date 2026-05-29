---
subcategory: "云密码机"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cspprivate_hsm_vendors"
sidebar_current: "docs-alibabacloudstack-datasource-cspprivate-hsm-vendors"
description: |-
  Query Alibaba Cloud Hardware Security Module (HSM) vendors and their products.
---

# alibabacloudstack_cspprivate_hsm_vendors

Query Alibaba Cloud Hardware Security Module (HSM) vendors and their products.

> **NOTE:** This data source is used to query HSM vendors and does not support create, modify, or delete operations.

## Example Usage

```hcl
data "alibabacloudstack_cspprivate_hsm_vendors" "default" {
}

output "first_vendor_code" {
  value = data.alibabacloudstack_cspprivate_hsm_vendors.default.vendors.0.code
}

output "first_vendor_name" {
  value = data.alibabacloudstack_cspprivate_hsm_vendors.default.vendors.0.name
}

output "first_vendor_first_product_code" {
  value = data.alibabacloudstack_cspprivate_hsm_vendors.default.vendors.0.products.0.code
}

output "first_vendor_first_product_name" {
  value = data.alibabacloudstack_cspprivate_hsm_vendors.default.vendors.0.products.0.name
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of vendor codes used to filter query results.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID.
* `vendors` - A list of HSM vendors. Each element contains the following attributes:
  * `code` - The code of the HSM vendor.
  * `name` - The name of the HSM vendor.
  * `products` - A list of HSM products. Each element contains the following attributes:
    * `code` - The code of the HSM vendor product.
    * `name` - The name of the HSM vendor product.
```