---
subcategory: "Cloud DNS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_dns_gtm_instances"
sidebar_current: "docs-Alibabacloudstack-datasource-dns_gtm_instances"
description: |-
  查询DNS GTM（全局流量管理）实例列表
---

# alibabacloudstack_dns_gtm_instances

查询DNS GTM（全局流量管理）实例列表，用于获取已创建的云解析调度实例信息。

## 示例用法

```hcl

variable "name" {
  default = "tfacc60948"
}

resource "alibabacloudstack_vpc_vpc" "default" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = "${var.name}_vpc0"
}

resource "alibabacloudstack_dns_private_domain" "default" {
  name   = "${var.name}.testtf."
  remark = "Created by Terraform for DNS GTM instance test"
  vpc_ids = [
    "${alibabacloudstack_vpc_vpc.default.id}"
  ]
}

resource "alibabacloudstack_dns_gtm_instance" "default" {
  name    = var.name
  prefix  = var.name
  zone_id = alibabacloudstack_dns_private_domain.default.id
  ttl     = 300
}

data "alibabacloudstack_dns_gtm_instances" "default" {
  name_regex = alibabacloudstack_dns_gtm_instance.default.name
}

```

## 参数说明
以下参数用于过滤查询结果：

- `ids` (列表, 可选)：按照一个或多个实例ID进行过滤。实例ID形如 `886c6195-36d7-4d97-8bb5-d7f4228ec55c`。列表中的多个ID表示逻辑或关系。

- `name_regex` (字符串, 可选)：按照实例名称进行正则表达式匹配过滤。例如，设置为 `^test.*` 将匹配所有以 "test" 开头的实例名称。

## 属性说明
以下属性被导出：

- `id` (字符串)：数据源ID，由实例ID列表计算得出的哈希值。

- `instances` (列表)：查询结果中的DNS GTM实例列表。每个元素包含以下属性：
  - `create_timestamp` (整数)：实例创建时间戳（秒）。
  - `id` (字符串)：DNS GTM实例的唯一标识符。
  - `name` (字符串)：DNS GTM实例的名称。
  - `prefix` (字符串)：调度域名前缀。
  - `ttl` (整数)：全局TTL（Time To Live）值，单位为秒。
  - `update_timestamp` (整数)：实例最后更新时间戳（秒）。
  - `zone_id` (字符串)：归属调度域的域名ID。
  - `zone_name` (字符串)：域名名称。