---
subcategory: "云解析 DNS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_universal_dns_lines"
sidebar_current: "docs-Alibabacloudstack-datasource-universal-dns-lines"
description: |-
  查询跨云解析路线列表
---

# alibabacloudstack_universal_dns_lines

查询阿里云跨云解析线路列表。

## 示例用法

```hcl

variable "name" {
  default = "tfacc20899"
}

resource "alibabacloudstack_universal_dns_line" "default" {
  name         = var.name
  v4_addresses = ["192.168.0.1"]
  v6_addresses = ["2020:148:2:28::", "2020:148:3:28::"]
}

data "alibabacloudstack_universal_dns_lines" "default" {
  name_regex = alibabacloudstack_universal_dns_line.default.name
}

```

## 参数说明
以下参数用于过滤查询结果：

* `name_regex` (字符串, 可选)：用于通过正则表达式过滤Universal DNS line名称。
* `ids` (字符串列表, 可选)：Universal DNS line ID列表，用于精确匹配指定ID的线路。

## 属性说明
以下属性被导出：

* `id` (字符串)：线路的唯一标识符。
* `create_timestamp` (整数)：创建时间戳（秒）。
* `name` (字符串)：线路名称。
* `priority` (整数)：线路优先级，1优先级最高，数值越大优先级越低。
* `update_timestamp` (整数)：更新时间戳（秒）。
* `v4_addresses` (字符串集合)：IPv4地址列表。
* `v6_addresses` (字符串集合)：IPv6地址列表。