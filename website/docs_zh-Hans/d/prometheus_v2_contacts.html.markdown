---
subcategory: "Prometheus"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_prometheus_v2_contact"
sidebar_current: "docs-Alibabacloudstack-datasource-prometheus_v2_contact"
description: |-
  查询阿里云Prometheus v2联系人信息
---

# alibabacloudstack_prometheus_v2_contact

查询阿里云Prometheus v2联系人信息，用于获取已创建的联系人列表。

## 示例用法

```hcl

variable "name" {
  default = "tfacc_prometheus31218"
}
resource "alibabacloudstack_prometheus_v2_contact" "default" {
  username = "tfacc-${var.name}"
  mobile   = "13812345678"
  mail     = "test@example.com"
}

data "alibabacloudstack_prometheus_v2_contacts" "default" {
  name_regex = alibabacloudstack_prometheus_v2_contact.default.username
}

```

## 参数说明
以下参数用于过滤查询结果：

* `ids` (可选)：联系人ID列表，用于精确匹配指定ID的联系人。
* `name_regex` (可选)：用户名正则表达式，用于模糊匹配联系人名称。

## 属性说明
以下属性被导出：

* `id` (字符串)：联系人的唯一标识符。
* `groups` (列表)：联系人所属的组ID列表。
* `mail` (字符串)：联系人的邮箱地址。
* `mobile` (字符串)：联系人的手机号码。
* `username` (字符串)：联系人的用户名。