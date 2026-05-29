---
subcategory: "云解析 DNS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_dns_forward_domains"
sidebar_current: "docs-Alibabacloudstack-datasource-dns-forward-domains"
description: |-
  查询DNS转发域名
---

# alibabacloudstack_dns_forward_domains

> 查询阿里云DNS全局转发域名列表

## 示例用法

```hcl

variable "name" {
  default = "tfacc6273812765840923665.test."
}

resource "alibabacloudstack_dns_forward_domain" "default" {
  name         = var.name
  remark       = "Created by Terraform"
  forward_mode = "FORWARD_FIRST"
  forwarders   = ["192.168.101.1"]
}

data "alibabacloudstack_dns_forward_domains" "default" {
  name_regex = alibabacloudstack_dns_forward_domain.default.name
}

```

## 参数说明
以下参数支持过滤结果：

* `forward_mode` (字符串, 可选)：转发模式，可选值为"FORWARD_FIRST"或"FORWARD_ONLY"。取值说明：
  - `FORWARD_ONLY`: 全部转发模式(建议)-全部请求只做转发(不作任何递归)。
  - `FORWARD_FIRST`: 优先转发模式-转发失败降级到互联网递归。
* `ids` (列表, 可选)：转发域名ID列表，用于过滤结果。只有ID在列表中的转发域名会被返回。
* `name` (字符串, 可选)：转发域名名称，用于过滤结果。匹配包含指定名称的转发域名。
* `name_regex` (字符串, 可选)：转发域名名称的正则表达式，用于过滤结果。只有名称匹配正则表达式的转发域名会被返回。

## 属性说明
以下属性被导出：

* `forward_domains` (列表)：匹配的转发域名列表。每个元素包含以下属性：
  * `id` (字符串)：转发域名ID。
  * `caller_uid` (字符串)：调用者UID。
  * `create_timestamp` (整数)：创建时间戳（秒）。
  * `forward_mode` (字符串)：转发模式。取值说明：
    - `FORWARD_ONLY`: 全部转发模式(建议)-全部请求只做转发(不作任何递归)。
    - `FORWARD_FIRST`: 优先转发模式-转发失败降级到互联网递归。
  * `forwarders` (集合)：转发器IP列表。
  * `name` (字符串)：转发域名名称。
  * `remark` (字符串)：备注。
  * `update_timestamp` (整数)：更新时间戳（秒）。
* `ids` (列表)：匹配的转发域名ID列表。