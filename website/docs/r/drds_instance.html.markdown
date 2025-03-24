---
subcategory: "DRDS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_drds_instance"
sidebar_current: "docs-Alibabacloudstack-drds-instance"
description: |- 
  Provides a drds Instance resource.
---

# alibabacloudstack_drds_instance

Provides a drds Instance resource.

For information about DRDS and how to use it, see [What is DRDS](https://www.alibabacloud.com/help/doc-detail/29659.htm).

-> **NOTE:** At present, DRDS instance only can be supported in the regions: cn-shenzhen, cn-beijing, cn-hangzhou, cn-hongkong, cn-qingdao, ap-southeast-1.

-> **NOTE:** Currently, this resource only support `Domestic Site Account`.

## Example Usage

```hcl

variable "name" {
	default = "tf-testaccDrdsdatabase-14880"
}

data "alibabacloudstack_zones" "default" {
	available_resource_creation = "VSwitch"
}

variable "instance_series" {
	default = "drds.sn2.4c16g"
}

resource "alibabacloudstack_vpc" "default" {
	name       = var.name
	cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vswitch" "default" {
	vpc_id            = alibabacloudstack_vpc.default.id
	cidr_block        = "172.16.0.0/24"
	availability_zone = data.alibabacloudstack_zones.default.zones.0.id
	name              = var.name
}

resource "alibabacloudstack_drds_instance" "default" {
	description          = var.name
	instance_charge_type = "PostPaid"
	zone_id              = alibabacloudstack_vswitch.default.availability_zone
	vswitch_id           = alibabacloudstack_vswitch.default.id
	instance_series      = var.instance_series
	specification        = "drds.sn2.4c16g.8C32G"
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Required) Description of the DRDS instance. This description can have a string of 2 to 256 characters.
* `zone_id` - (Required, ForceNew) The Zone to launch the DRDS instance.
* `specification` - (Required, ForceNew) User-defined DRDS instance specification. Value range:
    - For `drds.sn1.4c8g` (Starter version):
        - `drds.sn1.4c8g.8c16g`, `drds.sn1.4c8g.16c32g`, `drds.sn1.4c8g.32c64g`, `drds.sn1.4c8g.64c128g`
    - For `drds.sn1.8c16g` (Standard edition):
        - `drds.sn1.8c16g.16c32g`, `drds.sn1.8c16g.32c64g`, `drds.sn1.8c16g.64c128g`
    - For `drds.sn1.16c32g` (Enterprise Edition):
        - `drds.sn1.16c32g.32c64g`, `drds.sn1.16c32g.64c128g`
    - For `drds.sn1.32c64g` (Extreme Edition):
        - `drds.sn1.32c64g.128c256g`
* `instance_charge_type` - (Optional, ForceNew) Valid values are `PrePaid`, `PostPaid`. Default to `PostPaid`.
* `vswitch_id` - (Required, ForceNew) The VSwitch ID to launch in.
* `instance_series` - (Required, ForceNew) User-defined DRDS instance node spec. Value range:
    - `drds.sn1.4c8g` for DRDS instance Starter version;
    - `drds.sn1.8c16g` for DRDS instance Standard edition;
    - `drds.sn1.16c32g` for DRDS instance Enterprise Edition;
    - `drds.sn1.32c64g` for DRDS instance Extreme Edition;

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 10 mins) Used when creating the drds instance (until it reaches running status).
* `delete` - (Defaults to 10 mins) Used when terminating the drds instance.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The DRDS instance ID.

## Import

Distributed Relational Database Service (DRDS) can be imported using the id, e.g.

```bash
$ terraform import alibabacloudstack_drds_instance.example drds-abc123456
```