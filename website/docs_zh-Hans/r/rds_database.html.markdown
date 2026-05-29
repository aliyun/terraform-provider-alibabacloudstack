---
subcategory: "云数据库 RDS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_rds_database"
sidebar_current: "docs-Alibabacloudstack-resource-rds-database"
description: |- 
  编排RDS数据库
---

# alibabacloudstack_rds_database
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_db_database`

使用Provider配置的凭证在指定的资源集编排RDS数据库。

## 示例用法

```hcl
variable "name" {
    default = "tf-testAccDBdatabase_basic"
}

data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
  enable_details = true
}

resource "alibabacloudstack_vpc_vpc" "default" {
  vpc_name = "${var.name}_vpc"
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vpc_vswitch" "default" {
  name = "${var.name}_vsw"
  vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
  cidr_block = "172.16.0.0/24"
  zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
}

resource "alibabacloudstack_db_instance" "instance" {
     engine               = "MySQL"
     engine_version       = "5.6"
     instance_type        = "rds.mysql.s2.large"
     instance_storage     = "30"
     vswitch_id           = "${alibabacloudstack_vpc_vswitch.default.id}"
     instance_name        = "${var.name}"
     storage_type         = "local_ssd"
}

resource "alibabacloudstack_rds_database" "default" {
  data_base_instance_id = "${alibabacloudstack_db_instance.instance.id}"
  data_base_name        = "tftestdatabase"
  character_set_name    = "utf8"
  data_base_description = "This is a test database"
}
```

## 参数说明

支持以下参数：

* `data_base_instance_id` - (必填，变更时重建) 数据库所属实例的 ID。可调用 `DescribeDBInstances` 获取。修改此参数将强制创建新资源。
* `data_base_name` - (必填，变更时重建) 数据库的名称。长度为 2~64 个字符，以字母开头，以字母或数字结尾，由小写字母、数字、下划线或中划线组成。数据库名称在实例内必须是唯一的。修改此参数将强制创建新资源。
* `character_set_name` - (必填，变更时重建) 数据库的字符集。支持的值取决于数据库引擎：
  - **MySQL/MariaDB**: `utf8`, `gbk`, `latin1`, `utf8mb4`。
  - **SQLServer**: `Chinese_PRC_CI_AS`, `Chinese_PRC_CS_AS`, `SQL_Latin1_General_CP1_CI_AS`, `SQL_Latin1_General_CP1_CS_AS`, `Chinese_PRC_BIN`。
  - **PostgreSQL**: 需要配置字符集、Collate 和 Ctype，格式为`字符集,<Collate>,<Ctype>`，例如：`UTF8,C,en_US.utf8`。字符集取值包括：`KOI8U`, `UTF8`, `WIN866`, `WIN874`, `WIN1250`, `WIN1251`, `WIN1252`, `WIN1253`, `WIN1254`, `WIN1255`, `WIN1256`, `WIN1257`, `WIN1258`, `EUC_CN`, `EUC_KR`, `EUC_TW`, `EUC_JP`, `EUC_JIS_2004`, `KOI8R`, `MULE_INTERNAL`, `LATIN1`, `LATIN2`, `LATIN3`, `LATIN4`, `LATIN5`, `LATIN6`, `LATIN7`, `LATIN8`, `LATIN9`, `LATIN10`, `ISO_8859_5`, `ISO_8859_6`, `ISO_8859_7`, `ISO_8859_8`, `SQL_ASCII`。
* `data_base_description` - (可选) 数据库的描述。长度为 2~256 个字符。必须以中文字符或英文字母开头，可以包括中文和英文字符、下划线 (`_`)、短横线 (`-`)。不能以 `http://` 或 `https://` 开头。
* `instance_id` - (**已废弃**，变更时重建) 该字段已废弃，请使用 `data_base_instance_id` 替代。数据库所属实例的 ID。
* `name` - (**已废弃**，变更时重建) 该字段已废弃，请使用 `data_base_name` 替代。数据库的名称。
* `character_set` - (**已废弃**，变更时重建) 该字段已废弃，请使用 `character_set_name` 替代。数据库的字符集。
* `description` - (**已废弃**) 该字段已废弃，请使用 `data_base_description` 替代。数据库的描述。

-> **注意：** `data_base_instance_id`、`data_base_name` 和 `character_set_name` 字段为必填项，创建后不支持修改。修改这些参数将强制创建新资源。

## 属性说明

除了上述所有参数外，还导出了以下属性：

* `id` - 数据库资源的唯一标识符。它由实例 ID 和数据库名称组成，格式为 `<instance_id>:<db_name>`。
* `data_base_instance_id` - 创建数据库的 RDS 实例 ID。
* `data_base_name` - 数据库的名称。
* `character_set_name` - 数据库的字符集。
* `data_base_description` - 数据库的描述。
* `instance_id` - 创建数据库的 RDS 实例 ID。
* `name` - 数据库的名称。
* `character_set` - 数据库的字符集。
* `description` - 数据库的描述。

## Import

RDS 数据库可以使用实例 ID 和数据库名称进行导入，以冒号 (`:`) 分隔，例如：

```
$ terraform import alibabacloudstack_rds_database.example rm-uf6wjk5****:tftestdatabase
```
