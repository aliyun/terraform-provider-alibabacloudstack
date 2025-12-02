---
subcategory: "PolarDBX"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_polardbx_instances"
sidebar_current: "docs-Alibabacloudstack-datasource-polardbx-instances"
description: |-
  Provides a list of polardbx instances owned by an alibabacloudstack account.
---

# alibabacloudstack\_polardbx\_instances

This data source provides a list of polardbx instances in an alibabacloudstack account according to the specified filters.

## Example Usage
```
variable "name" {
  default = "tf-testAccPolardbxInstancesDataSource-3625795"
}

data "alibabacloudstack_zones" default {
  available_resource_creation = "VSwitch"
  enable_details = true
}

resource "alibabacloudstack_vpc_vpc" "default" {
  vpc_name = "${var.name}_vpc"
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vpc_vswitch" "default" {
  name = "${var.name}_vsw"
  vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
  cidr_block = "172.16.1.0/24"
  zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
}

resource "alibabacloudstack_polardbx_instance" "default" {
  description = "testtf1111"
	series = "enterprise"
	topology_type = "1azone"
	zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
	engine_version = "5.7"
	storage = "50"
	network_type = "vpc"
	vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
	vswitch_id = "${alibabacloudstack_vpc_vswitch.default.id}"
	cn_node_class = "polarx.x4.medium.2e"
	cn_node_count = "2"
	dn_node_class = "mysql.n4.medium.25"
	dn_node_count = "2"
}

data "alibabacloudstack_polardbx_instances" "default" {
  description_regex = "${alibabacloudstack_polardbx_instance.default.description}"
}
```

## Argument Reference

The following arguments are supported:
  * `ids` - (Optional) - A list of instance IDs to filter results.
  * `name_regex` - (Optional) - A name Regex of instance.
  * `description_regex` - (Optional) - A description Regex of instance.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
  * `polardbx_instances` - The list of PolarDB-X 2.0 instances.
    * `id` - The ID of the PolarDB-X 2.0 instance.
    * `series` - The series of the PolarDB-X 2.0 instance.
    * `storage` - The storage space of the PolarDB-X 2.0 instance.
    * `description` - The description of the PolarDB-X 2.0 instance.
    * `cpu_type` - The CPU type of the PolarDB-X 2.0 instance.
    * `cn_node_class` - The class of the computing node.
    * `cn_node_count` - Number of computing nodes.
    * `create_time` - The creation time of the resource
    * `db_node_class` - Node specifications:-**polarx.x4.medium.2e**:2 cores 8g-**polarx.x4.large.2e**:4 core 16g-**polarx.x8.large.2e**:4 core 32g-**polarx.x4.xlarge.2e**:8 cores 32g-**polarx.x8.xlarge.2e**:8 cores 64g-**polarx.x4.2xlarge.2e**:16 cores 64g-**polarx.x8.2xlarge.2e**:16 cores 128g-**polarx.x4.4xlarge.2e**:32 core 128g-**polarx.x8.4xlarge.2e**:32 cores 256G-**polarx.st.8xlarge.2e**:60 cores 470g-**polarx.st.12xlarge.2e**:90 core 720g
    * `db_node_count` - The number of instance nodes. The minimum number is 2.
    * `dn_node_class` - The node specification of storage nodes.
    * `dn_node_count` - The number of storage nodes.
    * `engine_version` - Fixed as 2.0 and cannot be changed.
    * `is_read_db_instance` - Is Read-Only Instance
    * `network_type` - The network type. Only the VPC network is supported.
    * `polardbx_instance_id` - The ID of the PolarDB-X 2.0 instance.
    * `primary_zone` - Primary Availability Zone.
    * `resource_type` - Resource type. Currently, only one type of resource for PolarDB-X 2.0 instance is supported.
    * `secondary_zone` - Secondary availability zone.
    * `status` - The status of the resource
    * `tertiary_zone` - Third Availability Zone.
    * `topology_type` - Topology type:-**3azones**: three available areas;-**1azone**: Single zone.
    * `vswitch_id` - The VSwitch ID.
    * `vpc_id` - The VPC ID.
    * `zone` - Instance availability zone.
