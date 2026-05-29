---
subcategory: "Distributed Relational Database Service(DRDS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_drds_rds_instances"
sidebar_current: "docs-Alibabacloudstack-datasource-drds-rds-instances"
description: |-
  Provides a list of DRDS RDS instances owned by an Alibabacloudstack account.
---

# alibabacloudstack_drds_rds_instances

This data source provides a list of DRDS RDS instances associated with a DRDS instance in an Alibabacloudstack account.

> **Note:** This resource can also be referred to by the following aliases:
> - `apsarastack_drds_rds_instances`

## Example Usage

```hcl
data "alibabacloudstack_drds_rds_instances" "example" {
  drds_instance_id = "drds-xxxxxxxxxxxx"
  ids              = ["rm-xxxxxxxxxxxx"]
}

output "rds_instance_ids" {
  value = data.alibabacloudstack_drds_rds_instances.example.ids
}
```

## Argument Reference

The following arguments are supported:

* `drds_instance_id` - (Required) The ID of the DRDS instance. You can call the `DescribeDrdsInstances` API to query the DRDS instance ID.
* `ids` - (Optional) A list of RDS instance IDs. This can be used to filter the results to specific RDS instances.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `ids` - A list of RDS instance IDs.
* `rds_instances` - A list of DRDS RDS instances. Each instance contains the following attributes:
  * `rds_instance_id` - The ID of the RDS instance.
  * `db_instance_storage` - The storage capacity of the RDS instance, unit: GB.
  * `create_time` - The expiration time of the RDS instance.
