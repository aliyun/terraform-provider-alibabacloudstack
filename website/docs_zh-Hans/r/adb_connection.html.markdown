---
subcategory: "ADB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_adb_connection"
sidebar_current: "docs-Alibabacloudstack-adb-connection"
description: |- 
  编排adb连接信息。
---

# alibabacloudstack_adb_connection

使用Provider配置的凭证在指定的资源集下编排adb连接信息。

## 示例用法

```hcl
variable "name" {
    default = "tf-testaccadbconnection96904"
}

resource "alibabacloudstack_vpc" "default" {
  vpc_name       = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vswitch" "default" {
  vpc_id            = alibabacloudstack_vpc.default.id
  cidr_block        = "172.16.0.0/24"
  zone_id           = data.alibabacloudstack_zones.default.zones[0].id
  vswitch_name      = var.name
}

resource "alibabacloudstack_adb_db_cluster" "cluster" {
  db_cluster_version  = "3.0"
  db_cluster_category = "Basic"
  db_node_class       = "C8"
  db_node_count       = 2
  db_node_storage     = 200
  mode                = "reserver"
  vswitch_id          = alibabacloudstack_vswitch.default.id
  description         = var.name
  cluster_type        = "analyticdb"
  cpu_type            = "intel"
}

resource "alibabacloudstack_adb_connection" "default" {
  db_cluster_id         = alibabacloudstack_adb_db_cluster.cluster.id
  connection_prefix     = "testabc"
}
```

## 参数说明

支持以下参数：
  * `db_cluster_id` - (必填, 变更时重建) ADB集群的ID。此字段是不可变的，创建后无法更改。
  * `connection_prefix` - (选填, 变更时重建) 集群公共端点的前缀。该前缀必须为6到30个字符长度，并且可以包含小写字母、数字和连字符(-)。它必须以字母开头并以数字或字母结尾。如果不指定，默认为`<db_cluster_id> + tf`。
  * `port` - (选填, 计算后返回) 用于连接到ADB集群的端口号。

## 属性说明

除了上述所有参数外，还导出了以下属性：
  * `id` - ADB连接资源的唯一标识符。它由集群ID和连接字符串组成，格式为`<db_cluster_id>:<connection_prefix>`。
  * `connection_prefix` - 用于连接字符串的前缀。
  * `port` - 用于连接到ADB集群的端口号。
  * `connection_string` - 用于访问ADB集群的完整连接字符串。
  * `ip_address` - 与连接字符串关联的IP地址。

## 导入

ADB连接可以通过id导入，例如：

```bash
$ terraform import alibabacloudstack_adb_connection.example am-12345678
```