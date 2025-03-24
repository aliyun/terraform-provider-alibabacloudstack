---
subcategory: "ECS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_security_group_rule"
sidebar_current: "docs-alibabacloudstack-resource-security-group-rule"
description: |-
  编排安全组规则
---

# alibabacloudstack_security_group_rule

使用Provider配置的凭证在指定的资源集编排安全组规则。
表示单个 `ingress` 或 `egress` 组规则，可以添加到外部安全组。

-> **注意:** 当安全组类型为 `vpc` 或指定 `source_security_group_id` 时，`nic_type` 应该设置为 `intranet`。在这种情况下，它不区分内网和外网，规则在两者上都有效。


## 示例用法

### 基础用法

```
resource "alibabacloudstack_vpc" "vpc" {
  cidr_block = "10.1.0.0/21"
}

resource "alibabacloudstack_security_group" "group" {
  vpc_id = "${alibabacloudstack_vpc.vpc.id}"
}

resource "alibabacloudstack_security_group_rule" "allow_all_tcp" {
  type              = "ingress"
  ip_protocol       = "tcp"
  nic_type          = "intranet"
  policy            = "accept"
  port_range        = "1/65535"
  priority          = 1
  security_group_id = "${alibabacloudstack_security_group.default.id}"
  cidr_ip           = "0.0.0.0/0"
}
```

## 参数参考

支持以下参数：

* `type` - (必填，变更时重建) 正在创建的规则类型。有效的选项是 `ingress`(入站)或 `egress`(出站)。
* `ip_protocol` - (必填，变更时重建) 协议。可以是 `tcp`、`udp`、`icmp`、`gre` 或 `all`。
* `port_range` - (变更时重建) 与 IP 协议相关的端口范围。默认为 "-1/-1"。当协议为 tcp 或 udp 时，每侧端口号范围从 1 到 65535，'-1/-1' 将无效。
  例如，`1/200` 表示端口号范围为 1-200。其他协议的 'port_range' 只能为 "-1/-1"，其他值将无效。
* `security_group_id` - (必填，变更时重建) 要应用此规则的安全组。
* `nic_type` - (可选，变更时重建) 网络类型，可以是 `internet` 或 `intranet`，默认值为 `internet`。
* `policy` - (可选，变更时重建) 授权策略，可以是 `accept` 或 `drop`，默认值为 `accept`。
* `priority` - (可选，变更时重建) 授权策略优先级，参数值范围：`1-100`，默认值：1。
* `cidr_ip` - (可选，变更时重建) 目标 IP 地址范围。默认值为 0.0.0.0/0(表示没有限制)。其他支持的格式包括 10.159.6.18/12。仅支持 IPv4。
* `source_security_group_id` - (可选，变更时重建) 同一区域内的目标安全组 ID。如果设置了此字段，则 `nic_type` 只能选择 `intranet`。
* `source_group_owner_account` - (可选，变更时重建) 跨账户授权时目标安全组所属的阿里云用户账号 ID。如果已设置 `cidr_ip` 参数，则此参数无效。
* `ipv6_cidr_ip` - (可选，强制新，自 v1.174.0 起可用)需要访问的源 IPv6 CIDR 地址块。支持 CIDR 格式和 IPv6 格式的 IP 地址范围。注意：此参数不能与 `cidr_ip` 参数同时设置。
* `description` - (可选) 安全组规则的描述。描述长度可以为 1 到 512 个字符，默认为 null。
* `port_range` - (必填，变更时重建) 指定与 IP 协议相关的端口范围。它是定义 TCP/UDP 协议特定端口或范围所必填的。

-> **注意:** 必须设置 `source_security_group_id` 或 `cidr_ip`。

## 属性参考

导出以下属性：

* `id` - 安全组规则的 ID
* `type` - 规则类型，`ingress` 或 `egress`
* `port_range` - 端口范围
* `ip_protocol` - 安全组规则的协议
* `nic_type` - 指示网络类型，可以是 `internet` 或 `intranet`。此属性根据提供的配置计算得出。