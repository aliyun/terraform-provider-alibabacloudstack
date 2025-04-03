---
subcategory: "AutoScaling"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_autoscaling_scalingrules"
sidebar_current: "docs-Alibabacloudstack-datasource-autoscaling-scalingrules"
description: |- 
  查询弹性伸缩规则
---

# alibabacloudstack_autoscaling_scalingrules
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_ess_scaling_rules`

根据指定过滤条件列出当前凭证权限可以访问的弹性伸缩规则列表。

## 示例用法

```hcl
variable "name" {
  default = "tf-testAccDataSourceEssScalingRules"
}

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

resource "alibabacloudstack_ess_scaling_rule" "default" {
  scaling_group_id = "${alibabacloudstack_ess_scaling_group.default.id}"
  adjustment_type = "TotalCapacity"
  adjustment_value = "1"
  cooldown = 0
}

data "alibabacloudstack_ess_scaling_rules" "default" {
  scaling_group_id = "${alibabacloudstack_ess_scaling_group.default.id}"
  ids = ["${alibabacloudstack_ess_scaling_rule.default.id}"]
  name_regex = "rule-name-*"
  type = "SimpleScalingConfiguration"

  output_file = "scaling_rules_output.json"
}

output "first_scaling_rule_id" {
  value = data.alibabacloudstack_ess_scaling_rules.default.rules[0].id
}
```

## 参数参考

以下参数是支持的：

* `scaling_group_id` - （选填，变更时重建）伸缩组的ID。用于筛选属于特定伸缩组的伸缩规则。
* `name_regex` - （选填，变更时重建）应用于伸缩规则名称的正则表达式字符串。这允许您仅检索那些名称与给定正则表达式匹配的规则。
* `ids` - （选填）伸缩规则ID列表。这可以用来过滤结果，只包括具有特定ID的伸缩规则。
* `type` - （选填，变更时重建）伸缩规则的类型。有效值包括：
  - `SimpleScalingConfiguration`: 简单伸缩规则。
  - `TargetTrackingConfiguration`: 目标跟踪伸缩规则。
  - `StepScalingConfiguration`: 阶梯伸缩规则。

## 属性参考

除了上述参数外，还导出以下属性：

* `names` - 所有匹配指定过滤器的伸缩规则的名称列表。
* `rules` - 伸缩规则列表。每个元素包含以下属性：
  * `id` - 伸缩规则的ID。
  * `scaling_group_id` - 伸缩规则所属的伸缩组的ID。
  * `name` - 伸缩规则的名称。
  * `type` - 伸缩规则的类型。可能的值包括 `SimpleScalingConfiguration`、`TargetTrackingConfiguration` 和 `StepScalingConfiguration`。
  * `cooldown` - 伸缩规则的冷却时间，在此期间不会触发额外的伸缩活动。这仅适用于简单伸缩规则。取值范围：0~86400 秒。
  * `adjustment_type` - 伸缩规则中使用的调整类型。可能的值为：
    - `QuantityChangeInCapacity`: 按指定数量增加或减少ECS实例的数量。
    - `PercentChangeInCapacity`: 按指定百分比增加或减少ECS实例的数量。
    - `TotalCapacity`: 将伸缩组中的ECS实例总数设置为指定值。
  * `adjustment_value` - 伸缩规则中使用的调整值。这指定了基于 `adjustment_type` 的调整幅度。
  * `min_adjustment_magnitude` - 伸缩规则的最小调整幅度。这确保任何伸缩活动都达到最低阈值。
  * `scaling_rule_ari` - 伸缩规则的唯一ARN（Aliyun资源名称）。
