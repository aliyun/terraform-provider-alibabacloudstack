---
subcategory: "硬件安全模块"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_hsm_instance"
sidebar_current: "docs-Alibabacloudstack-resource-hsm-instance"
description: |-
  创建和管理密码机实例
---

# alibabacloudstack_hsm_instance

创建和管理阿里云密码机实例。

## 示例用法

### 基础用法

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
```

## 参数说明

支持以下参数：

* `product_code` - (必填, 变更时重建) 密码机实例所用设备的设备型号编码。您可以通过DescribeProducts接口获取该参数。
* `vendor_code` - (必填, 变更时重建) 密码机实例所用设备的设备厂商编码。您可以通过DescribeVendors接口获取该参数。
* `vsm_type` - (必填, 变更时重建) 密码机实例的设备类型。取值：
  * `evsm`：金融数据密码机。
  * `gvsm`：通用服务器密码机。
  * `svsm`：签名验签服务器密码机。
* `zone_no` - (必填, 变更时重建) 密码机实例所在的可用区编号。您可以通过DescribeZones接口获取该参数。
* `hsm_id` - (可选, 可回读) 密码机ID。当指定此参数时，将创建指定设备的密码机实例。此属性由 API 返回，无法手动设置。
* `ip` - (可选, 可回读) 密码机实例对应的经典网络的IP地址。不传该参数时，系统将根据VPC和VSwitch自动分配一个可用IP。此属性由 API 返回，无法手动设置。
* `remark` - (可选) 密码机实例的别名。
* `vpc_id` - (可选, 可回读) 配置给密码机实例的VPC实例ID。您可以通过DescribeVpc接口获取该参数。此属性由 API 返回，无法手动设置。
* `vswitch_id` - (可选, 可回读) 配置给密码机实例的交换机实例ID。您可以通过DescribeVpc接口获取该参数。此属性由 API 返回，无法手动设置。
* `white_list` - (可选, 可回读) 可访问密码机实例的白名单IP。支持配置多个值，可以使用","分隔。此属性由 API 返回，无法手动设置。

## 属性说明

以下属性会从API响应中导出：

* `id` - 密码机实例ID。
* `cluster_id` - 密码机实例所在集群的ID。
* `cluster_name` - 密码机实例所在集群的名称。
* `instance_id` - 密码机实例ID（与id相同）。
* `is_master` - 当前密码机在集群中是否是主密码机。取值：
  * `0`：不是。
  * `1`：是。
* `product_name` - 密码机实例所用设备的设备型号名称。
* `release_protection` - 释放保护状态。
* `show_create_cluster` - 是否可以创建集群。取值：
  * `true`：是。
  * `false`：否。
* `status` - 密码机实例状态。取值：
  * `1`：未初始化。
  * `3`：已释放。
  * `4`：生产失败。
  * `5`：运行中。
  * `6`：同步中。
  * `7`：重置中。
  * `8`：已停用。
* `vendor_name` - 密码机实例所用设备的设备厂商名称。

## Import

密码机实例可以使用 instance_id 导入，例如：

```
$ terraform import alibabacloudstack_hsm_instance.example hsm-12345678
```