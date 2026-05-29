---
subcategory: "NAT网关"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_common_bandwidth_package"
sidebar_current: "docs-Alibabacloudstack-resource-common-bandwidth-package"
description: |-
  Provides a common Bandwidth package resource.
---

# alibabacloudstack\_common\_bandwidth_package

Provides a common Bandwidth package resource.

## 示例用法
```
variable "name" {
	default = "tf-testaccnat_gatewaybandwi59570"
}

data "alibabacloudstack_zones" "default" {
	available_resource_creation = "VSwitch"
}

resource "alibabacloudstack_vpc" "default" {
	name = "${var.name}"
	cidr_block = "172.16.0.0/12"
}

resource "alibabacloudstack_vswitch" "default" {
	vpc_id = "${alibabacloudstack_vpc.default.id}"
	cidr_block = "172.16.0.0/21"
	availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
	name = "${var.name}"
}

resource "alibabacloudstack_nat_gateway" "default" {
	vpc_id = "${alibabacloudstack_vswitch.default.vpc_id}"
	name = "${var.name}"
}


resource "alibabacloudstack_natgateway_bandwidth_package" "default" {
  name = "tf-testaccnat_gatewaybandwi59570"
  bandwidth = "5"
  natgateway_id = "${alibabacloudstack_nat_gateway.default.id}"
  description = "tf-testaccnat_gatewaybandwi59570"
  ip_count = "2"
}
```

## 参数参�?

支持以下参数�?
  * `name` - (选填) - 共享带宽的名称�?
  * `bandwidth` - (必填) -  共享带宽的带宽峰值， 单位：Mbps�?
  * `natgateway_id` - (必填) - natgateway �?id
  * `description` - (选填) -  共享带宽的描述信息�?
  * `ip_count` - (必填) - EIP的数�?
  * `status` - (选填) - 共享带宽实例的状态。默认取值：**Available**�?

## 属性参�?

除了上述所有参数外，还导出了以下属性：
  * `name` - 共享带宽的名称�?
  * `bandwidth_package_id` -  共享带宽的ID�?
  * `business_status` - 共享带宽实例的状态。取值：- **Normal**：正常状态�? **FinancialLocked**：欠费�? **Unactivated**：未激活�?
  * `description` -  共享带宽的描述信息�?
  * `instance_charge_type` - 共享带宽实例的计费类型。取值：<props="china">**PostPaid**：按量计费�?/props><props="china">**PrePaid**：包年包月�?/props><props="intl">**PostPaid**：按量计费�?/props>
  * `isp` - 线路类型，取值：- **BGP**：BGP（多线）线路�? **BGP_PRO**：BGP（多线）精品线路�?
  * `public_ip_addresses` - 共享带宽实例中的公网IP地址�?
    * `allocation_id` - 公网IP的实例ID�?
    * `ip_address` - 公网IP地址�?
  * `status` - 共享带宽实例的状态。默认取值：**Available**�?
