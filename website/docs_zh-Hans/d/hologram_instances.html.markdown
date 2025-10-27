---
subcategory: "Hologram"
layout: "alibabacloudstack"
page_title: "AlibabacloudStack: alibabacloudstack_hologram_instances"
sidebar_current: "docs-alibabacloudstack-datasource-hologram-instances"
description: |-
  提供Hologram实例列表
---

# alibabacloudstack_hologram_instances

该数据源根据指定的过滤条件提供阿里云账户中的Hologram实例列表。


## 示例用法

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

## 参数说明

以下参数是支持的：

* `ids` - （可选）实例ID列表。
* `name_regex` - （可选）用于按实例名称过滤结果的正则表达式。

## 属性导出

以下属性会被导出：

* `ids` - 实例ID列表。
* `instances` - Hologram实例列表。每个元素包含以下属性：
  * `id` - 实例ID。
  * `compute_type` - 实例类型。
  * `cpu` - CPU数量。
  * `node` - 节点数量。
  * `instance_name` - 实例名称。
  * `cluster` - 集群名称。
  * `instance_status` - 实例状态。
  * `creation_time` - 实例创建时间。
  * `version` - 实例版本。
  * `enable_hive_access` - 是否启用Hive访问。
  * `instance_type` - 实例类型。
  * `instance_charge_type` - 实例计费类型。
  * `cpu_arch` - CPU架构。
  * `cpu_brand` - CPU品牌。
  * `cpu_brand_i18n` - 国际化CPU品牌。
  * `apsara_ascm_cpu_brand` - ASCM CPU品牌。
  * `support_replica` - 是否支持副本。
  * `ascm_create_user` - 创建实例的用户。
  * `commodity_code` - 商品代码。
  * `endpoints` - 端点列表。每个端点包含：
    * `type` - 端点类型。
    * `endpoint` - 端点地址。
    * `enabled` - 是否启用端点。
    * `vpc_id` - 专有网络ID。
    * `vswitch_id` - 交换机ID。
    * `vpc_instance_id` - 专有网络实例ID。
