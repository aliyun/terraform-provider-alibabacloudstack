---
subcategory: "Graph Database(GDB)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_graph_database_db_instance"
sidebar_current: "docs-alibabacloudstack-resource-graph-database-db-instance"
description: |-
  Provides a Alibabacloudstack Graph Database Db Instance resource.
---

# alibabacloudstack\_graph\_database\_db\_instance

Provides a Graph Database Db Instance resource.

For information about Graph Database Db Instance and how to use it, see [What is Db Instance](https://help.aliyun.com/document_detail/102865.html).



## Example Usage

Basic Usage

```terraform
resource "alibabacloudstack_graph_database_db_instance" "example" {
  db_node_class            = "gdb.r.2xlarge"
  db_instance_network_type = "vpc"
  db_version               = "1.0"
  db_instance_category     = "HA"
  db_instance_storage_type = "cloud_ssd"
  db_node_storage          = "example_value"
  payment_type             = "PayAsYouGo"
  db_instance_description  = "example_value"
}

```

## Argument Reference

The following arguments are supported:

* `vpc_id` - (Required，ForceNew) The vpc id of the vpc authorization.
* `vswitch_id` - (Optional, ForceNew) The vswitch id.
* `zone_id` - (Optional, Computed, ForceNew) The zone ID of the resource.
* `db_instance_category` - (Required, ForceNew) The category of the db instance. Valid values: `HA`.
* `db_instance_description` - (Optional) According to the practical example or notes.
* `db_instance_network_type` - (Required, ForceNew) The network type of the db instance. Valid values: `vpc`.
* `db_instance_storage_type` - (Required) Disk storage type. Valid values: `cloud_essd`, `cloud_ssd`. Modification is not supported.
* `db_node_class` - (Required) The class of the db node. Valid values: `gdb.r.xlarge`, `gdb.r.2xlarge`, `gdb.r.4xlarge`, `gdb.r.8xlarge`, `gdb.r.16xlarge`.
* `db_node_storage` - (Required) Instance storage space, which is measured in GB.
* `db_version` - (Required, ForceNew) Kernel Version. Valid values: `1.0` or `1.0-OpenCypher`. `1.0`: represented as gremlin, `1.0-OpenCypher`: said opencypher.
* `payment_type` - (Required, ForceNew) The paymen type of the resource. Valid values: `PayAsYouGo`.
* `db_instance_ip_array` - (Optional, Computed) IP ADDRESS whitelist for the instance group list. See the following `Block db_instance_ip_array`.
  * `db_instance_ip_array_attribute` - (Optional) The default is empty. To distinguish between the different property console does not display a `hidden` label grouping.
  * `db_instance_ip_array_name` - (Optional) IP ADDRESS whitelist group name.
  * `security_ips` - (Optional) IP ADDRESS whitelist addresses in the IP ADDRESS list, and a maximum of 1000 comma-separated format is as follows: `0.0.0.0/0` and `10.23.12.24`(IP) or `10.23.12.24/24`(CIDR mode, CIDR (Classless Inter-Domain Routing)/24 represents the address prefixes in the length of the range [1,32]).

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Db Instance.
* `status` - Instance status. Value range: `Creating`, `Running`, `Deleting`, `Rebooting`, `DBInstanceClassChanging`, `NetAddressCreating` and `NetAddressDeleting`.

## Import

Graph Database Db Instance can be imported using the id, e.g.

```
$ terraform import alibabacloudstack_graph_database_db_instance.example <id>
```
