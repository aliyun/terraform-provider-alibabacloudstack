---
subcategory: "VPC"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_vpc_ipv6_egress_rules"
sidebar_current: "docs-Alibabacloudstack-datasource-vpc-ipv6-egress-rules"
description: |- 
  查询专有网络（VPC）IPv6地址出口规则
---

# alibabacloudstack_vpc_ipv6_egress_rules
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_vpc_ipv6_egressrules`

根据指定过滤条件列出当前凭证权限可以访问的专有网络（VPC）IPv6地址出口规则列表。

## 示例用法

以下示例展示了如何使用 `alibabacloudstack_vpc_ipv6_egress_rules` 数据源来查询和筛选 IPv6 出口规则。

### 基础用法

```terraform
data "alibabacloudstack_vpc_ipv6_egress_rules" "default" {
  ipv6_gateway_id = "example_value"
  ids             = ["example_value-1", "example_value-2"]
}

output "vpc_ipv6_egress_rule_id_1" {
  value = data.alibabacloudstack_vpc_ipv6_egress_rules.default.rules.0.id
}

data "alibabacloudstack_vpc_ipv6_egress_rules" "nameRegex" {
  ipv6_gateway_id = "example_value"
  name_regex      = "^my-Ipv6EgressRule"
}

output "vpc_ipv6_egress_rule_id_2" {
  value = data.alibabacloudstack_vpc_ipv6_egress_rules.nameRegex.rules.0.id
}

data "alibabacloudstack_vpc_ipv6_egress_rules" "status" {
  ipv6_gateway_id = "example_value"
  status          = "Available"
}

output "vpc_ipv6_egress_rule_id_3" {
  value = data.alibabacloudstack_vpc_ipv6_egress_rules.status.rules.0.id
}

data "alibabacloudstack_vpc_ipv6_egress_rules" "ipv6EgressRuleName" {
  ipv6_gateway_id       = "example_value"
  ipv6_egress_rule_name = "example_value"
}

output "vpc_ipv6_egress_rule_id_4" {
  value = data.alibabacloudstack_vpc_ipv6_egress_rules.ipv6EgressRuleName.rules.0.id
}
```

### 结合其他资源使用

```terraform
variable "name" {
  default = "tf-testacc-vpcipv6egressrule-4354090"
}

data "alibabacloudstack_instances" "default" {
  name_regex = "no-deleteing-ipv6-address"
  status     = "Running"
}

data "alibabacloudstack_vpc_ipv6_addresses" "default" {
  associated_instance_id = data.alibabacloudstack_instances.default.instances.0.id
  status                 = "Available"
}

resource "alibabacloudstack_vpc_ipv6_egress_rule" "default" {
  ipv6_egress_rule_name = var.name
  ipv6_gateway_id       = data.alibabacloudstack_vpc_ipv6_addresses.default.addresses.0.ipv6_gateway_id
  instance_id           = data.alibabacloudstack_vpc_ipv6_addresses.default.ids.0
  instance_type         = "Ipv6Address"
  description           = var.name
}

data "alibabacloudstack_vpc_ipv6_egress_rules" "default" {
  ipv6_gateway_id = alibabacloudstack_vpc_ipv6_egress_rule.default.ipv6_gateway_id
  name_regex      = alibabacloudstack_vpc_ipv6_egress_rule.default.ipv6_egress_rule_name
}
```

## 参数参考

以下参数是支持的：

* `instance_id` - (选填, 变更时重建) - 设置仅主动出规则的IPv6地址的ID。
* `ids` - (选填, 变更时重建) - IPv6出口规则ID列表，用于精确匹配特定规则。
* `name_regex` - (选填, 变更时重建) - 正则表达式字符串，用于通过名称筛选IPv6出口规则。
* `ipv6_egress_rule_name` - (选填, 变更时重建) - 出口规则的名称。名称长度必须在2到128个字符之间，可以包含字母、数字、下划线(_)和连字符(-)，并且必须以字母开头。
* `ipv6_gateway_id` - (必填, 变更时重建) - IPv6网关的ID。
* `status` - (选填, 变更时重建) - 资源的状态。有效值为：`Available`、`Deleting`、`Pending`。

## 属性参考

除了上述参数外，还导出以下属性：

* `names` - IPv6出口规则名称列表。
* `rules` - VPC IPv6出口规则列表。每个元素包含以下属性：
  * `description` - 出口规则的描述。描述长度必须在2到256个字符之间，不能以`http://`或`https://`开头。
  * `id` - IPv6出口规则的ID。格式为 `<ipv6_gateway_id>:<ipv6_egress_rule_id>`。
  * `instance_id` - 应用了出口规则的实例的ID。
  * `instance_type` - 需要设置出口规则的实例类型。有效值为：`Ipv6Address`(默认值)：一个IPv6地址。
  * `ipv6_egress_rule_id` - IPv6出口规则的ID。
  * `ipv6_egress_rule_name` - 出口规则的名称。名称长度必须在2到128个字符之间，可以包含字母、数字、下划线(_)和连字符(-)，并且必须以字母开头。
  * `status` - 资源的状态。有效值为：`Available`、`Pending`、`Deleting`。
  * `ipv6_gateway_id` - IPv6网关的ID。