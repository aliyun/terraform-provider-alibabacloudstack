---
subcategory: "Elastic Block Storage"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ebs_diskreplicagroups"
sidebar_current: "docs-Alibabacloudstack-datasource-ebs-diskreplicagroups"
description: |-
  Provides a list of ebs diskreplicagroups owned by an alibabacloudstack account.
---

# alibabacloudstack\_ebs\_diskreplicagroups

This data source provides a list of ebs diskreplicagroups in an alibabacloudstack account according to the specified filters.

## Example Usage
```
variable "name" {
	default = "tf-testAccEbsDiskReplicaGroupsDataSource-8901034"
}

variable "region_id" {
  default = ""
}

data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
  enable_details = true
}

resource "alibabacloudstack_ebs_diskreplicagroup" "default" {
    disk_replica_group_name = "${var.name}"
    description = "${var.name}"
    destination_region_id = "${var.region_id}"
    destination_zone_id ="${data.alibabacloudstack_zones.default.zones[1].id}"
    site = "production"
    source_region_id = "${var.region_id}"
    source_zone_id = "${data.alibabacloudstack_zones.default.zones[0].id}"
}
 
data "alibabacloudstack_ebs_diskreplicagroups" "default" {
  description_regex = "${alibabacloudstack_ebs_diskreplicagroup.default.description}"
}
```

## Argument Reference

The following arguments are supported:
  * `ids` - (Optional) - A list of Diskreplicagroup IDs.
  * `name_regex` - (Optional) - A name Regex of Diskreplicagroup.
  * `description_regex` - (Optional) - A description Regex of Diskreplicagroup.
  * `site` - (Optional) - Site information sources for replication pairs and consistent replication groups. Possible values:-production: production site.-backup: disaster recovery site.
  * `source_region_id` - (Optional) - The ID of the region to which the production site belongs.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
  * `disk_replica_groups` - A list of disk replica groups.
    * `id` - The ID of the consistent replication group.
    * `description` - The description of the consistent replication group.
    * `destination_region_id` - The ID of the region to which the disaster recovery site belongs.
    * `destination_zone_id` - The ID of the zone to which the disaster recovery site belongs.
    * `disk_replica_group_name` - Consistent replication group name.
    * `last_recover_point` - The time when the last asynchronous replication operation of the consistent replication group completed. This parameter provides the return value as a timestamp. Unit: seconds.
    * `one_shot` - Whether to synchronize immediately. Value range:-true: Start data synchronization immediately.-false: Data Synchronization starts after the RPO time period.Default value: false.
    * `pair_ids` - List of replication pair IDs contained in a consistent replication group.
    * `pair_number` - The number of replication pairs contained in a consistent replication group.
    * `rpo` - The RPO value of the consistent replication group.
    * `replica_group_id` - The ID of the consistent replication group.
    * `site` - Site information sources for replication pairs and consistent replication groups. Possible values:-production: production site.-backup: disaster recovery site.
    * `source_region_id` - The ID of the region to which the production site belongs.
    * `source_zone_id` - The ID of the zone to which the production site belongs.
    * `status` - The status of the consistent replication group. Possible values:-invalid: invalid. This state indicates that there is an exception to the replication pair in the consistent replication group.-creating: creating.-created: created.-create_failed: creation failed.-manual_syncing: in a single synchronization. If it is the first single synchronization, this status is also displayed in the synchronization.-syncing: synchronization. This state is the first time data is copied asynchronously between the master and slave disks.-normal: normal. When data replication is completed within the current cycle of asynchronous replication, it will be in this state.-stopping: stopping.-stopped: stopped.-stop_failed: Stop failed.-Failover: failover.-Failed: failover completed.-failover_failed: failover failed.-Reprotection: In reverse copy operation.-reprotect_failed: reverse replication failed.-deleting: deleting.-delete_failed: delete failed.-deleted: deleted.
    * `tags` - The tag of the resource
