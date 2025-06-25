---
subcategory: "VPC"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_vpc_ipv6_isps"
sidebar_current: "docs-Alibabacloudstack-datasource-vpc-ipv6-isps"
description: |- 
  Provides a list of vpc ipv6 isps in an alibabacloudstack account.
---


数据源：alibabacloudstack_vpc_ipv6_isps 数据源用于获取阿里云 VPC 中可用的 IPv6 ISP（互联网服务提供商）列表信息。

## 使用示例
```
hcl
data "alibabacloudstack_vpc_ipv6_isps" "example" {
  service_provider = "BGP"
  lock_status      = "unlocked"
}
```
## 参数说明
支持以下参数：

* `ids` - （可选）指定一组 IPv6 ISP 网段池 ID，用于过滤结果。
* `lock_status` - （可选）根据锁定状态进行过滤（locked 或 unlocked）。
* `service_provider` - （可选）根据服务商名称进行过滤，例如 China Mobile、China Telecom 等。

## 导出属性
以下属性被导出：

* `ipv6_isps` - IPv6 ISP 列表。每项包含：
  * `id` - IPv6 ISP 网段池的 ID（即 PoolId）。
  * `service_provider` - 服务提供商名称。
  * `zone_id` - 可用区 ID。
  * `type` - IPv6 ISP 类型。
  * `cidr_block` - 分配给该 ISP 的 CIDR 网段段。
  * `available_count` - 可用的 IPv6 网段数量。
  * `in_use_count` - 使用中 IPv6 网段数量。
  * `lock_status` - 当前网段池的锁定状态。
  * `need_declare` - 是否需要声明此 ISP。
  * `pool_id` - IPv6 网段池的唯一标识。
  * `ula` - 与该 ISP 关联的唯一本地网段（ULA）。
