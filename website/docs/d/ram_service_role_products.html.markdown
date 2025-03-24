---
subcategory: "ASCM"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ram_service_role_products"
sidebar_current: "docs-Alibabacloudstack-datasource-ram-service-role-products"
description: |- 
  Provides a list of RAM service role products.
---

# alibabacloudstack_ram_service_role_products

This data source provides a list of RAM service role products.

## Example Usage

```hcl
data "alibabacloudstack_ram_service_role_products" "example" {
  name_regex = "example-product"
}

output "products" {
  value = data.alibabacloudstack_ram_service_role_products.example.products
}
```

## Argument Reference
The following arguments are supported:

* `name_regex` - (Optional) A regex pattern to filter service role products by their ASCII name.

## Attributes Reference
The following attributes are exported:

* `products` - A list of RAM service role products. Each element contains the following attributes:
    * `chinese_name` - The Chinese name of the service role product.
    * `ascii_name` - The ASCII name of the service role product.
    * `key` - The key identifier of the service role product.
