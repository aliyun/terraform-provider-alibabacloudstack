---
subcategory: "Universal DNS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_dns_private_line"
sidebar_current: "docs-Alibabacloudstack-resource-dns-private-line"
description: |-
  管理阿里云云解析私有线路
---

# alibabacloudstack_dns_private_line

管理阿里云云解析私有线路，用于配置云解析服务的私有线路信息。

## 示例用法

### 基础用法

```hcl

variable "name" {
  default = "tfacc32375"
}


resource "alibabacloudstack_dns_private_line" "default" {
  name = var.name
  v4_addresses = [
    "192.168.0.1"
  ]
  v6_addresses = [
    "2020:148:2:28::",
    "2020:148:3:28::"
  ]
}
```

## 参数说明

支持以下参数：

* `name` - (必填) 云解析私有线路的名称。名称用于标识该线路，长度和格式限制请参考阿里云API文档。
* `v4_addresses` - (可选) IPv4地址列表。IPv4和IPv6地址至少需要填写一个。支持IPv4地址的模糊查询。
* `v6_addresses` - (可选) IPv6地址列表。IPv4和IPv6地址至少需要填写一个。IPv6地址只支持精确查询。

## 属性说明

以下属性会导出：

* `id` - 云解析私有线路的ID。
* `priority` - 云解析私有线路的优先级。1表示最高优先级，数值越大优先级越低。系统自动分配优先级值。