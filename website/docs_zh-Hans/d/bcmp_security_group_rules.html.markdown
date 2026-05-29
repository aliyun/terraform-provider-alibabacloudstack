---
subcategory: "裸金属计算平台"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_bcmp_security_group_rules"
sidebar_current: "docs-Alibabacloudstack-datasource-bcmp-security-group-rules"
description: |- 
  提供 BCMP 安全组规则列表。
---

# alibabacloudstack_bcmp_security_group_rules

此数据源根据指定的过滤器提供 AlibabaCloudStack 账户中的 BCMP 安全组规则列表。

## 示例用法

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
  description       = "test"
}

# 检索安全组的所有规则
data "alibabacloudstack_bcmp_security_group_rules" "default" {
  security_group_id = alibabacloudstack_bcmp_security_group.default.id
}

output "rules" {
  value = data.alibabacloudstack_bcmp_security_group_rules.default.rules
}
```

## 参数说明

支持以下参数：

* `security_group_id` - (必填) 安全组的 ID。
* `type` - (选填) 用于过滤的规则类型。有效值为 `ingress` 和 `egress`。
* `ip_protocol` - (选填) 用于过滤的 IP 协议类型。有效值为 `tcp`、`udp`、`icmp`、`gre` 和 `all`。
* `policy` - (选填) 用于过滤的策略。有效值为 `accept` 和 `drop`。

## 属性说明

除了上述列出的参数外，还导出以下属性：

* `rules` - 安全组规则列表。每个元素包含以下属性：
  * `type` - 规则的类型。
  * `ip_protocol` - IP 协议类型。
  * `port_range` - 端口范围。
  * `cidr_ip` - CIDR IP 地址。
  * `source_security_group_id` - 源安全组的 ID。
  * `policy` - 规则的策略。
  * `priority` - 规则的优先级。
  * `description` - 规则的描述。
