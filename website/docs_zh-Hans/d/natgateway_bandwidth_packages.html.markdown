---
subcategory: "NAT网关"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_natgateway_bandwidth_packages"
sidebar_current: "docs-Alibabacloudstack-datasource-natgateway-bandwidth-packages"
description: |-
  提供阿里云账号下拥有的NAT网关带宽包列表。
---

# alibabacloudstack\_natgateway\_bandwidth\_packages

此数据源提供根据指定过滤条件列出的阿里云账号下的NAT网关带宽包资源列表。

## 示例用法

```hcl
data "alibabacloudstack_natgateway_bandwidth_packages" "example" {
  name_regex = "^my-BWP"
}

output "first_bwp_id" {
  value = data.alibabacloudstack_natgateway_bandwidth_packages.example.bandwidth_packages.0.id
}
```

## 参数说明

支持以下参数：

* `ids` - (可选) NAT网关带宽包ID列表。
* `name_regex` - (可选) 用于按带宽包名称过滤结果的正则表达式字符串。
* `description_regex` - (可选) 用于按带宽包描述过滤结果的正则表达式字符串。

## 属性说明

除了上述参数外，还导出以下属性：

* `ids` - NAT网关带宽包ID列表。
* `bandwidth_packages` - NAT网关带宽包列表。每个元素包含以下属性：
  * `id` - NAT网关带宽包ID。
  * `bandwidth` - 带宽包的带宽峰值。单位：Mbps。
  * `bandwidth_package_id` - 带宽包ID。
  * `business_status` - 带宽包的业务状态。取值：`Normal`（正常）、`FinancialLocked`（财务锁定）、`Unactivated`（未激活）。
  * `creation_time` - 带宽包的创建时间。
  * `name` - 带宽包的名称。
  * `description` - 带宽包的描述信息。
  * `instance_charge_type` - 带宽包的计费类型。取值：`PostPaid`（按量付费）、`PrePaid`（包年包月）。
  * `internet_charge_type` - 网络计费类型。取值：`PayByBandwidth`（按带宽计费）、`PayBy95`（按95峰值计费）。
  * `natgateway_id` - 与带宽包关联的NAT网关ID。
  * `ip_count` - 绑定到带宽包的EIP数量。
  * `isp` - 线路类型。取值：`BGP`（BGP多线）、`BGP_PRO`（BGP精品线）。
  * `public_ip_addresses` - 公网IP地址列表。每个元素包含：
    * `ip_address` - 公网IP地址。
    * `allocation_id` - EIP的分配ID。
  * `status` - 带宽包的状态。取值：`Available`（可用）。
