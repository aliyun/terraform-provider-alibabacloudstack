---
subcategory: "Auto Scaling (ESS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_autoscaling_lifecyclehook"
sidebar_current: "docs-Alibabacloudstack-autoscaling-lifecyclehook"
description: |-
  Provides a autoscaling Lifecyclehook resource.
---

# alibabacloudstack_autoscaling_lifecyclehook
-> **NOTE:** Alias name has: `alibabacloudstack_ess_lifecycle_hook`

Provides a autoscaling Lifecyclehook resource.

## Example Usage
```
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
  //name_regex  = "arm_centos_7_6_20G_20211110.raw"
  //name_regex  = "^arm_centos_7"
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
		default = "tf-testAccEssLifecycleHook-81869"
	}
	
	resource "alibabacloudstack_vswitch" "default2" {
		  vpc_id = "${alibabacloudstack_vpc.default.id}"
		  cidr_block = "172.16.1.0/24"
		  availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
		  name = "${var.name}"
	}
	
	resource "alibabacloudstack_ess_scaling_group" "default" {
		min_size = 1
		max_size = 1
		scaling_group_name = "${var.name}"
		removal_policies = ["OldestInstance", "NewestInstance"]
		vswitch_ids = ["${alibabacloudstack_vswitch.default.id}","${alibabacloudstack_vswitch.default2.id}"]
	}
	
	resource "alibabacloudstack_ess_lifecycle_hook" "default"{
		scaling_group_id = "${alibabacloudstack_ess_scaling_group.default.id}"
		name = "${var.name}"
		lifecycle_transition = "SCALE_OUT"
		notification_metadata = "helloworld"
	}
```

## Argument Reference

The following arguments are supported:
  * `scaling_group_id` - (Required, ForceNew) - Telescopic group ID.
  * `name` - (Optional, ForceNew) - The name of the life cycle hook.
  * `lifecycle_hook_name` - (Optional, ForceNew) - Linked to the life cycle a name.
  * `lifecycle_transition` - (Required) - Life cycle linked to the corresponding expansion and contraction type of activity.
  * `heartbeat_timeout` - (Optional) - Life cycle linked to set the wait time for the FLEX group activities, will be the next step after the Timeout waiting for State action.
  * `default_result` - (Optional) - When scaling group contractile activity of elastic ( scale_in ) and multiple linked to the life cycle is triggered, defaultresult to abandon linked to the life cycle of trigger the end of the wait state, will end early will wait for the other early States.In other cases, the next
  
## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
  * `name` - The name of the life cycle hook.
  * `lifecycle_hook_name` - Linked to the life cycle a name.
  * `notification_arn` - Life cycle linked to the notification object identifier.
  * `notification_metadata` - Fixed string information for the expansion of the activities of the wait state.
  * `heartbeat_timeout` - Life cycle linked to set the wait time for the FLEX group activities, will be the next step after the Timeout waiting for State action.
  * `default_result` - When scaling group contractile activity of elastic ( scale_in ) and multiple linked to the life cycle is triggered, defaultresult to abandon linked to the life cycle of trigger the end of the wait state, will end early will wait for the other early States.In other cases, the next