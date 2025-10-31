---
subcategory: "Lindorm"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_lindorm_instances"
description: |-
  Provides a list of Lindorm Instances to the user.
---

# alibabacloudstack\_lindorm\_instances

This data source provides Lindorm instances available to the user.


## Example Usage

### Basic Usage

```terraform
data "alibabacloudstack_lindorm_instances" "example" {
  ids = ["i-12345678"]
}

output "first_instance_id" {
  value = data.alibabacloudstack_lindorm_instances.example.instances.0.id
}
```

### Filter by Description Regex

```terraform
data "alibabacloudstack_lindorm_instances" "example" {
  description_regex = "my-instance"
}

output "instance_ids" {
  value = data.alibabacloudstack_lindorm_instances.example.ids
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of instance IDs.
* `name_regex` - (Optional, Deprecated) A regex string to filter results by instance description. Field 'name_regex' is deprecated and will be removed in a future release. Please use new field 'description_regex' instead.
* `description_regex` - (Optional) A regex string to filter results by instance description.

## Attributes Reference

The following attributes are exported:

* `ids` - A list of instance IDs.
* `instances` - A list of Lindorm instances. Each element contains the following attributes:
  * `id` - The ID of the instance.
  * `instance_id` - The ID of the instance.
  * `cpu_brand` - The brand of CPU used by the instance.
  * `instance_storage` - The storage capacity of the instance.
  * `zone_id` - The ID of the zone where the instance is deployed.
  * `create_time` - The creation time of the instance.
  * `ascm_create_user` - The user who created the instance.
  * `instance_alias` - The alias of the instance.
  * `network_type` - The network type of the instance.
  * `service_type` - The service type of the instance.
  * `engine_type` - The engine type of the instance.
  * `ali_uid` - The UID of the Alibaba Cloud account.