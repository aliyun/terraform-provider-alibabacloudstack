---
subcategory: "Table Store (OTS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ots_instance_attachment"
sidebar_current: "docs-Alibabacloudstack-ots-instanceattachment"
description: |- 
  Provides an OTS (Open Table Service) resource to attach VPC to instance.
---

# alibabacloudstack_ots_instance_attachment
-> **NOTE:** Alias name has: `alibabacloudstack_ots_instanceattachment`

This resource will help you to bind a VPC to an OTS instance.

## Example Usage

```hcl
variable "name" {
  default = "tf-testAcc84399"
}

resource "alibabacloudstack_ots_instance" "default" {
  name        = "${var.name}"
  description = "${var.name}"
  accessed_by = "Vpc"
  instance_type = "Capacity"
}

data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alibabacloudstack_vpc" "default" {
  cidr_block = "172.16.0.0/16"
  name       = "${var.name}"
}

resource "alibabacloudstack_vswitch" "default" {
  vpc_id            = alibabacloudstack_vpc.default.id
  vswitch_name      = "${var.name}"
  cidr_block        = "172.16.1.0/24"
  zone_id           = data.alibabacloudstack_zones.default.zones[0].id
}

resource "alibabacloudstack_ots_instance_attachment" "default" {
  instance_name = alibabacloudstack_ots_instance.default.name
  vpc_name      = "test-attachment"
  vswitch_id    = alibabacloudstack_vswitch.default.id
}
```

## Argument Reference

The following arguments are supported:

* `instance_name` - (Required, ForceNew) The name of the OTS instance. This must match the name of an existing OTS instance.
* `vpc_name` - (Required, ForceNew) The name of the VPC being attached to the OTS instance. This is used for identification purposes.
* `vswitch_id` - (Required, ForceNew) The ID of the VSwitch being attached to the OTS instance. This must be within the same VPC as specified by `vpc_name`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID. The value is the same as `instance_name`.
* `instance_name` - The name of the OTS instance.
* `vpc_name` - The name of the VPC attached to the OTS instance.
* `vswitch_id` - The ID of the VSwitch attached to the OTS instance.
* `vpc_id` - The ID of the VPC attached to the OTS instance. This attribute is automatically derived from the `vswitch_id`.