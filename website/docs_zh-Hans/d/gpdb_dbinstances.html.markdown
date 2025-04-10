---
subcategory: "GPDB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_gpdb_dbinstances"
sidebar_current: "docs-Alibabacloudstack-datasource-gpdb-dbinstances"
description: |- 
  查询图数据库实例
---

# alibabacloudstack_gpdb_dbinstances
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_gpdb_instances`

根据指定过滤条件列出当前凭证权限可以访问的图数据库实例列表。

## 示例用法

```hcl
data "alibabacloudstack_zones" "default" {}

resource "alibabacloudstack_vpc" "default" {
	name = "testing"
	cidr_block = "10.0.0.0/8"
}

data "alibabacloudstack_gpdb_dbinstances" "example" {
  name_regex        = "gp-.+\\d+"
  availability_zone = data.alibabacloudstack_zones.default.zones.0.id
  vswitch_id        = alibabacloudstack_vswitch.default.id
  ids               = ["db-1234567890abcdefg"]
  output_file       = "dbinstances.txt"
}

resource "alibabacloudstack_vswitch" "default" {
	vpc_id         = alibabacloudstack_vpc.default.id
	cidr_block     = "10.1.0.0/16"
	name           = "apsara_vswitch"
	availability_zone = data.alibabacloudstack_zones.default.zones.0.id
}

resource "alibabacloudstack_gpdb_instance" "default" {
	vswitch_id           = alibabacloudstack_vswitch.default.id
	engine               = "gpdb"
	engine_version       = "6.0"
	instance_class       = "gpdb.group.segsdx2"
	instance_group_count = 2
	description          = "testing_01"
}

output "dbinstance_id" {
  value = "${data.alibabacloudstack_gpdb_dbinstances.example.instances.0.id}"
}
```

## 参数说明

以下参数是支持的：

* `name_regex` - (选填) 用于按名称筛选 GPDB 数据库实例的正则表达式字符串。
* `availability_zone` - (选填) 实例所在的可用区。可以通过 `data "alibabacloudstack_zones"` 获取可用区信息。
* `vswitch_id` - (选填) 用于检索属于指定 VSwitch 资源的实例。VSwitch 的 ID 可以通过 `alibabacloudstack_vswitch` 资源获取。
* `ids` - (选填) GPDB 数据库实例 ID 列表，用于精确匹配特定实例。

## 属性说明

除了上述参数外，还导出以下属性：

* `names` - 匹配的 GPDB 数据库实例名称列表。
* `ids` - 匹配的 GPDB 数据库实例 ID 列表。
* `instances` - GPDB 数据库实例列表。每个元素包含以下属性：
  * `id` - GPDB 数据库实例的 ID。
  * `description` - GPDB 数据库实例的描述。
  * `region_id` - GPDB 数据库实例所在的区域 ID。
  * `availability_zone` - GPDB 数据库实例所在的可用区。
  * `creation_time` - GPDB 数据库实例的创建时间（UTC 格式：YYYY-MM-DDThh:mm:ssZ）。
  * `status` - GPDB 数据库实例的当前状态。
  * `engine` - 数据库引擎类型。支持值为 `gpdb`。
  * `engine_version` - 数据库引擎版本。支持值包括 `6.0` 和 `7.0`。
  * `instance_class` - GPDB 数据库实例的规格。
  * `instance_group_count` - GPDB 数据库实例中的分组数量。
  * `instance_network_type` - GPDB 数据库实例的网络类型。支持值为 `VPC`。
  * `charge_type` - GPDB 数据库实例的计费方式。可能的值为 `PrePaid`（包年包月）和 `PostPaid`（按量付费）。