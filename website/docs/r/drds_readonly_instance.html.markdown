---
subcategory: "Distributed Relational Database Service(DRDS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_drds_readnoly_instance"
sidebar_current: "docs-Alibabacloudstack-drds-instance"
description: |- 
  Provides a drds Read Only Instance resource.
---

# alibabacloudstack_drds_readonly_instance

Provides a drds Read Only Instance resource.

For information about DRDS and how to use it, see [What is DRDS](https://www.alibabacloud.com/help/doc-detail/29659.htm).

-> **NOTE:** At present, DRDS instance only can be supported in the regions: cn-shenzhen, cn-beijing, cn-hangzhou, cn-hongkong, cn-qingdao, ap-southeast-1.

-> **NOTE:** Currently, this resource only support `Domestic Site Account`.

## Example Usage

```hcl
variable "name" {
  default = "tf-testacc-readonlyinstance-19500"
}

data "alibabacloudstack_drds_instance_specifications" "default" {
  sorted_by = "CPU"
}

resource "alibabacloudstack_drds_instance" "default" {
  description   = var.name
  zone_id       = alibabacloudstack_vpc_vswitch.default.availability_zone
  vswitch_id    = alibabacloudstack_vpc_vswitch.default.id
  specification = data.alibabacloudstack_drds_instance_specifications.default.specifications.0.id
}


data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
  enable_details              = true
}


resource "alibabacloudstack_vpc_vpc" "default" {
  vpc_name   = "${var.name}_vpc"
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vpc_vswitch" "default" {
  name       = "${var.name}_vsw"
  vpc_id     = alibabacloudstack_vpc_vpc.default.id
  cidr_block = "172.16.1.0/24"
  zone_id    = data.alibabacloudstack_zones.default.zones.0.id
}




resource "alibabacloudstack_drds_readonly_instance" "default" {
  master_instance_id   = alibabacloudstack_drds_instance.default.id
  zone_id              = alibabacloudstack_vpc_vswitch.default.availability_zone
  instance_charge_type = "PostPaid"
  vswitch_id           = alibabacloudstack_vpc_vswitch.default.id
  specification        = data.alibabacloudstack_drds_instance_specifications.default.specifications.0.id
  description          = var.name
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Optional) Description of the DRDS instance. This description can have a string of 2 to 256 characters.
* `zone_id` - (Required, ForceNew) The Zone to launch the DRDS instance.
* `specification` - (Required, ForceNew) User-defined DRDS instance specification. Value range:
* `instance_charge_type` - (Optional, ForceNew) Valid values are `PrePaid`, `PostPaid`. Default to `PostPaid`.
* `vswitch_id` - (Optional, ForceNew) The VSwitch ID to launch in. If not set while use classic network.
* `master_instance_id` - (Required, ForceNew) The maseter instance id.

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