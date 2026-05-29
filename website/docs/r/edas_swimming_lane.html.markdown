---
subcategory: "Enterprise Distributed Application Service"
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
* `apps` - (Required) List of application IDs associated with the lane. At least one application ID must be specified.
* `priority` - (Required) The priority of the swimming lane. Valid values: 1 to 100.
* `path` - (Required) The request path matched by the swimming lane.
* `condition` - (Required) The logical relationship between matching conditions. Valid values: `OR` (any condition satisfied), `ADD` (all conditions satisfied).
* `rest_items` - (Required) List of traffic rule items. Each item supports the following:
  * `type` - (Required) The type of the rule. Valid values: `cookie`, `header`, `param`.
  * `name` - (Required) The key name of the rule.
  * `value` - (Required) The value of the rule.
  * `cond` - (Required) The matching condition. Valid values: `==`, `!=`, `>`, `<`, `>=`, `<=`, `list`.
  * `operator` - (Optional) The type of the value. Valid values: `rawvalue` (original value), `mod` (modulo), `list` (list value). Default value: `rawvalue`.
* `logical_region_id` - (Optional) The logical region ID of the EDAS instance. Format: `physical-region-id:custom-namespace-identifier`, e.g., `cn-hangzhou:test`. If not specified, the provider's default region ID will be used.
* `enabled` - (Optional) Whether the swimming lane is enabled. Default value: `false`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the swimming lane. Format: `logical_region_id:group_id:lane_id`.
* `lane_id` - The unique identifier of the swimming lane in EDAS.

## Import

EDAS swimming lane can be imported using the `logical_region_id:group_id:lane_id` format ID, e.g.

```bash
$ terraform import alibabacloudstack_edas_swimming_lane.example cn-hangzhou:test:12345:88
```