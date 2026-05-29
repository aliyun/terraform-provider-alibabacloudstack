---
subcategory: "Cloud-Native Distributed Database PolarDB-X 2.0"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_polardbx_readonly_instance"
sidebar_current: "docs-alibabacloudstack-polardbx-readonly-instance"
description: |-
  Provides a PolarDB-X read-only instance resource.
---

# alibabacloudstack\_polardbx\_readonly\_instance

Provides a PolarDB-X read-only instance resource.

-> **Note:** This resource can also be referred to by the following aliases:
- `apsarastack_polardbx_readonly_instance`

## Example Usage

```hcl
variable "name" {
  default = "tf_acc_drds_polardb_readonly"
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

# Create primary instance first
resource "alibabacloudstack_polardbx_instance" "primary" {
  vswitch_id     = alibabacloudstack_vpc_vswitch.default.id
  cn_node_count  = 2
  dn_node_class  = "mysql.n4.medium.25"
  dn_node_count  = 2
  description    = "primary instance"
  storage        = 50
  cn_node_class  = "polarx.x4.medium.2e"
  engine_version = "5.7"
  zone_id        = data.alibabacloudstack_zones.default.zones.0.id
}

# Create read-only instance
resource "alibabacloudstack_polardbx_readonly_instance" "readonly" {
  primary_db_instance_id = alibabacloudstack_polardbx_instance.primary.id
  vswitch_id             = alibabacloudstack_vpc_vswitch.default.id
  cn_node_count          = 2
  dn_node_class          = "mysql.n4.medium.25"
  dn_node_count          = 2
  description            = "readonly instance"
  storage                = 50
  cn_node_class          = "polarx.x4.medium.2e"
  engine_version         = "5.7"
  zone_id                = data.alibabacloudstack_zones.default.zones.0.id
}
```

## Argument Reference

The following arguments are supported:

### Required

* `primary_db_instance_id` - (Required, ForceNew) The ID of the primary instance. Changing this parameter will force the creation of a new read-only instance.
* `storage` - (Required, ForceNew) The storage size of the PolarDB-X read-only instance, in GB.
* `cn_node_class` - (Required) The specification of the compute nodes (CN). Example values: `polarx.x4.medium.2e`, `polarx.x4.large.2e`, `polarx.x8.large.2e`, `polarx.x4.xlarge.2e`, `polarx.x8.xlarge.2e`, `polarx.x4.2xlarge.2e`, `polarx.x8.2xlarge.2e`, `polarx.x4.4xlarge.2e`, `polarx.x8.4xlarge.2e`.
* `cn_node_count` - (Required) The number of compute nodes (CN).
* `dn_node_class` - (Required) The specification of the data nodes (DN). Example values: `mysql.n4.medium.25`, `mysql.n4.large.25`, `mysql.n4.xlarge.25`, `mysql.n4.2xlarge.25`, `mysql.x4.medium.25`.
* `dn_node_count` - (Required) The number of data nodes (DN).
* `vswitch_id` - (Required, ForceNew) The ID of the VSwitch.

### Optional

* `cpu_type` - (Optional) The CPU architecture type of the PolarDB-X read-only instance.
* `description` - (Optional) The description of the PolarDB-X read-only instance.
* `engine_version` - (Optional, ForceNew) The engine version of the PolarDB-X read-only instance. Valid values: `5.7`, `8.0`.
* `gms_node_class` - (Optional) The specification of the GMS (Global Meta Service) node.
* `polardbx_instance_id` - (Optional) The ID of the PolarDB-X read-only instance.
* `primary_zone` - (Optional, ForceNew) The primary availability zone.
* `resource_type` - (Optional, ForceNew) The resource type. Currently only PolarDB-X 2.0 instance type is supported.
* `secondary_zone` - (Optional, ForceNew) The secondary availability zone.
* `tertiary_zone` - (Optional, ForceNew) The tertiary availability zone.
* `topology_type` - (Optional, ForceNew) The topology type of the instance. Valid values: `1azone` (single zone), `3azones` (three zones). Default: `1azone`.
* `enable_tde` - (Optional) Whether to enable TDE (Transparent Data Encryption). Valid values: `true`, `false`. Default: `false`. Once enabled, TDE cannot be disabled.
* `enable_ssl` - (Optional) Whether to enable SSL encryption. Valid values: `true`, `false`. Default: `false`.
* `zone_id` - (Optional, ForceNew) The ID of the zone to which the instance belongs.
* `compute_parameters` - (Optional) The compute resource configuration of the instance. This is a set of key-value pairs.
    * `name` - (Required) The name of the compute parameter.
    * `value` - (Required) The value of the compute parameter.
* `storage_parameters` - (Optional) The storage configuration of the instance. This is a set of key-value pairs.
    * `name` - (Required) The name of the storage parameter.
    * `value` - (Required) The value of the storage parameter.
* `security_groups` - (Optional) The security IP groups of the instance.
    * `group_name` - (Required) The name of the security IP group.
    * `ips` - (Required) The security IP addresses. Multiple IPs should be separated by ",", e.g., `192.168.0.1,192.168.0.0/24`.
* `private_connection_string_prefix` - (Optional) The prefix of the private connection string.
* `private_connection_port` - (Optional) The port of the private connection. Valid values: 3000-6000. Default: `3306`.
* `enable_public_connection` - (Optional) Whether to enable public connection. Valid values: `true`, `false`. Default: `false`.
* `public_connection_string_prefix` - (Optional) The prefix of the public connection string. Required when `enable_public_connection` is `true`.
* `public_connection_port` - (Optional) The port of the public connection. Valid values: 3000-6000. Default: `3306`.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The ID of the PolarDB-X read-only instance.
* `create_time` - The creation time of the instance.
* `network_type` - The network type of the instance.
* `status` - The status of the instance.
* `connection_string` - The private connection string of the instance.
* `private_connection_string` - The private connection string.
* `private_connection_port` - The private connection port.
* `public_connection_string` - The public connection string (if enabled).
* `public_connection_port` - The public connection port (if enabled).

## Import

PolarDB-X read-only instance can be imported using the instance ID, e.g.

```
$ terraform import alibabacloudstack_polardbx_readonly_instance.example px-12345678
```
