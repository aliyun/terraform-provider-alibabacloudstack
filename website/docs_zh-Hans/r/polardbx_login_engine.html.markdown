---
subcategory: "PolarDB-X"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_polardbx_log_engine"
sidebar_current: "docs-Alibabacloudstack-polardbx-log-engine"
description: |-
  提供一个polardbx日志引擎配置。
---

# alibabacloudstack_polardbx_log_engine

提供一个polardbx日志引擎配置。

## 使用示例
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

## 参数参考

支持以下参数：
  * `instance_id` - (必填) - PolarDBX实例的ID。
  * `cdc_node_class` - (可选) - 日志引擎节点规格。
  * `cdc_node_count` - (可选) - 日志引擎节点数量。
  * `multi_stream` - (可选) - 开通多流配置。
    * "instance_name"
    * `group_name` - (必填) - 流组名称。
    * `comment` - (可选) - 描述。
    * `hash_level` - (必填) - 拆分级别。
    * `node_class` - (必填) - 流组节点规格。
    * `node_count` - (必填) - 流组节点数量。

## 属性参考

除了上述参数外，还导出以下属性：
  * `multi_stream` - 开通多流配置。
    * "instance_name" - 流组实例名。
