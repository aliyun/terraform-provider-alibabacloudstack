---
subcategory: "Distributed Relational Database Service(DRDS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_drds_rds_instance"
description: |-
  Provides a DRDS RDS Instance resource.
---

# alibabacloudstack_drds_rds_instance

Provides a DRDS RDS Instance resource.

DRDS RDS Instance is a private RDS instance associated with a DRDS instance. It is used to store the data of DRDS databases.

## Example Usage

```hcl
variable "name" {
  default = "tf-testacc-drds-rds"
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

resource "alibabacloudstack_drds_instance" "default" {
  description          = "${var.name}"
  zone_id              = "${alibabacloudstack_vpc_vswitch.default.availability_zone}"
  instance_series      = "drds.sn2.4c16g"
  instance_charge_type = "PostPaid"
  vswitch_id           = "${alibabacloudstack_vpc_vswitch.default.id}"
  specification        = "drds.sn2.4c16g.8C32G"
}

resource "alibabacloudstack_drds_rds_instance" "default" {
  storage_type        = "local_ssd"
  category            = "HighAvailability"
  db_instance_class   = "rds.mysql.s1.small"
  drds_instance_id    = "${alibabacloudstack_drds_instance.default.id}"
  zone_id             = "${data.alibabacloudstack_zones.default.zones.0.id}"
  db_instance_storage = "20"
}
```

## Argument Reference

The following arguments are supported:

* `storage_type` - (Required, ForceNew) The storage type of the RDS instance. Valid values: `local_ssd`, `cloud_ssd`, `cloud_essd`. Changing this parameter will force a new resource to be created.
* `category` - (Required, ForceNew) The category of the RDS instance. Valid values: `HighAvailability`, `Finance`. Changing this parameter will force a new resource to be created.
* `drds_instance_id` - (Required, ForceNew) The ID of the DRDS instance to which the RDS instance belongs. Changing this parameter will force a new resource to be created.
* `zone_id` - (Required, ForceNew) The zone ID where the RDS instance will be created. Changing this parameter will force a new resource to be created.
* `db_instance_class` - (Required) The instance class (specification) of the RDS instance. Example: `rds.mysql.s1.small`.
* `db_instance_storage` - (Required) The storage capacity of the RDS instance, in GB.
* `force_remove` - (Optional) Whether to force remove the RDS instance when destroying. Valid values: `true`, `false`. Default to `false`.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 20 mins) Used when creating the DRDS RDS instance (until it reaches running status).
* `update` - (Defaults to 20 mins) Used when updating the DRDS RDS instance specification.
* `delete` - (Defaults to 20 mins) Used when deleting the DRDS RDS instance.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The ID of the resource. The format is `<drds_instance_id>:<rds_instance_id>`.
* `rds_instance_id` - The ID of the RDS instance.
* `create_time` - The creation time of the RDS instance.
* `status` - The status of the RDS instance. Valid values: `Creating`, `Running`, `Changing Specifications`, `Deleting`, `Stopping`, `Stopped`.

## Import

DRDS RDS Instance can be imported using the composite ID in the format `<drds_instance_id>:<rds_instance_id>`, e.g.

```bash
$ terraform import alibabacloudstack_drds_rds_instance.example drds-abc123456:rm-xyz789012
```
