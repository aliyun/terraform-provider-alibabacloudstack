---
subcategory: "GraphDatabase(GPDB)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_gpdb_dbinstances"
sidebar_current: "docs-Alibabacloudstack-datasource-gpdb-dbinstances"
description: |- 
  Provides a list of gpdb dbinstances owned by an alibabacloudstack account.
---

# alibabacloudstack_gpdb_dbinstances
-> **NOTE:** Alias name has: `alibabacloudstack_gpdb_instances`

This data source provides a list of GPDB DBInstances in an AlibabacloudStack account according to the specified filters.

## Example Usage

```hcl
data "alibabacloudstack_gpdb_dbinstances" "example" {
  name_regex        = "gp-.+\\d+"
  availability_zone = "cn-beijing-c"
  vswitch_id        = "vsw-1234567890abcdefg"
  output_file       = "dbinstances.txt"
}

output "dbinstance_id" {
  value = "${data.alibabacloudstack_gpdb_dbinstances.example.instances.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional) A regex string to filter GPDB DBInstances by name.
* `availability_zone` - (Optional) The Availability Zone of the instance.
* `vswitch_id` - (Optional) Used to retrieve instances belonging to a specified VSwitch resource.
* `ids` - (Optional) A list of GPDB DBInstance IDs.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `names` - A list of names of matching GPDB DBInstances.
* `ids` - A list of IDs of matching GPDB DBInstances.
* `instances` - A list of GPDB DBInstances. Each element contains the following attributes:
  * `id` - The ID of the GPDB DBInstance.
  * `description` - The description of the GPDB DBInstance.
  * `region_id` - The region ID where the GPDB DBInstance is located.
  * `availability_zone` - The Availability Zone of the GPDB DBInstance.
  * `creation_time` - The creation time of the GPDB DBInstance in UTC format (YYYY-MM-DDThh:mm:ssZ).
  * `status` - The current status of the GPDB DBInstance.
  * `engine` - The database engine type. Supported value is `gpdb`.
  * `engine_version` - The version of the database engine. Supported values include `6.0` and `7.0`.
  * `instance_class` - The class of the GPDB DBInstance.
  * `instance_group_count` - The number of groups in the GPDB DBInstance.
  * `instance_network_type` - The network type of the GPDB DBInstance. Supported value is `VPC`.
  * `charge_type` - The billing method of the GPDB DBInstance. Possible values are `PrePaid` (subscription) and `PostPaid` (pay-as-you-go).