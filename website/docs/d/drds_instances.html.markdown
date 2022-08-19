---
subcategory: "Distributed Relational Database Service (DRDS)"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_drds_instances"
sidebar_current: "docs-apsarastack-drds-instances"
description: |-
  Provides a collection of DRDS instances according to the specified filters.
---

# apsarastack_drds_instance

 The `apsarastack_drds_instance` data source provides a collection of DRDS instances available in Alibaba Cloud account.
Filters support regular expression for the instance name, searches by tags, and other filters which are listed below.

## Example Usage

 ```
data "apsarastack_drds_instances" "drds_instances_ds" {
  name_regex = "drds-\\d+"
  ids        = ["drdsabc123456"]
}
output "first_db_instance_id" {
  value = "${data.apsarastack_drds_instances.drds_instances_ds.instances.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional, Deprecated) A regex string to filter results by instance description. It is deprecated since v1.91.0 and will be removed in a future release, please use 'description_regex' instead.
* `description_regex` - (Optional) A regex string to filter results by instance description.
* `ids` - (Optional) A list of DRDS instance IDs.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

 * `ids` - A list of DRDS instance IDs.
 * `descriptions` - A list of DRDS descriptions. 
 * `instances` - A list of DRDS instances.
   * `id` - The ID of the DRDS instance.
   * `description` - The DRDS instance description.
   * `name` - The name of the RDS instance.
   * `status` - Status of the instance.
   * `type` - The DRDS Instance type.
   * `create_time` - Creation time of the instance.
   * `network_type` - `Classic` for public classic network or `VPC` for private network.
   * `zone_id` - Zone ID the instance belongs to.
   * `version` - The DRDS Instance version.
   * `ids` - A list of DRDS instance IDs.
