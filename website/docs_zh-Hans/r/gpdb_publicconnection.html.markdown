---
subcategory: "GraphDatabase(GPDB)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_gpdb_publicconnection"
sidebar_current: "docs-Alibabacloudstack-gpdb-publicconnection"
description: |- 
  编排图数据库连接
---

# alibabacloudstack_gpdb_publicconnection
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_gpdb_connection`

使用Provider配置的凭证在指定的资源集下编排图数据库连接。

## 示例用法

```hcl
data "alibabacloudstack_zones" "default" {
  available_resource_creation = "Gpdb"
}

variable "name" {
  default = "tf-testAccGpdbInstance"
}

resource "alibabacloudstack_vpc" "default" {
  name       = "testing"
  cidr_block = "10.0.0.0/8"
}

resource "alibabacloudstack_vswitch" "default" {
  vpc_id            = alibabacloudstack_vpc.default.id
  cidr_block        = "10.1.0.0/16"
  name              = "apsara_vswitch"
  availability_zone = data.alibabacloudstack_zones.default.zones[0].id
}

resource "alibabacloudstack_gpdb_instance" "default" {
  vswitch_id           = alibabacloudstack_vswitch.default.id
  engine               = "gpdb"
  engine_version       = "4.3"
  instance_class       = "gpdb.group.segsdx2"
  instance_group_count = "2"
  description          = var.name
}

resource "alibabacloudstack_gpdb_connection" "default" {
  instance_id       = alibabacloudstack_gpdb_instance.default.id
  connection_prefix = "tf-testacc10623"
  port              = 3306
}
```

## 参数参考

支持以下参数：

* `instance_id` - (必填，变更时重建) ：要为其创建公网连接的GPDB实例的ID。
* `connection_prefix` - (选填，变更时重建) ：公网连接字符串的前缀。它必须以字母开头，并且只能包含小写字母、数字和下划线(`_`)。长度不能超过30个字符。如果不指定，默认为`<instance_id>-tf`。
* `port` - (选填)：公网连接的端口号。有效值范围从`3200`到`3999`。默认值是`3306`。

### 超时设置

`timeouts`块允许您为某些操作指定[超时](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts)：

* `create` - (默认`10分钟`)用于创建公网连接。
* `update` - (默认`10分钟`)用于更新公网连接。
* `delete` - (默认`10分钟`)用于删除公网连接。

## 属性参考

除了上述所有参数外，还导出了以下属性：

* `id` - GPDB公网连接资源的唯一标识符。它由实例ID和连接前缀组成，格式为`<instance_id>:<connection_prefix>`。
* `connection_string` - 通过公网访问GPDB实例的完整连接字符串。
* `ip_address` - 与公网连接字符串关联的公网IP地址。