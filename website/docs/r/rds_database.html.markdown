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
  instance_id      = "${alibabacloudstack_db_instance.instance.id}"
  name             = "tftestdatabase"
  character_set    = "utf8"
  description      = "This is a test database"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) The ID of the RDS instance where the database will be created.
* `name` - (Required, ForceNew) The name of the database. It must start with a letter and can consist of lowercase letters, numbers, and underscores. The length should not exceed 64 characters.
* `character_set` - (Required, ForceNew) The character set of the database. Supported values depend on the database engine:
  - **MySQL**: `utf8`, `gbk`, `latin1`, `utf8mb4` (`utf8mb4` is only supported for versions 5.5 and 5.6).
  - **SQLServer**: `Chinese_PRC_CI_AS`, `Chinese_PRC_CS_AS`, `SQL_Latin1_General_CP1_CI_AS`, `SQL_Latin1_General_CP1_CS_AS`, `Chinese_PRC_BIN`.
  - **PostgreSQL**: `KOI8U`, `UTF8`, `WIN866`, `WIN874`, `WIN1250`, `WIN1251`, `WIN1252`, `WIN1253`, `WIN1254`, `WIN1255`, `WIN1256`, `WIN1257`, `WIN1258`, `EUC_CN`, `EUC_KR`, `EUC_TW`, `EUC_JP`, `EUC_JIS_2004`, `KOI8R`, `MULE_INTERNAL`, `LATIN1`, `LATIN2`, `LATIN3`, `LATIN4`, `LATIN5`, `LATIN6`, `LATIN7`, `LATIN8`, `LATIN9`, `LATIN10`, `ISO_8859_5`, `ISO_8859_6`, `ISO_8859_7`, `ISO_8859_8`, `SQL_ASCII`.
  Refer to the [API Docs](https://www.alibabacloud.com/help/zh/doc-detail/26258.htm) for more details.
* `description` - (Optional) A description of the database. It cannot start with `https://`. It must start with a Chinese character or an English letter and can include Chinese and English characters, underscores (`_`), hyphens (`-`), and numbers. The length should be between 2 and 256 characters.

**NOTE:** The `name` and `character_set` fields do not support modification after creation.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique identifier of the database resource. It is composed of the instance ID and the database name in the format `<instance_id>:<name>`.
* `instance_id` - The ID of the RDS instance where the database is created.
* `name` - The name of the database.
* `character_set` - The character set of the database.
* `description` - The description of the database.
