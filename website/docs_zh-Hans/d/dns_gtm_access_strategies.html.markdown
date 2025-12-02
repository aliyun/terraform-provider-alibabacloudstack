---
subcategory: "DNS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_dns_gtm_access_strategy"
sidebar_current: "docs-Alibabacloudstack-datasource-dns_gtm_access_strategy"
description: |-
  云解析全局调度实例访问策略
---

# alibabacloudstack_dns_gtm_access_strategy

查询云解析全局流量管理(GTM)实例的访问策略。

## 示例用法

```hcl

variable "name" {
  default = "tfacc3910"
}

resource "alibabacloudstack_dns_private_domain" "default" {
  name = "${var.name}.local."
}

resource "alibabacloudstack_dns_gtm_instance" "default" {
  name    = var.name
  prefix  = var.name
  zone_id = alibabacloudstack_dns_private_domain.default.id
  ttl     = 300
}

resource "alibabacloudstack_dns_private_line" "default" {
  name         = var.name
  v4_addresses = ["192.168.0.1"]
  v6_addresses = ["2020:148:2:28::", "2020:148:3:28::"]
}

resource "alibabacloudstack_dns_gtm_addresspool" "default" {
  name         = "${var.name}-default"
  type         = "A"
  lba_strategy = "ALL_RR"
  addrs {
    value = "1.1.1.1"
    mode  = "SMART"
  }
}

resource "alibabacloudstack_dns_gtm_addresspool" "failover" {
  name         = "${var.name}-failover"
  type         = "A"
  lba_strategy = "ALL_RR"
  addrs {
    value = "2.2.2.2"
    mode  = "SMART"
  }
}

resource "alibabacloudstack_dns_gtm_access_strategy" "default" {
  name                            = var.name
  gtm_instance_id                 = alibabacloudstack_dns_gtm_instance.default.id
  switch_mode                     = "BY_PROBE_RESULT"
  default_min_available_addr_num  = 1
  failover_min_available_addr_num = 1
  line_ids                        = [alibabacloudstack_dns_private_line.default.id]
  default_gtm_address_pool_id     = alibabacloudstack_dns_gtm_addresspool.default.id
  default_gtm_address_pool_type   = "IPV4"
  failover_gtm_address_pool_id    = alibabacloudstack_dns_gtm_addresspool.failover.id
  failover_gtm_address_pool_type  = "IPV4"
}

data "alibabacloudstack_dns_gtm_access_strategies" "default" {
  gtm_instance_id = alibabacloudstack_dns_gtm_instance.default.id
  name_regex      = alibabacloudstack_dns_gtm_access_strategy.default.name
}

```

## 参数说明
以下参数支持配置：

* `gtm_instance_id` (字符串, 必填)：调度实例的ID。

* `name_regex` (字符串, 可选)：用于过滤访问策略名称的正则表达式。

* `ids` (列表, 可选)：用于过滤特定访问策略ID的列表。

## 属性说明
以下属性被导出：

* `id` (字符串)：访问策略的唯一标识符，格式为"GtmInstanceId:AccessStrategyId"。

* `default_available_addr_num` (整数)：主地址池当前可用地址数量。

* `default_gtm_address_pool_id` (字符串)：主地址池ID。

* `default_gtm_address_pool_name` (字符串)：主地址池名称。

* `default_gtm_address_pool_type` (字符串)：主地址池类型。取值说明：DOMAIN（域名）、IPV6（IPV6类型）、IPV4（IPV4类型）。

* `default_min_available_addr_num` (整数)：主地址池最小可用地址数量。

* `failover_available_addr_num` (整数)：备地址池当前可用地址数量。

* `failover_gtm_address_pool_id` (字符串)：备地址池ID。

* `failover_gtm_address_pool_name` (字符串)：备地址池名称。

* `failover_gtm_address_pool_type` (字符串)：备地址池类型。取值说明：DOMAIN（域名）、IPV6（IPV6类型）、IPV4（IPV4类型）。

* `failover_min_available_addr_num` (整数)：备地址池最小可用地址数量。

* `gtm_instance_id` (字符串)：调度实例ID。

* `in_use_gtm_address_pool_id` (字符串)：当前使用的地址池ID。

* `in_use_gtm_address_pool_name` (字符串)：当前使用的地址池名称。

* `line_ids` (集合)：线路ID列表。

* `name` (字符串)：访问策略名称。

* `specified_gtm_address_pool` (字符串)：手动指定当前使用的地址池。取值说明：DEFAULT（主地址池）、FAILOVER（备地址池）。

* `switch_mode` (字符串)：生效地址池切换策略。取值说明：BY_HAND（手动）、BY_PROBE_RESULT（自动）。