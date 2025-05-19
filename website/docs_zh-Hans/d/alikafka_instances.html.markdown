---
subcategory: "AliKafka"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_alikafka_instances"
sidebar_current: "docs-alibabacloudstack-datasource-alikafka-instances"
description: |-
    查询Alikafka实例资源
---

# alibabacloudstack_alikakfa_instances

根据指定过滤条件列出当前凭证权限可以访问的alikafka实例列表。

## 示例用法

```terraform
variable "name" {
  default = "tf-testacc-alikafkainstance18734"
}



data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
  enable_details              = true
}



resource "alibabacloudstack_alikafka_instance" "default" {
  name      = var.name
  zone_id   = data.alibabacloudstack_zones.default.zones.0.id
  sasl      = true
  plaintext = true
  spec      = "Broker4C16G"
}



data "alibabacloudstack_alikafka_instances" "default" {
  enable_details = "true"
  name_regex     = alibabacloudstack_alikafka_instance.default.name
}
```

## 参数说明

以下参数是支持的：

* `ids` - (可选) 用于过滤结果的ID 列表。
* `name_regex` - (可选) 用于按名称筛选实例名称的正则表达式字符串。

## 属性说明

除了上述参数外，还导出以下属性：

* `ids` - 匹配条件的实例ID列表。
* `names` - 匹配条件的实例名称列表。
* `instances` - 匹配条件的实例列表。每个元素包含以下属性：
  * `id` - ID of the instance.
  * `name` - 实例名称，长度需在3-64字符之间。若未设置，将默认使用实例ID作为名称。
  * `zone_id` - 实例所属可用区ID。注意：当可用区资源不足时，可能会部署到其他可用区。
  * `selected_zones` - 实例部署的目标可用区列表。
  * `cpu_type` - 资源CPU类型，有效值：`intel`。
  * `spec_type` - 实例规格类型。
  * `replicas` - Broker节点数量。
  * `disk_num` - 每个Broker的磁盘数量。
  * `vpc_id` - VPC ID.
  * `vswitch_id` - Vswtich ID.
  * `sasl` - 启用SASL访问点类型。
  * `plaintext` - 启用PLAINTEXT访问点类型。PLAINTEXT类型无需认证即可访问，请注意安全防护。
  * `message_max_bytes` - server可以接收的消息最大尺寸。重要的是，consumer和producer有关这个属性的设置必须同步，否则producer发布的消息对consumer来说太大。消息大小限制不是设置得越大越好,需要视具体业务系统情况而定。
  * `num_partitions` - 如果创建topic时没有给出划分partitions个数，这个数字将是topic下partitions数目的默认数值。
  * `auto_create_topics_enable` - 是否允许自动创建topic。如果是true，则发送消息或者fetch不存在的topic时，会自动创建这个topic。否则需要使用命令行创建topic。
  * `num_io_threads` - server用来处理请求的I/O线程的数目；这个线程数目至少要等于硬盘的个数。
  * `queued_max_requests` - 在网络线程停止读取新请求之前，可以排队等待I/O线程处理的最大请求个数。
  * `replica_fetch_wait_max_ms` - replicas同leader之间通信的最大等待时间，如果失败则会重试
  * `replica_lag_time_max_ms` - 如果一个follower在这个时间内没有发送fetch请求，leader将从ISR中移除这个follower，并认为这个follower已经下线。
  * `num_network_threads` - 服务用来处理网络请求的网络线程数目。
  * `log_retention_bytes` - 每个topic下每个partition保存的字节数。注意，这是每个partition的上限，因此这个数值乘以partition的个数就是每个topic保存的数据总量。注意，如果log.retention.hours和log.retention.bytes都设置了，则超过了任何一个限制都会造成删除一个段文件。注意，这项设置可以由每个topic设置时进行覆盖。
  * `replica_fetch_max_bytes` - Follower从leader备份时每次fetch的最大字节数。
  * `num_replica_fetchers` - 从leader备份数据的线程数 
  * `default_replication_factor` - 默认备份份数，仅指自动创建的topics。
  * `offsets_retention_minutes` - 存在时间超过这个时间限制的offsets都将被标记为待删除。
  * `background_threads` - 用于后台处理的线程数目，例如文件删除。