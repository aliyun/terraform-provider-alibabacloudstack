---
subcategory: "PolarDB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_polardb_cluster_database"
sidebar_current: "docs-alibabacloudstack-resource-polardb-cluster-database"
description: |-
  提供 PolarDB 集群数据库资源。
---

# alibabacloudstack_polardb_cluster_database

提供 PolarDB 集群数据库资源。该资源允许您管理 PolarDB 集群中的数据库。

## 示例用法

```hcl
variable "name" {
  default = "tf-polar-test"
}

data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
  enable_details = true
}

resource "alibabacloudstack_vpc_vpc" "default" {
  vpc_name       = "${var.name}_vpc"
  cidr_block     = "172.16.0.0/16"
}

resource "alibabacloudstack_vpc_vswitch" "default" {
  vpc_id       = "${alibabacloudstack_vpc_vpc.default.id}"
  cidr_block   = "172.16.0.0/24"
  zone_id      = "${data.alibabacloudstack_zones.default.zones.0.id}"
  vswitch_name = "${var.name}_vsw"
}

resource "alibabacloudstack_polardb_cluster_instance" "default" {
  db_cluster_description = "${var.name}"
  zone_id                = "${data.alibabacloudstack_zones.default.zones.0.id}"
  db_type                = "${var.db_type}"
  db_version             = "${var.db_version}"
  storage_space          = "20"
  vpc_id                 = "${alibabacloudstack_vpc_vpc.default.id}"
  vswitch_id             = "${alibabacloudstack_vpc_vswitch.default.id}"
  db_node_class          = "${data.alibabacloudstack_polardb_cluster_instance_types.default.instance_types.0.id}"
  sub_category           = "${data.alibabacloudstack_polardb_cluster_instance_types.default.instance_types.0.sub_category}"
  storage_type           = "ESSDPL1"
}

resource "alibabacloudstack_polardb_cluster_database" "default" {
  db_cluster_id      = "${alibabacloudstack_polardb_cluster_instance.default.id}"
  db_name            = "${var.name}"
  character_set_name = "utf8"
  db_description     = "Test database for PolarDB cluster"
}
```

## 参数参考

支持以下参数：

* `db_cluster_id` - (必填, 变更后新建) PolarDB 集群的 ID。
* `db_name` - (必填, 变更后新建) 数据库的名称。
* `character_set_name` - (必填, 变更后新建) 数据库的字符集。
* `db_description` - (可选) 数据库的描述。

## 属性参考

导出以下属性：

* `id` - 数据库的 ID，格式为 `{DBClusterId}:{DBName}`。
* `engine` - 数据库引擎。
* `db_status` - 数据库的状态。

## 导入

可以使用 id 导入 PolarDB 集群数据库，例如：

```bash
$ terraform import alibabacloudstack_polardb_cluster_database.example pc-12345678:test_database
```