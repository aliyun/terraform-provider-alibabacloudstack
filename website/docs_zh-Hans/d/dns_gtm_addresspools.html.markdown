---
subcategory: "云解析 DNS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_dns_gtm_addresspools"
sidebar_current: "docs-Alibabacloudstack-datasource-dns-gtm-addresspools"
description: |-
  查询阿里云DNS GTM地址池列表
---

# alibabacloudstack_dns_gtm_addresspools

查询阿里云DNS GTM（全局流量管理）地址池列表。该数据源可用于检索已创建的地址池信息，支持通过ID、名称正则表达式、类型和负载均衡策略进行过滤。

## 示例用法

```hcl

variable "name" {
  default = "tfacc45249"
}

resource "alibabacloudstack_dns_gtm_addresspool" "default" {
  name         = var.name
  type         = "A"
  lba_strategy = "RATIO"

  addrs {
    value      = "192.168.1.1"
    mode       = "SMART"
    lba_weight = 20
  }

  addrs {
    value      = "127.0.0.1"
    mode       = "SMART"
    lba_weight = 80
  }
}

data "alibabacloudstack_dns_gtm_addresspools" "default" {
  name_regex = alibabacloudstack_dns_gtm_addresspool.default.name
}

```

## 参数说明
以下参数支持过滤查询结果：

### 可选参数
* `ids` (字符串列表)：地址池ID列表，用于精确匹配指定ID的地址池。
* `lba_strategy` (字符串)：负载均衡策略过滤条件。有效值：`ALL_RR`（返回全部地址，非权重），`RATIO`（按权重返回地址）。
* `name_regex` (字符串)：地址池名称的正则表达式过滤器，用于模糊匹配地址池名称。
* `type` (字符串)：地址池类型过滤条件。有效值：`A`（IPv4地址），`AAAA`（IPv6地址），`CNAME`（域名）。

## 属性说明
以下属性被导出：

* `id` (字符串)：数据源的唯一标识符，由匹配的地址池ID哈希生成。
* `address_pools` (列表)：匹配的地址池列表，每个元素包含以下属性：
  * `addrs` (列表)：地址池中的地址列表，每个元素包含：
    * `available` (布尔)：指示地址是否可用。
    * `lba_weight` (整数)：地址的负载均衡权重。
    * `mode` (字符串)：地址模式。有效值：`SMART`（智能返回），`ONLINE`（永远在线），`OFFLINE`（永远离线）。
    * `value` (字符串)：地址值（IP或域名）。
  * `create_timestamp` (整数)：地址池创建时间戳（秒）。
  * `id` (字符串)：地址池的唯一标识符。
  * `lba_strategy` (字符串)：地址池的负载均衡策略。有效值：`ALL_RR`，`RATIO`。
  * `name` (字符串)：地址池的名称。
  * `type` (字符串)：地址池的类型。有效值：`A`，`AAAA`，`CNAME`。
  * `update_timestamp` (整数)：地址池最后更新时间戳（秒）。
* `ids` (字符串列表)：匹配的地址池ID列表。