---
subcategory: "RDS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_db_read_write_splitting_connection"
sidebar_current: "docs-alibabacloudstack-resource-db-read-write-splitting-connection"
description: |-
  编排RDS读写分离连接
---

# alibabacloudstack_db_read_write_splitting_connection

使用Provider配置的凭证在指定的资源集下编排RDS读写分离连接，为RDS实例分配一个内网连接字符串。

## 示例用法

```
variable "creation" {
  default = "RDS"
}

variable "name" {
  default = "dbInstancevpc"
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

resource "alibabacloudstack_db_instance" "default" {
  engine               = "MySQL"
  engine_version       = "5.6"
  instance_type        = "rds.mysql.t1.small"
  instance_storage     = "20"
  instance_charge_type = "Postpaid"
  instance_name        = "${var.name}"
  vswitch_id           = "${alibabacloudstack_vswitch.default.id}"
  security_ips         = ["10.168.1.12", "100.69.7.112"]
}

resource "alibabacloudstack_db_readonly_instance" "default" {
  master_db_instance_id = "${alibabacloudstack_db_instance.default.id}"
  zone_id               = "${alibabacloudstack_db_instance.default.zone_id}"
  engine_version        = "${alibabacloudstack_db_instance.default.engine_version}"
  instance_type         = "${alibabacloudstack_db_instance.default.instance_type}"
  instance_storage      = "30"
  instance_name         = "${var.name}ro"
  vswitch_id            = "${alibabacloudstack_vswitch.default.id}"
}

resource "alibabacloudstack_db_read_write_splitting_connection" "default" {
  instance_id       = "${alibabacloudstack_db_instance.default.id}"
  connection_prefix = "t-con-123"
  distribution_type = "Standard"
  depends_on = ["alibabacloudstack_db_readonly_instance.default"]
}
```

-> **注意:** 资源 `alibabacloudstack_db_read_write_splitting_connection` 应该在 `alibabacloudstack_db_readonly_instance` 创建之后创建，因此需要使用 `depends_on` 语句。

## 参数说明

以下是支持的参数：

* `instance_id` - (必填，变更时重建) 可以运行数据库的实例ID。
* `distribution_type` - (必填) 读权重分配模式。取值如下：`Standard` 表示基于类型的自动权重分配，`Custom` 表示自定义权重分配。
* `connection_prefix` - (可选，变更时重建) 互联网连接字符串前缀。它必须检查唯一性。它可以由小写字母、数字和下划线组成，并且必须以字母开头，长度不超过30个字符。默认为 `<instance_id>` + 'rw'。
* `port` - (可选) 内网连接端口。有效值范围：[3001-3999]。默认为3306。
* `max_delay_time` - (可选) 延迟阈值，以秒为单位。取值范围是0到7200，默认为30。读请求不会被路由到延迟大于此阈值的只读实例。
* `weight` - (可选) 读权重分配。读权重以100为步长递增，最大为10,000。格式如下：{"Instanceid":"Weight","Instanceid":"Weight"}。当 `distribution_type` 设置为 Custom 时，此参数必须设置。
* `connection_string` - (可选) 连接实例字符串。

## 属性说明

以下属性将被导出：

* `id` - DB实例的ID。
* `connection_string` - 连接实例字符串。
* `port` - 内网连接端口。有效值范围：[3001-3999]。默认为3306。