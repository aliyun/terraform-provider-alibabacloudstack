---
subcategory: "ECS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_security_group_rules"
sidebar_current: "docs-alibabacloudstack-datasource-security-group-rules"
description: |-
    查询安全组规则
---

# alibabacloudstack_security_group_rules

根据指定过滤条件列出当前凭证权限可以访问的安全组规则列表
每个集合项表示一个单独的 `ingress` 或 `egress` 权限规则。
可以通过变量或通过其他数据源 `alibabacloudstack_security_groups` 的结果提供安全组的 ID。

## 示例用法

以下示例展示了如何获取安全组规则的详细信息，以及如何在实例启动时传递其数据。

```
# 从变量中获取安全组 id
variable "security_group_id" {}

# 或者从 alibabacloudstack_security_groups 数据源中获取它。
# 请注意，数据源参数必须足够过滤出一个安全组。
data "alibabacloudstack_security_groups" "groups_ds" {
  name_regex = "api"
}

# 按组筛选安全组规则
data "alibabacloudstack_security_group_rules" "ingress_rules_ds" {
  group_id    = "${data.alibabacloudstack_security_groups.groups_ds.groups.0.id}" # or ${var.security_group_id}
  nic_type    = "internet"
  direction   = "ingress"
  ip_protocol = "TCP"
}

output "security_group_rules" {
  value = data.alibabacloudstack_security_group_rules.ingress_rules_ds
}

```

## 参数说明

支持以下参数：

* `group_id` - (必填) 拥有规则的安全组的 ID。
* `nic_type` - (可选) 指网络类型。可以是 `internet` 或 `intranet`。默认值为 `internet`。
* `direction` - (可选) 授权方向。有效值为：`ingress` 或 `egress`。
* `ip_protocol` - (可选) IP 协议。有效值为：`tcp`、`udp`、`icmp`、`gre` 和 `all`。
* `policy` - (可选) 授权策略。可以是 `accept` 或 `drop`。默认值为 `accept`。
* `group_name` - (可选) 拥有规则的安全组的名称。
* `group_desc` - (可选) 拥有规则的安全组的描述。

## 属性说明

除了上述参数外，还导出以下属性：

* `group_name` - 拥有规则的安全组的名称。
* `group_desc` - 拥有规则的安全组的描述。
* `rules` - 安全组规则列表。每个元素包含以下属性：
  * `ip_protocol` - 规则所使用的协议类型，可为 `tcp`、`udp`、`icmp`、`gre` 或 `all`。
  * `port_range` - 端口范围，格式为“起始端口/结束端口”。
  * `source_cidr_ip` - 入站规则的源 IP 地址段。
  * `source_group_owner_account` - 入站规则中源安全组所属的阿里云账户。
  * `dest_cidr_ip` - 出站规则的目标 IP 地址段。
  * `dest_group_owner_account` - 出站规则中目标安全组所属的阿里云账户。
  * `policy` - 授权策略，可为 `accept`（允许）或 `drop`（拒绝）。
  * `nic_type` - 网络类型，可为 `internet`（公网）或 `intranet`（内网）。
  * `priority` - 规则优先级，数值越小优先级越高。
  * `direction` - 授权方向，可为 `ingress`（入站）或 `egress`（出站）。
  * `dest_group_id` - 出站规则中的目标安全组 ID。
  * `source_group_id` - 入站规则中的源安全组 ID。
  * `group_id` - 拥有该规则的安全组 ID。