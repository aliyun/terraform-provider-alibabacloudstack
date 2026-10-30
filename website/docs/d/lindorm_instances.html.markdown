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
  * `service_type` - The raw service type reported by `GetLindormInstanceList`, e.g. `"lindorm"` or `"lindorm_dual_zone"`.
  * `engine_type` - The normalized engine name, consistent with the `engine_type` argument of `alibabacloudstack_lindorm_instance`. For multi-zone instances (`lindorm_dual_zone`) this is reported as `"lindorm"`.
  * `ali_uid` - The UID of the Alibaba Cloud account.
  * `topology_type` - The topology type derived from `service_type`. Valid values: `"1azone"` (single zone) or `"multizone"` (dual zone).
  * `arch_version` - The architecture version. Note: this field is not present in the `GetLindormInstanceList` response and will be empty.
  * `standby_zone_id` - The secondary zone ID for multi-zone instances (sourced from `Idc2ZoneId`). Empty for single-zone instances.
  * `arbiter_zone_id` - The arbiter zone ID. Not returned by the list API; reserved for future use.
  * `standby_vswitch_id` - The vSwitch ID in the standby zone (sourced from `Idc2VSwitchId`). Empty for single-zone instances or when the list API does not return it.
  * `arbiter_vswitch_id` - The vSwitch ID in the arbiter zone. Not returned by the list API; reserved for future use.