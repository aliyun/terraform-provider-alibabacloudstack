---
subcategory: "AnalyticDB for PostgreSQL"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_gpdb_dbinstance"
sidebar_current: "docs-Alibabacloudstack-gpdb-dbinstance"
description: |- 
  Provides a GPDB DB Instance resource.
---

# alibabacloudstack_gpdb_dbinstance
-> **NOTE:** Alias name has: `alibabacloudstack_gpdb_instance`

Provides a GPDB DB Instance resource.

## Example Usage

### Create a GPDB Instance with VPC and Security IP List

```hcl
data "alibabacloudstack_zones" "default" {
  available_resource_creation = "Gpdb"
}

resource "alibabacloudstack_vpc" "default" {
  name       = "vpc-123456"
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vswitch" "default" {
  availability_zone = data.alibabacloudstack_zones.default.zones[0].id
  vpc_id           = alibabacloudstack_vpc.default.id
  cidr_block       = "172.16.0.0/24"
  name             = "vswitch-123456"
}

resource "alibabacloudstack_gpdb_dbinstance" "example" {
  db_instance_mode          = "Classic"
  db_instance_class         = "gpdb.group.segsdx2"
  instance_group_count      = "2"
  db_instance_description   = "Terraform Test GPDB Instance"
  availability_zone         = data.alibabacloudstack_zones.default.zones[0].id
  vswitch_id                = alibabacloudstack_vswitch.default.id
  security_ip_list          = ["10.168.1.12", "100.69.7.112"]
  engine                    = "gpdb"
  engine_version            = "4.3"
}
```

## Argument Reference

The following arguments are supported:

* `db_instance_mode` - (Required, ForceNew) The deployment mode of the GPDB instance. Valid values are `Classic` and `StorageReserver`. When set to `StorageReserver`, `db_instance_storage_type` and `seg_node_num` become required.
* `db_instance_class` - (Optional, ForceNew) The specification of the GPDB instance. For example, `gpdb.group.segsdx2`. Either `db_instance_class` or the deprecated `instance_class` must be specified.
* `instance_class` - (Optional, ForceNew, Deprecated) Deprecated. Use `db_instance_class` instead.
* `instance_group_count` - (Optional) The number of node groups in the GPDB instance. Valid values are `2`, `4`, `8`, `16`, `32`.
* `instance_charge_type` - (Optional, ForceNew, Deprecated) Deprecated. Use `payment_type` instead. The billing method of the instance. Valid value: `PostPaid` (Pay-As-You-Go).
* `payment_type` - (Optional, ForceNew) Alias for `instance_charge_type`. The billing method of the instance. Valid value: `PostPaid` (Pay-As-You-Go).
* `db_instance_description` - (Optional) A description of the GPDB instance. It can be 2 to 256 characters in length.
* `description` - (Optional, Deprecated) Deprecated. Use `db_instance_description` instead.
* `vswitch_id` - (Optional, ForceNew) The ID of the VSwitch. If specified, the instance will be created in the VPC associated with the VSwitch.
* `instance_inner_connection` - (Optional, ForceNew) The internal connection prefix of the GPDB instance.
* `engine` - (Optional, ForceNew) The database engine. Valid value: `gpdb`.
* `engine_version` - (Optional, ForceNew) The engine version of the GPDB instance. Valid values include `4.3`, `6.0`, `7.0`.
* `availability_zone` - (Optional, ForceNew) The zone ID where the instance will be created. If not specified, the system will automatically allocate one.
* `cpu_type` - (Optional, ForceNew) The CPU architecture type of the instance.
* `security_ip_list` - (Optional) A list of IP addresses allowed to access the instance.
* `db_instance_storage_type` - (Optional, ForceNew) The storage type of the instance. Required when `db_instance_mode` is set to `StorageReserver`.
* `seg_node_num` - (Optional) The number of compute nodes. Required when `db_instance_mode` is set to `StorageReserver`.
* `instance_pay_type` - (Optional) The payment type of the instance.
* `network_type` - (Optional, Deprecated) The network type of the instance. This field will be deprecated in version 3.21.0. It will be automatically determined based on whether `vswitch_id` is configured.
* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in format of `instance_id`.
* `db_instance_id` - The unique identifier of the GPDB instance.
* `instance_id` - (Deprecated) Alias for `db_instance_id`.
* `region_id` - The region ID of the instance.
* `status` - The status of the instance.
* `instance_network_type` - The network type of the instance.
* `network_type` - (Deprecated) Alias for `instance_network_type`.
* `instance_charge_type` - (Deprecated) The billing method.
* `payment_type` - The billing method of the instance.
* `description` - (Deprecated) The description of the instance.
* `db_instance_description` - The description of the instance.
* `vswitch_id` - The ID of the VSwitch.
* `instance_inner_connection` - The internal connection endpoint.
* `instance_inner_port` - (Deprecated) The internal port. Use `port` instead.
* `port` - The connection port of the instance.
* `vpc_id` - The VPC ID of the instance.
* `instance_vpc_id` - (Deprecated) Alias for `vpc_id`.
* `engine` - The database engine.
* `engine_version` - The engine version.
* `instance_class` - The instance specification.
* `db_instance_class` - The instance specification.
* `availability_zone` - The availability zone.
* `instance_group_count` - The number of node groups.
* `segment_instance_amount` - The number of segment instances.
* `port` - The connection port.

## Import

GPDB DB Instance can be imported using the instance ID, e.g.

```
$ terraform import alibabacloudstack_gpdb_dbinstance.example gp-gs5xxxxxxxxxx
```