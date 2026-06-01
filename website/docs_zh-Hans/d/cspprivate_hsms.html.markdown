---
subcategory: "云密码机"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cspprivate_hsms"
sidebar_current: "docs-alibabacloudstack-datasource-cspprivate-hsms"
description: |-
  提供阿里云专有云 CSP 密码机实例列表。
---

# alibabacloudstack\_cspprivate_hsms

该数据源提供阿里云专有云中 CSP 密码机（Hardware Security Module, HSM）实例列表。

## 示例用法

```hcl
data "alibabacloudstack_cspprivate_hsm_vendors" "default" {
}

data "alibabacloudstack_zones" "default" {
  enable_details = true
}

data "alibabacloudstack_cspprivate_hsms" "default" {
  vendor_code  = data.alibabacloudstack_cspprivate_hsm_vendors.default.vendors[0].code
  product_code = data.alibabacloudstack_cspprivate_hsm_vendors.default.vendors[0].products[0].code
  zone_id      = data.alibabacloudstack_zones.default.zones[0].id
}

output "first_hsm_id" {
  value = data.alibabacloudstack_cspprivate_hsms.default.hsms[0].id
}
```

## 参数说明

支持以下参数：

* `zone_id` - (必填) HSM 实例部署的可用区 ID。
* `vendor_code` - (必填) HSM 厂商代码。
* `product_code` - (必填) HSM 产品代码。

## 属性说明

以下属性导出为数据源属性：

* `ids` - CSP 密码机实例 ID 列表。
* `hsms` - CSP 密码机实例列表。每个元素包含以下属性：
  * `id` - CSP 密码机实例的 ID。
