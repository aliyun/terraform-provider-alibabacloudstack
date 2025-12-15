---
subcategory: "Hardware Security Module (HSM)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_hsm_vendors"
sidebar_current: "docs-Alibabacloudstack-datasource-hsm-vendors"
description: |-
  查询阿里云密码机厂商及其产品
---

# alibabacloudstack_hsm_vendors

查询阿里云密码机厂商及其产品。

> `注意` 该数据源用于查询密码机厂商，不支持创建、修改或删除操作。

## 示例用法

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

## 参数说明

以下参数用于查询：

* `ids` (列表, 可选) - 厂商代码列表，用于过滤查询结果。

## 属性说明

以下属性被导出：

* `id` (字符串) - 资源ID。
* `vendors` (列表) - HSM厂商列表。每个元素包含以下属性：
  * `code` (字符串) - HSM厂商代码。
  * `name` (字符串) - HSM厂商名称。
  * `products` (列表) - HSM产品列表。每个元素包含以下属性：
    * `code` (字符串) - HSM厂商产品代码。
    * `name` (字符串) - HSM厂商产品名称。