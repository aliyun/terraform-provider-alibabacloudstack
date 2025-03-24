---
subcategory: "VPC"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_vpc_ipv6_gateways"
sidebar_current: "docs-Alibabacloudstack-datasource-vpc-ipv6_gateways"
description: |- 
  查询专有网络（VPC）IPv6网关
---

# alibabacloudstack_vpc_ipv6_gateways

根据指定过滤条件列出当前凭证权限可以访问的专有网络（VPC）IPv6网关列表。

## 示例用法

### 基础用法：

```hcl
variable "name" {
  default = "tf-testacc-vpcipv6gateway-4944497"
}

resource "alibabacloudstack_vpc" "default" {
  vpc_name    = var.name
  enable_ipv6 = "true"
}

resource "alibabacloudstack_vpc_ipv6_gateway" "default" {
  vpc_id            = alibabacloudstack_vpc.default.id
  ipv6_gateway_name = var.name
  description       = var.name
}

data "alibabacloudstack_vpc_ipv6_gateways" "ids" {
  ids = [alibabacloudstack_vpc_ipv6_gateway.default.ipv6_gateway_id]
}

output "vpc_ipv6_gateway_id_1" {
  value = data.alibabacloudstack_vpc_ipv6_gateways.ids.gateways.0.id
}

data "alibabacloudstack_vpc_ipv6_gateways" "nameRegex" {
  name_regex = "${alibabacloudstack_vpc_ipv6_gateway.default.ipv6_gateway_name}"
}

output "vpc_ipv6_gateway_id_2" {
  value = data.alibabacloudstack_vpc_ipv6_gateways.nameRegex.gateways.0.id
}

data "alibabacloudstack_vpc_ipv6_gateways" "vpcId" {
  vpc_id = alibabacloudstack_vpc.default.id
}

output "vpc_ipv6_gateway_id_3" {
  value = data.alibabacloudstack_vpc_ipv6_gateways.vpcId.gateways.0.id
}

data "alibabacloudstack_vpc_ipv6_gateways" "status" {
  status = "Available"
}

output "vpc_ipv6_gateway_id_4" {
  value = data.alibabacloudstack_vpc_ipv6_gateways.status.gateways.0.id
}
```

## 参数参考

以下参数是支持的：

* `ids` - (可选，变更时重建) IPv6网关ID列表。用于通过ID筛选IPv6网关。
* `name_regex` - (可选，变更时重建) 用于按IPv6网关名称筛选结果的正则表达式字符串。
* `ipv6_gateway_name` - (可选，变更时重建) IPv6网关的名称。名称必须是2到128个字符长度，可以包含字母、数字、下划线(_)和连字符(-)。名称必须以字母开头，但不能以`http://`或`https://`开头。
* `status` - (可选，变更时重建) 资源的状态。有效值：`Available`(可用)、`Pending`(等待中)和`Deleting`(删除中)。
* `vpc_id` - (可选，变更时重建) IPv6网关所属虚拟私有云(VPC)的ID。用于筛选特定VPC下的IPv6网关。

## 属性参考

除了上述参数外，还导出以下属性：

* `gateways` - VPC IPv6网关列表。每个元素包含以下属性：
  * `business_status` - IPv6网关的状态。有效值：`Normal`(正常)、`FinancialLocked`(因逾期付款被锁定)和`SecurityLocked`(因安全原因被锁定)。
  * `create_time` - 资源的创建时间。
  * `description` - IPv6网关的描述。描述必须是2到256个字符长度，且不能以`http://`或`https://`开头。
  * `expired_time` - IPv6网关的过期时间。
  * `instance_charge_type` - IPv6网关的计费类型。有效值：`PayAsYouGo`(按量付费)。
  * `id` - IPv6网关的唯一标识符(ID)。
  * `ipv6_gateway_id` - IPv6网关的主要键属性字段。
  * `ipv6_gateway_name` - IPv6网关的名称。名称必须是2到128个字符长度，可以包含字母、数字、下划线(_)和连字符(-)。名称必须以字母开头，但不能以`http://`或`https://`开头。
  * `spec` - IPv6网关的规格。此参数不再使用，因为IPv6网关不分规格。
  * `status` - IPv6网关的状态。有效值：`Available`(可用)、`Pending`(等待中)和`Deleting`(删除中)。
  * `vpc_id` - IPv6网关所属虚拟私有云(VPC)的ID。