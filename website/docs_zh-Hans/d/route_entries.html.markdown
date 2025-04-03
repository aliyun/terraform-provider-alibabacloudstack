---
subcategory: "VPC"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_route_entries"
sidebar_current: "docs-alibabacloudstack-datasource-route-entries"
description: |-
    查询路由规则
---

# alibabacloudstack_route_entries

根据指定过滤条件列出当前凭证权限可以访问的路由规则列表


## 示例用法

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
  default = "tf-testAccRouteEntryConfig"
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

resource "alibabacloudstack_route_entry" "foo" {
  route_table_id        = "${alibabacloudstack_vpc.foo.route_table_id}"
  destination_cidrblock = "172.11.1.1/32"
  nexthop_type          = "Instance"
  nexthop_id            = "${alibabacloudstack_instance.foo.id}"
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

  vswitch_id         = "${alibabacloudstack_vswitch.foo.id}"
  allocate_public_ip = true

  # series III
  instance_type              = "${data.alibabacloudstack_instance_types.default.instance_types.0.id}"
  internet_max_bandwidth_out = 5

  system_disk_category = "cloud_efficiency"
  image_id             = "${data.alibabacloudstack_images.default.images.0.id}"
  instance_name        = "${var.name}"
}

data "alibabacloudstack_route_entries" "foo" {
  route_table_id = "${alibabacloudstack_route_entry.foo.route_table_id}"
}

output "route_entries" {
 data = data.alibabacloudstack_route_entries.foo
}

```

## 参数参考

支持以下参数：

* `route_table_id` - (必填，变更时重建) 路由表的 ID。
* `instance_id` - (可选) 下一跳的实例 ID。
* `type` - (可选) 路由条目的类型。
* `cidr_block` - (可选) 路由条目的目标 CIDR 块。

## 属性参考

除了上述参数外，还导出以下属性：

* `entries` - 路由条目列表。每个元素包含以下属性：
  * `type` - 路由条目的类型。
  * `next_hop_type` - 下一跳的类型。
  * `status` - 路由条目的状态。
  * `instance_id` - 下一跳的实例 ID。
  * `route_table_id` - 路由条所属的路由表的 ID。
  * `cidr_block` - 路由条目的目标 CIDR 块。