---
subcategory: "GPDB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_gpdb_dbinstance"
sidebar_current: "docs-Alibabacloudstack-gpdb-dbinstance"
description: |- 
  Provides a gpdb Dbinstance resource.
---

# alibabacloudstack_gpdb_dbinstance
-> **NOTE:** Alias name has: `alibabacloudstack_gpdb_instance`

Provides a gpdb Dbinstance resource.

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
  engine                  = "gpdb"
  engine_version          = "4.3"
  instance_class          = "gpdb.group.segsdx2"
  instance_group_count    = "2"
  description            = "Terraform Test GPDB Instance"
  db_instance_description = "Terraform Test GPDB Instance Description"
  availability_zone      = data.alibabacloudstack_zones.default.zones[0].id
  vswitch_id             = alibabacloudstack_vswitch.default.id
  security_ip_list       = ["10.168.1.12", "100.69.7.112"]
}
```

## Argument Reference

The following arguments are supported:

* `availability_zone` - (Optional, ForceNew) The availability zone where the DB instance will be created. If not specified, Terraform will use the first available zone from the `alibabacloudstack_zones` data source.
* `instance_class` - (Required, ForceNew) The specification of the GPDB instance. For example, `gpdb.group.segsdx2`. Refer to [Instance Specifications](https://www.alibabacloud.com/help/doc-detail/86942.htm) for more details.
* `db_instance_class` - (Optional, ForceNew) Deprecated. Use `instance_class` instead.
* `instance_group_count` - (Required) The number of groups in the GPDB instance. Valid values are `[2, 4, 8, 16, 32]`.
* `instance_charge_type` - (Optional, ForceNew) The billing method of the instance. Valid values are `PrePaid` (Subscription) and `PostPaid` (Pay-As-You-Go). Default value is `PostPaid`.
* `payment_type` - (Optional, ForceNew) Alias for `instance_charge_type`. Specifies the payment type of the instance.
* `description` - (Optional) A brief description of the GPDB instance. It can be up to 256 characters long.
* `db_instance_description` - (Optional) Alias for `description`. Provides a detailed description of the GPDB instance.
* `vswitch_id` - (Optional, ForceNew) The ID of the VSwitch in which the GPDB instance will be launched. This parameter is required if you want to launch the instance in a VPC.
* `instance_inner_connection` - (Optional, ForceNew) The internal connection endpoint of the GPDB instance.
* `instance_inner_port` - (Optional, ForceNew) The internal port of the GPDB instance.
* `port` - (Optional, ForceNew) The public port number of the GPDB instance. Default value is `5432`.
* `engine` - (Required, ForceNew) Database engine. Currently, only `gpdb` is supported.
* `engine_version` - (Required, ForceNew) The version of the GPDB engine. Valid values include `4.3`, `6.0`, and `7.0`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `availability_zone` - The availability zone where the GPDB instance resides.
* `instance_class` - The specification of the GPDB instance.
* `db_instance_class` - Alias for `instance_class`.
* `instance_id` - The unique identifier of the GPDB instance.
* `db_instance_id` - Alias for `instance_id`.
* `region_id` - The region where the GPDB instance is located.
* `status` - The current status of the GPDB instance.
* `instance_network_type` - The network type of the GPDB instance. Possible values include `VPC`.
* `network_type` - Alias for `instance_network_type`.
* `instance_charge_type` - The billing method of the GPDB instance.
* `payment_type` - Alias for `instance_charge_type`.
* `description` - The description of the GPDB instance.
* `db_instance_description` - Alias for `description`.
* `vswitch_id` - The ID of the VSwitch associated with the GPDB instance.
* `instance_inner_connection` - The internal connection endpoint of the GPDB instance.
* `instance_inner_port` - The internal port of the GPDB instance.
* `port` - The public port number of the GPDB instance.
* `instance_vpc_id` - The ID of the VPC associated with the GPDB instance.
* `vpc_id` - Alias for `instance_vpc_id`.
* `engine` - The database engine used by the GPDB instance.
* `engine_version` - The version of the database engine used by the GPDB instance.