---
subcategory: "RDS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_rds_database"
sidebar_current: "docs-Alibabacloudstack-rds-database"
description: |- 
  Provides a rds Database resource.
---

# alibabacloudstack_rds_database
-> **NOTE:** Alias name has: `alibabacloudstack_db_database`

Provides a rds Database resource.

## Example Usage

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

## Argument Reference

The following arguments are supported:

* `data_base_instance_id` - (Required, ForceNew) The ID of the RDS instance where the database will be created. Changing this parameter will force the creation of a new resource.
* `data_base_name` - (Required, ForceNew) The name of the database. It must start with a letter and can consist of lowercase letters, numbers, underscores, and hyphens. The length should be 2 to 64 characters. The database name must be unique within the instance.
* `character_set_name` - (Required, ForceNew) The character set of the database. Supported values depend on the database engine:
  - **MySQL/MariaDB**: `utf8`, `gbk`, `latin1`, `utf8mb4`.
  - **SQLServer**: `Chinese_PRC_CI_AS`, `Chinese_PRC_CS_AS`, `SQL_Latin1_General_CP1_CI_AS`, `SQL_Latin1_General_CP1_CS_AS`, `Chinese_PRC_BIN`.
  - **PostgreSQL**: The character set is configured in the format `CharacterSet,<Collate>,<Ctype>`, e.g., `UTF8,C,en_US.utf8`. Character set values include: `KOI8U`, `UTF8`, `WIN866`, `WIN874`, `WIN1250`, `WIN1251`, `WIN1252`, `WIN1253`, `WIN1254`, `WIN1255`, `WIN1256`, `WIN1257`, `WIN1258`, `EUC_CN`, `EUC_KR`, `EUC_TW`, `EUC_JP`, `EUC_JIS_2004`, `KOI8R`, `MULE_INTERNAL`, `LATIN1`, `LATIN2`, `LATIN3`, `LATIN4`, `LATIN5`, `LATIN6`, `LATIN7`, `LATIN8`, `LATIN9`, `LATIN10`, `ISO_8859_5`, `ISO_8859_6`, `ISO_8859_7`, `ISO_8859_8`, `SQL_ASCII`.
* `data_base_description` - (Optional) A description of the database. It must be 2 to 256 characters in length. It must start with a Chinese character or an English letter and can include Chinese and English characters, underscores (`_`), and hyphens (`-`). It cannot start with `http://` or `https://`.
* `instance_id` - (**Deprecated**, ForceNew) This field is deprecated. Use `data_base_instance_id` instead. The ID of the RDS instance where the database will be created.
* `name` - (**Deprecated**, ForceNew) This field is deprecated. Use `data_base_name` instead. The name of the database.
* `character_set` - (**Deprecated**, ForceNew) This field is deprecated. Use `character_set_name` instead. The character set of the database.
* `description` - (**Deprecated**) This field is deprecated. Use `data_base_description` instead. A description of the database.

-> **NOTE:** The `data_base_instance_id`, `data_base_name`, and `character_set_name` fields are required and do not support modification after creation. Modifying these parameters will force the creation of a new resource.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique identifier of the database resource. It is composed of the instance ID and the database name in the format `<instance_id>:<db_name>`.
* `data_base_instance_id` - The ID of the RDS instance where the database is created.
* `data_base_name` - The name of the database.
* `character_set_name` - The character set of the database.
* `data_base_description` - The description of the database.
* `instance_id` - The ID of the RDS instance where the database is created.
* `name` - The name of the database.
* `character_set` - The character set of the database.
* `description` - The description of the database.

## Import

RDS Database can be imported using the instance ID and database name, separated by a colon (`:`), e.g.

```
$ terraform import alibabacloudstack_rds_database.example rm-uf6wjk5****:tftestdatabase
```
