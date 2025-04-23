---
subcategory: "Redis And Memcache (KVStore)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_redis_tairinstances"
sidebar_current: "docs-Alibabacloudstack-datasource-redis-tairinstances"
description: |- 
  查询redis tair实例
---

# alibabacloudstack_redis_tairinstances
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_kvstore_instances`

根据指定过滤条件列出当前凭证权限可以访问的redis数据库tair实例列表。

## 示例用法

```hcl
variable "name" {
    default = "tf-testAccCheckAlibabacloudStackRKVInstancesDataSource46059"
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

resource "alibabacloudstack_kvstore_instance" "default" {
  zone_id = data.alibabacloudstack_zones.default.zones[0].id
  instance_class = "redis.master.small.default"
  instance_name  = var.name
  vswitch_id     = alibabacloudstack_vpc_vswitch.default.id
  security_ips   = ["10.0.0.1"]
  instance_type  = "Redis"
  engine_version = "4.0"
}

data "alibabacloudstack_redis_tairinstances" "default" {
    name_regex = "checkalibabacloudstacktairinstancesdatasource"
    status      = "Running"
    instance_type = "tair_rdb"
    ids         = [alibabacloudstack_kvstore_instance.default.id]
}

output "first_instance_name" {
    value = data.alibabacloudstack_redis_tairinstances.default.instances.0.name
}
```

## 参数参考

以下参数是支持的：

* `name_regex` - (可选) 应用到 Tair 实例名称的正则表达式字符串，用于筛选符合条件的实例名称。
* `ids` - (可选) Tair 实例 ID 列表，用于通过实例 ID 筛选特定实例。
* `status` - (可选) 资源的状态。有效值包括：`Creating`(创建中)、`Running`(运行中)、`Restarting`(重启中)、`ChangingConfig`(配置变更中)、`FlushingData`(数据刷盘中)、`Deleting`(删除中)、`NetworkChanging`(网络变更中)、`Abnormal`(异常)。
* `instance_type` - (可选) 实例的存储介质类型。有效值包括：`tair_rdb`(内存型)、`tair_scm`(持久内存型)、`tair_essd`(云盘型)。

## 属性参考

除了上述参数外，还导出以下属性：

* `names` - 匹配条件的 Tair 实例名称列表。
* `ids` - 匹配条件的 Tair 实例 ID 列表。
* `instances` - 匹配条件的 Tair 实例详细信息列表。每个元素包含以下属性：
  * `id` - Tair 实例的 ID。
  * `name` - Tair 实例的名称。
  * `charge_type` - 计费方式。选项值：`PostPaid` 表示按量付费，`PrePaid` 表示包年包月。
  * `region_id` - 实例所属的区域 ID。
  * `create_time` - 实例的创建时间，遵循 ISO 8601 标准，格式为 `yyyy-MM-ddTHH:mm:ssZ`，时间为 UTC 时间。
  * `expire_time` - 实例的过期时间，按量付费实例没有过期时间。
  * `status` - 资源的状态。有效值包括：`Creating`、`Running`、`Restarting`、`ChangingConfig`、`FlushingData`、`Deleting`、`NetworkChanging`、`Abnormal`。
  * `instance_type` - 实例的存储介质类型。有效值包括：`tair_rdb`(内存型)、`tair_scm`(持久内存型)、`tair_essd`(云盘型)。
  * `instance_class` - 实例的规格类型，更多信息请参见 [实例规格表](https://help.aliyun.com/document_detail/26350.htm)。
  * `availability_zone` - 实例所在的可用区。
  * `vpc_id` - 虚拟专用网络 (VPC) 的 ID。
  * `vswitch_id` - VSwitch 的 ID。
  * `private_ip` - 实例的私有 IP 地址。
  * `port` - Tair 服务端口，取值范围为 1024 到 65535，默认值为 6379。
  * `user_name` - 实例的用户名。
  * `capacity` - 实例的存储容量，单位为 MB。
  * `bandwidth` - 实例的带宽，单位为 Mbit/s。
  * `connections` - 实例的连接数量限制，单位为个数。
  * `connection_domain` - 实例的内网连接地址。