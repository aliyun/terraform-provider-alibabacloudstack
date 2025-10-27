---
subcategory: "Hologram"
layout: "alibabacloudstack"
page_title: "AlibabacloudStack: alibabacloudstack_hologram_instances"
sidebar_current: "docs-alibabacloudstack-datasource-hologram-instances"
description: |-
  Provides a list of Hologram Instances to the user.
---

# alibabacloudstack_hologram_instances

This data source provides a list of Hologram Instances in an Alibaba Cloud account according to the specified filters.

## Example Usage

```hcl

variable "name" {
  default = "tf-testacc"
}

data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
  enable_details = true
}

data "alibabacloudstack_hologram_clusters" "default" {
  zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
}

resource "alibabacloudstack_vpc_vpc" "default" {
  vpc_name = "${var.name}_vpc"
  cidr_block = "172.16.0.0/16"
  lifecycle {
      ignore_changes = [
        secondary_cidr_blocks,
        tags
      ]
  }
}

resource "alibabacloudstack_vpc_vswitch" "default" {
  vswitch_name = "${var.name}_vsw"
  vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
  cidr_block = "172.16.1.0/24"
  zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
  lifecycle {
      ignore_changes = [
        tags
      ]
  }
}

resource "alibabacloudstack_hologram_instance" "example" {
  zone_id =  "${data.alibabacloudstack_zones.default.zones.0.id}"
  instance_name = "${var.name}"
  compute_type = "Standard"
  cpu = "intel"
  node = 2
  cluster = "${data.alibabacloudstack_hologram_clusters.default.clusters.0.id}"
  vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
  vswitch_id = "${data.alibabacloudstack_vpc_vswitch.default.id}"
}

data "alibabacloudstack_hologram_instances" "example" {
  ids = ["${alibabacloudstack_hologram_instance.default.id}"]
}

output "first_instance_id" {
  value = data.alibabacloudstack_hologram_instances.example.instances.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of instance IDs.
* `name_regex` - (Optional) A regex string to filter results by instance name.

## Attributes Reference

The following attributes are exported:

* `ids` - A list of instance IDs.
* `instances` - A list of Hologram instances. Each element contains the following attributes:
  * `id` - ID of the instance.
  * `compute_type` - Type of the instance.
  * `cpu` - CPU count.
  * `node` - Number of nodes.
  * `instance_name` - Name of the instance.
  * `cluster` - Cluster name.
  * `instance_status` - Status of the instance.
  * `creation_time` - Creation time of the instance.
  * `version` - Version of the instance.
  * `enable_hive_access` - Whether Hive access is enabled.
  * `instance_type` - Type of the instance.
  * `instance_charge_type` - Charge type of the instance.
  * `cpu_arch` - CPU architecture.
  * `cpu_brand` - CPU brand.
  * `cpu_brand_i18n` - Internationalized CPU brand.
  * `apsara_ascm_cpu_brand` - ASCM CPU brand.
  * `support_replica` - Whether replica is supported.
  * `ascm_create_user` - User who created the instance.
  * `commodity_code` - Commodity code.
  * `endpoints` - List of endpoints. Each endpoint contains:
    * `type` - Type of the endpoint.
    * `endpoint` - Endpoint address.
    * `enabled` - Whether the endpoint is enabled.
    * `vpc_id` - VPC ID.
    * `vswitch_id` - VSwitch ID.
    * `vpc_instance_id` - VPC instance ID.