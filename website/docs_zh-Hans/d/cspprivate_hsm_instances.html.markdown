---
subcategory: "云密码机"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cspprivate_hsm_instances"
sidebar_current: "docs-alibabacloudstack-datasource-cspprivate-hsm-instances"
description: |-
  查询阿里云密码服务(HSM)实例。
---

# alibabacloudstack_cspprivate_hsm_instances

查询阿里云密码服务(HSM)实例。

> **注意：** 此数据源仅用于查询HSM实例，不支持创建、修改或删除操作。

## 典型用法

```hcl
variable "name" {
  default = "test-tf-cspprivate-hsm"
}

data "alibabacloudstack_zones" "default" {
  enable_details = true
}

data "alibabacloudstack_cspprivate_hsm_vendors" "default" {
}

resource "alibabacloudstack_cspprivate_hsm_instance" "default" {
  product_code = "${data.alibabacloudstack_cspprivate_hsm_vendors.default.vendors.0.products.0.code}"
  vendor_code  = "${data.alibabacloudstack_cspprivate_hsm_vendors.default.vendors.0.code}"
  vsm_type     = "gvsm"
  zone_id      = "${data.alibabacloudstack_zones.default.zones.0.id}"
  alias_name   = "${var.name}"
}

data "alibabacloudstack_cspprivate_hsm_instances" "default" {
  name_regex = "${alibabacloudstack_cspprivate_hsm_instance.default.alias_name}"
  ids        = ["${alibabacloudstack_cspprivate_hsm_instance.default.id}"]
}
```

## 参数说明

以下参数可用于过滤查询结果：

* `ids` - (可选) 用于过滤结果的实例ID列表。
* `name_regex` - (可选) 用于按实例名称(备注)过滤结果的正则表达式。

## 属性说明

以下属性会被导出：

* `id` - 数据源的ID。
* `ids` - 实例ID列表。
* `names` - 与返回实例对应的实例名称(备注)列表。
* `instances` - 符合过滤条件的HSM实例列表。每个元素包含以下属性：
  * `id` - 资源ID。
  * `instance_id` - HSM实例ID。
  * `hsm_status` - HSM实例的状态。1: 未初始化, 3: 已释放, 4: 失败, 5: 运行中, 6: 同步中, 7: 重置中, 8: 已禁用。
  * `vpc_id` - 分配给HSM实例的VPC ID。
  * `vswitch_id` - 分配给HSM实例的vSwitch ID。
  * `ip` - 与HSM实例关联的经典网络IP地址。
  * `alias_name` - HSM实例的别名或名称。
  * `product_code` - HSM实例所用设备的产品型号代码。
  * `vendor_code` - HSM实例所用设备的供应商代码。
  * `vsm_type` - HSM实例的类型。evsm: 金融数据HSM, gvsm: 通用服务器HSM, svsm: 签名验证服务器HSM。
  * `zone_id` - HSM实例所在的可用区ID。