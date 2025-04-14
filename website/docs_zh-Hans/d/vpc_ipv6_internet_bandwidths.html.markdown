---
subcategory: "VPC"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_vpc_ipv6_internet_bandwidths"
sidebar_current: "docs-Alibabacloudstack-datasource-vpc-ipv6-internet-bandwidths"
description: |- 
  查询专有网络（VPC）IPv6网络带宽
---

# alibabacloudstack_vpc_ipv6_internet_bandwidths
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_vpc_ipv6_internetbandwidths`

根据指定过滤条件列出当前凭证权限可以访问的专有网络（VPC）IPv6网络带宽列表。

## 示例用法

```hcl
variable "name" {
  default = "tf-testacc-vpcipv6internetbandwidth-7579550"
}

data "alibabacloudstack_instances" "default" {
  name_regex = "no-deleteing-ipv6-address"
  status     = "Running"
}

data "alibabacloudstack_vpc_ipv6_addresses" "default" {
  associated_instance_id = data.alibabacloudstack_instances.default.instances.0.id
  status                 = "Available"
}

resource "alibabacloudstack_vpc_ipv6_internet_bandwidth" "default" {
  ipv6_address_id      = data.alibabacloudstack_vpc_ipv6_addresses.default.addresses.0.id
  ipv6_gateway_id      = data.alibabacloudstack_vpc_ipv6_addresses.default.addresses.0.ipv6_gateway_id
  internet_charge_type = "PayByBandwidth"
  bandwidth            = "20"
}

data "alibabacloudstack_vpc_ipv6_internet_bandwidths" "default" {
  ids                  = [alibabacloudstack_vpc_ipv6_internet_bandwidth.default.id]
  ipv6_internet_bandwidth_id = alibabacloudstack_vpc_ipv6_internet_bandwidth.default.id
  ipv6_address_id      = data.alibabacloudstack_vpc_ipv6_addresses.default.addresses.0.id
  status               = "Normal"
  output_file          = "output.txt"
}

output "vpc_ipv6_internet_bandwidth_id_1" {
  value = data.alibabacloudstack_vpc_ipv6_internet_bandwidths.default.bandwidths.0.id
}
```

## 参数说明

以下参数是支持的：

* `ids` - (选填, 变更时重建) IPv6 Internet 带宽 ID 列表。可以通过此参数筛选特定的带宽资源。
* `ipv6_internet_bandwidth_id` - (选填, 变更时重建) IPv6 Internet 带宽的 ID。可以通过此参数精确匹配某个带宽资源。
* `ipv6_address_id` - (选填, 变更时重建) IPv6 地址实例的 ID。可以通过此参数筛选与特定 IPv6 地址关联的带宽资源。
* `status` - (选填, 变更时重建) 资源的状态。有效值：`Normal`、`FinancialLocked` 和 `SecurityLocked`。可以通过此参数筛选处于特定状态的带宽资源。

## 属性说明

除了上述参数外，还导出以下属性：

* `names` - IPv6 Internet 带宽名称列表。
* `bandwidths` - VPC IPv6 Internet 带宽列表。每个元素包含以下属性：
  * `bandwidth` - IPv6 地址的独享公网带宽，单位 Mbps。有效值范围为 `1` 到 `5000` Mbit/s。**注意**：如果 `internet_charge_type` 设置为 `PayByTraffic`，则 IPv6 地址的公网带宽资源将受到 IPv6 网关规格的限制。`Small`（默认）：表示免费版，公网带宽范围为 `1` 到 `500` Mbit/s；`Medium`：表示标准版，公网带宽范围为 `1` 到 `1000` Mbit/s；`Large`：表示高级版，公网带宽范围为 `1` 到 `2000` Mbit/s。
  * `id` - IPv6 Internet 带宽的 ID。
  * `internet_charge_type` - IPv6 网关的互联网带宽资源的计费方式。有效值：`PayByBandwidth`（按固定带宽计费）、`PayByTraffic`（按流量计费）。
  * `ipv6_address_id` - IPv6 地址实例的 ID。
  * `ipv6_gateway_id` - IPv6 地址所属的 IPv6 网关的 ID。
  * `ipv6_internet_bandwidth_id` - IPv6 Internet 带宽的 ID。
  * `payment_type` - 资源的支付类型。
  * `status` - 资源的状态。有效值：`Normal`（正常）、`FinancialLocked`（财务锁定）和 `SecurityLocked`（安全锁定）。