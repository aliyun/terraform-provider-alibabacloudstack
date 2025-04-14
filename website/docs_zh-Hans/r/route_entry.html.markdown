---
subcategory: "VPC"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_route_entry"
sidebar_current: "docs-alibabacloudstack-resource-route-entry"
description: |-
  集编排路由条目
---

# alibabacloudstack_route_entry

使用Provider配置的凭证在指定的资源集编排路由条目。路由条目表示 VPC 路由表中的一个路由项。

## 示例用法

### 基础用法

```
data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
}
data "alibabacloudstack_instance_types" "default" {
  availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
  cpu_core_count    = 1
  memory_size       = 2
}
data "alibabacloudstack_images" "default" {
  name_regex  = "^ubuntu_18.*64"
  most_recent = true
  owners      = "system"
}

variable "name" {
  default = "RouteEntryConfig"
}
resource "alibabacloudstack_vpc" "foo" {
  name       = "${var.name}"
  cidr_block = "10.1.0.0/21"
}

resource "alibabacloudstack_vswitch" "foo" {
  vpc_id            = "${alibabacloudstack_vpc.foo.id}"
  cidr_block        = "10.1.1.0/24"
  availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
  name              = "${var.name}"
}

resource "alibabacloudstack_security_group" "tf_test_foo" {
  name        = "${var.name}"
  description = "foo"
  vpc_id      = "${alibabacloudstack_vpc.foo.id}"
}

resource "alibabacloudstack_security_group_rule" "ingress" {
  type              = "ingress"
  ip_protocol       = "tcp"
  nic_type          = "intranet"
  policy            = "accept"
  port_range        = "22/22"
  priority          = 1
  security_group_id = "${alibabacloudstack_security_group.tf_test_foo.id}"
  cidr_ip           = "0.0.0.0/0"
}

resource "alibabacloudstack_instance" "foo" {
  security_groups = ["${alibabacloudstack_security_group.tf_test_foo.id}"]

  vswitch_id = "${alibabacloudstack_vswitch.foo.id}"
  instance_type              = "${data.alibabacloudstack_instance_types.default.instance_types.0.id}"
  internet_max_bandwidth_out = 5
  system_disk_category = "cloud_efficiency"
  image_id             = "${data.alibabacloudstack_images.default.images.0.id}"
  instance_name        = "${var.name}"
}
resource "alibabacloudstack_route_entry" "foo" {
  route_table_id        = "${alibabacloudstack_vpc.foo.route_table_id}"
  destination_cidrblock = "172.11.1.1/32"
  nexthop_type          = "Instance"
  nexthop_id            = "${alibabacloudstack_instance.foo.id}"
}
```


## 参数说明

支持以下参数：

* `route_table_id` - (必填，强制更新) 路由表的 ID。
* `router_id` - (已弃用) 此参数已被弃用，请使用其他参数创建自定义路由条目。
* `destination_cidrblock` - (强制更新) 目标网络段。
* `nexthop_type` - (强制更新) 下一跳类型。可用值：
    - `Instance` (默认值): 将目标 CIDR 块的流量路由到 VPC 中的 ECS 实例。
    - `RouterInterface`: 将目标 CIDR 块的流量路由到路由器接口。
    - `VpnGateway`: 将目标 CIDR 块的流量路由到 VPN 网关。
    - `HaVip`: 将目标 CIDR 块的流量路由到 HAVIP。
    - `NetworkInterface`: 将目标 CIDR 块的流量路由到 NetworkInterface。
    - `NatGateway`: 将目标 CIDR 块的流量路由到 Nat Gateway。
* `nexthop_id` - (强制更新) 路由条目的下一跳。ECS 实例 ID 或 VPC 路由器接口 ID。
* `name` - (可选，强制更新，1.55.1+ 版本中可用) 路由条目的名称。该名称可以包含 2 到 128 个字符，必须仅包含字母数字字符或连字符（例如 `-`、`.`、`_`），并且不能以连字符开头或结尾，也不能以 `http://` 或 `https://` 开头。
* `router_id` - (可选) 此参数已被弃用，请从模板中移除。

## 属性说明

导出以下属性：

* `id` - 路由条目 ID，格式为 `<route_table_id:router_id:destination_cidrblock:nexthop_type:nexthop_id>`。
* `route_table_id` - 路由表的 ID。
* `destination_cidrblock` - 目标网络段。
* `nexthop_type` - 下一跳类型。
* `nexthop_id` - 路由条目的下一跳。
* `router_id` - (计算) 与路由条目关联的路由器 ID。