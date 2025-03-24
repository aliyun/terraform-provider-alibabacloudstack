---
subcategory: "AutoScaling"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_autoscaling_scalinggroups"
sidebar_current: "docs-Alibabacloudstack-datasource-autoscaling-scalinggroups"
description: |- 
  查询弹性伸缩组
---

# alibabacloudstack_autoscaling_scalinggroups
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_ess_scaling_groups`

根据指定过滤条件列出当前凭证权限可以访问的弹性伸缩组列表。

## 示例用法

```hcl
data "alibabacloudstack_zones" "default" {
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
  default = "tf-test-108"
}

resource "alibabacloudstack_ess_scaling_group" "default" {
  min_size = 0
  max_size = 2
  scaling_group_name = "${var.name}"
  default_cooldown = 20
  removal_policies = ["OldestInstance", "NewestInstance"]
  vswitch_ids = ["${alibabacloudstack_vpc_vswitch.default.id}"]
}

data "alibabacloudstack_ess_scaling_groups" "default" {
  name_regex = "${alibabacloudstack_ess_scaling_group.default.scaling_group_name}"
}

output "scaling_group_id" {
  value = "${data.alibabacloudstack_ess_scaling_groups.default.groups.0.id}"
}
```

## 参数参考

以下参数是支持的：

* `name_regex` - (可选) 用于按名称过滤结果伸缩组的正则表达式字符串。
* `ids` - (可选) 伸缩组 ID 列表，用于精确匹配特定伸缩组。

## 属性参考

除了上述参数外，还导出以下属性：

* `ids` - 伸缩组 ID 列表。
* `names` - 伸缩组名称列表。
* `groups` - 伸缩组列表。每个元素包含以下属性：
  * `id` - 伸缩组的 ID。
  * `name` - 伸缩组的名称。
  * `active_scaling_configuration` - 伸缩组的活动伸缩配置。
  * `region_id` - 伸缩组所属的地域的 ID。
  * `min_size` - 伸缩组内 ECS 实例台数的最小值，当伸缩组内 ECS 实例数小于 MinSize 时，弹性伸缩会自动创建 ECS 实例。
  * `max_size` - 伸缩组内 ECS 实例台数的最大值，当伸缩组内 ECS 实例数大于 MaxSize 时，弹性伸缩会自动移出 ECS 实例。MaxSize 的取值范围和弹性伸缩使用情况有关，请前往配额中心查看单个伸缩组可以设置的组内最大实例数对应的配额值。例如，如果单个伸缩组可以设置的组内最大实例数对应的配额值为 2000，则 MaxSize 的取值范围为 0~2000。
  * `cooldown_time` - 伸缩组的默认冷却时间，单位为秒。在此期间，不会触发任何伸缩活动。
  * `removal_policies` - 减少实例数量时用于选择从伸缩组中移除的 ECS 实例的移除策略。
  * `load_balancer_ids` - 伸缩组中的 ECS 实例所附加的 SLB(服务器负载均衡器)实例 ID 列表。
  * `db_instance_ids` - 伸缩组中的 ECS 实例所附加的 RDS(关系型数据库服务)实例 ID 列表。
  * `vswitch_ids` - 启动 ECS 实例所在的交换机 ID 列表。
  * `lifecycle_state` - 伸缩组的生命周期状态。可能的值包括 `Active`、`Inactive` 等。
  * `total_capacity` - 伸缩组内的所有 ECS 实例的数量。
  * `active_capacity` - 已成功加入伸缩组并正常运行的 ECS 实例数量。
  * `pending_capacity` - 正在加入伸缩组但尚未完成相关配置的 ECS 实例数量。
  * `removing_capacity` - 正在从伸缩组中移除的 ECS 实例数量。
  * `creation_time` - 伸缩组的创建时间，格式为 ISO 8601。
