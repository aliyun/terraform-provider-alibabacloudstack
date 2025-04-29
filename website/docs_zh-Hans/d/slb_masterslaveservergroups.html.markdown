---
subcategory: "Server Load Balancer (SLB)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_slb_masterslaveservergroups"
sidebar_current: "docs-Alibabacloudstack-datasource-slb-masterslaveservergroups"
description: |- 
  查询负载均衡(SLB)主备服务器
---

# alibabacloudstack_slb_masterslaveservergroups
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_slb_master_slave_server_groups`

根据指定过滤条件列出当前凭证权限可以访问的负载均衡(SLB)主备服务器组列表。

## 示例用法

```hcl
data "alibabacloudstack_zones" default {
  available_resource_creation = "VSwitch"
  enable_details = true
}

resource "alibabacloudstack_vpc_vpc" "default" {
  vpc_name = "${var.name}_vpc"
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vpc_vswitch" "default" {
  name = "${var.name}_vsw"
  vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
  cidr_block = "172.16.0.0/24"
  zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
}

resource "alibabacloudstack_ecs_securitygroup" "default" {
  name   = "${var.name}_sg"
  vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
}

resource "alibabacloudstack_security_group_rule" "default" {
   type = "ingress"
   ip_protocol = "tcp"
   nic_type = "intranet"
   policy = "accept"
   port_range = "22/22"
   priority = 1
   security_group_id = "${alibabacloudstack_ecs_securitygroup.default.id}"
   cidr_ip = "172.16.0.0/24"
}

data "alibabacloudstack_images" "default" {
  name_regex  = "^ubuntu_"
  most_recent = true
  owners      = "system"
}

data "alibabacloudstack_instance_types" "all" {
  availability_zone = data.alibabacloudstack_zones.default.zones[0].id
}

data "alibabacloudstack_instance_types" "any_n4" {
  availability_zone = data.alibabacloudstack_zones.default.zones[0].id
  instance_type_family = "ecs.n4"
  sorted_by            = "Memory"
}

data "alibabacloudstack_instance_types" "default" {
  availability_zone = data.alibabacloudstack_zones.default.zones[0].id
  cpu_core_count       = 1
  memory_size          = 1
  instance_type_family = "ecs.n4"
  sorted_by            = "Memory"
}

locals {
  default_instance_type_id = try(element(sort(length(data.alibabacloudstack_instance_types.default.instance_types) > 0 ? data.alibabacloudstack_instance_types.default.ids : data.alibabacloudstack_instance_types.any_n4.ids), 0), sort(data.alibabacloudstack_instance_types.all.ids)[0])
}

resource "alibabacloudstack_ecs_instance" "default" {
  image_id             = "${data.alibabacloudstack_images.default.images.0.id}"
  instance_type        = "${local.default_instance_type_id}"
  system_disk_category = "${data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0}"
  system_disk_size     = 20
  system_disk_name     = "test_sys_disk"
  security_groups      = [alibabacloudstack_ecs_securitygroup.default.id]
  instance_name        = "${var.name}_ecs"
  vswitch_id           = alibabacloudstack_vpc_vswitch.default.id
  zone_id    = data.alibabacloudstack_zones.default.zones.0.id
  is_outdated          = false
  lifecycle {
    ignore_changes = [
      instance_type
    ]
  }
}

variable "name" {
  default = "tf-testAccslbmasterslaveservergroupsdatasourcebasic"
}

resource "alibabacloudstack_slb" "default" {
  name = "${var.name}"
  vswitch_id = "${alibabacloudstack_vpc_vswitch.default.id}"
}

resource "alibabacloudstack_ecs_instance" "new" {
  image_id             = "${data.alibabacloudstack_images.default.images.0.id}"
  instance_type        = "${local.default_instance_type_id}"
  system_disk_category = "${data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0}"
  system_disk_size     = 40
  system_disk_name     = "test_sys_diskv2"
  security_groups      = [alibabacloudstack_ecs_securitygroup.default.id]
  instance_name        = "${var.name}_ecs"
  vswitch_id           = alibabacloudstack_vpc_vswitch.default.id
  zone_id    = data.alibabacloudstack_zones.default.zones.0.id
  is_outdated          = false
  lifecycle {
    ignore_changes = [
      instance_type
    ]
  }
}

resource "alibabacloudstack_slb_master_slave_server_group" "default" {
  load_balancer_id = "${alibabacloudstack_slb.default.id}"
  name = "${var.name}"
  servers {
      server_id = "${alibabacloudstack_ecs_instance.default.id}"
      port = 80
      weight = 100
      server_type = "Master"
  }
  servers {
      server_id = "${alibabacloudstack_ecs_instance.new.id}"
      port = 80
      weight = 100
      server_type = "Slave"
  }
}

data "alibabacloudstack_slb_master_slave_server_groups" "default" {
  load_balancer_id = "${alibabacloudstack_slb.default.id}"
  ids              = ["${alibabacloudstack_slb_master_slave_server_group.default.id}"]
  name_regex       = "${var.name}"
  output_file      = "output.txt"
}
```

## 参数参考

以下参数是支持的：

* `load_balancer_id` - (必填) 关联的负载均衡实例ID。
* `ids` - (可选) 用于过滤结果的主备服务器组ID列表。
* `name_regex` - (可选，变更时重建) 用于通过主备服务器组名称过滤结果的正则表达式字符串。

## 属性参考

除了上述参数外，还导出以下属性：

* `names` - SLB主备服务器组名称列表。
* `groups` - SLB主备服务器组列表。每个元素包含以下属性：
  * `id` - 主备服务器组的ID。
  * `name` - 主备服务器组的名称。
  * `servers` - 与此组关联的ECS实例。每个元素包含以下属性：
    * `instance_id` - 已附加ECS实例的ID。
    * `weight` - 与ECS实例关联的权重。
    * `port` - 主备服务器组使用的端口。
    * `server_type` - 已附加ECS实例的服务器类型。