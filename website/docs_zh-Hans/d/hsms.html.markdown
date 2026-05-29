---
subcategory: "硬件安全模块"
layout: "alibabacloudstack"
page_title: "AlibabacloudStack: alibabacloudstack_hsms"
sidebar_current: "docs-alibabacloudstack-datasource-hsms"
description: |-
  提供阿里云专有云中可用的Hsm实例列表。
---

# alibabacloudstack\_hsms

该数据源提供阿里云专有云区域中的HSM实例列表。

## Example Usage

```hcl
data "alibabacloudstack_hsm_vendors" "default" {
}

data "alibabacloudstack_zones" "default" {
  enable_details = true
}

data "alibabacloudstack_hsms" "default" {
  vendor_code = "${data.alibabacloudstack_hsm_vendors.default.vendors[0].code}"
  product_code = "${data.alibabacloudstack_hsm_vendors.default.vendors[0].products[0].code}"
  zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
}
output "first_hsm_id" {
  value = data.alibabacloudstack_hsms.this.hsms.0.id
}
```

## Argument Reference

以下参数是支持的：

* `zone_id` - (必需) HSM实例所在的可用区。
* `vendor_code` - (必需) HSM供应商代码。
* `product_code` - (必需) HSM产品代码。

## Attributes Reference

以下属性会被导出：

* `ids` - HSM实例ID列表。
* `hsms` - HSM实例列表。每个元素包含以下属性：
  * `id` - HSM实例的ID。
