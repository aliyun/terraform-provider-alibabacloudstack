---
subcategory: "DNS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_dns_recursor_acls"
sidebar_current: "docs-Alibabacloudstack-datasource-dns-recursor-acls"
description: |-
  查询DNS Recursor ACL策略列表
---

# alibabacloudstack_dns_recursor_acls

查询DNS Recursor ACL策略列表，用于管理跨云解析域名解析的访问控制策略。

## 示例用法

```hcl

variable "name" {
  default = "tfacc63456"
}

resource "alibabacloudstack_dns_line" "ipv4" {
  name         = "${var.name}ipv4"
  v4_addresses = ["192.168.0.1"]
}

resource "alibabacloudstack_dns_recursor_acl" "default" {
  name     = var.name
  remark   = var.name
  policy   = "ALLOW"
  line_ids = ["${alibabacloudstack_dns_line.ipv4.id}"]
}

data "alibabacloudstack_dns_recursor_acls" "default" {
  name_regex = alibabacloudstack_dns_recursor_acl.default.name
}

```

## 参数说明
以下参数用于过滤查询结果：

* `ids` (可选)：指定要查询的Recursor ACL ID列表。如果提供了此参数，则只返回ID在列表中的Recursor ACL。

* `name_regex` (可选)：用于通过正则表达式过滤Recursor ACL名称。只有名称匹配该正则表达式的Recursor ACL会被返回。

## 属性说明
以下属性被导出：

* `id` (字符串)：Recursor ACL的唯一标识符。

* `create_timestamp` (整数)：Recursor ACL的创建时间戳（秒）。

* `line_ids` (列表)：请求线路ID列表，表示该ACL策略应用的线路。

* `name` (字符串)：Recursor ACL的名称。

* `policy` (字符串)：递归策略，取值为"ALLOW"表示允许递归，"FORBID"表示禁止递归。

* `remark` (字符串)：Recursor ACL的备注信息。

* `update_timestamp` (整数)：Recursor ACL的最后更新时间戳（秒）。