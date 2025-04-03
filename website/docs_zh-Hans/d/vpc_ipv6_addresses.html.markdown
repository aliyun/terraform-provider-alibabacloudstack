---
subcategory: "VPC"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_vpc_ipv6_addresses"
sidebar_current: "docs-Alibabacloudstack-datasource-vpc-ipv6-addresses"
description: |- 
  查询专有网络（VPC）IPv6地址
---

# alibabacloudstack_vpc_ipv6_addresses

根据指定过滤条件列出当前凭证权限可以访问的专有网络（VPC）IPv6地址列表。

## 示例用法

```terraform
variable "name" {
  default = "tf-testacc-vpcipv6address-2169925"
}

data "alibabacloudstack_instances" "default" {
  name_regex = "no-deleteing-ipv6-address"
  status     = "Running"
}

data "alibabacloudstack_vpcs" "default" {
  name_regex = "no-deleteing-ipv6-address"
}

data "alibabacloudstack_vpc_ipv6_addresses" "default" {
  associated_instance_id = data.alibabacloudstack_instances.default.instances.0.id
  vswitch_id             = "your-vswitch-id"
  vpc_id                = data.alibabacloudstack_vpcs.default.ids[0]
  status                = "Available"

  output_file = "output.txt"
}

output "ipv6_address_1" {
  value = data.alibabacloudstack_vpc_ipv6_addresses.default.addresses[0].ipv6_address
}

output "ipv6_address_name_1" {
  value = data.alibabacloudstack_vpc_ipv6_addresses.default.addresses[0].ipv6_address_name
}
```

## 参数参考

以下参数是支持的：

* `associated_instance_id` - (可选，变更时重建) 分配给 IPv6 地址的实例 ID。
* `ids` - (可选，变更时重建) IPv6 地址 ID 列表。
* `status` - (可选，变更时重建) 资源的状态。有效值：`Available`、`Pending` 和 `Deleting`。
* `vswitch_id` - (可选，变更时重建) 分配给该 IPv6 地址的交换机 ID。
* `vpc_id` - (可选，变更时重建) 分配给该 IPv6 地址的 VPC ID。

## 属性参考

除了上述参数外，还导出以下属性：

* `addresses` - VPC IPv6 地址列表。每个元素包含以下属性：
  * `associated_instance_id` - 分配给 IPv6 地址的实例 ID。
  * `associated_instance_type` - 分配给 IPv6 地址的实例类型。
  * `create_time` - 资源的创建时间。
  * `id` - IPv6 地址的 ID。
  * `ipv6_address` - IPv6 地址。
  * `ipv6_address_id` - 资源主键属性字段。
  * `ipv6_address_name` - IPv6 地址的名称。名称必须是 2 到 128 个字符长度，可以包含字母、数字、下划线 (_) 和连字符 (-)。名称必须以字母开头，但不能以 `http://` 或 `https://` 开头。
  * `ipv6_gateway_id` - 分配给该 IPv6 地址的 IPv6 网关 ID。
  * `network_type` - IPv6 地址支持的通信类型。有效值：`Private` 或 `Public`。
    - `Private`：私有网络内的通信。
    - `Public`：公共网络上的通信。
  * `status` - 资源状态。有效值：`Available`、`Pending` 和 `Deleting`。
  * `vswitch_id` - 分配给该 IPv6 地址的交换机 ID。
  * `vpc_id` - 分配给该 IPv6 地址的 VPC ID。

* `names` - IPv6 地址名称列表。