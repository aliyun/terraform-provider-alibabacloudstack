---
subcategory: "云解析 DNS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_universal_dns_domains"
sidebar_current: "docs-Alibabacloudstack-datasource-universal-dns-domains"
description: |-
  查询跨云解析域名列表
---

# alibabacloudstack_universal_dns_domains

查询阿里云跨云解析域名列表。

## 示例用法

```hcl

variable "name" {
  default = "tf-testacc6879360584559214283"
}

resource "alibabacloudstack_universal_dns_domain" "default" {
  name   = "${var.name}.example."
  remark = "Created by Terraform"
}

data "alibabacloudstack_universal_dns_domains" "default" {
  name_regex = alibabacloudstack_universal_dns_domain.default.name
}

```

## 参数说明
以下参数支持过滤查询结果：

* `ids` (可选)：用于过滤域名ID的列表。只有ID在该列表中的域名会被返回。

* `name_regex` (可选)：用于按域名名称过滤的正则表达式。只有名称匹配该正则表达式的域名会被返回。

## 属性说明
以下属性被导出：

* `id` (字符串)：数据源的唯一标识符，基于返回的域名ID列表的哈希值。

* `domains` (列表)：匹配条件的跨云DNS域名列表。每个域名包含以下属性：
  * `id` (字符串)：跨云DNS域名的ID。
  * `create_timestamp` (整数)：域名创建的时间戳（秒）。
  * `name` (字符串)：跨云DNS域名的名称。
  * `record_count` (整数)：域名中DNS记录的总数。
  * `remark` (字符串)：域名的备注或描述。
  * `update_timestamp` (整数)：域名最后更新的时间戳（秒）。

* `ids` (列表)：找到的域名对应的域名ID列表。

* `names` (列表)：找到的域名对应的域名名称列表。