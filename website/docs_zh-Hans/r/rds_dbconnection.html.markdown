---
subcategory: "RDS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_db_connection"
sidebar_current: "docs-alibabacloudstack-resource-db-connection"
description: |-
  编排RDS互联网连接字符串
---

# alibabacloudstack_db_connection

使用Provider配置的凭证在指定的资源集编排RDS互联网连接字符串。

-> **注意：** 每个RDS实例都会自动分配一个内网连接字符串，其前缀是RDS实例ID。为了避免不必要的冲突，请在应用资源之前指定一个互联网连接前缀。

## 示例用法

```
variable "creation" {
  default = "Rds"
}

variable "name" {
  default = "dbconnectionbasic"
}

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

resource "alibabacloudstack_db_instance" "instance" {
  engine           = "MySQL"
  engine_version   = "5.6"
  instance_type    = "rds.mysql.t1.small"
  instance_storage = "10"
  vswitch_id       = "${alibabacloudstack_vswitch.default.id}"
  instance_name    = "${var.name}"
}
 
resource "alibabacloudstack_db_connection" "foo" {
  instance_id       = "${alibabacloudstack_db_instance.instance.id}"
  connection_prefix = "testabc"
}
```

## 参数参考

以下参数被支持：

* `instance_id` - (必填，变更时重建) 可以运行数据库的实例Id。
* `connection_prefix` - (变更时重建) 互联网连接字符串的前缀。它必须检查唯一性。它可以由小写字母、数字和下划线组成，并且必须以字母开头，长度不超过30个字符。默认为<instance_id> + 'tf'。
* `port` - (可选) 互联网连接端口。有效值范围：[3001-3999]。默认为3306。
* `ip_address` - (可选) 连接字符串的IP地址。

## 属性参考

以下属性被导出：

* `id` - 当前实例连接资源ID。由实例ID和连接字符串组成，格式为 `<instance_id>:<connection_prefix>`。
* `connection_prefix` - 连接字符串的前缀。
* `port` - 连接实例端口。
* `connection_string` - 连接实例字符串。
* `ip_address` - 连接字符串的IP地址。