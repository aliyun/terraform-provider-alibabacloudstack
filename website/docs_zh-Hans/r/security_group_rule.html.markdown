---
subcategory: "云服务器 ECS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_security_group_rule"
sidebar_current: "docs-alibabacloudstack-resource-security-group-rule"
description: |-
  编排安全组规则
---

# alibabacloudstack_security_group_rule

使用 Provider 配置的凭证在指定的资源集编排安全组规则。
表示单个 `ingress` 或 `egress` 组规则，可以添加到外部安全组。

-> **注意:** 当安全组类型为 `vpc` 或指定 `source_security_group_id` 时，`nic_type` 应该设置为 `intranet`。在这种情况下，它不区分内网和外网，规则在两者上都有效。

> **注意:** 此资源也可以使用以下别名引用：
> - `alibabacloudstack_ecs_securitygrouprule`


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

## 参数说明

支持以下参数：

* `type` - (必填，变更时重建) 正在创建的规则类型。有效的选项是 `ingress`（入站）或 `egress`（出站）。
* `ip_protocol` - (必填，变更时重建) 协议类型。可以是 `tcp`、`udp`、`icmp`、`gre` 或 `all`。
* `port_range` - (必填，变更时重建) 指定与 IP 协议相关的端口范围。对于 TCP/UDP 协议，必须定义具体的端口或范围。当协议为 tcp 或 udp 时，每侧端口号范围从 1 到 65535，'-1/-1' 将无效。例如，`1/200` 表示端口号范围为 1-200。其他协议的 'port_range' 只能为 "-1/-1"，其他值将无效。
* `security_group_id` - (必填，变更时重建) 要应用此规则的安全组 ID。
* `cidr_ip` - (可选，变更时重建) 目标 IPv4 CIDR 地址段。与 `ipv6_cidr_ip` 和 `source_security_group_id` 冲突。
* `ipv6_cidr_ip` - (可选，变更时重建) 目标 IPv6 CIDR 地址段。与 `cidr_ip` 和 `source_security_group_id` 冲突。
* `source_security_group_id` - (可选，变更时重建) 同一区域内的目标安全组 ID。与 `cidr_ip` 和 `ipv6_cidr_ip` 冲突。如果设置了此字段，则 `nic_type` 必须设置为 `intranet`。
* `nic_type` - (可选，变更时重建) 网络类型，可以是 `internet` 或 `intranet`，默认值为 `intranet`。
* `policy` - (可选，变更时重建) 授权策略，可以是 `accept` 或 `drop`，默认值为 `accept`。
* `priority` - (可选，变更时重建) 授权策略优先级，参数值范围：`1-100`，默认值：1。较低的数值表示较高的优先级。
* `description` - (可选) 安全组规则的描述。描述长度可以为 1 到 512 个字符。
* `source_group_owner_account` - (可选，已弃用) 跨账户授权时目标安全组所属的阿里云用户账号 ID。该参数在专有云中无效，计划在 3.19.0 版本中移除。

-> **注意:** `source_security_group_id`、`cidr_ip`、`ipv6_cidr_ip` 中必须设置其中一项。

## 属性说明

导出以下属性：

* `id` - 安全组规则的唯一标识符，格式为 `<security_group_id>:<type>:<ip_protocol>:<port_range>:<nic_type>:<cidr_ip>:<policy>:<priority>`。
* `type` - 规则类型，`ingress`（入站）或 `egress`（出站）。
* `port_range` - 规则中指定的端口范围。
* `ip_protocol` - 安全组规则中使用的协议类型。
* `nic_type` - 指示网络类型，可以是 `internet` 或 `intranet`。
* `policy` - 授权策略，`accept` 或 `drop`。
* `priority` - 规则优先级。
* `description` - 安全组规则的描述。

## Import

安全组规则可以使用组合 ID 进行导入，例如：

```
$ terraform import alibabacloudstack_security_group_rule.example sg-12345678:ingress:tcp:22/22:intranet:10.0.0.0/8:accept:1
```