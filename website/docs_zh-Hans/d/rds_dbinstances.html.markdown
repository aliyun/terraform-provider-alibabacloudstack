---
subcategory: "RDS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_rds_dbinstances"
sidebar_current: "docs-Alibabacloudstack-datasource-rds-dbinstances"
description: |- 
  查询rds数据库实例
---

# alibabacloudstack_rds_dbinstances
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_db_instances`

根据指定过滤条件列出当前凭证权限可以访问的rds数据库实例列表。

## 示例用法

```hcl
data "alibabacloudstack_zones" "default" {
  available_resource_creation = "${var.creation}"
}

resource "alibabacloudstack_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vswitch" "default" {
  vpc_id            = "${alibabacloudstack_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
  name              = "${var.name}"
}

variable "name" {
  default = "tf-testAccDBInstanceConfig"
}

variable "creation" {
  default = "Rds"
}

resource "alibabacloudstack_db_instance" "default" {
  engine               = "MySQL"
  engine_version       = "5.6"
  instance_type        = "rds.mysql.s2.large"
  instance_storage     = "20"
  instance_name        = "${var.name}"
  vswitch_id          = "${alibabacloudstack_vswitch.default.id}"
  storage_type         = "local_ssd"
}

data "alibabacloudstack_rds_dbinstances" "db_instances_ds" {
  name_regex = "${alibabacloudstack_db_instance.default.instance_name}"
  ids        = ["${alibabacloudstack_db_instance.default.id}"]
  status     = "Running"
  engine     = "MySQL"
  tags       = {
    "type" = "database",
    "size" = "tiny"
  }
}

output "first_db_instance_id" {
  value = "${data.alibabacloudstack_rds_dbinstances.db_instances_ds.instances.0.id}"
}
```

## 参数参考

以下参数是支持的：

* `name_regex` - (可选) 用于按实例名称过滤结果的正则表达式字符串。
* `ids` - (可选) RDS 实例 ID 列表。
* `engine` - (可选) 数据库类型。有效值包括：`MySQL`、`PostgreSQL`、`SQLServer` 和 `MariaDB`。如果没有指定值，则返回所有类型。
* `status` - (可选) 资源的状态。例如：`Running`、`Stopped` 等。
* `db_type` - (可选) 数据库实例的类型。有效值：
  * `Primary`: 主实例。
  * `Readonly`: 只读实例。
  * `Guard`: 灾备实例。
  * `Temp`: 临时实例。
* `vpc_id` - (可选) 实例所属的 VPC 的 ID。
* `vswitch_id` - (可选) 实例所属的交换机的 ID。
* `connection_mode` - (可选) 实例的访问模式。有效值：
  * `Standard`: 标准访问模式。
  * `Safe`: 高安全访问模式(数据库代理模式)。
  > **注意**: SQL Server 2012、2016 和 2017 仅支持标准访问模式。
* `tags` - (可选) 分配给数据库实例的标签映射。

## 属性参考

除了上述参数外，还导出以下属性：

* `ids` - RDS 实例 ID 列表。
* `names` - RDS 实例名称列表。
* `instances` - RDS 实例列表。每个元素包含以下属性：
  * `id` - RDS 实例的 ID。
  * `uid` - `id` 的别名。
  * `name` - RDS 实例的名称。
  * `db_type` - 数据库实例的类型。有效值：`Primary`、`Readonly`、`Guard` 和 `Temp`。
  * `charge_type` - 资源的计费类型。有效值：`PrePaid` 或 `PostPaid`。
  * `region_id` - 实例所属的区域 ID。
  * `create_time` - 实例的创建时间。
  * `expire_time` - 实例的过期时间。格式为：`yyyy-MM-ddTHH:mm:ssZ`(UTC 时间)。按量付费实例没有过期时间。
  * `status` - 实例的状态。
  * `engine` - 数据库类型。有效值：`MySQL`、`PostgreSQL`、`SQLServer` 和 `MariaDB`。
  * `engine_version` - 数据库版本。
  * `net_type` - 网络类型。有效值：`Internet` 或 `Intranet`。
  * `connection_mode` - 实例的访问模式。有效值：`Standard` 或 `Safe`。
  * `instance_type` - RDS 实例的规格。
  * `availability_zone` - 可用区。
  * `master_instance_id` - 主实例的 ID。如果未返回此参数，则当前实例为主实例。
  * `guard_instance_id` - 如果当前实例附加了灾备实例，则适用灾备实例的 ID。
  * `temp_instance_id` - 如果当前实例附加了临时实例，则适用临时实例的 ID。
  * `readonly_instance_ids` - 附加到主实例的只读实例 ID 列表。
  * `vpc_id` - 实例所属的 VPC 的 ID。
  * `vswitch_id` - 实例所属的交换机的 ID。
  * `port` - RDS 数据库的连接端口。
  * `connection_string` - RDS 数据库的连接字符串。
  * `instance_storage` - 用户定义的 RDS 实例存储空间。