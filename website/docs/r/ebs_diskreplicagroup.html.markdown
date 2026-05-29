---
subcategory: "Elastic Block Storage"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ebs_diskreplicagroup"
sidebar_current: "docs-Alibabacloudstack-resource-ebs-diskreplicagroup"
description: |-
  Provides a ebs Diskreplicagroup resource.
---

# alibabacloudstack\_ebs\_diskreplicagroup

Provides a ebs Diskreplicagroup resource.

## Example Usage
```
variable "name" {
  default = "tf-testaccebs-diskreplicagroup43963"
}
variable "region" {
  default = ""
}

data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
  enable_details = true
}

resource "alibabacloudstack_ebs_diskreplicagroup" "default" {
  disk_replica_group_name = "${var.name}"
  description = "ebs_diskreplicagroup test"
  source_zone_id = "${data.alibabacloudstack_zones.default.zones[0].id}"
  source_region_id = "${var.region}"
  destination_zone_id = "${data.alibabacloudstack_zones.default.zones[1].id}"
  destination_region_id = "${var.region}"
  site = "production"
  rpo = "300"
}
```

## Argument Reference

The following arguments are supported:
  * `description` - (Optional) - The description of the consistent replication group.
  * `destination_region_id` - (Required) - The ID of the region to which the disaster recovery site belongs.
  * `destination_zone_id` - (Required) - The ID of the zone to which the disaster recovery site belongs.
  * `disk_replica_group_name` - (Optional) - Consistent replication group name.
  * `last_recover_point` - (Optional) - The time when the last asynchronous replication operation of the consistent replication group completed. This parameter provides the return value as a timestamp. Unit: seconds.
  * `rpo` - (Optional) - The RPO value of the consistent replication group.
  * `region_id` - (Optional) - The region ID of the consistent replication group, which is the same as the region of the production site.
  * `site` - (Optional) - Site information sources for replication pairs and consistent replication groups. Possible values:-production: production site.-backup: disaster recovery site.
  * `source_region_id` - (Required) - The ID of the region to which the production site belongs.
  * `source_zone_id` - (Required) - The ID of the zone to which the production site belongs.
  * `status` - (Optional) - The status of the consistent replication group. Possible values:-invalid: invalid. This state indicates that there is an exception to the replication pair in the consistent replication group.-creating: creating.-created: created.-create_failed: creation failed.-manual_syncing: in a single synchronization. If it is the first single synchronization, this status is also displayed in the synchronization.-syncing: synchronization. This state is the first time data is copied asynchronously between the master and slave disks.-normal: normal. When data replication is completed within the current cycle of asynchronous replication, it will be in this state.-stopping: stopping.-stopped: stopped.-stop_failed: Stop failed.-Failover: failover.-Failed: failover completed.-failover_failed: failover failed.-Reprotection: In reverse copy operation.-reprotect_failed: reverse replication failed.-deleting: deleting.-delete_failed: delete failed.-deleted: deleted.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
  * `last_recover_point` - The time when the last asynchronous replication operation of the consistent replication group completed. This parameter provides the return value as a timestamp. Unit: seconds.
  * `pair_ids` - List of replication pair IDs contained in a consistent replication group.
  * `pair_number` - The number of replication pairs contained in a consistent replication group.
  * `rpo` - The RPO value of the consistent replication group.
  * `replica_group_id` - The ID of the consistent replication group.
  * `status` - The status of the consistent replication group. Possible values:-invalid: invalid. This state indicates that there is an exception to the replication pair in the consistent replication group.-creating: creating.-created: created.-create_failed: creation failed.-manual_syncing: in a single synchronization. If it is the first single synchronization, this status is also displayed in the synchronization.-syncing: synchronization. This state is the first time data is copied asynchronously between the master and slave disks.-normal: normal. When data replication is completed within the current cycle of asynchronous replication, it will be in this state.-stopping: stopping.-stopped: stopped.-stop_failed: Stop failed.-Failover: failover.-Failed: failover completed.-failover_failed: failover failed.-Reprotection: In reverse copy operation.-reprotect_failed: reverse replication failed.-deleting: deleting.-delete_failed: delete failed.-deleted: deleted.
