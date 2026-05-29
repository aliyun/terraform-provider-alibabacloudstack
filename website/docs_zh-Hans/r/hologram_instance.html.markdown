---
subcategory: "实时数仓 Hologres"
layout: "alibabacloudstack"
page_title: "AlibabacloudStack: alibabacloudstack_hologram_instance"
sidebar_current: "docs-alibabacloudstack-resource-hologram-instance"
description: |-
  提供阿里云专有云Hologram实例资源
---

# alibabacloudstack_hologram_instance

提供Hologram实例资源。

## 示例用法

### 基础用法

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

## 参数说明

以下参数是支持的：

* `compute_type` - （必选，不可变）实例类型。有效值：`Standard`，`Follower`。
* `zone_id` - （必选，不可变）实例的可用区ID。
* `cpu` - （可选，不可变）CPU品牌。默认值：`intel`。
* `node` - （必选）节点数量。
* `vpc_id` - （必选，不可变）专有网络ID。
* `vswitch_id` - （必选，不可变）交换机ID。
* `leader_instance_id` - （可选）主实例ID。当`compute_type`为`Follower`时必选。
* `instance_name` - （必选）实例名称。
* `cluster` - （必选，不可变）集群名称。

## 属性导出

以下属性会被导出：

* `id` - 实例ID。
* `instance_id` - 实例ID。
* `instance_status` - 实例状态。
* `creation_time` - 实例创建时间。
* `version` - 实例版本。
* `enable_hive_access` - 是否启用Hive访问。
* `endpoints` - 实例的访问端点。

### endpoints块

endpoints映射包含以下属性：

* `type` - 端点类型。
* `endpoint` - 端点地址。
* `enabled` - 是否启用端点。
* `vpc_id` - 端点的专有网络ID。
* `vswitch_id` - 端点的交换机ID。
* `vpc_instance_id` - 端点的专有网络实例ID。

## 导入

Hologram实例可以使用实例ID进行导入，例如：

```shell
$ terraform import alibabacloudstack_hologram_instance.example <id>
```