---
subcategory: "ASCM"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ram_service_role_products"
sidebar_current: "docs-Alibabacloudstack-datasource-ram-service-role-products"
description: |- 
 查询RAM服务角色产品
---

# alibabacloudstack_ram_service_role_products

根据指定过滤条件列出当前凭证权限可以访问的RAM服务角色产品列表。

## 示例用法

```hcl
data "alibabacloudstack_ram_service_role_products" "example" {
  name_regex = "example-product"
}

output "products" {
  value = data.alibabacloudstack_ram_service_role_products.example.products
}
```

## 参数参考
支持以下参数：

* `name_regex` - (可选) 用于通过其 ASCII 名称过滤服务角色产品的正则表达式模式。

## 属性参考
导出以下属性：

* `products` - RAM 服务角色产品的列表。每个元素包含以下属性：
    * `chinese_name` - 服务角色产品的中文名称。
    * `ascii_name` - 服务角色产品的 ASCII 名称。
    * `key` - 服务角色产品的键标识符。