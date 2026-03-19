---
subcategory: "Cloud DNS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_dns_gtm_access_strategy"
sidebar_current: "docs-Alibabacloudstack-dns-dns_gtm_access_strategy"
description: |-
  云解析全局调度实例访问策略
---

# alibabacloudstack_dns_gtm_access_strategy

使用Provider配置的凭证在指定的资源集创建云解析全局调度实例访问策略。

## 示例用法

### 基础用法

```hcl

variable "name" {
  default = "tfacc35194"
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

resource "alibabacloudstack_dns_gtm_addresspool" "update" {
  name         = "${var.name}-update"
  type         = "A"
  lba_strategy = "ALL_RR"
  addrs {
    value = "3.3.3.3"
    mode  = "SMART"
  }
}


resource "alibabacloudstack_dns_gtm_access_strategy" "default" {
  default_gtm_address_pool_id    = alibabacloudstack_dns_gtm_addresspool.default.id
  failover_gtm_address_pool_id   = alibabacloudstack_dns_gtm_addresspool.failover.id
  failover_gtm_address_pool_type = "IPV4"
  default_min_available_addr_num = "1"
  line_ids = [
    "${alibabacloudstack_dns_private_line.default.id}"
  ]
  gtm_instance_id                 = alibabacloudstack_dns_gtm_instance.default.id
  default_gtm_address_pool_type   = "IPV4"
  failover_min_available_addr_num = "1"
  name                            = var.name
  switch_mode                     = "BY_PROBE_RESULT"
}
```

## 参数说明

支持以下参数：

* `default_gtm_address_pool_id` - (必填) 主地址池ID。
* `default_gtm_address_pool_type` - (必填) 主地址池类型。取值：DOMAIN（域名）、IPV6（IPV6类型）、IPV4（IPV4类型）。
* `default_min_available_addr_num` - (必填) 主地址池最小可用地址数量。
* `line_ids` - (必填) 线路ID列表，用于指定访问策略适用的线路。
* `name` - (必填) 访问策略的名称。
* `switch_mode` - (必填) 生效地址池切换策略。取值：BY_HAND（手动）、BY_PROBE_RESULT（自动）。
* `gtm_instance_id` - (必填, 变更时重建) 调度实例ID。
* `failover_gtm_address_pool_id` - (可选) 备用地址池ID。
* `failover_gtm_address_pool_type` - (可选) 备用地址池类型。取值：DOMAIN（域名）、IPV6（IPV6类型）、IPV4（IPV4类型）。
* `failover_min_available_addr_num` - (可选) 备用地址池最小可用地址数量。
* `specified_gtm_address_pool` - (可选) 手动指定当前使用的地址池。取值：DEFAULT（主地址池）、FAILOVER（备用地址池）。

## 属性说明

以下属性会从API中导出：

* `id` - 访问策略ID。
* `default_available_addr_num` - 主地址池当前可用地址数量。
* `default_gtm_address_pool_name` - 主地址池名称。
* `failover_available_addr_num` - 备用地址池当前可用地址数量。
* `failover_gtm_address_pool_name` - 备用地址池名称。
* `in_use_gtm_address_pool_id` - 当前使用的地址池ID。
* `in_use_gtm_address_pool_name` - 当前使用的地址池名称。