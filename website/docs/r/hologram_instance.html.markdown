---
subcategory: "Hologres"
layout: "alibabacloudstack"
page_title: "AlibabacloudStack: alibabacloudstack_hologram_instance"
sidebar_current: "docs-alibabacloudstack-resource-hologram-instance"
description: |-
  Provides a Alibaba Cloud Hologram Instance resource.
---

# alibabacloudstack_hologram_instance

Provides a Hologram Instance resource.


## Example Usage

### Basic Usage

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
  vswitch_id = "${alibabacloudstack_vpc_vswitch.default.id}"
}
```

## Argument Reference

The following arguments are supported:

* `compute_type` - (Required, ForceNew) The type of the instance. Valid values: `Standard`, `Follower`.
* `zone_id` - (Required, ForceNew) The zone ID of the instance.
* `cpu` - (Optional, ForceNew) The CPU brand. defaults to `intel`.
* `node` - (Required) The number of nodes.
* `vpc_id` - (Required, ForceNew) The ID of the VPC.
* `vswitch_id` - (Required, ForceNew) The ID of the vSwitch.
* `leader_instance_id` - (Optional) The ID of the leader instance. It is required when the `compute_type` is `Follower`.
* `instance_name` - (Required) The name of the instance.
* `cluster` - (Required, ForceNew) The cluster name.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the instance.
* `instance_id` - The ID of the instance.
* `instance_status` - The status of the instance.
* `creation_time` - The creation time of the instance.
* `version` - The version of the instance.
* `enable_hive_access` - Whether Hive access is enabled.
* `endpoints` - The endpoints of the instance.

### Block endpoints

The endpoints mapping supports the following:

* `type` - The type of endpoint.
* `endpoint` - The endpoint address.
* `enabled` - Whether the endpoint is enabled.
* `vpc_id` - The VPC ID of the endpoint.
* `vswitch_id` - The VSwitch ID of the endpoint.
* `vpc_instance_id` - The VPC instance ID of the endpoint.

## Import

Hologram Instance can be imported using the id, e.g.

```shell
$ terraform import alibabacloudstack_hologram_instance.example <id>
```
