---
subcategory: "DRDS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_drds_instances"
sidebar_current: "docs-Alibabacloudstack-datasource-drds-instances"
description: |- 
  Provides a list of DRDS instances owned by an Alibabacloudstack account.
---

# alibabacloudstack_drds_instances

This data source provides a list of DRDS instances in an Alibabacloudstack account according to the specified filters.

## Example Usage

```hcl
data "alibabacloudstack_drds_instances" "drds_instances_ds" {
  name_regex      = "drds-\\d+"
  ids             = ["drdsabc123456"]
  description_regex = "example-description.*"
}

output "first_db_instance_id" {
  value = "${data.alibabacloudstack_drds_instances.drds_instances_ds.instances.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional) A regex string to filter results by instance name.
* `description_regex` - (Optional) A regex string to filter results by instance description.
* `ids` - (Optional) A list of DRDS instance IDs. This can be used to limit the results to specific instances.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `ids` - A list of DRDS instance IDs.
* `descriptions` - A list of DRDS instance descriptions.
* `instances` - A list of DRDS instances. Each instance contains the following attributes:
  * `id` - The ID of the DRDS instance.
  * `description` - The description of the DRDS instance.
  * `status` - The status of the DRDS instance.
  * `type` - The type of the DRDS instance.
  * `create_time` - The creation time of the DRDS instance.
  * `network_type` - The network type of the DRDS instance. It can be `Classic` for public classic network or `VPC` for private network.
  * `zone_id` - The zone ID where the DRDS instance is located.
  * `version` - The version of the DRDS instance.