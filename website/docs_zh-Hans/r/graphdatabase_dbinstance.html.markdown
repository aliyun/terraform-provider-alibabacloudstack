---
subcategory: "GraphDatabase"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_graphdatabase_dbinstance"
sidebar_current: "docs-Alibabacloudstack-graphdatabase-dbinstance"
description: |- 
  编排图数据库实例
---

# alibabacloudstack_graphdatabase_dbinstance
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_graph_database_db_instance`

使用Provider配置的凭证在指定的资源集下编排图数据库实例。

## 示例用法

```hcl
variable "name" {
    default = "tf-testaccgraph_databasedb_instance98846"
}

resource "alibabacloudstack_graphdatabase_dbinstance" "default" {
  db_instance_network_type = "VPC"
  payment_type             = "PayAsYouGo"
  vswitch_id              = "vsw-bp152wgftimgq80eiii6k"
  zone_id                 = "cn-hangzhou-h"
  db_instance_storage_type = "cloud_ssd"
  db_instance_description  = "ssd测试"
  db_node_class           = "gdb.r.2xlarge"
  db_instance_category    = "ha"
  vpc_id                  = "vpc-bp1bvsykm9f9hkfeikfi5"
  db_version              = "1.0"
  region_id               = "cn-hangzhou"
  db_node_storage         = 100
  db_instance_ip_array = [
    {
      db_instance_ip_array_attribute = "hidden"
      db_instance_ip_array_name      = "default_whitelist"
      security_ips                   = ["0.0.0.0/0", "192.168.1.1"]
    }
  ]
}
```

## 参数说明

支持以下参数：

* `db_instance_ip_array` - (选填) - 实例组的IP白名单配置列表。列表中的每个元素是一个具有以下结构的对象：
  * `db_instance_ip_array_attribute` - (可选) IP白名单组的属性，默认值为 `hidden`。
  * `db_instance_ip_array_name` - (可选) IP白名单组的名称。
  * `security_ips` - (必填) IP白名单中的IP地址或CIDR块列表，以逗号分隔。最大条目数为1000。例如：`0.0.0.0/0`、`192.168.1.1` 或 `192.168.1.0/24`。

* `db_instance_category` - (必填, 变更时重建) - DB实例的类别。有效值：`ha`(高可用性)。
* `db_instance_description` - (选填) - 根据实际示例或注释的DB实例描述。
* `db_instance_network_type` - (必填, 变更时重建) - DB实例的网络类型。有效值：`VPC`。
* `db_instance_storage_type` - (必填, 变更时重建) - DB实例的存储类型。有效值：`cloud_ssd`、`cloud_essd`。
* `db_node_class` - (必填) - DB节点的类别。有效值包括：`gdb.r.xlarge`、`gdb.r.2xlarge`、`gdb.r.4xlarge`、`gdb.r.8xlarge`、`gdb.r.16xlarge`。
* `db_node_storage` - (必填) - DB实例的存储空间，单位为GB。
* `db_version` - (必填, 变更时重建) - DB实例的内核版本。有效值：`1.0`(Gremlin)、`1.0-OpenCypher`(OpenCypher)。
* `payment_type` - (选填, 变更时重建) - 资源的付费类型。有效值：`PayAsYouGo`。
* `vswitch_id` - (选填, 变更时重建) - 创建DB实例所在的交换机ID。
* `vpc_id` - (选填, 变更时重建) - 创建DB实例所在的专有网络ID。
* `zone_id` - (选填, 变更时重建) - 创建DB实例所在的可用区ID。

## 属性说明

除了上述所有参数外，还导出了以下属性：

* `id` - 资源的唯一标识符。
* `status` - DB实例的状态。可能的值包括：`Creating`（创建中）、`Running`（运行中）、`Deleting`（删除中）、`Rebooting`（重启中）、`DBInstanceClassChanging`（规格变更中）、`NetAddressCreating`（网络地址创建中）和 `NetAddressDeleting`（网络地址删除中）。
* `vswitch_id` - 与DB实例关联的交换机ID。
* `vpc_id` - 与DB实例关联的专有网络ID。
* `zone_id` - DB实例所在的可用区ID。
* `db_instance_ip_array` - 实例的IP白名单配置信息，包含白名单组的属性、名称及具体的IP地址或CIDR块列表。