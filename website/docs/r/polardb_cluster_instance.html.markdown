---
subcategory: "PolarDB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_polardb_cluster_instance"
sidebar_current: "docs-Alibabacloudstack-resource-polardb-cluster-instance"
description: |-
  Provides a PolarDB cluster instance resource.
---

# alibabacloudstack_polardb_cluster_instance

Provides a PolarDB cluster instance resource

## Example Usage

### Create a PolarDB MySQL cluster instance

```hcl
variable "name" {
  default = "tf-polar-test"
}

data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
  enable_details = true
}

resource "alibabacloudstack_vpc_vpc" "default" {
  vpc_name       = "${var.name}_vpc"
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vpc_vswitch" "default" {
  vpc_id            = "${alibabacloudstack_vpc_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  zone_id           = "${data.alibabacloudstack_zones.default.zones.0.id}"
  vswitch_name      = "${var.name}_vsw"
}

resource "alibabacloudstack_polardb_cluster_instance" "default" {
	db_cluster_description 	=  "${var.name}"
	zone_id 				= "${data.alibabacloudstack_zones.default.zones.0.id}"
	db_type 				= "${var.db_type}"
	db_version 				= "${var.db_version}"
	storage_space 			= "20"
	vpc_id 					= "${alibabacloudstack_vpc_vpc.default.id}"
	vswitch_id				= "${alibabacloudstack_vpc_vswitch.default.id}"
	db_node_class 			= "${data.alibabacloudstack_polardb_cluster_instance_types.default.instance_types.0.id}"
	sub_category 			= "${data.alibabacloudstack_polardb_cluster_instance_types.default.instance_types.0.sub_category}"
	storage_type 			= "ESSDPL1"
}
```

## Argument Reference

The following arguments are supported:

* `db_type` - (Required, ForceNew) The database engine type. Valid values: [MySQL](file://d:\terraform\terraform-provider-apsarastack\alibabacloudstack\extension_rds.go#L5-L5), `PolarDB_PG`, `PolarDB_PPAS`.
* `db_version` - (Required, ForceNew) The database engine version.
* `db_node_class` - (Required) The node class of the cluster instance.
* `zone_id` - (Required, ForceNew) The zone ID of the cluster instance.
* `vpc_id` - (Required, ForceNew) The ID of the VPC.
* `vswitch_id` - (Required, ForceNew) The ID of the vSwitch.
* `db_cluster_description` - (Optional) The description of the cluster instance.
* `storage_type` - (Optional, ForceNew) The storage type of the cluster instance.
* `storage_space` - (Optional, ForceNew) The storage space of the cluster instance.
* `readonly_node_num` - (Optional) The number of read-only nodes. Valid values: 0 to 15.
* `hot_standby_cluster` - (Optional) Specifies whether to enable the hot standby cluster. Valid values: `standby`, `off`.
* `cpu_type` - (Optional, ForceNew) The CPU type of the cluster instance. Valid values: `intel`, `hygon`.
* `sub_category` - (Optional) The sub category of the cluster instance. Valid values: `normal_general`, `normal_exclusive`.
* `deletion_lock` - (Optional) The deletion lock of the cluster instance. Valid values: 0, 1.
* `maintain_time` - (Optional) The maintenance time of the cluster instance.
* `security_ips_groups` - (Optional) The security IP groups of the cluster instance.
* `security_groups` - (Optional) The security groups of the cluster instance.
* `ssl_enabled` - (Optional) Whether to enable SSL encryption.
* `tde_enabled` - (Optional) Whether to enable Transparent Data Encryption (TDE).
* `encryption_key` - (Optional) The ID of the KMS key.
* `encrypt_algorithm` - (Optional) The encryption algorithm.
* `parameters` - (Optional) The parameters of the cluster instance.
* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the cluster instance.
* `db_cluster_id` - The ID of the cluster instance.
* `category` - The category of the cluster instance.
* `db_cluster_network_type` - The network type of the cluster instance.
* `storage_max` - The maximum storage space of the cluster instance.
* `engine` - The database engine.
* `zone_ids` - The zone IDs of the cluster instance.
* `db_nodes` - The database nodes information.
  * `db_node_status` - The status of the node.
  * `zone_id` - The zone ID of the node.
  * `max_connections` - The maximum connections of the node.
  * `added_cpu_cores` - The added CPU cores of the node.
  * `db_node_role` - The role of the node.
  * `imci_switch` - The IMCI switch of the node.
  * `db_node_id` - The ID of the node.
  * `max_iops` - The maximum IOPS of the node.
  * `db_node_class` - The class of the node.
  * `creation_time` - The creation time of the node.
  * `scc_mode` - The SCC mode of the node.
  * `failover_priority` - The failover priority of the node.
  * `hot_replica_mode` - The hot replica mode of the node.
  * `server_weight` - The server weight of the node.
* `proxy_status` - The proxy status of the cluster instance.
* `architecture` - The architecture of the cluster instance.
* `db_cluster_status` - The status of the cluster instance.
* `vip` - The VIP of the cluster instance.
* `lock_mode` - The lock mode of the cluster instance.
* `creation_time` - The creation time of the cluster instance.
* `sql_size` - The SQL size of the cluster instance.
* `proxy_cpu_cores` - The proxy CPU cores of the cluster instance.
* `pay_type` - The payment type of the cluster instance.

## Import

PolarDB cluster instance can be imported using the id, e.g.

```bash
$ terraform import alibabacloudstack_polardb_cluster_instance.example pc-12345678
```