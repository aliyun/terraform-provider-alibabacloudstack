---
subcategory: "云解析 DNS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_dns_gtm_addresspool"
sidebar_current: "docs-Alibabacloudstack-dns-dns_gtm_addresspool"
description: |-
  云解析全局调度地址池
---

# alibabacloudstack_dns_gtm_addresspool

使用Provider配置的凭证在指定的资源集创建云解析全局调度地址池。

## 示例用法

### 基础用法

```hcl

variable "name" {
  default = "tfacc16606"
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
```

## 参数说明

支持以下参数：

* `name` - (必填) 地址池的名称。
* `type` - (必填) 地址池的类型。取值：`A`（IPv4地址）、`AAAA`（IPv6地址）、`CNAME`（域名）。
* `lba_strategy` - (必填) 负载均衡策略。取值：`ALL_RR`（返回全部地址，非权重）、`RATIO`（按权重返回地址）。
* `addrs` - (必填) 地址列表。每个地址包含以下属性：
  * `value` - (必填) 地址值。
  * `mode` - (必填) 模式。取值：`SMART`（智能返回）、`ONLINE`（永远在线）、`OFFLINE`（永远离线）。
  * `lba_weight` - (可选) 权重。取值范围：0-100。仅当`lba_strategy`为`RATIO`时有效。

## 属性说明

以下属性会从API响应中导出：

* `id` - 地址池的ID。

## Import

云解析全局调度地址池可以使用地址池 ID 导入，例如：

```
$ terraform import alibabacloudstack_dns_gtm_addresspool.example <address_pool_id>
```