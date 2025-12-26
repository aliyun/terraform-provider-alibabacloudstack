---
subcategory: "Cloud Firewall"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cloudfw_address_books"
sidebar_current: "docs-Alibabacloudstack-datasource-cloudfw-address-books"
description: |-
  查询云防火墙地址簿信息
---

# alibabacloudstack_cloudfw_address_books

云防火墙地址簿数据源，用于查询云防火墙中已创建的地址簿信息。

## 示例用法

```hcl

variable "name" {
  default = "tf-testacc-addrbook-11865"
}

resource "alibabacloudstack_cloudfw_address_book" "default" {
  group_type   = "ip"
  group_name   = var.name
  address_list = ["100.100.100.100/30"]
  description  = "test address book"
}

data "alibabacloudstack_cloudfw_address_books" "default" {
  group_type = "ip"
  ids        = ["${alibabacloudstack_cloudfw_address_book.default.id}"]
}

```

## 参数说明
以下参数支持过滤查询结果：

* `contain_port` (可选) - 过滤包含特定端口的地址簿。
* `group_type` (可选) - 地址簿类型，可选值为`ip`或`port`。
* `ids` (可选) - 地址簿UUID列表，用于过滤结果。
* `name_regex` (可选) - 用于通过正则表达式过滤地址簿名称。
* `query` (可选) - 通过名称关键字查询地址簿。

## 属性说明
以下属性被导出：

* `id` (字符串) - 地址簿ID，格式为`{group_type}:{group_uuid}`。
* `address_list` (列表) - 地址簿中的IP地址或端口列表。
* `address_list_count` (整数) - 地址列表中的条目数量。
* `auto_add_tag_ecs` (整数) - 是否自动添加标签ECS实例（0: 否，1: 是）。
* `description` (字符串) - 地址簿描述信息。
* `global` (整数) - 是否为全局地址簿（0: 否，1: 是）。
* `group_name` (字符串) - 地址簿名称。
* `group_type` (字符串) - 地址簿类型（ip或port）。
* `group_uuid` (字符串) - 地址簿唯一标识UUID。
* `reference_count` (整数) - 地址簿被引用的次数。
* `tag_relation` (字符串) - 标签关系配置信息。