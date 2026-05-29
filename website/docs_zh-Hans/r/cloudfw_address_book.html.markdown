---
subcategory: "云防火墙"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cloudfw_address_book"
sidebar_current: "docs-Alibabacloudstack-resource-cloudfw-address-book"
description: |-
  云防火墙地址簿资源，用于管理云防火墙中的IP或端口地址簿。
---

# alibabacloudstack_cloudfw_address_book

云防火墙地址簿资源，用于在阿里云云防火墙中创建和管理IP或端口地址簿。

## 示例用法

### 基础用法

```hcl

variable "name" {
  default = "tf-testacc-addrbook-75750"
}


resource "alibabacloudstack_cloudfw_address_book" "default" {
  address_list = [
    "8888",
    "9999"
  ]
  group_type  = "port"
  group_name  = var.name
  description = var.name
}
```

## 参数说明

支持以下参数：

* `group_name` - (必填, 变更时重建) 地址簿名称。名称长度需符合云防火墙要求，不能包含特殊字符。
* `group_type` - (必填, 变更时重建) 地址簿类型。有效值：`ip`（IP地址簿）、`port`（端口地址簿）。
* `description` - (必填) 地址簿描述信息。用于说明地址簿的用途或内容。
* `address_list` - (必填) 地址列表。至少包含一个地址项，IP地址簿格式为CIDR（如192.168.1.0/24），端口地址簿格式为端口号（如80,443）。

## 属性说明

以下属性会从云防火墙API中导出：

* `id` - 资源ID，格式为"{group_type}:{group_uuid}"。
* `address_list_count` - 地址列表总数，表示地址簿中包含的地址数量。
* `auto_add_tag_ecs` - ECS自动添加标记状态。0表示未启用，1表示已启用。
* `global` - 全局标记。0表示非全局地址簿，1表示全局地址簿。
* `group_uuid` - 地址簿唯一标识符，由云防火墙系统生成的UUID。
* `reference_count` - 引用总计，表示该地址簿被规则引用的次数。