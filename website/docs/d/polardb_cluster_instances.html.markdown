---
subcategory: "PolarDB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_polardb_cluster_instances"
sidebar_current: "docs-alibabacloudstack-datasource-polardb-cluster-instances"
description: |-
  Provides a list of PolarDB cluster instances.
---

# alibabacloudstack_polardb_cluster_instances

This data source provides a list of PolarDB cluster instances in an Alibaba Cloud Stack environment.

## Example Usage

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
  db_version 			= "${var.db_version}"
  storage_space 		= "20"
  vpc_id 				= "${alibabacloudstack_vpc_vpc.default.id}"
  vswitch_id			= "${alibabacloudstack_vpc_vswitch.default.id}"
  db_node_class 		= "${data.alibabacloudstack_polardb_cluster_instance_types.default.instance_types.0.id}"
  sub_category 			= "${data.alibabacloudstack_polardb_cluster_instance_types.default.instance_types.0.sub_category}"
  storage_type 			= "ESSDPL1"
}
data "alibabacloudstack_polardb_cluster_instances" "default" {
  description_regex = "test"
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of instance IDs.
* `name_regex` - (Optional, Deprecated) A regex string to filter results by instance description. **Field 'name_regex' is deprecated and will be removed in a future release. Please use new field 'description_regex' instead.**
* `description_regex` - (Optional) A regex string to filter results by instance description.

## Attributes Reference

The following attributes are exported:

* `ids` - A list of instance IDs.
* `db_cluster_instances` - A list of PolarDB cluster instances. Each element contains the following attributes:
  * `deletion_lock` - The deletion lock status.
  * `category` - The category of the cluster instance.
  * `resource_group_id` - The ID of the resource group.
  * `storage_pay_type` - The storage payment type.
  * `hot_standby_cluster` - The hot standby cluster status.
  * `db_cluster_id` - The ID of the cluster instance.
  * `db_type` - The database type.
  * `db_cluster_network_type` - The network type of the cluster instance.
  * `db_proxy_cluster_class` - The proxy cluster class.
  * `db_version` - The database version.
  * `department` - The department information.
  * `db_nodes` - The database nodes information.
    * `zone_id` - The zone ID of the node.
    * `db_node_role` - The role of the node.
    * `db_node_id` - The ID of the node.
    * `region_id` - The region ID of the node.
    * `db_node_class` - The class of the node.
  * `engine` - The database engine.
  * `tags` - The tags of the cluster instance.
    * `tag_key` - The tag key.
    * `tag_value` - The tag value.
  * `resource_group` - The resource group information.
  * `architecture` - The architecture of the cluster instance.
  * `zone_id` - The zone ID of the cluster instance.
  * `db_cluster_status` - The status of the cluster instance.
  * `create_time` - The creation time of the cluster instance.
  * `db_cluster_description` - The description of the cluster instance.
  * `expired` - Whether the cluster instance is expired.
  * `pay_type` - The payment type of the cluster instance.
  * `lock_mode` - The lock mode of the cluster instance.
  * `vswitch_id` - The ID of the vSwitch.
  * `ascm_create_user` - The ASCM create user.
  * `db_node_class` - The node class of the cluster instance.
  * `storage_used` - The used storage space.
  * `db_node_number` - The number of database nodes.
  * `vpc_id` - The ID of the VPC.
  * `storage_space` - The storage space of the cluster instance.
  * `serverless_type` - The serverless type.
  * `department_name` - The department name.
  * `region_id` - The region ID of the cluster instance.
  * `expire_time` - The expiration time of the cluster instance.
  * `vip` - The VIP of the cluster instance.
  * `resource_group_name` - The name of the resource group.
```