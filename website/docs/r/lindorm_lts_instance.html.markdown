---
subcategory: "Lindorm"
layout: "alibabacloudstack"
page_title: "AlibabacloudStack: alibabacloudstack_lindorm_lts_instance"
sidebar_current: "docs-alibabacloudstack-resource-lindorm-lts-instance"
description: |-
  Provides a AlibabacloudStack Lindorm LTS Instance resource.
---

# alibabacloudstack_lindorm_lts_instance

Provides a Lindorm LTS Instance resource.


## Example Usage

Basic Usage

```terraform

variable "name" {
    default = "tf-testacc97984"
}

data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
  enable_details = true
}

data "alibabacloudstack_lindorm_instance_types" "sortbycpu" {
	sorted_by = "CPU"
	engine_type = "lts"
}

resource "alibabacloudstack_lindorm_lts_instance" "example" {
  zone_id        = "${data.alibabacloudstack_zones.default.zones.0.id}"
  instance_alias = "${var.name}"
  cpu_brand      = "Intel"
  instance_type  = "${data.alibabacloudstack_lindorm_instance_types.sortbycpu.instance_types[0].name}"
  lts_num        = 2
}

```

## Argument Reference

The following arguments are supported:

* `zone_id` - (Required, ForceNew) The zone ID of the instance.
* `instance_alias` - (Required) The alias of the instance.
* `cpu_brand` - (Required, ForceNew) The CPU brand. 
* `instance_type` - (Required) The specification of the instance.
* `lts_num` - (Required) The number of LTS nodes.
* `deletion_protection` - (Optional, Computed) Specifies whether to enable deletion protection for the instance. Default value: `false`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the instance.
* `instance_id` - The ID of the instance.
* `instance_status` - The status of the instance.
* `create_time` - The creation time of the instance.
* `service_type` - The service type of the instance.
* `network_type` - The network type of the instance.

## Import

Lindorm LTS Instance can be imported using the id, e.g.

```bash
$ terraform import alibabacloudstack_lindorm_lts_instance.example <id>
```
