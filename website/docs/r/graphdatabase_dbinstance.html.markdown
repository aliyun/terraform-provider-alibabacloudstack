---
subcategory: "GraphDatabase(GPDB)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_graphdatabase_dbinstance"
sidebar_current: "docs-Alibabacloudstack-graphdatabase-dbinstance"
description: |- 
  Provides a graphdatabase Dbinstance resource.
---

# alibabacloudstack_graphdatabase_dbinstance
-> **NOTE:** Alias name has: `alibabacloudstack_graph_database_db_instance`

Provides a graphdatabase Dbinstance resource.

## Example Usage

```hcl
variable "name" {
    default = "tf-testaccgraph_databasedb_instance98846"
}

resource "alibabacloudstack_graphdatabase_dbinstance" "default" {
  db_instance_network_type = "VPC"
  payment_type             = "PayAsYouGo"
  vswitch_id              = "vsw-bp152wgftimgq80eiii6k"
  zone_id                 = "cn-hangzhou-h"
  db_instance_storage_type = "cloud_ssd"
  db_instance_description  = "ssd测试"
  db_node_class           = "gdb.r.2xlarge"
  db_instance_category    = "ha"
  vpc_id                  = "vpc-bp1bvsykm9f9hkfeikfi5"
  db_version              = "1.0"
  region_id               = "cn-hangzhou"
  db_node_storage         = 100
  db_instance_ip_array = [
    {
      db_instance_ip_array_attribute = "hidden"
      db_instance_ip_array_name      = "default_whitelist"
      security_ips                   = ["0.0.0.0/0", "192.168.1.1"]
    }
  ]
}
```

## Argument Reference

The following arguments are supported:

* `db_instance_ip_array` - (Optional) The IP whitelist configuration for the instance group list. Each element in the list is an object with the following structure:
  * `db_instance_ip_array_attribute` - (Optional) The attribute of the IP whitelist group. The default value is `hidden`.
  * `db_instance_ip_array_name` - (Optional) The name of the IP whitelist group.
  * `security_ips` - (Required) A list of IP addresses or CIDR blocks in the IP whitelist, separated by commas. The maximum number of entries is 1000. For example: `0.0.0.0/0`, `192.168.1.1`, or `192.168.1.0/24`.

* `db_instance_category` - (Required, ForceNew) The category of the DB instance. Valid values: `ha` (High Availability).
* `db_instance_description` - (Optional) The description of the DB instance according to practical examples or notes.
* `db_instance_network_type` - (Required, ForceNew) The network type of the DB instance. Valid values: `VPC`.
* `db_instance_storage_type` - (Required, ForceNew) The storage type of the DB instance. Valid values: `cloud_ssd`, `cloud_essd`.
* `db_node_class` - (Required) The class of the DB node. Valid values include: `gdb.r.xlarge`, `gdb.r.2xlarge`, `gdb.r.4xlarge`, `gdb.r.8xlarge`, `gdb.r.16xlarge`.
* `db_node_storage` - (Required) The storage space of the DB instance, measured in GB.
* `db_version` - (Required, ForceNew) The kernel version of the DB instance. Valid values: `1.0` (Gremlin), `1.0-OpenCypher` (OpenCypher).
* `payment_type` - (Optional, ForceNew) The payment type of the resource. Valid values: `PayAsYouGo`.
* `vswitch_id` - (Optional, ForceNew) The ID of the VSwitch where the DB instance will be created.
* `vpc_id` - (Optional, ForceNew) The ID of the VPC where the DB instance will be created.
* `zone_id` - (Optional, ForceNew) The ID of the availability zone where the DB instance will be created.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the resource.
* `status` - The status of the DB instance. Possible values include: `Creating`, `Running`, `Deleting`, `Rebooting`, `DBInstanceClassChanging`, `NetAddressCreating`, and `NetAddressDeleting`.
* `vswitch_id` - The ID of the VSwitch associated with the DB instance.
* `vpc_id` - The ID of the VPC associated with the DB instance.
* `zone_id` - The ID of the availability zone where the DB instance resides.