---
subcategory: "Enterprise Distributed Application Service (EDAS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_edas_swimming_lane"
sidebar_current: "docs-alibabacloudstack-resource-edas-swimming-lane"
description: |-
  Provides a AlibabacloudStack EDAS Swimming Lane resource.
---

# alibabacloudstack\_edas\_swimming\_lane

Provides an EDAS swimming lane resource.


## Example Usage

```terraform
resource "alibabacloudstack_edas_swimming_lane" "example" {
  name             = "example_value"
  group_id         = "12345"
  apps             = ["example_app_id"]
  priority         = 1
  path             = "/example"
  condition        = "OR"
  enabled          = true

  rest_items {
    type     = "header"
    name     = "example_header"
    value    = "example_value"
    cond     = "=="
    operator = "rawvalue"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the swimming lane.
* `group_id` - (Required) The ID of the swimming lane group.
* `apps` - (Required) The list of application IDs. At least one application ID must be specified.
* `priority` - (Required) The priority of the swimming lane.
* `path` - (Required) The path of the swimming lane.
* `condition` - (Required) The condition of the swimming lane. Valid values: `OR`, `ADD`.
* `rest_items` - (Required) The REST items of the swimming lane. Each item supports the following:
  * `type` - (Required) The type of the REST item. Valid values: `cookie`, `header`, `param`.
  * `name` - (Required) The name of the REST item.
  * `value` - (Required) The value of the REST item.
  * `cond` - (Required) The condition of the REST item. Valid values: `==`, `!=`, `>`, `<`, `>=`, `<=`, `list`.
  * `operator` - (Optional) The operator of the REST item. Valid values: `rawvalue`, `mod`, `list`. Default value: `rawvalue`.
* `logical_region_id` - (Optional) The logical region ID of the EDAS instance. If not specified, the default region ID of the provider will be used.
* `enabled` - (Optional) Whether the swimming lane is enabled. Default value: `false`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the swimming lane.
* `lane_id` - The lane ID of the swimming lane.

## Import

EDAS swimming lane can be imported using the id, e.g.

```bash
$ terraform import alibabacloudstack_edas_swimming_lane.example logical_region_id:group_id:lane_id
```