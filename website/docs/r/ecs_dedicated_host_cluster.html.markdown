---
subcategory: "Elastic Compute Service"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ecs_dedicated_host_cluster"
description: |-
  Provides a ecs Dedicatedhostcluster resource.
---

# alibabacloudstack\_ecs\_dedicatedhostcluster

Provides a ecs Dedicatedhostcluster resource.

## Example Usage
```
variable "name" {
    default = "tf-testaccecsdedicated_hostcluster20348"
}



data "alibabacloudstack_zones" default {
  available_resource_creation = "VSwitch"
  enable_details = true
}





resource "alibabacloudstack_ecs_dedicated_host_cluster" "default" {
  zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
  dedicated_host_cluster_name = "${var.name}"
}
```

## Argument Reference

The following arguments are supported:
  * `dedicated_host_cluster_name` - (Optional) - dedicated host cluster name
  * `description` - (Optional) - description
  * `zone_id` - (Optional) - zone id

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
  * `dedicated_host_cluster_id` - dedicated host cluster id
