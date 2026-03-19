---
subcategory: "Anti-Brute Force Security (AQS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_aqs_anti_brute_force_rule"
sidebar_current: "docs-Alibabacloudstack-aqs-anti_brute_force_rule"
description: |-
  安骑士主机防暴力破解规则
---

# alibabacloudstack_aqs_anti_brute_force_rule

使用Provider配置的凭证在指定的资源集配置安骑士主机防暴力破解规则。

## 示例用法

### 基础用法

```hcl

variable "name" {
  default = "tf-testacc9908"
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
resource "alibabacloudstack_ecs_instance" "default0" {
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

resource "alibabacloudstack_ecs_instance" "default1" {
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
  default_rule = "true"
  instance_ids = [
    "${alibabacloudstack_ecs_instance.default0.id}",
    "${alibabacloudstack_ecs_instance.default1.id}"
  ]
  name           = var.name
  span           = "10"
  fail_count     = "80"
  forbidden_time = "360"
}
```

## 参数说明

支持以下参数：

* `default_rule` - (可选) 是否默认规则。默认值为`false`。
* `fail_count` - (可选) 登录失败次数阈值。当登录失败次数达到此值时，触发防暴力破解规则。
* `forbidden_time` - (可选) 禁止登录时间，单位分钟。触发规则后，禁止登录的时长。
* `instance_ids` - (可选) 关联的ECS实例ID列表。指定应用此规则的ECS实例（需转换为UUID格式）。
* `name` - (可选) 规则名称。用于标识规则。
* `span` - (可选) 登录失败次数统计时间范围，单位分钟。在此时间范围内统计登录失败次数。

## 属性说明

导出以下属性：

* `id` - 规则ID。
* `create_timestamp` - 规则创建时间戳（Unix时间戳，毫秒）。
* `enable_smart_rule` - 是否启用智能规则。
* `machine_count` - 关联的机器数量。