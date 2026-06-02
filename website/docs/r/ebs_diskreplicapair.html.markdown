---
subcategory: "Elastic Block Storage"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ebs_diskreplicapair"
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
  * `source_region_id` - (Required) The ID of the region to which the primary disk belongs.
  * `source_zone_id` - (Required) The ID of the zone to which the primary disk belongs.
  * `source_disk_id` - (Required) The ID of the primary disk.
  * `destination_region_id` - (Required) The ID of the region to which the secondary disk belongs.
  * `destination_zone_id` - (Required) The ID of the zone to which the secondary disk belongs.
  * `destination_disk_id` - (Required) The ID of the secondary disk.
  * `disk_replica_pair_name` - (Optional) The name of the async replication pair. It must be 2 to 128 characters in length, start with a letter, and cannot start with http:// or https://. It can contain letters, digits, colons (:), underscores (_), periods (.), and hyphens (-).
  * `description` - (Optional) The description of the async replication pair. It must be 2 to 256 characters in length and cannot start with http:// or https://.
  * `rpo` - (Optional) The Recovery Point Objective (RPO) of the async replication pair. Unit: seconds. Valid values: 300 to 86400. Default value: 300.
  * `one_shot` - (Optional) Specifies whether to perform immediate synchronization. Valid values: `true` (start data synchronization immediately), `false` (start data synchronization after the RPO period). Default value: `false`.
  * `replica_group_id` - (Optional) The ID of the replication group to which the pair belongs.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
  * `id` - The ID of the async replication pair.
  * `replica_pair_id` - The ID of the async replication pair.
  * `create_time` - The creation time of the async replication pair in UTC format (e.g., 2006-01-02T15:04:05-07:00).
  * `last_recover_point` - The timestamp when data was last replicated from the primary disk to the secondary disk. Unit: seconds.
  * `status` - The status of the async replication pair. Valid values: `creating`, `created`, `syncing`, `normal`, `stopped`, `failovered`, `deleting`, `failed`, etc.

## Import

EBS Disk Replica Pair can be imported using the replica pair ID, e.g.

```
$ terraform import alibabacloudstack_ebs_diskreplicapair.example rp-xxxxxxxxx
```
