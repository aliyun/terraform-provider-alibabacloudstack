---
subcategory: "Auto Scaling (ESS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ess_scaling_group"
sidebar_current: "docs-alibabacloudstack-resource-ess-scaling-group"
description: |-
  编排ESS伸缩组资源
---

# alibabacloudstack_ess_scaling_group

使用Provider配置的凭证在指定的资源集下编排ESS伸缩组资源，该资源是一个具有相同应用场景的ECS实例集合。

它定义了伸缩组中ECS实例的最大和最小数量，以及它们关联的负载均衡实例、RDS实例和其他属性。

-> **注意：** 您可以通过指定参数`vswitch_ids`在VPC网络中启动一个ESS伸缩组。

## 示例用法

```
variable "name" {
  default = "essscalinggroupconfig"
}

data "alibabacloudstack_zones" "default" {
  available_disk_category     = "cloud_efficiency"
  available_resource_creation = "VSwitch"
}

data "alibabacloudstack_instance_types" "default" {
  availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
  cpu_core_count    = 2
  memory_size       = 4
}

data "alibabacloudstack_images" "default" {
  name_regex  = "^ubuntu_18.*64"
  most_recent = true
  owners      = "system"
}

resource "alibabacloudstack_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vswitch" "default" {
  vpc_id            = "${alibabacloudstack_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
  name              = "${var.name}"
}

resource "alibabacloudstack_security_group" "default" {
  name   = "${var.name}"
  vpc_id = "${alibabacloudstack_vpc.default.id}"
}

resource "alibabacloudstack_security_group_rule" "default" {
  type              = "ingress"
  ip_protocol       = "tcp"
  nic_type          = "intranet"
  policy            = "accept"
  port_range        = "22/22"
  priority          = 1
  security_group_id = "${alibabacloudstack_security_group.default.id}"
  cidr_ip           = "172.16.0.0/24"
}

resource "alibabacloudstack_vswitch" "default2" {
  vpc_id            = "${alibabacloudstack_vpc.default.id}"
  cidr_block        = "172.16.1.0/24"
  availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
  name              = "${var.name}-bar"
}

resource "alibabacloudstack_ess_scaling_group" "default" {
  min_size           = 1
  max_size           = 1
  scaling_group_name = "${var.name}"
  default_cooldown   = 20
  vswitch_ids        = ["${alibabacloudstack_vswitch.default.id}", "${alibabacloudstack_vswitch.default2.id}"]
  removal_policies   = ["OldestInstance", "NewestInstance"]
}
```

## 参数参考

以下是支持的参数：

* `min_size` - (必填) 伸缩组中的ECS实例最小数量。取值范围：[0, 100]。
* `max_size` - (必填) 伸缩组中的ECS实例最大数量。取值范围：[0, 100]。
* `scaling_group_name` - (可选) 伸缩组显示名称，必须包含2-40个字符(英文或中文)，以数字、英文字母或中文字符开头，可以包含数字、下划线 `_`、连字符 `-` 和点 `.`。如果未指定此参数，默认值为ScalingGroupId。
* `default_cooldown` - (可选) 伸缩组的默认冷却时间(秒)。取值范围：[0, 86400]。默认值为300s。
* `vswitch_ids` - (可选) 要启动ECS实例的虚拟交换机ID列表。
* `removal_policies` - (可选) RemovalPolicy用于选择从伸缩组中移除的ECS实例，当存在多个候选实例时。可选值：
    - OldestInstance: 移除最早附加到伸缩组的ECS实例。
    - NewestInstance: 移除最新附加到伸缩组的ECS实例。
    - OldestScalingConfiguration: 移除具有最旧伸缩配置的ECS实例。
    - 默认值：OldestScalingConfiguration 和 OldestInstance。最多可以输入两个移除策略。
* `db_instance_ids` - (可选) 如果在伸缩组中指定了RDS实例，伸缩组会自动将ECS实例的内网IP地址添加到RDS访问白名单中。
    - 指定的RDS实例必须处于运行状态。
    - 指定的RDS实例的白名单必须有空间容纳更多的IP地址。
* `loadbalancer_ids` - (可选) 如果在伸缩组中指定了负载均衡实例，伸缩组会自动将ECS实例附加到负载均衡实例。
    - 负载均衡实例必须已启用。
    - 每个负载均衡实例至少必须配置一个监听器，并且其健康检查必须开启。否则，创建将失败(可能需要添加一个`depends_on`参数，目标是您的`alibabacloudstack_slb_listener`，以确保监听器及其健康检查配置就绪后再创建伸缩组)。
    - 附加VPC类型ECS实例的负载均衡实例不能附加到伸缩组。
    - 附加到负载均衡实例的ECS实例的默认权重为50。
* `multi_az_policy` - (可选, 变更时重建) 多可用区伸缩组的ECS实例扩缩策略。可选值：PRIORITY（优先级）、BALANCE（平衡）或 COST_OPTIMIZED（成本优化）。

-> **注意：** 当分离负载均衡器时，组中的实例将从负载均衡器的`默认服务器组`中移除；相反，当附加负载均衡器时，组中的实例将被添加到负载均衡器的`默认服务器组`。

-> **注意：** 当分离数据库实例时，组中实例的私有IP将从数据库实例的`白名单`中移除；相反，当附加数据库实例时，组中实例的私有IP将被添加到数据库实例的`白名单`中。


## 属性参考

以下属性将导出：

* `id` - 伸缩组ID。
* `min_size` - ECS实例的最小数量。
* `max_size` - ECS实例的最大数量。
* `scaling_group_name` - 伸缩组名称。
* `default_cooldown` - 伸缩组的默认冷却时间。
* `removal_policies` - 用于选择从伸缩组中移除ECS实例的移除策略。
* `db_instance_ids` - ECS实例附加到的数据库实例ID。
* `loadbalancer_ids` - ECS实例附加到的负载均衡实例ID。
* `vswitch_ids` - 启动ECS实例的交换机ID。