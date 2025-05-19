---
subcategory: "AliKafka"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_alikafka_instance"
sidebar_current: "docs-Alibabacloudstack-resource-alikafka-instance"
description: |-
  编排Alikafka实例资源。 
---

# alibabacloudstack_alikafka_instance

提供阿里云消息队列Kafka版实例资源。

## Example Usage

基础用法


```terraform
data "alibabacloudstack_zones" default {
  available_resource_creation = "VSwitch"
  enable_details = true
}



resource "alibabacloudstack_alikafka_instance" "default" {
  sasl = "true"
  plaintext = "true"
  spec = "Broker4C16G"
  cup_type = "Intel"
  name = "tf-testacc-alikafkainstancebasic14494"
  zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
}
```

## 参数说明

以下参数是支持的：

* `name` - (可选， 变更时重建) 实例名称，长度需在3-64字符之间。若未设置，将默认使用实例ID作为名称。
* `zone_id` - (可选， 变更时重建) 实例所属可用区ID。注意：当可用区资源不足时，可能会部署到其他可用区。
* `selected_zones` - (可选， 变更时重建， 列表) 实例部署的目标可用区列表。
* `cpu_type` - (必填， 变更时重建) 资源CPU类型，有效值：`intel`。
* `spec_type` - (可选， 变更时重建) 实例规格类型。
* `replicas` - (可选， 变更时重建) Broker节点数量。
* `disk_num` - (可选， 变更时重建) 每个Broker的磁盘数量。
* `vpc_id` - (可选， 变更时重建) VPC ID.
* `vswitch_id` - (可选， 变更时重建) Vswtich ID.
* `sasl` - (可选， 变更时重建) 启用SASL访问点类型。
* `plaintext` - (可选， 变更时重建) 启用PLAINTEXT访问点类型。PLAINTEXT类型无需认证即可访问，请注意安全防护。
* `message_max_bytes` - (可选) server可以接收的消息最大尺寸。重要的是，consumer和producer有关这个属性的设置必须同步，否则producer发布的消息对consumer来说太大。消息大小限制不是设置得越大越好,需要视具体业务系统情况而定。
* `num_partitions` - (可选) 如果创建topic时没有给出划分partitions个数，这个数字将是topic下partitions数目的默认数值。
* `auto_create_topics_enable` - (可选) 是否允许自动创建topic。如果是true，则发送消息或者fetch不存在的topic时，会自动创建这个topic。否则需要使用命令行创建topic。
* `num_io_threads` - (可选) server用来处理请求的I/O线程的数目；这个线程数目至少要等于硬盘的个数。
* `queued_max_requests` - (可选) 在网络线程停止读取新请求之前，可以排队等待I/O线程处理的最大请求个数。
* `replica_fetch_wait_max_ms` - (可选) replicas同leader之间通信的最大等待时间，如果失败则会重试
* `replica_lag_time_max_ms` - (可选) 如果一个follower在这个时间内没有发送fetch请求，leader将从ISR中移除这个follower，并认为这个follower已经下线。
* `num_network_threads` - (可选) 服务用来处理网络请求的网络线程数目。
* `log_retention_bytes` - (可选) 每个topic下每个partition保存的字节数。注意，这是每个partition的上限，因此这个数值乘以partition的个数就是每个topic保存的数据总量。注意，如果log.retention.hours和log.retention.bytes都设置了，则超过了任何一个限制都会造成删除一个段文件。注意，这项设置可以由每个topic设置时进行覆盖。
* `replica_fetch_max_bytes` - (可选) Follower从leader备份时每次fetch的最大字节数。
* `num_replica_fetchers` - (可选) 从leader备份数据的线程数 
* `default_replication_factor` - (可选) 默认备份份数，仅指自动创建的topics。
* `offsets_retention_minutes` - (可选) 存在时间超过这个时间限制的offsets都将被标记为待删除。
* `background_threads` - (可选) 用于后台处理的线程数目，例如文件删除。

## 属性说明

除了上述参数外，还导出以下属性：

* `id` - 实例在Terraform中的资源ID。
* `sasl_plaintext_endpoint` - 实例的SASL明文访问端点（域名模式）。
* `sasl_ssl_endpoint` - 实例的SASL_SSL访问端点（域名模式）。
* `status` - 实例状态。

## Timeouts

* `create` - （默认60分钟）创建资源时使用。
* `update` - （默认120分钟）更新资源时使用。
* `delete` - （默认30分钟）删除资源时使用。

## Import

可通过ID导入AliKafka实例，例如：

```shell
$ terraform import alibabacloudstack_alikafka_instance.instance <id>
```
