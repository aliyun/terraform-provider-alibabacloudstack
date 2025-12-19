---
subcategory: "Bare Metal Server (BMS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_bms_keypairs"
sidebar_current: "docs-Alibabacloudstack-datasource-bms-keypairs"
description: |-
  查询阿里云裸金属服务器(BMS)的密钥对信息
---

# alibabacloudstack_bms_keypairs

查询阿里云裸金属服务器(BMS)的密钥对信息。

## 示例用法

```hcl

variable "name" {
  default = "tf-bms-keypair12432"
}

resource "alibabacloudstack_bms_keypair" "default" {
  name = var.name
}

data "alibabacloudstack_bms_keypairs" "default" {
  name_regex = alibabacloudstack_bms_keypair.default.name
}

```

## 参数说明
以下参数可供配置：

* `ids` (列表)：密钥对名称列表，用于过滤结果。如果指定，只会返回名称在列表中的密钥对。

* `name_regex` (字符串)：名称正则表达式，用于过滤结果。如果指定，只会返回名称匹配正则表达式的密钥对。

## 属性说明
以下属性被导出：

* `id` (字符串)：密钥对ID，与名称相同。在BMS中，密钥对名称是唯一的，因此用作ID。

* `create_time` (字符串)：密钥对创建时间。

* `key_pair_fingerprint` (字符串)：密钥对指纹信息。

* `name` (字符串)：密钥对名称。

* `public_key` (字符串)：密钥对的公钥内容。

* `shared` (整数)：密钥对共享状态，0表示不共享，1表示共享。

* `update_time` (字符串)：密钥对最后更新时间。