---
subcategory: "ECS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ecs_instances"
sidebar_current: "docs-Alibabacloudstack-datasource-ecs-instances"
description: |- 
  查询云服务器实例资

---

# alibabacloudstack_ecs_instances
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_instances`

根据指定过滤条件列出当前凭证权限可以访问的查询云服务器实例列表。

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
  type                = "ingress"
  ip_protocol         = "tcp"
  nic_type           = "intranet"
  policy             = "accept"
  port_range         = "22/22"
  priority           = 1
  security_group_id  = "${alibabacloudstack_ecs_securitygroup.default.id}"
  cidr_ip            = "172.16.0.0/24"
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
  availability_zone    = data.alibabacloudstack_zones.default.zones[0].id
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
  zone_id             = data.alibabacloudstack_zones.default.zones.0.id
  is_outdated          = false
  lifecycle {
    ignore_changes = [
      instance_type
    ]
  }
}

variable "name" {
  default = "Tf-EcsInstanceDataSource"
}

data "alibabacloudstack_ecs_instances" "default" {
  ids              = ["${alibabacloudstack_ecs_instance.default.id}"]
  name_regex       = "web_server"
  status           = "Running"
  vpc_id           = "${alibabacloudstack_vpc_vpc.default.id}"
  vswitch_id       = "${alibabacloudstack_vpc_vswitch.default.id}"
  availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
  tags = {
    Environment = "Production"
    Owner      = "TeamA"
  }
  ram_role_name    = "exampleRole"
  output_file     = "instances_output.txt"
}

output "first_instance_id" {
  value = "${data.alibabacloudstack_ecs_instances.default.instances.0.id}"
}

output "instance_ids" {
  value = "${data.alibabacloudstack_ecs_instances.default.ids}"
}
```

## 参数参考

以下参数是支持的：

* `ids` - (可选) ECS实例ID列表。如果指定，数据源将仅返回匹配ID的实例。
* `name_regex` - (可选) 用于通过实例名称过滤结果的正则表达式字符串。这允许您检索名称与指定模式匹配的实例。
* `image_id` - (可选) 一些ECS实例使用的镜像ID。可以根据此参数进行过滤。
* `status` - (可选) 实例状态。有效值包括："Creating"（创建中）、"Starting"（启动中）、"Running"（运行中）、"Stopping"（停止中）和 "Stopped"（已停止）。如果不指定，则考虑所有状态。
* `vpc_id` - (可选) 实例关联的VPC的ID。可以根据此参数进行过滤。
* `vswitch_id` - (可选) 实例关联的交换机的ID。可以根据此参数进行过滤。
* `availability_zone` - (可选) 实例所在的可用区。可以根据此参数进行过滤。
* `tags` - (可选) 分配给ECS实例的标签映射。可以根据这些标签进行过滤。例如：
* `ram_role_name` - (可选) 实例附加的RAM角色名称。可以根据此参数进行过滤。

## 属性说明

除了上述参数外，还导出以下属性：

* `ids` - ECS实例ID列表。
* `names` - 实例名称列表。
* `instances` - 实例列表。每个元素包含以下属性：
  * `id` - 实例的ID。
  * `region_id` - 实例所属的区域ID。
  * `availability_zone` - 实例所属的可用区。
  * `status` - 实例的当前状态。
  * `name` - 实例的名称。
  * `description` - 实例的描述。
  * `instance_type` - 实例的类型。
  * `instance_charge_type` - 实例的计费类型。
  * `vpc_id` - 实例所属的VPC的ID。
  * `vswitch_id` - 实例所属的交换机的ID。
  * `image_id` - 实例正在使用的镜像ID。
  * `private_ip` - 实例的私有IP地址。
  * `eip` - VPC实例正在使用的弹性公网IP地址。
  * `security_groups` - 实例所属的安全组ID列表。
  * `key_name` - 实例正在使用的密钥对。
  * `creation_time` - 实例的创建时间。
  * `internet_max_bandwidth_out` - 互联网最大输出带宽。
  * `tags` - 分配给ECS实例的标签映射。
  * `disk_device_mappings` - 已挂载磁盘的描述。
    * `device` - 创建的磁盘设备信息，例如 `/dev/xvdb`。
    * `size` - 创建的磁盘大小。
    * `category` - 云盘的类别。
    * `type` - 云盘的类型：系统盘或数据盘。
  * `ram_role_name` - 实例附加的RAM角色名称。