---
subcategory: "NATGateway"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_natgateway_forwardentry"
sidebar_current: "docs-Alibabacloudstack-natgateway-forwardentry"
description: |- 
  编排专有网络的NAT网关DNAT表规则
---

# alibabacloudstack_natgateway_forwardentry
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_forward_entry`

使用Provider配置的凭证在指定的资源集编排专有网络的NAT网关DNAT表规则。

## 示例用法

```hcl
variable "name" {
  default = "tf-testAccForwardEntryConfig17430"
}

variable "number" {
  default = "2"
}

data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alibabacloudstack_vpc" "default" {
  cidr_block = "172.16.0.0/12"
}

resource "alibabacloudstack_vswitch" "default" {
  vpc_id            = alibabacloudstack_vpc.default.id
  cidr_block        = "172.16.0.0/21"
  availability_zone = data.alibabacloudstack_zones.default.zones.0.id
}

resource "alibabacloudstack_nat_gateway" "default" {
  vpc_id        = alibabacloudstack_vswitch.default.vpc_id
  specification = "Small"
}

resource "alibabacloudstack_eip" "default" {
  count = var.number
}

resource "alibabacloudstack_eip_association" "default" {
  count          = var.number
  allocation_id  = alibabacloudstack_eip.default[count.index].id
  instance_id    = alibabacloudstack_nat_gateway.default.id
}

resource "alibabacloudstack_forward_entry" "default" {
  name              = var.name
  forward_table_id = alibabacloudstack_nat_gateway.default.forward_table_ids[0]
  external_ip      = alibabacloudstack_eip.default[0].ip_address
  external_port    = "80"
  ip_protocol      = "tcp"
  internal_ip     = "172.16.0.4"
  internal_port   = "8080"
}
```

## 参数说明

支持以下参数：

* `forward_table_id` - (必填，变更时重建) DNAT 条目所属的 DNAT 表的 ID。此字段在创建后不可更改。
* `external_ip` - (必填) 公网 IP 地址。该地址用于 ECS 实例接收来自互联网的请求。
* `external_port` - (必填) 外部端口。该端口用于 ECS 实例接收来自互联网的请求。有效值为 1 到 65535 或 "any"。
* `ip_protocol` - (必填) 协议类型。有效值为 `tcp`、`udp` 或 `any`。
* `name` - (可选) DNAT 条目的名称。如果未提供，则默认使用 `forward_entry_name` 的值。
* `forward_entry_name` - (可选) DNAT 条目的名称。如果未提供，则默认使用 `name` 的值。
* `internal_ip` - (必填) 私网 IP 地址。它必须是 VPC 内的有效私有 IP 地址。
* `internal_port` - (必填) 目标私网端口。有效值为 1 到 65535 或 "any"。

## 属性说明

除了上述所有参数外，还导出了以下属性：

* `id` - DNAT 条目的唯一标识符。格式为 `<forward_table_id>:<forward_entry_id>`。
* `forward_entry_id` - DNAT 条目在服务器上的唯一标识符。
* `forward_entry_name` - DNAT 条目的名称。如果未明确设置，默认为 `name` 字段的值。