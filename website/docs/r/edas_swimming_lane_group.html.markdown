---
subcategory: "Enterprise Distributed Application Service (EDAS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_edas_swimming_lane_group"
description: |-
  Provides a AlibabacloudStack EDAS Swimming Lane Group resource.
---

# alibabacloudstack\_edas\_swimming\_lane\_group

Provides an EDAS swimming lane group resource.


## Example Usage

```terraform
resource "alibabacloudstack_edas_swimming_lane_group" "example" {
  name             = "example_value"
  entry_app_id     = "example_app_id"
  apps             = ["example_app_id"]
  logical_region_id = "cn-beijing:test"
  strategy_type    = "CONTENT"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the swimming lane group.
* `entry_app_id` - (Required) The ID of the entry application.
* `apps` - (Required) The list of application IDs.
* `logical_region_id` - (Optional) The logical region ID of the EDAS instance. If not specified, the default region ID of the provider will be used.
* `strategy_type` - (Optional) The strategy type of the swimming lane group. Valid values: `CONTENT`, `PERCENT`. Default value: `CONTENT`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the swimming lane group.
* `group_id` - The group ID of the swimming lane group.

## Import

EDAS swimming lane group can be imported using the id, e.g.

```bash
$ terraform import alibabacloudstack_edas_swimming_lane_group.example logical_region_id:group_id
```