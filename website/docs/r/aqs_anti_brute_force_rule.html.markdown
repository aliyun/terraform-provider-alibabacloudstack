---
subcategory: "Anti-Brute Force Security (AQS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_aqs_anti_brute_force_rule"
sidebar_current: "docs-Alibabacloudstack-aqs-anti_brute_force_rule"
description: |-
  Anti-brute force rule for Security Center (AQS)
---

# alibabacloudstack_aqs_anti_brute_force_rule

Configure Anti-brute force rule for Security Center (AQS) with the credentials configured in the provider for the specified resource set.

## Example Usage

### Basic Usage

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

## Argument Reference

The following arguments are supported:

* `name` - (Optional) Rule name. Used to identify the rule.
* `span` - (Optional) Time range for counting login failures in minutes. The time range within which login failures are counted. Valid values: `1`, `2`, `5`, `10`, `15`.
* `fail_count` - (Optional) Threshold of login failure count. When the login failure count reaches this value, the anti-brute force rule is triggered. Valid values: `2`, `3`, `4`, `5`, `10`, `50`, `80`, `100`.
* `forbidden_time` - (Optional) Login ban duration in minutes. The duration of login ban after the rule is triggered. Valid values: `5`, `15`, `30`, `60`, `120`, `360`, `720`, `1440`, `10080`, `52560000` (permanent).
* `default_rule` - (Optional) Whether it is a default rule. When the asset is not in any other rule, the default rule will be used. Default value is `false`.
* `instance_ids` - (Optional) List of associated ECS instance IDs. Specifies the ECS instances to which this rule applies (needs to be converted to UUID format).

-> **Note**: The `span`, `fail_count`, and `forbidden_time` parameters combine to form an anti-brute force rule, meaning that if an account fails to log in more than XX times within XX minutes, the account will be locked out for XX minutes.

## Attributes Reference

The following attributes are exported:

* `id` - Rule ID.
* `create_timestamp` - Rule creation timestamp (Unix timestamp in milliseconds).
* `enable_smart_rule` - Whether smart rule is enabled.
* `machine_count` - Number of associated machines.

## Import

Anti-brute force rule can be imported using the rule ID, e.g.

```
$ terraform import alibabacloudstack_aqs_anti_brute_force_rule.example 65778
```