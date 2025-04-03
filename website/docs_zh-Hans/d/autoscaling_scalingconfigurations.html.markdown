---
subcategory: "AutoScaling"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_autoscaling_scalingconfigurations"
sidebar_current: "docs-Alibabacloudstack-datasource-autoscaling-scalingconfigurations"
description: |- 
  查询弹性伸缩配置
---

# alibabacloudstack_autoscaling_scalingconfigurations
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_ess_scaling_configurations`

根据指定过滤条件列出当前凭证权限可以访问的弹性伸缩配置列表。

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
  default = "tf-testscalconf-218"
}

resource "alibabacloudstack_ess_scaling_group" "default" {
  min_size = 0
  max_size = 2
  default_cooldown = 20
  removal_policies = ["OldestInstance", "NewestInstance"]
  scaling_group_name = "${var.name}"
  vswitch_ids = ["${alibabacloudstack_vpc_vswitch.default.id}"]
}

resource "alibabacloudstack_ecs_deployment_set" "default" {
  strategy            = "Availability"
  domain              = "Default"
  granularity         = "Host"
  deployment_set_name = "example_value"
  description         = "example_value"
}

resource "alibabacloudstack_ess_scaling_configuration" "default" {
  scaling_group_id = "${alibabacloudstack_ess_scaling_group.default.id}"
  image_id = "${data.alibabacloudstack_images.default.images.0.id}"
  instance_type = "ecs.e4.small"
  security_group_ids = [alibabacloudstack_ecs_securitygroup.default.id]
  force_delete = true
  active = true
  enable = true
  deployment_set_id = alibabacloudstack_ecs_deployment_set.default.id
}

data "alibabacloudstack_ess_scaling_configurations" "default" {
  scaling_group_id = "${alibabacloudstack_ess_scaling_group.default.id}"
  name_regex = "${alibabacloudstack_ess_scaling_configuration.default.scaling_configuration_name}"
  output_file = "scaling_configurations_output.txt"
}

output "first_scaling_configuration_id" {
  value = data.alibabacloudstack_ess_scaling_configurations.default.configurations.0.id
}
```

## 参数参考

以下参数是支持的：

* `scaling_group_id` - (可选，变更时重建) 伸缩配置所属伸缩组的ID。
* `name_regex` - (可选，变更时重建) 按名称筛选结果的正则表达式字符串。
* `ids` - (可选) 用于过滤结果的伸缩配置ID列表。

## 属性参考

除了上述参数外，还导出以下属性：

* `ids` - 伸缩配置ID列表。
* `names` - 伸缩配置名称列表。
* `configurations` - 伸缩配置列表。每个元素包含以下属性：
  * `id` - 伸缩配置的ID。
  * `name` - 伸缩配置的名称。
  * `scaling_group_id` - 该伸缩配置所属的伸缩组的ID。
  * `image_id` - 使用此伸缩配置创建ECS实例时使用的镜像文件ID。
  * `instance_type` - 使用此伸缩配置创建的ECS实例的规格。
  * `security_group_id` - ECS实例所属的安全组ID。同一安全组内的实例可以互相访问。
  * `internet_max_bandwidth_in` - ECS实例的互联网入方向带宽值，单位为Mbps(Mega bit per second)，取值范围为1~200。如果不指定，默认设置为200Mbps。
  * `internet_max_bandwidth_out` - ECS实例的互联网出方向带宽值，单位为Mbps。
  * `system_disk_category` - 伸缩配置中使用的系统盘类别。
  * `system_disk_size` - 系统盘大小，单位为GB。
  * `data_disks` - 伸缩配置中配置的数据盘列表。每个数据盘具有以下属性：
    * `size` - 数据盘大小，单位为GB。
    * `category` - 数据盘类别。
    * `snapshot_id` - 用于创建数据盘的快照ID。
    * `device` - 数据盘的设备属性。
    * `delete_with_instance` - 实例释放时是否删除数据盘。
  * `lifecycle_state` - 伸缩配置的生命周期状态。
  * `creation_time` - 伸缩配置的创建时间。
