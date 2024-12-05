---
subcategory: "Table Store (OTS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ots_instances"
sidebar_current: "docs-alibabacloudstack-datasource-ots-instances"
description: |-
    Provides a list of ots instances to the user.
---

# alibabacloudstack\_ots\_instances

This data source provides the ots instances of the current Alibaba Cloud user.


## Example Usage

```terraform
data "alibabacloudstack_ots_instances" "instances_ds" {
  output_file = "instances.txt"
}

output "first_instance_id" {
  value = "${data.alibabacloudstack_ots_instances.instances_ds.instances.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew) A list of instance IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by instance name.
* `tags` - (Optional) A map of tags assigned to the instance. It must be in the format:
  ```terraform
  data "alibabacloudstack_ots_instances" "instances_ds" {
    tags = {
      tagKey1 = "tagValue1",
      tagKey2 = "tagValue2"
    }
  }
  ```
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of instance IDs.
* `names` - A list of instance names.
* `instances` - A list of instances. Each element contains the following attributes:
  * `id` - ID of the instance.
  * `name` - Instance name.
  * `status` - Instance status. Possible values: `Running`, `Disabled`, `Deleting`.
  * `write_capacity` - Reserve write throughput. The unit is CU. Only high-performance instances have this return value.
  * `read_capacity` - Reserve read throughput. The unit is CU. Only high-performance instances have this return value.
  * `cluster_type` - The cluster type of the instance. Possible values: `SSD`, `HYBRID`.
  * `create_time` - The create time of the instance.
  * `user_id` - The user id of the instance.
  * `network` - The network type of the instance. Possible values: `NORMAL`, `VPC`, `VPC_CONSOLE`.
  * `description` - The description of the instance.
  * `entity_quota` - The instance quota which indicating the maximum number of tables.
  * `tags` - The tags of the instance.
	
