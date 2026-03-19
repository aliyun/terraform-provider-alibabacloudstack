---
subcategory: "Anti-Brute Force Security (AQS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_aqs_anti_brute_force_rule"
sidebar_current: "docs-Alibabacloudstack-datasource-aqs-anti_brute_force_rule"
description: |-
  Query Anti-Brute Force Security rules for servers
---

# alibabacloudstack_aqs_anti_brute_force_rule

Query Anti-Brute Force Security rules for servers.

## Example Usage

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

## Argument Reference

The following arguments are supported:

* `ids` (Optional): A list of rule IDs used to filter results.
* `name_regex` (Optional): A regular expression used to filter results by rule name.

## Attributes Reference

The following attributes are exported:

* `id` (String): The unique identifier of the rule.
* `create_timestamp` (Integer): The timestamp when the rule was created.
* `default_rule` (Boolean): Whether it is a default rule.
* `enable_smart_rule` (Boolean): Whether the smart rule is enabled.
* `fail_count` (Integer): The threshold for the number of failed attempts.
* `forbidden_time` (Integer): The login ban duration in seconds.
* `instance_ids` (Set): A list of associated instance IDs.
* `machine_count` (Integer): The number of associated machines.
* `name` (String): The name of the rule.
* `span` (Integer): The time window in minutes.