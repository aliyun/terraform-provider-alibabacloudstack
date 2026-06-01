---
subcategory: "云企业网 CEN"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cen_transit_router_multicast_domain_sources"
sidebar_current: "docs-Alibabacloudstack-datasource-cen-transit-router-multicast-domain-sources"
description: |-
  提供阿里云账户拥有的 cen transitroutermulticastdomainsources 列表。
---

# alibabacloudstack\_cen\_transitroutermulticastdomainsources

此数据源根据指定的过滤条件提供阿里云账户中的 cen transit router multicast domain sources 列表。

## Example Usage

```hcl
variable "name" {
  default = "tf-testaccmulticast_domain_source"
}

data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
  enable_details              = true
}

resource "alibabacloudstack_vpc_vpc" "default" {
  vpc_name   = "${var.name}_vpc"
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vpc_vswitch" "default" {
  name       = "${var.name}_vsw"
  vpc_id     = alibabacloudstack_vpc_vpc.default.id
  cidr_block = "172.16.1.0/24"
  zone_id    = data.alibabacloudstack_zones.default.zones.0.id
}

resource "alibabacloudstack_cen_instance" "default" {
  cen_instance_name       = var.name
  description             = var.name
  transit_router_name     = var.name
  transit_router_description = var.name
}

resource "alibabacloudstack_cen_transit_router_vpc_attachment" "default" {
  transit_router_attachment_name        = var.name
  transit_router_attachment_description = var.name
  cen_id                                = alibabacloudstack_cen_instance.default.id
  vpc_id                                = alibabacloudstack_vpc_vswitch.default.vpc_id
  transit_router_id                     = alibabacloudstack_cen_instance.default.transit_router_id
  zone_mappings {
    vswitch_id = alibabacloudstack_vpc_vswitch.default.id
    zone_id    = alibabacloudstack_vpc_vswitch.default.zone_id
  }
}

resource "alibabacloudstack_cen_transit_router_multicast_domain" "default" {
  transit_router_multicast_domain_description = var.name
  transit_router_multicast_domain_name        = var.name
  transit_router_id                           = alibabacloudstack_cen_instance.default.transit_router_id
}

resource "alibabacloudstack_cen_transit_router_multicast_domain_association" "default" {
  transit_router_attachment_id       = alibabacloudstack_cen_transit_router_vpc_attachment.default.transit_router_attachment_id
  transit_router_multicast_domain_id = alibabacloudstack_cen_transit_router_multicast_domain.default.transit_router_multicast_domain_id
  vswitch_id                         = alibabacloudstack_vpc_vswitch.default.id
}

data "alibabacloudstack_cen_transit_router_multicast_domain_sources" "default" {
  transit_router_multicast_domain_id = alibabacloudstack_cen_transit_router_multicast_domain.default.transit_router_multicast_domain_id
}
```

## Argument Reference

支持以下参数：

* `transit_router_multicast_domain_id` - (必选) 转发路由器组播域的 ID。
* `ids` - (可选, 可回读) 用于过滤结果的 ID 列表。
* `transit_router_attachment_id` - (可选) 转发路由器附件的 ID。
* `vswitch_id` - (可选) 交换机的 ID。

## Attributes Reference

除上述参数外，还导出以下属性：

* `transit_router_multicast_groups` - 转发路由器组播组列表。每个元素包含以下属性：
  * `id` - 组播组源的 ID，格式为 `<group_ip_address>:<transit_router_multicast_domain_id>:<resource_type>:<resource_data>`。
  * `group_ip_address` - 组播 IP 地址。
  * `network_interface_id` - 弹性网卡的 ID。
  * `transit_router_multicast_domain_id` - 转发路由器组播域的 ID。
  * `status` - 组播组源的状态。
  * `transit_router_attachment_id` - 转发路由器附件的 ID。
  * `vswitch_id` - 交换机的 ID。
  * `resource_type` - 资源类型。取值：`VPC`、`Connect`。
  * `resource_id` - 资源的 ID。
  * `source_type` - 源类型。
  * `group_source` - 是否为组播源。
  * `group_member` - 是否为组播成员。
* `ids` - 转发路由器组播域 ID 列表。
