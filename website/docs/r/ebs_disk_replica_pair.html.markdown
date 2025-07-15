---
subcategory: "EBS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ebs_diskreplicapair"
sidebar_current: "docs-Alibabacloudstack-ebs-diskreplicapair"
description: |-
  Provides a ebs Diskreplicapair resource.
---

# alibabacloudstack\_ebs\_diskreplicapair

Provides a ebs Diskreplicapair resource.

## Example Usage
```
variable "name" {
  default = "tf-testaccebs-diskreplicapair59477"
}

variable "region" {
  default = ""
}

data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
  enable_details = true
}



resource "alibabacloudstack_ecs_disk" "disk1" {
	availability_zone = "${data.alibabacloudstack_zones.default.zones[0].id}"
	size = "20"
	name = "${var.name}"
	category = "${data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0}"

	lifecycle {
		ignore_changes = ["tags"]
	}
}

resource "alibabacloudstack_ecs_disk" "disk2" {
	availability_zone = "${data.alibabacloudstack_zones.default.zones[1].id}"
	size = "20"
	name = "${var.name}"
	category = "${data.alibabacloudstack_zones.default.zones.1.available_disk_categories.0}"
	lifecycle {
		ignore_changes = ["tags"]
	}
}




resource "alibabacloudstack_ebs_diskreplicapair" "default" {
  description = "ebs_diskreplicapair test"
  source_zone_id = "${data.alibabacloudstack_zones.default.zones[0].id}"
  source_region_id = "${var.region}"
  source_disk_id = "${alibabacloudstack_ecs_disk.disk1.id}"
  destination_zone_id = "${data.alibabacloudstack_zones.default.zones[1].id}"
  destination_disk_id = "${alibabacloudstack_ecs_disk.disk2.id}"
  disk_replica_pair_name = "${var.name}"
  destination_region_id = "${var.region}"
  rpo = "300"
}
```

## Argument Reference

The following arguments are supported:
  * `description` - (Optional) - The description of the asynchronous replication relationship. 2 to 256 English or Chinese characters in length and cannot start with' http:// 'or' https.
  * `destination_disk_id` - (Required) - The ID of the standby disk.
  * `destination_region_id` - (Required) - The ID of the region to which the disaster recovery site belongs.
  * `destination_zone_id` - (Required) - The ID of the zone to which the disaster recovery site belongs.
  * `source_disk_id` - (Required) - The ID of the disk to be replicated.
  * `disk_replica_pair_name` - (Optional) - The name of the asynchronous replication relationship. The length must be 2 to 128 characters in length and must start with a letter or Chinese name. It cannot start with http:// or https. It can contain Chinese, English, numbers, half-width colons (:), underscores (_), half-width periods (.), or dashes (-).
  * `last_recover_point` - (Optional) - The time when data was last replicated from the primary disk to the secondary disk in the replication pair. 
  * `one_shot` - (Optional) - Whether to synchronize immediately. Value range:-true: Start data synchronization immediately.-false: Data Synchronization starts after the RPO time period.Default value: false.
  * `rpo` - (Optional) - The Rpo of the asynchronous replication relationship. The unit is seconds. Value range: 300 to 86400.Default value: 300.
  * `source_region_id` - (Required) - The ID of the region to which the production site belongs.
  * `replica_group_id` - (Optional) - The ID of the replication group.
  * `source_zone_id` - (Required) - The ID of the zone to which the production site belongs.
  * `status` - (Optional) - The status of the resource

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
  * `last_recover_point` - The time when data was last replicated from the primary disk to the secondary disk in the replication pair. 
  * `replica_pair_id` - The first ID of the resource
  * `create_time` - The creation time of the resource
  * `status` - The status of the resource
