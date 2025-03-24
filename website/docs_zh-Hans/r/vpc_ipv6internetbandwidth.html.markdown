---
subcategory: "VPC"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_vpc_ipv6internetbandwidth"
sidebar_current: "docs-Alibabacloudstack-vpc-ipv6internetbandwidth"
description: |- 
  集编排VPC的IPv6网络带宽
---

# alibabacloudstack_vpc_ipv6internetbandwidth
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_vpc_ipv6_internet_bandwidth`

使用Provider配置的凭证在指定的资源集编排VPC的IPv6网络带宽。

## 示例用法

以下示例展示了如何配置 `alibabacloudstack_vpc_ipv6_internet_bandwidth` 资源的### 基础用法：

```terraform
variable "name" {
  default = "tf-testaccvpcipv6internetbandwidth46716"
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
```

### 示例说明

1. **变量定义**：通过 `variable` 定义了变量 `name`，用于命名资源。
2. **数据源引用**：
   - 使用 `data "alibabacloudstack_instances"` 获取运行中的实例信息。
   - 使用 `data "alibabacloudstack_vpc_ipv6_addresses"` 获取与实例关联的可用 IPv6 地址。
3. **资源配置**：
   - `ipv6_address_id` 和 `ipv6_gateway_id` 分别从 `data "alibabacloudstack_vpc_ipv6_addresses"` 中获取。
   - 配置 `internet_charge_type` 为 `PayByBandwidth`，并设置带宽为 `20 Mbps`。

## 参数参考

支持以下参数：

* `bandwidth` - (必填) IPv6 地址的互联网带宽资源量，单位：`Mbit/s`。有效值范围：`1` 到 `5000`。  
  **注意**：如果将 `internet_charge_type` 设置为 `PayByTraffic`，则 IPv6 地址的互联网带宽资源受 IPv6 网关规格的限制：
  * `Small`(默认值)：指定免费版，互联网带宽为 `1` 到 `500` Mbit/s。
  * `Medium`：指定中型版，互联网带宽为 `1` 到 `1000` Mbit/s。
  * `Large`：指定大型版，互联网带宽为 `1` 到 `2000` Mbit/s。

* `internet_charge_type` - (可选，变更时重建) IPv6 网关的互联网带宽资源的计费方式。有效值：`PayByBandwidth`，`PayByTraffic`。

* `ipv6_address_id` - (必填，变更时重建) IPv6 地址实例的 ID。

* `ipv6_gateway_id` - (必填，变更时重建) 属于该 IPv6 地址的 IPv6 网关的 ID。

### 参数详细说明

- **`bandwidth`**：指定 IPv6 地址的独享公网带宽，单位为 Mbps。必须在有效范围内配置，具体范围取决于 `internet_charge_type` 的值。
- **`internet_charge_type`**：定义带宽的计费方式。可以选择按固定带宽计费(`PayByBandwidth`)或按流量计费(`PayByTraffic`)。
- **`ipv6_address_id`**：指定要绑定带宽的 IPv6 地址实例的唯一标识符。
- **`ipv6_gateway_id`**：指定与 IPv6 地址关联的网关实例的唯一标识符。

## 属性参考

除了上述所有参数外，还导出了以下属性：

* `id` - Terraform 中 IPv6 Internet Bandwidth 资源的唯一标识符。
* `status` - 资源的状态。有效值：`Normal`、`FinancialLocked` 和 `SecurityLocked`。
* `internet_charge_type` - IPv6 网关的互联网带宽资源的计费方式。有效值：`PayByBandwidth`，`PayByTraffic`。

### 属性详细说明

- **`id`**：Terraform 自动生成的资源唯一标识符，用于区分不同的 IPv6 Internet Bandwidth 资源。
- **`status`**：表示资源当前的状态。可能的值包括：
  - `Normal`：资源正常运行。
  - `FinancialLocked`：由于财务问题(如欠费)，资源被锁定。
  - `SecurityLocked`：由于安全原因，资源被锁定。
- **`internet_charge_type`**：再次确认资源的计费方式，确保配置的一致性。

此文档提供了更详细的参数和属性描述，帮助用户更好地理解和使用 `alibabacloudstack_vpc_ipv6internetbandwidth` 资源。
```