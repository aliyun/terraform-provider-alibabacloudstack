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

## 参数参考

支持以下参数：

* `group_id` - (必填) 拥有规则的安全组的 ID。
* `nic_type` - (可选) 指网络类型。可以是 `internet` 或 `intranet`。默认值为 `internet`。
* `direction` - (可选) 授权方向。有效值为：`ingress` 或 `egress`。
* `ip_protocol` - (可选) IP 协议。有效值为：`tcp`、`udp`、`icmp`、`gre` 和 `all`。
* `policy` - (可选) 授权策略。可以是 `accept` 或 `drop`。默认值为 `accept`。
* `group_name` - (可选) 拥有规则的安全组的名称。
* `group_desc` - (可选) 拥有规则的安全组的描述。

## 属性参考

除了上述参数外，还导出以下属性：

* `group_name` - 拥有规则的安全组的名称。
* `group_desc` - 拥有规则的安全组的描述。
* `rules` - 安全组规则列表。每个元素包含以下属性：
  * `ip_protocol` - 协议。可以是 `tcp`、`udp`、`icmp`、`gre` 或 `all`。
  * `port_range` - 端口号范围。
  * `source_cidr_ip` - 入站授权的源 IP 地址段。
  * `source_group_owner_account` - 源安全组的阿里云账户。
  * `dest_cidr_ip` - 出站授权的目标 IP 地址段。
  * `dest_group_owner_account` - 目标安全组的阿里云账户。
  * `policy` - 授权策略。可以是 `accept` 或 `drop`。
  * `nic_type` - 网络类型，`internet` 或 `intranet`。
  * `priority` - 规则优先级。
  * `direction` - 授权方向，`ingress` 或 `egress`。
  * `dest_group_id` - 入站授权的目标安全组 id。
  * `source_group_id` - 入站授权的源安全组 ID。
  * `group_id` - 拥有规则的安全组的 ID。
  * `output_file` - 保存数据源结果的文件名(运行 `terraform plan` 后)。