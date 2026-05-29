---
subcategory: "硬件安全模块"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_hsm_instances"
sidebar_current: "docs-Alibabacloudstack-datasource-hsm-instances"
description: |-
  查询阿里云密码机实例
---

# alibabacloudstack_hsm_instances

查询阿里云密码机实例。

> `注意` 该数据源用于查询密码机实例，不支持创建、修改或删除操作。

## 示例用法

```hcl

variable "name" {
  default = "test-tf-hsm-instance78983"
}

data "alibabacloudstack_zones" "default" {
  enable_details = true
}

resource "alibabacloudstack_vpc_vpc" "default" {
  vpc_name = "${var.name}_vpc"
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vpc_vswitch" "default" {
  name = "${var.name}_vsw"
  vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
  cidr_block = "172.16.0.0/24"
  zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
}

resource "alibabacloudstack_hsm_instance" "default" {
  product_code = "jnta.SJJ1528"
  vendor_code  = "jnta"
  vsm_type     = "gvsm"
  zone_no      = data.alibabacloudstack_zones.default.zones.0.id
  remark       = var.name
  vpc_id       = "${alibabacloudstack_vpc_vpc.default.id}"
  vswitch_id   ="${alibabacloudstack_vpc_vswitch.default.id}"
  ip           = "172.16.1.100"
  white_list   = "192.168.1.0/24"
}

data "alibabacloudstack_hsm_instances" "default" {
  status     = alibabacloudstack_hsm_instance.default.status
  name_regex = "test-tf-hsm-instance"
}

```

## 参数说明
以下参数支持过滤查询结果：

* `cluster_id` (字符串, 可选) - 密码机实例所属集群的ID。

* `ids` (列表, 可选) - 实例ID列表，用于过滤查询结果。

* `instance_id` (字符串, 可选) - 密码机实例ID。

* `name_regex` (字符串, 可选) - 用于通过实例名称（Remark）过滤结果的正则表达式。

* `status` (字符串, 可选) - 密码机实例状态。多个状态可以用逗号分隔。取值：1（未初始化）、3（已释放）、4（生产失败）、5（运行中）、6（同步中）、7（重置中）、8（已停用）。

* `vpc_ip` (字符串, 可选) - 密码机实例所注册的经典网络IP地址。

* `vsm_type` (字符串, 可选) - 密码机实例的设备类型。取值：evsm（金融数据密码机）、gvsm（通用服务器密码机）、svsm（签名验签服务器密码机）。

* `zone_no` (字符串, 可选) - 密码机实例所在的可用区编号。

* `enable_details` (布尔, 可选, 默认: false) - 是否检索每个实例的详细信息。

## 属性说明
以下属性被导出：

* `id` (字符串) - 资源ID。

* `cluster_id` (字符串) - 密码机实例所属集群的ID。

* `cluster_name` (字符串) - 密码机实例所属集群的名称。

* `hsm_id` (字符串) - 与实例关联的物理HSM ID。

* `hsm_status` (整数) - 密码机实例状态。取值：1（未初始化）、3（已释放）、4（生产失败）、5（运行中）、6（同步中）、7（重置中）、8（已停用）。

* `ip` (字符串) - 与密码机实例关联的经典网络IP地址。

* `is_master` (整数) - 当前密码机在集群中是否是主密码机。取值：0（不是）、1（是）。

* `product_code` (字符串) - 密码机实例所用设备的设备型号编码。

* `product_name` (字符串) - 密码机实例所用设备的设备型号名称。

* `release_protection` (整数) - 是否启用释放保护。取值：0（禁用）、1（启用）。

* `remark` (字符串) - 密码机实例的别名或备注。

* `show_create_cluster` (布尔) - 是否允许为此密码机实例创建集群。

* `vpc_id` (字符串) - 分配给密码机实例的VPC ID。

* `vendor_code` (字符串) - 密码机实例所用设备的设备厂商编码。

* `vendor_name` (字符串) - 密码机实例所用设备的设备厂商名称。

* `vsm_type` (字符串) - 密码机实例的设备类型。取值：evsm（金融数据密码机）、gvsm（通用服务器密码机）、svsm（签名验签服务器密码机）。

* `vswitch_id` (字符串) - 分配给密码机实例的交换机ID。

* `white_list` (列表) - 密码机实例的白名单IP列表。如果实例属于集群，则白名单将与集群保持一致。

* `zone_no` (字符串) - 密码机实例所在的可用区编号。