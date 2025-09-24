---
subcategory: "云企业网(CEN)"
layout: "alibabacloudstack"
page_title: "阿里云专有云: alibabacloudstack_cen_transit_router_multicast_domain_members"
sidebar_current: "docs-alibabacloudstack-datasource-cen-transit-router-multicast-domain-members"
description: |-
  提供阿里云专有云账户的CEN转发路由器组播域成员列表。
---

# alibabacloudstack\_cen\_transit\_router\_multicast\_domain\_members

本数据源根据指定的过滤条件，提供阿里云专有云账户中的CEN转发路由器组播域成员列表。

## 示例用法

```hcl
variable "name" {
  default = "tf-testAccMulticastDomainMember"
}

data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
  enable_details = true
}

resource "alibabacloudstack_vpc" "default" {
  vpc_name = "${var.name}_vpc"
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vswitch" "default" {
  name = "${var.name}_vsw"
  vpc_id = "${alibabacloudstack_vpc.default.id}"
  cidr_block = "172.16.1.0/24"
  zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
}

resource "alibabacloudstack_cen_instance" "default" {
  cen_instance_name = "${var.name}"
  description = "${var.name}"
}

resource "alibabacloudstack_cen_transit_router" "default" {
  cen_id = "${alibabacloudstack_cen_instance.default.id}"
}

resource "alibabacloudstack_cen_transit_router_multicast_domain" "default" {
  transit_router_multicast_domain_name = "${var.name}"
  transit_router_id = "${alibabacloudstack_cen_transit_router.default.transit_router_id}"
  cen_id = "${alibabacloudstack_cen_instance.default.id}"
}

resource "alibabacloudstack_network_interface" "default" {
  vswitch_id = "${alibabacloudstack_vswitch.default.id}"
}

resource "alibabacloudstack_cen_transit_router_multicast_domain_member" "default" {
  group_ip_address = "224.0.0.1"
  network_interface_id = "${alibabacloudstack_network_interface.default.id}"
  transit_router_multicast_domain_id = "${alibabacloudstack_cen_transit_router_multicast_domain.default.id}"
  vswitch_id = "${alibabacloudstack_vswitch.default.id}"
}

data "alibabacloudstack_cen_transit_router_multicast_domain_members" "default" {
  transit_router_multicast_domain_id = "${alibabacloudstack_cen_transit_router_multicast_domain.default.id}"
}
```

## 参数参考

以下参数被支持：

* `ids` - (可选) 组播域成员ID列表。
* `transit_router_multicast_domain_id` - (必填) 组播成员所属的组播域ID。
* `transit_router_attachment_id` - (可选) 转发路由器附件ID。
* `vswitch_id` - (可选) 交换机ID。

## 属性参考

以下属性会被导出：

* `transit_router_multicast_groups` - 组播域成员列表。每个元素包含以下属性：
  * `id` - 组播成员的ID。
  * `group_ip_address` - 组播IP地址。
  * `network_interface_id` - 网络接口ID。
  * `status` - 组播成员的状态。
  * `transit_router_multicast_domain_id` - 组播成员所属的组播域ID。
  * `transit_router_attachment_id` - 转发路由器附件ID。
  * `vswitch_id` - 交换机ID。
  * `resource_type` - 资源类型。
  * `member_type` - 成员类型。
  * `resource_id` - 资源ID。
  * `group_source` - 组是否为源组。
  * `group_member` - 组是否为成员组。