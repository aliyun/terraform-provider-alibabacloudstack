---
subcategory: "PolarDB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_polardb_cluster_databases"
sidebar_current: "docs-alibabacloudstack-datasource-polardb-cluster-databases"
description: |-
  提供 PolarDB 集群数据库列表。
---

# alibabacloudstack_polardb_cluster_databases

本数据源提供阿里云平台环境中的 PolarDB 集群数据库列表。

## 示例用法

```hcl
variable "name" {
	default = "tfAccountsName"
}

variable "password" {
}

data  "alibabacloudstack_zones" "default" {
	available_resource_creation = "PolarDB"
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

resource "alibabacloudstack_polardb_cluster_database" "default" {
    db_cluster_id		= "${alibabacloudstack_polardb_cluster_instance.default.id}"
    db_name 			= "${var.name}"
    character_set_name 	= "utf8"
}

data "alibabacloudstack_polardb_cluster_databases" "example" {
  db_cluster_id = "pc-12345678"
  db_name       = "test_database"
}
```

## 参数参考

支持以下参数：

* `ids` - (可选) 数据库 ID 列表。
* `db_name` - (可选) 数据库名称。
* `db_cluster_id` - (必填) PolarDB 集群的 ID。
* `name_regex` - (可选, 已弃用) 用于通过数据库名称过滤结果的正则表达式。**字段 'name_regex' 已弃用，并将在未来版本中移除。请使用新字段 'description_regex' 替代。**
* `description_regex` - (可选) 用于通过数据库名称过滤结果的正则表达式。

## 属性参考

导出以下属性：

* `ids` - 数据库 ID 列表。
* `databases` - PolarDB 集群数据库列表。每个元素包含以下属性：
  * `id` - 数据库的 ID。
  * `db_name` - 数据库名称。
  * `character_set_name` - 数据库的字符集。
  * `db_description` - 数据库的描述。
  * `engine` - 数据库引擎。
  * `db_status` - 数据库的状态。