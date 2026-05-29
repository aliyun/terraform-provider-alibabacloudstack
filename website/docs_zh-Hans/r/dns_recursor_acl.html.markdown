---
subcategory: "Cloud DNS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_dns_recursor_acl"
sidebar_current: "docs-Alibabacloudstack-dns-recursor_acl"
description: |-
  跨云解析域名解析
---

# alibabacloudstack_dns_recursor_acl

使用Provider配置的凭证在指定的资源集配置DNS递归ACL策略。

## 示例用法

### 基础用法

```hcl

variable "name" {
  default = "tfacc67879"
}

resource "alibabacloudstack_dns_line" "ipv4" {
  name         = "${var.name}ipv4"
  v4_addresses = ["192.168.0.1"]
}

resource "alibabacloudstack_dns_line" "ipv6" {
  name         = "${var.name}ipv6"
  v6_addresses = ["2020:148:2:28::"]
}

resource "alibabacloudstack_dns_recursor_acl" "default" {
  line_ids = [
    "${alibabacloudstack_dns_line.ipv4.id}",
    "${alibabacloudstack_dns_line.ipv6.id}"
  ]
  name   = var.name
  remark = var.name
  policy = "ALLOW"
}
```

## 参数说明

支持以下参数：

* `line_ids` - (必填) 请求线路ID列表。至少需要指定1个线路ID。
* `name` - (必填) ACL策略名称。
* `policy` - (必填) 允许递归策略。取值：`ALLOW`（允许递归）或`FORBID`（不允许递归）。
* `remark` - (可选) 备注信息。

## 属性说明

以下属性会从API响应中导出：

* `id` - ACL策略的ID。

## 导入

DNS递归ACL策略可以通过ACL ID导入，例如：

```
$ terraform import alibabacloudstack_dns_recursor_acl.example acl-12345678
```