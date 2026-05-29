---
subcategory: "裸金属算力平台 BMCP"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_bcmp_security_group_rule"
sidebar_current: "docs-Alibabacloudstack-resource-bcmp-security-group-rule"
description: |- 
  提供 BCMP 安全组规则资源。
---

# alibabacloudstack_bcmp_security_group_rule

提供 BCMP 安全组规则资源。

> **注意:** 此资源也可以使用以下别名引用：
> - `alibabacloudstack_bcmp_security_group_rules`（已弃用）

## 示例用法

基本用法

```hcl
resource "alibabacloudstack_vpc" "default" {
  name       = "tf-test-vpc"
  cidr_block = "172.16.0.0/12"
}

resource "alibabacloudstack_bcmp_security_group" "default" {
  vpc_id = alibabacloudstack_vpc.default.id
  name   = "tf-test-sg"
}

resource "alibabacloudstack_bcmp_security_group_rule" "default" {
  type              = "ingress"
  ip_protocol       = "tcp"
  policy            = "accept"
  port_range        = "22/22"
  priority          = 1
  security_group_id = alibabacloudstack_bcmp_security_group.default.id
  cidr_ip           = "0.0.0.0/0"
  description       = "允许 SSH 访问"
}
```

使用源安全组

```hcl
resource "alibabacloudstack_bcmp_security_group" "new" {
  vpc_id = alibabacloudstack_vpc.default.id
  name   = "tf-test-sg-new"
}

resource "alibabacloudstack_bcmp_security_group_rule" "with_source_sg" {
  type                     = "ingress"
  ip_protocol              = "tcp"
  policy                   = "drop"
  port_range               = "22/22"
  priority                 = 100
  security_group_id        = alibabacloudstack_bcmp_security_group.default.id
  source_security_group_id = alibabacloudstack_bcmp_security_group.new.id
  description              = "阻止来自其他安全组的 SSH"
}
```

## 参数说明

支持以下参数：

* `type` - (必填, ForceNew) 规则类型，`ingress`（入方向）或 `egress`（出方向）。修改此参数会强制重新创建资源。
* `ip_protocol` - (必填) IP 协议类型。有效值为 `tcp`、`udp`、`icmp`、`gre` 和 `all`。
* `policy` - (选填) 规则的策略。有效值为 `accept` 和 `drop`。默认为 `accept`。
* `port_range` - (必填) 端口范围。对于 `tcp` 和 `udp`，格式为 `start/end`（例如 `22/22`）。对于 `icmp`、`gre` 和 `all`，使用 `-1/-1`。
* `priority` - (选填) 规则的优先级。有效范围：1-100。默认为 1。
* `security_group_id` - (必填) 安全组的 ID。
* `cidr_ip` - (选填) CIDR IP 地址。与 `source_security_group_id` 冲突。必须指定 `cidr_ip` 或 `source_security_group_id` 中的一个。
* `source_security_group_id` - (选填) 源安全组的 ID。与 `cidr_ip` 冲突。必须指定 `cidr_ip` 或 `source_security_group_id` 中的一个。
* `description` - (选填) 规则的描述。

> **注意:** 必须指定 `cidr_ip` 或 `source_security_group_id` 中的一个。

## 属性说明

除了上述所有参数外，还导出了以下属性：

* `id` - 安全组规则的 ID。格式为 `<security_group_id>:<sgr_id>:<type>`。
* `type` - 规则的类型。
* `ip_protocol` - IP 协议类型。
* `policy` - 规则的策略。
* `port_range` - 端口范围。
* `priority` - 规则的优先级。
* `security_group_id` - 安全组的 ID。
* `cidr_ip` - CIDR IP 地址。
* `source_security_group_id` - 源安全组的 ID。
* `description` - 规则的描述。

## Import

BCMP 安全组规则可以使用 `<security_group_id>:<sgr_id>:<type>` 导入，例如：

```
$ terraform import alibabacloudstack_bcmp_security_group_rule.example sg-12345678:sgr-abcdef:ingress
```
