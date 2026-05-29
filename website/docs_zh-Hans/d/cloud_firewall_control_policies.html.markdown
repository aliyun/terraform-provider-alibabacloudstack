---
subcategory: "云防火墙"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cloud_firewall_control_policies"
sidebar_current: "docs-Alibabacloudstack-datasource-cloud-firewall-control-policies"
description: |- 
  查询防火墙控制策略
---

# Data Source: alibabacloudstack_cloud_firewall_control_policies

根据指定过滤条件列出当前凭证权限可以访问的防火墙控制策略列表。

-> **NOTE:** 自 v1.213.0+ 版本可用

## 示例用法

```terraform
data "alibabacloudstack_cloud_firewall_control_policies" "example" {
  direction   = "in"
  acl_action  = "accept"
  source      = "192.168.0.0/16"
  destination = "10.0.0.0/8"
  proto       = "TCP"
}
```

## 参数说明

以下参数是支持的：

* `direction` - (必填) 安全访问控制策略的流量方向。有效值：
  * **in**：外对内流量访问控制
  * **out**：内对外流量访问控制

* `acl_action` - (选填) 安全访问控制策略中设置的流量通过云防火墙的方式。有效值：
  * **accept**：放行
  * **drop**：拒绝
  * **log**：观察

* `acl_uuid` - (选填) 安全访问控制策略的唯一标识ID。

* `description` - (选填) 安全访问控制策略的描述信息。

* `destination` - (选填) 安全访问控制策略中的目的地址。

* `proto` - (选填) 安全访问控制策略中流量访问的安全协议类型。有效值：`TCP`、`UDP`、`ANY`、`ICMP`。

* `source` - (选填) 安全访问控制策略中的源地址。

* `source_ip` - (选填) 请求的源IP地址。

* `output_file` - (选填, 已废弃) 该字段已废弃，计划在 3.19.0 版本中移除。如需将内容写入文件，请使用 `local_file` Provider。

## 属性说明

除了上述参数外，还导出以下属性：

* `ids` - 控制策略ID列表。每个ID格式为 `<acl_uuid>:<direction>`。

* `policies` - 云防火墙控制策略列表。每个元素包含以下属性：
  * `id` - 控制策略的ID。格式为 `<acl_uuid>:<direction>`。
  * `acl_uuid` - 安全访问控制策略的唯一标识ID。
  * `acl_action` - 安全访问控制策略中设置的流量通过云防火墙的方式。
  * `application_id` - 安全访问控制策略中设置的应用ID。
  * `application_name` - 安全访问控制策略中的应用名称。
  * `description` - 安全访问控制策略的描述信息。
  * `dest_port` - 安全访问控制策略中的目的端口。
  * `dest_port_group` - 安全访问控制策略中的目的端口地址簿名称。
  * `dest_port_group_ports` - 目的端口地址簿中包含的端口列表。
  * `dest_port_type` - 安全访问控制策略中的目的端口类型。
  * `destination` - 安全访问控制策略中的目的地址。
  * `destination_group_cidrs` - 安全访问控制策略中的目的地址簿中的网段列表。
  * `destination_group_type` - 安全访问控制策略中的目的地址簿类型。
  * `destination_type` - 安全访问控制策略中的目的地址类型。
  * `direction` - 安全访问控制策略的流量方向。
  * `hit_times` - 安全访问控制策略命中次数统计。
  * `order` - 安全访问控制策略生效的优先级。
  * `proto` - 安全访问控制策略中流量访问的安全协议类型。
  * `release` - 安全访问控制策略是否生效。
  * `source` - 安全访问控制策略中的源地址。
  * `source_group_cidrs` - 安全访问控制策略中的源地址簿中的网段列表。
  * `source_group_type` - 安全访问控制策略中的源地址簿类型。
  * `source_type` - 安全访问控制策略中的源地址类型。
