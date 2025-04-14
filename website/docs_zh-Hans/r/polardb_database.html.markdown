---
subcategory: "PolarDB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_polardb_database"
sidebar_current: "docs-Alibabacloudstack-polardb-database"
description: |- 
  编排polardb数据库表
---

# alibabacloudstack_polardb_database

使用Provider配置的凭证在指定的资源集编排polardb数据库表。

## 示例用法

```hcl
variable "name" {
    default = "tf-testaccdbdatabase_basic"
}

data "alibabacloudstack_zones" default {
  available_resource_creation = "VSwitch"
  enable_details = true
}

resource "alibabacloudstack_vpc_vpc" "default" {
  vpc_name   = "${var.name}_vpc"
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vpc_vswitch" "default" {
  name       = "${var.name}_vsw"
  vpc_id     = "${alibabacloudstack_vpc_vpc.default.id}"
  cidr_block = "172.16.0.0/24"
  zone_id    = "${data.alibabacloudstack_zones.default.zones.0.id}"
}

resource "alibabacloudstack_polardb_instance" "instance" {
  engine                  = "MySQL"
  engine_version          = "5.7"
  instance_name           = "${var.name}"
  db_instance_storage_type= "local_ssd"
  db_instance_storage     = 5
  db_instance_class       = "rds.mysql.t1.small"
  zone_id                = "${data.alibabacloudstack_zones.default.zones.0.id}"
  vswitch_id             = "${alibabacloudstack_vpc_vswitch.default.id}"
}

resource "alibabacloudstack_polardb_database" "default" {
  data_base_name         = "tf-testaccdbdatabase_basic"
  character_set_name     = "utf8"
  data_base_instance_id  = "${alibabacloudstack_polardb_instance.instance.id}"

  accounts {
    account               = "test_account"
    account_privilege     = "ReadWrite"
    account_privilege_detail = "Full access to the database"
  }

  data_base_description   = "This is a test database"
}
```

## 参数说明

支持以下参数：

  * `accounts` - (选填) - 数据库账号信息详情。当集群为PolarDB MySQL引擎时，不含高权限账号。
    
    * `account` - (选填) - 数据库账户名称。
    
    * `account_privilege` - (选填) - 账户权限。取值范围如下：
      * **ReadWrite**: 读写权限
      * **ReadOnly**: 只读权限
      * **DMLOnly**: 仅允许DML操作
      * **DDLOnly**: 仅允许DDL操作
      * **ReadIndex**: 只读+索引
    
    * `account_privilege_detail` - (选填) - 账户权限详细信息。

  * `character_set_name` - (必填) - 字符集，详情请参见[字符集表](~~99716~~)。

  * `data_base_description` - (选填) - 数据库的描述。

  * `data_base_instance_id` - (必填) - 将关联的PolarDB实例ID。

  * `data_base_name` - (必填) - 数据库名称。

  * `engine` - (选填) - 数据库引擎类型，取值范围如下：
    * **MySQL**
    * **Oracle**
    * **PostgreSQL**

  * `status` - (选填) - 资源状态。

## 属性说明

除了上述所有参数外，还导出了以下属性：

  * `accounts` - 数据库账号信息详情。当集群为PolarDB MySQL引擎时，不含高权限账号。
    * `account` - 数据库账户名称。
    * `account_privilege` - 账号权限，取值范围如下：
      * **ReadWrite**: 读写
      * **ReadOnly**: 只读
      * **DMLOnly**: 仅允许DML
      * **DDLOnly**: 仅允许DDL
      * **ReadIndex**: 只读+索引
    * `account_privilege_detail` - 账户权限详细信息。

  * `data_base_description` - 数据库的描述。

  * `engine` - 数据库引擎类型，取值范围如下：
    * **MySQL**
    * **Oracle**
    * **PostgreSQL**

  * `status` - 资源状态。
  * `data_base_instance_id` - 关联的PolarDB实例ID。