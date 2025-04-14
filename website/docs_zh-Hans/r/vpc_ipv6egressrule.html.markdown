---
subcategory: "VPC"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_vpc_ipv6egressrule"
sidebar_current: "docs-Alibabacloudstack-vpc-ipv6egressrule"
description: |- 
  编排VPC的IPv6 出口规则
---

# alibabacloudstack_vpc_ipv6_egress_rule
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_vpc_ipv6egressrule`

使用Provider配置的凭证在指定的资源集编排VPC的IPv6出口规则。

有关 VPC IPv6 出口规则及其使用方法的信息，请参阅 [什么是 IPv6 出口规则](https://www.alibabacloud.com/help/doc-detail/102200.htm)。

## 示例用法

### 基础用法

```terraform
variable "name" {
  default = "tf-testaccvpcipv6egressrule88807"
}

resource "alibabacloudstack_vpc" "default" {
  vpc_name    = "example_value"
  enable_ipv6 = true
}

resource "alibabacloudstack_vpc_ipv6_gateway" "example" {
  ipv6_gateway_name = "example_value"
  vpc_id            = alibabacloudstack_vpc.default.id
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
  ipv6_gateway_id       = alibabacloudstack_vpc_ipv6_gateway.example.id
  instance_id           = data.alibabacloudstack_vpc_ipv6_addresses.default.ids.0
  instance_type         = "Ipv6Address"
  description           = var.name
}
```

## 参数说明

支持以下参数：

* `description` - (可选，强制新建) 出口规则的描述。描述必须在 `2` 到 `256` 个字符之间，不能以 `http://` 或 `https://` 开头。
* `instance_id` - (必填，强制新建) 要应用出口规则的 IPv6 地址的 ID。
* `instance_type` - (可选，强制新建) 要应用出口规则的实例类型。有效值：`Ipv6Address`（默认值），表示一个 IPv6 地址。
* `ipv6_egress_rule_name` - (可选，强制新建) 出口规则的名称。名称必须在 `2` 到 `128` 个字符之间，可以包含字母、数字、下划线 (`_`) 和连字符 (`-`)。名称必须以字母开头，但不能以 `http://` 或 `https://` 开头。
* `ipv6_gateway_id` - (必填，强制新建) IPv6 网关的 ID。

## 属性说明

除了上述所有参数外，还导出了以下属性：

* `id` - Terraform 中 IPv6 出口规则的资源 ID。格式为 `<ipv6_gateway_id>:<ipv6_egress_rule_id>`。
* `status` - 资源的状态。有效值：`Available`（可用）、`Pending`（等待中）和 `Deleting`（删除中）。
* `instance_type` - 实例类型，表示需要设置出口规则的实例类型。取值为 `Ipv6Address`（默认值），表示一个 IPv6 地址。

### 超时设置

`timeouts` 块允许您为某些操作指定 [超时](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts)：

* `create` - (默认为 1 分钟)用于创建 IPv6 出口规则时。
* `delete` - (默认为 1 分钟)用于删除 IPv6 出口规则时。

## 导入

VPC IPv6 出口规则可以通过 id 导入，例如：

```bash
$ terraform import alibabacloudstack_vpc_ipv6_egress_rule.example <ipv6_gateway_id>:<ipv6_egress_rule_id>
```