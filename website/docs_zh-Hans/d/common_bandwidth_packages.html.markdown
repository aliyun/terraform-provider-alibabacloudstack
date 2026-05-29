---
subcategory: "NAT网关"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_common_bandwidth_packages"
sidebar_current: "docs-Alibabacloudstack-datasource-natgateway-bandwidthpackages"
description: |-
  提供阿里云账号下拥有的natgateway bandwidthpackages列表�?
---

# alibabacloudstack\_common\_bandwidth_packages

此数据源提供根据指定过滤条件列出的阿里云账号下的natgateway bandwidthpackages资源列表�?

## 示例用法
```
variable "name" {
  default = "tf-testAccNatGatewaysBandwidthPackagesDatasource19756"
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
    name = "${var.name}"
	bandwidth = "5"
    natgateway_id = "${alibabacloudstack_nat_gateway.default.id}"
	description = "${var.name}"
	ip_count = "2"
}

data "alibabacloudstack_natgateway_bandwidth_packages" "default" {
	name_regex = "${alibabacloudstack_natgateway_bandwidth_package.default.name}"
}
```

## 参数参�?
以下参数是支持的�?
  * `ids` - (选填) - 共享带宽的id列表，用于过滤结果�?
  * `name_regex` - (选填) -  用于按实例名称过滤结果的正则表达式字符串�?
  * `description_regex` - (选填) -  用于按实例描述过滤结果的正则表达式字符串�?

## Attributes Reference
除了上述参数外，还导出以下属性：
  * `bandwidth_packages` - 共享宽带列表�?
    * `id` - 共享带宽的id�?
    * `bandwidth` -  共享带宽的带宽峰值， 单位：Mbps�?
    * `bandwidth_package_id` -  共享带宽的ID�?
    * `business_status` - 共享带宽实例的状态。取值：- **Normal**：正常状态�? **FinancialLocked**：欠费�? **Unactivated**：未激活�?
    * `creation_time` - 共享宽带的创建时间�?
    * `name` - 共享宽带的名称�?
    * `description` -  共享带宽的描述信息�?
    * `instance_charge_type` - 共享带宽实例的计费类型。取值：<props="china">**PostPaid**：按量计费�?/props><props="china">**PrePaid**：包年包月�?/props><props="intl">**PostPaid**：按量计费�?/props>
    * `internet_charge_type` - 共享带宽实例的计费类型�?
    * `natgateway_id` - natgateway的id
    * `ip_count` - EIP的数量�?
    * `isp` - 线路类型，取值：- **BGP**：BGP（多线）线路�? **BGP_PRO**：BGP（多线）精品线路�?
    * `public_ip_addresses` - 共享带宽实例中的公网IP地址�?
    * `status` - 共享带宽实例的状态。默认取值：**Available**�?
