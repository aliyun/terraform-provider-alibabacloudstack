---
subcategory: "Cloud-Native Distributed Database PolarDB-X 2.0"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_polardbx_log_engine"
description: |-
  Provides a PolarDBX log engine configuration.
---

# alibabacloudstack_polardbx_log_engine

Provides configuration management for PolarDBX log engines.

## Example Usage
```
variable "name" {
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




variable "existed_polardbx_id" {
  type    = string
  default = ""
}

data "alibabacloudstack_polardbx_instance_types" "cn" {
  sorted_by = "CPU"
  spec_type = "CN"
}

data "alibabacloudstack_polardbx_instance_types" "dn" {
  sorted_by = "CPU"
  spec_type = "DN"
}

data "alibabacloudstack_polardbx_instances" "default" {
  ids = var.existed_polardbx_id == "" ? [" ", ] : ["${var.existed_polardbx_id}", ]
}

resource "alibabacloudstack_polardbx_instance" "default" {
  count          = length(data.alibabacloudstack_polardbx_instances.default.polardbx_instances) == 0 ? 1 : 0
  zone_id        = data.alibabacloudstack_zones.default.zones.0.id
  engine_version = "5.7"
  storage        = 50
  vswitch_id     = alibabacloudstack_vpc_vswitch.default.id
  cn_node_class  = data.alibabacloudstack_polardbx_instance_types.cn.instance_types.0.id
  cn_node_count  = "2"
  dn_node_class  = data.alibabacloudstack_polardbx_instance_types.dn.instance_types.0.id
  dn_node_count  = "2"
}
locals {
  polardbx_instance = length(data.alibabacloudstack_polardbx_instances.default.polardbx_instances) == 0 ? alibabacloudstack_polardbx_instance.default.0 : data.alibabacloudstack_polardbx_instances.default.polardbx_instances.0
}




data "alibabacloudstack_polardbx_cdc_classes" "default" {
  instance_id = local.polardbx_instance.id
  sorted_by   = "CPU"
}

resource "alibabacloudstack_polardbx_log_engine" "default" {
  cdc_node_class = data.alibabacloudstack_polardbx_cdc_classes.default.cdc_classes.0.id
  cdc_node_count = 2
  multi_stream {
    hash_level = "RECORD"
    node_class = data.alibabacloudstack_polardbx_cdc_classes.default.cdc_classes.0.id
    node_count = 3
    group_name = "group83t1"
    comment    = "test1"
  }
  multi_stream {
    group_name = "group83t2"
    comment    = "test2"
    hash_level = "RECORD"
    node_class = data.alibabacloudstack_polardbx_cdc_classes.default.cdc_classes.1.id
    node_count = 2
  }

  instance_id = local.polardbx_instance.id
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required) ID of the PolarDBX instance.
* `cdc_node_class` - (Optional) Specification for log engine nodes.
* `cdc_node_count` - (Optional) Number of log engine nodes.
* `multi_stream` - (Optional) Multi-stream configuration block containing:
  * `group_name` - (Required) Name of the stream group.
  * `comment` - (Optional) Description of the configuration.
  * `hash_level` - (Required) Sharding level for data distribution.
  * `node_class` - (Required) Specification for stream group nodes.
  * `node_count` - (Required) Number of nodes in the stream group.

## Attributes Reference

The following attributes are exported in addition to the arguments above:

* `multi_stream` - Multi-stream configuration containing:
  * `instance_name` - Name of the stream group instance.
