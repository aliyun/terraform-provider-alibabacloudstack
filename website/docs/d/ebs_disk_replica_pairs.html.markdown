---
subcategory: "Elastic Block Storage"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ebs_diskreplicapairs"
sidebar_current: "docs-Alibabacloudstack-datasource-ebs-diskreplicapairs"
description: |-
  Provides a list of ebs diskreplicapairs owned by an alibabacloudstack account.
---

# alibabacloudstack\_ebs\_diskreplicapairs

This data source provides a list of ebs diskreplicapairs in an alibabacloudstack account according to the specified filters.

## Example Usage
```
variable "name" {
	default = "tf-testAccEbsDiskReplicaPairsDataSource-2021428"
}

variable "region" {
  default = ""
}

resource "alibabacloudstack_ecs_disk" "disk1" {
	availability_zone = "${data.alibabacloudstack_zones.default.zones[0].id}"
	size = "20"
	name = "${var.name}"
	category = "${data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0}"
}

resource "alibabacloudstack_ecs_disk" "disk2" {
	availability_zone = "${data.alibabacloudstack_zones.default.zones[1].id}"
	size = "20"
	name = "${var.name}"
	category = "${data.alibabacloudstack_zones.default.zones.1.available_disk_categories.0}"
}
	
resource "alibabacloudstack_ebs_diskreplicapair" "default" {
	disk_replica_pair_name = "${var.name}"
	description =            "${var.name}"
	source_zone_id =         "${data.alibabacloudstack_zones.default.zones[0].id}"
	source_region_id =       "${var.region}"
	source_disk_id =         "${alibabacloudstack_ecs_disk.disk1.id}$"
	destination_zone_id =    "${data.alibabacloudstack_zones.default.zones[1].id}"
	destination_region_id =  "${var.region}"
	destination_disk_id =    "${alibabacloudstack_ecs_disk.disk2.id}$"
	rpo =                    300
}
 

data "alibabacloudstack_ebs_diskreplicapairs" "default" {
  description_regex = "${alibabacloudstack_ebs_diskreplicapair.default.description}"
}
```

## Argument Reference

The following arguments are supported:
  * `name_regex` - (Optional) - A regex string to filter results by diskreplicapair name.
  * `description_regex` - (Optional) - A regex string to filter results by diskreplicapair description.
  * `ids` - (Optional) - A list of diskreplicapair ids to filter results.
  * `source_region_id` - (Optional) - The ID of the source region.
  * `replica_group_id` - (Optional) - The ID of the replication group.
  * `max_results` - (Optional) - MaxItems specifies the maximum number of records returned by this request.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
  * `disk_replica_pairs` - A list of diskreplicapairs. Each element contains the following attributes:
    * `id` - The ID of the diskreplicapair.
    * `bandwidth` - The bandwidth for asynchronous data replication between cloud disks. The unit is Kbps. Value range:-10240 Kbps: equal to 10 Mbps.-20480 Kbps: equal to 20 Mbps.-51200 Kbps: equal to 50 Mbps.-102400 Kbps: equal to 100 Mbps.Default value: 10240.This parameter cannot be specified when the ChargeType value is PayAsYouGo The system value is 0, which indicates that the disk is dynamically allocated according to data write changes during asynchronous replication.
    * `description` - The description of the asynchronous replication relationship. 2 to 256 English or Chinese characters in length and cannot start with' http:'or' https.
    * `destination_disk_id` - The ID of the standby disk.
    * `destination_region_id` - The ID of the region to which the disaster recovery site belongs.
    * `destination_zone_id` - The ID of the zone to which the disaster recovery site belongs.
    * `source_disk_id` - The ID of the zone to which the production site belongs.
    * `disk_replica_pair_name` - The name of the asynchronous replication relationship. The length must be 2 to 128 characters in length and must start with a letter or Chinese name. It cannot start with http:or https. It can contain Chinese, English, numbers, half-width colons (:), underscores (_), half-width periods (.), or dashes (-).
    * `last_recover_point` - The time when data was last replicated from the primary disk to the secondary disk in the replication pair. 
    * `one_shot` - Whether to synchronize immediately. Value range:-true: Start data synchronization immediately.-false: Data Synchronization starts after the RPO time period.Default value: false.
    * `rpo` - The Rpo of the asynchronous replication relationship. The unit is seconds. Value range: 300 to 86400.Default value: 300.
    * `source_region_id` - The ID of the region to which the production site belongs.
    * `replica_group_id` - The ID of the replication group.
    * `replica_pair_id` - The first ID of the resource
    * `source_zone_id` - The ID of the zone to which the production site belongs.
    * `status` - The status of the resource
    * `tags` - The tag of the resource
