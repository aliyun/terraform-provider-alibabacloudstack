---
subcategory: "RocketMQ"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ons_instances"
sidebar_current: "docs-alibabacloudstack-datasource-ons-instances"
description: |-
    Provides a list of ons instances available to the user.
---

# alibabacloudstack\_ons\_instances

This data source provides a list of ONS Instances in an Apsara Stack Cloud account according to the specified filters.


## Example Usage

```
variable "name" {
  default = "onsInstanceDatasourceName"
}

resource "alibabacloudstack_ons_instance" "default" {
  tps_receive_max = 500
  tps_send_max = 500
  topic_capacity = 50
  cluster = "cluster1"
  independent_naming = "true"
  name = var.name
  remark = "Ons Instance"
}

data "alibabacloudstack_ons_instances" "instances_ds" {
  name_regex = alibabacloudstack_ons_instance.inst.name
  output_file = "instances.txt"
}

output "first_instance_id" {
  value = data.alibabacloudstack_ons_instances.instances_ds.*
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of instance IDs to filter results.
* `name_regex` - (Optional) A regex string to filter results by the instance name. 
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of instance IDs.
* `names` - A list of instance names.
* `instances` - A list of instances. Each element contains the following attributes:
  * `id` - ID of the instance.
  * `instance_id` - ID of the instance.
  * `instance_name` - Name of the instance.
  * `instance_type` - The type of the instance. 
  * `instance_status` - The status of the instance.
  * `independent_naming` - Indicates whether namespaces are available.
  * `tps_receive_max` - This attribute is used to set the message receiving transactions per second (TPS) of the topic during a certain period of time.
  * `tps_send_max` - This attribute is used to set the message sending transactions per second (TPS) of the topic during a certain period of time.
  * `topic_capacity` - This attribute is used to set the topic capacity.
  * `cluster` - This attribute is a used to add cluster name.
  * `create_time` - Create time of the instance.