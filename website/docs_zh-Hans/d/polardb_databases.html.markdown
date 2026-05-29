---
subcategory: "云原生数据库 PolarDB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_polardb_databases"
sidebar_current: "docs-Alibabacloudstack-datasource-polardb-databases"
description: |-
  提供阿里云账号下拥有的polardb databases列表。
---

# alibabacloudstack\_polardb\_databases

此数据源提供根据指定过滤条件列出的阿里云账号下的polardb databases资源列表。

## 示例用法
```
variable "name" {
	default = "tf-testAccPolardbBackups14307"
}

data "alibabacloudstack_zones" default {
  available_resource_creation = "VSwitch"
  enable_details = true
}

resource "alibabacloudstack_polardb_dbinstance" "instance" {
	engine            = "MySQL"
	engine_version    = "5.7"
	instance_name = "${var.name}"
	db_instance_storage_type= "local_ssd"
	db_instance_storage = 5
	db_instance_class = "rds.mysql.t1.small"
	zone_id= "${data.alibabacloudstack_zones.default.zones.0.id}"
}
resource "alibabacloudstack_polardb_database" "default" {
	data_base_instance_id = "${alibabacloudstack_polardb_dbinstance.instance.id}"
	data_base_description = "自动化生成测试"
	data_base_name        = "tftest"
	character_set_name = "utf8"
}

data "alibabacloudstack_polardb_databases" "default" {
  ids = [
          "${alibabacloudstack_polardb_database.default.id}"
        ]
  data_base_instance_id = "${alibabacloudstack_polardb_dbinstance.instance.id}"
}
```

## 参数参考
以下参数是支持的：
  * `ids` - (选填) - 用于筛选的PolarDB数据库ID列表。
  * `data_base_instance_id` - (必填) - 数据库实例ID。
  * `data_base_name` - (选填) - 数据库名称
  * `name_regex` - (选填) - 用于筛选的PolarDB数据库名称。
  * `description_regex` - (选填) - 用于筛选的PolarDB数据库描述信息。

## Attributes Reference
除了上述参数外，还导出以下属性：
  * `databases` - 数据库列表。
    * `id` - 数据库ID。
    * `accounts` - 数据库账号信息详情。> 当集群为PolarDB MySQL引擎时，不含高权限账号。
    * `character_set_name` - 字符集，详情请参见[字符集表](~~99716~~)。  
    * `data_base_description` - 数据库描述信息。
    * `data_base_instance_id` - 数据库实例ID.
    * `data_base_name` - 数据库名称。
    * `engine` - 数据库引擎类型，取值范围如下：* **MySQL*** **Oracle*** **PostgreSQL**
    * `page_number` - 页码。
    * `page_size` - 每页记录数，取值范围如下：* **30*** **50*** **100**默认值为**30**。
    * `status` - 代表资源状态的资源属性字段
