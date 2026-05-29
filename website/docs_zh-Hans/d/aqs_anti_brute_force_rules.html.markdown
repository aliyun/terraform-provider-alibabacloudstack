---
subcategory: "防暴力破解安全服务"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_aqs_anti_brute_force_rule"
sidebar_current: "docs-Alibabacloudstack-datasource-aqs-anti_brute_force_rule"
description: |-
  查询安骑士主机防暴力破解规则
---

# alibabacloudstack_aqs_anti_brute_force_rule

安骑士主机防暴力破解规则查询。

## 示例用法

```hcl

variable "name" {
  default = "tf-testacc16314"
}

data "alibabacloudstack_zones" "default" {
  provider                    = alibabacloudstack-common
  available_resource_creation = "VSwitch"
  enable_details              = true
}

resource "alibabacloudstack_vpc_vpc" "default" {
  provider   = alibabacloudstack-common
  vpc_name   = "${var.name}_vpc"
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vpc_vswitch" "default" {
  provider   = alibabacloudstack-common
  name       = "${var.name}_vsw"
  vpc_id     = alibabacloudstack_vpc_vpc.default.id
  cidr_block = "172.16.0.0/24"
  zone_id    = data.alibabacloudstack_zones.default.zones.0.id
}

resource "alibabacloudstack_ecs_securitygroup" "default" {
  provider = alibabacloudstack-common
  name     = "${var.name}_sg"
  vpc_id   = alibabacloudstack_vpc_vpc.default.id
}

data "alibabacloudstack_images" "default" {
  provider    = alibabacloudstack-common
  name_regex  = "^ubuntu_"
  most_recent = true
  owners      = "system"
}

data "alibabacloudstack_instance_types" "all" {
  provider          = alibabacloudstack-common
  sorted_by         = "Memory"
  availability_zone = data.alibabacloudstack_zones.default.zones[0].id
}

resource "alibabacloudstack_ecs_instance" "default" {
  provider                      = alibabacloudstack-common
  system_disk_category          = data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0
  instance_name                 = var.name
  user_data                     = "I_am_user_data"
  security_groups               = ["${alibabacloudstack_ecs_securitygroup.default.id}"]
  vswitch_id                    = alibabacloudstack_vpc_vswitch.default.id
  image_id                      = data.alibabacloudstack_images.default.images.0.id
  security_enhancement_strategy = "Active"
  instance_type                 = data.alibabacloudstack_instance_types.all.instance_types.0.id
  availability_zone             = data.alibabacloudstack_zones.default.zones[0].id
}

resource "alibabacloudstack_aqs_anti_brute_force_rule" "default" {
  provider       = alibabacloudstack-common
  name           = "${var.name}_rule"
  span           = 10
  fail_count     = 80
  forbidden_time = 360
  default_rule   = true
  instance_ids = [
    alibabacloudstack_ecs_instance.default.id,
  ]
}

data "alibabacloudstack_aqs_anti_brute_force_rules" "default" {
  ids = ["${alibabacloudstack_aqs_anti_brute_force_rule.default.id}"]
}

```

## 参数说明
以下参数用于过滤结果：

* `ids` (可选)：规则ID列表，用于过滤结果。
* `name_regex` (可选)：用于按规则名称过滤结果的正则表达式。

## 属性说明
以下属性被导出：

* `id` (字符串)：规则的唯一标识符。
* `create_timestamp` (整数)：规则创建时间戳。
* `default_rule` (布尔值)：是否为默认规则。
* `enable_smart_rule` (布尔值)：是否启用智能规则。
* `fail_count` (整数)：失败尝试次数阈值。
* `forbidden_time` (整数)：禁止登录时间（秒）。
* `instance_ids` (集合)：关联的实例ID列表。
* `machine_count` (整数)：关联的机器数量。
* `name` (字符串)：规则名称。
* `span` (整数)：时间窗口（分钟）。