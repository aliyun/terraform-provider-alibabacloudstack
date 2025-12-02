---
subcategory: "PolarDBX"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_polardbx_instance"
sidebar_current: "docs-alibabacloudstack-polardbx-instance"
description: |-
  Provides a Polardbx instance resource.
---

# alibabacloudstack\_polardbx\_instance

Provides a Polardbx instance resource.

## Example Usage
```
variable "name" {
		default = "tf_acc_drds_polardb_16807"
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
  vswitch_id = "${alibabacloudstack_vpc_vswitch.default.id}"
  cn_node_count = "2"
  dn_node_class = "mysql.n4.medium.25"
  dn_node_count = "2"
  description = "testtf1111"
  storage = "50"
  cn_node_class = "polarx.x4.medium.2e"
  engine_version = "5.7"
  network_type = "vpc"
  zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
  vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
}
```

## Argument Reference

The following arguments are supported:
  * `series` - (Optional) - The series of PolarDB-X 2.0 instance.
  * `cpu_type` - (Optional) - The CPU type of PolarDB-X 2.0 instance.
  * `storage` - (Required) - The storage size of PolarDB-X 2.0 instance.
  * `spec_series` - (Optional) - The specification series of PolarDB-X 2.0 instance.
  * `cidr_block` - (Optional) - The CIDR block of PolarDB-X 2.0 instance.
  * `description` - (Optional) - The description of PolarDB-X 2.0 instance.
  * `cn_node_class` - (Required) - Class of computing nodes.
  * `cn_node_count` - (Required) - Number of computing nodes.
  * `create_time` - (Optional) - The creation time of the resource
  * `db_node_class` - (Optional) - Node specifications:-**polarx.x4.medium.2e**:2 cores 8g-**polarx.x4.large.2e**:4 core 16g-**polarx.x8.large.2e**:4 core 32g-**polarx.x4.xlarge.2e**:8 cores 32g-**polarx.x8.xlarge.2e**:8 cores 64g-**polarx.x4.2xlarge.2e**:16 cores 64g-**polarx.x8.2xlarge.2e**:16 cores 128g-**polarx.x4.4xlarge.2e**:32 core 128g-**polarx.x8.4xlarge.2e**:32 cores 256G-**polarx.st.8xlarge.2e**:60 cores 470g-**polarx.st.12xlarge.2e**:90 core 720g
  * `db_node_count` - (Optional) - The number of instance nodes. The minimum number is 2.
  * `dn_node_class` - (Required) - Class of computing nodes.
  * `dn_node_count` - (Required) - The number of storage nodes.
  * `gms_node_class` - (Optional) - Class of GMS nodes.
  * `is_read_db_instance` - (Optional) - If true, the instance is a read-only instance.
  * `polardbx_instance_id` - (Optional) - The ID of the PolarDB-X 2.0 instance.
  * `primary_zone` - (Optional) - Primary Availability Zone.
  * `resource_type` - (Optional) - Resource type. Currently, only one type of resource for PolarDB-X 2.0 instance is supported.
  * `secondary_zone` - (Optional) - Secondary availability zone.
  * `tertiary_zone` - (Optional) - Third Availability Zone.
  * `topology_type` - (Optional) - Topology type:-**3azones**: three available areas;-**1azone**: Single zone.
  * `vswitch_id` - (Required) -The VSwitch ID.
  * `zone_id` - (Optional) - The ID of the zone to which the instance belongs.
  * `compute_parameters` - (Optional) - The compute resource configuration of the instance.
    * `name` - (Required) - The name of the compute resource configuration.
    * `value` - (Required) - The value of the compute resource configuration.
  * `storage_parameters` - (Optional) - The storage configuration of the instance.
    * `name` - (Required) - The name of the storage configuration.
    * `value` - (Required) - The value of the storage configuration.
  * `security_groups` - (Optional) - The security ip group of the instance.
    * `group_name` - (Required) - The name of the security ip group.
    * `ips` - (Required) - The security ips. Multiple IPs need to be separated by ",", example: `192.168.0.1,192.168.0.0/24`

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
  * `series` - The series of the PolarDB-X instance.
  * `cpu_type` - The CPU type of the PolarDB-X instance.
  * `spec_series` - The series of the PolarDB-X instance.
  * `description` - The description of the PolarDB-X instance.
  * `create_time` - The creation time of the resource
  * `db_node_class` - Node specifications:-**polarx.x4.medium.2e**:2 cores 8g-**polarx.x4.large.2e**:4 core 16g-**polarx.x8.large.2e**:4 core 32g-**polarx.x4.xlarge.2e**:8 cores 32g-**polarx.x8.xlarge.2e**:8 cores 64g-**polarx.x4.2xlarge.2e**:16 cores 64g-**polarx.x8.2xlarge.2e**:16 cores 128g-**polarx.x4.4xlarge.2e**:32 core 128g-**polarx.x8.4xlarge.2e**:32 cores 256G-**polarx.st.8xlarge.2e**:60 cores 470g-**polarx.st.12xlarge.2e**:90 core 720g
  * `db_node_count` - The number of instance nodes. The minimum number is 2.
  * `status` - The status of the resource
