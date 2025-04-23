---
subcategory: "Auto Scaling (ESS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_autoscaling_scalingrule"
sidebar_current: "docs-Alibabacloudstack-autoscaling-scalingrule"
description: |- 
  编排弹性伸缩规则
---

# alibabacloudstack_autoscaling_scalingrule
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_autoscaling_rule` `alibabacloudstack_ess_scaling_rule`

使用Provider配置的凭证在指定的资源集下编排弹性伸缩规则。

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
  type              = "ingress"
  ip_protocol       = "tcp"
  nic_type          = "intranet"
  policy            = "accept"
  port_range        = "22/22"
  priority          = 1
  security_group_id = "${alibabacloudstack_ecs_securitygroup.default.id}"
  cidr_ip           = "172.16.0.0/24"
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
  default = "tf-testAccEssScalingRule-109553"
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

resource "alibabacloudstack_autoscaling_scalingrule" "default" {
  scaling_group_id = "${alibabacloudstack_ess_scaling_group.default.id}"
  adjustment_type = "TotalCapacity"
  adjustment_value = "1"
  cooldown = 0
  scaling_rule_name = "example-scaling-rule"
}
```

## 参数参考

支持以下参数：

* `scaling_group_id` - (必填, 变更时重建) - 伸缩组的ID。
* `adjustment_type` - (必填) - 伸缩规则的调整方式，取值范围：
  * `QuantityChangeInCapacity`：增加或减少指定数量的ECS实例。
  * `PercentChangeInCapacity`：增加或减少指定比例的ECS实例。
  * `TotalCapacity`：将当前伸缩组的ECS实例数量调整到指定数量。
* `adjustment_value` - (必填) - 伸缩规则的调整值。取值范围取决于`adjustment_type`：
  * `QuantityChangeInCapacity`：(-500, 0] U (0, 500]
  * `PercentChangeInCapacity`：[-100, 0] U [0, 10000]
  * `TotalCapacity`：[0, 1000]
* `scaling_rule_name` - (选填) - 伸缩规则的名称，2~64个英文或中文字符，以数字、大小字母或中文开头，可包含数字、下划线(_)、连字符(-)或点号(.)。同一用户账号同一地域同一伸缩组内唯一。如果没有指定该参数，默认为`ScalingRuleId`的值。
* `cooldown` - (选填) - 伸缩规则的冷却时间，仅适用于简单伸缩规则。取值范围：0~86400，单位：秒。默认值为0。

## 属性参考

除了上述所有参数外，还导出了以下属性：

* `id` - 伸缩规则的ID。
* `ari` - 伸缩规则的唯一标识符。
* `scaling_rule_aris` - 伸缩规则的唯一标识符列表。
* `scaling_rule_name` - 伸缩规则的名称，2~64个英文或中文字符，以数字、大小字母或中文开头，可包含数字、下划线(_)、连字符(-)或点号(.)。同一用户账号同一地域同一伸缩组内唯一。如果没有指定该参数，默认为`ScalingRuleId`的值。
```