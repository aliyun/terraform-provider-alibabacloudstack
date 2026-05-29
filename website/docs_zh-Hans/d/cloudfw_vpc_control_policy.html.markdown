---
subcategory: "云防火墙"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cloudfw_vpc_control_policy"
sidebar_current: "docs-Alibabacloudstack-datasource-cloudfw-vpc-control-policy"
description: |-
  查询云防火墙VPC控制策略
---

# alibabacloudstack_cloudfw_vpc_control_policy

查询云防火墙VPC控制策略，用于获取已配置的VPC防火墙访问控制规则信息。

## 示例用法

```hcl

variable "name" {
  default = "tf-testacc-vpccontrolpolicies-12994"
}

resource "alibabacloudstack_cloudfw_vpc_control_policy" "default" {
  destination      = "0.0.0.0/16"
  description      = var.name
  application_name = "ANY"
  source_type      = "net"
  dest_port        = "33/33"
  acl_action       = "log"
  destination_type = "net"

  source         = "0.0.0.0/16"
  dest_port_type = "port"
  proto          = "UDP"
  application_id = "0"
  release        = true
}

data "alibabacloudstack_cloudfw_vpc_control_policies" "default" {
  ids = ["${alibabacloudstack_cloudfw_vpc_control_policy.default.id}"]
}

```

## 参数说明
以下参数支持过滤查询结果：

* `destination` (字符串, 可选)：目标地址或地址组，用于过滤目标地址匹配的策略。

* `ids` (列表, 可选)：策略ID列表，用于过滤指定ID的策略。如果不指定，将返回所有策略。

* `name_regex` (字符串, 可选)：策略描述的正则表达式，用于过滤描述匹配的策略。

* `proto` (字符串, 可选)：协议类型，例如"TCP"、"UDP"，用于过滤指定协议的策略。

* `source` (字符串, 可选)：源地址或地址组，用于过滤源地址匹配的策略。

* `acl_action` (字符串, 可选)：访问控制规则的动作，例如"accept"、"drop"、"log"，用于过滤指定动作的策略。

## 属性说明
以下属性被导出：

* `id` (字符串)：数据源的唯一标识符，基于匹配到的策略ID哈希生成。

* `ids` (列表)：匹配到的策略ID列表。

* `names` (列表)：匹配到的策略描述列表。

* `policies` (列表)：匹配到的策略列表，每个策略包含以下属性：
  * `acl_uuid` (字符串)：策略的唯一ID。
  * `application_id` (字符串)：关联应用的ID。
  * `application_name` (字符串)：关联应用的名称。
  * `description` (字符串)：策略的描述信息。
  * `dest_port` (字符串)：目标端口或端口范围。
  * `dest_port_group` (字符串)：目标端口组名称。
  * `dest_port_group_ports` (列表)：目标端口组包含的端口列表。
  * `dest_port_type` (字符串)：目标端口类型，例如"port"。
  * `destination` (字符串)：目标地址或地址组。
  * `destination_group_cidrs` (列表)：当目标类型为组时，包含的目标CIDR列表。
  * `destination_type` (字符串)：目标类型，例如"net"。
  * `direction` (字符串)：流量方向，例如"inout"。
  * `hit_times` (整数)：策略匹配次数。
  * `order` (整数)：策略在规则列表中的顺序。
  * `proto` (字符串)：使用的协议，例如"TCP"、"UDP"。
  * `release` (字符串)：策略是否已发布（启用）。
  * `source` (字符串)：源地址或地址组。
  * `source_group_cidrs` (列表)：当源类型为组时，包含的源CIDR列表。
  * `source_type` (字符串)：源类型，例如"net"。